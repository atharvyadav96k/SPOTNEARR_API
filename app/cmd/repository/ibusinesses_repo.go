package repository

import (
	"context"

	"github.com/atharvyadav96k/SPOTNEARR_API/models"
)

type IBusinessesRepository interface {
	Create(ctx context.Context, business *models.Business) (*models.Business, error)
	GetByID(ctx context.Context, id uint) (*models.Business, error)
	Update(ctx context.Context, business *models.Business) (*models.Business, error)
	Delete(ctx context.Context, id uint) error

	GetByUserID(ctx context.Context, userID uint) (*models.Business, error)
	GetByEmail(ctx context.Context, email string) (*models.Business, error)
	GetByPhone(ctx context.Context, phone string) (*models.Business, error)

	VerifyBusiness(ctx context.Context, id uint) error
	ToggleActiveStatus(ctx context.Context, id uint, isActive bool) error
}
