package tests

import (
	"strings"
	"testing"

	"github.com/udaycmd/hamlet/core"
)

func lexToken(content string) *core.Token {
	lexer := core.NewLexer(strings.NewReader(content))
	lexer.Next()
	return lexer.Token
}

func TestTokens(t *testing.T) {
	expected := []struct {
		str  string
		kind core.Tok
	}{
		{"", core.EOF},
		{"\x00", core.INVALID},

		// - Keywords -
		{"break", core.BREAK},
		{"case", core.CASE},
		{"const", core.CONST},
		{"continue", core.CONTINUE},
		{"default", core.DEFAULT},
		{"else", core.ELSE},
		{"enum", core.ENUM},
		{"fn", core.FN},
		{"for", core.FOR},
		{"import", core.IMPORT},
		{"interface", core.INTERFACE},
		{"if", core.IF},
		{"in", core.IN},
		{"map", core.MAP},
		{"return", core.RETURN},
		{"str", core.STR},
		{"struct", core.STRUCT},
		{"switch", core.SWITCH},
		{"type", core.TYPE},
		{"var", core.VAR},
		{"weak", core.WEAK},

		// - Operators -
		{"+", core.PLUS},
		{"-", core.MINUS},
		{"*", core.STAR},
		{"/", core.SLASH},
		{"%", core.PERCENT},
		{"&", core.AMPERSAND},
		{"|", core.PIPE},
		{"~", core.TILDE},
		{"^", core.CARET},
		{"<<", core.LEFT_SHIFT},
		{">>", core.RIGHT_SHIFT},
		{"+=", core.PLUS_EQ},
		{"-=", core.MINUS_EQ},
		{"*=", core.STAR_EQ},
		{"/=", core.SLASH_EQ},
		{"%=", core.PERCENT_EQ},
		{"&=", core.AND_EQ},
		{"|=", core.OR_EQ},
		{"~=", core.NOT_EQ_BIT},
		{"<<=", core.LSHIFT_EQ},
		{">>=", core.RSHIFT_EQ},
		{":=", core.WALRUS},
		{"=", core.EQUAL},
		{"&&", core.AND},
		{"||", core.OR},
		{"!", core.BANG},
		{"!=", core.BANG_EQ},
		{"==", core.EQUAL_EQUAL},
		{"<", core.LESS},
		{">", core.GREATER},
		{"<=", core.LESS_EQ},
		{">=", core.GREATER_EQ},
		{"?", core.QUESTION},

		// - Punctuations -
		{"(", core.LEFT_PAREN},
		{")", core.RIGHT_PAREN},
		{"[", core.LEFT_BRACKET},
		{"]", core.RIGHT_BRACKET},
		{"{", core.LEFT_BRACE},
		{"}", core.RIGHT_BRACE},
		{",", core.COMMA},
		{";", core.SEMICOLON},
		{":", core.COLON},
		{"::", core.DOUBLE_COLON},
		{".", core.DOT},
		{"..", core.DOT_DOT},
	}

	for _, exp := range expected {
		res := lexToken(exp.str)
		t.Run(exp.str, func(t *testing.T) {
			if res.Kind != exp.kind {
				t.Errorf("Failed '%s' != '%s'\n", res.Kind, exp.kind)
			}
		})
	}
}
