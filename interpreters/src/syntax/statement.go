package syntax

type StatementType int

const (
	IF StatementType = iota
	PRINT
	EXPRESSION
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
	return EXPRESSION
}

func (p PrintStatement) Type() StatementType {
	return PRINT
}
