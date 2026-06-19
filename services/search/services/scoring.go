package services

import (
	"math"
	"sort"

	searchdb "github.com/Developer-Aadesh/spotnearr-database/search"
)

const (
	scoreCoverageMax = 1000
	scoreCategoryMax = 300

	bucketKm0 = 0.4
	bucketKm1 = 0.8
	bucketKm2 = 2.0
	bucketKm3 = 5.0

	noGeoTier  = math.MaxInt32
	maxResults = 100
)

// SearchResult is the per-item shape returned to the client.
// Full product detail (store address, availability, categories) is available
// via GET /api/v1/products/{invProductId}/detail on the vendor service.
type SearchResult struct {
	InvProductID uint    `json:"inv_product_id"`
	ProductID    uint    `json:"product_id"`
	Name         string  `json:"name"`
	Price        float64 `json:"price"`
	PriceUnit    string  `json:"price_unit"`
	Lat          float64 `json:"lat"`
	Long         float64 `json:"long"`
	DistanceKm   float64 `json:"distance_km,omitempty"`
}

type rankedEntry struct {
	row        searchdb.SearchRow
	score      int
	distanceKm float64
	tier       int
}

func distanceTierOf(km float64) int {
	switch {
	case km < bucketKm0:
		return 0
	case km < bucketKm1:
		return 1
	case km < bucketKm2:
		return 2
	case km < bucketKm3:
		return 3
	default:
		return 4
	}
}

// buildCategoryScores sums stored frequency counts per category for the query tokens.
func buildCategoryScores(queryTokens []string, freqs []searchdb.TokenCategoryFreq) map[uint]int64 {
	tokenSet := make(map[string]struct{}, len(queryTokens))
	for _, t := range queryTokens {
		tokenSet[t] = struct{}{}
	}
	scores := make(map[uint]int64)
	for _, f := range freqs {
		if _, ok := tokenSet[f.Token]; ok {
			scores[f.CategoryID] += f.Count
		}
	}
	return scores
}

// rankProducts scores every search row, deduplicates by ProductID (keeping the
// highest-scored entry per product), sorts, caps at 100, and returns SearchResults.
func rankProducts(
	rows []searchdb.SearchRow,
	queryTokens []string,
	freqs []searchdb.TokenCategoryFreq,
	lat, long *float64,
) []SearchResult {
	if len(rows) == 0 {
		return []SearchResult{}
	}

	catScores := buildCategoryScores(queryTokens, freqs)
	var maxCatScore int64
	for _, v := range catScores {
		if v > maxCatScore {
			maxCatScore = v
		}
	}

	best := make(map[uint]rankedEntry, len(rows))
	for _, row := range rows {
		e := scoreEntry(row, len(queryTokens), catScores, maxCatScore, lat, long)
		if prev, ok := best[e.row.ProductID]; !ok || isBetter(e, prev) {
			best[e.row.ProductID] = e
		}
	}

	ranked := make([]rankedEntry, 0, len(best))
	for _, e := range best {
		ranked = append(ranked, e)
	}
	sort.SliceStable(ranked, func(i, j int) bool {
		a, b := ranked[i], ranked[j]
		if a.tier != b.tier {
			return a.tier < b.tier
		}
		if a.score != b.score {
			return a.score > b.score
		}
		return a.distanceKm < b.distanceKm
	})

	if len(ranked) > maxResults {
		ranked = ranked[:maxResults]
	}

	results := make([]SearchResult, len(ranked))
	for i, e := range ranked {
		results[i] = SearchResult{
			InvProductID: e.row.ID,
			ProductID:    e.row.ProductID,
			Name:         e.row.Name,
			Price:        e.row.Price,
			PriceUnit:    e.row.PriceUnit,
			Lat:          e.row.Lat,
			Long:         e.row.Long,
			DistanceKm:   e.distanceKm,
		}
	}
	return results
}

func isBetter(a, b rankedEntry) bool {
	if a.tier != b.tier {
		return a.tier < b.tier
	}
	if a.score != b.score {
		return a.score > b.score
	}
	return a.distanceKm < b.distanceKm
}

func scoreEntry(
	row searchdb.SearchRow,
	queryTokenCount int,
	catScores map[uint]int64,
	maxCatScore int64,
	lat, long *float64,
) rankedEntry {
	score := 0

	if queryTokenCount > 0 {
		coverage := float64(row.TokenMatchCnt) / float64(queryTokenCount)
		score += int(coverage * float64(scoreCoverageMax))
	}

	if maxCatScore > 0 {
		var bestCat int64
		for _, catID := range row.CategoryIDs {
			if s := catScores[catID]; s > bestCat {
				bestCat = s
			}
		}
		score += int(float64(bestCat) / float64(maxCatScore) * float64(scoreCategoryMax))
	}

	var distKm float64
	tier := noGeoTier
	if lat != nil && long != nil && (row.Lat != 0 || row.Long != 0) {
		distKm = haversineKm(*lat, *long, row.Lat, row.Long)
		tier = distanceTierOf(distKm)
	}

	return rankedEntry{row: row, score: score, distanceKm: distKm, tier: tier}
}

func haversineKm(lat1, lon1, lat2, lon2 float64) float64 {
	const R = 6371.0
	dLat := (lat2 - lat1) * math.Pi / 180
	dLon := (lon2 - lon1) * math.Pi / 180
	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1*math.Pi/180)*math.Cos(lat2*math.Pi/180)*
			math.Sin(dLon/2)*math.Sin(dLon/2)
	return R * 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
}
