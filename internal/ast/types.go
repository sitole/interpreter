package ast

type ContextError struct {
	Message string
	Line    int
	Column  int
}

type Context interface {
	GetVariableValue(key string) interface{}
	SetVariableValue(key string, value interface{})

	SetError(err ContextError)
}

type Node interface {
	Visitor(ctx Context)
	Position() (int, int)
}

// todo
var RESULT_REGISTER = "RESULT_REGISTER"

var (
	MATH_OPERATION_MINUS = "MATH_OPERATION_MINUS"
	MATH_OPERATION_PLUS  = "MATH_OPERATION_PLUS"
	MATH_OPERATION_DIV   = "MATH_OPERATION_DIV"
	MATH_OPERATION_MUL   = "MATH_OPERATION_MUL"
)

type MathExpression struct {
	PositionLine   int
	PositionColumn int

	Operation string
	Left      Node
	Right     Node
}

func (m MathExpression) Visitor(ctx Context) {
	mathResult := 0

	// todo: check state
	m.Left.Visitor(ctx)
	leftLiteral := ctx.GetVariableValue(RESULT_REGISTER).(int)

	// todo: check state
	m.Right.Visitor(ctx)
	rightLiteral := ctx.GetVariableValue(RESULT_REGISTER).(int)

	if m.Operation == MATH_OPERATION_PLUS {
		mathResult = leftLiteral + rightLiteral
	} else if m.Operation == MATH_OPERATION_MINUS {
		mathResult = leftLiteral - rightLiteral
	} else if m.Operation == MATH_OPERATION_MUL {
		mathResult = leftLiteral * rightLiteral
	} else if m.Operation == MATH_OPERATION_DIV {
		if rightLiteral == 0 {
			ctx.SetError(
				ContextError{
					Message: "Division by zero",
					Line:    m.PositionLine,
					Column:  m.PositionColumn,
				},
			)
		}

		mathResult = leftLiteral / rightLiteral
	} else {
		panic("Unknown operation")
	}

	ctx.SetVariableValue(RESULT_REGISTER, mathResult)
}

func (m MathExpression) Position() (int, int) {
	return m.PositionLine, m.PositionColumn
}

var (
	VARIABLE_INT    = "VARIABLE_INT"
	VARIABLE_STRING = "VARIABLE_STRING"
)

type LiteralDefinition struct {
	PositionLine   int
	PositionColumn int

	Type  string
	Value interface{}
}

func (l LiteralDefinition) Visitor(ctx Context) {
	ctx.SetVariableValue(RESULT_REGISTER, l.Value)
}

func (l LiteralDefinition) Position() (int, int) {
	return l.PositionLine, l.PositionColumn
}

type VariableDefinition struct {
	PositionLine   int
	PositionColumn int

	Identifier string
	ValueType  string
	Value      Node
}

func (v VariableDefinition) Visitor(ctx Context) {
	/*if ctx.GetVariableValue(v.Identifier) != nil {
		ctx.SetError(ContextError{Message: "Variable already defined"})
		return
	}*/

	v.Value.Visitor(ctx)
	variableValue := ctx.GetVariableValue(RESULT_REGISTER)
	ctx.SetVariableValue(v.Identifier, variableValue)
}

func (v VariableDefinition) Position() (int, int) {
	return v.PositionLine, v.PositionColumn
}

type Program struct {
	Roots []Node
}
