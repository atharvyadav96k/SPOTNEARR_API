package services

import (
	"context"
	"errors"

	vendormodel "github.com/Developer-Aadesh/spotnearr-database/vendordb"
	pkgdtos "github.com/atharvyadav96k/spotnearr/pkg/dtos"
	"github.com/atharvyadav96k/spotnearr/pkg/httputil"
	"github.com/atharvyadav96k/spotnearr/vendor-svc/connections/cache"
	"gorm.io/gorm"
)

type OfferService struct {
	baseService
}

func NewOfferService(db *gorm.DB, c *cache.Cache) *OfferService {
	return &OfferService{baseService: newBaseService(db, c)}
}

func (s *OfferService) CreateOffer(bizID uint, dto pkgdtos.OfferCreateRequest) httputil.Res {
	offer := vendormodel.Offer{
		BusinessID:    bizID,
		Title:         dto.Title,
		Description:   dto.Description,
		DiscountType:  vendormodel.DiscountType(dto.DiscountType),
		DiscountValue: dto.DiscountValue,
		MinOrderValue: dto.MinOrderValue,
		Code:          dto.Code,
		MaxUsage:      dto.MaxUsage,
		ExpiresAt:     dto.ExpiresAt,
		Active:        true,
	}
	created, err := s.RepoOffer().Create(context.Background(), offer)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return s.ResponseConflict("coupon code already in use")
		}
		return s.ResponseInternalServer("failed to create offer")
	}
	return s.ResponseCreated("offer created", created)
}

func (s *OfferService) ListOffers(bizID uint) httputil.Res {
	offers, err := s.RepoOffer().GetByBusiness(context.Background(), bizID)
	if err != nil {
		return s.ResponseInternalServer("failed to fetch offers")
	}
	return s.ResponseOK("offers", offers)
}

func (s *OfferService) GetOffer(bizID uint, offerID uint) httputil.Res {
	offer, err := s.RepoOffer().GetByID(context.Background(), offerID, bizID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return s.ResponseNotFound("offer not found")
		}
		return s.ResponseInternalServer("failed to fetch offer")
	}
	return s.ResponseOK("offer", offer)
}

func (s *OfferService) UpdateOffer(bizID uint, offerID uint, dto pkgdtos.OfferUpdateRequest) httputil.Res {
	offer := vendormodel.Offer{
		ID:            offerID,
		Title:         dto.Title,
		Description:   dto.Description,
		DiscountType:  vendormodel.DiscountType(dto.DiscountType),
		DiscountValue: dto.DiscountValue,
		MinOrderValue: dto.MinOrderValue,
		Code:          dto.Code,
		MaxUsage:      dto.MaxUsage,
		ExpiresAt:     dto.ExpiresAt,
	}
	updated, err := s.RepoOffer().Update(context.Background(), offer, bizID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return s.ResponseNotFound("offer not found")
		}
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return s.ResponseConflict("coupon code already in use")
		}
		return s.ResponseInternalServer("failed to update offer")
	}
	return s.ResponseOK("offer updated", updated)
}

func (s *OfferService) DeleteOffer(bizID uint, offerID uint) httputil.Res {
	if err := s.RepoOffer().Delete(context.Background(), offerID, bizID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return s.ResponseNotFound("offer not found")
		}
		return s.ResponseInternalServer("failed to delete offer")
	}
	return s.ResponseOK("offer deleted", nil)
}
