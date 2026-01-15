// Copyright 2025 Uday Tiwari. All rights reserved.
// Use of this source code is governed by MIT
// license that can be found in the LICENSE file.

package lexer_test

import (
	"strings"
	"testing"

	"github.com/udaycmd/hamlet/src/errors"
	. "github.com/udaycmd/hamlet/src/lexer"
	. "github.com/udaycmd/hamlet/src/token"
)

func lexSingleToken(content string) *Token {
	lexer := NewLexer(strings.NewReader(content))
	lexer.Next()
	return lexer.NxtToken
}

func lexErrSingleToken(content string) errors.Error {
	lexer := NewLexer(strings.NewReader(content))
	lexer.Next()
	return lexer.Err
}

func lexTokens(content string) Lexer {
	return NewLexer(strings.NewReader(content))
}

func TestTokens(t *testing.T) {
	cases := []struct {
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
		{"true", TRUE},
		{"false", FALSE},
		{"empty", EMPTY},

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

	for _, tc := range cases {
		res := lexSingleToken(tc.str).Kind

		t.Run(tc.str, func(t *testing.T) {
			if res != tc.kind {
				t.Errorf("Test case failed because: Tok.Kind('%s') != Expected.Kind('%s')\n", res, tc.kind)
			}
		})
	}
}

func lexChar(t *testing.T, content, expected string) {
	t.Run(content, func(t *testing.T) {
		char := lexSingleToken(content).Value
		if char != expected {
			t.Errorf("Test case failed because: Tok.character('%s') != Expected.character('%s')\n", char, expected)
		}
	})
}

func lexErrChar(t *testing.T, content, lexErr string) {
	t.Run(content, func(t *testing.T) {
		err := lexErrSingleToken(content).Msg
		if err != lexErr {
			t.Errorf("Test case failed because: Tok.lexErr(\"%s\") != Expected.lexErr(\"%s\")\n", err, lexErr)
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

func TestImplicitSemiColon(t *testing.T) {
	cases := []struct {
		str   string
		kinds []Tok
	}{
		// - Keywords that trigger implicit semicolon -
		{"return\n", []Tok{RETURN, IMPLICIT_SEMICOLON, EOF}},
		{"break\n", []Tok{BREAK, IMPLICIT_SEMICOLON, EOF}},
		{"continue\n", []Tok{CONTINUE, IMPLICIT_SEMICOLON, EOF}},

		// - Literals -
		{"123\n", []Tok{INTEGER, IMPLICIT_SEMICOLON, EOF}},
		{"45.67\n", []Tok{REAL, IMPLICIT_SEMICOLON, EOF}},
		{"\"string\"\n", []Tok{STRING, IMPLICIT_SEMICOLON, EOF}},
		{"'c'\n", []Tok{CHAR, IMPLICIT_SEMICOLON, EOF}},
		{"ident\n", []Tok{IDENTIFIER, IMPLICIT_SEMICOLON, EOF}},

		// - Brackets / Parens -
		{")\n", []Tok{RIGHT_PAREN, IMPLICIT_SEMICOLON, EOF}},
		{"]\n", []Tok{RIGHT_BRACKET, IMPLICIT_SEMICOLON, EOF}},
		{"}\n", []Tok{RIGHT_BRACE, IMPLICIT_SEMICOLON, EOF}},

		// - Operator -
		{"?\n", []Tok{QUESTION, IMPLICIT_SEMICOLON, EOF}},

		// - Cases NOT triggering implicit semicolon -
		{"+\n", []Tok{PLUS, EOL, EOF}},
		{"+=\n", []Tok{PLUS_EQ, EOL, EOF}},
		{"-\n", []Tok{MINUS, EOL, EOF}},
		{"(\n", []Tok{LEFT_PAREN, EOL, EOF}},
		{"{\n", []Tok{LEFT_BRACE, EOL, EOF}},
		{",\n", []Tok{COMMA, EOL, EOF}},
		{".\n", []Tok{DOT, EOL, EOF}},

		// - Multiple newlines -
		{"return\n\n", []Tok{RETURN, IMPLICIT_SEMICOLON, EOL, EOF}},
		{"return \n  \n", []Tok{RETURN, IMPLICIT_SEMICOLON, EOL, EOF}},
		{"return # this is a comment after keyword\n \n", []Tok{RETURN, IMPLICIT_SEMICOLON, EOL, EOF}},

		// - Mixed -
		{"a = 1\n b = 2", []Tok{IDENTIFIER, EQUAL, INTEGER, IMPLICIT_SEMICOLON, IDENTIFIER, EQUAL, INTEGER, EOF}},
		{"return\n 1 + 2", []Tok{RETURN, IMPLICIT_SEMICOLON, INTEGER, PLUS, INTEGER, EOF}},
	}

	for _, tc := range cases {
		t.Run(tc.str, func(t *testing.T) {
			lexer := lexTokens(tc.str)
			lexer.Next()

			for _, kind := range tc.kinds {
				if lexer.NxtToken.Kind != kind {
					t.Errorf("Test Case failed because: Tok.Kind('%s') != Expected.Kind('%s')", lexer.NxtToken.Kind, kind)
				}
				lexer.Next()
			}
		})
	}
}

func TestPositions(t *testing.T) {
	cases := []struct {
		str       string
		positions []Position
	}{
		{
			"var x = 1",
			[]Position{
				{Line: 0, Column: 1, Offset: 1},
				{Line: 0, Column: 5, Offset: 5},
				{Line: 0, Column: 7, Offset: 7},
				{Line: 0, Column: 9, Offset: 9},
				{Line: 0, Column: 10, Offset: 9},
			},
		},
		{
			"\ta",
			[]Position{
				{Line: 0, Column: 6, Offset: 2},
				{Line: 0, Column: 7, Offset: 2},
			},
		},
		{
			"a\nb",
			[]Position{
				{Line: 0, Column: 1, Offset: 1},
				{Line: 1, Column: 0, Offset: 2},
				{Line: 1, Column: 1, Offset: 3},
				{Line: 1, Column: 2, Offset: 3},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.str, func(t *testing.T) {
			lexer := lexTokens(tc.str)
			lexer.Next()

			for _, pos := range tc.positions {
				if lexer.NxtToken.Pos != pos {
					t.Errorf("Test Case failed because: Pos(%v) != Expected(%v). Token: %s", lexer.NxtToken.Pos, pos, lexer.NxtToken.Kind)
				}
				lexer.Next()
			}
		})
	}
}
