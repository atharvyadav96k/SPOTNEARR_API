package tsclient

import (
	"context"
	"fmt"

	"github.com/typesense/typesense-go/v2/typesense"
)

// ProductDoc is the document shape stored in the Typesense products collection.
type ProductDoc struct {
	ID           string    `json:"id"`
	InvProductID int64     `json:"inv_product_id"`
	ProductID    int64     `json:"product_id"`
	Name         string    `json:"name"`
	Price        float64   `json:"price"`
	PriceUnit    string    `json:"price_unit"`
	CategoryIDs  []int64   `json:"category_ids"`
	Available    bool      `json:"available"`
	Location     []float64 `json:"location"` // [lat, lng]
}

type Indexer struct {
	client *typesense.Client
}

func NewIndexer(client *typesense.Client) *Indexer {
	return &Indexer{client: client}
}

func (idx *Indexer) Upsert(ctx context.Context, doc ProductDoc) error {
	_, err := idx.client.Collection(CollectionName).Documents().Upsert(ctx, doc)
	if err != nil {
		return fmt.Errorf("typesense: upsert %s: %w", doc.ID, err)
	}
	return nil
}

func (idx *Indexer) Delete(ctx context.Context, invProductID uint) error {
	id := fmt.Sprintf("%d", invProductID)
	_, err := idx.client.Collection(CollectionName).Document(id).Delete(ctx)
	if err != nil {
		return fmt.Errorf("typesense: delete %s: %w", id, err)
	}
	return nil
}
