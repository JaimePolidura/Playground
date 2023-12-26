package syntax

import (
	"errors"
	"interpreters/src/lex"
	"reflect"
)

func Evaluate(rootExpression Expr) (Expr, error) {
	return evaluateRecursive(rootExpression)
}

func evaluateRecursive(rootExpression Expr) (Expr, error) {
	switch rootExpression.Type() {
	case BINARY:
		return evaluateBinary(rootExpression.(BinaryExpression))
	case UNARY:
		return evaluateUnary(rootExpression.(UnaryExpression))
	case GROUPING:
		return evaluateGrouping(rootExpression.(GroupingExpression))
	case LITERAL:
		return rootExpression, nil
	}

	return nil, nil
}

func evaluateUnary(unaryExpression UnaryExpression) (Expr, error) {
	evaluatedExpression, err := evaluateRecursive(unaryExpression.Right)
	if err != nil {
		return nil, err
	}

	if evaluatedExpression.Type() != LITERAL {
		unaryExpression.Right = evaluatedExpression
		return unaryExpression, nil
	}

	unaryResult, err := calculateUnaryExpression(evaluatedExpression.(LiteralExpression).Literal, unaryExpression.Token.Type)
	if err != nil {
		return nil, err
	}

	return CreateLiteralExpression(unaryResult), nil
}

func evaluateGrouping(expression GroupingExpression) (Expr, error) {
	return evaluateRecursive(expression)
}

func evaluateBinary(expression BinaryExpression) (Expr, error) {
	left, errLeft := evaluateRecursive(expression.Left)
	right, errRight := evaluateRecursive(expression.Left)
	if errLeft != nil {
		return nil, errLeft
	}
	if errRight != nil {
		return nil, errRight
	}

	if !isLiteral(left, right) {
		expression.Left = left
		expression.Right = right
		return expression, nil
	}

	literalEvaluationResult, err := calculateBinaryExpression(left.(LiteralExpression).Literal,
		right.(LiteralExpression).Literal, expression.Token.Type)
	if err != nil {
		return nil, err
	}

	return CreateLiteralExpression(literalEvaluationResult), nil
}

func calculateUnaryExpression(literal any, tokenType lex.TokenType) (any, error) {
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

func calculateBinaryExpression(left any, right any, operationTokenType lex.TokenType) (any, error) {
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
		return calculateComparativeOperation(left, right, operationTokenType)
	} else if isArithmeticOperation {
		return calculateArithmeticOperation(left, right, operationTokenType)
	} else {
		return calculateLogicalOperation(left, right, operationTokenType)
	}
}

func calculateComparativeOperation(left any, right any, tokenType lex.TokenType) (bool, error) {
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

func calculateLogicalOperation(left any, right any, tokenType lex.TokenType) (bool, error) {
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

func calculateArithmeticOperation(left any, right any, tokenType lex.TokenType) (float64, error) {
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

func castBoolean(value any) (bool, error) {
	switch value.(type) {
	case float64:
		return value.(float64) > 0, nil
	case bool:
		return value.(bool), nil
	default:
		return false, errors.New("Cannot take " + reflect.TypeOf(value).Name() + " as boolean")
	}
}

func castNumber(value any) (float64, error) {
	switch value.(type) {
	case float64:
		return value.(float64), nil
	case bool:
		return value.(float64), nil
	default:
		return -1, errors.New("Cannot take " + reflect.TypeOf(value).Name() + " as number")
	}
}

func isLiteral(expressions ...Expr) bool {
	for _, expression := range expressions {
		if expression.Type() != LITERAL {
			return false
		}
	}

	return true
}
