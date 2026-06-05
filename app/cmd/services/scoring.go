package services

import (
	"math"
	"sort"

	"github.com/atharvyadav96k/SPOTNEARR_API/factories/token"
	"github.com/atharvyadav96k/SPOTNEARR_API/models"
)

const (
	scoreExactName     = 1000
	scoreCompleteName  = 800
	scoreCoverageMax   = 1000
	scoreCategoryMax   = 300
	scorePhraseExact   = 200
	scorePhrasePartial = 50
	scoreDistanceMax   = 100
)

type rankedProduct struct {
	product    models.ProductResult
	score      int
	coverage   float64
	distanceKm float64
}

func rankProducts(
	products []models.ProductResult,
	parsed token.ParseResult,
	categoryFreqs map[uint]int,
	productCategories map[uint][]uint,
	lat, long *float64,
) []models.ProductResult {
	ranked := make([]rankedProduct, len(products))
	for i, p := range products {
		ranked[i] = scoreProduct(p, parsed, categoryFreqs, productCategories[p.ID], lat, long)
	}

	sort.SliceStable(ranked, func(i, j int) bool {
		a, b := ranked[i], ranked[j]
		if a.score != b.score {
			return a.score > b.score
		}
		if a.coverage != b.coverage {
			return a.coverage > b.coverage
		}
		return a.distanceKm < b.distanceKm
	})

	result := make([]models.ProductResult, len(ranked))
	for i, r := range ranked {
		r.product.DistanceKm = r.distanceKm
		result[i] = r.product
	}
	return result
}

func scoreProduct(
	p models.ProductResult,
	parsed token.ParseResult,
	categoryFreqs map[uint]int,
	productCategoryIDs []uint,
	lat, long *float64,
) rankedProduct {
	score := 0
	coverage := 0.0

	nameTokens := token.TokenParser(p.Name)
	nameSet := toSet(nameTokens)
	queryTokens := parsed.Tokens

	if len(queryTokens) > 0 {
		// Factor 1 — Coverage: fraction of query tokens found in product name × 1000
		nameHits := 0
		for _, qt := range queryTokens {
			if _, ok := nameSet[qt]; ok {
				nameHits++
			}
		}
		coverage = float64(nameHits) / float64(len(queryTokens))
		score += int(coverage * float64(scoreCoverageMax))

		// Factor 2 — Exact Match
		if len(nameTokens) == len(queryTokens) && allInSet(queryTokens, nameSet) {
			score += scoreExactName
		} else if allInSet(queryTokens, nameSet) {
			score += scoreCompleteName
		}

		// Factor 4 — Phrase Order
		if isExactPhrase(nameTokens, queryTokens) {
			score += scorePhraseExact
		} else if countConsecutivePairs(nameTokens, queryTokens) > 0 {
			score += scorePhrasePartial
		}
	}

	// Factor 3 — Category Intent
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

	// Factor 5 — Distance
	var distKm float64
	if lat != nil && long != nil && (p.Lat != 0 || p.Long != 0) {
		distKm = haversineKm(*lat, *long, p.Lat, p.Long)
		bonus := scoreDistanceMax - int(distKm)
		if bonus > 0 {
			score += bonus
		}
	}

	return rankedProduct{
		product:    p,
		score:      score,
		coverage:   coverage,
		distanceKm: distKm,
	}
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
			found := false
			for _, pb := range positions[b] {
				if pb == pa+1 {
					found = true
					break
				}
			}
			if found {
				count++
				break
			}
		}
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
