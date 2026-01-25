// Copyright 2026 Uday Tiwari. All rights reserved.
// Use of this source code is governed by MIT
// license that can be found in the LICENSE file.

package token

type Tok uint8

const (
	EOF Tok = iota
	INVALID
	COMMENT

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
	TRUE
	FALSE
	EMPTY

	_OpBegin

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
	EQUAL
	ASSIGN
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

	_OpEnd

	_PunBegin

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
	DOT
	DOT_DOT
	DOT_DOT_DOT

	_PunEnd

	// - Literals -
	IDENTIFIER
	INTEGER
	REAL
	CHAR
	STRING
)

var keywords = map[string]Tok{
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
	"true":      TRUE,
	"false":     FALSE,
	"empty":     EMPTY,
}

func (t Tok) IsOperator() bool {
	return t < _OpEnd && t > _OpBegin
}

func (t Tok) IsPunctuator() bool {
	return t < _PunEnd && t > _PunBegin
}

var tokens = [...]string{
	ARROW:         "->",
	PLUS:          "+",
	MINUS:         "-",
	STAR:          "*",
	SLASH:         "/",
	PERCENT:       "%",
	AMPERSAND:     "&",
	PIPE:          "|",
	TILDE:         "~",
	CARET:         "^",
	LEFT_SHIFT:    "<<",
	RIGHT_SHIFT:   ">>",
	PLUS_EQ:       "+=",
	MINUS_EQ:      "-=",
	STAR_EQ:       "*=",
	SLASH_EQ:      "/=",
	PERCENT_EQ:    "%=",
	AND_EQ:        "&=",
	OR_EQ:         "|=",
	NOT_EQ_BIT:    "~=",
	LSHIFT_EQ:     "<<=",
	RSHIFT_EQ:     ">>=",
	EQUAL:         "=",
	ASSIGN:        "::",
	AND:           "&&",
	OR:            "||",
	BANG:          "!",
	BANG_EQ:       "!=",
	LESS:          "<",
	GREATER:       ">",
	LESS_EQ:       "<=",
	GREATER_EQ:    ">=",
	QUESTION:      "?",
	LEFT_PAREN:    "(",
	RIGHT_PAREN:   ")",
	LEFT_BRACKET:  "[",
	RIGHT_BRACKET: "]",
	LEFT_BRACE:    "{",
	RIGHT_BRACE:   "}",
	COMMA:         ",",
	SEMICOLON:     ";",
	COLON:         ":",
	DOT:           ".",
	DOT_DOT:       "..",
	DOT_DOT_DOT:   "...",
	IDENTIFIER:    "[identifier]",
	INTEGER:       "[integer]",
	REAL:          "[real]",
	CHAR:          "[char]",
	STRING:        "[string]",
}

func (t Tok) String() string {
	if int(t) >= len(tokens) || tokens[t] == "" {
		panic("unknown token encountered")
	}

	return tokens[t]
}

func IsKeyword(name string) Tok {
	if kw, ok := keywords[name]; ok {
		return kw
	}

	return IDENTIFIER
}
