package tokenizer

import "strings"

var stopWords = map[string]struct{}{
	"is": {}, "the": {}, "a": {}, "an": {}, "and": {}, "or": {}, "for": {},
	"to": {}, "in": {}, "on": {}, "at": {}, "by": {}, "with": {}, "from": {},
	"of": {}, "about": {}, "as": {}, "into": {}, "like": {}, "than": {},
	"this": {}, "that": {}, "these": {}, "those": {}, "it": {}, "its": {},
	"you": {}, "your": {}, "me": {}, "my": {}, "we": {}, "our": {},
	"buy": {}, "get": {}, "where": {}, "can": {}, "do": {}, "does": {},
}

func StopWordFilter(tokens []string) []string {
	cleaned := make([]string, 0, len(tokens))
	for _, t := range tokens {
		if len(t) <= 1 {
			continue
		}
		if _, isStopWord := stopWords[t]; isStopWord {
			continue
		}
		cleaned = append(cleaned, stripGrammarSuffixes(t))
	}
	return cleaned
}

func stripGrammarSuffixes(word string) string {
	if len(word) <= 4 {
		return word
	}
	if strings.HasSuffix(word, "ing") {
		base := word[:len(word)-3]
		if len(base) > 2 && base[len(base)-1] == base[len(base)-2] {
			if !strings.HasSuffix(base, "ll") && !strings.HasSuffix(base, "ss") {
				return base[:len(base)-1]
			}
		}
		return base
	}
	if strings.HasSuffix(word, "ed") {
		return word[:len(word)-2]
	}
	if strings.HasSuffix(word, "ly") {
		return word[:len(word)-2]
	}
	return word
}
