package services

import (
	"math"
	"sort"

	"github.com/atharvyadav96k/spotnearr/pkg/tokenizer"
	"github.com/atharvyadav96k/spotnearr/search-svc/models"
	"github.com/atharvyadav96k/spotnearr/search-svc/repository"
)

const (
	scoreExactName     = 1000
	scoreCompleteName  = 800
	scoreCoverageMax   = 1000
	scoreCategoryMax   = 300
	scorePhraseExact   = 200
	scorePhrasePartial = 50

	bucketKm0 = 0.4
	bucketKm1 = 0.8
	bucketKm2 = 2.0
	bucketKm3 = 5.0

	noGeoTier = math.MaxInt32

	diversityMaxPerStore = 5
	diversityMinRatio    = 0.60
)

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

type rankedProduct struct {
	product    models.ProductResult
	score      int
	coverage   float64
	distanceKm float64
	tier       int
}

// buildCategoryFreqs approximates the monolith's product_tokens frequency table by
// distributing each search row's token_match_cnt across its category_ids.
func buildCategoryFreqs(rows []repository.SearchRow) map[uint]int {
	freqs := make(map[uint]int)
	for _, r := range rows {
		for _, catID := range r.CategoryIDs {
			freqs[catID] += r.TokenMatchCnt
		}
	}
	return freqs
}

func rankProducts(
	rows []repository.SearchRow,
	parsed tokenizer.ParseResult,
	lat, long *float64,
) []models.ProductResult {
	categoryFreqs := buildCategoryFreqs(rows)

	ranked := make([]rankedProduct, len(rows))
	for i, r := range rows {
		pr := models.ProductResult{
			ID:            r.ID,
			ProductID:     r.ProductID,
			BusinessID:    r.BusinessID,
			ProductName:   r.ProductName,
			Price:         r.Price,
			PriceUnit:     r.PriceUnit,
			Quantity:      r.Quantity,
			QuantityUnit:  r.QuantityUnit,
			Description:   r.Description,
			StoreID:       r.StoreID,
			StoreName:     r.StoreName,
			StreetAddress: r.StreetAddress,
			Lat:           r.Lat,
			Long:          r.Long,
			GeoHash:       r.GeoHash,
		}
		ranked[i] = scoreProduct(pr, parsed, categoryFreqs, r.CategoryIDs, lat, long)
	}

	sort.SliceStable(ranked, func(i, j int) bool {
		a, b := ranked[i], ranked[j]
		if a.tier != b.tier {
			return a.tier < b.tier
		}
		if a.score != b.score {
			return a.score > b.score
		}
		if a.coverage != b.coverage {
			return a.coverage > b.coverage
		}
		return a.distanceKm < b.distanceKm
	})

	ranked = diversifyResults(ranked)

	result := make([]models.ProductResult, len(ranked))
	for i, r := range ranked {
		r.product.DistanceKm = r.distanceKm
		result[i] = r.product
	}
	return result
}

func diversifyResults(ranked []rankedProduct) []rankedProduct {
	if len(ranked) <= 1 {
		return ranked
	}
	bestScore := ranked[0].score
	threshold := int(float64(bestScore) * diversityMinRatio)
	promoted := make([]rankedProduct, 0, len(ranked))
	deferred := make([]rankedProduct, 0)
	storeCount := make(map[uint]int)

	for _, r := range ranked {
		sid := r.product.StoreID
		count := storeCount[sid]
		switch {
		case count == 0:
			promoted = append(promoted, r)
			storeCount[sid]++
		case count < diversityMaxPerStore && r.score >= threshold:
			promoted = append(promoted, r)
			storeCount[sid]++
		default:
			deferred = append(deferred, r)
		}
	}
	return append(promoted, deferred...)
}

func scoreProduct(
	p models.ProductResult,
	parsed tokenizer.ParseResult,
	categoryFreqs map[uint]int,
	productCategoryIDs []uint,
	lat, long *float64,
) rankedProduct {
	score := 0
	coverage := 0.0

	nameTokens := tokenizer.TokenParser(p.ProductName)
	nameSet := toSet(nameTokens)
	queryTokens := parsed.Tokens

	if len(queryTokens) > 0 {
		nameHits := 0
		for _, qt := range queryTokens {
			if _, ok := nameSet[qt]; ok {
				nameHits++
			}
		}
		coverage = float64(nameHits) / float64(len(queryTokens))
		score += int(coverage * float64(scoreCoverageMax))

		if len(nameTokens) == len(queryTokens) && allInSet(queryTokens, nameSet) {
			score += scoreExactName
		} else if allInSet(queryTokens, nameSet) {
			score += scoreCompleteName
		}

		if isExactPhrase(nameTokens, queryTokens) {
			score += scorePhraseExact
		} else if pairs := countConsecutivePairs(nameTokens, queryTokens); pairs > 0 {
			score += pairs * scorePhrasePartial
		}
	}

	if len(categoryFreqs) > 0 && len(productCategoryIDs) > 0 {
		maxFreq := 0
		for _, v := range categoryFreqs {
			if v > maxFreq {
				maxFreq = v
			}
		}
		if maxFreq > 0 {
			best := 0
			for _, catID := range productCategoryIDs {
				if freq, ok := categoryFreqs[catID]; ok {
					v := int(float64(freq) / float64(maxFreq) * float64(scoreCategoryMax))
					if v > best {
						best = v
					}
				}
			}
			score += best
		}
	}

	var distKm float64
	tier := noGeoTier
	if lat != nil && long != nil && (p.Lat != 0 || p.Long != 0) {
		distKm = haversineKm(*lat, *long, p.Lat, p.Long)
		tier = distanceTierOf(distKm)
	}

	return rankedProduct{product: p, score: score, coverage: coverage, distanceKm: distKm, tier: tier}
}

func isExactPhrase(nameTokens, queryTokens []string) bool {
	if len(queryTokens) == 0 || len(queryTokens) > len(nameTokens) {
		return false
	}
	for i := 0; i <= len(nameTokens)-len(queryTokens); i++ {
		match := true
		for j, qt := range queryTokens {
			if nameTokens[i+j] != qt {
				match = false
				break
			}
		}
		if match {
			return true
		}
	}
	return false
}

func countConsecutivePairs(nameTokens, queryTokens []string) int {
	if len(queryTokens) < 2 {
		return 0
	}
	positions := make(map[string][]int, len(nameTokens))
	for i, t := range nameTokens {
		positions[t] = append(positions[t], i)
	}
	count := 0
	for i := 0; i < len(queryTokens)-1; i++ {
		a, b := queryTokens[i], queryTokens[i+1]
		for _, pa := range positions[a] {
			for _, pb := range positions[b] {
				if pb == pa+1 {
					count++
					goto nextPair
				}
			}
		}
	nextPair:
	}
	return count
}

func allInSet(tokens []string, set map[string]struct{}) bool {
	for _, t := range tokens {
		if _, ok := set[t]; !ok {
			return false
		}
	}
	return true
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

func toSet(tokens []string) map[string]struct{} {
	s := make(map[string]struct{}, len(tokens))
	for _, t := range tokens {
		s[t] = struct{}{}
	}
	return s
}
