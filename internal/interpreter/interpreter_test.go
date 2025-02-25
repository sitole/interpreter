package interpreter

import (
	"github.com/sitole/interpreter/internal/ast"
	"github.com/sitole/interpreter/internal/tokenizer"
	"testing"
)

func TestDemo(t *testing.T) {
	programAst := ast.Program{
		Roots: []ast.Node{
			ast.VariableDefinition{
				Identifier: "xx",
				ValueType:  ast.VARIABLE_INT,
				Value: ast.MathExpression{
					Operation: ast.MATH_OPERATION_PLUS,
					Left: ast.LiteralDefinition{
						Type:  ast.VARIABLE_INT,
						Value: 33,
					},
					Right: ast.LiteralDefinition{
						Type:  ast.VARIABLE_INT,
						Value: 13,
					},
				},
			},

			ast.VariableDefinition{
				Identifier: "yy",
				ValueType:  ast.VARIABLE_INT,
				Value: ast.LiteralDefinition{
					Type:  ast.VARIABLE_INT,
					Value: 33,
				},
			},
		},
	}

	programError := func(err ast.ContextError) {
		t.Fatalf(tokenizer.ColorRed+"Error: %s (line %d, column %d)"+tokenizer.ColorReset, err.Message, err.Line, err.Column)
	}

	program := CreateProgramWithThrowback(programAst, programError)
	program.run()
}
