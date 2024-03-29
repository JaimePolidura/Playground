package interpreter

import "interpreters/src/syntax"

type Interpreter struct {
	statements []syntax.Stmt

	environment *Environment
	global      *Environment //Used for functions,

	Log []string
}

func CreateInterpreter(statements []syntax.Stmt) *Interpreter {
	return &Interpreter{
		statements:  statements,
		environment: createRootEnvironment(),
		global:      createGlobalEnvironment(),
		Log:         make([]string, 0),
	}
}

func createGlobalEnvironment() *Environment {
	global := createRootEnvironment()
	global.variables["clock"] = ClockNativeFunction{}

	return global
}

func (i *Interpreter) Interpret() error {
	for _, statement := range i.statements {
		if err := i.interpretStatement(statement); err != nil {
			return err
		}
	}

	return nil
}
