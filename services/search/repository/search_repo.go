package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"strings"

	"github.com/atharvyadav96k/spotnearr/search-svc/models"
	"gorm.io/gorm"
)

type SearchRepository struct{ db *gorm.DB }

func NewSearchRepository(db *gorm.DB) *SearchRepository {
	return &SearchRepository{db: db}
}

type SearchRow struct {
	models.SearchEntry
	TokenMatchCnt int `gorm:"column:token_match_cnt"`
}

func buildTokenFilter(tokens []string) string {
	parts := make([]string, len(tokens))
	for i, tok := range tokens {
		parts[i] = fmt.Sprintf(`se.search_tokens @> '["%s"]'::jsonb`, tok)
	}
	return "(" + strings.Join(parts, " OR ") + ")"
}

func (r *SearchRepository) Search(ctx context.Context, tokens []string, lat, long *float64, rangeKm float64) ([]SearchRow, error) {
	if len(tokens) == 0 {
		return nil, nil
	}
	if lat == nil || long == nil {
		return r.searchNoGeo(ctx, tokens)
	}
	return r.searchWithGeo(ctx, tokens, *lat, *long, rangeKm)
}

func (r *SearchRepository) searchNoGeo(ctx context.Context, tokens []string) ([]SearchRow, error) {
	tokenArray := "{" + strings.Join(tokens, ",") + "}"
	tokenFilter := buildTokenFilter(tokens)
	var rows []SearchRow
	err := r.db.WithContext(ctx).Raw(`
		SELECT se.*, mc.cnt AS token_match_cnt
		FROM search_entries se
		CROSS JOIN LATERAL (
			SELECT COUNT(DISTINCT v)::int AS cnt
			FROM jsonb_array_elements_text(se.search_tokens) t(v)
			WHERE v = ANY(?::text[])
		) mc
		WHERE se.deleted_at IS NULL AND se.available = true AND `+tokenFilter+`
		ORDER BY mc.cnt DESC LIMIT 200
	`, tokenArray).Scan(&rows).Error
	return rows, err
}

func (r *SearchRepository) searchWithGeo(ctx context.Context, tokens []string, lat, long, rangeKm float64) ([]SearchRow, error) {
	latDelta := rangeKm / 111.32
	lonDelta := rangeKm / (111.32 * math.Cos(lat*math.Pi/180))
	minLat, maxLat := lat-latDelta, lat+latDelta
	minLon, maxLon := long-lonDelta, long+lonDelta

	tokenArray := "{" + strings.Join(tokens, ",") + "}"
	tokenFilter := buildTokenFilter(tokens)
	var rows []SearchRow
	err := r.db.WithContext(ctx).Raw(`
		SELECT se.*, mc.cnt AS token_match_cnt
		FROM search_entries se
		CROSS JOIN LATERAL (
			SELECT COUNT(DISTINCT v)::int AS cnt
			FROM jsonb_array_elements_text(se.search_tokens) t(v)
			WHERE v = ANY(?::text[])
		) mc
		WHERE se.deleted_at IS NULL AND se.available = true
		  AND se.lat BETWEEN ? AND ? AND se.long BETWEEN ? AND ?
		  AND `+tokenFilter+`
		ORDER BY mc.cnt DESC LIMIT 200
	`, tokenArray, minLat, maxLat, minLon, maxLon).Scan(&rows).Error
	return rows, err
}

func (r *SearchRepository) Upsert(ctx context.Context, entry models.SearchEntry) error {
	tokensJSON, err := json.Marshal(entry.SearchTokens)
	if err != nil {
		return fmt.Errorf("marshal search_tokens: %w", err)
	}
	catsJSON, err := json.Marshal(entry.CategoryIDs)
	if err != nil {
		return fmt.Errorf("marshal category_ids: %w", err)
	}
	return r.db.WithContext(ctx).Exec(`
		INSERT INTO search_entries (
			id, product_id, name, price, price_unit,
			search_tokens, category_ids,
			lat, long, geo_hash, available, updated_at, deleted_at
		) VALUES (?, ?, ?, ?, ?, ?::jsonb, ?::jsonb, ?, ?, ?, ?, NOW(), NULL)
		ON CONFLICT (id) DO UPDATE SET
			product_id    = EXCLUDED.product_id,
			name          = EXCLUDED.name,
			price         = EXCLUDED.price,
			price_unit    = EXCLUDED.price_unit,
			search_tokens = EXCLUDED.search_tokens,
			category_ids  = EXCLUDED.category_ids,
			lat           = EXCLUDED.lat,
			long          = EXCLUDED.long,
			geo_hash      = EXCLUDED.geo_hash,
			available     = EXCLUDED.available,
			updated_at    = NOW(),
			deleted_at    = NULL
	`, entry.ID, entry.ProductID, entry.Name, entry.Price, entry.PriceUnit,
		string(tokensJSON), string(catsJSON),
		entry.Lat, entry.Long, entry.GeoHash, entry.Available).Error
}

func (r *SearchRepository) SoftDelete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Exec(`
		UPDATE search_entries
		SET deleted_at = NOW(), available = false, updated_at = NOW()
		WHERE id = ?
	`, id).Error
}
