package implementation

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/atharvyadav96k/SPOTNEARR_API/factories/token"
	"github.com/atharvyadav96k/SPOTNEARR_API/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ProductRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (p *ProductRepository) AddProduct(ctx context.Context, product models.Product, bizID uint) (models.Product, error) {
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
			// Index only name tokens — desc tokens are for text search, not frequency.
			// Sort so concurrent transactions acquire row locks in the same order (no deadlock).
			nameTokens := token.TokenParser(product.Name)
			sort.Strings(nameTokens)
			sortedTokens := nameTokens

			for _, cat := range categories {
				for _, tok := range sortedTokens {
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
		return models.Product{}, err
	}
	return product, nil
}

func (p *ProductRepository) UpdateProduct(ctx context.Context, product models.Product, bizID uint) (models.Product, error) {
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
		return models.Product{}, fmt.Errorf("nothing to update")
	}

	var updatedProduct models.Product

	db := p.db.WithContext(ctx).Model(&models.Product{}).
		Clauses(clause.Returning{}).
		Where(
			"id = ? AND business_id = ?",
			product.ID,
			bizID,
		).
		Updates(updates).
		Scan(&updatedProduct)

	if db.Error != nil {
		return models.Product{}, db.Error
	}

	if db.RowsAffected == 0 {
		return models.Product{}, gorm.ErrRecordNotFound
	}

	return updatedProduct, nil
}

func (p *ProductRepository) GetProductsByBusiness(ctx context.Context, bizID uint) ([]models.Product, error) {

	var products []models.Product

	err := p.db.WithContext(ctx).
		Where("business_id = ?", bizID).
		Order("created_at DESC").
		Find(&products).Error

	if err != nil {
		return nil, err
	}

	return products, nil
}

func (p *ProductRepository) GetProductById(ctx context.Context, id uint, bizID uint) (models.Product, error) {

	var product models.Product

	err := p.db.WithContext(ctx).
		First(&product, id).Error

	if err != nil {
		return models.Product{}, err
	}

	return product, nil
}

func (p *ProductRepository) DeleteProduct(ctx context.Context, productID uint, bizID uint) error {
	return p.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var product models.Product
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

		db := tx.Where("id = ? AND business_id = ?", productID, bizID).Delete(&models.Product{})
		if db.Error != nil {
			return db.Error
		}
		if db.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}
		return nil
	})
}

func (p *ProductRepository) GetProductCategoryIDs(ctx context.Context, productIDs []uint) (map[uint][]uint, error) {
	if len(productIDs) == 0 {
		return map[uint][]uint{}, nil
	}

	idStrs := make([]string, len(productIDs))
	for i, id := range productIDs {
		idStrs[i] = fmt.Sprintf("%d", id)
	}
	idArray := "{" + strings.Join(idStrs, ",") + "}"

	type row struct {
		ProductID  uint
		CategoryID uint
	}
	var rows []row

	err := p.db.WithContext(ctx).Raw(`
		SELECT product_id, category_id
		FROM product_categories
		WHERE product_id = ANY(?::int[])
	`, idArray).Scan(&rows).Error

	if err != nil {
		return nil, err
	}

	result := make(map[uint][]uint, len(productIDs))
	for _, r := range rows {
		result[r.ProductID] = append(result[r.ProductID], r.CategoryID)
	}
	return result, nil
}

func (p *ProductRepository) SearchProducts(ctx context.Context, tokens []string) ([]models.Product, error) {
	if len(tokens) == 0 {
		return []models.Product{}, nil
	}

	// Build a PostgreSQL text-array literal: {token1,token2,...}
	// Tokens are safe (lowercase alphanumeric only) so embedding directly is fine.
	tokenArray := "{" + strings.Join(tokens, ",") + "}"

	var products []models.Product
	err := p.db.WithContext(ctx).Raw(`
		SELECT p.*
		FROM products p,
		     jsonb_array_elements_text(p.search_tokens::jsonb) AS elem
		WHERE p.deleted_at IS NULL
		  AND elem = ANY(?::text[])
		GROUP BY p.id
		ORDER BY COUNT(DISTINCT elem) DESC
		LIMIT 50
	`, tokenArray).Scan(&products).Error

	return products, err
}
