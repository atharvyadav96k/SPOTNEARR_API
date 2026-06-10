package repository

import (
	"context"

	"github.com/atharvyadav96k/spotnearr/vendor-svc/models"
)

type ICategoryRepository interface {
	CreateCategory(ctx context.Context, category models.Category) (models.Category, error)
	GetAllCategories(ctx context.Context) ([]models.Category, error)
	GetCategoryByIDs(ctx context.Context, ids []uint) ([]models.Category, error)
}
