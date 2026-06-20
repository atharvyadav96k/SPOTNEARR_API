package tsclient

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/typesense/typesense-go/v2/typesense"
	"github.com/typesense/typesense-go/v2/typesense/api"
	"github.com/typesense/typesense-go/v2/typesense/api/pointer"
)

const CollectionName = "products"

func New(host, port, apiKey string) *typesense.Client {
	return typesense.NewClient(
		typesense.WithServer(fmt.Sprintf("http://%s:%s", host, port)),
		typesense.WithAPIKey(apiKey),
		typesense.WithConnectionTimeout(5*time.Second),
		typesense.WithRetryInterval(2*time.Second),
	)
}

// EnsureCollection creates the products collection if it does not already exist.
func EnsureCollection(ctx context.Context, client *typesense.Client) error {
	_, err := client.Collection(CollectionName).Retrieve(ctx)
	if err == nil {
		return nil
	}

	schema := &api.CollectionSchema{
		Name: CollectionName,
		Fields: []api.Field{
			{Name: "id",             Type: "string"},
			{Name: "inv_product_id", Type: "int64"},
			{Name: "product_id",     Type: "int64"},
			{Name: "name",           Type: "string"},
			{Name: "price",          Type: "float"},
			{Name: "price_unit",     Type: "string"},
			{Name: "category_ids",   Type: "int64[]", Optional: pointer.True()},
			{Name: "available",      Type: "bool"},
			{Name: "location",       Type: "geopoint"},
		},
	}

	if _, err := client.Collections().Create(ctx, schema); err != nil {
		return fmt.Errorf("typesense: create collection: %w", err)
	}
	log.Printf("typesense: collection %q created", CollectionName)
	return nil
}
