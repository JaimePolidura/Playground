package lex

type Token struct {
	Type    TokenType
	Lexeme  string
	Literal any
	Line    int
}

func CreateEOF(line int) Token {
	return Token{
		Type:    EOF,
		Lexeme:  "",
		Literal: nil,
		Line:    line,
	}
}
