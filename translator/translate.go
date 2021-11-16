package translator

import (
	"../generator"
	"../lexer"
	"../parser"
	"../semantic"
)

func Translate(code string) (genCode string) {
	tokens := lexer.NewLexer(code).Tokenize()
	ast := parser.NewParser(tokens).Parse()
	_, semErr := semantic.NewAnalyzer(ast).Analyze()
	parseErr := ast.Errors

	if len(semErr) > 0 || len(parseErr) > 0 {
		if len(parseErr) > 0 {
			genCode = "Syntax errors:\n" + parseErr.String()
		}
		if len(semErr) > 0 {
			genCode += "Semantic errors:\n" + semErr.String()
		}
	} else {
		genCode = generator.NewGenerator(ast).Generate()
	}
	return
}
