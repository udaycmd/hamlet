// Copyright 2025 Uday Tiwari. All rights reserved.
// Use of this source code is governed by MIT
// license that can be found in the LICENSE file.

package parser

// import (
// 	"io"

// 	"github.com/udaycmd/hamlet/src/errors"
// 	"github.com/udaycmd/hamlet/src/lexer"
// 	"github.com/udaycmd/hamlet/src/token"
// )

// type parser struct {
// 	lexer  lexer.Lexer
// 	Errors []errors.Error
// }

// func NewParser(source io.Reader) *parser {
// 	// TODO: remove this
// 	l := lexer.NewLexer(source)
// 	l.Next()

// 	return &parser{
// 		lexer: l,
// 	}
// }

// func (p *parser) over() bool {
// 	return p.lexer.Token.Kind == token.EOF
// }

// func (p *parser) curr() *token.Token {
// 	return p.lexer.Token
// }

// func (p *parser) advance() {
// 	if p.over() {
// 		return
// 	}

// 	p.lexer.Next()
// }

// func (p *parser) expectNext(expected token.Tok) bool {
// 	if p.lexer.Token.Kind == expected {
// 		p.advance()
// 		return true
// 	}

// 	return false
// }

// func (p *parser) parseDecl() Decl {
// 	curr := p.curr()

// 	switch curr.Kind {
// 	case token.FN:
// 		return p.parseFnDecl(false)
// 	case token.EXPORT:
// 		return p.parseExportDecl()
// 	default:
// 		panic("unimplemented")
// 	}
// }

// func (p *parser) parseExportDecl() Decl {
// 	p.advance()
// 	curr := p.curr()

// 	switch curr.Kind {
// 	case token.FN:
// 		return p.parseFnDecl(true)
// 	default:
// 		panic("unimplemented")
// 	}
// }

// func (p *parser) parseFnDecl(exported bool) Decl {
// 	decl := FuncDecl{IsExported: exported}

// 	if !p.expectNext(token.IDENTIFIER) {
// 		// TODO: Better error report!
// 		return nil
// 	}

// 	if !p.expectNext(token.LEFT_PAREN) {
// 		// TODO: Better error report!
// 		return nil
// 	}

// 	// TODO: use this in decl above
// 	_ = p.parseFnParamDecl()

// 	if !p.expectNext(token.ARROW) {
// 		// TODO: Better error report!
// 		return nil
// 	}

// 	if !p.expectNext(token.LEFT_BRACE) {
// 		// TODO: Better error report!
// 		return nil
// 	}

// 	decl.Body = p.parseBlockStmt()

// 	return nil
// }

// func (p *parser) parseFnParamDecl() Decl {
// 	return nil
// }

// func (p *parser) parseBlockStmt() *BlockStmt {
// 	return nil
// }

// // Expressions

// func (p *parser) parseExpr() Expr {
// 	return p.parseEqualityExpr()
// }

// func (p *parser) parseEqualityExpr() Expr {
// 	return nil
// }

// func (p *parser) parseComparisonExpr() Expr {
// 	return nil
// }

// func (p *parser) parseTermExpr() Expr {
// 	return nil
// }

// func (p *parser) parseFactorizedExpr() Expr {
// 	return nil
// }

// func (p *parser) parseUnaryExpr() Expr {
// 	return nil
// }

// func (p *parser) parsePrimaryExpr() Expr {
// 	var x Expr
// 	curr := p.curr()

// 	switch curr.Kind {
// 	case token.TRUE, token.FALSE, token.EMPTY, token.REAL,
// 		token.INTEGER, token.STRING, token.CHAR:
// 		x = BasicLit{
// 			Lit: *curr,
// 		}
// 	case token.LEFT_PAREN:
// 		expr := p.parseExpr()

// 		// expect RPAREN
// 		p.expectNext(token.RIGHT_PAREN)
// 		x = GroupExpr{X: expr}
// 	}

// 	return x
// }
