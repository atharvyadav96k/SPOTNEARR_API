package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	pkgdtos "github.com/atharvyadav96k/spotnearr/pkg/dtos"
	pkghttputil "github.com/atharvyadav96k/spotnearr/pkg/httputil"
	"github.com/atharvyadav96k/spotnearr/vendor-svc/config"
	"github.com/atharvyadav96k/spotnearr/vendor-svc/connections/cache"
	vendorpostgres "github.com/Developer-Aadesh/spotnearr-database/vendordb/postgres"
	"gorm.io/gorm"
)

// ClaimSummary is the shape returned by the user-service internal claim endpoint.
type ClaimSummary struct {
	ID                 uint            `json:"id"`
	UserID             uint            `json:"userId"`
	InventoryProductID uint            `json:"inventoryProductId"`
	Status             string          `json:"status"`
	Product            json.RawMessage `json:"product,omitempty"`
	CreatedAt          time.Time       `json:"createdAt"`
	UpdatedAt          time.Time       `json:"updatedAt"`
}

type ClaimService struct {
	baseService
	invProductRepo *vendorpostgres.InvProductRepository
	httpClient     *http.Client
}

func NewClaimService(db *gorm.DB, c *cache.Cache) *ClaimService {
	return &ClaimService{
		baseService:    newBaseService(db, c),
		invProductRepo: vendorpostgres.NewInvProductRepository(db),
		httpClient:     &http.Client{Timeout: 2 * time.Second},
	}
}

// ListClaims fetches all incoming claims for the vendor's inventory products.
func (s *ClaimService) ListClaims(bizID uint) pkghttputil.Res {
	ids, err := s.invProductRepo.GetIDsByBusiness(context.Background(), bizID)
	if err != nil {
		log.Printf("claim-svc: get inv product ids for biz %d: %v", bizID, err)
		return s.ResponseInternalServer("failed to fetch inventory")
	}
	if len(ids) == 0 {
		return s.ResponseOK("no claims", []ClaimSummary{})
	}

	parts := make([]string, len(ids))
	for i, id := range ids {
		parts[i] = strconv.FormatUint(uint64(id), 10)
	}
	url := fmt.Sprintf("%s/internal/claims?product_ids=%s", config.C.UserServiceURL, strings.Join(parts, ","))

	resp, err := s.httpClient.Get(url)
	if err != nil {
		log.Printf("claim-svc: user-svc unreachable: %v", err)
		return s.ResponseInternalServer("user service unavailable")
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return s.ResponseInternalServer("failed to fetch claims")
	}

	var claims []ClaimSummary
	if err := json.Unmarshal(body, &claims); err != nil {
		return s.ResponseInternalServer("failed to parse claims")
	}
	return s.ResponseOK("claims", claims)
}

// UpdateClaimStatus accepts or rejects a claim after validating it belongs to the vendor.
func (s *ClaimService) UpdateClaimStatus(bizID uint, claimID uint, dto pkgdtos.ClaimStatusUpdateRequest) pkghttputil.Res {
	// Fetch the claim from user-service to validate ownership.
	claimURL := fmt.Sprintf("%s/internal/claims/%d", config.C.UserServiceURL, claimID)
	resp, err := s.httpClient.Get(claimURL)
	if err != nil {
		return s.ResponseInternalServer("user service unavailable")
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return s.ResponseNotFound("claim not found")
	}
	if resp.StatusCode != http.StatusOK {
		return s.ResponseInternalServer("failed to fetch claim")
	}

	body, _ := io.ReadAll(resp.Body)
	var claim ClaimSummary
	if err := json.Unmarshal(body, &claim); err != nil {
		return s.ResponseInternalServer("failed to parse claim")
	}

	// Verify the claimed inventory product belongs to this business.
	ids, err := s.invProductRepo.GetIDsByBusiness(context.Background(), bizID)
	if err != nil {
		return s.ResponseInternalServer("failed to validate ownership")
	}
	owned := false
	for _, id := range ids {
		if id == claim.InventoryProductID {
			owned = true
			break
		}
	}
	if !owned {
		return s.ResponseUnauthorized()
	}

	// Patch status on user-service.
	payload, _ := json.Marshal(map[string]string{"status": dto.Status})
	patchURL := fmt.Sprintf("%s/internal/claims/%d/status", config.C.UserServiceURL, claimID)
	req, _ := http.NewRequest(http.MethodPatch, patchURL, strings.NewReader(string(payload)))
	req.Header.Set("Content-Type", "application/json")

	patchResp, err := s.httpClient.Do(req)
	if err != nil {
		return s.ResponseInternalServer("user service unavailable")
	}
	defer patchResp.Body.Close()

	if patchResp.StatusCode != http.StatusOK {
		return s.ResponseInternalServer("failed to update claim status")
	}
	return s.ResponseOK("claim "+dto.Status, nil)
}
