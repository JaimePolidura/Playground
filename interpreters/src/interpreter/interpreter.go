package interpreter

import "interpreters/src/syntax"

type Interpreter struct {
	statements []syntax.Stmt

	environment *Environment

	Log []string
}

func CreateInterpreter(statements []syntax.Stmt) *Interpreter {
	return &Interpreter{
		statements:  statements,
		environment: createEnvironment(),
		Log:         make([]string, 0),
	}
}

func (i *Interpreter) Interpret() error {
	for _, statement := range i.statements {
		if err := i.interpretStatement(statement); err != nil {
			return err
		}
	}

	return nil
}
