package interpreter

import "interpreters/src/syntax"

type LoxFunction struct {
	FunctionStmt syntax.FunctionStatement
}

func (l LoxFunction) Call(interpreter *Interpreter, args []any) (syntax.Expr, error) {
	environmentForFunction := createGlobalEnvironment().CopyInto(interpreter.environment)

	for i := 0; i < len(l.FunctionStmt.Params); i++ {
		environmentForFunction.Define(l.FunctionStmt.Params[i].Lexeme, args[i])
	}

	return nil, interpreter.interpretBlockStmt(l.FunctionStmt.Body, environmentForFunction)
}

func (l LoxFunction) Arity() int {
	return len(l.FunctionStmt.Params)
}
