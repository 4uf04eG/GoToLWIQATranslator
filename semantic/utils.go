package semantic

import (
	"../parser"
	"sort"
)

func SortMapKeys(m Variables) []int {
	keys := make([]int, 0, len(m))

	for k := range m {
		keys = append(keys, int(k))
	}

	sort.Ints(keys)

	return keys
}

func isComparison(operator string) bool {
	switch operator {
	case parser.GetType(parser.Eq), parser.GetType(parser.Geq), parser.GetType(parser.Leq),
		parser.GetType(parser.Greater), parser.GetType(parser.Less):
		return true
	default:
		return false
	}
}

func isNumber(exprType string) bool {
	return exprType == Float || exprType == Int
}

func newAssignError(variable *Variable, realType string) *parser.Error {
	message := "Cannot use type " + realType + " in variable '" +
		variable.Name + "' of type " + variable.Type
	return &parser.Error{Type: parser.AssignError, Message: message}
}

func newExpressionError(leftType string, rightType string) *parser.Error {
	message := "Mismatched types: " + leftType + " and " + rightType
	return &parser.Error{Type: parser.MismatchedTypesError, Message: message}
}

func newNonBoolError() *parser.Error {
	msg := "Non-bool type used as condition"
	return &parser.Error{Type: parser.MismatchedTypesError, Message: msg}
}

func newAlreadyDefinedError(identifier parser.Identifier) *parser.Error {
	msg := "Variable '" + identifier.Name + "' is already defined"
	return &parser.Error{Type: parser.AssignError, Message: msg}
}

func newNotDefinedError(identifier parser.Identifier) *parser.Error {
	msg := "Variable '" + identifier.Name + "' is not defined"
	return &parser.Error{Type: parser.AssignError, Message: msg}
}
