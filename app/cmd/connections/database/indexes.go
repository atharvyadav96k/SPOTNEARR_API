package database

import "gorm.io/gorm"

func RunIndexes(db *gorm.DB) error {
	stmts := []string{
		// Convert search_tokens text → jsonb so the GIN index below can be created.
		`DO $$ BEGIN
		   IF (SELECT data_type FROM information_schema.columns
		       WHERE table_name='products' AND column_name='search_tokens') = 'text' THEN
		       ALTER TABLE products ALTER COLUMN search_tokens TYPE jsonb USING search_tokens::jsonb;
		   END IF;
		 END $$`,

		// GIN index (jsonb_path_ops) supports @> — O(log n) lookup instead of full scan.
		`CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_products_search_tokens_gin
		 ON products USING GIN (search_tokens jsonb_path_ops)`,

		// Partial composite on stores for bbox filter.
		`CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_stores_lat_long
		 ON stores (lat, long) WHERE deleted_at IS NULL`,

		// Partial covering index for the hot inventory join path.
		`CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_inv_products_avail
		 ON inventory_products (product_id, store_id)
		 WHERE deleted_at IS NULL AND available = true`,
	}
	for _, s := range stmts {
		if err := db.Exec(s).Error; err != nil {
			return err
		}
	}
	return nil
}
