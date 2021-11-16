package generator

import (
	"../parser"
	"fmt"
	"strconv"
	"strings"
)

type Generator struct {
	syntaxTree parser.File
}

type index []int

func NewGenerator(ast parser.File) *Generator {
	return &Generator{ast}
}

func (generator *Generator) Generate() string {
	var genCode string
	var curIndex = index([]int{1, 1})

	genCode += fmt.Sprintf("Z1 %s\n", generator.syntaxTree.Package.Name)

	for _, stmt := range generator.syntaxTree.Declarations {
		function := stmt.(parser.FuncDeclaration)

		header := fmt.Sprintf("%sQ%s PROCEDURE &%s&\n", curIndex.Indentation(), curIndex.String(), function.Name.Name)
		body := generator.generateStatement(function, append(curIndex, 1))

		curIndex[1]++
		genCode += header + body
	}

	return genCode
}

func (generator *Generator) generateStatement(statement parser.Statement, index index) string {
	if stmt, status := statement.(parser.SwitchStatement); status {                  // Switch
		return generator.generateCaseStatements(stmt.Expression, stmt.Body, index)
	} else if stmt, status := statement.(parser.AssignStatement); status {           // Assign
		return generator.generateAssignStatement(stmt, index)
	} else if stmt, status := statement.(parser.BlockStatement); status {            // Block
		return generator.generateStatement(stmt.Statements, index)
	} else if stmt, status := statement.(parser.BranchStatement); status {
		return fmt.Sprintf("%sQ%s %s\n",
			index.Indentation(), index.String(), stmt.Keyword)
	} else if stmt, status := statement.(parser.IfStatement); status {               // If
		header := fmt.Sprintf("%sQ%s IF %s THEN BEGIN\n",
			index.Indentation(), index.String(), generator.generateExpression(stmt.Condition))
		return header + generator.generateIfStatement(stmt, append(index, 1))
	} else if decl, status := statement.(parser.FuncDeclaration); status {           // Procedure
		body := generator.generateStatement(decl.Body, index)
		closure := fmt.Sprintf("%sQ%s ENDPROC &%s&\n",
			index.Indentation(), index.String(), decl.Name.Name)
		return body + closure
	} else if stmts, status := statement.(parser.Statements); status {               // Statements
		var str string

		for _, stmt := range stmts {
			newStr := generator.generateStatement(stmt, index)

			if newStr != "" {
				str += generator.generateStatement(stmt, index)
				index[len(index)-1]++
			}
		}

		return str
	}

	return "!!!Error!!!"
}

func (generator *Generator) generateAssignStatement(assign parser.AssignStatement, index index) string {
	var expr string
	ident := fmt.Sprintf("%sQ%s %s",
		index.Indentation(), index.String(), generator.generateExpression(assign.Identifier))

	// Basic types should be on the next line like an answer.
	// But complex unary and binary expressions should be on the same line
	if lit, status := assign.Expression.(parser.Literal); status {
		expr = fmt.Sprintf("\n%sA%s %s\n",
			index.Indentation(), index.String(), generateLiteral(lit))
	} else {
		expr = fmt.Sprintf(" := %s\n",
			generator.generateExpression(assign.Expression))
	}

	return ident + expr
}

func (generator *Generator) generateExpression(expression parser.Expression) string {
	if expr, status := expression.(parser.UnaryExpression); status {
		operand := generator.generateExpression(expr.Operand)
		return expr.Operator + " " + operand
	} else if expr, status := expression.(parser.BinaryExpression); status {
		left := generator.generateExpression(expr.LeftOperand)
		right := generator.generateExpression(expr.RightOperand)
		return left + " " + expr.Operator + " " + right
	} else if lit, status := expression.(parser.Literal); status {
		return generateLiteral(lit)
	} else if ident, status := expression.(parser.Identifier); status {
		return "&" + ident.Name + "&"
	}

	return "!!!Error!!!"
}

func (generator *Generator) generateIfStatement(stmt parser.IfStatement, index index) string {
	ifBody := generator.generateStatement(stmt.IfBody, index)
	var elseBody string

	if stmt.ElseBody.Statements != nil {
		elseBody = fmt.Sprintf("%sQ%s END ELSE BEGIN\n",
			index.Indentation(), index.String())
		index[len(index)-1]++
		elseBody += generator.generateStatement(stmt.ElseBody, index)
	} else {
		elseBody = generator.generateStatement(stmt.ElseBody, index)
	}

	closure := fmt.Sprintf("%sQ%s END\n", index.Indentation(), index.String())
	return ifBody + elseBody + closure
}

func (generator *Generator) generateCaseStatements(
	parentExpr parser.Expression,
	stmts parser.CaseStatements,
	index index,
) string {
	var expr parser.Expression
	var ifs []parser.IfStatement
	var str string

	for _, stmt := range stmts {
		if !isExpressionNil(parentExpr) {
			expr = parser.BinaryExpression{
				LeftOperand: parentExpr, Operator: "==", RightOperand: stmt.Expression}
		} else { // Switch without condition
			expr = stmt.Expression
		}

		body := stmt.Body
		curIf := parser.IfStatement{
			Condition: expr, IfBody: body, ElseBody: parser.BlockStatement{}}

		if isExpressionNil(stmt.Expression){ // If 'default' case
			curIf.Condition = parser.Literal{Type: parser.BooleanLiteral, Value: true}
		}

		ifs = append(ifs, curIf)
	}


	for i := len(ifs) - 1; i > 0; i-- {
		tmpStmt := parser.BlockStatement{}
		tmpStmt.Statements = append(tmpStmt.Statements, ifs[i])
		ifs[i-1].ElseBody = tmpStmt
	}

	if len(ifs) > 0 {
		str += generator.generateStatement(ifs[0], index)
	}

	return str
}

func generateLiteral(literal parser.Literal) string {
	if str, status := literal.Value.(string); status {
		return "\"" + str + "\""
	}

	return fmt.Sprintf("%v", literal.Value)
}

func (index index) Indentation() string  {
	return strings.Repeat("\t", len(index) - 1)
}

func (index index) String() string {
	var str string

	for _, i := range index {
		str += strconv.Itoa(i) + "."
	}

	return str
}

func isExpressionNil(expression parser.Expression) bool {
	if expr, status := expression.(parser.UnaryExpression); status {
		if expr.Operand == nil {
			return true
		}
	}

	return false
}