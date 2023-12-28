package interpreter

import (
	"github.com/stretchr/testify/assert"
	"interpreters/src/lex"
	"interpreters/src/syntax"
	"testing"
)

func TestInterpreter_NoError_Error(t *testing.T) {
	_, err := interpretTest("2 + \"hola\"") //2 * 0
	assert.NotNil(t, err)
}

func TestInterpreter_NoError_FinalValueNotLiteral(t *testing.T) {
	//evaluate("(1 + 34) + PI")

	//TODO Pending. Expect BINARY(LITERAL(35), IDENTIFIER?)
}

func TestInterpreter_NoError_FinalValueLiteral(t *testing.T) {
	result, err := interpretTest("2 * (1 + 2) == 3")

	assert.Nil(t, err)
	assert.Equal(t, result.Type(), syntax.LITERAL)
	assert.Equal(t, result.(syntax.LiteralExpression).Literal, false)

	result, err = interpretTest("(2 + 1) * (2 + 1)")

	assert.Nil(t, err)
	assert.Equal(t, result.Type(), syntax.LITERAL)
	assert.Equal(t, result.(syntax.LiteralExpression).Literal, float64(9))

	result, err = interpretTest("1 != 2")

	assert.Nil(t, err)
	assert.Equal(t, result.Type(), syntax.LITERAL)
	assert.Equal(t, result.(syntax.LiteralExpression).Literal, true)

	result, err = interpretTest("2 * false") //2 * 0
	assert.Nil(t, err)
	assert.Equal(t, result.Type(), syntax.LITERAL)
	assert.Equal(t, result.(syntax.LiteralExpression).Literal, float64(0))
}

func interpretTest(code string) (syntax.Expr, error) {
	lexer := lex.CreateLexer(code)
	tokens, _ := lexer.ScanTokens()
	parser := syntax.CreateParser(tokens)
	expr, _ := parser.ParseExpression()

	return interpretExpression(expr)
}
