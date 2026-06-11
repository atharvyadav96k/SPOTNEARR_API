package sync

import (
	"context"
	"encoding/json"
	"log"

	searchdb "github.com/Developer-Aadesh/spotnearr-database/search"
	searchpostgres "github.com/Developer-Aadesh/spotnearr-database/search/postgres"
	"gorm.io/gorm"
)

// SyncPayload is the JSON body pushed by the vendor service to /internal/sync.
// It mirrors models.OutboxPayload from the vendor service.
type SyncPayload struct {
	EventType    string   `json:"event_type"`
	InvProductID uint     `json:"inv_product_id"`
	ProductID    uint     `json:"product_id"`
	BusinessID   uint     `json:"business_id"`
	ProductName  string   `json:"product_name"`
	Price        float64  `json:"price"`
	PriceUnit    string   `json:"price_unit"`
	Quantity     *float64 `json:"quantity"`
	QuantityUnit *string  `json:"quantity_unit"`
	Desc         string   `json:"desc"`
	SearchTokens []string `json:"search_tokens"`
	Categories   []struct {
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

// Applier applies a vendor-pushed sync payload to the search database.
// It never accesses the vendor database.
type Applier struct {
	repo *searchpostgres.SearchRepository
}

func NewApplier(db *gorm.DB) *Applier {
	return &Applier{repo: searchpostgres.NewSearchRepository(db)}
}

func (a *Applier) Apply(ctx context.Context, data []byte) error {
	var p SyncPayload
	if err := json.Unmarshal(data, &p); err != nil {
		return err
	}
	switch p.EventType {
	case "upsert":
		catIDs := make([]uint, len(p.Categories))
		for i, c := range p.Categories {
			catIDs[i] = c.ID
		}
		return a.repo.Upsert(ctx, searchdb.SearchEntry{
			ID:            p.InvProductID,
			ProductID:     p.ProductID,
			BusinessID:    p.BusinessID,
			ProductName:   p.ProductName,
			Price:         p.Price,
			PriceUnit:     p.PriceUnit,
			Quantity:      p.Quantity,
			QuantityUnit:  p.QuantityUnit,
			Description:   p.Desc,
			SearchTokens:  p.SearchTokens,
			CategoryIDs:   catIDs,
			StoreID:       p.StoreID,
			StoreName:     p.StoreName,
			StreetAddress: p.StreetAddress,
			Lat:           p.Lat,
			Long:          p.Long,
			GeoHash:       p.GeoHash,
			Available:     p.Available,
		})
	case "delete":
		return a.repo.SoftDelete(ctx, p.InvProductID)
	default:
		log.Printf("sync: unknown event_type %q", p.EventType)
		return nil
	}
}
