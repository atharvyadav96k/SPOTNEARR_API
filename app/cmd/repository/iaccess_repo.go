package repository

import (
	"context"

	"github.com/atharvyadav96k/SPOTNEARR_API/models"
)

type IAccessRepository interface {
	GetAccessByUserId(ctx context.Context, id uint) (models.BusinessAccess, error)
	GetAccessByBusinessId(ctx context.Context, id uint) (models.BusinessAccess, error)
	CreateNewAccess(ctx context.Context, access models.BusinessAccess) error
	RemoveAccessByUserId(ctx context.Context, id uint) error
}
