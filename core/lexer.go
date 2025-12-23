package core

import (
	"bufio"
	"errors"
	"io"
	"unicode"
)

const (
	maxLexBufSize = 8 * 1024 // 8 Kilobytes
	eof           = rune(-1) // eof
)

var (
	errUnexpectedUnicodeChar = errors.New("unexpected unicode codepoint")
	errInvalidIdentStart     = errors.New("invalid start of an identifier name")
)

type Tok uint8

const (
	EOF Tok = iota
	INVALID

	// - Keywords -
	BREAK
	CASE
	CONST
	CONTINUE
	DEFAULT
	ELSE
	ENUM
	FN
	FOR
	IMPORT
	INTERFACE
	IF
	IN
	MAP
	RETURN
	STR
	STRUCT
	SWITCH
	TYPE
	VAR
	WEAK

	// - Operators -
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
	case CONST:
		s = "const"
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
	case STR:
		s = "str"
	case STRUCT:
		s = "struct"
	case SWITCH:
		s = "switch"
	case TYPE:
		s = "type"
	case VAR:
		s = "var"
	case WEAK:
		s = "weak"
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
		s = "identifier"
	case INTEGER:
		s = "integer"
	case REAL:
		s = "real"
	case CHAR:
		s = "char"
	case STRING:
		s = "string"
	default:
		s = "unknown"
	}

	return s
}

func (t Tok) IsNumber() bool {
	return t == INTEGER || t == REAL
}

type Position struct {
	Line   uint32
	Column uint32
}

type Token struct {
	Offset uint32
	Start  uint32
	Width  uint32
	Kind   Tok
}

type Lexer struct {
	F         *bufio.Reader
	Pos       *Position
	Ppos      *Position
	Token     *Token
	Cursor    uint32
	Pcursor   uint32
	CodePoint rune
}

func isNewLineTerminator(r rune) bool {
	switch r {
	case '\u2028', '\u2029', '\n':
		return true
	default:
		return false
	}
}

func isDigit(r rune) bool {
	return unicode.IsDigit(r)
}

func isHexDigit(r rune) bool {
	return unicode.Is(unicode.ASCII_Hex_Digit, r)
}

func isWhitespace(r rune) bool {
	return unicode.Is(unicode.White_Space, r)
}

func isIdentStart(r rune) bool {
	return r == '$' || r == '_' || unicode.IsLetter(r)
}

func isIdentContinue(r rune) bool {
	return isIdentStart(r) || isDigit(r)
}

func NewLexer(file io.ReadCloser) *Lexer {
	return &Lexer{
		F: bufio.NewReaderSize(file, maxLexBufSize),
	}
}

func (l *Lexer) step() {
	l.Ppos = l.Pos
	l.Pcursor = l.Cursor

	codePoint, width, _ := l.F.ReadRune()
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

// Do not call twice without a l.step() in between
func (l *Lexer) back() {
	err := l.F.UnreadRune()
	if err != nil {
		panic("hamlet_crash: " + err.Error())
	}

	l.Pos = l.Ppos
	l.Cursor = l.Pcursor
}

func (l *Lexer) Next() error {
	t := &Token{}

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
		if l.CodePoint == '=' {
			t.Kind = MINUS_EQ
		} else {
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
		case '&':
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
	default:
		// skip whitespaces
		for isWhitespace(l.CodePoint) {
			l.step()
		}

		// identifiers
		if isIdentStart(l.CodePoint) {
			l.step()
			for isIdentContinue(l.CodePoint) {
				l.step()
			}

			t.Kind = IDENTIFIER
		} else {
			return errInvalidIdentStart
		}
	}

	l.Token = t
	return nil
}
