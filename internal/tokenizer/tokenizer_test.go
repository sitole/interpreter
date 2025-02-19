package tokenizer

import (
	"reflect"
	"strings"
	"testing"
)

func TestVariableDefinition(t *testing.T) {
	code := "var xx = 1"
	codeLines := strings.Split(code, "\n")

	tokens, err := tokenizer(codeLines)
	if err != nil {
		errorPretty := tokenizationErrorFormatter(code, *err)
		for _, line := range errorPretty {
			t.Error(line)
		}
	}

	tokensExpected := []Token{
		{Type: TOKEN_VARIABLE_DEFINITION, Line: 1, Column: 1},
		{Type: TOKEN_VARIABLE_IDENTIFIER, Line: 1, Column: 5, Literal: "xx"},
		{Type: TOKEN_VARIABLE_ASSIGN, Line: 1, Column: 8},
		{Type: TOKEN_NUMBER, Line: 1, Column: 10, Literal: 1},
		{Type: TOKEN_EOL, Line: 1, Column: 9},
	}

	if !reflect.DeepEqual(tokensExpected, tokens) {
		t.Errorf("Expected %+v but got %+v", tokensExpected, tokens)
	}
}
