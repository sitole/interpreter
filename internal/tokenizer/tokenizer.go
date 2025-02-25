package tokenizer

import (
	"fmt"
	"strconv"
)

var (
	TOKEN_MATH_MINUS       = "TOKEN_MATH_MINUS"
	TOKEN_MATH_PLUS        = "TOKEN_MATH_PLUS"
	TOKEN_MATH_EQUAL       = "TOKEN_MATH_EQUAL"
	TOKEN_MATH_EQUAL_EQUAL = "TOKEN_MATH_EQUAL_EQUAL"

	TOKEN_SEMICOLON = "TOKEN_SEMICOLON"
	TOKEN_STRING    = "TOKEN_STRING"
	TOKEN_NUMBER    = "TOKEN_NUMBER"

	TOKEN_FUNCTION_CALL       = "TOKEN_FUNCTION_CALL"
	TOKEN_VARIABLE_DEFINITION = "TOKEN_VARIABLE_DEFINITION"
	TOKEN_VARIABLE_IDENTIFIER = "TOKEN_VARIABLE_IDENTIFIER"
	TOKEN_VARIABLE_ASSIGN     = "TOKEN_VARIABLE_ASSIGN"

	TOKEN_EOL = "TOKEN_EOL"
)

type TokenError struct {
	Err    string
	Line   int
	Column int
}

type Token struct {
	Type    string
	Line    int
	Column  int
	Literal interface{} // value for strings, numbers etc
}

func Tokenizer(lines []string) ([]Token, *TokenError) {
	tokens := make([]Token, 0)

	for lineIndex, line := range lines {
		lineNum := lineIndex + 1
		lineRunes := []rune(line)

		for i := 0; i < len(lineRunes); i++ {
			c := string(lineRunes[i])
			columnNum := i + 1

			// ---
			// Skipping spaces
			// ---
			if c == " " {
				continue

				// ---
				// Variable assignation
				// ---
			} else if c == "=" {
				_, err := lastToken(tokens)
				if err != nil {
					return nil, &TokenError{
						Err:    fmt.Sprintf("previous token expected to be variable identifier, got nothing"),
						Line:   lineNum,
						Column: columnNum,
					}
				}

				tokens = append(
					tokens,
					Token{
						Type:   TOKEN_VARIABLE_ASSIGN,
						Line:   lineNum,
						Column: columnNum,
					},
				)
			} else if c == "+" {
				tokens = append(
					tokens,
					Token{
						Type:   TOKEN_MATH_PLUS,
						Line:   lineNum,
						Column: columnNum,
					},
				)

				// ---
				// Variable definition
				// ---
			} else if c == "v" {
				// todo: what if var identifier will start with "war"? Will match this case
				varSequence := "var"
				if !matchSequence(lineRunes, i, varSequence) {
					return nil, &TokenError{
						Err:    fmt.Sprintf("unexpected token found"),
						Line:   lineNum,
						Column: columnNum,
					}
				}

				// check if previous token is new line or empty (var can only be first on the line)
				lastToken, err := lastToken(tokens)
				if err == nil && lastToken.Type != TOKEN_EOL {
					return nil, &TokenError{
						Err:    fmt.Sprintf("unexpected token found, var definitino needs to be first token on the line"),
						Line:   lineNum,
						Column: columnNum,
					}
				}

				// check if there is a space after var
				if string(lineRunes[i+len(varSequence)]) != " " {
					return nil, &TokenError{
						Err:    fmt.Sprintf("unexpected token found, space expected after var is missing"),
						Line:   lineNum,
						Column: columnNum + len(varSequence),
					}
				}

				tokens = append(
					tokens,
					Token{
						Type:   TOKEN_VARIABLE_DEFINITION,
						Line:   lineNum,
						Column: columnNum,
					},
				)

				i += len(varSequence) - 1

				// ---
				// Number literal
				// ---
			} else if isInteger(c) {
				_, err := lastToken(tokens)
				if err != nil {
					return nil, &TokenError{
						Err:    fmt.Sprintf("expecting prevoius token, got nothing"),
						Line:   lineNum,
						Column: columnNum,
					}
				}

				numberLiteralEnd := firstNonNumberIndex(string(lineRunes[i:]))
				if numberLiteralEnd == -1 {
					return nil, &TokenError{
						Err:    fmt.Sprintf("number literal not valid"),
						Line:   lineNum,
						Column: columnNum,
					}
				}

				numberLiteral, err := strconv.Atoi(string(lineRunes[i : i+numberLiteralEnd]))
				if err != nil {
					return nil, &TokenError{
						Err:    fmt.Sprintf("number literal not valid during parsing"),
						Line:   lineNum,
						Column: columnNum,
					}
				}

				tokens = append(
					tokens,
					Token{
						Type:    TOKEN_NUMBER,
						Line:    lineNum,
						Column:  columnNum,
						Literal: numberLiteral,
					},
				)

				i += numberLiteralEnd - 1

				// ---
				// String literal
				// ---
			} else if c == "\"" {
				panic("string literal, not implemented yet")

				// ---
				// Variable names
				// ---
			} else {
				last, err := lastToken(tokens)
				if err != nil {
					return nil, &TokenError{
						Err:    fmt.Sprintf("unexpected token found"),
						Line:   lineNum,
						Column: columnNum,
					}
				}

				if last.Type == TOKEN_VARIABLE_DEFINITION || last.Type == TOKEN_VARIABLE_ASSIGN {
					variableNameDefEnd := firstNonAlphabetIndex(string(lineRunes[i:]))
					if variableNameDefEnd == -1 {
						return nil, &TokenError{
							Err:    fmt.Sprintf("variable identifier not valid"),
							Line:   lineNum,
							Column: columnNum,
						}
					}

					variableName := string(lineRunes[i : i+variableNameDefEnd])
					tokens = append(
						tokens,
						Token{
							Type:    TOKEN_VARIABLE_IDENTIFIER,
							Line:    lineNum,
							Column:  columnNum,
							Literal: variableName,
						},
					)

					i += variableNameDefEnd - 1
				} else {
					return nil, &TokenError{
						Err:    fmt.Sprintf("not expected token received"),
						Line:   lineNum,
						Column: columnNum,
					}
				}
			}
		}

		tokens = append(
			tokens,
			Token{
				Type:   TOKEN_EOL,
				Line:   lineNum,
				Column: len(lineRunes) - 1,
			},
		)
	}

	return tokens, nil
}

func lastToken(tokens []Token) (*Token, error) {
	if len(tokens) == 0 {
		return nil, fmt.Errorf("no token found")
	}

	token := tokens[len(tokens)-1]
	return &token, nil
}

func matchSequence(runes []rune, runesStartIndex int, sequence string) bool {
	runes = runes[runesStartIndex:]

	for i := 0; i < len(sequence); i++ {
		if runes[i] != rune(sequence[i]) {
			return false
		}
	}

	return true
}
