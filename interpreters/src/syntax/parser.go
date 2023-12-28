package syntax

import "interpreters/src/lex"

type Parser struct {
	Tokens []lex.Token

	current int
}

func CreateParser(Tokens []lex.Token) *Parser {
	return &Parser{
		Tokens:  Tokens,
		current: 0,
	}
}

func (p *Parser) Parse() ([]Stmt, error) {
	statements := make([]Stmt, 0)
	for !p.atTheEnd() {
		if statement, err := p.declaration(); err != nil {
			return statements, err
		} else {
			statements = append(statements, statement)
		}
	}

	return statements, nil
}

func (p *Parser) parseExpression() (Expr, error) {
	return p.expression(), nil
}

func (p *Parser) declaration() (Stmt, error) {
	if p.match(lex.VAR) {
		return p.varDeclaration()
	} else {
		return p.statement()
	}
}

func (p *Parser) varDeclaration() (Stmt, error) {
	name := p.consume(lex.IDENTIFIER, "Expect variable name")

	var initializer Expr
	if p.match(lex.EQUAL) {
		initializer = p.expression()
	}

	p.consume(lex.SEMICOLON, "Expect ;")

	return CreateVarStatement(name, initializer), nil
}

func (p *Parser) statement() (Stmt, error) {
	if p.match(lex.PRINT) {
		return p.printStatement()
	}
	if p.match(lex.OPEN_BRACE) {
		return p.block()
	}

	return p.expressionStatement()
}

func (p *Parser) block() (Stmt, error) {
	statements := make([]Stmt, 0)
	for !p.check(lex.CLOSE_BRACE) && !p.atTheEnd() {
		if declaration, err := p.declaration(); err != nil {
			return nil, err
		} else {
			statements = append(statements, declaration)
		}
	}

	p.consume(lex.CLOSE_BRACE, "Expected } at the end of the statement")

	return CreateBlockStatement(statements), nil
}

func (p *Parser) printStatement() (Stmt, error) {
	expr := p.expression()
	p.consume(lex.SEMICOLON, "Expect ; after value.")
	return CreatePrintStatement(expr), nil
}

func (p *Parser) expressionStatement() (Stmt, error) {
	expr := p.expression()
	p.consume(lex.SEMICOLON, "Expect ; after value.")
	return CreateExpressionStatement(expr), nil
}

// Every expression starts with assigment
func (p *Parser) expression() Expr {
	return p.assignment()
}

func (p *Parser) assignment() Expr {
	exp := p.equality() //Actual expression

	if p.match(lex.EQUAL) {
		if exp.Type() != VARIABLE_EXPR {
			panic("Invalid assigment")
		}
		variableName := exp.(VariableExpression).Name
		variableValueExpr := p.assignment() //We get the other expression

		return CreateAssignExpression(variableName, variableValueExpr)
	}

	return exp
}

// equality → comparison ( ( "!=" | "==" ) comparison )*
func (p *Parser) equality() Expr {
	expr := p.comparison()

	for p.match(lex.BANG_EQUAL, lex.EQUAL_EQUAL) {
		operator := p.previousToken()
		right := p.comparison()
		expr = CreateBinaryExpression(expr, right, operator)
	}

	return expr
}

// comparison → term ( ( ">" | ">=" | "<" | "<=" ) term )*
func (p *Parser) comparison() Expr {
	expr := p.term()

	for p.match(lex.GREATER, lex.GREATER_EQUAL, lex.LESS, lex.LESS_EQUAL) {
		operator := p.previousToken()
		right := p.factor()
		expr = CreateBinaryExpression(expr, right, operator)
	}

	return expr
}

// term → factor ( ( "-" | "+" ) factor )*
func (p *Parser) term() Expr {
	expr := p.factor()

	for p.match(lex.MINUS, lex.PLUS) {
		operator := p.previousToken()
		right := p.factor()
		expr = CreateBinaryExpression(expr, right, operator)
	}

	return expr
}

// factor → unary ( ( "/" | "*" ) unary )*
func (p *Parser) factor() Expr {
	expr := p.unary()

	for p.match(lex.SLASH, lex.STAR) {
		operator := p.previousToken()
		right := p.unary()
		expr = CreateBinaryExpression(expr, right, operator)
	}

	return expr
}

// unary → ( "!" | "-" ) unary | primary
func (p *Parser) unary() Expr {
	if p.match(lex.MINUS, lex.BANG) {
		prevToken := p.previousToken()
		right := p.unary()
		return CreateUnaryExpression(right, prevToken)
	} else {
		return p.primary()
	}
}

// primary → NUMBER | STRING | "true" | "false" | "nil" | "(" expression ")"
func (p *Parser) primary() Expr {
	if p.match(lex.FALSE) {
		return CreateLiteralExpression(false)
	}
	if p.match(lex.TRUE) {
		return CreateLiteralExpression(true)
	}
	if p.match(lex.NIL) {
		return CreateLiteralExpression(nil)
	}
	if p.match(lex.NUMBER, lex.STRING) {
		return CreateLiteralExpression(p.previousToken().Literal)
	}
	if p.match(lex.IDENTIFIER) {
		return CreateVariableExpression(p.previousToken())
	}
	if p.match(lex.OPEN_PAREN) {
		expr := p.expression()
		p.consume(lex.CLOSE_PAREN, "Expect ')' after expression")
		return CreateGroupingExpression(expr)
	}

	panic("Unexpected token")
}

func (p *Parser) consume(expected lex.TokenType, errorMsg string) lex.Token {
	if p.check(expected) {
		return p.advance()
	}

	panic(errorMsg)
}

func (p *Parser) match(tokens ...lex.TokenType) bool {
	for _, tokenType := range tokens {
		if p.check(tokenType) {
			p.advance()
			return true
		}
	}

	return false
}

func (p *Parser) advance() lex.Token {
	if !p.atTheEnd() {
		p.current++
	}

	return p.previousToken()
}

func (p *Parser) check(tokenTypeToCheck lex.TokenType) bool {
	if p.atTheEnd() {
		return false
	}

	return p.peek().Type == tokenTypeToCheck
}

func (p *Parser) atTheEnd() bool {
	return p.peek().Type == lex.EOF
}

func (p *Parser) peek() lex.Token {
	return p.Tokens[p.current]
}

func (p *Parser) previousToken() lex.Token {
	return p.Tokens[p.current-1]
}
