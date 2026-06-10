package postgres

import (
	"context"

	"github.com/Developer-Aadesh/spotnearr-database/user"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ReviewRepository struct{ db *gorm.DB }

func NewReviewRepository(db *gorm.DB) *ReviewRepository {
	return &ReviewRepository{db: db}
}

func (r *ReviewRepository) Add(ctx context.Context, review user.Review) (user.Review, error) {
	err := r.db.WithContext(ctx).Create(&review).Error
	return review, err
}

func (r *ReviewRepository) Update(ctx context.Context, userID uint, targetType user.ReviewTarget, targetID uint, stars uint8, comment string) (user.Review, error) {
	var updated user.Review
	db := r.db.WithContext(ctx).
		Model(&user.Review{}).
		Clauses(clause.Returning{}).
		Where("user_id = ? AND target_type = ? AND target_id = ? AND deleted_at IS NULL", userID, targetType, targetID).
		Updates(map[string]any{"stars": stars, "comment": comment}).
		Scan(&updated)
	if db.Error != nil {
		return user.Review{}, db.Error
	}
	if db.RowsAffected == 0 {
		return user.Review{}, gorm.ErrRecordNotFound
	}
	return updated, nil
}

func (r *ReviewRepository) Delete(ctx context.Context, userID uint, targetType user.ReviewTarget, targetID uint) error {
	db := r.db.WithContext(ctx).
		Where("user_id = ? AND target_type = ? AND target_id = ? AND deleted_at IS NULL", userID, targetType, targetID).
		Delete(&user.Review{})
	if db.Error != nil {
		return db.Error
	}
	if db.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *ReviewRepository) GetByTarget(ctx context.Context, targetType user.ReviewTarget, targetID uint) ([]user.Review, error) {
	var reviews []user.Review
	err := r.db.WithContext(ctx).
		Where("target_type = ? AND target_id = ? AND deleted_at IS NULL", targetType, targetID).
		Order("created_at DESC").
		Find(&reviews).Error
	return reviews, err
}

func (r *ReviewRepository) GetByUserAndTarget(ctx context.Context, userID uint, targetType user.ReviewTarget, targetID uint) (user.Review, error) {
	var review user.Review
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND target_type = ? AND target_id = ? AND deleted_at IS NULL", userID, targetType, targetID).
		First(&review).Error
	return review, err
}
