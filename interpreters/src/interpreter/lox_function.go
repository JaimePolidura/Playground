package interpreter

import "interpreters/src/syntax"

type LoxReturn struct {
	Value syntax.Expr
}

func (l LoxReturn) Error() string { return "" }

func isErrTypeOfLoxReturn(err error) (bool, LoxReturn) {
	if err == nil {
		return false, LoxReturn{}
	}

	if loxReturn, isLoxReturn := err.(LoxReturn); isLoxReturn {
		return true, loxReturn
	} else {
		return false, LoxReturn{}
	}
}

type LoxFunction struct {
	FunctionStmt syntax.FunctionStatement
}

func (l LoxFunction) Call(interpreter *Interpreter, args []any) (syntax.Expr, error) {
	environmentForFunction := createGlobalEnvironment().CopyInto(interpreter.environment)

	for i := 0; i < len(l.FunctionStmt.Params); i++ {
		environmentForFunction.Define(l.FunctionStmt.Params[i].Lexeme, args[i])
	}

	err := interpreter.interpretBlockStmt(l.FunctionStmt.Body, environmentForFunction)

	if isLoxReturn, loxReturn := isErrTypeOfLoxReturn(err); isLoxReturn {
		return loxReturn.Value, nil
	} else {
		return nil, err
	}
}

func (l LoxFunction) Arity() int {
	return len(l.FunctionStmt.Params)
}
