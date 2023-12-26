package syntax

import "interpreters/src/lex"

type ExpressionType int

const (
	BINARY ExpressionType = iota
	UNARY
	GROUPING
	LITERAL
)

type Expr interface {
	Type() ExpressionType
}

type Expression struct {
	Expr
	ExpressionType ExpressionType
}

func (e Expression) Type() ExpressionType {
	return e.ExpressionType
}

// The infix arithmetic (+, -, *, /) and logic (==, !=, <=, < , >, >=)
type BinaryExpression struct {
	Expression

	Left  Expr
	Right Expr
	Token lex.Token
}

// A prefix ! to perform a logical not
type UnaryExpression struct {
	Expression

	Right Expr
	Token lex.Token
}

// Parentheses – A pair of ( and ) wrapped around an expression
type GroupingExpression struct {
	Expression

	OtherExpression Expr
}

// Numbers, strings, Booleans, and nil.
type LiteralExpression struct {
	Expression

	Literal any
}

func CreateLiteralExpression(literal any) Expr {
	return LiteralExpression{
		Expression: Expression{ExpressionType: LITERAL},
		Literal:    literal,
	}
}

func CreateGroupingExpression(otherExpression Expr) Expr {
	return GroupingExpression{
		Expression:      Expression{ExpressionType: GROUPING},
		OtherExpression: otherExpression,
	}
}

func CreateBinaryExpression(left Expr, right Expr, token lex.Token) Expr {
	return BinaryExpression{
		Expression: Expression{ExpressionType: BINARY},
		Left:       left,
		Right:      right,
		Token:      token,
	}
}

func CreateUnaryExpression(right Expr, token lex.Token) Expr {
	return BinaryExpression{
		Expression: Expression{ExpressionType: UNARY},
		Right:      right,
		Token:      token,
	}
}
