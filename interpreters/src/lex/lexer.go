package lex

import (
	"interpreters/src/utils"
	"strconv"
	"strings"
)

type Lexer struct {
	code string

	current int //Offset
	start   int //Offset
	line    int

	tokens []Token
}

func CreateLexer(inputCode string) *Lexer {
	return &Lexer{code: inputCode}
}

func CreateLexerFromLines(lines ...string) *Lexer {
	return &Lexer{code: strings.Join(lines[:], "\n")}
}

func (l *Lexer) ScanTokens() ([]Token, error) {
	for !l.atTheEnd() {
		l.start = l.current
		if err := l.scanToken(); err != nil {
			return []Token{}, err
		}
	}

	l.tokens = append(l.tokens, CreateEOF(l.line))

	return l.tokens, nil
}

func (l *Lexer) scanToken() error {
	actualChar := l.advance()

	switch actualChar {
	case '(':
		l.addToken(OPEN_PAREN, nil)
		break
	case ')':
		l.addToken(CLOSE_PAREN, nil)
		break
	case '{':
		l.addToken(OPEN_BRACE, nil)
		break
	case '}':
		l.addToken(CLOSE_BRACE, nil)
		break
	case ',':
		l.addToken(COMMA, nil)
		break
	case '.':
		l.addToken(DOT, nil)
		break
	case '-':
		l.addToken(MINUS, nil)
		break
	case '+':
		l.addToken(PLUS, nil)
		break
	case ';':
		l.addToken(SEMICOLON, nil)
		break
	case '*':
		l.addToken(STAR, nil)
		break
	case '/':
		if l.isNext('/') { //Ignore comments
			for l.getActual() != '\n' && !l.atTheEnd() {
				l.advance()
			}
		} else {
			l.addToken(SLASH, nil)
		}

		break
	case ' ':
	case '\r':
	case '\t':
		break
	case '\n':
		l.line++
		break
	case '!':
		l.addToken(l.matchNext('=', BANG_EQUAL, BANG), nil)
		break
	case '=':
		l.addToken(l.matchNext('=', EQUAL_EQUAL, EQUAL), nil)
		break
	case '>':
		l.addToken(l.matchNext('=', GREATER_EQUAL, GREATER), nil)
		break
	case '<':
		l.addToken(l.matchNext('=', LESS_EQUAL, LESS), nil)
		break
	case '"':
		if err := l.addStringToken(); err != nil {
			return err
		}
		break
	default:
		if l.isDigit(actualChar) {
			if err := l.addNumberToken(); err != nil {
				return err
			}
		} else if l.isAlpha(actualChar) {
			l.addIdentifierToken()
		} else {
			return utils.LoxError{Line: l.line, Where: "", Message: "Unexpected token: " + string(actualChar)}
		}
	}

	return nil
}

func (l *Lexer) addIdentifierToken() {
	for l.isAlphaNumeric(l.getActual()) {
		l.advance()
	}

	value := l.code[l.start:l.current]
	keywordTokenType, isKeyword := GetKeyword(value)

	if isKeyword {
		l.addToken(keywordTokenType, nil)
	} else {
		l.addToken(IDENTIFIER, value)
	}
}

func (l *Lexer) isAlphaNumeric(char rune) bool {
	return l.isAlpha(char) || l.isDigit(char)
}

func (l *Lexer) isAlpha(char rune) bool {
	return (char >= 'a' && char <= 'z') ||
		(char >= 'A' && char <= 'Z') ||
		char == '_'
}

func (l *Lexer) addNumberToken() error {
	for l.isDigit(l.getActual()) {
		l.advance()
	}

	if (l.getActual() == '.' || l.getActual() == ',') && l.isDigit(l.getNext()) {
		l.advance()
		for l.isDigit(l.getActual()) {
			l.advance()
		}
	}

	textNumber := l.code[l.start:l.current]
	number, err := strconv.ParseFloat(textNumber, 64)
	if err != nil {
		return err
	}

	l.addToken(NUMBER, float64(number))

	return nil
}

func (l *Lexer) addStringToken() error {
	for l.getActual() != '"' && !l.atTheEnd() {
		if l.getActual() == '\n' {
			l.line++
		}

		l.advance()
	}

	if l.atTheEnd() {
		return utils.LoxError{Line: l.line, Where: "", Message: "Unexpected end of string"}
	}

	l.advance()

	value := l.code[l.start+1 : l.current-1] //Get rid of double quotes
	l.addToken(STRING, value)

	return nil
}

func (l *Lexer) getActual() rune {
	if !l.atTheEnd() {
		return rune(l.code[l.current])
	} else {
		return rune(0)
	}
}

func (l *Lexer) getNext() rune {
	if l.current+1 >= len(l.code) {
		return rune(0)
	} else {
		return rune(l.code[l.current+1])
	}
}

func (l *Lexer) matchNext(expected rune, valueIfTrue TokenType, valueIfFalse TokenType) TokenType {
	if l.isNext(expected) {
		return valueIfTrue
	} else {
		return valueIfFalse
	}
}

func (l *Lexer) isNext(expected rune) bool {
	if l.atTheEnd() {
		return false
	}
	if rune(l.code[l.current]) != expected {
		return false
	}

	l.current = l.current + 1
	return true
}

func (l *Lexer) addToken(tokenType TokenType, literal any) {
	lexeme := l.code[l.start:l.current]

	l.tokens = append(l.tokens, Token{
		Type:    tokenType,
		Line:    l.line,
		Literal: literal,
		Lexeme:  lexeme,
	})
}

func (l *Lexer) advance() rune {
	l.current = l.current + 1
	return rune(l.code[l.current-1])
}

func (l *Lexer) atTheEnd() bool {
	return l.current >= len(l.code)
}

func (l *Lexer) isDigit(char rune) bool {
	return char >= '0' && char <= '9'
}

func IsAnyType(target TokenType, types ...TokenType) bool {
	for _, tokenType := range types {
		if tokenType == target {
			return true
		}
	}

	return false
}
