package postgres

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Developer-Aadesh/spotnearr-database/vendordb"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type InvProductRepository struct{ db *gorm.DB }

func NewInvProductRepository(db *gorm.DB) *InvProductRepository {
	return &InvProductRepository{db: db}
}

// outboxQueryRow holds all fields needed to build an OutboxPayload from a single JOIN query.
type outboxQueryRow struct {
	InvProductID  uint     `gorm:"column:inv_product_id"`
	ProductID     uint     `gorm:"column:product_id"`
	BusinessID    uint     `gorm:"column:business_id"`
	ProductName   string   `gorm:"column:product_name"`
	Price         float64  `gorm:"column:price"`
	PriceUnit     string   `gorm:"column:price_unit"`
	Quantity      *float64 `gorm:"column:quantity"`
	QuantityUnit  *string  `gorm:"column:quantity_unit"`
	Desc          string   `gorm:"column:desc"`
	SearchTokens  []string `gorm:"column:search_tokens;serializer:json"`
	StoreID       uint     `gorm:"column:store_id"`
	StoreName     string   `gorm:"column:store_name"`
	StreetAddress string   `gorm:"column:street_address"`
	Lat           float64  `gorm:"column:lat"`
	Long          float64  `gorm:"column:long"`
	GeoHash       string   `gorm:"column:geo_hash"`
	Available     bool     `gorm:"column:available"`
}

func fetchOutboxData(ctx context.Context, tx *gorm.DB, invProductID uint) (*outboxQueryRow, []vendordb.CatRef, error) {
	var row outboxQueryRow
	err := tx.WithContext(ctx).Raw(`
		SELECT
			ip.id        AS inv_product_id,
			p.id         AS product_id,
			p.business_id,
			p.name       AS product_name,
			p.price,
			p.price_unit,
			p.quantity,
			p.quantity_unit,
			p.desc,
			p.search_tokens,
			s.id         AS store_id,
			s.name       AS store_name,
			s.street_address,
			s.lat,
			s.long,
			s.geo_hash,
			ip.available
		FROM inventory_products ip
		JOIN products p ON p.id = ip.product_id AND p.deleted_at IS NULL
		JOIN stores  s ON s.id = ip.store_id   AND s.deleted_at IS NULL
		WHERE ip.id = ?
	`, invProductID).Scan(&row).Error
	if err != nil {
		return nil, nil, fmt.Errorf("fetch outbox data: %w", err)
	}
	if row.InvProductID == 0 {
		return nil, nil, fmt.Errorf("inv_product %d not found for outbox", invProductID)
	}

	type catIDRow struct{ CategoryID uint }
	var catRows []catIDRow
	if err := tx.WithContext(ctx).Raw(`
		SELECT category_id FROM product_categories WHERE product_id = ?
	`, row.ProductID).Scan(&catRows).Error; err != nil {
		return nil, nil, fmt.Errorf("fetch categories for outbox: %w", err)
	}
	cats := make([]vendordb.CatRef, len(catRows))
	for i, c := range catRows {
		cats[i] = vendordb.CatRef{ID: c.CategoryID}
	}
	return &row, cats, nil
}

func buildUpsertPayload(row *outboxQueryRow, cats []vendordb.CatRef) ([]byte, error) {
	payload := vendordb.OutboxPayload{
		EventType:     "upsert",
		InvProductID:  row.InvProductID,
		ProductID:     row.ProductID,
		BusinessID:    row.BusinessID,
		ProductName:   row.ProductName,
		Price:         row.Price,
		PriceUnit:     row.PriceUnit,
		Quantity:      row.Quantity,
		QuantityUnit:  row.QuantityUnit,
		Desc:          row.Desc,
		SearchTokens:  row.SearchTokens,
		Categories:    cats,
		StoreID:       row.StoreID,
		StoreName:     row.StoreName,
		StreetAddress: row.StreetAddress,
		Lat:           row.Lat,
		Long:          row.Long,
		GeoHash:       row.GeoHash,
		Available:     row.Available,
	}
	return json.Marshal(payload)
}

func (i *InvProductRepository) AddProduct(ctx context.Context, invProduct vendordb.InventoryProduct, bizID uint) error {
	var count int64
	if err := i.db.WithContext(ctx).
		Model(&vendordb.Store{}).
		Where("id = ? AND business_id = ? AND deleted_at IS NULL", invProduct.StoreID, bizID).
		Count(&count).Error; err != nil {
		return err
	}
	if count == 0 {
		return gorm.ErrRecordNotFound
	}

	return i.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&invProduct).Error; err != nil {
			return err
		}
		row, cats, err := fetchOutboxData(ctx, tx, invProduct.ID)
		if err != nil {
			return err
		}
		payloadBytes, err := buildUpsertPayload(row, cats)
		if err != nil {
			return fmt.Errorf("marshal upsert outbox payload: %w", err)
		}
		return tx.Create(&vendordb.SearchSyncOutbox{
			EventType: "upsert",
			Payload:   payloadBytes,
		}).Error
	})
}

func (i *InvProductRepository) GetInvProductList(ctx context.Context, storeID uint, bizID uint) ([]vendordb.Product, error) {
	var products []vendordb.Product
	err := i.db.WithContext(ctx).
		Model(&vendordb.Product{}).
		Distinct("products.*").
		Joins("JOIN inventory_products ip ON ip.product_id = products.id").
		Joins("JOIN stores s ON s.id = ip.store_id").
		Where("ip.store_id = ?", storeID).
		Where("ip.deleted_at IS NULL").
		Where("s.business_id = ?", bizID).
		Find(&products).Error
	return products, err
}

func (i *InvProductRepository) UpdateProduct(ctx context.Context, invProduct vendordb.InventoryProduct, bizID uint) (vendordb.InventoryProduct, error) {
	updates := map[string]any{}
	availabilityChanged := false

	if invProduct.Count != nil {
		updates["count"] = *invProduct.Count
	}
	if invProduct.Available != nil {
		updates["available"] = *invProduct.Available
		availabilityChanged = true
	}
	if len(updates) == 0 {
		return vendordb.InventoryProduct{}, fmt.Errorf("nothing to update")
	}

	var updated vendordb.InventoryProduct
	err := i.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		db := tx.Model(&vendordb.InventoryProduct{}).
			Clauses(clause.Returning{}).
			Where("id = ?", invProduct.ID).
			Where(
				"store_id IN (?)",
				tx.Model(&vendordb.Store{}).
					Select("id").
					Where("business_id = ? AND deleted_at IS NULL", bizID),
			).
			Updates(updates).
			Scan(&updated)
		if db.Error != nil {
			return db.Error
		}
		if db.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}

		if !availabilityChanged {
			return nil
		}

		row, cats, err := fetchOutboxData(ctx, tx, updated.ID)
		if err != nil {
			return err
		}
		payloadBytes, err := buildUpsertPayload(row, cats)
		if err != nil {
			return fmt.Errorf("marshal upsert outbox payload: %w", err)
		}
		return tx.Create(&vendordb.SearchSyncOutbox{
			EventType: "upsert",
			Payload:   payloadBytes,
		}).Error
	})
	return updated, err
}

func (i *InvProductRepository) RemoveProduct(ctx context.Context, invProdID uint, bizID uint) error {
	return i.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		db := tx.
			Where("id = ?", invProdID).
			Where(
				"store_id IN (?)",
				tx.Model(&vendordb.Store{}).
					Select("id").
					Where("business_id = ? AND deleted_at IS NULL", bizID),
			).
			Delete(&vendordb.InventoryProduct{})
		if db.Error != nil {
			return db.Error
		}
		if db.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}

		payloadBytes, err := json.Marshal(vendordb.OutboxPayload{
			EventType:    "delete",
			InvProductID: invProdID,
		})
		if err != nil {
			return fmt.Errorf("marshal delete outbox payload: %w", err)
		}
		return tx.Create(&vendordb.SearchSyncOutbox{
			EventType: "delete",
			Payload:   payloadBytes,
		}).Error
	})
}

// WriteUpsertOutboxesForProduct writes a search sync upsert entry for every
// active inventory_product that belongs to productID. Call this after a
// product's name, price, or description is updated so the search index stays current.
func (i *InvProductRepository) WriteUpsertOutboxesForProduct(ctx context.Context, productID uint) error {
	var ids []uint
	if err := i.db.WithContext(ctx).Raw(`
		SELECT ip.id FROM inventory_products ip
		WHERE ip.product_id = ? AND ip.deleted_at IS NULL
	`, productID).Scan(&ids).Error; err != nil {
		return fmt.Errorf("list inv_products for product %d: %w", productID, err)
	}
	for _, invID := range ids {
		if err := i.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
			row, cats, err := fetchOutboxData(ctx, tx, invID)
			if err != nil {
				return err
			}
			payload, err := buildUpsertPayload(row, cats)
			if err != nil {
				return err
			}
			return tx.Create(&vendordb.SearchSyncOutbox{EventType: "upsert", Payload: payload}).Error
		}); err != nil {
			return fmt.Errorf("upsert outbox for inv_product %d: %w", invID, err)
		}
	}
	return nil
}

// WriteDeleteOutboxesForProduct writes a search sync delete entry for every
// inventory_product (including soft-deleted) that belongs to productID.
// Call this after a product is deleted so the search index removes the entries.
func (i *InvProductRepository) WriteDeleteOutboxesForProduct(ctx context.Context, productID uint) error {
	var ids []uint
	if err := i.db.WithContext(ctx).Raw(`
		SELECT id FROM inventory_products WHERE product_id = ?
	`, productID).Scan(&ids).Error; err != nil {
		return fmt.Errorf("list inv_products for delete %d: %w", productID, err)
	}
	for _, invID := range ids {
		payload, err := json.Marshal(vendordb.OutboxPayload{
			EventType:    "delete",
			InvProductID: invID,
		})
		if err != nil {
			return err
		}
		if err := i.db.WithContext(ctx).Create(&vendordb.SearchSyncOutbox{
			EventType: "delete",
			Payload:   payload,
		}).Error; err != nil {
			return fmt.Errorf("delete outbox for inv_product %d: %w", invID, err)
		}
	}
	return nil
}

// GetIDsByBusiness returns all inventory_product IDs that belong to a business,
// used by the vendor-side claim listing to scope claims to the right products.
func (i *InvProductRepository) GetIDsByBusiness(ctx context.Context, bizID uint) ([]uint, error) {
	var ids []uint
	err := i.db.WithContext(ctx).Raw(`
		SELECT ip.id FROM inventory_products ip
		JOIN stores s ON s.id = ip.store_id AND s.deleted_at IS NULL
		WHERE s.business_id = ? AND ip.deleted_at IS NULL
	`, bizID).Scan(&ids).Error
	return ids, err
}

func (i *InvProductRepository) GetByID(ctx context.Context, id uint) (*vendordb.InvProductDetail, error) {
	var result vendordb.InvProductDetail
	err := i.db.WithContext(ctx).Raw(`
		SELECT
			ip.id,
			ip.available,
			p.name       AS product_name,
			p.price,
			p.price_unit,
			s.name       AS store_name,
			s.street_address,
			s.lat,
			s.long
		FROM inventory_products ip
		JOIN products p ON p.id = ip.product_id AND p.deleted_at IS NULL
		JOIN stores  s ON s.id = ip.store_id   AND s.deleted_at IS NULL
		WHERE ip.id = ? AND ip.deleted_at IS NULL
	`, id).Scan(&result).Error
	if err != nil {
		return nil, err
	}
	if result.ID == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return &result, nil
}
