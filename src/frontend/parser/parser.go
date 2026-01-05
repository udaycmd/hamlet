// Copyright 2025 Uday Tiwari. All rights reserved.
// Use of this source code is governed by MIT
// license that can be found in the LICENSE file.

package parser

import (
	"io"

	. "github.com/udaycmd/hamlet/src/frontend"
)

type parser struct {
	lexer  Lexer
	Errors []Error
}

func NewParser(source io.Reader) *parser {
	return &parser{
		lexer: NewLexer(source),
	}
}

func (p *parser) curr() Tok {
	return p.lexer.PrevToken.Kind
}

func (p *parser) advance() {
	if p.lexer.Token.Kind == EOF {
		return
	}

	p.lexer.Next()
}

func (p *parser) expectNext(expected Tok) bool {
	if p.lexer.Token.Kind == expected {
		p.advance()
		return true
	}

	return false
}

func (p *parser) parseDecl() Decl {
	switch p.curr() {
	case FN:
		return p.parseFnDecl(false)
	case EXPORT:
		return p.parseExportDecl()
	default:
		panic("unimplemented")
	}
}

func (p *parser) parseExportDecl() Decl {
	p.advance()

	switch p.curr() {
	case FN:
		return p.parseFnDecl(true)
	default:
		panic("unimplemented")
	}
}

func (p *parser) parseFnDecl(exported bool) Decl {
	decl := FuncDecl{IsExported: exported}

	if !p.expectNext(IDENTIFIER) {
		// TODO: Better error report!
		return nil
	}

	if !p.expectNext(LEFT_PAREN) {
		// TODO: Better error report!
		return nil
	}

	// TODO: use this in decl above
	_ = p.parseFnParamDecl()

	if !p.expectNext(ARROW) {
		// TODO: Better error report!
		return nil
	}

	if !p.expectNext(LEFT_BRACE) {
		// TODO: Better error report!
		return nil
	}

	decl.Body = p.parseBlockStmt()

	return decl
}

func (p *parser) parseFnParamDecl() Decl {
	return nil
}

func (p *parser) parseBlockStmt() *BlockStmt {
	return nil
}
