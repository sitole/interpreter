package interpreter

import (
	"github.com/sitole/interpreter/internal/ast"
)

type Environment struct {
	Variables      map[string]interface{}
	Error          *ast.ContextError
	ErrorThrowback *func(err ast.ContextError)
}

func (e Environment) GetVariableValue(key string) interface{} {
	keyValue, keyFound := e.Variables[key]
	if keyFound {
		return keyValue
	}
	return false
}

func (e Environment) SetVariableValue(key string, value interface{}) {
	e.Variables[key] = value
}

func (e Environment) SetError(err ast.ContextError) {
	e.Error = &err

	if e.ErrorThrowback != nil {
		(*e.ErrorThrowback)(err)
	}
}

type Program struct {
	env     Environment
	program ast.Program
}

func CreateProgramWithThrowback(ast ast.Program, throwback func(err ast.ContextError)) Program {
	return Program{
		program: ast,
		env: Environment{
			Variables:      make(map[string]interface{}),
			ErrorThrowback: &throwback,
			Error:          nil,
		},
	}
}

func (p Program) run() {
	for _, node := range p.program.Roots {
		node.Visitor(p.env)
	}
}
