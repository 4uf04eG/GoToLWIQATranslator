package lexer

import (
	"regexp"
 	"strings"
)

type Lexer struct {
	text          string
	patterns      []Pair
	remainingText string
}

func NewLexer(text string) *Lexer {
	var compiledPatterns []Pair

	for _, pair := range patterns {
		tokenType := TokenType(pair.Key.(string))
		regex := pair.Value.(string)

		value := Pair{tokenType, *regexp.MustCompile(regex)}
		compiledPatterns = append(compiledPatterns, value)
	}

	return &Lexer{text, compiledPatterns, text}
}

func (lexer *Lexer) Tokenize() []Token {
	var tokens []Token

	for lexer.remainingText != "" {
		token := lexer.nextToken()

		if token.Text == "" {
			break
		}

		if token.TokenType != Comment {
			tokens = append(tokens, token)
		} else if token.Text[0:2] == "/*" {
			lexer.consumeMultilineComment()
		}
	}

	return tokens
}

func (lexer *Lexer) nextToken() Token {
	var token Token
	var nearestIndices []int
	var nearestType TokenType

	for _, pair := range lexer.patterns {
		tokenType := pair.Key.(TokenType)
		regex := pair.Value.(regexp.Regexp)
		indices := regex.FindStringIndex(lexer.remainingText)

		if len(indices) > 0 && (len(nearestIndices) == 0 || nearestIndices[0] > indices[0]) {
			nearestIndices = indices
			nearestType = tokenType
		}
	}

	if len(nearestIndices) > 0 {
		token.TokenType = nearestType
		token.Text = lexer.remainingText[nearestIndices[0]:nearestIndices[1]]
		lexer.remainingText = lexer.remainingText[nearestIndices[1]:]
	}

	return token
}

func (lexer *Lexer) consumeMultilineComment() {
	index := strings.Index(lexer.remainingText, "*/")

	if index != -1 {
		lexer.remainingText = lexer.remainingText[index + 2:]
	} else {
		lexer.remainingText = ""
	}
}
