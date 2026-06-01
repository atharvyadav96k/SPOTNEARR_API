package repository

import (
	"context"

	"github.com/atharvyadav96k/SPOTNEARR_API/models"
)

type IClaimRepository interface {
	Create(ctx context.Context, claim models.Claim) (models.Claim, error)
	GetByID(ctx context.Context, claimID uint) (models.Claim, error)
	GetByUserID(ctx context.Context, userID uint) ([]models.Claim, error)
	GetByUserAndInvProduct(ctx context.Context, userID uint, invProductID uint) (models.Claim, error)
	GetInvProductForClaim(ctx context.Context, invProductID uint) (models.InventoryProduct, error)
	Delete(ctx context.Context, claimID uint, userID uint) error
}
