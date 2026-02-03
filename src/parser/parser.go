// Copyright 2025 Uday Tiwari. All rights reserved.
// Use of this source code is governed by MIT
// license that can be found in the LICENSE file.

package parser

import (
	"fmt"
	"io"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/udaycmd/hamlet/src/lexer"
	"github.com/udaycmd/hamlet/src/token"
)

type (
	bailout struct{}

	ParseError struct {
		position token.SrcPos
		msg      string
	}

	// Slice of *[ParseError]
	ParseErrors []*ParseError
)

var (
	stmtStarters = map[token.Tok]bool{
		token.PROC:     true,
		token.DECL:     true,
		token.FOR:      true,
		token.IF:       true,
		token.RETURN:   true,
		token.BREAK:    true,
		token.CONTINUE: true,
		token.EXPORT:   true,
	}
)

func (e *ParseError) Error() string {
	if e.position.FileName != "" || e.position.IsValid() {
		return fmt.Sprintf("Parsing Error: %s\n\tat %s", e.msg, e.position)
	}

	return fmt.Sprintf("Parsing Error: %s", e.msg)
}

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

	return x.Column < y.Column
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
	file     *token.SourceHandle // file feeded to the parser
	lexer    *lexer.Lexer        // lexer
	pos      token.Position      // current position in the [token.SourceManager]
	kind     token.Tok           // current token
	tokenLit string              // current token's literal value

	// errors
	errors         ParseErrors    // all parsing errors
	maxReportError int            // maximum number of errors to report before parsing termination
	syncPos        token.Position // last sync position
	syncCount      int            // no of synchronize() calls without progress

	// tracing
	tracing     bool      // do tracing?
	traceW      io.Writer // trace output
	traceIndent int       // trace indentation
}

func NewParser(
	file *token.SourceHandle,
	src []byte,
	maxReportError int,
	mode lexer.LexMode,
	traceW io.Writer,
) *Parser {
	p := &Parser{file: file, maxReportError: maxReportError, tracing: traceW != nil, traceW: traceW}

	lexerErrorHandlerfunc := func(msg string, pos token.SrcPos) {
		p.errors.Extend(pos, msg)
	}
	p.lexer = lexer.NewLexer(p.file, src, lexerErrorHandlerfunc, mode)
	p.next()
	return p
}

func (p *Parser) tracePrint(stringer ...any) {
	srcPos := p.file.SrcPos(p.pos)
	fmt.Fprintf(p.traceW, "%5d: line: %d column: %d ", p.pos, srcPos.Line, srcPos.Column)
	fmt.Fprint(p.traceW, strings.Repeat(" ", 2*p.traceIndent))
	fmt.Fprintln(p.traceW, stringer...)
}

func trace(p *Parser, msg string) *Parser {
	p.tracePrint(msg)
	p.traceIndent++
	return p
}

func untrace(p *Parser) {
	p.traceIndent--
}

func (p *Parser) next() {
	p.kind, p.tokenLit, p.pos = p.lexer.Lex()
}

func (p *Parser) synchronize(to map[token.Tok]bool) {
	for ; p.kind != token.EOF; p.next() {
		if to[p.kind] {
			if p.pos == p.syncPos && p.syncCount < 10 {
				p.syncCount++
				return
			}
			if p.pos > p.syncPos {
				p.syncPos = p.pos
				p.syncCount = 0
				return
			}
		}
	}
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

func (p *Parser) errorExpected(pos token.Position, msg string) {
	if pos == p.pos {
		switch {
		case p.kind == token.SEMICOLON && p.tokenLit == "\n":
			msg += ", found newline"
		default:
			msg += ", found '" + p.kind.String() + "'"
		}
	}

	p.error(pos, msg)
}

func (p *Parser) expect(kind token.Tok) token.Position {
	pos := p.pos

	if p.kind != kind {
		p.errorExpected(pos, "expected "+"'"+kind.String()+"'")
	}

	p.next()
	return pos
}

func (p *Parser) expectSemicolon() {
	switch p.kind {
	case token.SEMICOLON:
		p.next()
	case token.RIGHT_BRACE, token.RIGHT_PAREN:
		// semicolon is optional before a '}' or ')'
	default:
		p.errorExpected(p.pos, "';'")
		p.synchronize(stmtStarters)
	}
}

func (p *Parser) expectComma(gotAfter token.Tok, wantAfter string) bool {
	if p.kind == token.COMMA {
		p.next()

		if p.kind == gotAfter {
			p.errorExpected(p.pos, wantAfter)
			return false
		}
		return true
	}

	if p.kind == token.SEMICOLON && p.tokenLit == "\n" {
		p.next()
	}
	return false
}

func (p *Parser) parseSimpleStmt() Stmt {
	return nil
}

func (p *Parser) parseStmt() Stmt {
	if p.tracing {
		defer untrace(trace(p, "Stmt"))
	}

	switch p.kind {
	case token.PROC:
		return p.parseProcStmt()
	case token.RETURN:
		return p.parseReturnStmt()
	case token.SEMICOLON:
		s := &EmptyStmt{Semicolon: p.pos, IsImplicit: p.tokenLit == "\n"}
		p.next()
		return s
	case token.BREAK, token.CONTINUE:
		return p.parseBranchStmt(p.kind)
	}

	return nil
}

func (p *Parser) parseStmtList() []Stmt {
	var stmts []Stmt
	for p.kind != token.RIGHT_BRACE && p.kind != token.EOF {
		stmts = append(stmts, p.parseStmt())
	}

	return stmts
}

func (p *Parser) parseIdent() *Ident {
	if p.tracing {
		defer untrace(trace(p, "Identifier"))
	}

	pos := p.pos
	name := ""

	if p.kind == token.IDENTIFIER {
		name = p.tokenLit
		p.next()
	} else {
		return nil
	}

	return &Ident{
		Name: name,
		Pos:  pos,
	}
}

func (p *Parser) parseReturnStmt() Stmt {
	if p.tracing {
		defer untrace(trace(p, "ReturnStmt"))
	}

	pos := p.pos
	p.expect(token.RETURN)

	var x Expr
	if p.kind != token.SEMICOLON && p.kind != token.RIGHT_BRACE {
		x = p.parseExpr()
	}

	p.expectSemicolon()
	return &ReturnStmt{
		Ret: pos,
		e:   x,
	}
}

func (p *Parser) parseBranchStmt(tok token.Tok) Stmt {
	if p.tracing {
		defer untrace(trace(p, "BranchStmt"))
	}

	pos := p.expect(tok)
	var label *Ident
	if p.kind == token.IDENTIFIER {
		label = p.parseIdent()
	}

	p.expectSemicolon()
	return &BranchStmt{
		Kind:  p.kind,
		Pos:   pos,
		Label: label,
	}
}

func (p *Parser) parseIdentList() *IdentList {
	if p.tracing {
		defer untrace(trace(p, "IdentList"))
	}

	var params []*Ident
	lparen := p.expect(token.LEFT_PAREN)
	VarArgs := false

	if p.kind != token.RIGHT_PAREN {
		if p.kind == token.DOT_DOT_DOT {
			VarArgs = true
			p.next()
		}

		params = append(params, p.parseIdent())
		for !VarArgs && p.kind == token.COMMA {
			p.next()
			if p.kind == token.DOT_DOT_DOT {
				VarArgs = true
				p.next()
			}

			params = append(params, p.parseIdent())
		}
	}

	rparen := p.expect(token.RIGHT_PAREN)
	return &IdentList{
		LParen:  lparen,
		RParen:  rparen,
		VarArgs: VarArgs,
		List:    params,
	}
}

func (p *Parser) parseBlockStmt() *BlockStmt {
	if p.tracing {
		defer untrace(trace(p, "BlockStmt"))
	}

	lbrace := p.expect(token.LEFT_BRACE)
	stmts := p.parseStmtList()
	rbrace := p.expect(token.RIGHT_BRACE)

	return &BlockStmt{
		LBrace: lbrace,
		Stmts:  stmts,
		RBrace: rbrace,
	}
}

func (p *Parser) parseProcStmt() *ProcStmt {
	if p.tracing {
		defer untrace(trace(p, "ProcStmt"))
	}

	proc := p.expect(token.PROC)
	procName := p.parseIdent()
	params := p.parseIdentList()
	body := p.parseBlockStmt()

	return &ProcStmt{
		Proc:     proc,
		ProcName: procName,
		Params:   params,
		Body:     body,
	}
}

func (p *Parser) parseExpr() Expr {
	return nil
}

func (p *Parser) parseCallExpr(x Expr) *CallExpr {
	if p.tracing {
		defer untrace(trace(p, "CallExpr"))
	}

	lparen := p.expect(token.LEFT_PAREN)

	var list []Expr
	var ellipsis token.Position
	for p.kind != token.RIGHT_PAREN && p.kind != token.EOF && !ellipsis.IsValid() {
		list = append(list, p.parseExpr())
		if p.kind == token.DOT_DOT_DOT {
			ellipsis = p.pos
			p.next()
		}

		// catches a trailing comma with no following element
		if !p.expectComma(token.RIGHT_PAREN, "<call_argument>") {
			break
		}
	}

	rparen := p.expect(token.RIGHT_PAREN)
	return &CallExpr{
		Proc:     x,
		LParen:   lparen,
		RParen:   rparen,
		Ellipsis: ellipsis,
		Args:     list,
	}
}

func (p *Parser) parseImportExpr() Expr {
	if p.tracing {
		defer untrace(trace(p, "ImportExpr"))
	}

	pos := p.pos
	p.next()
	p.expect(token.LEFT_PAREN)

	if p.kind != token.STRING {
		p.errorExpected(p.pos, "<module_name>")
		p.synchronize(stmtStarters)
		return &BadExpr{From: pos, To: p.pos}
	}

	unquotedName, _ := strconv.Unquote(p.tokenLit)
	expr := &ImportExpr{
		ModuleName: unquotedName,
		Pos:        pos,
	}

	p.next()
	p.expect(token.RIGHT_PAREN)
	return expr
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

	if p.tracing {
		fullPath, _ := filepath.Abs(p.file.Name)
		fmt.Fprintf(p.traceW, "AST Trace of (%s)\n\n", fullPath)
		defer untrace(trace(p, "File"))
	}

	// if p.next() fails in NewParser()
	if p.errors.Len() > 0 {
		return nil, p.errors.err()
	}

	stmts := p.parseStmtList()
	if p.errors.Len() > 0 {
		return nil, p.errors.err()
	}

	return &File{SrcFile: p.file, Statements: stmts}, err
}
