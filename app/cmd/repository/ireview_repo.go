package repository

import (
	"context"

	"github.com/atharvyadav96k/SPOTNEARR_API/models"
)

type IReviewRepository interface {
	Add(ctx context.Context, review models.Review) (models.Review, error)
	Update(ctx context.Context, userID uint, targetType models.ReviewTarget, targetID uint, stars uint8, comment string) (models.Review, error)
	Delete(ctx context.Context, userID uint, targetType models.ReviewTarget, targetID uint) error
	GetByTarget(ctx context.Context, targetType models.ReviewTarget, targetID uint) ([]models.Review, error)
	GetByUserAndTarget(ctx context.Context, userID uint, targetType models.ReviewTarget, targetID uint) (models.Review, error)
}
