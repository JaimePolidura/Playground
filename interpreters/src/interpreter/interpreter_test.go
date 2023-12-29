package interpreter

import (
	"github.com/stretchr/testify/assert"
	"interpreters/src/lex"
	"interpreters/src/syntax"
	"strconv"
	"testing"
)

func TestInterpreter_Interpret_WhileLoop(t *testing.T) {
	interpreter, err := interpret(
		"var numero1 = 0;",
		"while (numero1 < 10) {",
		"print numero1;",
		"numero1 = numero1 + 1;",
		"}",
		"while(numero1 == 0){",
		"print numero1;",
		"}",
	)

	assert.Nil(t, err)
	assert.Equal(t, len(interpreter.Log), 10)
	for i := 0; i < 10; i++ {
		assert.Equal(t, interpreter.Log[i], strconv.Itoa(i))
	}
}

func TestInterpreter_Interpret_IfWithLogicalOperators(t *testing.T) {
	interpreter, err := interpret(
		"var numero1 = 2;",
		"var numero2 = 2;",
		"if(((numero1 + 1) == numero2) or (numero1 == numero2)) {",
		"print \"1.if\";",
		"}",
		"if((numero1 == numero2) and numero2 > numero1) {",
		"print \"2.if\";",
		"}",
	)

	assert.Nil(t, err)
	assert.Equal(t, len(interpreter.Log), 1)
	assert.Equal(t, interpreter.Log[0], "1.if")
}

func TestInterpreter_Interpret_IfCondition(t *testing.T) {
	interpreter, err := interpret(
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
	interpreter, err := interpret("print \"hola\";",
		"print 1 + 2 == 2;",
		"print (7 + 7 + 7);")

	assert.Nil(t, err)
	assert.Equal(t, len(interpreter.Log), 3)
	assert.Equal(t, interpreter.Log[0], "hola")
	assert.Equal(t, interpreter.Log[1], "false")
	assert.Equal(t, interpreter.Log[2], "21")
}

func TestInterpreter_Interpret_Variable(t *testing.T) {
	interpreter, err := interpret("var numero = 1;",
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
	interpreter, err := interpret("var a = 1;",
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

func interpret(code ...string) (*Interpreter, error) {
	lexer := lex.CreateLexerFromLines(code...)
	tokens, _ := lexer.ScanTokens()
	parser := syntax.CreateParser(tokens)
	statements, _ := parser.Parse()
	interpreter := CreateInterpreter(statements)

	err := interpreter.Interpret()

	return interpreter, err
}
