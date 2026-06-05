package implementation

import (
	"context"

	"github.com/atharvyadav96k/SPOTNEARR_API/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ReviewRepository struct {
	db *gorm.DB
}

func NewReviewRepository(db *gorm.DB) *ReviewRepository {
	return &ReviewRepository{db: db}
}

func (r *ReviewRepository) Add(ctx context.Context, review models.Review) (models.Review, error) {
	err := r.db.WithContext(ctx).Create(&review).Error
	return review, err
}

func (r *ReviewRepository) Update(ctx context.Context, userID uint, targetType models.ReviewTarget, targetID uint, stars uint8, comment string) (models.Review, error) {
	var updated models.Review
	db := r.db.WithContext(ctx).
		Model(&models.Review{}).
		Clauses(clause.Returning{}).
		Where("user_id = ? AND target_type = ? AND target_id = ? AND deleted_at IS NULL", userID, targetType, targetID).
		Updates(map[string]any{"stars": stars, "comment": comment}).
		Scan(&updated)

	if db.Error != nil {
		return models.Review{}, db.Error
	}
	if db.RowsAffected == 0 {
		return models.Review{}, gorm.ErrRecordNotFound
	}
	return updated, nil
}

func (r *ReviewRepository) Delete(ctx context.Context, userID uint, targetType models.ReviewTarget, targetID uint) error {
	db := r.db.WithContext(ctx).
		Where("user_id = ? AND target_type = ? AND target_id = ? AND deleted_at IS NULL", userID, targetType, targetID).
		Delete(&models.Review{})

	if db.Error != nil {
		return db.Error
	}
	if db.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *ReviewRepository) GetByTarget(ctx context.Context, targetType models.ReviewTarget, targetID uint) ([]models.Review, error) {
	var reviews []models.Review
	err := r.db.WithContext(ctx).
		Where("target_type = ? AND target_id = ? AND deleted_at IS NULL", targetType, targetID).
		Order("created_at DESC").
		Find(&reviews).Error
	return reviews, err
}

func (r *ReviewRepository) GetByUserAndTarget(ctx context.Context, userID uint, targetType models.ReviewTarget, targetID uint) (models.Review, error) {
	var review models.Review
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND target_type = ? AND target_id = ? AND deleted_at IS NULL", userID, targetType, targetID).
		First(&review).Error
	return review, err
}
