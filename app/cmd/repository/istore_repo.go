package repository

import (
	"context"

	"github.com/atharvyadav96k/SPOTNEARR_API/models"
)

type IStoreRepository interface {
	GetStoreByBusinessId(ctx context.Context, id uint) ([]models.Store, error)
	CreateStoreByBusinessId(ctx context.Context, store *models.Store) error
	UpdateStoreByBusinessId(ctx context.Context, store models.Store) (models.Store, error)
}
