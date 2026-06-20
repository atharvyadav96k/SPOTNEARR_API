package services

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/atharvyadav96k/spotnearr/pkg/httputil"
	"github.com/atharvyadav96k/spotnearr/pkg/tokenizer"
	"github.com/atharvyadav96k/spotnearr/search-svc/cache"
	tsclient "github.com/atharvyadav96k/spotnearr/search-svc/typesense"
	"github.com/typesense/typesense-go/v2/typesense"
	"github.com/typesense/typesense-go/v2/typesense/api"
)

const searchTimeout = 5 * time.Second

// SearchFilters holds optional query-level filters.
type SearchFilters struct {
	CategoryIDs []uint
	MinPrice    *float64
	MaxPrice    *float64
}

// SearchResult is the per-item shape returned to the client.
type SearchResult struct {
	InvProductID int64   `json:"inv_product_id"`
	ProductID    int64   `json:"product_id"`
	Name         string  `json:"name"`
	Price        float64 `json:"price"`
	PriceUnit    string  `json:"price_unit"`
	Lat          float64 `json:"lat"`
	Long         float64 `json:"long"`
	DistanceKm   float64 `json:"distance_km,omitempty"`
}

type SearchService struct {
	client  *typesense.Client
	tcCache *cache.TokenCategoryCache
}

func NewSearchService(client *typesense.Client, tcCache *cache.TokenCategoryCache) *SearchService {
	return &SearchService{client: client, tcCache: tcCache}
}

func (s *SearchService) Search(ctx context.Context, query string, lat, long *float64, rangeKm float64, filters SearchFilters) httputil.Res {
	if strings.TrimSpace(query) == "" {
		return httputil.NewResponse("no valid search terms", http.StatusBadRequest, nil)
	}

	ctx, cancel := context.WithTimeout(ctx, searchTimeout)
	defer cancel()

	// ── Filters ────────────────────────────────────────────────────────────
	filterParts := []string{"available:=true"}
	if lat != nil && long != nil && rangeKm > 0 {
		filterParts = append(filterParts, fmt.Sprintf("location:(%v,%v,%v km)", *lat, *long, rangeKm))
	}
	if len(filters.CategoryIDs) > 0 {
		ids := make([]string, len(filters.CategoryIDs))
		for i, id := range filters.CategoryIDs {
			ids[i] = fmt.Sprintf("%d", id)
		}
		filterParts = append(filterParts, "category_ids:["+strings.Join(ids, ",")+"]")
	}
	if filters.MinPrice != nil {
		filterParts = append(filterParts, fmt.Sprintf("price:>=%v", *filters.MinPrice))
	}
	if filters.MaxPrice != nil {
		filterParts = append(filterParts, fmt.Sprintf("price:<=%v", *filters.MaxPrice))
	}
	filterBy := strings.Join(filterParts, " && ")

	// ── Sort: text match → category affinity boost → geo distance ──────────
	tokens := tokenizer.TokenParser(query)
	topCats := s.tcCache.GetTopCategories(ctx, tokens, 3)

	sortBy := "_text_match:desc"
	if len(topCats) > 0 {
		catStrs := make([]string, len(topCats))
		for i, id := range topCats {
			catStrs[i] = fmt.Sprintf("%d", id)
		}
		sortBy += ",_eval(category_ids:[" + strings.Join(catStrs, ",") + "]):desc"
	}
	if lat != nil && long != nil {
		sortBy += fmt.Sprintf(",location(%v,%v):asc", *lat, *long)
	}

	queryBy := "name"
	perPage := 100
	searchRes, err := s.client.Collection(tsclient.CollectionName).Documents().Search(ctx, &api.SearchCollectionParams{
		Q:        &query,
		QueryBy:  &queryBy,
		FilterBy: &filterBy,
		SortBy:   &sortBy,
		PerPage:  &perPage,
	})
	if err != nil {
		if ctx.Err() != nil {
			log.Printf("search: timeout after %s: %v", searchTimeout, err)
			return httputil.NewResponse("search timed out", http.StatusGatewayTimeout, nil)
		}
		log.Printf("search: typesense error: %v", err)
		return httputil.NewResponse("search failed", http.StatusInternalServerError, nil)
	}

	results := hitsToResults(searchRes.Hits)
	return httputil.NewResponse("search results", http.StatusOK, results)
}

func hitsToResults(hits *[]api.SearchResultHit) []SearchResult {
	if hits == nil || len(*hits) == 0 {
		return []SearchResult{}
	}
	results := make([]SearchResult, 0, len(*hits))
	for _, hit := range *hits {
		if hit.Document == nil {
			continue
		}
		doc := *hit.Document
		r := SearchResult{
			InvProductID: int64Field(doc, "inv_product_id"),
			ProductID:    int64Field(doc, "product_id"),
			Name:         stringField(doc, "name"),
			Price:        floatField(doc, "price"),
			PriceUnit:    stringField(doc, "price_unit"),
			Lat:          geoField(doc, "location", 0),
			Long:         geoField(doc, "location", 1),
		}
		if hit.GeoDistanceMeters != nil {
			if dist, ok := (*hit.GeoDistanceMeters)["location"]; ok {
				r.DistanceKm = float64(dist) / 1000.0
			}
		}
		results = append(results, r)
	}
	return results
}

func int64Field(doc map[string]interface{}, key string) int64 {
	if v, ok := doc[key]; ok {
		if n, ok := v.(float64); ok {
			return int64(n)
		}
	}
	return 0
}

func floatField(doc map[string]interface{}, key string) float64 {
	if v, ok := doc[key]; ok {
		if n, ok := v.(float64); ok {
			return n
		}
	}
	return 0
}

func stringField(doc map[string]interface{}, key string) string {
	if v, ok := doc[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

func geoField(doc map[string]interface{}, key string, idx int) float64 {
	if v, ok := doc[key]; ok {
		if arr, ok := v.([]interface{}); ok && idx < len(arr) {
			if f, ok := arr[idx].(float64); ok {
				return f
			}
		}
	}
	return 0
}
