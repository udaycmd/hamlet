// Copyright 2025 Uday Tiwari. All rights reserved.
// Use of this source code is governed by MIT
// license that can be found in the LICENSE file.

package lexer

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/udaycmd/hamlet/src/frontend/token"
)

const (
	eof = rune(-1) // eof
	bom = rune(0xFEFF)
)

const (
	NoAsi LexMode = iota + 1
	ParseComment
)

const (
	ErrUnexpectedNullChar        = "unexpected NULL character"
	ErrUnexpectedUnicodeChar     = "unexpected unicode codepoint"
	ErrIllegalBOM                = "illegal byte order mark"
	ErrUnterminatedStrLiteral    = "unterminated string literal"
	ErrUnterminatedCharLiteral   = "unterminated character literal"
	ErrUnterminatedEscapeSeq     = "unterminated escape sequence"
	ErrUnknownEscapeSeq          = "unknown escape sequence"
	ErrIllegalUnicodeEscape      = "illegal unicode escape inside byte sequence"
	ErrEmptyCharLiteral          = "zero width character literal encountered"
	ErrCharLiteralTooWide        = "character literal too wide"
	ErrNoDigitAfterExponent      = "expected digit(s) after exponent"
	ErrNoDigitAfterBaseSpecifier = "expected digit(s) after base specifier"
	ErrInvalidEllipsisMark       = "expected a '.' after '..'"
)

type (
	ErrorHandler func(message string, pos token.SrcPos)

	// lexing mode
	LexMode int

	// Based upon Go's [scanner] package
	//
	// [scanner]: https://github.com/golang/go/blob/master/src/go/scanner/scanner.go
	Lexer struct {
		file     *token.SourceHandle // source file handle
		err      ErrorHandler
		errCount int    // number of lexical errors
		path     string // abs path of file

		src      []byte // actual source
		cc       rune   // current character
		offset   int    // cc's offset in source
		rdOffset int    // position after current char
		asi      bool   // assume a ';' instead of '\n'
		mode     LexMode
	}
)

func NewLexer(file *token.SourceHandle, src []byte, err ErrorHandler, mode LexMode) *Lexer {
	if file.Len != len(src) {
		panic(fmt.Sprintf("file size %d does not match with source length %d", file.Len, len(src)))
	}

	l := &Lexer{
		file: file,
		src:  src,
		err:  err,
		cc:   ' ',
		mode: mode,
	}

	l.next()
	if l.cc == bom {
		l.next() // skip BOM
	}

	return l
}

func isDigit(r rune) bool {
	return ('0' <= r && r <= '9') || (r >= utf8.RuneSelf && unicode.IsDigit(r))
}

func digitVal(ch rune) int {
	switch {
	case '0' <= ch && ch <= '9':
		return int(ch - '0')
	case 'a' <= ch && ch <= 'f':
		return int(ch - 'a' + 10)
	case 'A' <= ch && ch <= 'F':
		return int(ch - 'A' + 10)
	}
	return 16
}

func isLetter(r rune) bool {
	return 'a' <= r && r <= 'z' || 'A' <= r && r <= 'Z' || r == '_' || (r >= utf8.RuneSelf && unicode.IsLetter(r))
}

func (l *Lexer) peek() byte {
	if l.rdOffset < len(l.src) {
		return l.src[l.rdOffset]
	}

	return 0
}

func (l *Lexer) skipWhitespace() {
	for l.cc == ' ' || l.cc == '\t' || l.cc == '\r' || (l.cc == '\n' && !l.asi) {
		l.next()
	}
}

func (l *Lexer) error(offset int, msg string) {
	if l.err != nil {
		l.err(msg, l.file.SrcPos(l.file.TapePos(offset)))
	}

	l.errCount++
}

func (l *Lexer) ErrCount() int {
	return l.errCount
}

func (l *Lexer) lexIdent() string {
	offset := l.offset
	for isLetter(l.cc) || isDigit(l.cc) {
		l.next()
	}

	return string(l.src[offset:l.offset])
}

func (l *Lexer) lexDigitSeq(base int) {
	for l.cc == '_' || digitVal(l.cc) < base {
		l.next()
	}
}

func (l *Lexer) lexNum() (tok token.Tok, lit string) {
	tok = token.INTEGER
	offset := l.offset
	base := 10

	switch peek := strings.ToLower(string(l.peek())); {
	case l.cc == '0' && peek == "b":
		base = 2
		l.next()
		l.next()
	case l.cc == '0' && peek == "o":
		base = 8
		l.next()
		l.next()
	case l.cc == '0' && peek == "x":
		base = 16
		l.next()
		l.next()
	}

	if base != 10 && digitVal(rune(l.peek())) == 16 {
		l.error(offset, ErrNoDigitAfterBaseSpecifier)
		goto ret
	}

	// lex whole number
	l.lexDigitSeq(base)

	// lex fractional
	if l.cc == '.' && base == 10 {
		tok = token.REAL
		l.next()
		l.lexDigitSeq(base)
	}

	// lex exponent
	if l.cc == 'e' || l.cc == 'E' {
		tok = token.REAL
		l.next()

		// lex exponent sign
		if l.cc == '-' || l.cc == '+' {
			l.next()
		}

		offset := l.offset
		l.lexDigitSeq(10)

		if offset == l.offset {
			l.error(offset, ErrNoDigitAfterExponent)
		}
	}

ret:
	return tok, string(l.src[offset:l.offset])
}

func (l *Lexer) lexEscape(quote rune) bool {
	offset := l.offset

	var n int
	var base, max uint32

	switch l.cc {
	case 'a', 'b', 'f', 'n', 'r', 't', 'v', '\\', quote:
		l.next()
		return true
	case '0', '1', '2', '3', '4', '5', '6', '7':
		n, base, max = 3, 8, 255
	case 'x':
		l.next()
		n, base, max = 2, 16, 255
	case 'u':
		l.next()
		n, base, max = 4, 16, unicode.MaxRune
	case 'U':
		l.next()
		n, base, max = 8, 16, unicode.MaxRune
	default:
		msg := ErrUnknownEscapeSeq
		if l.cc == eof {
			msg = ErrUnterminatedEscapeSeq
		}
		l.error(offset, msg)
		return false
	}

	var x uint32 = 0
	for n > 0 {
		d := uint32(digitVal(l.cc))
		if d >= base {
			msg := ErrIllegalUnicodeEscape
			if l.cc == eof {
				msg = ErrUnterminatedEscapeSeq
			}

			l.error(l.offset, msg)
			return false
		}

		x = x*base + d
		l.next()
		n--
	}

	if x > max || 0xD800 <= x && x < 0xE000 {
		l.error(offset, ErrIllegalUnicodeEscape)
		return false
	}

	return true
}

func (l *Lexer) lexStr() string {
	offset := l.offset - 1 // opening '"' already consumed

	for {
		c := l.cc
		if c == '\n' || c < 0 {
			l.error(offset, ErrUnterminatedStrLiteral)
			break
		}
		l.next()
		if c == '"' {
			break
		}
		if c == '\\' {
			l.lexEscape('"')
		}
	}

	return string(l.src[offset:l.offset])
}

func (l *Lexer) lexChar() string {
	offset := l.offset - 1 // opening '\'' already consumed

	valid := true
	n := 0

	for {
		c := l.cc
		if c == '\n' || c == eof {
			// only report error if we don't have one already
			if valid {
				l.error(offset, ErrUnterminatedCharLiteral)
				valid = false
			}
			break
		}
		l.next()
		if c == '\'' {
			if n == 0 {
				l.error(offset, ErrEmptyCharLiteral)
				valid = false
			}
			break
		}
		n++
		if c == '\\' {
			valid = l.lexEscape('\'')
		}
		// continue to read until closing quote
	}

	if valid && n != 1 {
		l.error(offset, ErrCharLiteralTooWide)
	}

	return string(l.src[offset:l.offset])
}

func (l *Lexer) lexComment() string {
	offset := l.offset - 1 // '#' already consumed

	for l.cc != '\n' && l.cc >= 0 {
		l.next()
	}

	return string(l.src[offset:l.offset])
}

// forwards the [Lexer.offset] and [Lexer.rdOffset]
func (l *Lexer) next() {
	if l.rdOffset < len(l.src) {
		l.offset = l.rdOffset

		if l.cc == '\n' {
			l.file.AddLine(l.offset)
		}

		r, w := rune(l.src[l.rdOffset]), 1
		switch {
		case r == 0:
			l.error(l.offset, ErrUnexpectedNullChar)
		case r >= utf8.RuneSelf:
			r, w = utf8.DecodeRune(l.src[l.rdOffset:])

			if r == utf8.RuneError && w == 1 {
				l.error(l.offset, ErrUnexpectedUnicodeChar)
			} else if r == bom && l.offset > 0 {
				l.error(l.offset, ErrIllegalBOM)
			}
		}

		l.rdOffset += w
		l.cc = r
	} else {
		l.offset = len(l.src)
		if l.cc == '\n' {
			l.file.AddLine(l.offset)
		}

		l.cc = eof
	}
}

func (l *Lexer) Lex() (token.Tok, string, token.Position) {
	var lit string
	var tok token.Tok

	l.skipWhitespace()

	// calculate tape position
	pos := l.file.TapePos(l.offset)

	asi := false

	switch c := l.cc; {
	case isLetter(c):
		lit = l.lexIdent()
		tok = token.IsIdent(lit)
		switch tok {
		case token.BREAK, token.CONTINUE, token.RETURN, token.IDENTIFIER,
			token.TRUE, token.FALSE, token.EMPTY:
			asi = true
		}
	case ('0' <= c && c <= '9') || (c == '.' && '0' <= l.peek() && l.peek() <= '9'):
		asi = true
		tok, lit = l.lexNum()
	default:
		l.next()

		switch c {
		case eof:
			if l.asi {
				l.asi = false // eof consumed
				return token.SEMICOLON, "\n", pos
			}

			tok = token.EOF
		case '\n':
			l.asi = false
			return token.SEMICOLON, "\n", pos
		case '"':
			asi = true
			tok = token.STRING
			lit = l.lexStr()
		case '\'':
			asi = true
			tok = token.CHAR
			lit = l.lexChar()
		case '#':
			// check if the comment is just after a 'asi' trigger
			if l.asi {
				l.cc = '#'
				l.offset = l.file.Offset(pos)
				l.rdOffset = l.offset + 1
				l.asi = false
				return token.SEMICOLON, "\n", pos
			}

			// parse comment
			comment := l.lexComment()
			if l.mode&ParseComment == 0 {
				l.asi = false
				return l.Lex()
			}

			tok = token.COMMENT
			lit = comment
		case '?':
			tok = token.QUESTION
		case '(':
			tok = token.LEFT_PAREN
		case ')':
			tok = token.RIGHT_PAREN
			asi = true
		case '[':
			tok = token.LEFT_BRACKET
		case ']':
			tok = token.RIGHT_BRACKET
			asi = true
		case '{':
			tok = token.LEFT_BRACE
		case '}':
			tok = token.RIGHT_BRACE
			asi = true
		case ',':
			tok = token.COMMA
		case ';':
			tok = token.SEMICOLON
			lit = ";"
		case ':':
			tok = token.COLON
			if l.cc == ':' {
				tok = token.DOUBLE_COLON
				l.next()
			}
		case '=':
			tok = token.ASSIGN
			if l.cc == '=' {
				tok = token.EQUALS
				l.next()
			}
		case '+':
			tok = token.PLUS
			if l.cc == '=' {
				tok = token.PLUS_EQ
				l.next()
			}
		case '*':
			tok = token.STAR
			if l.cc == '=' {
				tok = token.STAR_EQ
				l.next()
			}
		case '/':
			tok = token.SLASH
			if l.cc == '=' {
				tok = token.SLASH_EQ
				l.next()
			}
		case '%':
			tok = token.PERCENT
			if l.cc == '=' {
				tok = token.PERCENT_EQ
				l.next()
			}
		case '!':
			tok = token.BANG
			if l.cc == '=' {
				tok = token.BANG_EQ
				l.next()
			}
		case '^':
			tok = token.XOR
			if l.cc == '=' {
				tok = token.XOR_EQ
				l.next()
			}
		case '-':
			tok = token.MINUS
			switch l.cc {
			case '=':
				tok = token.MINUS_EQ
				l.next()
			case '>':
				tok = token.ARROW
				l.next()
			}
		case '>':
			tok = token.GREATER
			switch l.cc {
			case '>':
				tok = token.RIGHT_SHIFT
				l.next()
				if l.cc == '=' {
					tok = token.RSHIFT_EQ
					l.next()
				}
			case '=':
				tok = token.GREATER_EQ
				l.next()
			}
		case '<':
			tok = token.LESS
			switch l.cc {
			case '<':
				tok = token.LEFT_SHIFT
				l.next()
				if l.cc == '=' {
					tok = token.LSHIFT_EQ
					l.next()
				}
			case '=':
				tok = token.LESS_EQ
				l.next()
			}
		case '&':
			tok = token.AMPERSAND
			switch l.cc {
			case '=':
				tok = token.AND_EQ
				l.next()
			case '&':
				tok = token.AND
				l.next()
			}
		case '|':
			tok = token.PIPE
			switch l.cc {
			case '=':
				tok = token.OR_EQ
				l.next()
			case '|':
				tok = token.OR
				l.next()
			}
		case '.':
			tok = token.DOT
			switch l.cc {
			case '.':
				l.next()
				if l.cc == '.' {
					tok = token.DOT_DOT_DOT
					l.next()
				} else {
					l.error(l.file.Offset(pos), ErrInvalidEllipsisMark)
				}
			}
		default:
			if c != bom {
				l.error(l.file.Offset(pos), ErrUnexpectedUnicodeChar)
			}

			// preserve asi info in case of illegal token
			asi = l.asi
			lit = string(c)
			tok = token.INVALID
		}
	}

	// calculate literals for punctuators and operators
	if tok.IsOperator() || tok.IsPunctuator() {
		lit = tok.String()
	}

	if l.mode&NoAsi == 0 {
		l.asi = asi
	}
	return tok, lit, pos
}
