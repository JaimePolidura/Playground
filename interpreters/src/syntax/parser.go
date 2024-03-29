package syntax

import (
	"interpreters/src/lex"
)

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
	} else if p.match(lex.FUN) {
		return p.function()
	} else if p.match(lex.CLASS) {
		return p.class()
	} else {
		return p.statement()
	}
}

func (p *Parser) class() (Stmt, error) {
	name := p.consume(lex.IDENTIFIER, "Class name expected")

	superClassName := ""
	if p.check(lex.LESS) { //Extends another class
		p.advance() //Consume <
		superClassName = p.consume(lex.IDENTIFIER, "Expected class name after <").Lexeme
	}

	p.consume(lex.OPEN_BRACE, "Expected '{' after class name")
	methods := make([]FunctionStatement, 0)

	for p.check(lex.FUN) && !p.atTheEnd() {
		p.advance() //Consume fun
		if method, err := p.function(); err == nil {
			methods = append(methods, method.(FunctionStatement))
		} else {
			return nil, err
		}
	}

	p.consume(lex.CLOSE_BRACE, "Expected '}' after class declaration")

	return CreateClassStatement(name, methods, superClassName), nil
}

func (p *Parser) function() (Stmt, error) {
	name := p.consume(lex.IDENTIFIER, "Expect name of the function")
	p.consume(lex.OPEN_PAREN, "'(' expected after function name")
	parameters := make([]lex.Token, 0)

	if !p.match(lex.CLOSE_PAREN) {
		for {
			parameterName := p.consume(lex.IDENTIFIER, "Expect parameter name in function")
			parameters = append(parameters, parameterName)

			if !p.match(lex.COMMA) {
				break
			}
		}
		p.consume(lex.CLOSE_PAREN, "Expect ')' after function parameters")
	}
	p.consume(lex.OPEN_BRACE, "Expect '{' after function declaration")
	body, err := p.blockStatement()
	if err != nil {
		return nil, err
	}

	return CreateFunctionStatement(name, parameters, body.(BlockStatement).Statements), nil
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
	if p.match(lex.WHILE) {
		return p.whileStatement()
	}
	if p.match(lex.OPEN_BRACE) {
		return p.blockStatement()
	}
	if p.match(lex.IF) {
		return p.ifStatement()
	}
	if p.match(lex.FOR) {
		return p.forStatement()
	}
	if p.match(lex.RETURN) {
		return p.returnStatement()
	}
	return p.expressionStatement()
}

func (p *Parser) returnStatement() (Stmt, error) {
	keyword := p.previousToken()
	var value Expr
	if !p.check(lex.SEMICOLON) {
		valueParsed, err := p.parseExpression()
		if err != nil {
			return nil, err
		}

		value = valueParsed
	}
	p.consume(lex.SEMICOLON, "Expect ';' after return statement")

	return CreateReturnStatement(keyword, value), nil
}

func (p *Parser) forStatement() (Stmt, error) {
	p.consume(lex.OPEN_PAREN, "Expect '(' after 'for'.")
	p.consume(lex.VAR, "Expect var in for initializer")
	initializer, err := p.varDeclaration()
	if err != nil {
		return nil, err
	}
	condition := p.expression()
	p.consume(lex.SEMICOLON, "Expect ; after for condition")
	increment := p.expression()
	p.consume(lex.CLOSE_PAREN, "Expect ) after for increment")
	body, err := p.statement()
	if err != nil {
		return nil, err
	}

	whileBody := CreateBlockStatement([]Stmt{body, CreateExpressionStatement(increment)})
	while := CreateWhileStatement(condition, whileBody)

	return CreateBlockStatement([]Stmt{initializer, while}), nil
}

func (p *Parser) whileStatement() (Stmt, error) {
	p.consume(lex.OPEN_PAREN, "Expect '(' at beginning of while loop")
	condition := p.expression()
	p.consume(lex.CLOSE_PAREN, "Expect ')' at the end of while loop")
	body, err := p.statement()
	if err != nil {
		return nil, err
	}

	return CreateWhileStatement(condition, body), nil
}

func (p *Parser) ifStatement() (Stmt, error) {
	p.consume(lex.OPEN_PAREN, "Expect '(' at beginning of if condition")
	ifCondition := p.expression()
	p.consume(lex.CLOSE_PAREN, "Expect ')' at end of if condition")

	thenBranch, err := p.statement()
	if err != nil {
		return nil, err
	}
	var elseBranch Stmt
	if p.match(lex.ELSE) {
		if newElseBranch, err := p.statement(); err != nil {
			return nil, err
		} else {
			elseBranch = newElseBranch
		}
	}

	return CreateIfStatement(ifCondition, thenBranch, elseBranch), nil
}

func (p *Parser) blockStatement() (Stmt, error) {
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
	expr := p.or() //Actual expression

	if p.match(lex.EQUAL) {
		variableValueExpr := p.assignment() //We get the other expression

		if expr.Type() == VARIABLE_EXPR {
			return CreateAssignExpression(expr.(VariableExpression).Name, variableValueExpr)
		} else if expr.Type() == GET_EXPR {
			return CreateSetExpression(expr.(GetExpression).Object, expr.(GetExpression).Name, variableValueExpr)
		} else {
			panic("Invalid assigment")
		}
	}

	return expr
}

func (p *Parser) or() Expr {
	expr := p.and()

	for p.match(lex.OR) {
		operator := p.previousToken()
		right := p.equality()
		expr = CreateLogicalExpression(operator, expr, right)
	}

	return expr
}

func (p *Parser) and() Expr {
	expr := p.equality()

	for p.match(lex.AND) {
		operator := p.previousToken()
		right := p.equality()
		expr = CreateLogicalExpression(operator, expr, right)
	}

	return expr
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
		return p.call()
	}
}

// Page 196
func (p *Parser) call() Expr {
	expr := p.primary()

	for {
		if p.match(lex.OPEN_PAREN) {
			expr = p.finishCall(expr)
		} else if p.match(lex.DOT) { //Accessing object property
			name := p.consume(lex.IDENTIFIER, "Expect property name after .")
			expr = CreateGetExpression(expr, name)
		} else {
			break
		}
	}

	return expr
}

func (p *Parser) finishCall(callee Expr) Expr {
	args := make([]Expr, 0)

	if !p.check(lex.CLOSE_PAREN) {
		args = append(args, p.expression())

		for p.match(lex.COMMA) {
			args = append(args, p.expression())

			if len(args) >= 255 {
				panic("cannot handle more than 255 arguments")
			}
		}
	}

	parent := p.consume(lex.CLOSE_PAREN, "Expect ')' after arguments.")

	return CreateCallExpression(callee, parent, args)
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
