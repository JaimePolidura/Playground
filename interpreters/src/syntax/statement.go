package syntax

import "interpreters/src/lex"

type StatementType int

const (
	PRINT_STMT StatementType = iota
	EXPRESSION_STMT
	VAR_STMT
)

type Stmt interface {
	Type() StatementType
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

func (p ExpressionStatement) Type() StatementType {
	return EXPRESSION_STMT
}
func (p PrintStatement) Type() StatementType {
	return PRINT_STMT
}
func (p VarStatement) Type() StatementType {
	return VAR_STMT
}
