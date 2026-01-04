// Copyright 2025 Uday Tiwari. All rights reserved.
// Use of this source code is governed by MIT
// license that can be found in the LICENSE file.

package frontend

import (
	"bufio"
	"io"
	"strings"
	"unicode"
)

const (
	maxLexBufSize = 32 * 1024 // 32 Kilobytes
	eof           = rune(-1)  // eof
)

var (
	ErrUnexpectedUnicodeChar    = "unexpected unicode codepoint"
	ErrUnterminatedStrLiteral   = "unterminated string literal"
	ErrUnterminatedCharLiteral  = "unterminated character literal"
	ErrEmptyCharLiteral         = "zero width character literal encountered"
	ErrCharLiteralTooWide       = "character literal too wide"
	ErrNoDigitsAfterBasePrefix  = "expected digit(s) after base prefix"
	ErrNoDigitAfterDecimalPoint = "expected digit(s) after decimal point"
	ErrNoDigitAfterExponent     = "expected digit(s) after exponent"
)

type (
	Tok   uint8
	Radix uint8

	Position struct {
		Line   uint32
		Column uint32
	}

	Token struct {
		Value string
		Kind  Tok
		Pos   Position
	}

	Lexer struct {
		F                  *bufio.Reader
		Buf                *strings.Builder
		Pos                Position
		Token, PrevToken   *Token
		Cursor, PrevCursor uint32
		Err                Error
		CodePoint          rune
	}
)

const (
	base2 Radix = iota
	base8
	base10
	base16
)

const (
	EOF Tok = iota
	INVALID

	// - Keywords -
	BREAK
	CASE
	CONTINUE
	DEFAULT
	ELSE
	ENUM
	FN
	FOR
	IMPORT
	EXPORT
	INTERFACE
	IF
	IN
	MAP
	RETURN
	STRUCT
	SWITCH
	TYPE
	VAR

	// - Operators -
	ARROW
	PLUS
	MINUS
	STAR
	SLASH
	PERCENT
	AMPERSAND
	PIPE
	TILDE
	CARET
	LEFT_SHIFT
	RIGHT_SHIFT
	PLUS_EQ
	MINUS_EQ
	STAR_EQ
	SLASH_EQ
	PERCENT_EQ
	AND_EQ
	OR_EQ
	NOT_EQ_BIT
	LSHIFT_EQ
	RSHIFT_EQ
	WALRUS
	EQUAL
	AND
	OR
	BANG
	BANG_EQ
	EQUAL_EQUAL
	LESS
	GREATER
	LESS_EQ
	GREATER_EQ
	QUESTION

	// - Punctuations -
	LEFT_PAREN
	RIGHT_PAREN
	LEFT_BRACKET
	RIGHT_BRACKET
	LEFT_BRACE
	RIGHT_BRACE
	COMMA
	SEMICOLON
	COLON
	DOUBLE_COLON
	DOT
	DOT_DOT

	// - Literals -
	IDENTIFIER
	INTEGER
	REAL
	CHAR
	STRING
)

var Keywords = map[string]Tok{
	"break":     BREAK,
	"case":      CASE,
	"continue":  CONTINUE,
	"default":   DEFAULT,
	"else":      ELSE,
	"enum":      ENUM,
	"fn":        FN,
	"for":       FOR,
	"import":    IMPORT,
	"export":    EXPORT,
	"interface": INTERFACE,
	"if":        IF,
	"in":        IN,
	"map":       MAP,
	"return":    RETURN,
	"struct":    STRUCT,
	"switch":    SWITCH,
	"type":      TYPE,
	"var":       VAR,
}

func (t Tok) String() string {
	var s string

	switch t {
	case EOF:
		s = "eof"
	case INVALID:
		s = "invalid"
	case BREAK:
		s = "break"
	case CASE:
		s = "case"
	case CONTINUE:
		s = "continue"
	case DEFAULT:
		s = "default"
	case ELSE:
		s = "else"
	case ENUM:
		s = "enum"
	case FN:
		s = "fn"
	case FOR:
		s = "for"
	case IMPORT:
		s = "import"
	case EXPORT:
		s = "export"
	case INTERFACE:
		s = "interface"
	case IF:
		s = "if"
	case IN:
		s = "in"
	case MAP:
		s = "map"
	case RETURN:
		s = "return"
	case STRUCT:
		s = "struct"
	case SWITCH:
		s = "switch"
	case TYPE:
		s = "type"
	case VAR:
		s = "var"
	case ARROW:
		s = "->"
	case PLUS:
		s = "+"
	case MINUS:
		s = "-"
	case STAR:
		s = "*"
	case SLASH:
		s = "/"
	case PERCENT:
		s = "%"
	case AMPERSAND:
		s = "&"
	case PIPE:
		s = "|"
	case TILDE:
		s = "~"
	case CARET:
		s = "^"
	case LEFT_SHIFT:
		s = "<<"
	case RIGHT_SHIFT:
		s = ">>"
	case PLUS_EQ:
		s = "+="
	case MINUS_EQ:
		s = "-="
	case STAR_EQ:
		s = "*="
	case SLASH_EQ:
		s = "/="
	case PERCENT_EQ:
		s = "%="
	case AND_EQ:
		s = "&="
	case OR_EQ:
		s = "|="
	case NOT_EQ_BIT:
		s = "~="
	case LSHIFT_EQ:
		s = "<<="
	case RSHIFT_EQ:
		s = ">>="
	case WALRUS:
		s = ":="
	case EQUAL:
		s = "="
	case AND:
		s = "&&"
	case OR:
		s = "||"
	case BANG:
		s = "!"
	case BANG_EQ:
		s = "!="
	case EQUAL_EQUAL:
		s = "=="
	case LESS:
		s = "<"
	case GREATER:
		s = ">"
	case LESS_EQ:
		s = "<="
	case GREATER_EQ:
		s = ">="
	case QUESTION:
		s = "?"
	case LEFT_PAREN:
		s = "("
	case RIGHT_PAREN:
		s = ")"
	case LEFT_BRACKET:
		s = "["
	case RIGHT_BRACKET:
		s = "]"
	case LEFT_BRACE:
		s = "{"
	case RIGHT_BRACE:
		s = "}"
	case COMMA:
		s = ","
	case SEMICOLON:
		s = ";"
	case COLON:
		s = ":"
	case DOUBLE_COLON:
		s = "::"
	case DOT:
		s = "."
	case DOT_DOT:
		s = ".."
	case IDENTIFIER:
		s = "[identifier]"
	case INTEGER:
		s = "[integer]"
	case REAL:
		s = "[real]"
	case CHAR:
		s = "[char]"
	case STRING:
		s = "[string]"
	default:
		s = "[unknown]"
	}

	return s
}

func (t Tok) IsNumber() bool {
	return t == INTEGER || t == REAL
}

func isNewLineTerminator(r rune) bool {
	switch r {
	case '\u2028', '\u2029', '\n':
		return true
	default:
		return false
	}
}

func isHexDigit(r rune) bool {
	return unicode.Is(unicode.ASCII_Hex_Digit, r)
}

func isDigit(r rune, base Radix) bool {
	switch base {
	case base2:
		return r == '0' || r == '1'
	case base8:
		return r >= '0' && r <= '7'
	case base10:
		return unicode.IsDigit(r)
	case base16:
		return isHexDigit(r)
	}

	return false
}

func isWhitespace(r rune) bool {
	return unicode.Is(unicode.White_Space, r)
}

func isIdentStart(r rune) bool {
	return r == '$' || r == '_' || unicode.IsLetter(r)
}

func isIdentContinue(r rune) bool {
	return isIdentStart(r) || isDigit(r, base10)
}

func NewLexer(file io.Reader) Lexer {
	return Lexer{
		F:         bufio.NewReaderSize(file, maxLexBufSize),
		Buf:       &strings.Builder{},
		Token:     &Token{Kind: INVALID},
		PrevToken: &Token{Kind: INVALID},
		Pos:       Position{Line: 1, Column: 1},
	}
}

// Forwards the [Lexer.F] by one rune.
func (l *Lexer) step() {
	l.PrevCursor = l.Cursor
	l.PrevToken = l.Token

	codePoint, width, err := l.F.ReadRune()
	if err != nil && err != io.EOF {
		panic("hamlet_crash: " + err.Error())
	}

	if width == 0 {
		codePoint = eof
	}

	if isNewLineTerminator(codePoint) {
		l.Pos.Line++
		l.Pos.Column = 0
	} else {
		if codePoint == '\t' {
			l.Pos.Column += 4
		}
		l.Pos.Column++
	}

	l.CodePoint = codePoint
	l.Cursor += uint32(width)
}

// Back tracks the [Lexer.F] by one rune.
// Do not call twice without a [Lexer.step] in between.
func (l *Lexer) back() {
	if l.CodePoint == eof {
		return
	}

	err := l.F.UnreadRune()
	if err != nil {
		panic("hamlet_crash: " + err.Error())
	}

	l.Pos = l.PrevToken.Pos
	l.Cursor = l.PrevCursor
}

func (l *Lexer) lexIdent() (string, string) {
	l.Buf.Reset()

	for isIdentContinue(l.CodePoint) {
		_, err := l.Buf.WriteRune(l.CodePoint)
		if err != nil {
			return l.Buf.String(), err.Error()
		}
		l.step()
	}

	return l.Buf.String(), ""
}

// [Lexer.lexStr] only identifies any lexical error while scanning a string literal.
// It does not do any kind of string validation and parsing.
func (l *Lexer) lexStr() (string, string) {
	// skip the starting `"`
	l.step()

	// skip the ending `"`
	defer l.step()

	l.Buf.Reset()

	for l.CodePoint != '"' {
		if l.CodePoint == eof {
			return l.Buf.String(), ErrUnterminatedStrLiteral
		}

		if l.CodePoint == unicode.ReplacementChar {
			return l.Buf.String(), ErrUnexpectedUnicodeChar
		}

		_, err := l.Buf.WriteRune(l.CodePoint)
		if err != nil {
			return l.Buf.String(), err.Error()
		}
		l.step()
	}

	return l.Buf.String(), ""
}

func (l *Lexer) lexChar() (string, string) {
	// skip the starting `'`
	l.step()

	// skip the ending `'`
	defer l.step()

	l.Buf.Reset()
	charCount := 0
	for l.CodePoint != '\'' {
		if l.CodePoint == eof {
			return l.Buf.String(), ErrUnterminatedCharLiteral
		}

		if l.CodePoint == unicode.ReplacementChar {
			return l.Buf.String(), ErrUnexpectedUnicodeChar
		}

		_, err := l.Buf.WriteRune(l.CodePoint)
		if err != nil {
			return l.Buf.String(), err.Error()
		}

		charCount++
		l.step()
	}

	if l.Buf.Len() == 0 {
		return l.Buf.String(), ErrEmptyCharLiteral
	}

	if !(l.Buf.Len() <= 4 && charCount == 1) {
		return l.Buf.String(), ErrCharLiteralTooWide
	}

	return l.Buf.String(), ""
}

func (l *Lexer) lexDigitSeq(base Radix) string {
	for isDigit(l.CodePoint, base) || l.CodePoint == '_' {
		_, err := l.Buf.WriteRune(l.CodePoint)
		if err != nil {
			return err.Error()
		}
		l.step()
	}

	return ""
}

// Parse a number (int or real) and return it as string or a parse error.
// Referenced from [umka-lang].
//
// [umka-lang]: https://github.com/vtereshkov/umka-lang.git
func (l *Lexer) lexNum() (string, Tok, string) {
	l.Buf.Reset()
	base := base10
	isReal := false

	if l.CodePoint == '0' {
		_, err := l.Buf.WriteRune(l.CodePoint)
		if err != nil {
			return l.Buf.String(), INVALID, err.Error()
		}

		l.step()
		nbase := base10
		switch l.CodePoint {
		case 'x', 'X':
			nbase = base16
		case 'o', 'O':
			nbase = base8
		case 'b', 'B':
			nbase = base2
		}

		if nbase != base10 {
			base = nbase
			_, err := l.Buf.WriteRune(l.CodePoint)
			if err != nil {
				return l.Buf.String(), INVALID, err.Error()
			}

			l.step()
			if !isDigit(l.CodePoint, base) {
				return l.Buf.String(), INVALID, ErrNoDigitsAfterBasePrefix
			}
		}
	}

	l.lexDigitSeq(base)

	if base != base10 {
		return l.Buf.String(), INTEGER, ""
	}

	if l.CodePoint == '.' {
		l.step()

		// found `..` (range) operator
		if l.CodePoint == '.' {
			l.back()
			return l.Buf.String(), INTEGER, ""
		}

		isReal = true
		_, err := l.Buf.WriteRune('.')
		if err != nil {
			return l.Buf.String(), INVALID, ""
		}

		if !isDigit(l.CodePoint, base10) {
			return l.Buf.String(), INVALID, ErrNoDigitAfterDecimalPoint
		}

		l.lexDigitSeq(base10)
	}

	if l.CodePoint == 'e' || l.CodePoint == 'E' {
		isReal = true
		_, err := l.Buf.WriteRune(l.CodePoint)
		if err != nil {
			return l.Buf.String(), INVALID, ""
		}

		l.step()
		if l.CodePoint == '+' || l.CodePoint == '-' {
			_, err := l.Buf.WriteRune(l.CodePoint)
			if err != nil {
				return l.Buf.String(), INVALID, ""
			}
			l.step()
		}

		if !isDigit(l.CodePoint, base10) {
			return l.Buf.String(), INVALID, ErrNoDigitAfterExponent
		}

		l.lexDigitSeq(base10)
	}

	if isReal {
		return l.Buf.String(), REAL, ""
	}

	return l.Buf.String(), INTEGER, ""
}

func (l *Lexer) lexWhiteSpaceAndComment() {
	l.step()

	for isWhitespace(l.CodePoint) || l.CodePoint == '#' {
		if l.CodePoint == '#' {
			for !isNewLineTerminator(l.CodePoint) && l.CodePoint != eof {
				l.step()
			}
		}

		l.step()
	}
}

func (l *Lexer) Next() {
	t := &Token{Kind: INVALID}

	// skip whitespaces and new lines
	l.lexWhiteSpaceAndComment()

	switch l.CodePoint {
	case eof:
		t.Kind = EOF
	case '+':
		l.step()
		if l.CodePoint == '=' {
			t.Kind = PLUS_EQ
		} else {
			l.back()
			t.Kind = PLUS
		}
	case '-':
		l.step()
		switch l.CodePoint {
		case '=':
			t.Kind = MINUS_EQ
		case '>':
			t.Kind = ARROW
		default:
			l.back()
			t.Kind = MINUS
		}
	case '*':
		l.step()
		if l.CodePoint == '=' {
			t.Kind = STAR_EQ
		} else {
			l.back()
			t.Kind = STAR
		}
	case '/':
		l.step()
		if l.CodePoint == '=' {
			t.Kind = SLASH_EQ
		} else {
			l.back()
			t.Kind = SLASH
		}
	case '%':
		l.step()
		if l.CodePoint == '=' {
			t.Kind = PERCENT_EQ
		} else {
			l.back()
			t.Kind = PERCENT
		}
	case '&':
		l.step()
		switch l.CodePoint {
		case '=':
			t.Kind = AND_EQ
		case '&':
			t.Kind = AND
		default:
			l.back()
			t.Kind = AMPERSAND
		}
	case '|':
		l.step()
		switch l.CodePoint {
		case '=':
			t.Kind = OR_EQ
		case '|':
			t.Kind = OR
		default:
			l.back()
			t.Kind = PIPE
		}
	case '~':
		l.step()
		switch l.CodePoint {
		case '=':
			t.Kind = NOT_EQ_BIT
		default:
			l.back()
			t.Kind = TILDE
		}
	case '=':
		l.step()
		if l.CodePoint == '=' {
			t.Kind = EQUAL_EQUAL
		} else {
			l.back()
			t.Kind = EQUAL
		}
	case '!':
		l.step()
		if l.CodePoint == '=' {
			t.Kind = BANG_EQ
		} else {
			l.back()
			t.Kind = BANG
		}
	case '.':
		l.step()
		if l.CodePoint == '.' {
			t.Kind = DOT_DOT
		} else {
			l.back()
			t.Kind = DOT
		}
	case ':':
		l.step()
		switch l.CodePoint {
		case ':':
			t.Kind = DOUBLE_COLON
		case '=':
			t.Kind = WALRUS
		default:
			l.back()
			t.Kind = COLON
		}
	case '<':
		l.step()
		switch l.CodePoint {
		case '<':
			l.step()
			if l.CodePoint == '=' {
				t.Kind = LSHIFT_EQ
			} else {
				l.back()
				t.Kind = LEFT_SHIFT
			}
		case '=':
			t.Kind = LESS_EQ
		default:
			l.back()
			t.Kind = LESS
		}
	case '>':
		l.step()
		switch l.CodePoint {
		case '>':
			l.step()
			if l.CodePoint == '=' {
				t.Kind = RSHIFT_EQ
			} else {
				l.back()
				t.Kind = RIGHT_SHIFT
			}
		case '=':
			t.Kind = GREATER_EQ
		default:
			l.back()
			t.Kind = GREATER
		}
	case '^':
		t.Kind = CARET
	case '?':
		t.Kind = QUESTION
	case '(':
		t.Kind = LEFT_PAREN
	case ')':
		t.Kind = RIGHT_PAREN
	case '[':
		t.Kind = LEFT_BRACKET
	case ']':
		t.Kind = RIGHT_BRACKET
	case '{':
		t.Kind = LEFT_BRACE
	case '}':
		t.Kind = RIGHT_BRACE
	case ',':
		t.Kind = COMMA
	case ';':
		t.Kind = SEMICOLON
	case '"':
		v, err := l.lexStr()
		if err != "" {
			l.Err.Msg = err
		}

		t.Kind = STRING
		t.Value = v
	case '\'':
		v, err := l.lexChar()
		if err != "" {
			l.Err.Msg = err
		}

		t.Kind = CHAR
		t.Value = v
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		v, kind, err := l.lexNum()
		if err != "" {
			l.Err.Msg = err
		}

		t.Kind = kind
		t.Value = v
	default:
		if isIdentStart(l.CodePoint) { // identifiers
			name, err := l.lexIdent()
			if err != "" {
				l.Err.Msg = err
			}

			t.Value = name
			if Keywords[name] != 0 {
				t.Kind = Keywords[name]
			} else {
				t.Kind = IDENTIFIER
			}
		} else if l.CodePoint == unicode.ReplacementChar { // invalid unicode codepoint
			l.Err.Msg = ErrUnexpectedUnicodeChar
			l.step()
		}
	}

	l.Token = t
}
