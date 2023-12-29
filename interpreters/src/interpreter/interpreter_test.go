package interpreter

import (
	"github.com/stretchr/testify/assert"
	"interpreters/src/lex"
	"interpreters/src/syntax"
	"testing"
)

func TestInterpreter_Interpret_IfCondition(t *testing.T) {
	interpreter, err := interprete(
		"var numero1 = 2;",
		"var numero2 = 2;",
		"if(numero1 == numero2) {",
		"print \"1.if\";",
		"if(numero1 + 1 == numero2) {",
		"print \"2.if\";",
		"} else {",
		"print \"2.else\";",
		"}",
		"} else {",
		"print \"1.else\";",
		"}",
		"print \"exit\";",
	)

	assert.Nil(t, err)
	assert.Equal(t, len(interpreter.Log), 3)
	assert.Equal(t, interpreter.Log[0], "1.if")
	assert.Equal(t, interpreter.Log[1], "2.else")
	assert.Equal(t, interpreter.Log[2], "exit")
}

func TestInterpreter_Interpret_OnlyPrint(t *testing.T) {
	interpreter, err := interprete("print \"hola\";",
		"print 1 + 2 == 2;",
		"print (7 + 7 + 7);")

	assert.Nil(t, err)
	assert.Equal(t, len(interpreter.Log), 3)
	assert.Equal(t, interpreter.Log[0], "hola")
	assert.Equal(t, interpreter.Log[1], "false")
	assert.Equal(t, interpreter.Log[2], "21")
}

func TestInterpreter_Interpret_Variable(t *testing.T) {
	interpreter, err := interprete("var numero = 1;",
		"print numero;",
		"numero = 2;",
		"print numero;",
		"var numero2 = numero + 2;",
		"print numero2;")

	assert.Nil(t, err)
	assert.Equal(t, len(interpreter.Log), 3)
	assert.Equal(t, interpreter.Log[0], "1")
	assert.Equal(t, interpreter.Log[1], "2")
	assert.Equal(t, interpreter.Log[2], "4")
}

func TestInterpreter_Interpret_ScopeVariables(t *testing.T) {
	interpreter, err := interprete("var a = 1;",
		"{",
		"a = 10;",
		"var b = 2;",
		"print a;",
		"print b;",
		"}",
		"print a;")

	assert.Nil(t, err)
	assert.Equal(t, len(interpreter.Log), 3)
	assert.Equal(t, interpreter.Log[0], "10")
	assert.Equal(t, interpreter.Log[1], "2")
	assert.Equal(t, interpreter.Log[2], "10")
}

func interprete(code ...string) (*Interpreter, error) {
	lexer := lex.CreateLexerFromLines(code...)
	tokens, _ := lexer.ScanTokens()
	parser := syntax.CreateParser(tokens)
	statements, _ := parser.Parse()
	interpreter := CreateInterpreter(statements)

	err := interpreter.Interpret()

	return interpreter, err
}
