package lex

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLexer_ScanTokens_Success(t *testing.T) {
	lexer := CreateLexer("//Mi comentario\n" +
		"var numeroA = 10.2;\n" +
		"var numeroB = 103;\n" +
		"var suma = numeroA + numeroB;\n\n" +
		"if (suma > 10) {\n" +
		"\tprint \"hola\";\n" +
		"}\n")

	tokens, err := lexer.ScanTokens()

	assert.Nil(t, err)

	assert.Equal(t, tokens[0].Type, VAR)
	assert.Equal(t, tokens[1].Type, IDENTIFIER)
	assert.Equal(t, tokens[2].Type, EQUAL)
	assert.Equal(t, tokens[3].Type, NUMBER)
	assert.Equal(t, tokens[4].Type, SEMICOLON)

	assert.Equal(t, tokens[5].Type, VAR)
	assert.Equal(t, tokens[6].Type, IDENTIFIER)
	assert.Equal(t, tokens[7].Type, EQUAL)
	assert.Equal(t, tokens[8].Type, NUMBER)
	assert.Equal(t, tokens[9].Type, SEMICOLON)

	assert.Equal(t, tokens[10].Type, VAR)
	assert.Equal(t, tokens[11].Type, IDENTIFIER)
	assert.Equal(t, tokens[12].Type, EQUAL)
	assert.Equal(t, tokens[13].Type, IDENTIFIER)
	assert.Equal(t, tokens[14].Type, PLUS)
	assert.Equal(t, tokens[15].Type, IDENTIFIER)
	assert.Equal(t, tokens[16].Type, SEMICOLON)

	assert.Equal(t, tokens[17].Type, IF)
	assert.Equal(t, tokens[18].Type, OPEN_PAREN)
	assert.Equal(t, tokens[19].Type, IDENTIFIER)
	assert.Equal(t, tokens[20].Type, GREATER)
	assert.Equal(t, tokens[21].Type, NUMBER)
	assert.Equal(t, tokens[22].Type, CLOSE_PAREN)
	assert.Equal(t, tokens[23].Type, OPEN_BRACE)

	assert.Equal(t, tokens[24].Type, PRINT)
	assert.Equal(t, tokens[25].Type, STRING)
	assert.Equal(t, tokens[26].Type, SEMICOLON)

	assert.Equal(t, tokens[27].Type, CLOSE_BRACE)
	assert.Equal(t, tokens[28].Type, EOF)

	fmt.Println(tokens)
}
