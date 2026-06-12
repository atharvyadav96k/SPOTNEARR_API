package postgres

import (
	"context"

	"github.com/Developer-Aadesh/spotnearr-database/user"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type SpotlightEngagementRepository struct{ db *gorm.DB }

func NewSpotlightEngagementRepository(db *gorm.DB) *SpotlightEngagementRepository {
	return &SpotlightEngagementRepository{db: db}
}

func (r *SpotlightEngagementRepository) Like(ctx context.Context, userID, spotlightID uint) error {
	return r.db.WithContext(ctx).
		Clauses(clause.OnConflict{DoNothing: true}).
		Create(&user.SpotlightLike{UserID: userID, SpotlightID: spotlightID}).Error
}

func (r *SpotlightEngagementRepository) Unlike(ctx context.Context, userID, spotlightID uint) error {
	return r.db.WithContext(ctx).
		Where("user_id = ? AND spotlight_id = ?", userID, spotlightID).
		Delete(&user.SpotlightLike{}).Error
}

func (r *SpotlightEngagementRepository) Save(ctx context.Context, userID, spotlightID uint) error {
	return r.db.WithContext(ctx).
		Clauses(clause.OnConflict{DoNothing: true}).
		Create(&user.SpotlightSave{UserID: userID, SpotlightID: spotlightID}).Error
}

func (r *SpotlightEngagementRepository) Unsave(ctx context.Context, userID, spotlightID uint) error {
	return r.db.WithContext(ctx).
		Where("user_id = ? AND spotlight_id = ?", userID, spotlightID).
		Delete(&user.SpotlightSave{}).Error
}
