package services

import (
	"context"

	"github.com/atharvyadav96k/SPOTNEARR_API/connections/cache"
	"github.com/atharvyadav96k/SPOTNEARR_API/factories/token"
	"github.com/atharvyadav96k/SPOTNEARR_API/utils/response"
	"gorm.io/gorm"
)

type SearchService struct {
	base_service
}

func NewSearchService(db *gorm.DB, cache *cache.Cache) *SearchService {
	return &SearchService{
		base_service: NewBaseService(db, cache),
	}
}

func (s *SearchService) Search(query string, lat, long *float64, rangeKm float64) response.Res {
	parsed := token.QueryTokenParse(query)
	if len(parsed.Tokens) == 0 {
		return s.ResponseBadRequest("no valid search terms")
	}
	products, err := s.RepoInvProduct().SearchProduct(context.Background(), parsed.Tokens, lat, long, rangeKm)
	if err != nil {
		return s.ResponseInternalServer("search failed")
	}
	return s.ResponseOK("search results", rankProducts(products, parsed, lat, long))
}
