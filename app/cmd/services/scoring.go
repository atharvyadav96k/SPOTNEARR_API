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

	// Distance bucket boundaries in km.
	// All products within the same bucket are treated as equally proximate;
	// ranking within a bucket is determined purely by semantic score.
	bucketKm0 = 0.4 // tier 0: 0–400 m
	bucketKm1 = 0.8 // tier 1: 400–800 m
	bucketKm2 = 2.0 // tier 2: 800 m–2 km
	bucketKm3 = 5.0 // tier 3: 2–5 km
	// tier 4: 5 km+

	noGeoTier = math.MaxInt32 // assigned when no coordinate data is available

	// Store diversity — prevent a single merchant from monopolising top results.
	// A store's first product is always promoted. Each subsequent product from the
	// same store is only promoted if it scores above diversityMinRatio × bestScore
	// and the store hasn't yet filled its per-store cap.
	diversityMaxPerStore = 5    // max products per store in the promoted window
	diversityMinRatio    = 0.60 // n-th product (n>1) must reach this fraction of top score
)

// distanceTierOf maps a Haversine distance (km) to a discrete tier index.
// Lower tier = closer bucket = higher priority in sort.
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
		if a.tier != b.tier {
			return a.tier < b.tier // closer bucket always wins
		}
		if a.score != b.score {
			return a.score > b.score // better semantic match wins within same bucket
		}
		if a.coverage != b.coverage {
			return a.coverage > b.coverage
		}
		return a.distanceKm < b.distanceKm // exact meters as final tiebreaker
	})

	ranked = diversifyResults(ranked)

	result := make([]models.ProductResult, len(ranked))
	for i, r := range ranked {
		r.product.DistanceKm = r.distanceKm
		result[i] = r.product
	}
	return result
}

// diversifyResults reorders ranked products so that no single store occupies
// more than diversityMaxPerStore slots in the promoted window.
//
// Walk order (tier → score) is preserved: each store's first product is always
// promoted regardless of score. A store's subsequent products are promoted only
// if their score is at least diversityMinRatio × bestScore AND the store has not
// yet filled its cap. Everything else is deferred and appended afterwards in their
// original relative order.
//
// Fallback: if only one store has relevant products, all its products end up in
// the promoted window naturally — no diversity penalty is applied.
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
			// First product from this store — always include for diversity.
			promoted = append(promoted, r)
			storeCount[sid]++
		case count < diversityMaxPerStore && r.score >= threshold:
			// Within per-store cap and still relevant — include.
			promoted = append(promoted, r)
			storeCount[sid]++
		default:
			// Store at cap or product below relevance threshold — defer.
			deferred = append(deferred, r)
		}
	}

	return append(promoted, deferred...)
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
		// Each consecutive ordered pair earns scorePhrasePartial points;
		// a full exact phrase earns scorePhraseExact (always > any partial sum).
		if isExactPhrase(nameTokens, queryTokens) {
			score += scorePhraseExact
		} else if pairs := countConsecutivePairs(nameTokens, queryTokens); pairs > 0 {
			score += pairs * scorePhrasePartial
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

	// Factor 5 — Distance tier (no score contribution; shapes the sort key only)
	var distKm float64
	tier := noGeoTier
	if lat != nil && long != nil && (p.Lat != 0 || p.Long != 0) {
		distKm = haversineKm(*lat, *long, p.Lat, p.Long)
		tier = distanceTierOf(distKm)
	}

	return rankedProduct{
		product:    p,
		score:      score,
		coverage:   coverage,
		distanceKm: distKm,
		tier:       tier,
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
