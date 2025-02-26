package tokenizer

import (
	"reflect"
	"strings"
	"testing"
)

func TestVariableDefinition(t *testing.T) {
	code := "var xx = 1"
	codeLines := strings.Split(code, "\n")

	tokens, err := Tokenizer(codeLines)
	if err != nil {
		errorPretty := ErrorFormatter(code, *err)
		for _, line := range errorPretty {
			t.Error(line)
		}

		t.Fatal()
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

func TestMultipleVariablesDefinition(t *testing.T) {
	code := "var xx = 1\nvar yy = xx + 10"

	tokens, err := Tokenizer(strings.Split(code, "\n"))
	if err != nil {
		errorPretty := ErrorFormatter(code, *err)
		for _, line := range errorPretty {
			t.Error(line)
		}

		t.Fatal()
	}

	tokensExpected := []Token{
		// line 1
		{Type: TOKEN_VARIABLE_DEFINITION, Line: 1, Column: 1},
		{Type: TOKEN_VARIABLE_IDENTIFIER, Line: 1, Column: 5, Literal: "xx"},
		{Type: TOKEN_VARIABLE_ASSIGN, Line: 1, Column: 8},
		{Type: TOKEN_NUMBER, Line: 1, Column: 10, Literal: 1},
		{Type: TOKEN_EOL, Line: 1, Column: 9},

		// line 2
		{Type: TOKEN_VARIABLE_DEFINITION, Line: 2, Column: 1},
		{Type: TOKEN_VARIABLE_IDENTIFIER, Line: 2, Column: 5, Literal: "yy"},
		{Type: TOKEN_VARIABLE_ASSIGN, Line: 2, Column: 8},
		{Type: TOKEN_VARIABLE_IDENTIFIER, Line: 2, Column: 10, Literal: "xx"},
		{Type: TOKEN_MATH_PLUS, Line: 2, Column: 13},
		{Type: TOKEN_NUMBER, Line: 2, Column: 15, Literal: 10},
		{Type: TOKEN_EOL, Line: 2, Column: 15},
	}

	if !reflect.DeepEqual(tokensExpected, tokens) {
		t.Errorf("Expected %+v but got %+v", tokensExpected, tokens)
	}
}

func TestVariableDefinition_UnknownTokenError(t *testing.T) {
	code := "var xx != 1"
	codeLines := strings.Split(code, "\n")
	_, err := Tokenizer(codeLines)

	if err == nil {
		t.Fatalf("Expected an error but got nil")
	}

	expectedError := TokenError{
		Err:    "not expected token received",
		Line:   1,
		Column: 8,
	}

	if !reflect.DeepEqual(&expectedError, err) {
		t.Errorf("Expected %+v but got %+v", expectedError, err)
	}
}

func TestVariableDefinition_UnknownTokenErrorFormatted(t *testing.T) {
	code := "var xx != 1"
	codeLines := strings.Split(code, "\n")
	_, err := Tokenizer(codeLines)

	if err == nil {
		t.Fatalf("Expected an error but got nil")
	}

	errorLines := ErrorFormatter(code, *err)
	errorLinesExpected := []string{
		"var xx != 1",
		ColorRed + "       ^ syntax error (line 1, column 8): not expected token received" + ColorReset,
	}

	if !reflect.DeepEqual(errorLinesExpected, errorLines) {
		t.Errorf("Expected %+v but got %+v", errorLinesExpected, errorLines)
	}
}
