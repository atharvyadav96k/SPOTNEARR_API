package tokenizer

import "strings"

type ParseResult struct {
	Tokens []string
	Price  PriceFilter
	Units  []UnitMatch
}

func QueryTokenParse(val string) ParseResult {
	clean := strings.ToLower(val)
	tokes := SplitAlphaNumeric(clean)
	units, tokes := ExtractUnitFilter(tokes)
	tokes = StopWordFilter(tokes)
	tokes = MakeTokensSingular(tokes)
	tokes = CleanPunctuation(tokes)
	price, tokes := ExtractPriceFilter(tokes)
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
