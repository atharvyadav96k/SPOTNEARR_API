package implementation

import (
	"context"
	"fmt"

	"github.com/atharvyadav96k/SPOTNEARR_API/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type InvProductRepository struct {
	db *gorm.DB
}

func NewInvProductRepository(db *gorm.DB) *InvProductRepository {
	return &InvProductRepository{db: db}
}

func (i *InvProductRepository) AddProduct(ctx context.Context, invProduct models.InventoryProduct, bizID uint) error {
	var count int64
	err := i.db.WithContext(ctx).
		Model(&models.Store{}).
		Where(
			"id = ? AND business_id = ? AND deleted_at IS NULL",
			invProduct.StoreID,
			bizID,
		).
		Count(&count).Error

	if err != nil {
		return err
	}

	if count == 0 {
		return gorm.ErrRecordNotFound
	}

	return i.db.WithContext(ctx).
		Create(&invProduct).Error
}

func (i *InvProductRepository) GetInvProductList(ctx context.Context, invID uint, bizID uint) ([]models.Product, error) {
	var products []models.Product

	err := i.db.WithContext(ctx).
		Model(&models.Product{}).
		Distinct("products.*").
		Joins("JOIN inventory_products ip ON ip.product_id = products.id").
		Joins("JOIN stores s ON s.id = ip.store_id").
		Where("ip.store_id = ?", invID).
		Where("ip.deleted_at IS NULL").
		Where("s.business_id = ?", bizID).
		Find(&products).Error

	if err != nil {
		return nil, err
	}

	return products, nil
}

func (i *InvProductRepository) UpdateProduct(ctx context.Context, invProduct models.InventoryProduct, bizID uint) (models.InventoryProduct, error) {
	updates := map[string]any{}

	if invProduct.Count != nil {
		updates["count"] = *invProduct.Count
	}

	if invProduct.Available != nil {
		updates["available"] = *invProduct.Available
	}

	if len(updates) == 0 {
		return models.InventoryProduct{}, fmt.Errorf("nothing to update")
	}

	var updatedInvProduct models.InventoryProduct

	db := i.db.WithContext(ctx).
		Model(&models.InventoryProduct{}).
		Clauses(clause.Returning{}).
		Where("id = ?", invProduct.ID).
		Where(
			"store_id IN (?)",
			i.db.Model(&models.Store{}).
				Select("id").
				Where(
					"business_id = ? AND deleted_at IS NULL",
					bizID,
				),
		).
		Updates(updates).
		Scan(&updatedInvProduct)

	if db.Error != nil {
		return models.InventoryProduct{}, db.Error
	}

	if db.RowsAffected == 0 {
		return models.InventoryProduct{}, gorm.ErrRecordNotFound
	}

	return updatedInvProduct, nil
}

func (i *InvProductRepository) RemoveProduct(ctx context.Context, invProdID uint, bizID uint) error {
	db := i.db.WithContext(ctx).
		Where("id = ?", invProdID).
		Where(
			"store_id IN (?)",
			i.db.Model(&models.Store{}).
				Select("id").
				Where(
					"business_id = ? AND deleted_at IS NULL",
					bizID,
				),
		).
		Delete(&models.InventoryProduct{})

	if db.Error != nil {
		return db.Error
	}

	if db.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}
