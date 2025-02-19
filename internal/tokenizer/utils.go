package tokenizer

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

func tokenizationErrorFormatter(code string, err TokenError) []string {
	codeLines := strings.Split(code, "\n")

	errorRedColor := "\033[31m"
	errorLine := codeLines[err.Line-1]
	errorSpaces := strings.Repeat(" ", err.Column-1)
	errorMessage := fmt.Sprintf(errorRedColor+"%s^ syntax error (line %d, column %d): %s", errorSpaces, err.Line, err.Column, err.Err)

	return []string{errorLine, errorMessage}
}

func tokenStringLiteral(t Token) string {
	return t.Literal.(string)
}

func tokenNumberLiteral(t Token) int {
	return t.Literal.(int)
}

func isVariableDefinition(lex string) bool {
	// support variables with numbers on +1 positions
	return regexp.MustCompile("^[A-Za-z_]*$").MatchString(lex)
}

func isInteger(s string) bool {
	_, err := strconv.Atoi(s)
	return err == nil
}

func firstNonAlphabetIndex(s string) int {
	for i, r := range s {
		if (r < 'A' || r > 'Z') && (r < 'a' || r > 'z') {
			return i // return the byte index of the first non-letter
		}
	}

	return -1 // return -1 if all characters are A-Z or a-z
}

func firstNonNumberIndex(s string) int {
	for i, _ := range s {
		_, err := strconv.Atoi(s[i : i+1])
		if err != nil {
			return i // return the byte index of the first non-number
		}
	}

	return len(s)
}
