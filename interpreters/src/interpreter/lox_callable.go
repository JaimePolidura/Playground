package interpreter

import "interpreters/src/syntax"

type LoxCallable interface {
	Call(interpreter *Interpreter, args []any) (syntax.Expr, error)
	Arity() int
}
