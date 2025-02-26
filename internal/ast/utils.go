package ast

import (
	"fmt"
	"github.com/sitole/interpreter/internal/tokenizer"
	"strings"
)

func ErrorFormatter(code string, err ParserError) []string {
	codeLines := strings.Split(code, "\n")

	errorLine := codeLines[err.Line-1]
	errorSpaces := strings.Repeat(" ", err.Column-1)
	errorMessage := fmt.Sprintf(tokenizer.ColorRed+"%s^ syntax error (line %d, column %d): %s"+tokenizer.ColorReset, errorSpaces, err.Line, err.Column, err.Message)

	return []string{errorLine, errorMessage}
}
