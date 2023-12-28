package interpreter

import (
	"errors"
	"fmt"
	"interpreters/src/syntax"
)

func (i *Interpreter) interpretStatement(statement syntax.Stmt) error {
	switch statement.Type() {
	case syntax.PRINT_STMT:
		return i.interpretPrintStmt(statement.(syntax.PrintStatement))
	case syntax.VAR_STMT:
		return i.interpretVarStmt(statement.(syntax.VarStatement))
	case syntax.EXPRESSION_STMT:
		return i.interpretExprStmt(statement.(syntax.ExpressionStatement))
	case syntax.BLOCK_STMT:
		return i.interpretBlockStmt(statement.(syntax.BlockStatement), i.environment)
	}

	return errors.New("unhandled statement")
}

func (i *Interpreter) interpretBlockStmt(blockStatement syntax.BlockStatement, environmentParent *Environment) error {
	newEnvironment := createChildEnvironment(environmentParent)
	prevEnvironment := i.environment
	i.environment = newEnvironment
	for _, statement := range blockStatement.Statements {
		if err := i.interpretStatement(statement); err != nil {
			return err
		}
	}

	i.environment = prevEnvironment
	return nil
}

func (i *Interpreter) interpretExprStmt(statement syntax.ExpressionStatement) error {
	_, err := i.interpretExpression(statement.Expression)
	return err
}

func (i *Interpreter) interpretVarStmt(statement syntax.VarStatement) error {
	initializer := statement.Initializer
	if initializer != nil {
		if initializerInterpreted, err := i.interpretExpression(initializer); err != nil {
			return err
		} else {
			initializer = initializerInterpreted
		}
	}

	if initializer.Type() != syntax.LITERAL_EXPR {
		return errors.New("invalid initializer variable. It should return a literal")
	}

	i.environment.Define(statement.Name.Literal.(string), initializer.(syntax.LiteralExpression).Literal)

	return nil
}

func (i *Interpreter) interpretPrintStmt(statement syntax.PrintStatement) error {
	expr, err := i.interpretExpression(statement.Expression)
	if err != nil {
		return err
	}

	stringValueToPrint, err := castString(expr.(syntax.LiteralExpression).Literal)
	if err != nil {
		return err
	}

	i.Log = append(i.Log, stringValueToPrint)

	fmt.Println(stringValueToPrint)
	return nil
}
