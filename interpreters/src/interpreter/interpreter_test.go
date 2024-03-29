package interpreter

import (
	"github.com/stretchr/testify/assert"
	"interpreters/src/lex"
	"interpreters/src/syntax"
	"strconv"
	"testing"
)

func TestInterpreter_Interpret_Class_Superclass(t *testing.T) {
	interpreter, err := interpret("class Base {",
		"fun hello() {",
		"print \"Base\";",
		"}",
		"}",
		"class Derived < Base {",
		"}",
		"var derived = Derived();",
		"derived.hello();")

	assert.Nil(t, err)
	assert.Equal(t, len(interpreter.Log), 1)
	assert.Equal(t, interpreter.Log[0], "Base")
}

func TestInterpreter_Interpret_Class_Constructor(t *testing.T) {
	interpreter, err := interpret("class Point {",
		"fun init(xCons, yCons) {",
		"this.x = xCons;",
		"this.y = yCons;",
		"}",
		"}",
		"var punto = Point(10, 10);",
		"print punto.x;",
		"print punto.y;")

	assert.Nil(t, err)
	assert.Equal(t, len(interpreter.Log), 2)
	assert.Equal(t, interpreter.Log[0], "10")
	assert.Equal(t, interpreter.Log[1], "10")
}

func TestInterpreter_Interpret_Class_Methods_WithThis(t *testing.T) {
	interpreter, err := interpret("class Person {",
		"fun printName() {",
		"print this.name;",
		"}",
		"}",
		"var person = Person();",
		"person.name = \"pedro\";",
		"person.printName();")

	assert.Nil(t, err)
	assert.Equal(t, len(interpreter.Log), 1)
	assert.Equal(t, interpreter.Log[0], "pedro")
}

func TestInterpreter_Interpret_Class_Methods(t *testing.T) {
	interpreter, err := interpret(
		"class Car {",
		"fun start() {",
		"print 1;",
		"}",
		"fun get() {",
		"return 10;",
		"}",
		"}",
		"var car = Car();",
		"car.start();",
		"print car.get();",
	)

	assert.Nil(t, err)
	assert.Equal(t, len(interpreter.Log), 2)
	assert.Equal(t, interpreter.Log[0], "1")
	assert.Equal(t, interpreter.Log[1], "10")
}

func TestInterpreter_Interpret_Class_Properties(t *testing.T) {
	interpreter, err := interpret(
		"class Car {",
		"}",
		"class Shock{",
		"}",
		"var car = Car();",
		"car.numero = \"Motor\";",
		"car.shock = Shock();",
		"car.shock.price = 10;",
		"print car.numero;",
		"car.numero = 12;",
		"print car.numero;",
		"print car.shock.price;")

	assert.Nil(t, err)
	assert.Equal(t, len(interpreter.Log), 3)
	assert.Equal(t, interpreter.Log[0], "Motor")
	assert.Equal(t, interpreter.Log[1], "12")
	assert.Equal(t, interpreter.Log[2], "10")
}

func TestInterpreter_Interpret_RecursiveReturnFunction(t *testing.T) {
	interpreter, err := interpret(
		"fun factRecursive(next, act) {",
		"var newNext = next - 1;",
		"if(newNext == 0) {",
		"return act;",
		"}",
		"return factRecursive(newNext, act * newNext);",
		"}",
		"fun fact(value) {",
		"return factRecursive(value, value);",
		"}",
		"print fact(5);")
	assert.Nil(t, err)
	assert.Equal(t, len(interpreter.Log), 1)
	assert.Equal(t, interpreter.Log[0], "120")
}

func TestInterpreter_Interpret_ReturnFunction(t *testing.T) {
	interpreter, err := interpret("fun max(a, b) {",
		"if(a > b) {",
		"return a;",
		"} else {",
		"return b;",
		"}",
		"}",
		"print max(2, 4);",
		"print max(3, 7);")

	assert.Nil(t, err)
	assert.Equal(t, len(interpreter.Log), 2)
	assert.Equal(t, interpreter.Log[0], "4")
	assert.Equal(t, interpreter.Log[1], "7")
}

func TestInterpreter_Interpret_VoidFunction(t *testing.T) {
	interpreter, err := interpret("fun sumar(a, b) {",
		"var resultado = a + b;",
		"print resultado;",
		"}",
		"sumar(2, 4);",
		"sumar(2, 5);")

	assert.Nil(t, err)
	assert.Equal(t, len(interpreter.Log), 2)
	assert.Equal(t, interpreter.Log[0], "6")
	assert.Equal(t, interpreter.Log[1], "7")
}

func TestInterpreter_Interpret_RecursiveFunction(t *testing.T) {
	interpreter, err := interpret("fun repeat(times, actual) {",
		"print actual;",
		"actual = actual + 1;",
		"if(actual < times) {",
		"repeat(times, actual);",
		"}",
		"}",
		"repeat(3, 0);")

	assert.Nil(t, err)
	assert.Equal(t, len(interpreter.Log), 3)
	assert.Equal(t, interpreter.Log[0], "0")
	assert.Equal(t, interpreter.Log[1], "1")
	assert.Equal(t, interpreter.Log[2], "2")
}

func TestInterpreter_Interpret_ForLoop(t *testing.T) {
	interpreter, err := interpret("for(var i = 0; i < 10; i = i + 1){",
		"print i;",
		"}")

	assert.Nil(t, err)
	assert.Equal(t, len(interpreter.Log), 10)
	for i := 0; i < 10; i++ {
		assert.Equal(t, interpreter.Log[i], strconv.Itoa(i))
	}
}

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
