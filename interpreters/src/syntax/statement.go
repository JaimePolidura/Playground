package syntax

import "interpreters/src/lex"

type StatementType int

const (
	PRINT_STMT StatementType = iota
	EXPRESSION_STMT
	VAR_STMT
	BLOCK_STMT
	IF_STMT
	WHILE_STMT
	FUNCTION_STMT
	RETURN_STMT
	CLASS_STMT
)

type Stmt interface {
	Type() StatementType
}

type WhileStatement struct {
	Condition Expr
	Body      Stmt
}

type PrintStatement struct {
	Expression Expr
}

type ExpressionStatement struct {
	Expression Expr
}

type VarStatement struct {
	Name        lex.Token
	Initializer Expr
}

type BlockStatement struct {
	Statements []Stmt
}

type FunctionStatement struct {
	Name   lex.Token
	Params []lex.Token
	Body   []Stmt
}

type IfStatement struct {
	Condition  Expr
	ThenBranch Stmt
	ElseBranch Stmt
}

type ReturnStatement struct {
	Keyword lex.Token
	Value   Expr
}

type ClassStatement struct {
	Name    lex.Token
	Methods []Stmt
}

func CreateIfStatement(condition Expr, thenBranch Stmt, elseBranch Stmt) IfStatement {
	return IfStatement{
		Condition:  condition,
		ThenBranch: thenBranch,
		ElseBranch: elseBranch,
	}
}

func CreateVarStatement(name lex.Token, initializer Expr) VarStatement {
	return VarStatement{
		Name:        name,
		Initializer: initializer,
	}
}

func CreateExpressionStatement(other Expr) ExpressionStatement {
	return ExpressionStatement{
		Expression: other,
	}
}

func CreatePrintStatement(other Expr) PrintStatement {
	return PrintStatement{
		Expression: other,
	}
}

func CreateBlockStatement(statements []Stmt) BlockStatement {
	return BlockStatement{
		Statements: statements,
	}
}

func CreateWhileStatement(condition Expr, body Stmt) WhileStatement {
	return WhileStatement{
		Condition: condition,
		Body:      body,
	}
}

func CreateReturnStatement(keyword lex.Token, value Expr) ReturnStatement {
	return ReturnStatement{
		Keyword: keyword,
		Value:   value,
	}
}

func CreateClassStatement(name lex.Token, methods []Stmt) ClassStatement {
	return ClassStatement{
		Name:    name,
		Methods: methods,
	}
}

func CreateFunctionStatement(name lex.Token, params []lex.Token, body []Stmt) FunctionStatement {
	return FunctionStatement{
		Name:   name,
		Params: params,
		Body:   body,
	}
}

func (p ExpressionStatement) Type() StatementType {
	return EXPRESSION_STMT
}
func (p PrintStatement) Type() StatementType {
	return PRINT_STMT
}
func (p VarStatement) Type() StatementType {
	return VAR_STMT
}
func (p BlockStatement) Type() StatementType {
	return BLOCK_STMT
}
func (p WhileStatement) Type() StatementType {
	return WHILE_STMT
}
func (p IfStatement) Type() StatementType {
	return IF_STMT
}
func (p FunctionStatement) Type() StatementType {
	return FUNCTION_STMT
}
func (p ReturnStatement) Type() StatementType {
	return RETURN_STMT
}
func (p ClassStatement) Type() StatementType {
	return CLASS_STMT
}
