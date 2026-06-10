package repository

import (
	"context"

	"github.com/atharvyadav96k/spotnearr/vendor-svc/models"
)

type IBusinessesRepository interface {
	Create(ctx context.Context, business *models.Business) (*models.Business, error)
	GetByID(ctx context.Context, id uint) (*models.Business, error)
	Update(ctx context.Context, businessID uint, businessName string, desc string) (*models.Business, error)
	Delete(ctx context.Context, id uint) error
	HardDelete(ctx context.Context, id uint) error
	GetByEmail(ctx context.Context, email string) (*models.Business, error)
	GetByPhone(ctx context.Context, phone string) (*models.Business, error)
	VerifyBusiness(ctx context.Context, id uint) error
	ToggleActiveStatus(ctx context.Context, id uint, isActive bool) error
}
