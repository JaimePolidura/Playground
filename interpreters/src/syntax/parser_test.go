package syntax

import (
	"github.com/stretchr/testify/assert"
	"interpreters/src/lex"
	"testing"
)

/*
				!=
		       /  \
		      -   false
			 / \
			+   /
	       /\  /\
		  1 2 1 2
*/
func TestParser_Parse(t *testing.T) {
	lexer := lex.CreateLexer("(1 + 2 - (1 / 2)) != false")
	tokens, _ := lexer.ScanTokens()

	parser := CreateParser(tokens)
	expr, _ := parser.Parse()

	assert.Equal(t, expr.Type(), BINARY)

	assert.Equal(t, expr.(BinaryExpression).Right.Type(), LITERAL)
	assert.Equal(t, expr.(BinaryExpression).Right.(LiteralExpression).Literal, false)

	assert.Equal(t, expr.(BinaryExpression).Left.Type(), GROUPING)
	minus := expr.(BinaryExpression).Left.(GroupingExpression).OtherExpression.(BinaryExpression)
	assert.Equal(t, minus.Type(), BINARY)

	assert.Equal(t, minus.Left.Type(), BINARY)
	assert.Equal(t, minus.Right.Type(), GROUPING)
}
