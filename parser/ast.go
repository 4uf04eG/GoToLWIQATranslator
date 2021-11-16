package parser

import (
	"fmt"
)
type Ast interface {
	String() string
}

type Package struct {
	Name string
}

type Declarations []Declaration

type File struct {
	Package      Package
	Declarations Declarations
	Errors       Errors
}

func (i Package) String() string {
	return fmt.Sprintf("  Name: '%s'", i.Name)
}

func (i Declarations) String() string {
	var str string

	for _, statement := range i {
		stmt := statement.String()

		if stmt != "" {
			str += "  " + statement.String() + "\n"
		}
	}

	return str
}

func (i File) String() string {
	return fmt.Sprintf("Package:\n%s\nDeclarations:\n%s\nErrors:\n%s\n",
		i.Package.String(), i.Declarations.String(), i.Errors.String())
}

//-----------------------------------------------------------------------------
// Declarations
type Declaration interface {
	Ast
}

type Identifier struct {
	Name string
}

type Literal struct {
	Type  TokenType
	Value interface{}
}

type FuncDeclaration struct {
	Name Identifier
	Body BlockStatement
}

func (i Identifier) String() string {
	return fmt.Sprintf("\nIdentifier\n  Name: '%s'\n",
		i.Name)
}

func (i Literal) String() string {
	var litType string

	switch i.Type {
	case 0:
		litType = "Integer"
	case 1:
		litType = "Float"
	case 2:
		litType = "String"
	default:
		litType = "Boolean"
	}

	return fmt.Sprintf("\nLiteral\n  Type:%s\n  Value: '%s'\n",
		litType, i.Value)
}

func (i FuncDeclaration) String() string {
	return fmt.Sprintf("\nFunction declaration\n  Type: %s\n  Value: %s",
		i.Name.String(), i.Body.String())
}

//------------------------------------------------------------------------------
// Expressions
type Expression interface {
	Ast
}

type UnaryExpression struct {
	Operator string
	Operand  Expression
}

type BinaryExpression struct {
	LeftOperand  Expression
	Operator     string
	RightOperand Expression
}

func (i UnaryExpression) String() string {
	if i.Operand == nil {
		return ""
	}

	return fmt.Sprintf("\nUnary expression\n  Operator:%s\n  Operand:\n%s",
		i.Operator, i.Operand.String())
}

func (i BinaryExpression) String() string {
	if i.LeftOperand == nil || i.RightOperand == nil {
		return ""
	}

	return fmt.Sprintf("\nBinary expression:\n  Left operand:\n%s\n  Operator: %s\n  Right operand:\n%s",
		i.LeftOperand.String(), i.Operator, i.RightOperand.String())
}

//------------------------------------------------------------------------------
// Statements
type Statement interface {
	Ast
}

type AssignStatement struct {
	Identifier Identifier
	Operator   string
	Expression Expression
}

type Statements []Statement

// Statements grouped by braces
type BlockStatement struct {
	Statements Statements
}

type BranchStatement struct {
	Keyword string
}

type CaseStatement struct {
	Expression Expression
	Body       BlockStatement
}

type CaseStatements []CaseStatement

type SwitchStatement struct {
	Expression Expression
	Body       CaseStatements
}

type IfStatement struct {
	Condition Expression
	IfBody    BlockStatement
	ElseBody  BlockStatement
}

func (i AssignStatement) String() string {
	return fmt.Sprintf("\nAssign statement\n  Identifier:\n%s\n  Operator: %s\n  Expression:\n%s",
		i.Identifier.String(), i.Operator, i.Expression.String())
}

func (i Statements) String() string {
	var str string

	for _, statement := range i {
		stmt := statement.String()

		if stmt != "" {
			str += "\n" + stmt + "\n"
		}
	}

	return str
}

func (i BlockStatement) String() string {
	return fmt.Sprintf("\nBlock statement:\n  Statements:%s", i.Statements.String())
}

func (i BranchStatement) String() string {
	return fmt.Sprintf("\nBranch statement:\n  Keyword: '%s'", i.Keyword)
}

func (i CaseStatements) String() string {
	var str string

	for _, statement := range i {
		stmt := statement.String()

		if stmt != "" {
			str += "\n" + statement.String()
		}
	}

	return str
}

func (i CaseStatement) String() string {
	return fmt.Sprintf("\nCase statement:\n  Case:\n%s  Body:%s",
		i.Expression.String(), i.Body.String())
}

func (i SwitchStatement) String() string {
	return fmt.Sprintf("\nSwitch statement:\n  Expression:%s  Body:%s",
		i.Expression.String(), i.Body.String())
}

func (i IfStatement) String() string {
	return fmt.Sprintf("\nIf statement:\n Condition:%s  If body:%s  Else body:%s",
		i.Condition.String(), i.IfBody.String(), i.ElseBody.String())
}
