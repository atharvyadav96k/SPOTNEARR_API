package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/atharvyadav96k/SPOTNEARR_API/config"
	"github.com/atharvyadav96k/SPOTNEARR_API/utils/response"
	usermodel "github.com/Developer-Aadesh/spotnearr-database/user"
	"gorm.io/gorm"
)

type ClaimService struct {
	base_service
	httpClient *http.Client
}

func NewClaimService(b base_service) *ClaimService {
	return &ClaimService{
		base_service: b,
		httpClient:   &http.Client{Timeout: 500 * time.Millisecond},
	}
}

func (c *ClaimService) fetchInvProduct(invProductID uint) (json.RawMessage, bool, error) {
	url := fmt.Sprintf("%s/internal/inventory-products/%d", config.C.VendorServiceURL, invProductID)
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, false, fmt.Errorf("vendor service unreachable: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, false, nil
	}
	if resp.StatusCode != http.StatusOK {
		return nil, false, fmt.Errorf("vendor service returned %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, false, err
	}

	var detail struct {
		Available bool `json:"available"`
	}
	if err := json.Unmarshal(body, &detail); err != nil {
		return nil, false, err
	}
	return json.RawMessage(body), detail.Available, nil
}

func (c *ClaimService) ClaimProduct(userID uint, invProductID uint) response.Res {
	ctx := context.Background()

	snapshot, available, err := c.fetchInvProduct(invProductID)
	if err != nil {
		log.Printf("claim: vendor service error for inv_product %d: %v", invProductID, err)
		return response.Res{Message: "Product service unavailable, try again shortly", StatusCode: http.StatusServiceUnavailable}
	}
	if snapshot == nil {
		return c.ResponseNotFound("Product not found in inventory")
	}
	if !available {
		return c.ResponseBadRequest("Product is not available for claiming")
	}

	_, err = c.RepoClaim().GetByUserAndInvProduct(ctx, userID, invProductID)
	if err == nil {
		return c.ResponseConflict("You have already claimed this product")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return c.ResponseInternalServer("Failed to process claim")
	}

	claim := usermodel.NewClaim(userID, invProductID, snapshot)
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
	if err := c.RepoClaim().Delete(ctx, claimID, userID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.ResponseNotFound("Claim not found")
		}
		return c.ResponseInternalServer("Failed to remove claim")
	}
	return c.ResponseOK("Claim removed successfully", nil)
}
