package token

import (
	"strconv"
)

type PriceFilter struct {
	Min         int
	Max         int
	Approximate int
	Found       bool
}

var maxKeywords = map[string]struct{}{
	"under": {}, "below": {}, "upto": {}, "within": {},
}

var minKeywords = map[string]struct{}{
	"above": {}, "over": {}, "atleast": {},
}

var lessThanKeywords = map[string]struct{}{
	"less": {}, "cheaper": {},
}

var moreThanKeywords = map[string]struct{}{
	"more": {}, "expensive": {},
}

var approxKeywords = map[string]struct{}{
	"around": {}, "near": {}, "about": {}, "approximately": {},
}

var rangeKeywords = map[string]struct{}{
	"between": {},
}

var ignoredPriceWords = map[string]struct{}{
	"than": {}, "rs": {}, "inr": {}, "rupee": {}, "rupees": {},
	"dollar": {}, "dollars": {}, "usd": {}, "to": {},
}

func ExtractPriceFilter(tokens []string) (PriceFilter, []string) {
	pf := PriceFilter{}
	kept := make([]string, 0, len(tokens))

	skip := make([]bool, len(tokens))

	for i, t := range tokens {
		if _, ok := rangeKeywords[t]; ok {
			n1, ok1 := lookAheadNumber(tokens, i+1)
			n2, ok2 := lookAheadNumber(tokens, i+3)
			if !ok2 {
				n2, ok2 = lookAheadNumber(tokens, i+2)
			}
			if ok1 && ok2 {
				pf.Min, pf.Max, pf.Found = n1, n2, true
				skip[i] = true
				markNumberSkip(tokens, i+1, skip)
				markNumberSkip(tokens, i+2, skip)
				markNumberSkip(tokens, i+3, skip)
				continue
			}
		}

		if isNumericToken(t) {
			if i+2 < len(tokens) && tokens[i+1] == "to" && isNumericToken(tokens[i+2]) {
				n1, _ := strconv.Atoi(t)
				n2, _ := strconv.Atoi(tokens[i+2])
				pf.Min, pf.Max, pf.Found = n1, n2, true
				skip[i], skip[i+1], skip[i+2] = true, true, true
				continue
			}
		}

		if _, ok := maxKeywords[t]; ok {
			if n, ok2 := lookAheadNumber(tokens, i+1); ok2 {
				pf.Max, pf.Found = n, true
				skip[i] = true
				markNumberSkip(tokens, i+1, skip)
				continue
			}
		}

		if _, ok := minKeywords[t]; ok {
			if n, ok2 := lookAheadNumber(tokens, i+1); ok2 {
				pf.Min, pf.Found = n, true
				skip[i] = true
				markNumberSkip(tokens, i+1, skip)
				continue
			}
		}

		if _, ok := lessThanKeywords[t]; ok {
			offset := 1
			if i+1 < len(tokens) && tokens[i+1] == "than" {
				offset = 2
			}
			if n, ok2 := lookAheadNumber(tokens, i+offset); ok2 {
				pf.Max, pf.Found = n, true
				skip[i] = true
				if offset == 2 {
					skip[i+1] = true
				}
				markNumberSkip(tokens, i+offset, skip)
				continue
			}
		}

		if _, ok := moreThanKeywords[t]; ok {
			offset := 1
			if i+1 < len(tokens) && tokens[i+1] == "than" {
				offset = 2
			}
			if n, ok2 := lookAheadNumber(tokens, i+offset); ok2 {
				pf.Min, pf.Found = n, true
				skip[i] = true
				if offset == 2 {
					skip[i+1] = true
				}
				markNumberSkip(tokens, i+offset, skip)
				continue
			}
		}

		if _, ok := approxKeywords[t]; ok {
			if n, ok2 := lookAheadNumber(tokens, i+1); ok2 {
				pf.Approximate, pf.Found = n, true
				skip[i] = true
				markNumberSkip(tokens, i+1, skip)
				continue
			}
		}

		if _, ok := ignoredPriceWords[t]; ok {
			skip[i] = true
			continue
		}
	}

	for i, t := range tokens {
		if !skip[i] {
			kept = append(kept, t)
		}
	}

	return pf, kept
}

func lookAheadNumber(tokens []string, idx int) (int, bool) {
	for idx < len(tokens) {
		t := tokens[idx]
		if _, noise := ignoredPriceWords[t]; noise {
			idx++
			continue
		}
		if n, err := strconv.Atoi(t); err == nil {
			return n, true
		}
		return 0, false
	}
	return 0, false
}

func markNumberSkip(tokens []string, idx int, skip []bool) {
	for idx < len(tokens) {
		t := tokens[idx]
		if _, noise := ignoredPriceWords[t]; noise {
			skip[idx] = true
			idx++
			continue
		}
		if isNumericToken(t) {
			skip[idx] = true
		}
		return
	}
}

func isNumericToken(s string) bool {
	if len(s) == 0 {
		return false
	}
	_, err := strconv.Atoi(s)
	return err == nil
}
