package sync

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/atharvyadav96k/spotnearr/search-svc/models"
	"github.com/atharvyadav96k/spotnearr/search-svc/repository"
	"gorm.io/gorm"
)

// outboxRow mirrors search_sync_outbox in the vendor DB. We read it here;
// we never write to the vendor DB except to mark rows processed.
type outboxRow struct {
	ID        int64     `gorm:"column:id"`
	EventType string    `gorm:"column:event_type"`
	Payload   []byte    `gorm:"column:payload"`
	CreatedAt time.Time `gorm:"column:created_at"`
}

// outboxPayload mirrors models.OutboxPayload from the vendor service.
type outboxPayload struct {
	EventType     string   `json:"event_type"`
	InvProductID  uint     `json:"inv_product_id"`
	ProductID     uint     `json:"product_id"`
	BusinessID    uint     `json:"business_id"`
	ProductName   string   `json:"product_name"`
	Price         float64  `json:"price"`
	PriceUnit     string   `json:"price_unit"`
	Quantity      *float64 `json:"quantity"`
	QuantityUnit  *string  `json:"quantity_unit"`
	Desc          string   `json:"desc"`
	SearchTokens  []string `json:"search_tokens"`
	Categories    []struct {
		ID uint `json:"id"`
	} `json:"categories"`
	StoreID       uint    `json:"store_id"`
	StoreName     string  `json:"store_name"`
	StreetAddress string  `json:"street_address"`
	Lat           float64 `json:"lat"`
	Long          float64 `json:"long"`
	GeoHash       string  `json:"geo_hash"`
	Available     bool    `json:"available"`
}

// Poller reads unprocessed rows from the vendor DB outbox and applies them to
// search_entries. It runs every pollInterval as a fallback; the Redis subscriber
// triggers it immediately after each vendor write.
type Poller struct {
	vendorDB *gorm.DB
	searchDB *gorm.DB
	repo     *repository.SearchRepository
}

func NewPoller(vendorDB, searchDB *gorm.DB) *Poller {
	return &Poller{
		vendorDB: vendorDB,
		searchDB: searchDB,
		repo:     repository.NewSearchRepository(searchDB),
	}
}

// Run starts the 30-second fallback poll loop. Call in a goroutine.
func (p *Poller) Run(ctx context.Context, pollInterval time.Duration) {
	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			p.ProcessPending(ctx)
		}
	}
}

// ProcessPending fetches and applies all unprocessed outbox rows. Safe to call
// concurrently (rows are claimed with a SELECT FOR UPDATE SKIP LOCKED).
func (p *Poller) ProcessPending(ctx context.Context) {
	var rows []outboxRow
	err := p.vendorDB.WithContext(ctx).Raw(`
		SELECT id, event_type, payload, created_at
		FROM search_sync_outbox
		WHERE processed = false
		ORDER BY id ASC
		LIMIT 100
		FOR UPDATE SKIP LOCKED
	`).Scan(&rows).Error
	if err != nil {
		log.Printf("sync: fetch outbox rows: %v", err)
		return
	}
	if len(rows) == 0 {
		return
	}

	for _, row := range rows {
		if err := p.applyRow(ctx, row); err != nil {
			log.Printf("sync: apply outbox row %d: %v", row.ID, err)
			continue
		}
		// Mark processed — best-effort; if this fails the row will be retried.
		p.vendorDB.WithContext(ctx).Exec(`
			UPDATE search_sync_outbox
			SET processed = true, processed_at = NOW()
			WHERE id = ?
		`, row.ID)
	}
}

func (p *Poller) applyRow(ctx context.Context, row outboxRow) error {
	var payload outboxPayload
	if err := json.Unmarshal(row.Payload, &payload); err != nil {
		return err
	}

	switch payload.EventType {
	case "upsert":
		catIDs := make([]uint, len(payload.Categories))
		for i, c := range payload.Categories {
			catIDs[i] = c.ID
		}
		entry := models.SearchEntry{
			ID:            payload.InvProductID,
			ProductID:     payload.ProductID,
			BusinessID:    payload.BusinessID,
			ProductName:   payload.ProductName,
			Price:         payload.Price,
			PriceUnit:     payload.PriceUnit,
			Quantity:      payload.Quantity,
			QuantityUnit:  payload.QuantityUnit,
			Description:   payload.Desc,
			SearchTokens:  payload.SearchTokens,
			CategoryIDs:   catIDs,
			StoreID:       payload.StoreID,
			StoreName:     payload.StoreName,
			StreetAddress: payload.StreetAddress,
			Lat:           payload.Lat,
			Long:          payload.Long,
			GeoHash:       payload.GeoHash,
			Available:     payload.Available,
		}
		return p.repo.Upsert(ctx, entry)
	case "delete":
		return p.repo.SoftDelete(ctx, payload.InvProductID)
	default:
		log.Printf("sync: unknown event_type %q in row %d", payload.EventType, row.ID)
		return nil
	}
}
