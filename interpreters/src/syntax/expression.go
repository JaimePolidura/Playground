package syntax

import "interpreters/src/lex"

type ExpressionType int

const (
	BINARY ExpressionType = iota
	UNARY
	GROUPING
)

type Expression struct {
	Left  *Expression
	Right *Expression
	Token lex.Token

	ExpressionType ExpressionType
}

type BinaryExpression struct {
	Expression
}

type UnaryExpression struct {
	Expression
}

type GroupingExpression struct {
}

func CreateBinaryExpression(left *Expression, right *Expression, token lex.Token) *BinaryExpression {
	return &BinaryExpression{
		Expression: Expression{
			Left:           left,
			Right:          right,
			Token:          token,
			ExpressionType: BINARY,
		},
	}
}

func CreateUnaryExpression(right *Expression, token lex.Token) *BinaryExpression {
	return &BinaryExpression{
		Expression: Expression{
			Right:          right,
			Token:          token,
			ExpressionType: UNARY,
		},
	}
}

func CreateGrouping(right *Expression, token lex.Token) *BinaryExpression {
	return &BinaryExpression{
		Expression: Expression{
			Right:          right,
			Token:          token,
			ExpressionType: UNARY,
		},
	}
}
