package search

import "gorm.io/gorm"

func AutoMigrate(db *gorm.DB) error {
	// Convert category_ids from integer[] to jsonb (preserves data via array_to_json).
	// Safe to run repeatedly — the DO block is a no-op once the column is already jsonb.
	db.Exec(`
		DO $$
		BEGIN
			IF EXISTS (
				SELECT 1 FROM information_schema.columns
				WHERE table_name  = 'search_entries'
				  AND column_name = 'category_ids'
				  AND udt_name    = '_int4'
			) THEN
				ALTER TABLE search_entries
				ALTER COLUMN category_ids TYPE jsonb
				USING array_to_json(category_ids)::jsonb;
			END IF;
		END $$;
	`)

	if err := db.AutoMigrate(&SearchEntry{}); err != nil {
		return err
	}
	db.Exec(`
		CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_se_tokens_gin
		ON search_entries USING GIN (search_tokens jsonb_path_ops)
		WHERE deleted_at IS NULL AND available = true
	`)
	db.Exec(`
		CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_se_lat_long
		ON search_entries (lat, long)
		WHERE deleted_at IS NULL AND available = true
	`)
	return nil
}
