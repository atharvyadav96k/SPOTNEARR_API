package implementation

import (
	"context"
	"fmt"
	"log"
	"math"
	"strings"

	"github.com/atharvyadav96k/SPOTNEARR_API/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type InvProductRepository struct {
	db *gorm.DB
}

func NewInvProductRepository(db *gorm.DB) *InvProductRepository {
	return &InvProductRepository{db: db}
}

func (i *InvProductRepository) AddProduct(ctx context.Context, invProduct models.InventoryProduct, bizID uint) error {
	log.Default().Println("db start")
	var count int64
	err := i.db.WithContext(ctx).
		Model(&models.Store{}).
		Where(
			"id = ? AND business_id = ? AND deleted_at IS NULL",
			invProduct.StoreID,
			bizID,
		).
		Count(&count).Error

	if err != nil {
		return err
	}

	if count == 0 {
		return gorm.ErrRecordNotFound
	}
	log.Default().Println("db end")
	return i.db.WithContext(ctx).
		Create(&invProduct).Error
}

func (i *InvProductRepository) GetInvProductList(ctx context.Context, invID uint, bizID uint) ([]models.Product, error) {
	var products []models.Product

	err := i.db.WithContext(ctx).
		Model(&models.Product{}).
		Distinct("products.*").
		Joins("JOIN inventory_products ip ON ip.product_id = products.id").
		Joins("JOIN stores s ON s.id = ip.store_id").
		Where("ip.store_id = ?", invID).
		Where("ip.deleted_at IS NULL").
		Where("s.business_id = ?", bizID).
		Find(&products).Error

	if err != nil {
		return nil, err
	}

	return products, nil
}

func (i *InvProductRepository) UpdateProduct(ctx context.Context, invProduct models.InventoryProduct, bizID uint) (models.InventoryProduct, error) {
	updates := map[string]any{}

	if invProduct.Count != nil {
		updates["count"] = *invProduct.Count
	}

	if invProduct.Available != nil {
		updates["available"] = *invProduct.Available
	}

	if len(updates) == 0 {
		return models.InventoryProduct{}, fmt.Errorf("nothing to update")
	}

	var updatedInvProduct models.InventoryProduct

	db := i.db.WithContext(ctx).
		Model(&models.InventoryProduct{}).
		Clauses(clause.Returning{}).
		Where("id = ?", invProduct.ID).
		Where(
			"store_id IN (?)",
			i.db.Model(&models.Store{}).
				Select("id").
				Where(
					"business_id = ? AND deleted_at IS NULL",
					bizID,
				),
		).
		Updates(updates).
		Scan(&updatedInvProduct)

	if db.Error != nil {
		return models.InventoryProduct{}, db.Error
	}

	if db.RowsAffected == 0 {
		return models.InventoryProduct{}, gorm.ErrRecordNotFound
	}

	return updatedInvProduct, nil
}

// buildTokenFilter returns a SQL fragment that matches products whose search_tokens jsonb
// array contains at least one of the given tokens, using the @> (contains) operator so the
// GIN (jsonb_path_ops) index is used. Tokens are embedded directly — they are safe because
// token.TokenParser normalises them to lowercase alphanumeric before this is called.
func buildTokenFilter(tokens []string) string {
	parts := make([]string, len(tokens))
	for i, tok := range tokens {
		parts[i] = fmt.Sprintf(`p.search_tokens @> '["%s"]'::jsonb`, tok)
	}
	return "(" + strings.Join(parts, " OR ") + ")"
}

func (i *InvProductRepository) SearchProduct(ctx context.Context, tokens []string, lat, long *float64, rangeKm float64) ([]models.ProductResult, error) {
	if len(tokens) == 0 {
		return []models.ProductResult{}, nil
	}
	tokenArray := "{" + strings.Join(tokens, ",") + "}"
	if lat == nil || long == nil {
		return i.searchNoGeo(ctx, tokens, tokenArray)
	}
	return i.searchWithGeo(ctx, tokens, tokenArray, *lat, *long, rangeKm)
}

func (i *InvProductRepository) searchNoGeo(ctx context.Context, tokens []string, tokenArray string) ([]models.ProductResult, error) {
	tokenFilter := buildTokenFilter(tokens)
	var results []models.ProductResult
	err := i.db.WithContext(ctx).Raw(`
		WITH matching AS (
			SELECT p.id, mc.cnt
			FROM products p
			JOIN inventory_products ip ON ip.product_id = p.id
			                           AND ip.deleted_at IS NULL
			                           AND ip.available  = true
			JOIN stores s ON s.id = ip.store_id AND s.deleted_at IS NULL
			CROSS JOIN LATERAL (
				SELECT COUNT(DISTINCT v)::int AS cnt
				FROM jsonb_array_elements_text(p.search_tokens) t(v)
				WHERE v = ANY(?::text[])
			) mc
			WHERE p.deleted_at IS NULL
			  AND `+tokenFilter+`
			GROUP BY p.id, mc.cnt
			ORDER BY mc.cnt DESC
			LIMIT 200
		)
		SELECT p.*,
		       loc.id       AS store_id,
		       loc.name     AS store_name,
		       loc.lat,
		       loc.long,
		       loc.geo_hash AS geo_hash
		FROM products p
		JOIN LATERAL (
			SELECT s.id, s.name, s.lat, s.long, s.geo_hash
			FROM inventory_products ip
			JOIN stores s ON s.id = ip.store_id AND s.deleted_at IS NULL
			WHERE ip.product_id = p.id
			  AND ip.deleted_at IS NULL
			  AND ip.available  = true
			LIMIT 1
		) loc ON TRUE
		WHERE p.id IN (SELECT id FROM matching)
	`, tokenArray).Scan(&results).Error
	return results, err
}

func (i *InvProductRepository) searchWithGeo(ctx context.Context, tokens []string, tokenArray string, lat, long, rangeKm float64) ([]models.ProductResult, error) {
	// Bounding box pre-filter (index-friendly); Haversine in LATERAL for closest-store ordering.
	latDelta := rangeKm / 111.32
	lonDelta := rangeKm / (111.32 * math.Cos(lat*math.Pi/180))
	minLat, maxLat := lat-latDelta, lat+latDelta
	minLon, maxLon := long-lonDelta, long+lonDelta

	tokenFilter := buildTokenFilter(tokens)
	var results []models.ProductResult
	err := i.db.WithContext(ctx).Raw(`
		WITH matching AS (
			SELECT p.id, mc.cnt
			FROM products p
			JOIN inventory_products ip ON ip.product_id = p.id
			                           AND ip.deleted_at IS NULL
			                           AND ip.available  = true
			JOIN stores s ON s.id = ip.store_id
			              AND s.deleted_at IS NULL
			              AND s.lat  BETWEEN ? AND ?
			              AND s.long BETWEEN ? AND ?
			CROSS JOIN LATERAL (
				SELECT COUNT(DISTINCT v)::int AS cnt
				FROM jsonb_array_elements_text(p.search_tokens) t(v)
				WHERE v = ANY(?::text[])
			) mc
			WHERE p.deleted_at IS NULL
			  AND `+tokenFilter+`
			GROUP BY p.id, mc.cnt
			ORDER BY mc.cnt DESC
			LIMIT 200
		)
		SELECT p.*,
		       loc.id       AS store_id,
		       loc.name     AS store_name,
		       loc.lat,
		       loc.long,
		       loc.geo_hash AS geo_hash
		FROM products p
		JOIN LATERAL (
			SELECT s.id, s.name, s.lat, s.long, s.geo_hash
			FROM inventory_products ip
			JOIN stores s ON s.id = ip.store_id AND s.deleted_at IS NULL
			WHERE ip.product_id = p.id
			  AND ip.deleted_at IS NULL
			  AND ip.available  = true
			  AND s.lat  BETWEEN ? AND ?
			  AND s.long BETWEEN ? AND ?
			ORDER BY acos(LEAST(1.0,
				sin(radians(?)) * sin(radians(s.lat)) +
				cos(radians(?)) * cos(radians(s.lat)) * cos(radians(s.long - ?))
			)) * 6371 ASC
			LIMIT 1
		) loc ON TRUE
		WHERE p.id IN (SELECT id FROM matching)
	`,
		minLat, maxLat, minLon, maxLon,
		tokenArray,
		minLat, maxLat, minLon, maxLon,
		lat, lat, long,
	).Scan(&results).Error
	return results, err
}

func (i *InvProductRepository) RemoveProduct(ctx context.Context, invProdID uint, bizID uint) error {
	db := i.db.WithContext(ctx).
		Where("id = ?", invProdID).
		Where(
			"store_id IN (?)",
			i.db.Model(&models.Store{}).
				Select("id").
				Where(
					"business_id = ? AND deleted_at IS NULL",
					bizID,
				),
		).
		Delete(&models.InventoryProduct{})

	if db.Error != nil {
		return db.Error
	}

	if db.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}
