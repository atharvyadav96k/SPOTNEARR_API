package services

import (
	"context"
	"log"

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
	log.Default().Println(parsed)
	ctx := context.Background()

	products, err := s.RepoInvProduct().SearchProduct(ctx, parsed.Tokens, lat, long, rangeKm)
	if err != nil {
		return s.ResponseInternalServer("search failed")
	}
	log.Default().Println(products)
	categoryFreqs, _ := s.RepoProductToken().GetTokenCategoryFreqs(ctx, parsed.Tokens)

	productIDs := make([]uint, len(products))
	for i, p := range products {
		productIDs[i] = p.ID
	}
	productCategories, _ := s.RepoProduct().GetProductCategoryIDs(ctx, productIDs)

	ranked := rankProducts(products, parsed, categoryFreqs, productCategories, lat, long)
	if len(ranked) > 100 {
		ranked = ranked[:100]
	}
	return s.ResponseOK("search results", ranked)
}
