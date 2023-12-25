package syntax

import "interpreters/src/lex"

type Parser struct {
	Tokens []lex.Token

	current int
}

// equality → comparison ( ( "!=" | "==" ) comparison )*
func (p *Parser) equality() Expression {
	expr := p.comparation()

	for p.match(lex.BANG_EQUAL, lex.EQUAL_EQUAL) {
		operator := p.previous()
		right := p.comparation()

		expr = 
	}

	return expr
}

// comparison → term ( ( ">" | ">=" | "<" | "<=" ) term )*
func (p *Parser) comparasion() Expression {
}

func (p *Parser) match(tokens ...lex.TokenType) bool {
	for _, tokenType := range tokens {
		if p.check(tokenType) {
			return true
		}
	}

	return false
}

func (p *Parser) advance() lex.Token {
	if !p.atTheEnd() {
		p.current++
	}

	return p.previous()
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

func (p *Parser) previous() lex.Token {
	return p.Tokens[p.current-1]
}
