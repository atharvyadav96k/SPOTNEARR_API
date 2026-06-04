package token

import "strings"

// singularGuard: tokens ending in "s" that must never be singularised.
// Covers brand model slugs, units, and words where "s" is part of the root.
var singularGuard = map[string]struct{}{
	// brands / models
	"adidas": {}, "levis": {}, "crocs": {}, "boss": {},
	"xps": {}, "asus": {}, "bose": {},
	// common product words
	"plus":    {}, // already in noStripWords but double-guard here
	"class":   {},
	"glass":   {},
	"bass":    {},
	"press":   {},
	"access":  {},
	"unless":  {},
	"across":  {},
	"process": {},
	"canvas":  {},
	"bonus":   {},
	"genius":  {},
	"status":  {},
	"versus":  {},
	"nexus":   {},
	"campus":  {},
	"chorus":  {},
	"corpus":  {},
	"focus":   {},
	"radius":  {},
	"series":  {}, // "series" → "serie" is wrong
	"species": {},
	// units that end in s
	"mbps": {}, "gbps": {}, "kbps": {},
	"watts": {}, // handled as unit, but guard anyway
}

// MakeTokensSingular loops through a slice of tokens and converts each to singular.
func MakeTokensSingular(tokens []string) []string {
	out := make([]string, 0, len(tokens))
	for _, t := range tokens {
		out = append(out, toSingular(t))
	}
	return out
}

func toSingular(word string) string {
	if len(word) <= 2 {
		return word
	}

	// Guard: known words that must not be singularised
	if _, ok := singularGuard[word]; ok {
		return word
	}

	// Irregular
	switch word {
	case "geese":
		return "goose"
	case "feet":
		return "foot"
	case "men":
		return "man"
	case "women":
		return "woman"
	case "data":
		return "datum"
	case "sheep", "deer", "fish":
		return word
	}

	// -ves → -f / -fe
	if strings.HasSuffix(word, "ves") {
		base := word[:len(word)-3]
		if base == "kni" || base == "li" || base == "wi" {
			return base + "fe"
		}
		return base + "f"
	}

	// -ies → -y
	if strings.HasSuffix(word, "ies") {
		return word[:len(word)-3] + "y"
	}

	// -es
	if strings.HasSuffix(word, "es") {
		if strings.HasSuffix(word, "oes") {
			return word[:len(word)-2]
		}
		base := word[:len(word)-2]
		if strings.HasSuffix(base, "ch") || strings.HasSuffix(base, "sh") ||
			strings.HasSuffix(base, "x") || strings.HasSuffix(base, "ss") ||
			strings.HasSuffix(base, "z") {
			return base
		}
		return word[:len(word)-1]
	}

	// -s  (plain)
	if strings.HasSuffix(word, "s") {
		// Don't strip if the word is all-consonant slug like "xps", "tvs" etc.
		// Heuristic: if stripping leaves something ≤ 1 char, skip.
		base := word[:len(word)-1]
		if len(base) <= 1 {
			return word
		}
		return base
	}

	return word
}
