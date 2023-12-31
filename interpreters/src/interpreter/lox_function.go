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
	FunctionStmt      syntax.FunctionStatement
	BindedLoxInstance *LoxInstance
	ThisBinded        bool
}

func (l LoxFunction) BindThis(instanceToBind *LoxInstance) LoxFunction {
	l.BindedLoxInstance = instanceToBind
	l.ThisBinded = true
	return l
}

func (l LoxFunction) Call(interpreter *Interpreter, args []any) (syntax.Expr, error) {
	mostParentEnvironment := interpreter.environment.getMostParent()
	functionEnvironment := createChildEnvironment(mostParentEnvironment)
	prevEnv := interpreter.environment

	if l.ThisBinded {
		functionEnvironment.Define("this", l.BindedLoxInstance)
	}

	for i := 0; i < len(l.FunctionStmt.Params); i++ {
		functionEnvironment.Define(l.FunctionStmt.Params[i].Lexeme, args[i])
	}

	err := interpreter.interpretBlockStmt(l.FunctionStmt.Body, functionEnvironment)

	interpreter.environment = prevEnv

	if isLoxReturn, loxReturn := isErrTypeOfLoxReturn(err); isLoxReturn {
		return loxReturn.Value, nil
	} else {
		return nil, err
	}
}

func (l LoxFunction) Arity() int {
	return len(l.FunctionStmt.Params)
}
