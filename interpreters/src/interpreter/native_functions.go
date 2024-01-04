package interpreter

import (
	"interpreters/src/syntax"
	"time"
)

type ClockNativeFunction struct{}

func (c ClockNativeFunction) Call(interpreter *Interpreter, args []any) (syntax.Expr, error) {
	return syntax.CreateLiteralExpression(time.Now().UnixMilli()), nil
}

func (c ClockNativeFunction) Arity() int {
	return 0
}
