package services

import (
	"context"
	"log"
	"net/http"
	"time"

	searchpostgres "github.com/Developer-Aadesh/spotnearr-database/search/postgres"
	"github.com/atharvyadav96k/spotnearr/pkg/httputil"
	"github.com/atharvyadav96k/spotnearr/pkg/tokenizer"
	searchcache "github.com/atharvyadav96k/spotnearr/search-svc/cache"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

const searchTimeout = 5 * time.Second

type SearchService struct {
	repo  *searchpostgres.SearchRepository
	cache *searchcache.SearchCache
}

func NewSearchService(db *gorm.DB, rdb *redis.Client) *SearchService {
	return &SearchService{
		repo:  searchpostgres.NewSearchRepository(db),
		cache: searchcache.New(rdb),
	}
}

// Search tokenizes the query, loads category frequency data (Redis-cached), queries
// the index, and returns a ranked list of product cards capped at 100 results.
func (s *SearchService) Search(ctx context.Context, query string, lat, long *float64, rangeKm float64) httputil.Res {
	parsed := tokenizer.QueryTokenParse(query)
	if len(parsed.Tokens) == 0 {
		return httputil.NewResponse("no valid search terms", http.StatusBadRequest, nil)
	}

	ctx, cancel := context.WithTimeout(ctx, searchTimeout)
	defer cancel()

	freqs, ok := s.cache.GetTokenFreqs(ctx, parsed.Tokens)
	if !ok {
		var err error
		freqs, err = s.repo.GetTokenCategoryFreqs(ctx, parsed.Tokens)
		if err != nil {
			log.Printf("search: load freqs: %v", err)
			// Non-fatal: category boost degrades to zero; token coverage still ranks results.
		} else {
			s.cache.SetTokenFreqs(ctx, parsed.Tokens, freqs)
		}
	}

	rows, err := s.repo.Search(ctx, parsed.Tokens, lat, long, rangeKm)
	if err != nil {
		if ctx.Err() != nil {
			log.Printf("search: timeout after %s: %v", searchTimeout, err)
			return httputil.NewResponse("search timed out", http.StatusGatewayTimeout, nil)
		}
		log.Printf("search: query error: %v", err)
		return httputil.NewResponse("search failed", http.StatusInternalServerError, nil)
	}

	results := rankProducts(rows, parsed.Tokens, freqs, lat, long)
	return httputil.NewResponse("search results", http.StatusOK, results)
}
