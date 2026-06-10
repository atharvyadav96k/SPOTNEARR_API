package services

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/atharvyadav96k/spotnearr/pkg/httputil"
	"github.com/atharvyadav96k/spotnearr/pkg/jwtutil"
	"github.com/atharvyadav96k/spotnearr/vendor-svc/connections/cache"
	"github.com/atharvyadav96k/spotnearr/vendor-svc/dtos"
	"github.com/atharvyadav96k/spotnearr/vendor-svc/models"
	"gorm.io/gorm"
)

type BusinessService struct {
	baseService
	userServiceURL string
	httpClient     *http.Client
}

func NewBusinessService(db *gorm.DB, c *cache.Cache, userServiceURL string) *BusinessService {
	return &BusinessService{
		baseService:    newBaseService(db, c),
		userServiceURL: userServiceURL,
		httpClient:     &http.Client{Timeout: 3 * time.Second},
	}
}

func (b *BusinessService) RegisterBusiness(dto *dtos.RegisterBusiness, userID uint) httputil.Res {
	biz, err := b.RepoBusiness().Create(context.Background(), models.NewBusiness(dto.Name, dto.Email, dto.Phone, dto.Desc))
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return b.ResponseConflict("Phone number or email is already in use")
		}
		if errors.Is(err, gorm.ErrForeignKeyViolated) {
			return b.ResponseConflict("Invalid user account")
		}
		return b.ResponseBadRequest(err.Error())
	}

	if err := b.RepoAccess().CreateNewAccess(context.Background(),
		models.NewBusinessAccess(userID, biz.ID, jwtutil.RoleAdmin)); err != nil {
		log.Default().Println(err)
		b.RepoBusiness().HardDelete(context.Background(), biz.ID)
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return b.ResponseConflict("This account is already associated with a business")
		}
		return b.ResponseBadRequest(err.Error())
	}

	// Fire-and-forget: ask User Service to drop the refresh token so the next login
	// generates a JWT with business_id. Non-fatal if it fails.
	go b.invalidateRefreshToken(userID)

	return b.ResponseCreated("Business successfully registered", biz)
}

func (b *BusinessService) invalidateRefreshToken(userID uint) {
	url := fmt.Sprintf("%s/internal/users/%d/invalidate-refresh", b.userServiceURL, userID)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(nil))
	if err != nil {
		log.Printf("business: build invalidate-refresh request: %v", err)
		return
	}
	resp, err := b.httpClient.Do(req)
	if err != nil {
		log.Printf("business: invalidate-refresh call failed (user %d): %v", userID, err)
		return
	}
	resp.Body.Close()
}

func (b *BusinessService) GetBusinessByID(id uint) httputil.Res {
	business, err := b.RepoBusiness().GetByID(context.Background(), id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return b.ResponseNotFound("Business not found")
		}
		return b.ResponseInternalServer("Failed to get business")
	}
	return b.ResponseOK("Business info", business)
}

func (b *BusinessService) UpdateBusiness(businessID uint, dto *dtos.UpdateBusiness) httputil.Res {
	updated, err := b.RepoBusiness().Update(context.Background(), businessID, dto.Name, dto.Desc)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return b.ResponseNotFound("Business not found")
		}
		return b.ResponseBadRequest("Failed to update the business")
	}
	return b.ResponseOK("Business updated successfully", updated)
}
