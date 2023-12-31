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
	LOGICAL_EXPR
	CALL_EXPR
	GET_EXPR
	SET_EXPR
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

type LogicalExpression struct {
	Operator lex.Token
	Left     Expr
	Right    Expr
}

type CallExpression struct {
	Callee Expr
	Parent lex.Token
	Args   []Expr
}

type GetExpression struct {
	Object Expr
	Name   lex.Token
}

type SetExpression struct {
	Object Expr
	Name   lex.Token
	Value  Expr
}

func CreateSetExpression(object Expr, name lex.Token, value Expr) SetExpression {
	return SetExpression{
		Object: object,
		Name:   name,
		Value:  value,
	}
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

func CreateLogicalExpression(operator lex.Token, left Expr, right Expr) Expr {
	return LogicalExpression{
		Operator: operator,
		Left:     left,
		Right:    right,
	}
}

func CreateAssignExpression(name lex.Token, value Expr) Expr {
	return AssignExpression{
		Name:  name,
		Value: value,
	}
}

func CreateCallExpression(callee Expr, parent lex.Token, args []Expr) Expr {
	return CallExpression{
		Args:   args,
		Parent: parent,
		Callee: callee,
	}
}

func CreateGetExpression(object Expr, Name lex.Token) GetExpression {
	return GetExpression{
		Object: object,
		Name:   Name,
	}
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
func (e LogicalExpression) Type() ExpressionType {
	return LOGICAL_EXPR
}
func (e AssignExpression) Type() ExpressionType {
	return ASSIGN_EXPR
}
func (e CallExpression) Type() ExpressionType {
	return CALL_EXPR
}
func (e GetExpression) Type() ExpressionType {
	return GET_EXPR
}
func (e SetExpression) Type() ExpressionType {
	return SET_EXPR
}
