package core

import (
	"strings"
	"testing"
)

func lexToken(content string) *Token {
	lexer := NewLexer(strings.NewReader(content))
	lexer.Next()
	return lexer.Token
}

func TestTokens(t *testing.T) {
	expected := []struct {
		str  string
		kind Tok
	}{
		{"", EOF},
		{"\x00", INVALID},

		// - Keywords -
		{"break", BREAK},
		{"case", CASE},
		{"const", CONST},
		{"continue", CONTINUE},
		{"default", DEFAULT},
		{"else", ELSE},
		{"enum", ENUM},
		{"fn", FN},
		{"for", FOR},
		{"import", IMPORT},
		{"interface", INTERFACE},
		{"if", IF},
		{"in", IN},
		{"map", MAP},
		{"return", RETURN},
		{"str", STR},
		{"struct", STRUCT},
		{"switch", SWITCH},
		{"type", TYPE},
		{"var", VAR},
		{"weak", WEAK},

		// - Operators -
		{"+", PLUS},
		{"-", MINUS},
		{"*", STAR},
		{"/", SLASH},
		{"%", PERCENT},
		{"&", AMPERSAND},
		{"|", PIPE},
		{"~", TILDE},
		{"^", CARET},
		{"<<", LEFT_SHIFT},
		{">>", RIGHT_SHIFT},
		{"+=", PLUS_EQ},
		{"-=", MINUS_EQ},
		{"*=", STAR_EQ},
		{"/=", SLASH_EQ},
		{"%=", PERCENT_EQ},
		{"&=", AND_EQ},
		{"|=", OR_EQ},
		{"~=", NOT_EQ_BIT},
		{"<<=", LSHIFT_EQ},
		{">>=", RSHIFT_EQ},
		{":=", WALRUS},
		{"=", EQUAL},
		{"&&", AND},
		{"||", OR},
		{"!", BANG},
		{"!=", BANG_EQ},
		{"==", EQUAL_EQUAL},
		{"<", LESS},
		{">", GREATER},
		{"<=", LESS_EQ},
		{">=", GREATER_EQ},
		{"?", QUESTION},

		// - Punctuations -
		{"(", LEFT_PAREN},
		{")", RIGHT_PAREN},
		{"[", LEFT_BRACKET},
		{"]", RIGHT_BRACKET},
		{"{", LEFT_BRACE},
		{"}", RIGHT_BRACE},
		{",", COMMA},
		{";", SEMICOLON},
		{":", COLON},
		{"::", DOUBLE_COLON},
		{".", DOT},
		{"..", DOT_DOT},
		{"\"hello, \nWorld\"", STRING},
	}

	for _, exp := range expected {
		res := lexToken(exp.str)
		t.Run(exp.str, func(t *testing.T) {
			if res.Kind != exp.kind {
				t.Errorf("Test case failed because: '%s' != '%s'\n", res.Kind, exp.kind)
			}
		})
	}
}
