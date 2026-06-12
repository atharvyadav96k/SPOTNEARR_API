package postgres

import (
	"context"

	"github.com/Developer-Aadesh/spotnearr-database/user"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ProductEngagementRepository struct{ db *gorm.DB }

func NewProductEngagementRepository(db *gorm.DB) *ProductEngagementRepository {
	return &ProductEngagementRepository{db: db}
}

func (r *ProductEngagementRepository) Like(ctx context.Context, userID, invProductID uint) error {
	return r.db.WithContext(ctx).
		Clauses(clause.OnConflict{DoNothing: true}).
		Create(&user.ProductLike{UserID: userID, InvProductID: invProductID}).Error
}

func (r *ProductEngagementRepository) Unlike(ctx context.Context, userID, invProductID uint) error {
	return r.db.WithContext(ctx).
		Where("user_id = ? AND inv_product_id = ?", userID, invProductID).
		Delete(&user.ProductLike{}).Error
}

func (r *ProductEngagementRepository) Save(ctx context.Context, userID, invProductID uint) error {
	return r.db.WithContext(ctx).
		Clauses(clause.OnConflict{DoNothing: true}).
		Create(&user.ProductSave{UserID: userID, InvProductID: invProductID}).Error
}

func (r *ProductEngagementRepository) Unsave(ctx context.Context, userID, invProductID uint) error {
	return r.db.WithContext(ctx).
		Where("user_id = ? AND inv_product_id = ?", userID, invProductID).
		Delete(&user.ProductSave{}).Error
}
