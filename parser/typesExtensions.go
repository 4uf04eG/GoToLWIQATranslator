package parser

import "../lexer"

// Due to bad program design there's two 'TokenType's:
// one in lexer and one's there with different type.
// So there's might be some misunderstandings.
type TokenType int

const (
	IntegerLiteral = iota
	FloatLiteral
	StringLiteral
	BooleanLiteral

	Minus
	Plus
	Mul
	Div
	Mod

	And
	Or
	Eq
	Less
	Greater

	Not
	Neq
	Leq
	Geq

	Define
	Assign

	Comma
	Colon
	Semicolon

	LeftBrace
	RightBrace
	LeftParen
	RightParen

	Pkg
	Func
	Var
	Case
	Default
	Switch
	If
	Else
	For
	Break
	Continue
	Return
)

var tokenTypes = map[TokenType]string{
	Minus: "-",
	Plus:  "+",
	Mul:   "*",
	Div:   "/",
	Mod:   "%",

	And:     "&&",
	Or:      "||",
	Eq:      "==",
	Less:    "<",
	Greater: ">",

	Not: "!",
	Neq: "!=",
	Leq: "<=",
	Geq: ">=",

	Define: ":=",
	Assign: "=",

	Comma:     ",",
	Colon:     ":",
	Semicolon: ";",

	LeftBrace:  "{",
	RightBrace: "}",
	LeftParen:  "(",
	RightParen: ")",

	Pkg:      "package",
	Func:     "func",
	Var:      "var",
	Case:     "case",
	Default:  "default",
	Switch:   "switch",
	If:       "if",
	Else:     "else",
	For:      "for",
	Break:    "break",
	Continue: "continue",
	Return:   "return",
}

const (
	LowestPrecedence = 0
)

func GetType(tokenType TokenType) string {
	return tokenTypes[tokenType]
}

func precedence(token lexer.Token) int {
	switch token.Text {
	case GetType(Or):
		return 1
	case GetType(And):
		return 2
	case GetType(Eq), GetType(Neq), GetType(Less), GetType(Leq), GetType(Greater), GetType(Geq):
		return 3
	case GetType(Plus), GetType(Minus):
		return 4
	case GetType(Mul), GetType(Div), GetType(Mod):
		return 5
	}

	return LowestPrecedence
}


