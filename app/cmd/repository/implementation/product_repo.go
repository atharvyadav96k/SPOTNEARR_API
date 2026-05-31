package implementation

import (
	"context"
	"fmt"

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

func (p *ProductRepository) AddProduct(ctx context.Context, product models.Product, bizID uint) error {
	product.BusinessID = bizID
	return p.db.WithContext(ctx).Create(&product).Error
}

func (p *ProductRepository) UpdateProduct(ctx context.Context, product models.Product, bizID uint) (models.Product, error) {
	updates := map[string]any{}
	if product.Name != "" {
		updates["name"] = product.Name
	}

	if product.Price > 0 {
		updates["price"] = product.Price
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

	db := p.db.WithContext(ctx).
		Where(
			"id = ? AND business_id = ?",
			productID,
			bizID,
		).
		Delete(&models.Product{})

	if db.Error != nil {
		return db.Error
	}

	if db.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}
