package parser

import (
	"../lexer"
	"strings"
)

const (
	TypeError = iota
	WrongLiteralError
	ExpectError
	AssignError
	MismatchedTypesError
)

type Error struct {
	Type    int
	Message string
}

func NewTypeError(expectedType string, realType lexer.Token) *Error {
	message := expectedType + " expected, " +
		"got '" + realType.Text +
		"' of type '" + string(realType.TokenType) + "'"
	return &Error{TypeError, message}
}

func NewExpectError(expectedToken string, realToken string) *Error {
	realToken = strings.Replace(realToken, "\n", "end of file", -1)
	message := expectedToken + " expected, " + "got '" + realToken + "'"
	return &Error{ExpectError, message}
}

func (err *Error) String() string {
	return err.Message
}

type Errors []*Error

func (slice Errors) String() string {
	var str string

	for _, item := range slice {
		str += "\t" + item.String() + "\n"
	}

	return str
}
