package ast

import (
	"github.com/sitole/interpreter/internal/tokenizer"
	"reflect"
	"strings"
	"testing"
)

func TestParser(t *testing.T) {
	code := "var xx = 1"
	tokens := []tokenizer.Token{
		{Type: tokenizer.TOKEN_VARIABLE_DEFINITION, Line: 1, Column: 1},
		{Type: tokenizer.TOKEN_VARIABLE_IDENTIFIER, Line: 1, Column: 5, Literal: "xx"},
		{Type: tokenizer.TOKEN_VARIABLE_ASSIGN, Line: 1, Column: 8},
		{Type: tokenizer.TOKEN_NUMBER, Line: 1, Column: 10, Literal: 1},
		{Type: tokenizer.TOKEN_EOL, Line: 1, Column: 9},
	}

	parser := BuildParser(tokens)
	ast, err := parser.Parse()
	if err != nil {
		errorPretty := ErrorFormatter(code, *err)
		for _, line := range errorPretty {
			t.Error(line)
		}

		t.Fatal()
	}

	astExpected := Program{
		Roots: []Node{
			VariableDefinition{
				PositionLine:   1,
				PositionColumn: 1,

				Identifier: "xx",
				ValueType:  VARIABLE_INT,
				Value: LiteralDefinition{
					PositionLine:   1,
					PositionColumn: 10,

					Type:  VARIABLE_INT,
					Value: 1,
				},
			},
		},
	}

	if !reflect.DeepEqual(astExpected, *ast) {
		t.Errorf("Expected %+v but got %+v", astExpected, *ast)
	}
}

func TestParserFromCode(t *testing.T) {
	code := "var xx = 1"
	codeLines := strings.Split(code, "\n")

	tokens, err := tokenizer.Tokenizer(codeLines)
	if err != nil {
		errorPretty := tokenizer.ErrorFormatter(code, *err)
		for _, line := range errorPretty {
			t.Error(line)
		}

		t.Fatal()
	}

	parser := BuildParser(tokens)
	ast, astErr := parser.Parse()
	if astErr != nil {
		errorPretty := ErrorFormatter(code, *astErr)
		for _, line := range errorPretty {
			t.Error(line)
		}

		t.Fatal()
	}

	astExpected := Program{
		Roots: []Node{
			VariableDefinition{
				PositionLine:   1,
				PositionColumn: 1,

				Identifier: "xx",
				ValueType:  VARIABLE_INT,
				Value: LiteralDefinition{
					PositionLine:   1,
					PositionColumn: 10,

					Type:  VARIABLE_INT,
					Value: 1,
				},
			},
		},
	}

	if !reflect.DeepEqual(astExpected, *ast) {
		t.Errorf("Expected %+v but got %+v", astExpected, *ast)
	}
}
