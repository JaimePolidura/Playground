package syntax

import (
	"github.com/stretchr/testify/assert"
	"interpreters/src/lex"
	"testing"
)

func TestParser_Parse2(t *testing.T) {
	lexer := lex.CreateLexer("print \"hola\";\n" +
		"print 1 + 2 == 2;\n" +
		"print (7 * 3) == (7 + 7 + 7);")
	tokens, _ := lexer.ScanTokens()
	parser := CreateParser(tokens)

	statements, err := parser.Parse()

	assert.Nil(t, err)
	assert.Equal(t, len(statements), 3)

	assert.True(t, statements[0].Type() == PRINT && statements[1].Type() == PRINT && statements[2].Type() == PRINT)
	assert.Equal(t, statements[0].(PrintStatement).Expression.Type(), LITERAL)
	assert.Equal(t, statements[1].(PrintStatement).Expression.Type(), BINARY)
	assert.Equal(t, statements[2].(PrintStatement).Expression.Type(), BINARY)
}

/*
				!=
		       /  \
		      -   false
			 / \
			+   /
	       /\  / \
		  1 2 1  2
*/
func TestParser_Parse(t *testing.T) {
	lexer := lex.CreateLexer("(1 + 2 - (1 / 2)) != false")
	tokens, _ := lexer.ScanTokens()

	parser := CreateParser(tokens)
	expr, _ := parser.ParseExpression()

	assert.Equal(t, expr.Type(), BINARY)

	assert.Equal(t, expr.(BinaryExpression).Right.Type(), LITERAL)
	assert.Equal(t, expr.(BinaryExpression).Right.(LiteralExpression).Literal, false)

	assert.Equal(t, expr.(BinaryExpression).Left.Type(), GROUPING)
	minus := expr.(BinaryExpression).Left.(GroupingExpression).OtherExpression.(BinaryExpression)
	assert.Equal(t, minus.Type(), BINARY)

	assert.Equal(t, minus.Left.Type(), BINARY)
	assert.Equal(t, minus.Right.Type(), GROUPING)

	minus_slash := minus.Right.(GroupingExpression).OtherExpression.(BinaryExpression)
	assert.Equal(t, minus_slash.Type(), BINARY)
	assert.Equal(t, minus_slash.Token.Type, lex.SLASH) //NOT
	assert.Equal(t, minus_slash.Left.Type(), LITERAL)
	assert.Equal(t, minus_slash.Right.Type(), LITERAL)
	assert.Equal(t, minus_slash.Left.(LiteralExpression).Literal, float64(1))
	assert.Equal(t, minus_slash.Right.(LiteralExpression).Literal, float64(2))

	minus_plus := minus.Left.(BinaryExpression)
	assert.Equal(t, minus_plus.Type(), BINARY)
	assert.Equal(t, minus_plus.Token.Type, lex.PLUS)
	assert.Equal(t, minus_plus.Left.Type(), LITERAL)
	assert.Equal(t, minus_plus.Right.Type(), LITERAL)
	assert.Equal(t, minus_plus.Left.(LiteralExpression).Literal, float64(1))
	assert.Equal(t, minus_plus.Right.(LiteralExpression).Literal, float64(2))
}
