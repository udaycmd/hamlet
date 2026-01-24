// Copyright 2025 Uday Tiwari. All rights reserved.
// Use of this source code is governed by MIT
// license that can be found in the LICENSE file.

package parser

import (
	"fmt"

	"github.com/udaycmd/hamlet/src/lexer"
	"github.com/udaycmd/hamlet/src/token"
)

type bailout struct{}

type ParseError struct {
	position token.SrcPos
	msg      string
}

type ParseErrors []*ParseError

func (pe ParseErrors) Extend(pos token.SrcPos, msg string) {
	pe = append(pe, &ParseError{position: pos, msg: msg})
}

func (pe ParseErrors) Error() string {
	n := len(pe)
	switch n {
	case 0:
		return ""
	case 1:
		return "TODO: error should be here"
	}

	// TODO: change this
	return fmt.Sprintf("%s (and %d more errors)", pe[0].msg, n-1)
}

func (pe ParseErrors) err() error {
	if len(pe) == 0 {
		return nil
	}

	return pe
}

// Based upon Go's [parser] package
//
// [parser]: https://github.com/golang/go/blob/master/src/go/parser/parser.go
type Parser struct {
	file           *token.SourceHandle // file feeded to the parser
	errors         ParseErrors         // all parsing errors
	lexer          *lexer.Lexer        // lexer
	pos            token.Position      // current position in the SourceManager
	kind           token.Tok           // current token
	maxReportError int                 // maximum number of errors to report before parsing termination
	tokenLit       string              // current token's literal value
}

func NewParser(file *token.SourceHandle, src []byte, maxReportError int, mode lexer.LexMode) *Parser {
	p := &Parser{file: file, maxReportError: maxReportError}

	lexerErrorHandlerfunc := func(msg string, pos token.SrcPos) {
		p.errors.Extend(pos, msg)
	}
	p.lexer = lexer.NewLexer(p.file, src, lexerErrorHandlerfunc, mode)
	p.next()
	return p
}

func (p *Parser) next() {
	p.kind, p.tokenLit, p.pos = p.lexer.Lex()
}

func (p *Parser) error(pos token.Position, msg string) {
	srcpos := p.file.SrcPos(pos)

	n := len(p.errors)
	if n > 0 && p.errors[n-1].position.Line == srcpos.Line {
		return
	}

	if n > p.maxReportError {
		panic(bailout{})
	}
	p.errors.Extend(srcpos, msg)
}

func (p *Parser) Parse() (*File, error) {
	var err error
	defer func() {
		if v := recover(); v != nil {
			if _, ok := v.(bailout); !ok {
				panic(v)
			}
		}

		// sort errors based upon filename or line & column number
		err = p.errors.err()
	}()

	// TODO: change this
	file := &File{
		SrcFile:    p.file,
		Statements: nil,
	}
	return file, err
}
