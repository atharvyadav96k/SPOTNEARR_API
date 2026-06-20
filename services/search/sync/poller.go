package sync

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/atharvyadav96k/spotnearr/pkg/tokenizer"
	"github.com/atharvyadav96k/spotnearr/search-svc/cache"
	tsclient "github.com/atharvyadav96k/spotnearr/search-svc/typesense"
)

// SyncPayload is the JSON body published by the vendor service on product changes.
type SyncPayload struct {
	EventType    string  `json:"event_type"`
	InvProductID uint    `json:"inv_product_id"`
	ProductID    uint    `json:"product_id"`
	ProductName  string  `json:"product_name"`
	Price        float64 `json:"price"`
	PriceUnit    string  `json:"price_unit"`
	Categories   []struct {
		ID uint `json:"id"`
	} `json:"categories"`
	Lat       float64 `json:"lat"`
	Long      float64 `json:"long"`
	Available bool    `json:"available"`
}

// Applier applies vendor sync events to the Typesense index and maintains
// token→category frequency counts in Redis for category-affinity boosting.
type Applier struct {
	indexer *tsclient.Indexer
	tcCache *cache.TokenCategoryCache
}

func NewApplier(indexer *tsclient.Indexer, tcCache *cache.TokenCategoryCache) *Applier {
	return &Applier{indexer: indexer, tcCache: tcCache}
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
	catIDs := make([]int64, len(p.Categories))
	for i, c := range p.Categories {
		catIDs[i] = int64(c.ID)
	}

	if err := a.indexer.Upsert(ctx, tsclient.ProductDoc{
		ID:           fmt.Sprintf("%d", p.InvProductID),
		InvProductID: int64(p.InvProductID),
		ProductID:    int64(p.ProductID),
		Name:         p.ProductName,
		Price:        p.Price,
		PriceUnit:    p.PriceUnit,
		CategoryIDs:  catIDs,
		Available:    p.Available,
		Location:     []float64{p.Lat, p.Long},
	}); err != nil {
		return err
	}

	// Update token→category frequency counts and store reverse-lookup for future deletes.
	tokens := tokenizer.TokenParser(p.ProductName)
	if err := a.tcCache.IncrFreqs(ctx, tokens, catIDs); err != nil {
		log.Printf("sync: IncrFreqs %d: %v", p.InvProductID, err)
	}
	if err := a.tcCache.SetDocMeta(ctx, p.InvProductID, tokens, catIDs); err != nil {
		log.Printf("sync: SetDocMeta %d: %v", p.InvProductID, err)
	}

	return nil
}

func (a *Applier) applyDelete(ctx context.Context, id uint) error {
	// Retrieve stored metadata before deleting so we can decrement frequencies.
	tokens, catIDs := a.tcCache.GetDocMeta(ctx, id)
	if len(tokens) > 0 {
		if err := a.tcCache.DecrFreqs(ctx, tokens, catIDs); err != nil {
			log.Printf("sync: DecrFreqs %d: %v", id, err)
		}
	}
	a.tcCache.DelDocMeta(ctx, id)

	return a.indexer.Delete(ctx, id)
}
