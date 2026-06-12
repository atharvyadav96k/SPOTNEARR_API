package postgres

import (
	"context"

	"github.com/Developer-Aadesh/spotnearr-database/user"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type BusinessFollowRepository struct{ db *gorm.DB }

func NewBusinessFollowRepository(db *gorm.DB) *BusinessFollowRepository {
	return &BusinessFollowRepository{db: db}
}

func (r *BusinessFollowRepository) Follow(ctx context.Context, userID, businessID uint) error {
	return r.db.WithContext(ctx).
		Clauses(clause.OnConflict{DoNothing: true}).
		Create(&user.BusinessFollow{UserID: userID, BusinessID: businessID}).Error
}

func (r *BusinessFollowRepository) Unfollow(ctx context.Context, userID, businessID uint) error {
	result := r.db.WithContext(ctx).
		Where("user_id = ? AND business_id = ?", userID, businessID).
		Delete(&user.BusinessFollow{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *BusinessFollowRepository) IsFollowing(ctx context.Context, userID, businessID uint) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&user.BusinessFollow{}).
		Where("user_id = ? AND business_id = ?", userID, businessID).
		Count(&count).Error
	return count > 0, err
}

func (r *BusinessFollowRepository) GetFollowedBusinessIDs(ctx context.Context, userID uint) ([]uint, error) {
	var follows []user.BusinessFollow
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&follows).Error; err != nil {
		return nil, err
	}
	ids := make([]uint, len(follows))
	for i, f := range follows {
		ids[i] = f.BusinessID
	}
	return ids, nil
}
