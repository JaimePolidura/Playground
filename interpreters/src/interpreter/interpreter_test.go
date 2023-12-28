package interpreter

import (
	"github.com/stretchr/testify/assert"
	"interpreters/src/lex"
	"interpreters/src/syntax"
	"testing"
)

func TestInterpreter_Interpret_OnlyPrint(t *testing.T) {
	lexer := lex.CreateLexerFromLines("print \"hola\";",
		"print 1 + 2 == 2;",
		"print (7 + 7 + 7);")
	tokens, _ := lexer.ScanTokens()
	parser := syntax.CreateParser(tokens)
	statements, _ := parser.Parse()
	interpreter := CreateInterpreter(statements)

	err := interpreter.Interpret()

	assert.Nil(t, err)
	assert.Equal(t, len(interpreter.Log), 3)
	assert.Equal(t, interpreter.Log[0], "hola")
	assert.Equal(t, interpreter.Log[1], "false")
	assert.Equal(t, interpreter.Log[2], "21")
}

func TestInterpreter_Interpret_Variable(t *testing.T) {
	lexer := lex.CreateLexerFromLines("var numero = 1;",
		"print numero;",
		"numero = 2;",
		"print numero;"+
			"var numero2 = numero + 2;",
		"print numero2;")
	tokens, _ := lexer.ScanTokens()
	parser := syntax.CreateParser(tokens)
	statements, _ := parser.Parse()
	interpreter := CreateInterpreter(statements)

	err := interpreter.Interpret()

	assert.Nil(t, err)
	assert.Equal(t, len(interpreter.Log), 3)
	assert.Equal(t, interpreter.Log[0], "1")
	assert.Equal(t, interpreter.Log[1], "2")
	assert.Equal(t, interpreter.Log[2], "4")
}
