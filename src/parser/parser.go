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
		return fmt.Sprintf("Parsing Error: %s\tat %s", e.msg, e.position)
	}

	return fmt.Sprintf("Parsing Error: %s", e.msg)
}

func (pe *ParseErrors) Extend(pos token.SrcPos, msg string) {
	*pe = append(*pe, &ParseError{position: pos, msg: msg})
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
		return "empty"
	case 1:
		return pe[0].Error()
	default:
		return fmt.Sprintf("%s and %d more error(s)", pe[0].Error(), n-1)
	}
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
	msg = "expected " + msg

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
		p.errorExpected(pos, "'"+kind.String()+"'")
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
	case token.EXPORT:
		return p.parseExportStmt()
	case token.PROC:
		return p.parseProcStmt()
	case token.DECL:
		return p.parseDeclStmt()
	case token.RETURN:
		return p.parseReturnStmt()
	case token.SEMICOLON, token.RIGHT_BRACE:
		s := &EmptyStmt{Semicolon: p.pos, IsImplicit: (p.kind == token.RIGHT_BRACE || p.tokenLit == "\n")}
		p.next()
		return s
	case token.BREAK, token.CONTINUE:
		return p.parseBranchStmt(p.kind)
	default:
		pos := p.pos
		p.errorExpected(pos, "<statement>")
		p.synchronize(stmtStarters)
		return &BadStmt{From: pos, To: pos}
	}
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
		p.expect(token.IDENTIFIER)
	}

	return &Ident{
		Name: name,
		Pos:  pos,
	}
}

func (p *Parser) parseReturnStmt() *ReturnStmt {
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

func (p *Parser) parseDeclStmt() *DeclStmt {
	if p.tracing {
		defer untrace(trace(p, "DeclStmt"))
	}

	decl := p.expect(token.DECL)
	ident := p.parseIdent()
	equals := p.expect(token.ASSIGN)
	x := p.parseExpr()

	return &DeclStmt{
		Decl:  decl,
		Ident: ident,
		Equal: equals,
		Val:   x,
	}
}

func (p *Parser) parseBranchStmt(tok token.Tok) *BranchStmt {
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

func (p *Parser) parseProcLit() Expr {
	if p.tracing {
		defer untrace(trace(p, "ProcLit"))
	}

	proc := p.expect(token.PROC)
	params := p.parseIdentList()
	body := p.parseBlockStmt()

	return &ProcLit{
		Proc:   proc,
		Params: params,
		Body:   body,
	}
}

func (p *Parser) parseExpr() Expr {
	if p.tracing {
		defer untrace(trace(p, "Expr"))
	}

	x := p.parseEqualityExpr()

	if p.kind == token.QUESTION {
		return p.parseTernaryExpr(x)
	}

	return x
}

func (p *Parser) parseEqualityExpr() Expr {
	if p.tracing {
		defer untrace(trace(p, "EqualityExpr"))
	}

	lhs := p.parseComparisonExpr()
	for p.kind == token.EQUALS || p.kind == token.BANG_EQ {
		opKind, opPos := p.kind, p.pos
		p.next()

		rhs := p.parseComparisonExpr()
		lhs = &BinaryExpr{
			Lhs:   lhs,
			Op:    opKind,
			OpPos: opPos,
			Rhs:   rhs,
		}
	}

	return lhs
}

func (p *Parser) parseComparisonExpr() Expr {
	if p.tracing {
		defer untrace(trace(p, "ComparisonExpr"))
	}

	lhs := p.parseLogicalExpr()
	for p.kind == token.LESS || p.kind == token.GREATER ||
		p.kind == token.LESS_EQ || p.kind == token.GREATER_EQ {
		opKind, opPos := p.kind, p.pos
		p.next()

		rhs := p.parseLogicalExpr()
		lhs = &BinaryExpr{
			Lhs:   lhs,
			Op:    opKind,
			OpPos: opPos,
			Rhs:   rhs,
		}
	}

	return lhs
}

func (p *Parser) parseLogicalExpr() Expr {
	if p.tracing {
		defer untrace(trace(p, "LogicalExpr"))
	}

	lhs := p.parseTermExpr()
	for p.kind == token.AND || p.kind == token.OR ||
		p.kind == token.CARET || p.kind == token.AMPERSAND || p.kind == token.PIPE {
		opKind, opPos := p.kind, p.pos
		p.next()

		rhs := p.parseTermExpr()
		lhs = &BinaryExpr{
			Lhs:   lhs,
			Op:    opKind,
			OpPos: opPos,
			Rhs:   rhs,
		}
	}

	return lhs
}

func (p *Parser) parseTermExpr() Expr {
	if p.tracing {
		defer untrace(trace(p, "TermExpr"))
	}

	lhs := p.parseFactorizedExpr()
	for p.kind == token.PLUS || p.kind == token.MINUS {
		opKind, opPos := p.kind, p.pos
		p.next()

		rhs := p.parseFactorizedExpr()
		lhs = &BinaryExpr{
			Lhs:   lhs,
			Op:    opKind,
			OpPos: opPos,
			Rhs:   rhs,
		}
	}

	return lhs
}

func (p *Parser) parseFactorizedExpr() Expr {
	if p.tracing {
		defer untrace(trace(p, "FactorizedExpr"))
	}

	lhs := p.parseUnaryExpr()
	for p.kind == token.STAR || p.kind == token.SLASH {
		opKind, opPos := p.kind, p.pos
		p.next()

		rhs := p.parseUnaryExpr()
		lhs = &BinaryExpr{
			Lhs:   lhs,
			Op:    opKind,
			OpPos: opPos,
			Rhs:   rhs,
		}
	}

	return lhs
}

func (p *Parser) parseUnaryExpr() Expr {
	if p.tracing {
		defer untrace(trace(p, "UnaryExpr"))
	}

	switch p.kind {
	case token.TILDE, token.BANG:
		opKind, opPos := p.kind, p.pos
		p.next()

		x := p.parseUnaryExpr()
		return &UnaryExpr{
			Op:    opKind,
			OpPos: opPos,
			X:     x,
		}
	}

	return p.parsePrimaryExpr()
}

func (p *Parser) parsePrimaryExpr() Expr {
	if p.tracing {
		defer untrace(trace(p, "PrimaryExpr"))
	}

	x := p.parseLiteralOrSubExpr()

loop:
	for {
		switch p.kind {
		case token.LEFT_BRACKET:
			x = p.parseIndexerExpr(x)
		case token.LEFT_PAREN:
			x = p.parseCallExpr(x)
		case token.DOT:
			p.next()

			switch p.kind {
			case token.IDENTIFIER:
				x = p.parseReceiverExpr(x)
			default:
				pos := p.pos
				p.errorExpected(pos, "<receiver>")
				p.synchronize(stmtStarters)
				return &BadExpr{From: pos, To: p.pos}
			}
		default:
			break loop
		}
	}

	return x
}

func (p *Parser) parseIndexerExpr(x Expr) Expr {
	if p.tracing {
		defer untrace(trace(p, "IndexerExpr"))
	}

	lbrack := p.expect(token.LEFT_BRACKET)

	// [e1:e2]
	var indices [2]Expr
	if p.kind != token.COLON {
		// parse e1, if present
		indices[0] = p.parseExpr()
	}

	colons := 0
	if p.kind == token.COLON {
		colons += 1
		p.next()

		// parse e2, if present
		if p.kind != token.LEFT_BRACKET && p.kind != token.EOF {
			indices[1] = p.parseExpr()
		}
	}

	rbrack := p.expect(token.RIGHT_BRACKET)

	// a slice, if it has a colon present
	if colons > 0 {
		return &SliceExpr{
			X:     x,
			LBrac: lbrack,
			Lo:    indices[0],
			Hi:    indices[1],
			Rbrac: rbrack,
		}
	}

	return &IndexExpr{
		X:     x,
		LBrac: lbrack,
		Index: indices[0],
		Rbrac: rbrack,
	}
}

func (p *Parser) parseReceiverExpr(x Expr) Expr {
	if p.tracing {
		defer untrace(trace(p, "ReceiverExpr"))
	}

	recv := p.parseIdent()
	return &ReceiverExpr{
		X:  x,
		Id: recv,
	}
}

func (p *Parser) parseTernaryExpr(cond Expr) *TernaryExpr {
	question := p.expect(token.QUESTION)
	e1 := p.parseExpr()
	colon := p.expect(token.COLON)
	e2 := p.parseExpr()

	return &TernaryExpr{
		X:     cond,
		Ques:  question,
		True:  e1,
		Colon: colon,
		False: e2,
	}
}

func (p *Parser) parseLiteralOrSubExpr() Expr {
	if p.tracing {
		defer untrace(trace(p, "LiteralOrSubExpr"))
	}

	switch p.kind {
	case token.EMPTY:
		x := &EmptyLit{Pos: p.pos}

		p.next()
		return x
	case token.INTEGER:
		i, err := strconv.ParseInt(p.tokenLit, 0, 64)
		if err != nil {
			if err == strconv.ErrRange {
				p.error(p.pos, "<integer_too_big>")
			} else {
				p.error(p.pos, "<invalid_integer>")
			}
		}
		x := &IntLit{
			Val:     i,
			Literal: p.tokenLit,
			Pos:     p.pos,
		}

		p.next()
		return x
	case token.REAL:
		f, err := strconv.ParseFloat(p.tokenLit, 64)
		if err != nil {
			if err == strconv.ErrRange {
				p.error(p.pos, "<number_too_big>")
			} else {
				p.error(p.pos, "<invalid_number>")
			}
		}
		x := &FloatLit{
			Val:     f,
			Literal: p.tokenLit,
			Pos:     p.pos,
		}

		p.next()
		return x
	case token.CHAR:
		return p.parseCharLit()
	case token.STRING:
		unquoted, _ := strconv.Unquote(p.tokenLit)
		x := &StringLit{
			Val:     unquoted,
			Literal: p.tokenLit,
			Pos:     p.pos,
		}

		p.next()
		return x
	case token.TRUE, token.FALSE:
		x := &BoolLit{
			Val:     (p.kind == token.TRUE),
			Literal: p.tokenLit,
			Pos:     p.pos,
		}

		p.next()
		return x
	case token.PROC:
		return p.parseProcLit()
	case token.IMPORT:
		return p.parseImportExpr()
	case token.IDENTIFIER:
		return p.parseIdent()
	case token.LEFT_PAREN:
		lparen := p.pos
		p.next()
		x := p.parseExpr()
		rparen := p.expect(token.RIGHT_PAREN)

		return &GroupedExpr{
			Lparen: lparen,
			X:      x,
			Rparen: rparen,
		}
	case token.LEFT_BRACE:
		return p.parseMapLiteral()
	case token.LEFT_BRACKET:
		return p.parseArrayLiteral()
	default:
		p.errorExpected(p.pos, "<expr>")
	}

	// something is wrong!
	pos := p.pos
	p.synchronize(stmtStarters)
	return &BadExpr{From: pos, To: p.pos}
}

func (p *Parser) parseKvLit() *KvLit {
	if p.tracing {
		defer untrace(trace(p, "KvLit"))
	}

	pos := p.pos
	key := ""

	// valid map keys are just like ecmascript's object keys
	switch p.kind {
	case token.IDENTIFIER:
		key = p.tokenLit
	case token.STRING:
		v, _ := strconv.Unquote(p.tokenLit)
		key = v
	default:
		p.errorExpected(pos, "<valid_key>")
	}

	p.next()
	colon := p.expect(token.COLON)
	value := p.parseExpr()
	return &KvLit{
		KeyPos: pos,
		Key:    key,
		Colon:  colon,
		Value:  value,
	}
}

func (p *Parser) parseMapLiteral() *MapLit {
	if p.tracing {
		defer untrace(trace(p, "MapLit"))
	}

	lbrace := p.expect(token.LEFT_BRACE)

	var kv []*KvLit
	for p.kind != token.RIGHT_BRACE && p.kind != token.EOF {
		kv = append(kv, p.parseKvLit())

		// catches a trailing comma with no following KV pair
		if !p.expectComma(token.RIGHT_BRACE, "<key_val_pair>") {
			break
		}
	}

	rbrace := p.expect(token.RIGHT_BRACE)
	return &MapLit{
		LBrace: lbrace,
		Kvs:    kv,
		RBrace: rbrace,
	}
}

func (p *Parser) parseArrayLiteral() *ArrayLit {
	if p.tracing {
		defer untrace(trace(p, "ArrayLit"))
	}

	lbrack := p.expect(token.LEFT_BRACKET)

	var items []Expr
	for p.kind != token.RIGHT_BRACKET && p.kind != token.EOF {
		items = append(items, p.parseExpr())

		// catches a trailing comma with no following array element
		if !p.expectComma(token.RIGHT_BRACKET, "<array_element>") {
			break
		}
	}

	rbrack := p.expect(token.RIGHT_BRACKET)
	return &ArrayLit{
		LBrack: lbrack,
		Items:  items,
		RBrack: rbrack,
	}
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

func (p *Parser) parseCharLit() Expr {
	if n := len(p.tokenLit); n >= 3 {
		cp, _, _, err := strconv.UnquoteChar(p.tokenLit[1:n-1], '\'')
		if err == nil {
			x := &CharLit{
				Val:     cp,
				Literal: p.tokenLit,
				Pos:     p.pos,
			}

			p.next()
			return x
		}
	}

	pos := p.pos
	p.error(pos, "invalid character literal")
	p.next()
	return &BadExpr{
		From: pos,
		To:   p.pos,
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

func (p *Parser) parseExportStmt() Stmt {
	if p.tracing {
		defer untrace(trace(p, "ExportStmt"))
	}

	pos := p.pos
	p.expect(token.EXPORT)
	x := p.parseExpr()
	p.expectSemicolon()
	return &ExportStmt{
		Pos: pos,
		e:   x,
	}
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
		fmt.Fprintf(p.traceW, "AST Trace of [%s]\n\n", fullPath)
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
