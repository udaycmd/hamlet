// Copyright 2025 Uday Tiwari. All rights reserved.
// Use of this source code is governed by MIT
// license that can be found in the LICENSE file.

package parser

import "github.com/udaycmd/hamlet/src/frontend"

type Node interface {
	Start() frontend.Position
	End() frontend.Position
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
	BinaryExpr struct {
		Rhs Expr
		Op  frontend.Tok
		Lhs Expr
	}

	UnaryExpr struct {
		X  Expr
		Op frontend.Tok
	}

	GroupExpr struct {
		lparen frontend.Position
		X      Expr
		rparen frontend.Position
	}
)

func (e BinaryExpr) exprNode()               {}
func (e UnaryExpr) exprNode()                {}
func (e GroupExpr) exprNode()                {}
func (e GroupExpr) Start() frontend.Position { return e.lparen }
func (e GroupExpr) End() frontend.Position   { return e.rparen }
