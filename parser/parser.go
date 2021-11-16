package parser

import (
	"../lexer"
	"strconv"
)

type Parser struct {
	tokens       []lexer.Token
	currentIndex int
	currentToken lexer.Token
	errors       []*Error
}

func NewParser(tokens []lexer.Token) *Parser {
	return &Parser{tokens, 0, tokens[0], []*Error{}}
}

func (parser *Parser) Parse() File {
	ast := File{}
	pkg := parser.parsePackage()

	if pkg.Name != "" {
		ast.Package = pkg
	} else { // Keyword 'package' not found
		ast.Errors = parser.errors
		return ast
	}

	for !parser.foundEndOfFile() {
		parser.skipLineEndings()
		function := parser.parseFunctionDeclaration()

		if function.Name.Name != "" {
			ast.Declarations = append(ast.Declarations, function)
		}
	}

	ast.Errors = parser.errors
	return ast
}

func (parser *Parser) isErrorFound(error *Error) bool {
	if error != nil {
		parser.errors = append(parser.errors, error)
		return true
	}

	return false
}

func (parser *Parser) foundEndOfFile() bool {
	return parser.currentIndex >= len(parser.tokens)
}

func (parser *Parser) skipLineEndings() {
	for parser.currentToken.TokenType == lexer.EndOfLine && !parser.foundEndOfFile() {
		parser.nextToken()
	}
}

//------------------------------------------------------------------------------------------
// Expect functions (Checking if current token has right type)
func (parser *Parser) expect(tokenType TokenType) *Error {
	parser.skipLineEndings()

	if !parser.isTokenOfType(tokenType) {
		token := parser.currentToken.Text
		parser.nextToken()
		return NewExpectError(tokenTypes[tokenType], token)
	}

	parser.nextToken()

	return nil
}

func (parser *Parser) expectType(tokenType lexer.TokenType) *Error {
	if parser.currentToken.TokenType != tokenType {
		return NewExpectError(string(tokenType), parser.currentToken.Text)
	}

	// Used when calculating binary statements.
	// Switching to next token is performed there.

	return nil
}

// Used to either check if statements on different lines
// or on the same, but delimited with semicolon.
func (parser *Parser) expectSemicolon() *Error {
	tok := parser.currentToken

	if tok.Text != GetType(Semicolon) && tok.TokenType != lexer.EndOfLine {
		return NewExpectError("';' or a new line", tok.Text)
	}

	parser.nextToken()

	return nil
}

func (parser *Parser) isTokenOfType(tokenType TokenType) bool {
	parser.skipLineEndings()
	return parser.currentToken.Text == GetType(tokenType)
}

func (parser *Parser) nextToken() {
	parser.currentIndex++

	if !parser.foundEndOfFile() {
		parser.currentToken = parser.tokens[parser.currentIndex]
	}
}

func (parser *Parser) parsePackage() Package {
	errKey := parser.expect(Pkg)
	ident, errIdent := parser.parseIdentifier()
	errSemi := parser.expectSemicolon()

	if parser.isErrorFound(errKey) { // Keyword 'package' not found
		return Package{}
	}

	parser.isErrorFound(errIdent)
	parser.isErrorFound(errSemi)

	return Package{ident.Name}
}

//------------------------------------------------------------------------------------------
// Parsing declarations
func (parser *Parser) parseIdentifier() (Identifier, *Error) {
	token := parser.currentToken

	if token.TokenType != lexer.Identifier {
		return Identifier{}, NewTypeError(lexer.Identifier, token)
	}

	parser.nextToken()

	return Identifier{token.Text}, nil
}

func (parser *Parser) parseLiteral(token lexer.Token) (Expression, *Error) {
	if token.TokenType != lexer.Literal {
		return Literal{}, NewTypeError(lexer.Literal, token)
	}

	if res, err := strconv.ParseInt(token.Text, 10, 32); err == nil {
		return Literal{IntegerLiteral, res}, nil
	} else if res, err := strconv.ParseBool(token.Text); err == nil {
		return Literal{BooleanLiteral, res}, nil
	} else if res, err := strconv.ParseFloat(token.Text, 64); err == nil {
		return Literal{FloatLiteral, res}, nil
	} else if len(token.Text) > 1 { // String literal with quotes removed
		return Literal{StringLiteral, token.Text[1 : len(token.Text)-1]}, nil
	}

	// For some reason switching to next token here
	// Causes index exception when dealing with block statements.
	// So switching occurs only when parsing unary expressions.

	return Literal{}, &Error{WrongLiteralError, "Unexpected literal '" + token.Text + "'"}
}

func (parser *Parser) parseFunctionDeclaration() FuncDeclaration {
	if parser.currentToken.TokenType == lexer.EndOfLine {
		return FuncDeclaration{}
	}

	err := parser.expect(Func)

	if parser.isErrorFound(err) {
		return FuncDeclaration{}
	}

	name, err := parser.parseIdentifier()
	parser.isErrorFound(err)

	parser.isErrorFound(parser.expect(LeftParen))
	parser.isErrorFound(parser.expect(RightParen))

	parser.isErrorFound(parser.expect(LeftBrace))
	body := parser.parseBlockStatement()
	parser.isErrorFound(parser.expect(RightBrace))

	return FuncDeclaration{name, body}
}

//------------------------------------------------------------------------------------------
// Parsing expressions
func (parser *Parser) parseExpression() (Expression, *Error) {
	return parser.parseBinaryExpression(LowestPrecedence)
}

// Recursively parsing expressions until non-operator found
func (parser *Parser) parseBinaryExpression(prec int) (Expression, *Error) {
	left, _ := parser.parseUnaryExpression()

	for {
		otherPrec := precedence(parser.currentToken)
		operator := parser.currentToken.Text

		err := parser.expectType(lexer.Operator)

		if prec > otherPrec || err != nil {
			return left, nil
		}

		parser.nextToken()
		right, err := parser.parseBinaryExpression(otherPrec)

		if parser.isErrorFound(err) {
			return left, err
		}

		left = BinaryExpression{left, operator, right}
	}
}

func (parser *Parser) parseUnaryExpression() (Expression, *Error) {
	if parser.currentToken.TokenType == lexer.Operator {
		operator := parser.currentToken.Text
		parser.nextToken()
		x, err := parser.parseUnaryExpression()

		if err != nil {
			return UnaryExpression{}, err
		}

		return UnaryExpression{operator, x}, nil
	}

	return parser.parseSimpleExpression()
}

func (parser *Parser) parseSimpleExpression() (Expression, *Error) {
	switch parser.currentToken.TokenType {
	case lexer.Identifier:
		ident, err := parser.parseIdentifier()
		parser.isErrorFound(err)
		return ident, err
	case lexer.Literal:
		lit, err := parser.parseLiteral(parser.currentToken)
		parser.isErrorFound(err)
		parser.nextToken()
		return lit, err
	}

	return UnaryExpression{}, NewExpectError("Expression", parser.currentToken.Text)
}

//------------------------------------------------------------------------------------------
// Parsing statements
func (parser *Parser) parseStatement() (Statement, *Error) {
	switch parser.currentToken.TokenType {
	case lexer.Identifier:
		return parser.parseAssignStatement(), nil
	case lexer.EndOfLine:
		parser.nextToken()
		return parser.parseStatement()
	}

	switch parser.currentToken.Text {
	case GetType(Var):
		return parser.parseAssignStatement(), nil
	case GetType(Switch):
		return parser.parseSwitchStatement(), nil
	case GetType(If):
		return parser.parseIfStatement(), nil
	case GetType(Break), GetType(Continue), GetType(Return):
		keyword := parser.currentToken.Text
		parser.nextToken()
		return BranchStatement{keyword}, nil
	}

	return nil, NewExpectError("Statement", parser.currentToken.Text)
}

// Assign statement can look like:
//   var a = 2 (Initializing)
//   a := 2    (Initializing)
//   a = 2
func (parser *Parser) parseAssignStatement() AssignStatement {
	if parser.isTokenOfType(Var) {
		parser.nextToken()
	}

	ident, err := parser.parseIdentifier()
	parser.isErrorFound(err)

	if !parser.isTokenOfType(Assign) && !parser.isTokenOfType(Define) {
		err := NewExpectError("':=' or '='", parser.currentToken.Text)
		parser.isErrorFound(err)
		return AssignStatement{Expression: UnaryExpression{}}
	}

	operator := parser.currentToken.Text
	parser.nextToken()
	expr, _ := parser.parseExpression()

	err = parser.expectSemicolon()
	parser.isErrorFound(err)

	return AssignStatement{ident, operator, expr}
}

func (parser *Parser) parseBlockStatement() BlockStatement {
	var statements []Statement

	for !parser.isTokenOfType(Default) && !parser.isTokenOfType(Case) &&
		!parser.isTokenOfType(RightBrace) && !parser.foundEndOfFile() {
		stmt, err := parser.parseStatement()

		if parser.isErrorFound(err) {
			parser.nextToken()
			continue
		}

		statements = append(statements, stmt)
	}

	return BlockStatement{statements}
}

func (parser *Parser) parseIfStatement() IfStatement {
	parser.nextToken()

	cond, _ := parser.parseExpression()
	parser.isErrorFound(parser.expect(LeftBrace))
	ifBody := parser.parseBlockStatement()
	parser.isErrorFound(parser.expect(RightBrace))
	elseBody := BlockStatement{}

	if parser.isTokenOfType(Else) {
		parser.nextToken()
		parser.isErrorFound(parser.expect(LeftBrace))
		elseBody = parser.parseBlockStatement()
		parser.isErrorFound(parser.expect(RightBrace))
	}

	return IfStatement{cond, ifBody, elseBody}
}

func (parser *Parser) parseCaseStatement() CaseStatement {
	expr := Expression(UnaryExpression{})

	if parser.isTokenOfType(Case) {
		parser.nextToken()
		expr, _ = parser.parseExpression()

		if expr.String() == "" {
			err := NewExpectError("Expression", parser.currentToken.Text)
			parser.isErrorFound(err)
		}
	} else {
		parser.nextToken()
	}

	parser.isErrorFound(parser.expect(Colon))
	body := parser.parseBlockStatement()

	return CaseStatement{expr, body}
}

func (parser *Parser) parseSwitchStatement() SwitchStatement {
	parser.nextToken()
	expression := Expression(UnaryExpression{})

	if !parser.isTokenOfType(LeftBrace) { // If switch has parameters
		expression, _ = parser.parseExpression()
	}

	parser.isErrorFound(parser.expect(LeftBrace))
	var cases []CaseStatement

	for parser.isTokenOfType(Case) || parser.isTokenOfType(Default) {
		cases = append(cases, parser.parseCaseStatement())
	}

	parser.isErrorFound(parser.expect(RightBrace))
	parser.expectSemicolon()

	return SwitchStatement{expression, cases}
}
