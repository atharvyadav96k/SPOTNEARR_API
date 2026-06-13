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

type SpotlightService struct {
	baseService
}

func NewSpotlightService(db *gorm.DB, c *cache.Cache) *SpotlightService {
	return &SpotlightService{baseService: newBaseService(db, c)}
}

func (s *SpotlightService) PostSpotlight(bizID uint, dto pkgdtos.SpotlightCreateRequest) httputil.Res {
	spotlight := vendormodel.Spotlight{
		BusinessID: bizID,
		Type:       vendormodel.SpotlightType(dto.Type),
		Title:      dto.Title,
		Caption:    dto.Caption,
		MediaURL:   dto.MediaURL,
		MediaType:  vendormodel.SpotlightMediaType(dto.MediaType),
		ProductID:  dto.ProductID,
		OfferID:    dto.OfferID,
	}
	created, err := s.RepoSpotlight().Create(context.Background(), spotlight)
	if err != nil {
		return s.ResponseInternalServer("failed to post spotlight")
	}
	return s.ResponseCreated("spotlight posted", created)
}

func (s *SpotlightService) ListSpotlights(bizID uint) httputil.Res {
	spotlights, err := s.RepoSpotlight().GetByBusiness(context.Background(), bizID)
	if err != nil {
		return s.ResponseInternalServer("failed to fetch spotlights")
	}
	return s.ResponseOK("spotlights", spotlights)
}

func (s *SpotlightService) GetSpotlight(bizID uint, spotlightID uint) httputil.Res {
	spotlight, err := s.RepoSpotlight().GetByID(context.Background(), spotlightID, bizID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return s.ResponseNotFound("spotlight not found")
		}
		return s.ResponseInternalServer("failed to fetch spotlight")
	}
	return s.ResponseOK("spotlight", spotlight)
}

func (s *SpotlightService) DeleteSpotlight(bizID uint, spotlightID uint) httputil.Res {
	if err := s.RepoSpotlight().Delete(context.Background(), spotlightID, bizID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return s.ResponseNotFound("spotlight not found")
		}
		return s.ResponseInternalServer("failed to delete spotlight")
	}
	return s.ResponseOK("spotlight deleted", nil)
}
