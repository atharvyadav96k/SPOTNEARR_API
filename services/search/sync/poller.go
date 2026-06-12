package sync

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	searchdb "github.com/Developer-Aadesh/spotnearr-database/search"
	searchpostgres "github.com/Developer-Aadesh/spotnearr-database/search/postgres"
	"github.com/atharvyadav96k/spotnearr/pkg/tokenizer"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

const freqDeltaKey = "search:tf_delta"

// SyncPayload is the JSON body pushed by the vendor service to POST /internal/sync.
// Only the fields required by the search service are parsed; extra fields are silently ignored.
type SyncPayload struct {
	EventType    string  `json:"event_type"`
	InvProductID uint    `json:"inv_product_id"`
	ProductID    uint    `json:"product_id"`
	ProductName  string  `json:"product_name"`
	Price        float64 `json:"price"`
	PriceUnit    string  `json:"price_unit"`
	Desc         string  `json:"desc"`
	Categories   []struct {
		ID uint `json:"id"`
	} `json:"categories"`
	Lat       float64 `json:"lat"`
	Long      float64 `json:"long"`
	GeoHash   string  `json:"geo_hash"`
	Available bool    `json:"available"`
}

// Applier applies vendor-pushed sync events to the search index and maintains
// token-category frequency counters in Redis.
type Applier struct {
	repo *searchpostgres.SearchRepository
	rdb  *redis.Client
}

func NewApplier(db *gorm.DB, rdb *redis.Client) *Applier {
	return &Applier{
		repo: searchpostgres.NewSearchRepository(db),
		rdb:  rdb,
	}
}

func (a *Applier) Apply(ctx context.Context, data []byte) error {
	var p SyncPayload
	if err := json.Unmarshal(data, &p); err != nil {
		return err
	}
	switch p.EventType {
	case "upsert":
		return a.applyUpsert(ctx, p)
	case "delete":
		return a.applyDelete(ctx, p.InvProductID)
	default:
		log.Printf("sync: unknown event_type %q", p.EventType)
		return nil
	}
}

func (a *Applier) applyUpsert(ctx context.Context, p SyncPayload) error {
	// Search service owns tokenization — re-derive tokens from name + desc.
	tokens := tokenizer.MergeTokens(p.ProductName, p.Desc)

	catIDs := make([]uint, len(p.Categories))
	for i, c := range p.Categories {
		catIDs[i] = c.ID
	}

	// Read the existing entry (if any) to compute the freq delta.
	old, err := a.repo.GetByID(ctx, p.InvProductID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("upsert: read old entry: %w", err)
	}

	if err := a.repo.Upsert(ctx, searchdb.SearchEntry{
		ID:           p.InvProductID,
		ProductID:    p.ProductID,
		Name:         p.ProductName,
		Price:        p.Price,
		PriceUnit:    p.PriceUnit,
		SearchTokens: tokens,
		CategoryIDs:  catIDs,
		Lat:          p.Lat,
		Long:         p.Long,
		GeoHash:      p.GeoHash,
		Available:    p.Available,
	}); err != nil {
		return fmt.Errorf("upsert: write: %w", err)
	}

	a.adjustFreqs(ctx, old, tokens, catIDs)
	return nil
}

func (a *Applier) applyDelete(ctx context.Context, id uint) error {
	old, err := a.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil // already gone
		}
		return fmt.Errorf("delete: read entry: %w", err)
	}

	if err := a.repo.SoftDelete(ctx, id); err != nil {
		return fmt.Errorf("delete: soft-delete: %w", err)
	}

	a.removeFreqs(ctx, old)
	return nil
}

// adjustFreqs decrements the old entry's token×category pairs and increments the new ones.
func (a *Applier) adjustFreqs(ctx context.Context, old *searchdb.SearchEntry, newTokens []string, newCatIDs []uint) {
	pipe := a.rdb.Pipeline()
	if old != nil {
		for _, tok := range old.SearchTokens {
			for _, catID := range old.CategoryIDs {
				pipe.HIncrBy(ctx, freqDeltaKey, fmt.Sprintf("%s:%d", tok, catID), -1)
			}
		}
	}
	for _, tok := range newTokens {
		for _, catID := range newCatIDs {
			pipe.HIncrBy(ctx, freqDeltaKey, fmt.Sprintf("%s:%d", tok, catID), 1)
		}
	}
	if _, err := pipe.Exec(ctx); err != nil {
		log.Printf("sync: freq adjust: %v", err)
	}
}

func (a *Applier) removeFreqs(ctx context.Context, entry *searchdb.SearchEntry) {
	if entry == nil || len(entry.SearchTokens) == 0 {
		return
	}
	pipe := a.rdb.Pipeline()
	for _, tok := range entry.SearchTokens {
		for _, catID := range entry.CategoryIDs {
			pipe.HIncrBy(ctx, freqDeltaKey, fmt.Sprintf("%s:%d", tok, catID), -1)
		}
	}
	if _, err := pipe.Exec(ctx); err != nil {
		log.Printf("sync: freq remove: %v", err)
	}
}
