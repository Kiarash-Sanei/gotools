package common

import (
	"fmt"
	"strings"
	"unicode"
)

func ParseTokens(input string) ([]string, error) {
	var tokens []string
	var token strings.Builder
	var quote rune
	escaped := false
	tokenStarted := false

	for _, character := range input {
		if escaped {
			token.WriteRune(character)
			escaped = false
			tokenStarted = true
			continue
		}

		if character == '\\' && quote != '\'' {
			escaped = true
			tokenStarted = true
			continue
		}

		if quote != 0 {
			if character == quote {
				quote = 0
			} else {
				token.WriteRune(character)
			}
			tokenStarted = true
			continue
		}

		switch character {
		case '\'', '"':
			quote = character
			tokenStarted = true
		case '\n', '\r':
			continue
		default:
			if unicode.IsSpace(character) {
				if tokenStarted {
					tokens = append(tokens, token.String())
					token.Reset()
					tokenStarted = false
				}
			} else {
				token.WriteRune(character)
				tokenStarted = true
			}
		}
	}

	if escaped {
		return nil, fmt.Errorf("input ends with an escape character")
	}
	if quote != 0 {
		return nil, fmt.Errorf("unclosed quote")
	}
	if tokenStarted {
		tokens = append(tokens, token.String())
	}
	return tokens, nil
}