package interpreter

import (
	"errors"
	"interpreters/src/lex"
	"interpreters/src/syntax"
)

func interpretExpression(rootExpression syntax.Expr) (syntax.Expr, error) {
	return interpretRecursiveExpression(rootExpression)
}

func interpretRecursiveExpression(rootExpression syntax.Expr) (syntax.Expr, error) {
	switch rootExpression.Type() {
	case syntax.BINARY:
		return interpretBinaryExpression(rootExpression.(syntax.BinaryExpression))
	case syntax.UNARY:
		return interpretUnaryExpression(rootExpression.(syntax.UnaryExpression))
	case syntax.GROUPING:
		return interpretGroupingExpression(rootExpression.(syntax.GroupingExpression))
	case syntax.LITERAL:
		return rootExpression, nil
	}

	return nil, nil
}

func interpretUnaryExpression(unaryExpression syntax.UnaryExpression) (syntax.Expr, error) {
	interpretedExpression, err := interpretRecursiveExpression(unaryExpression.Right)
	if err != nil {
		return nil, err
	}

	if interpretedExpression.Type() != syntax.LITERAL {
		unaryExpression.Right = interpretedExpression
		return unaryExpression, nil
	}

	unaryResult, err := calculateUnaryExpression(interpretedExpression.(syntax.LiteralExpression).Literal, unaryExpression.Token.Type)
	if err != nil {
		return nil, err
	}

	return syntax.CreateLiteralExpression(unaryResult), nil
}

func interpretGroupingExpression(expression syntax.GroupingExpression) (syntax.Expr, error) {
	return interpretRecursiveExpression(expression.OtherExpression)
}

func interpretBinaryExpression(expression syntax.BinaryExpression) (syntax.Expr, error) {
	left, errLeft := interpretRecursiveExpression(expression.Left)
	right, errRight := interpretRecursiveExpression(expression.Right)
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

	literalInterpretedResult, err := calculateBinaryExpression(left.(syntax.LiteralExpression).Literal,
		right.(syntax.LiteralExpression).Literal, expression.Token.Type)
	if err != nil {
		return nil, err
	}

	return syntax.CreateLiteralExpression(literalInterpretedResult), nil
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

func isLiteral(expressions ...syntax.Expr) bool {
	for _, expression := range expressions {
		if expression.Type() != syntax.LITERAL {
			return false
		}
	}

	return true
}
