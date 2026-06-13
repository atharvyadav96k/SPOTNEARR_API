package vendordb

import "gorm.io/gorm"

func AutoMigrate(db *gorm.DB) error {
	if err := db.AutoMigrate(
		&Business{},
		&Store{},
		&BusinessAccess{},
		&Category{},
		&Product{},
		&InventoryProduct{},
		&ProductToken{},
		&BusinessAccount{},
		&SearchSyncOutbox{},
		&Offer{},
		&Spotlight{},
	); err != nil {
		return err
	}
	return RunIndexes(db)
}

func RunIndexes(db *gorm.DB) error {
	stmts := []string{
		`DO $$ BEGIN
		   IF (SELECT data_type FROM information_schema.columns
		       WHERE table_name='products' AND column_name='search_tokens') = 'text' THEN
		       ALTER TABLE products ALTER COLUMN search_tokens TYPE jsonb USING search_tokens::jsonb;
		   END IF;
		 END $$`,
		`CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_products_search_tokens_gin
		 ON products USING GIN (search_tokens jsonb_path_ops)`,
		`CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_stores_lat_long
		 ON stores (lat, long) WHERE deleted_at IS NULL`,
		`CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_inv_products_avail
		 ON inventory_products (product_id, store_id)
		 WHERE deleted_at IS NULL AND available = true`,
		`CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_search_sync_outbox_unprocessed
		 ON search_sync_outboxes (id) WHERE processed = false`,
	}
	for _, s := range stmts {
		if err := db.Exec(s).Error; err != nil {
			return err
		}
	}
	return nil
}
