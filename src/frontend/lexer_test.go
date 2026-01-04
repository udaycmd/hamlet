// Copyright 2025 Uday Tiwari. All rights reserved.
// Use of this source code is governed by MIT
// license that can be found in the LICENSE file.

package frontend_test

import (
	"strings"
	"testing"

	. "github.com/udaycmd/hamlet/src/frontend"
)

func lexToken(content string) Lexer {
	lexer := NewLexer(strings.NewReader(content))
	lexer.Next()
	return lexer
}

func TestTokens(t *testing.T) {
	expected := []struct {
		str  string
		kind Tok
	}{
		{"", EOF},
		{"\uFFFD", INVALID},

		// - Keywords -
		{"break", BREAK},
		{"case", CASE},
		{"continue", CONTINUE},
		{"default", DEFAULT},
		{"else", ELSE},
		{"enum", ENUM},
		{"fn", FN},
		{"for", FOR},
		{"import", IMPORT},
		{"export", EXPORT},
		{"interface", INTERFACE},
		{"if", IF},
		{"in", IN},
		{"map", MAP},
		{"return", RETURN},
		{"struct", STRUCT},
		{"switch", SWITCH},
		{"type", TYPE},
		{"var", VAR},

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
		{"->", ARROW},
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

		// - Literals -
		{"Hamlet", IDENTIFIER},
		{"$Hamlet", IDENTIFIER},
		{"_Hamlet", IDENTIFIER},
		{"\"Hamlet\"", STRING},
		{"'H'", CHAR},
		{"69", INTEGER},
		{"69.42", REAL},

		// - Comment -
		{"    # This is a single line comment with leading spaces.", EOF},
	}

	for _, exp := range expected {
		res := lexToken(exp.str).Token
		t.Run(exp.str, func(t *testing.T) {
			if res.Kind != exp.kind {
				t.Errorf("Test case failed because: '%s' != '%s'\n", res.Kind, exp.kind)
			}
		})
	}
}

func lexChar(t *testing.T, content, expected string) {
	t.Run(content, func(t *testing.T) {
		char := lexToken(content).Token.Value
		if char != expected {
			t.Errorf("Test case failed because: character('%s') != character('%s')\n", char, expected)
		}
	})
}

func lexErrChar(t *testing.T, content, lexErr string) {
	t.Run(content, func(t *testing.T) {
		err := lexToken(content).Err.Msg
		if err != lexErr {
			t.Errorf("Test case failed because: lexErr(\"%s\") != lexErr(\"%s\")\n", err, lexErr)
		}
	})
}

func TestCharLiterals(t *testing.T) {
	lexChar(t, "'H'", "H")
	lexChar(t, "'\n'", "\n")
	lexChar(t, "'\t'", "\t")
	lexChar(t, "'\r'", "\r")
	lexChar(t, "'\b'", "\b")

	lexErrChar(t, "'À'", "")
	lexErrChar(t, "'Б'", "")
	lexErrChar(t, "'\uFFFD'", ErrUnexpectedUnicodeChar)
	lexErrChar(t, "''", ErrEmptyCharLiteral)
	lexErrChar(t, "'", ErrUnterminatedCharLiteral)
	lexErrChar(t, "'👩🏻‍🤝‍👨🏾'", ErrCharLiteralTooWide)
}
