package search

import "gorm.io/gorm"

func AutoMigrate(db *gorm.DB) error {
	// Legacy: convert category_ids from integer[] to jsonb if the table already existed.
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

	// Create / update tables. AutoMigrate adds columns; it never drops them.
	if err := db.AutoMigrate(&SearchEntry{}, &TokenCategoryFreq{}); err != nil {
		return err
	}

	// Drop columns that no longer belong in the lean search index.
	// IF EXISTS makes this a no-op on fresh installs.
	db.Exec(`
		ALTER TABLE search_entries
			DROP COLUMN IF EXISTS business_id,
			DROP COLUMN IF EXISTS product_name,
			DROP COLUMN IF EXISTS price,
			DROP COLUMN IF EXISTS price_unit,
			DROP COLUMN IF EXISTS quantity,
			DROP COLUMN IF EXISTS quantity_unit,
			DROP COLUMN IF EXISTS description,
			DROP COLUMN IF EXISTS store_id,
			DROP COLUMN IF EXISTS store_name,
			DROP COLUMN IF EXISTS street_address
	`)

	db.Exec(`CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_se_tokens_gin
		ON search_entries USING GIN (search_tokens jsonb_path_ops)
		WHERE deleted_at IS NULL AND available = true`)
	db.Exec(`CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_se_lat_long
		ON search_entries (lat, long)
		WHERE deleted_at IS NULL AND available = true`)
	return nil
}
