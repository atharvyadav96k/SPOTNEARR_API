package postgres

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"

	"github.com/Developer-Aadesh/spotnearr-database/vendordb"
	"github.com/atharvyadav96k/spotnearr/pkg/tokenizer"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ProductRepository struct{ db *gorm.DB }

func NewProductRepository(db *gorm.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (p *ProductRepository) AddProduct(ctx context.Context, product vendordb.Product, bizID uint) (vendordb.Product, error) {
	product.BusinessID = bizID
	categories := product.Categories
	product.Categories = nil

	err := p.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Omit("Categories.*").Create(&product).Error; err != nil {
			return err
		}
		if len(categories) > 0 {
			if err := tx.Model(&product).Association("Categories").Append(categories); err != nil {
				return err
			}
			nameTokens := tokenizer.TokenParser(product.Name)
			sort.Strings(nameTokens)
			for _, cat := range categories {
				for _, tok := range nameTokens {
					if err := tx.Exec(`
						INSERT INTO product_tokens (token, category_id, count, updated_at)
						VALUES (?, ?, 1, NOW())
						ON CONFLICT (token, category_id) DO UPDATE
						SET count = product_tokens.count + 1, updated_at = NOW()
					`, tok, cat.ID).Error; err != nil {
						return err
					}
				}
			}
		}
		return nil
	})
	if err != nil {
		return vendordb.Product{}, err
	}
	product.Categories = categories
	return product, nil
}

func (p *ProductRepository) UpdateProduct(ctx context.Context, product vendordb.Product, bizID uint) (vendordb.Product, error) {
	updates := map[string]any{}
	if product.Name != "" {
		updates["name"] = product.Name
	}
	if product.Price > 0 {
		updates["price"] = product.Price
	}
	if product.PriceUnit != "" {
		updates["price_unit"] = product.PriceUnit
	}
	if product.Quantity != nil {
		updates["quantity"] = product.Quantity
	}
	if product.QuantityUnit != nil {
		updates["quantity_unit"] = product.QuantityUnit
	}
	if product.Desc != "" {
		updates["desc"] = product.Desc
	}
	if len(updates) == 0 {
		return vendordb.Product{}, fmt.Errorf("nothing to update")
	}

	var updated vendordb.Product
	db := p.db.WithContext(ctx).Model(&vendordb.Product{}).
		Clauses(clause.Returning{}).
		Where("id = ? AND business_id = ?", product.ID, bizID).
		Updates(updates).
		Scan(&updated)
	if db.Error != nil {
		return vendordb.Product{}, db.Error
	}
	if db.RowsAffected == 0 {
		return vendordb.Product{}, gorm.ErrRecordNotFound
	}
	return updated, nil
}

func (p *ProductRepository) GetProductsByBusiness(ctx context.Context, bizID uint) ([]vendordb.Product, error) {
	var products []vendordb.Product
	err := p.db.WithContext(ctx).
		Where("business_id = ?", bizID).
		Order("created_at DESC").
		Find(&products).Error
	return products, err
}

func (p *ProductRepository) GetProductById(ctx context.Context, id uint, bizID uint) (vendordb.Product, error) {
	var product vendordb.Product
	err := p.db.WithContext(ctx).First(&product, id).Error
	return product, err
}

func (p *ProductRepository) DeleteProduct(ctx context.Context, productID uint, bizID uint) error {
	return p.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var product vendordb.Product
		if err := tx.Preload("Categories").First(&product, "id = ? AND business_id = ?", productID, bizID).Error; err != nil {
			return err
		}

		for _, cat := range product.Categories {
			for _, tok := range product.SearchTokens {
				tx.Exec(`
					UPDATE product_tokens
					SET count = GREATEST(count - 1, 0), updated_at = NOW()
					WHERE token = ? AND category_id = ?
				`, tok, cat.ID)
			}
		}

		var invIDs []uint
		if err := tx.Model(&vendordb.InventoryProduct{}).
			Where("product_id = ? AND deleted_at IS NULL", productID).
			Pluck("id", &invIDs).Error; err != nil {
			return err
		}

		for _, invID := range invIDs {
			payload, err := json.Marshal(vendordb.OutboxPayload{
				EventType:    "delete",
				InvProductID: invID,
			})
			if err != nil {
				return fmt.Errorf("marshal delete outbox payload: %w", err)
			}
			if err := tx.Create(&vendordb.SearchSyncOutbox{
				EventType: "delete",
				Payload:   payload,
			}).Error; err != nil {
				return fmt.Errorf("write delete outbox: %w", err)
			}
		}

		if len(invIDs) > 0 {
			if err := tx.Where("product_id = ?", productID).Delete(&vendordb.InventoryProduct{}).Error; err != nil {
				return err
			}
		}

		db := tx.Where("id = ? AND business_id = ?", productID, bizID).Delete(&vendordb.Product{})
		if db.Error != nil {
			return db.Error
		}
		if db.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}
		return nil
	})
}
