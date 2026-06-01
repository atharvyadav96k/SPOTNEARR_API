package implementation

import (
	"context"

	"github.com/atharvyadav96k/SPOTNEARR_API/models"
	"gorm.io/gorm"
)

type ClaimRepository struct {
	db *gorm.DB
}

func NewClaimRepository(db *gorm.DB) *ClaimRepository {
	return &ClaimRepository{db: db}
}

func (c *ClaimRepository) Create(ctx context.Context, claim models.Claim) (models.Claim, error) {
	if err := c.db.WithContext(ctx).Create(&claim).Error; err != nil {
		return models.Claim{}, err
	}
	return claim, nil
}

func (c *ClaimRepository) GetByID(ctx context.Context, claimID uint) (models.Claim, error) {
	var claim models.Claim
	err := c.db.WithContext(ctx).
		Preload("InventoryProduct.Product").
		Preload("InventoryProduct.Store").
		First(&claim, claimID).Error
	return claim, err
}

func (c *ClaimRepository) GetByUserID(ctx context.Context, userID uint) ([]models.Claim, error) {
	var claims []models.Claim
	err := c.db.WithContext(ctx).
		Preload("InventoryProduct.Product").
		Preload("InventoryProduct.Store").
		Where("user_id = ?", userID).
		Find(&claims).Error
	return claims, err
}

func (c *ClaimRepository) GetByUserAndInvProduct(ctx context.Context, userID uint, invProductID uint) (models.Claim, error) {
	var claim models.Claim
	err := c.db.WithContext(ctx).
		Where("user_id = ? AND inventory_product_id = ?", userID, invProductID).
		First(&claim).Error
	return claim, err
}

func (c *ClaimRepository) GetInvProductForClaim(ctx context.Context, invProductID uint) (models.InventoryProduct, error) {
	var invProduct models.InventoryProduct
	err := c.db.WithContext(ctx).
		Preload("Product").
		Preload("Store").
		Where("id = ? AND deleted_at IS NULL", invProductID).
		First(&invProduct).Error
	return invProduct, err
}

func (c *ClaimRepository) Delete(ctx context.Context, claimID uint, userID uint) error {
	db := c.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", claimID, userID).
		Delete(&models.Claim{})
	if db.Error != nil {
		return db.Error
	}
	if db.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
