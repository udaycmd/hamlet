// Copyright 2026 Uday Tiwari. All rights reserved.
// Use of this source code is governed by MIT
// license that can be found in the LICENSE file.

package token

import "fmt"

type Tok uint8

const (
	EOF Tok = iota
	INVALID
	COMMENT

	// - Keywords -

	_KwBegin

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
	UNTIL
	DECL
	CONST
	TRUE
	FALSE
	EMPTY

	_KwEnd

	_OpBegin

	// - Operators -
	PLUS
	MINUS
	STAR
	SLASH
	PERCENT
	AMPERSAND
	PIPE
	XOR
	LEFT_SHIFT
	RIGHT_SHIFT
	PLUS_EQ
	MINUS_EQ
	STAR_EQ
	SLASH_EQ
	PERCENT_EQ
	XOR_EQ
	AND_EQ
	OR_EQ
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
	DOT_DOT_DOT

	_PuncEnd

	// - Literals -
	IDENTIFIER
	INTEGER
	REAL
	CHAR
	STRING
)

func (t Tok) IsOperator() bool {
	return t < _OpEnd && t > _OpBegin
}

func (t Tok) IsPunctuator() bool {
	return t < _PuncEnd && t > _PuncBegin
}

func (t Tok) IsKeyword() bool {
	return t < _KwEnd && t > _KwBegin
}

var tokens = [...]string{
	BREAK:         "break",
	CASE:          "case",
	CONTINUE:      "continue",
	DEFAULT:       "default",
	ELSE:          "else",
	ENUM:          "enum",
	PROC:          "proc",
	FOR:           "for",
	IMPORT:        "import",
	EXPORT:        "export",
	IF:            "if",
	IN:            "in",
	MAP:           "map",
	RETURN:        "return",
	CONCEPT:       "concept",
	SWITCH:        "switch",
	TYPE_INFO:     "type_info",
	UNTIL:         "until",
	DECL:          "decl",
	CONST:         "const",
	TRUE:          "true",
	FALSE:         "false",
	EMPTY:         "empty",
	ARROW:         "->",
	PLUS:          "+",
	MINUS:         "-",
	STAR:          "*",
	SLASH:         "/",
	PERCENT:       "%",
	AMPERSAND:     "&",
	PIPE:          "|",
	XOR:           "^",
	LEFT_SHIFT:    "<<",
	RIGHT_SHIFT:   ">>",
	PLUS_EQ:       "+=",
	MINUS_EQ:      "-=",
	STAR_EQ:       "*=",
	XOR_EQ:        "^=",
	SLASH_EQ:      "/=",
	PERCENT_EQ:    "%=",
	AND_EQ:        "&=",
	OR_EQ:         "|=",
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
	DOT_DOT_DOT:   "...",
	IDENTIFIER:    "identifier",
	INTEGER:       "integer",
	REAL:          "real",
	CHAR:          "char",
	STRING:        "string",
	EOF:           "eof",
}

var keywords map[string]Tok

func (t Tok) String() string {
	if int(t) >= len(tokens) || tokens[t] == "" {
		x := fmt.Sprintf("%d", t)
		panic("unknown token encountered" + x)
	}

	return tokens[t]
}

func IsIdent(name string) Tok {
	if kw, ok := keywords[name]; ok {
		return kw
	}

	return IDENTIFIER
}

func init() {
	keywords = make(map[string]Tok)
	for i := _KwBegin + 1; i < _KwEnd; i++ {
		keywords[tokens[i]] = i
	}
}
