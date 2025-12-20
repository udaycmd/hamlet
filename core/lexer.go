package core

import (
	"bufio"
	"io"
	"unicode"
)

const (
	maxLexBufSize = 8 * 1024 // 8 Kilobytes
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
	PLUS_PLUS
	MINUS_MINUS
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
	case PLUS_PLUS:
		s = "++"
	case MINUS_MINUS:
		s = "--"
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

type Token struct {
	Offset uint32
	Start  uint32
	Width  uint32
	Kind   Tok
}

type Lexer struct {
	IsModule bool
	F        *bufio.Reader
	Crune    rune // current rune
	Cursor   uint32
	Len      uint32
	Line     uint32
	Column   uint32
	Ctoken   *Tok // current
}

func (l *Lexer) isNewLineTerminator() bool {
	switch l.Crune {
	case '\u2028', '\u2029', '\n', '\r':
		return true
	default:
		return false
	}
}

func (l *Lexer) isDigit() bool {
	return unicode.IsDigit(l.Crune)
}

func (l *Lexer) isHexDigit() bool {
	return unicode.Is(unicode.ASCII_Hex_Digit, l.Crune)
}

func (l *Lexer) isWhitespace() bool {
	return unicode.Is(unicode.White_Space, l.Crune)
}

func (l *Lexer) isIdentStart() bool {
	return l.Crune == '$' || l.Crune == '_' || unicode.IsLetter(l.Crune)
}

func (l *Lexer) isIdentContinue() bool {
	return l.isIdentStart() || l.isDigit()
}

func NewLexer(file io.ReadCloser, isMod bool) *Lexer {
	return &Lexer{
		IsModule: isMod,
		F:        bufio.NewReaderSize(file, maxLexBufSize),
	}
}

func (l *Lexer) Next() *Token {

}
