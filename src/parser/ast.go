// Copyright 2025 Uday Tiwari. All rights reserved.
// Use of this source code is governed by MIT
// license that can be found in the LICENSE file.

package parser

import "github.com/udaycmd/hamlet/src/token"

type Node interface {
	Start() token.Position
	End() token.Position
}

type (
	Expr interface {
		Node
		exprNode()
	}

	// Statement is an ast node representing
	// a complete instruction or action.
	Stmt interface {
		Node
		stmtNode()
	}

	// Declaration is a ast node that represents the introduction of an
	// identifier / type / function / imports of a particular module.
	Decl interface {
		Node
		declNode()
	}
)

type (
	BlockStmt struct {
		Statements []Stmt
	}
)

func (s *BlockStmt) NodeVal() string { return "block" }
func (s *BlockStmt) stmtNode()       {}

type (
	FuncDecl struct {
		Name       string
		IsExported bool
		IsMain     bool
		Body       *BlockStmt
	}

	ExportDecl struct {
	}
)

func (d FuncDecl) NodeVal() string { return d.Name }
func (d FuncDecl) declNode()       {}

type Module struct {
	Name    string
	Decls   []Decl
	Imports []string
	Exports []string
}

type (
	BasicLit struct {
		Lit token.Token
	}

	BinaryExpr struct {
		Rhs Expr
		Op  token.Tok
		Lhs Expr
	}

	UnaryExpr struct {
		X  Expr
		Op token.Tok
	}

	GroupExpr struct {
		Lparen token.Position
		X      Expr
		Rparen token.Position
	}
)

func (e BasicLit) exprNode()             {}
func (e BasicLit) Start() token.Position { return e.Lit.Pos }
func (e BasicLit) End() token.Position {
	width := len(e.Lit.Value)

	return token.Position{
		Line:   e.Lit.Pos.Line,
		Column: e.Lit.Pos.Column + width,
		Offset: e.Lit.Pos.Offset + width,
	}
}

func (e BinaryExpr) exprNode()            {}
func (e UnaryExpr) exprNode()             {}
func (e GroupExpr) exprNode()             {}
func (e GroupExpr) Start() token.Position { return e.Lparen }
func (e GroupExpr) End() token.Position   { return e.Rparen }
