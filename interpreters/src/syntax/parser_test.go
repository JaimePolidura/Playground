package syntax

import (
	"github.com/stretchr/testify/assert"
	"interpreters/src/lex"
	"testing"
)

func TestParser_Parse_WithLogicalOperations(t *testing.T) {
	lexer := lex.CreateLexer("((1 + 2) == 3) and ((1 + 2) != 4 or (1 + 2) != 3)")
	tokens, _ := lexer.ScanTokens()
	parser := CreateParser(tokens)
	expression, err := parser.parseExpression()

	assert.Nil(t, err)
	assert.Equal(t, expression.Type(), LOGICAL_EXPR)
	assert.Equal(t, expression.(LogicalExpression).Operator.Type, lex.AND)

	right := expression.(LogicalExpression).Right
	left := expression.(LogicalExpression).Left

	assert.Equal(t, left.Type(), GROUPING_EXPR)
	assert.Equal(t, left.(GroupingExpression).OtherExpression.Type(), BINARY_EXPR)
	assert.Equal(t, left.(GroupingExpression).OtherExpression.(BinaryExpression).Token.Type, lex.EQUAL_EQUAL)

	assert.Equal(t, right.Type(), GROUPING_EXPR)
	assert.Equal(t, right.(GroupingExpression).OtherExpression.(LogicalExpression).Operator.Type, lex.OR)
	assert.Equal(t, right.(GroupingExpression).OtherExpression.(LogicalExpression).Left.Type(), BINARY_EXPR)
	assert.Equal(t, right.(GroupingExpression).OtherExpression.(LogicalExpression).Right.Type(), BINARY_EXPR)
}

func TestParser_Parse_StatementsWithExpressions(t *testing.T) {
	lexer := lex.CreateLexer("print \"hola\";\n" +
		"print 1 + 2 == 2;\n" +
		"print (7 * 3) == (7 + 7 + 7);")
	tokens, _ := lexer.ScanTokens()
	parser := CreateParser(tokens)

	statements, err := parser.Parse()

	assert.Nil(t, err)
	assert.Equal(t, len(statements), 3)

	assert.True(t, statements[0].Type() == PRINT_STMT && statements[1].Type() == PRINT_STMT && statements[2].Type() == PRINT_STMT)
	assert.Equal(t, statements[0].(PrintStatement).Expression.Type(), LITERAL_EXPR)
	assert.Equal(t, statements[1].(PrintStatement).Expression.Type(), BINARY_EXPR)
	assert.Equal(t, statements[2].(PrintStatement).Expression.Type(), BINARY_EXPR)
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
func TestParser_Parse_OnlyExpression(t *testing.T) {
	lexer := lex.CreateLexer("(1 + 2 - (1 / 2)) != false")
	tokens, _ := lexer.ScanTokens()

	parser := CreateParser(tokens)
	expr, _ := parser.parseExpression()

	assert.Equal(t, expr.Type(), BINARY_EXPR)

	assert.Equal(t, expr.(BinaryExpression).Right.Type(), LITERAL_EXPR)
	assert.Equal(t, expr.(BinaryExpression).Right.(LiteralExpression).Literal, false)

	assert.Equal(t, expr.(BinaryExpression).Left.Type(), GROUPING_EXPR)
	minus := expr.(BinaryExpression).Left.(GroupingExpression).OtherExpression.(BinaryExpression)
	assert.Equal(t, minus.Type(), BINARY_EXPR)

	assert.Equal(t, minus.Left.Type(), BINARY_EXPR)
	assert.Equal(t, minus.Right.Type(), GROUPING_EXPR)

	minus_slash := minus.Right.(GroupingExpression).OtherExpression.(BinaryExpression)
	assert.Equal(t, minus_slash.Type(), BINARY_EXPR)
	assert.Equal(t, minus_slash.Token.Type, lex.SLASH) //NOT
	assert.Equal(t, minus_slash.Left.Type(), LITERAL_EXPR)
	assert.Equal(t, minus_slash.Right.Type(), LITERAL_EXPR)
	assert.Equal(t, minus_slash.Left.(LiteralExpression).Literal, float64(1))
	assert.Equal(t, minus_slash.Right.(LiteralExpression).Literal, float64(2))

	minus_plus := minus.Left.(BinaryExpression)
	assert.Equal(t, minus_plus.Type(), BINARY_EXPR)
	assert.Equal(t, minus_plus.Token.Type, lex.PLUS)
	assert.Equal(t, minus_plus.Left.Type(), LITERAL_EXPR)
	assert.Equal(t, minus_plus.Right.Type(), LITERAL_EXPR)
	assert.Equal(t, minus_plus.Left.(LiteralExpression).Literal, float64(1))
	assert.Equal(t, minus_plus.Right.(LiteralExpression).Literal, float64(2))
}
