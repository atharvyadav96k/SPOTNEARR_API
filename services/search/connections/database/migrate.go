package database

import (
	"github.com/atharvyadav96k/spotnearr/search-svc/models"
	"gorm.io/gorm"
)

func AutoMigrate(db *gorm.DB) error {
	if err := db.AutoMigrate(&models.SearchEntry{}); err != nil {
		return err
	}
	// GIN index for token containment queries (@> operator).
	db.Exec(`
		CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_se_tokens_gin
		ON search_entries USING GIN (search_tokens jsonb_path_ops)
		WHERE deleted_at IS NULL AND available = true
	`)
	// Bounding-box geo index for radius filtering.
	db.Exec(`
		CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_se_lat_long
		ON search_entries (lat, long)
		WHERE deleted_at IS NULL AND available = true
	`)
	return nil
}
