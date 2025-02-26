package ast

import (
	"fmt"
	"github.com/sitole/interpreter/internal/tokenizer"
)

type ParserError struct {
	Message string
	Line    int
	Column  int
}

type Parser struct {
	tokens []tokenizer.Token
	cursor int
}

func BuildParser(tokens []tokenizer.Token) Parser {
	return Parser{
		tokens: tokens,
		cursor: 0,
	}
}

func (p *Parser) Parse() (*Program, *ParserError) {
	program := Program{
		Roots: make([]Node, 0),
	}

	for _, token := range p.tokens {
		switch token.Type {
		case tokenizer.TOKEN_VARIABLE_DEFINITION:
			def, err := p.parseVariableDefinition()
			if err != nil {
				return nil, err
			}

			program.Roots = append(program.Roots, *def)
			break
		default:
			// todo: implement other types
			p.setCursorToNext()
			break
		}
	}

	return &program, nil
}

func (p *Parser) parseVariableDefinition() (*VariableDefinition, *ParserError) {
	current := p.currentToken()
	next := p.nextToken()

	if next == nil || next.Type != tokenizer.TOKEN_VARIABLE_IDENTIFIER {
		return nil, &ParserError{
			Message: "Expected variable identifier",
			Line:    current.Line,
			Column:  current.Column,
		}
	}

	variableName := next.Literal.(string)
	p.setCursorToNext()

	next = p.nextToken()
	if next == nil || next.Type != tokenizer.TOKEN_VARIABLE_ASSIGN {
		return nil, &ParserError{
			Message: "Expected variable assignment",
			Line:    current.Line,
			Column:  current.Column,
		}
	}

	p.setCursorToNext()
	next = p.nextToken()
	if next == nil {
		return nil, &ParserError{
			Message: "After assignment value is expected",
			Line:    current.Line,
			Column:  current.Column,
		}
	}

	switch next.Type {
	case tokenizer.TOKEN_NUMBER:
		return &VariableDefinition{
			PositionLine:   current.Line,
			PositionColumn: current.Column,

			Identifier: variableName,
			ValueType:  VARIABLE_INT,
			Value: LiteralDefinition{
				PositionLine:   next.Line,
				PositionColumn: next.Column,

				Type:  VARIABLE_INT,
				Value: next.Literal.(int),
			},
		}, nil
	case tokenizer.TOKEN_STRING:
		return &VariableDefinition{
			PositionLine:   current.Line,
			PositionColumn: current.Column,

			Identifier: variableName,
			ValueType:  VARIABLE_STRING,
			Value: LiteralDefinition{
				PositionLine:   next.Line,
				PositionColumn: next.Column,

				Type:  VARIABLE_STRING,
				Value: next.Literal.(string),
			},
		}, nil
	default:
		return nil, &ParserError{
			Message: fmt.Sprintf("Unexpected token when variable parsing %s", next.Type),
			Line:    next.Line,
			Column:  next.Column,
		}
	}
}

func (p *Parser) currentToken() tokenizer.Token {
	return p.tokens[p.cursor]
}

func (p *Parser) nextToken() *tokenizer.Token {
	if p.cursor+1 < len(p.tokens) {
		return &p.tokens[p.cursor+1]
	}

	return nil
}

func (p *Parser) setCursorToNext() {
	p.cursor = p.cursor + 1
}
