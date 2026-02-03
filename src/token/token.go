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
	PROC
	FOR
	IMPORT
	EXPORT
	IF
	IN
	MAP
	RETURN
	CONCEPT
	SWITCH
	TYPE_INFO
	DECL
	TRUE
	FALSE
	EMPTY

	_OpBegin

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
	ASSIGN
	EQUALS
	AND
	OR
	BANG
	BANG_EQ
	LESS
	GREATER
	LESS_EQ
	GREATER_EQ
	QUESTION

	_OpEnd

	_PuncBegin

	// - Punctuations -
	ARROW
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
	DOT_DOT_DOT

	_PuncEnd

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
	"proc":      PROC,
	"for":       FOR,
	"import":    IMPORT,
	"export":    EXPORT,
	"if":        IF,
	"in":        IN,
	"map":       MAP,
	"return":    RETURN,
	"concept":   CONCEPT,
	"switch":    SWITCH,
	"type_info": TYPE_INFO,
	"decl":      DECL,
	"true":      TRUE,
	"false":     FALSE,
	"empty":     EMPTY,
}

func (t Tok) IsOperator() bool {
	return t < _OpEnd && t > _OpBegin
}

func (t Tok) IsPunctuator() bool {
	return t < _PuncEnd && t > _PuncBegin
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
	EQUALS:        "==",
	ASSIGN:        "=",
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
	DOUBLE_COLON:  "::",
	DOT:           ".",
	DOT_DOT:       "..",
	DOT_DOT_DOT:   "...",
	IDENTIFIER:    "[identifier]",
	INTEGER:       "[integer]",
	REAL:          "[real]",
	CHAR:          "[char]",
	STRING:        "[string]",
	EOF:           "[EOF]",
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
