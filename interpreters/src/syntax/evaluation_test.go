package syntax

import (
	"github.com/stretchr/testify/assert"
	"interpreters/src/lex"
	"testing"
)

func TestEvaluate_NoError_Error(t *testing.T) {
	_, err := evaluate("2 + \"hola\"") //2 * 0
	assert.NotNil(t, err)
}

func TestEvaluate_NoError_FinalValueNotLiteral(t *testing.T) {
	//evaluate("(1 + 34) + PI")

	//TODO Pending. Expect BINARY(LITERAL(35), IDENTIFIER?)
}

func TestEvaluate_NoError_FinalValueLiteral(t *testing.T) {
	result, err := evaluate("2 * (1 + 2) == 3")

	assert.Nil(t, err)
	assert.Equal(t, result.Type(), LITERAL)
	assert.Equal(t, result.(LiteralExpression).Literal, false)

	result, err = evaluate("(2 + 1) * (2 + 1)")

	assert.Nil(t, err)
	assert.Equal(t, result.Type(), LITERAL)
	assert.Equal(t, result.(LiteralExpression).Literal, float64(9))

	result, err = evaluate("1 != 2")

	assert.Nil(t, err)
	assert.Equal(t, result.Type(), LITERAL)
	assert.Equal(t, result.(LiteralExpression).Literal, true)

	result, err = evaluate("2 * false") //2 * 0
	assert.Nil(t, err)
	assert.Equal(t, result.Type(), LITERAL)
	assert.Equal(t, result.(LiteralExpression).Literal, float64(0))
}

func evaluate(code string) (Expr, error) {
	lexer := lex.CreateLexer(code)
	tokens, _ := lexer.ScanTokens()
	parser := CreateParser(tokens)
	expr, _ := parser.Parse()

	return Evaluate(expr)
}
