package services

import (
	"math"
	"sort"

	"github.com/atharvyadav96k/SPOTNEARR_API/factories/token"
	"github.com/atharvyadav96k/SPOTNEARR_API/models"
)

const (
	bonusDistanceVeryClose = 6
	bonusDistanceClose     = 5
	bonusDistanceNear      = 4
	bonusDistanceModerate  = 3
	bonusDistanceFar       = 2
	bonusDistanceVeryFar   = 1
)

const (
	scoreNameToken   = 100
	scoreUnitToken   = 80
	scoreFilterToken = 60
	scoreDescToken   = 10

	bonusCoverage100  = 300
	bonusCoverage90   = 250
	bonusCoverage80   = 200
	bonusCoverage70   = 150
	bonusCoverage60   = 100
	bonusCoverage50   = 50
	bonusExactPhrase  = 500
	bonusConsecutive  = 50
	bonusPosZero      = 100
	bonusPosOne       = 80
	bonusPosTwo       = 60
	bonusPosThree     = 40
	bonusPosOther     = 20
	bonusCompleteName = 400
	bonusExactUnit    = 150
	bonusExactFilter  = 100
	bonusAvailability = 1000
)

type rankedProduct struct {
	product        models.ProductResult
	score          int
	coverage       float64
	nameMatchCount int
	exactPhrase    bool
	filterMatches  int
	distanceKm     float64
}

func rankProducts(products []models.ProductResult, parsed token.ParseResult, lat, long *float64) []models.ProductResult {
	ranked := make([]rankedProduct, len(products))
	for i, p := range products {
		ranked[i] = scoreProduct(p, parsed, lat, long)
	}

	sort.SliceStable(ranked, func(i, j int) bool {
		a, b := ranked[i], ranked[j]
		if a.score != b.score {
			return a.score > b.score
		}
		if a.coverage != b.coverage {
			return a.coverage > b.coverage
		}
		if a.nameMatchCount != b.nameMatchCount {
			return a.nameMatchCount > b.nameMatchCount
		}
		if a.exactPhrase != b.exactPhrase {
			return a.exactPhrase
		}
		if a.filterMatches != b.filterMatches {
			return a.filterMatches > b.filterMatches
		}
		// Closer store wins as final tie-breaker (only meaningful when location provided)
		if a.distanceKm != b.distanceKm {
			return a.distanceKm < b.distanceKm
		}
		return false
	})

	result := make([]models.ProductResult, len(ranked))
	for i, r := range ranked {
		r.product.DistanceKm = r.distanceKm
		result[i] = r.product
	}
	return result
}

func scoreProduct(p models.ProductResult, parsed token.ParseResult, lat, long *float64) rankedProduct {
	score := bonusAvailability

	nameTokens := token.TokenParser(p.Name)
	descTokens := token.TokenParser(p.Desc)
	nameSet := toSet(nameTokens)
	descSet := toSet(descTokens)
	queryTokens := parsed.Tokens

	totalSignals := len(queryTokens) + len(parsed.Units)
	if parsed.Price.Found {
		totalSignals++
	}

	matched := 0
	nameMatchCount := 0

	// Rule 1 — field weights
	for _, qt := range queryTokens {
		if _, ok := nameSet[qt]; ok {
			score += scoreNameToken
			matched++
			nameMatchCount++
		} else if _, ok := descSet[qt]; ok {
			score += scoreDescToken
			matched++
		}
	}

	// Rule 1 (unit weight) + Rule 7 (exact unit bonus)
	for _, u := range parsed.Units {
		if p.QuantityUnit != nil && token.NormalizeUnit(*p.QuantityUnit) == u.RawUnit {
			score += scoreUnitToken + bonusExactUnit
			matched++
		}
	}

	// Rule 1 (filter weight) + Rule 8 (exact filter bonus)
	filterMatches := 0
	if parsed.Price.Found && priceMatchesFilter(p.Price, parsed.Price) {
		score += scoreFilterToken + bonusExactFilter
		matched++
		filterMatches++
	}

	// Rule 2 — coverage bonus
	var coverage float64
	if totalSignals > 0 {
		coverage = float64(matched) / float64(totalSignals)
		score += coverageBonus(coverage)
	}

	// Rule 3 — exact phrase bonus
	exactPhrase := isExactPhrase(nameTokens, queryTokens)
	if exactPhrase {
		score += bonusExactPhrase
	}

	// Rule 4 — consecutive token bonus
	score += countConsecutivePairs(nameTokens, queryTokens) * bonusConsecutive

	// Rule 5 — position bonus (first matched query token in product name)
	score += firstMatchPositionBonus(nameTokens, toSet(queryTokens))

	// Rule 6 — complete name match bonus
	if len(queryTokens) > 0 && allInSet(queryTokens, nameSet) {
		score += bonusCompleteName
	}

	// Distance bonus — only applied when user location is provided
	var distKm float64
	if lat != nil && long != nil && (p.Lat != 0 || p.Long != 0) {
		distKm = haversineKm(*lat, *long, p.Lat, p.Long)
		score += distanceBonus(distKm)
	}

	return rankedProduct{
		product:        p,
		score:          score,
		coverage:       coverage,
		nameMatchCount: nameMatchCount,
		exactPhrase:    exactPhrase,
		filterMatches:  filterMatches,
		distanceKm:     distKm,
	}
}

func coverageBonus(c float64) int {
	switch {
	case c >= 1.0:
		return bonusCoverage100
	case c >= 0.9:
		return bonusCoverage90
	case c >= 0.8:
		return bonusCoverage80
	case c >= 0.7:
		return bonusCoverage70
	case c >= 0.6:
		return bonusCoverage60
	case c >= 0.5:
		return bonusCoverage50
	default:
		return 0
	}
}

// isExactPhrase returns true if all queryTokens appear consecutively and in
// order anywhere within nameTokens.
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

// countConsecutivePairs returns the number of adjacent query token pairs that
// also appear adjacent and in the same order in nameTokens.
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

// firstMatchPositionBonus returns the position bonus based on where the first
// query token match appears in the product name token list.
func firstMatchPositionBonus(nameTokens []string, querySet map[string]struct{}) int {
	for i, t := range nameTokens {
		if _, ok := querySet[t]; ok {
			switch i {
			case 0:
				return bonusPosZero
			case 1:
				return bonusPosOne
			case 2:
				return bonusPosTwo
			case 3:
				return bonusPosThree
			default:
				return bonusPosOther
			}
		}
	}
	return 0
}

func allInSet(tokens []string, set map[string]struct{}) bool {
	for _, t := range tokens {
		if _, ok := set[t]; !ok {
			return false
		}
	}
	return true
}

func priceMatchesFilter(price float64, pf token.PriceFilter) bool {
	if pf.Approximate > 0 {
		return math.Abs(price-float64(pf.Approximate))/float64(pf.Approximate) <= 0.1
	}
	if pf.Min > 0 && pf.Max > 0 {
		return price >= float64(pf.Min) && price <= float64(pf.Max)
	}
	if pf.Max > 0 {
		return price <= float64(pf.Max)
	}
	if pf.Min > 0 {
		return price >= float64(pf.Min)
	}
	return false
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

func distanceBonus(km float64) int {
	switch {
	case km < 0.5:
		return bonusDistanceVeryClose
	case km < 1:
		return bonusDistanceClose
	case km < 2:
		return bonusDistanceNear
	case km < 5:
		return bonusDistanceModerate
	case km < 10:
		return bonusDistanceFar
	case km < 20:
		return bonusDistanceVeryFar
	default:
		return 0
	}
}

func toSet(tokens []string) map[string]struct{} {
	s := make(map[string]struct{}, len(tokens))
	for _, t := range tokens {
		s[t] = struct{}{}
	}
	return s
}
