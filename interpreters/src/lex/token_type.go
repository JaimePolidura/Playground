package lex

type TokenType string

const (
	OPEN_PAREN  TokenType = "OPEN_PAREN"
	CLOSE_PAREN           = "CLOSE_PAREN"
	OPEN_BRACE            = "OPEN_BRACE"
	CLOSE_BRACE           = "CLOSE_BRACE"
	COMMA                 = "COMMA"
	DOT                   = "DOT"
	MINUS                 = "MINUS"
	PLUS                  = "PLUS"
	SEMICOLON             = "SEMICOLON"
	SLASH                 = "SLASH"
	STAR                  = "STAR"

	// One or two character tokens.
	BANG          = "BANG"
	BANG_EQUAL    = "BANG_EQUAL"
	EQUAL         = "EQUAL"
	EQUAL_EQUAL   = "EQUAL_EQUAL"
	GREATER       = "GREATER"
	GREATER_EQUAL = "GREATER_EQUAL"
	LESS          = "LESS"
	LESS_EQUAL    = "LESS_EQUAL"

	// Literals.
	IDENTIFIER = "IDENTIFIER"
	STRING     = "STRING"
	NUMBER     = "NUMBER"

	// Keywords.
	AND    = "AND"
	CLASS  = "CLASS"
	ELSE   = "ELSE"
	FALSE  = "FALSE"
	FUN    = "FUN"
	FOR    = "FOR"
	IF     = "IF"
	NIL    = "NIL"
	OR     = "OR"
	PRINT  = "PRINT"
	RETURN = "RETURN"
	SUPER  = "SUPER"
	THIS   = "THIS"
	TRUE   = "TRUE"
	VAR    = "VAR"
	WHILE  = "WHILE"

	EOF = "EOF"
)
