package semantic

import "../parser"

type Analyzer struct {
	variables  Variables
	syntaxTree parser.File
	errors     parser.Errors
}

func NewAnalyzer(tree parser.File) *Analyzer {
	return &Analyzer{Variables{}, tree, parser.Errors{}}
}

func (analyzer *Analyzer) Analyze() (Variables, parser.Errors) {
	for _, function := range analyzer.syntaxTree.Declarations {
		stmt := function.(parser.FuncDeclaration).Body
		analyzer.traverseStatement(stmt, 0)
	}

	return analyzer.variables, analyzer.errors
}

func (analyzer *Analyzer) traverseStatement(statement parser.Statement, scope Scope) {
	if stmt, status := statement.(parser.SwitchStatement); status {         // Switch
		analyzer.getExpressionType(stmt.Expression, scope) // Validating condition
		analyzer.traverseStatement(stmt.Body, scope+1)
	} else if stmt, status := statement.(parser.CaseStatement); status {    // Case
		analyzer.getExpressionType(stmt.Expression, scope) // Validating condition
		analyzer.traverseStatement(stmt.Body, scope+1)
	} else if stmt, status := statement.(parser.AssignStatement); status {  // Assign
		analyzer.validateAssignStatement(stmt, scope)
	} else if stmt, status := statement.(parser.BlockStatement); status {   // Block
		analyzer.traverseStatement(stmt.Statements, scope)
	} else if stmts, status := statement.(parser.CaseStatements); status {  // Cases
		for _, stmt := range stmts {
			analyzer.traverseStatement(stmt, scope)
		}
	} else if stmts, status := statement.(parser.Statements); status {      // Statements
		for _, stmt := range stmts {
			analyzer.traverseStatement(stmt, scope)
		}
	} else if stmt, status := statement.(parser.IfStatement); status {      // If
		if analyzer.getExpressionType(stmt.Condition, scope) != Bool {
			analyzer.errors = append(analyzer.errors, newNonBoolError())
		}

		analyzer.traverseStatement(stmt.IfBody, scope+1)
		analyzer.traverseStatement(stmt.ElseBody, scope+1)
	}
}

func (analyzer *Analyzer) validateAssignStatement(assign parser.AssignStatement, scope Scope) {
	identifier := assign.Identifier

	if assign.Operator == parser.GetType(parser.Assign) {
		variable := analyzer.findVariableSomewhere(identifier, scope)

		if variable != nil {
			analyzer.assignVariable(variable, assign.Expression, scope)
		}
	} else if assign.Operator == parser.GetType(parser.Define) {
		if analyzer.findVariableAtScope(identifier, scope) != nil {
			analyzer.errors = append(analyzer.errors, newAlreadyDefinedError(identifier))
			return
		}

		analyzer.defineVariable(assign, scope)
	}
}

// Fairly bad name for a function which finds
// vars in scopes which are less or equal than current
func (analyzer *Analyzer) findVariableSomewhere(ident parser.Identifier, scope Scope) *Variable {
	scopes := SortMapKeys(analyzer.variables)

	// Iterating in reverse order is used because variable with highest scope
	// overlaps the one with the same name in lower scope.
	for index := len(scopes) - 1; index >= 0; index-- {
		otherScope := Scope(scopes[index])

		if scope < otherScope {
			continue
		}

		variable := analyzer.findVariableAtScope(ident, otherScope)

		if variable != nil {
			return variable
		}
	}

	analyzer.errors = append(analyzer.errors, newNotDefinedError(ident))
	return nil
}

func (analyzer *Analyzer) findVariableAtScope(ident parser.Identifier, scope Scope) *Variable {
	for _, variable := range analyzer.variables[scope] {
		if variable.Name == ident.Name {
			return variable
		}
	}

	return nil
}

func (analyzer *Analyzer) defineVariable(assign parser.AssignStatement, scope Scope) {
	identifier := assign.Identifier.Name
	varType := analyzer.getExpressionType(assign.Expression, scope)

	variable := &Variable{identifier, varType}
	analyzer.variables[scope] = append(analyzer.variables[scope], variable)
}

func (analyzer *Analyzer) assignVariable(
	variable *Variable,
	expression parser.Expression,
	scope Scope,
) bool {
	varType := analyzer.getExpressionType(expression, scope)

	if variable.Type != varType {
		analyzer.errors = append(analyzer.errors, newAssignError(variable, varType))
		return false
	}

	return true
}

func (analyzer *Analyzer) getExpressionType(expression parser.Expression, scope Scope) string {
	if expr, status := expression.(parser.UnaryExpression); status {
		return analyzer.getExpressionType(expr.Operand, scope)
	} else if expr, status := expression.(parser.BinaryExpression); status {
		return analyzer.getBinaryExpressionType(expr, scope)
	} else if lit, status := expression.(parser.Literal); status {
		return intToType(int(lit.Type))
	} else if ident, status := expression.(parser.Identifier); status {
		variable := analyzer.findVariableSomewhere(ident, scope)

		if variable == nil {
			return Undefined
		}

		return variable.Type
	}

	return Undefined
}

func (analyzer *Analyzer) getBinaryExpressionType(expr parser.BinaryExpression, scope Scope) string {
	left := analyzer.getExpressionType(expr.LeftOperand, scope)
	right := analyzer.getExpressionType(expr.RightOperand, scope)

	if left != right {
		if isNumber(left) && isNumber(right) {
			if isComparison(expr.Operator) {
				return Bool
			}

			// Non-comparison operation between int and float
			// results in float type
			return Float
		}

		analyzer.errors = append(analyzer.errors, newExpressionError(left, right))
		return Undefined
	} else { // Operand types are equal
		if isComparison(expr.Operator) {
			return Bool
		}

		return left
	}
}
