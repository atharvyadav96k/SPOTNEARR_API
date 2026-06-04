package token

import "strings"

// ParseResult is returned by TokenParse.
type ParseResult struct {
	Tokens []string
	Price  PriceFilter
	Units  []UnitMatch
}

func QueryTokenParse(val string) ParseResult {
	clean := strings.ToLower(val)
	tokes := SplitAlphaNumeric(clean) // step 2

	units, tokes := ExtractUnitFilter(tokes) // step 3 — before anything mangles units

	tokes = StopWordFilter(tokes)     // step 4
	tokes = MakeTokensSingular(tokes) // step 5
	tokes = CleanPunctuation(tokes)   // step 6

	price, tokes := ExtractPriceFilter(tokes) // step 7

	return ParseResult{
		Tokens: tokes,
		Price:  price,
		Units:  units,
	}
}

func TokenParser(val string) []string {
	clean := strings.ToLower(val)
	tokes := SplitAlphaNumeric(clean)
	tokes = StopWordFilter(tokes)
	tokes = MakeTokensSingular(tokes)
	tokes = CleanPunctuation(tokes)
	return tokes
}

// MergeTokens tokenizes each field and returns a deduplicated union.
// Name tokens appear first so they stay higher in the slice for readability.
func MergeTokens(fields ...string) []string {
	seen := make(map[string]struct{})
	var result []string
	for _, f := range fields {
		for _, t := range TokenParser(f) {
			if _, exists := seen[t]; !exists {
				seen[t] = struct{}{}
				result = append(result, t)
			}
		}
	}
	return result
}
