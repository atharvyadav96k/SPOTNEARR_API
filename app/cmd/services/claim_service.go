package services

import (
	"context"
	"errors"

	"github.com/atharvyadav96k/SPOTNEARR_API/models"
	"github.com/atharvyadav96k/SPOTNEARR_API/utils/response"
	"gorm.io/gorm"
)

type ClaimService struct {
	base_service
}

func NewClaimService(b base_service) *ClaimService {
	return &ClaimService{base_service: b}
}

func (c *ClaimService) ClaimProduct(userID uint, invProductID uint) response.Res {
	ctx := context.Background()

	invProduct, err := c.RepoClaim().GetInvProductForClaim(ctx, invProductID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.ResponseNotFound("Product not found in inventory")
		}
		return c.ResponseInternalServer("Failed to process claim")
	}

	if invProduct.Available != nil && !*invProduct.Available {
		return c.ResponseBadRequest("Product is not available for claiming")
	}

	_, err = c.RepoClaim().GetByUserAndInvProduct(ctx, userID, invProductID)
	if err == nil {
		return c.ResponseConflict("You have already claimed this product")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return c.ResponseInternalServer("Failed to process claim")
	}

	claim := models.NewClaim(userID, invProductID)
	created, err := c.RepoClaim().Create(ctx, claim)
	if err != nil {
		return c.ResponseInternalServer("Failed to create claim")
	}

	return c.ResponseCreated("Claim submitted successfully", created)
}

func (c *ClaimService) GetUserClaims(userID uint) response.Res {
	claims, err := c.RepoClaim().GetByUserID(context.Background(), userID)
	if err != nil {
		return c.ResponseInternalServer("Failed to fetch claims")
	}
	return c.ResponseOK("Claims fetched successfully", claims)
}

func (c *ClaimService) RemoveClaim(userID uint, claimID uint) response.Res {
	ctx := context.Background()

	err := c.RepoClaim().Delete(ctx, claimID, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.ResponseNotFound("Claim not found")
		}
		return c.ResponseInternalServer("Failed to remove claim")
	}

	return c.ResponseOK("Claim removed successfully", nil)
}
