package interpreter

import (
	"errors"
	"interpreters/src/lex"
	"interpreters/src/syntax"
)

func (i *Interpreter) interpretExpression(rootExpression syntax.Expr) (syntax.Expr, error) {
	return i.interpretRecursiveExpression(rootExpression)
}

func (i *Interpreter) interpretRecursiveExpression(rootExpression syntax.Expr) (syntax.Expr, error) {
	switch rootExpression.Type() {
	case syntax.BINARY_EXPR:
		return i.interpretBinaryExpression(rootExpression.(syntax.BinaryExpression))
	case syntax.UNARY_EXPR:
		return i.interpretUnaryExpression(rootExpression.(syntax.UnaryExpression))
	case syntax.GROUPING_EXPR:
		return i.interpretGroupingExpression(rootExpression.(syntax.GroupingExpression))
	case syntax.VARIABLE_EXPR:
		return i.interpretVariableExpression(rootExpression.(syntax.VariableExpression))
	case syntax.LITERAL_EXPR:
		return rootExpression, nil
	case syntax.ASSIGN_EXPR:
		return i.interpretAssignExpression(rootExpression.(syntax.AssignExpression))
	}

	return nil, nil
}

func (i *Interpreter) interpretAssignExpression(assignExpression syntax.AssignExpression) (syntax.Expr, error) {
	valueExpr, err := i.interpretExpression(assignExpression.Value)
	if err != nil {
		return nil, err
	}

	if valueExpr.Type() != syntax.LITERAL_EXPR {
		return nil, errors.New("unexpected assigment expression type")
	}

	value := valueExpr.(syntax.LiteralExpression).Literal

	if err = i.environment.Assign(assignExpression.Name, value); err != nil {
		return nil, err
	} else {
		return syntax.CreateLiteralExpression(value), nil //We return the assigned value
	}
}

func (i *Interpreter) interpretUnaryExpression(unaryExpression syntax.UnaryExpression) (syntax.Expr, error) {
	interpretedExpression, err := i.interpretRecursiveExpression(unaryExpression.Right)
	if err != nil {
		return nil, err
	}

	if interpretedExpression.Type() != syntax.LITERAL_EXPR {
		unaryExpression.Right = interpretedExpression
		return unaryExpression, nil
	}

	unaryResult, err := i.calculateUnaryExpression(interpretedExpression.(syntax.LiteralExpression).Literal, unaryExpression.Token.Type)
	if err != nil {
		return nil, err
	}

	return syntax.CreateLiteralExpression(unaryResult), nil
}

func (i *Interpreter) interpretGroupingExpression(expression syntax.GroupingExpression) (syntax.Expr, error) {
	return i.interpretRecursiveExpression(expression.OtherExpression)
}

func (i *Interpreter) interpretVariableExpression(expression syntax.VariableExpression) (syntax.Expr, error) {
	variableValue, err := i.environment.Get(expression.Name.Literal.(string))
	if err != nil {
		return nil, err
	}

	return syntax.CreateLiteralExpression(variableValue), nil
}

func (i *Interpreter) interpretBinaryExpression(expression syntax.BinaryExpression) (syntax.Expr, error) {
	left, errLeft := i.interpretRecursiveExpression(expression.Left)
	right, errRight := i.interpretRecursiveExpression(expression.Right)
	if errLeft != nil {
		return nil, errLeft
	}
	if errRight != nil {
		return nil, errRight
	}

	if !i.isLiteral(left, right) {
		expression.Left = left
		expression.Right = right
		return expression, nil
	}

	literalInterpretedResult, err := i.calculateBinaryExpression(left.(syntax.LiteralExpression).Literal,
		right.(syntax.LiteralExpression).Literal, expression.Token.Type)
	if err != nil {
		return nil, err
	}

	return syntax.CreateLiteralExpression(literalInterpretedResult), nil
}

func (i *Interpreter) calculateUnaryExpression(literal any, tokenType lex.TokenType) (any, error) {
	switch tokenType {
	case lex.BANG:
		if castedBool, err := castBoolean(literal); err != nil {
			return 0, err
		} else {
			return !castedBool, nil
		}
	default:
		return false, errors.New("Unhandled unary operation: " + string(tokenType))
	}
}

func (i *Interpreter) calculateBinaryExpression(left any, right any, operationTokenType lex.TokenType) (any, error) {
	if lex.IsAnyType(operationTokenType, lex.OPEN_PAREN, lex.CLOSE_PAREN, lex.OPEN_BRACE, lex.CLOSE_BRACE, lex.COMMA, lex.DOT, lex.SEMICOLON,
		lex.IDENTIFIER, lex.STRING, lex.NUMBER, lex.CLASS, lex.ELSE, lex.FUN, lex.IF, lex.FOR, lex.NIL, lex.PRINT,
		lex.RETURN, lex.SUPER, lex.THIS, lex.TRUE, lex.VAR, lex.WHILE, lex.EOF, lex.BANG, lex.FALSE) {
		return nil, errors.New(string(operationTokenType) + " cannot be used as an operation")
	}

	isComparativeOperation := lex.IsAnyType(operationTokenType, lex.LESS, lex.BANG_EQUAL, lex.LESS_EQUAL, lex.EQUAL_EQUAL, lex.GREATER_EQUAL, lex.GREATER)
	isArithmeticOperation := lex.IsAnyType(operationTokenType, lex.MINUS, lex.PLUS, lex.SLASH, lex.STAR)
	isLogicalOperation := lex.IsAnyType(operationTokenType, lex.OR, lex.AND)

	if !isComparativeOperation && !isArithmeticOperation && !isLogicalOperation {
		return nil, errors.New("Unhandled token: " + string(operationTokenType))
	}

	if isComparativeOperation {
		return i.calculateComparativeOperation(left, right, operationTokenType)
	} else if isArithmeticOperation {
		return i.calculateArithmeticOperation(left, right, operationTokenType)
	} else {
		return i.calculateLogicalOperation(left, right, operationTokenType)
	}
}

func (i *Interpreter) calculateComparativeOperation(left any, right any, tokenType lex.TokenType) (bool, error) {
	numberLeft, errLeft := castNumber(left)
	numberRight, errRight := castNumber(right)
	if errLeft != nil {
		return false, errLeft
	}
	if errRight != nil {
		return false, errRight
	}

	switch tokenType {
	case lex.LESS:
		return numberLeft < numberRight, nil
	case lex.LESS_EQUAL:
		return numberLeft <= numberRight, nil
	case lex.BANG_EQUAL:
		return numberLeft != numberRight, nil
	case lex.EQUAL_EQUAL:
		return numberLeft == numberRight, nil
	case lex.GREATER:
		return numberLeft > numberRight, nil
	case lex.GREATER_EQUAL:
		return numberLeft >= numberRight, nil
	default:
		return false, errors.New("Unhandled comparative operation: " + string(tokenType))
	}
}

func (i *Interpreter) calculateLogicalOperation(left any, right any, tokenType lex.TokenType) (bool, error) {
	boolLeft, errLeft := castBoolean(left)
	boolRight, errRight := castBoolean(right)
	if errLeft != nil {
		return false, errLeft
	}
	if errRight != nil {
		return false, errRight
	}

	switch tokenType {
	case lex.AND:
		return boolLeft && boolRight, nil
	case lex.OR:
		return boolLeft || boolRight, nil
	default:
		return false, errors.New("Unhandled logical operation: " + string(tokenType))
	}
}

func (i *Interpreter) calculateArithmeticOperation(left any, right any, tokenType lex.TokenType) (float64, error) {
	numberLeft, errLeft := castNumber(left)
	numberRight, errRight := castNumber(right)
	if errLeft != nil {
		return 0, errLeft
	}
	if errRight != nil {
		return 0, errRight
	}

	switch tokenType {
	case lex.MINUS:
		return numberLeft - numberRight, nil
	case lex.PLUS:
		return numberLeft + numberRight, nil
	case lex.STAR:
		return numberLeft * numberRight, nil
	case lex.SLASH:
		return numberLeft / numberRight, nil
	default:
		return 0, errors.New("Unhandled arithmetic operation: " + string(tokenType))
	}
}

func (i *Interpreter) isLiteral(expressions ...syntax.Expr) bool {
	for _, expression := range expressions {
		if expression.Type() != syntax.LITERAL_EXPR {
			return false
		}
	}

	return true
}