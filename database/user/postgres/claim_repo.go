package postgres

import (
	"context"

	"github.com/Developer-Aadesh/spotnearr-database/user"
	"gorm.io/gorm"
)

type ClaimRepository struct{ db *gorm.DB }

func NewClaimRepository(db *gorm.DB) *ClaimRepository {
	return &ClaimRepository{db: db}
}

func (c *ClaimRepository) Create(ctx context.Context, claim user.Claim) (user.Claim, error) {
	if err := c.db.WithContext(ctx).Create(&claim).Error; err != nil {
		return user.Claim{}, err
	}
	return claim, nil
}

func (c *ClaimRepository) GetByID(ctx context.Context, claimID uint) (user.Claim, error) {
	var claim user.Claim
	err := c.db.WithContext(ctx).First(&claim, claimID).Error
	return claim, err
}

func (c *ClaimRepository) GetByUserID(ctx context.Context, userID uint) ([]user.Claim, error) {
	var claims []user.Claim
	err := c.db.WithContext(ctx).Where("user_id = ?", userID).Find(&claims).Error
	return claims, err
}

func (c *ClaimRepository) GetByUserAndInvProduct(ctx context.Context, userID uint, invProductID uint) (user.Claim, error) {
	var claim user.Claim
	err := c.db.WithContext(ctx).
		Where("user_id = ? AND inventory_product_id = ?", userID, invProductID).
		First(&claim).Error
	return claim, err
}

func (c *ClaimRepository) Delete(ctx context.Context, claimID uint, userID uint) error {
	db := c.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", claimID, userID).
		Delete(&user.Claim{})
	if db.Error != nil {
		return db.Error
	}
	if db.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
