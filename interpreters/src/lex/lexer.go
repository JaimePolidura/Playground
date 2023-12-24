package lex

type Lexer struct {
	inputCode string
}

func CreateLexer(inputCode string) *Lexer {
	return &Lexer{inputCode: inputCode}
}

func (lexer *Lexer) HasNext() bool {
}

func (lexer *Lexer) Next() (Token, error) {
}
