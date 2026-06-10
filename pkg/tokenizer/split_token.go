package tokenizer

import "unicode"

func SplitAlphaNumeric(cleanQuery string) []string {
	var tokens []string
	var currentToken []rune

	const (
		typeSpace = iota
		typeLetter
		typeNumber
	)
	lastType := typeSpace

	for _, r := range cleanQuery {
		var currentType int
		if unicode.IsSpace(r) || (!unicode.IsLetter(r) && !unicode.IsDigit(r)) {
			currentType = typeSpace
		} else if unicode.IsDigit(r) {
			currentType = typeNumber
		} else {
			currentType = typeLetter
		}

		if currentType == typeSpace {
			if len(currentToken) > 0 {
				tokens = append(tokens, string(currentToken))
				currentToken = []rune{}
			}
			lastType = typeSpace
			continue
		}

		if lastType != typeSpace && currentType != lastType {
			if len(currentToken) > 0 {
				tokens = append(tokens, string(currentToken))
				currentToken = []rune{}
			}
		}

		currentToken = append(currentToken, r)
		lastType = currentType
	}

	if len(currentToken) > 0 {
		tokens = append(tokens, string(currentToken))
	}

	return tokens
}

func CleanPunctuation(tokens []string) []string {
	cleanedTokens := make([]string, 0, len(tokens))
	for _, tok := range tokens {
		filteredRunes := make([]rune, 0, len(tok))
		for _, r := range tok {
			if r == '-' || r == '_' {
				continue
			}
			filteredRunes = append(filteredRunes, r)
		}
		if cleaned := string(filteredRunes); cleaned != "" {
			cleanedTokens = append(cleanedTokens, cleaned)
		}
	}
	return cleanedTokens
}
