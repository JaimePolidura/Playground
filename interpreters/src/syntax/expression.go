package syntax

import "interpreters/src/lex"

type ExpressionType int

const (
	BINARY_EXPR ExpressionType = iota
	UNARY_EXPR
	GROUPING_EXPR
	LITERAL_EXPR
	VARIABLE_EXPR
	ASSIGN_EXPR
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

type VariableExpression struct {
	Name lex.Token
}

type AssignExpression struct {
	Name  lex.Token
	Value Expr
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
	return BINARY_EXPR
}
func (e GroupingExpression) Type() ExpressionType {
	return GROUPING_EXPR
}
func (e UnaryExpression) Type() ExpressionType {
	return UNARY_EXPR
}
func (e LiteralExpression) Type() ExpressionType {
	return LITERAL_EXPR
}
func (e VariableExpression) Type() ExpressionType {
	return VARIABLE_EXPR
}

func (e AssignExpression) Type() ExpressionType {
	return ASSIGN_EXPR
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

func CreateVariableExpression(name lex.Token) Expr {
	return VariableExpression{
		Name: name,
	}
}

func CreateAssignExpression(name lex.Token, value Expr) Expr {
	return AssignExpression{
		Name:  name,
		Value: value,
	}
}
