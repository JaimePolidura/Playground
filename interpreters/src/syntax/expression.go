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

// The infix arithmetic (+, -, *, /) and logic (==, !=, <=, < , >, >=)
type BinaryExpression struct {
	Left  Expr
	Right Expr
	Token lex.Token
}

// A prefix ! to perform a logical not
type UnaryExpression struct {
	Right Expr
	Token lex.Token
}

// Parentheses – A pair of ( and ) wrapped around an expression
type GroupingExpression struct {
	OtherExpression Expr
}

// Numbers, strings, Booleans, and nil.
type LiteralExpression struct {
	Literal any
}

func (e BinaryExpression) Type() ExpressionType {
	return BINARY
}
func (e GroupingExpression) Type() ExpressionType {
	return GROUPING
}
func (e UnaryExpression) Type() ExpressionType {
	return UNARY
}
func (e LiteralExpression) Type() ExpressionType {
	return LITERAL
}

func CreateLiteralExpression(literal any) Expr {
	return LiteralExpression{
		Literal: literal,
	}
}

func CreateGroupingExpression(otherExpression Expr) Expr {
	return GroupingExpression{
		OtherExpression: otherExpression,
	}
}

func CreateBinaryExpression(left Expr, right Expr, token lex.Token) Expr {
	return BinaryExpression{
		Left:  left,
		Right: right,
		Token: token,
	}
}

func CreateUnaryExpression(right Expr, token lex.Token) Expr {
	return BinaryExpression{
		Right: right,
		Token: token,
	}
}
