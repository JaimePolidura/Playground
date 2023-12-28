package interpreter

import (
	"errors"
	"fmt"
	"interpreters/src/syntax"
)

func (i *Interpreter) interpretStatement(statement syntax.Stmt) error {
	switch statement.Type() {
	case syntax.PRINT:
		return i.interpretPrint(statement.(syntax.PrintStatement))
	}

	return errors.New("unhandled statement")
}

func (i *Interpreter) interpretPrint(statement syntax.PrintStatement) error {
	expr, err := interpretExpression(statement.Expression)
	if err != nil {
		return err
	}

	//TODO Lookup variables
	stringValueToPrint, err := castString(expr.(syntax.LiteralExpression).Literal)
	if err != nil {
		return err
	}

	i.Log = append(i.Log, stringValueToPrint)

	fmt.Println(stringValueToPrint)
	return nil
}
