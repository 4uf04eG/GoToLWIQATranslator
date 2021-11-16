package lexer

type TokenType string

type Token struct {
	TokenType TokenType
	Text      string
}
