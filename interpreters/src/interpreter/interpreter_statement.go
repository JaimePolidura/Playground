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
		return i.interpretBlockStmt(statement.(syntax.BlockStatement).Statements, i.environment)
	case syntax.IF_STMT:
		return i.interpretIfStmt(statement.(syntax.IfStatement))
	case syntax.WHILE_STMT:
		return i.interpretWhileStmt(statement.(syntax.WhileStatement))
	case syntax.FUNCTION_STMT:
		return i.interpretFunctionStmt(statement.(syntax.FunctionStatement))
	case syntax.RETURN_STMT:
		return i.interpretReturnStmt(statement.(syntax.ReturnStatement))
	case syntax.CLASS_STMT:
		return i.interpretClassStmt(statement.(syntax.ClassStatement))
	}

	return errors.New("unhandled statement")
}

func (i *Interpreter) interpretClassStmt(statement syntax.ClassStatement) error {
	className := statement.Name.Lexeme
	superClassName := statement.Superclass
	class := &LoxClass{Name: className, Methods: make(map[string]LoxFunction)}

	if superClassName != "" { //Has superclass
		superClassAny, err := i.environment.Get(superClassName)
		if err != nil {
			return errors.New("class " + superClassName + " not found")
		}
		superClass, isSuperClassType := superClassAny.(*LoxClass)
		if !isSuperClassType {
			return errors.New("variable names/methods cannot have the same name as classes")
		}

		class.SuperClass = superClass
	}

	for _, method := range statement.Methods {
		class.Methods[method.Name.Lexeme] = LoxFunction{FunctionStmt: method}
	}
	i.environment.Define(className, class)

	return nil
}

func (i *Interpreter) interpretReturnStmt(statement syntax.ReturnStatement) error {
	value, err := i.interpretExpression(statement.Value)
	if err != nil {
		return err
	}

	return LoxReturn{Value: value}
}

func (i *Interpreter) interpretFunctionStmt(functionStmt syntax.FunctionStatement) error {
	loxFunction := LoxFunction{FunctionStmt: functionStmt}
	i.environment.Define(functionStmt.Name.Lexeme, loxFunction)
	return nil
}

func (i *Interpreter) interpretWhileStmt(statement syntax.WhileStatement) error {
	for {
		conditionBool, err := i.interpretExprAndGetBool(statement.Condition)
		if err != nil {
			return err
		}
		if !conditionBool {
			return nil
		}

		i.interpretStatement(statement.Body)
	}
}

func (i *Interpreter) interpretIfStmt(ifStmt syntax.IfStatement) error {
	boolIfCondition, err := i.interpretExprAndGetBool(ifStmt.Condition)
	if err != nil {
		return errors.New("the result of a if statement should yield a boolean value")
	}

	if boolIfCondition {
		return i.interpretStatement(ifStmt.ThenBranch)
	} else if ifStmt.ElseBranch != nil {
		return i.interpretStatement(ifStmt.ElseBranch)
	}

	return nil
}

func (i *Interpreter) interpretBlockStmt(blockStatements []syntax.Stmt, environmentParent *Environment) error {
	newEnvironment := createChildEnvironment(environmentParent)
	prevEnvironment := i.environment
	i.environment = newEnvironment
	var errorToReturn error
	for _, statement := range blockStatements {
		if err := i.interpretStatement(statement); err != nil {
			errorToReturn = err
			break
		}
	}

	i.environment = prevEnvironment

	return errorToReturn
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
