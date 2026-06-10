package tokenizer

import "strings"

var singularGuard = map[string]struct{}{
	"adidas": {}, "levis": {}, "crocs": {}, "boss": {},
	"xps": {}, "asus": {}, "bose": {},
	"plus":    {}, "class": {}, "glass": {}, "bass": {}, "press": {},
	"access":  {}, "unless": {}, "across": {}, "process": {}, "canvas": {},
	"bonus":   {}, "genius": {}, "status": {}, "versus": {}, "nexus": {},
	"campus":  {}, "chorus": {}, "corpus": {}, "focus": {}, "radius": {},
	"series":  {}, "species": {},
	"mbps": {}, "gbps": {}, "kbps": {},
	"watts": {},
}

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
	if _, ok := singularGuard[word]; ok {
		return word
	}
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
	if strings.HasSuffix(word, "ves") {
		base := word[:len(word)-3]
		if base == "kni" || base == "li" || base == "wi" {
			return base + "fe"
		}
		return base + "f"
	}
	if strings.HasSuffix(word, "ies") {
		return word[:len(word)-3] + "y"
	}
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
	if strings.HasSuffix(word, "s") {
		base := word[:len(word)-1]
		if len(base) <= 1 {
			return word
		}
		return base
	}
	return word
}
