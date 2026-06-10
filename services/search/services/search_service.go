package services

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/atharvyadav96k/spotnearr/pkg/httputil"
	"github.com/atharvyadav96k/spotnearr/pkg/tokenizer"
	searchpostgres "github.com/Developer-Aadesh/spotnearr-database/search/postgres"
	"gorm.io/gorm"
)

type SearchService struct {
	repo *searchpostgres.SearchRepository
}

func NewSearchService(db *gorm.DB) *SearchService {
	return &SearchService{repo: searchpostgres.NewSearchRepository(db)}
}

func (s *SearchService) Search(query string, lat, long *float64, rangeKm float64) httputil.Res {
	parsed := tokenizer.QueryTokenParse(query)
	if len(parsed.Tokens) == 0 {
		return httputil.NewResponse("no valid search terms", http.StatusBadRequest, nil)
	}

	rows, err := s.repo.Search(context.Background(), parsed.Tokens, lat, long, rangeKm)
	if err != nil {
		log.Printf("search: query error: %v", err)
		return httputil.NewResponse("search failed", http.StatusInternalServerError, nil)
	}

	ranked := rankProducts(rows, parsed, lat, long)
	if len(ranked) > 100 {
		ranked = ranked[:100]
	}

	return httputil.NewResponse("search results", http.StatusOK, ranked)
}

// ParseFloat64Param parses a float64 from a query param, returning nil if absent or invalid.
func ParseFloat64Param(val string) *float64 {
	if val == "" {
		return nil
	}
	f, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return nil
	}
	return &f
}
