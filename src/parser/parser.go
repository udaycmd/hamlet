// Copyright 2025 Uday Tiwari. All rights reserved.
// Use of this source code is governed by MIT
// license that can be found in the LICENSE file.

package parser

import (
	"fmt"
	"sort"

	"github.com/udaycmd/hamlet/src/lexer"
	"github.com/udaycmd/hamlet/src/token"
)

type bailout struct{}

type ParseError struct {
	position token.SrcPos
	msg      string
}

func (e ParseError) Error() string {
	if e.position.FileName != "" || e.position.IsValid() {
		return fmt.Sprintf("Parsing Error: %s\n\tat %s", e.msg, e.position)
	}

	return fmt.Sprintf("Parsing Error: %s", e.msg)
}

// Slice of [ParseError]
type ParseErrors []*ParseError

func (pe ParseErrors) Extend(pos token.SrcPos, msg string) {
	pe = append(pe, &ParseError{position: pos, msg: msg})
}

func (pe ParseErrors) Len() int {
	return len(pe)
}

func (pe ParseErrors) Less(i, j int) bool {
	x := &pe[i].position
	y := &pe[j].position

	if x.FileName != y.FileName {
		return x.FileName < y.FileName
	}

	if x.Line != y.Line {
		return x.Line < y.Line
	}

	if x.Column != y.Column {
		return x.Column < y.Column
	}

	return false
}

func (pe ParseErrors) Swap(i, j int) {
	pe[i], pe[j] = pe[j], pe[i]
}

// [sort.Interface] implementation for [ParseErrors]
func (pe ParseErrors) Sort() {
	sort.Sort(pe)
}

func (pe ParseErrors) Error() string {
	n := len(pe)
	switch n {
	case 0:
		return ""
	case 1:
		return "TODO: error should be here"
	}

	return fmt.Sprintf("%s and %d more error(s)", pe[0].Error(), n-1)
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
	pos            token.Position      // current position in the [token.SourceManager]
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

func (p *Parser) parseStmtList() []Stmt {
	return nil
}

func (p *Parser) Parse() (*File, error) {
	var err error
	defer func() {
		if v := recover(); v != nil {
			if _, ok := v.(bailout); !ok {
				panic(v)
			}
		}

		p.errors.Sort()
		err = p.errors.err()
	}()

	stmts := p.parseStmtList()
	if p.errors.Len() > 0 {
		return nil, p.errors.err()
	}

	return &File{SrcFile: p.file, Statements: stmts}, err
}
