package services

import (
	"context"
	"errors"
	"log"

	vendormodel "github.com/Developer-Aadesh/spotnearr-database/vendordb"
	pkgdtos "github.com/atharvyadav96k/spotnearr/pkg/dtos"
	"github.com/atharvyadav96k/spotnearr/pkg/httputil"
	"github.com/atharvyadav96k/spotnearr/pkg/mq"
	"github.com/atharvyadav96k/spotnearr/vendor-svc/connections/cache"
	"github.com/mmcloughlin/geohash"
	"gorm.io/gorm"
)

type OfferService struct {
	baseService
	pub *mq.Publisher
}

func NewOfferService(db *gorm.DB, c *cache.Cache, pub *mq.Publisher) *OfferService {
	return &OfferService{baseService: newBaseService(db, c), pub: pub}
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

	go s.publishDealSync(bizID, created)

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
	offer, err := s.RepoOffer().GetByID(context.Background(), offerID, bizID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return s.ResponseNotFound("offer not found")
		}
		return s.ResponseInternalServer("failed to fetch offer")
	}
	if err := s.RepoOffer().Delete(context.Background(), offerID, bizID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return s.ResponseNotFound("offer not found")
		}
		return s.ResponseInternalServer("failed to delete offer")
	}
	go func() {
		if err := s.pub.Publish(context.Background(), mq.TopicDealDelete, mq.DealDeletePayload{
			DealID: offer.ID,
		}); err != nil {
			log.Printf("offer: publish deal.delete %d: %v", offer.ID, err)
		}
	}()
	return s.ResponseOK("offer deleted", nil)
}

func (s *OfferService) publishDealSync(bizID uint, offer vendormodel.Offer) {
	stores, err := s.RepoStore().GetStoreByBusinessId(context.Background(), bizID)
	if err != nil {
		log.Printf("offer: deal sync get stores biz %d: %v", bizID, err)
		return
	}
	seen := map[string]struct{}{}
	locs := make([]mq.DealStoreLocation, 0, len(stores))
	for _, st := range stores {
		gh5 := geohash.EncodeWithPrecision(st.Lat, st.Long, 5)
		if _, ok := seen[gh5]; ok {
			continue
		}
		seen[gh5] = struct{}{}
		locs = append(locs, mq.DealStoreLocation{StoreID: st.ID, GeoHash5: gh5})
	}
	payload := mq.DealSyncPayload{
		DealID:        offer.ID,
		BusinessID:    offer.BusinessID,
		Title:         offer.Title,
		DiscountType:  string(offer.DiscountType),
		DiscountValue: offer.DiscountValue,
		MinOrderValue: offer.MinOrderValue,
		Active:        offer.Active,
		ExpiresAt:     offer.ExpiresAt,
		Locations:     locs,
	}
	if err := s.pub.Publish(context.Background(), mq.TopicDealSync, payload); err != nil {
		log.Printf("offer: publish deal.sync %d: %v", offer.ID, err)
	}
}
