// Copyright 2025 Uday Tiwari. All rights reserved.
// Use of this source code is governed by MIT
// license that can be found in the LICENSE file.

package lexer_test

import (
	"fmt"
	"math/rand"
	"reflect"
	"strings"
	"testing"

	"github.com/udaycmd/hamlet/src/lexer"
	"github.com/udaycmd/hamlet/src/token"
)

var testSrcManager = token.NewSourceManager()

type lexResult struct {
	Lit    string
	Kind   token.Tok
	Line   int
	Column int
}

func fail(t *testing.T, msg string) {
	t.Errorf("%s", msg)
}

func expectEqual(t *testing.T, x, y any, msg string) {
	if !reflect.DeepEqual(x, y) {
		if msg == "" {
			msg = "unspecified failure!"
		}
		fail(t, msg)
	}
}

func countLines(s string) int {
	if s == "" {
		return 0
	}
	n := 1
	for i := 0; i < len(s); i++ {
		if s[i] == '\n' {
			n++
		}
	}
	return n
}

func lexExpect(t *testing.T, input string, mode lexer.LexMode, expected []lexResult) {
	testFile := testSrcManager.AddFile("testfile", -1, len(input))

	l := lexer.NewLexer(testFile, []byte(input), func(msg string, _ token.Position) { fail(t, msg) }, mode)

	for i, e := range expected {
		tok, literal, pos := l.Lex()

		srcPos := testFile.SrcPos(pos)
		expectEqual(t, e.Kind, tok, fmt.Sprintf("(%d) expected: %s, got: %s", i, e.Kind, tok))
		expectEqual(t, e.Lit, literal, "literal value not equal")
		expectEqual(t, e.Line, srcPos.Line, "line number not synchronized")
		expectEqual(t, e.Column, srcPos.Column, "column number not synchronized")
	}

	tok, _, _ := l.Lex()
	expectEqual(t, token.EOF, tok, "more tokens left")
	expectEqual(t, 0, l.ErrCount(), "error count not correct")
}

func TestTokens(t *testing.T) {
	testcases := []struct {
		token token.Tok
		lit   string
	}{
		{token.COMMENT, "# a comment \n"},
		{token.COMMENT, "#\r\n"},
		{token.IDENTIFIER, "foobar"},
		{token.IDENTIFIER, "a۰۱۸"},
		{token.IDENTIFIER, "foo६४"},
		{token.IDENTIFIER, "bar９８７６"},
		{token.IDENTIFIER, "ŝ"},
		{token.IDENTIFIER, "ŝfoo"},
		{token.INTEGER, "0"},
		{token.INTEGER, "1"},
		{token.INTEGER, "123456789012345678890"},
		{token.INTEGER, "01234567"},
		{token.INTEGER, "0xcafebabe"},
		{token.REAL, "0."},
		{token.REAL, ".0"},
		{token.REAL, "3.14159265"},
		{token.REAL, "1e0"},
		{token.REAL, "1e+100"},
		{token.REAL, "1e-100"},
		{token.REAL, "2.71828e-1000"},
		{token.CHAR, "'a'"},
		{token.CHAR, "'\\000'"},
		{token.CHAR, "'\\xFF'"},
		{token.CHAR, "'\\uff16'"},
		{token.CHAR, "'\\U0000ff16'"},
		{token.PLUS, "+"},
		{token.MINUS, "-"},
		{token.STAR, "*"},
		{token.SLASH, "/"},
		{token.PERCENT, "%"},
		{token.AMPERSAND, "&"},
		{token.PIPE, "|"},
		{token.CARET, "^"},
		{token.LEFT_SHIFT, "<<"},
		{token.RIGHT_SHIFT, ">>"},
		{token.PLUS_EQ, "+="},
		{token.MINUS_EQ, "-="},
		{token.STAR_EQ, "*="},
		{token.SLASH_EQ, "/="},
		{token.PERCENT_EQ, "%="},
		{token.AND_EQ, "&="},
		{token.OR_EQ, "|="},
		{token.NOT_EQ_BIT, "~="},
		{token.LSHIFT_EQ, "<<="},
		{token.RSHIFT_EQ, ">>="},
		{token.AND, "&&"},
		{token.OR, "||"},
		{token.EQUAL_EQUAL, "=="},
		{token.LESS, "<"},
		{token.GREATER, ">"},
		{token.EQUAL, "="},
		{token.BANG, "!"},
		{token.BANG_EQ, "!="},
		{token.LESS_EQ, "<="},
		{token.GREATER_EQ, ">="},
		{token.DOT, "."},
		{token.DOT_DOT, ".."},
		{token.LEFT_PAREN, "("},
		{token.LEFT_BRACKET, "["},
		{token.LEFT_BRACE, "{"},
		{token.COMMA, ","},
		{token.RIGHT_PAREN, ")"},
		{token.RIGHT_BRACKET, "]"},
		{token.RIGHT_BRACE, "}"},
		{token.SEMICOLON, ";"},
		{token.COLON, ":"},
		{token.BREAK, "break"},
		{token.CONTINUE, "continue"},
		{token.ELSE, "else"},
		{token.FOR, "for"},
		{token.FN, "fn"},
		{token.IF, "if"},
		{token.RETURN, "return"},
		{token.EXPORT, "export"},
	}

	var lines []string
	var lineSum int
	lineNumbers := make([]int, len(testcases))
	columnNumbers := make([]int, len(testcases))

	for i, tc := range testcases {
		// add extra lines before each test case
		emptyLines := rand.Intn(4)
		for j := 0; j < emptyLines; j++ {
			lines = append(lines, strings.Repeat(" ", rand.Intn(10)))
		}

		// add extra columns around each test case
		emptyColumns := rand.Intn(10)
		lines = append(lines, fmt.Sprintf("%s%s%s",
			strings.Repeat(" ", emptyColumns),
			tc.lit,
			strings.Repeat(" ", rand.Intn(10))))

		lineNumbers[i] = lineSum + emptyLines + 1
		lineSum += emptyLines + countLines(tc.lit)
		columnNumbers[i] = emptyColumns + 1
	}

	// expected results
	var expected []lexResult
	var expectedSkipComments []lexResult

	for i, tc := range testcases {
		// expected literal
		var expectedLiteral string
		switch tc.token {
		case token.COMMENT:
			// # style comment literal doesn't contain a '\n'
			expectedLiteral = tc.lit[:len(tc.lit)-1]
		case token.IDENTIFIER:
			expectedLiteral = tc.lit
		case token.SEMICOLON:
			expectedLiteral = ";"
		default:
			expectedLiteral = tc.lit
		}

		res := lexResult{
			Lit:    expectedLiteral,
			Kind:   tc.token,
			Line:   lineNumbers[i],
			Column: columnNumbers[i],
		}

		expected = append(expected, res)
		if tc.token != token.COMMENT {
			expectedSkipComments = append(expectedSkipComments, res)
		}
	}

	lexExpect(t, strings.Join(lines, "\n"), lexer.ParseComment|lexer.NoAsi, expected)
	lexExpect(t, strings.Join(lines, "\n"), lexer.NoAsi, expectedSkipComments)
}
