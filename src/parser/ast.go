// Copyright 2025 Uday Tiwari. All rights reserved.
// Use of this source code is governed by MIT
// license that can be found in the LICENSE file.

package parser

import "github.com/udaycmd/hamlet/src/lexer"

type Node interface {
	Start() lexer.Position
	End() lexer.Position
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
		Op  lexer.Tok
		Lhs Expr
	}

	UnaryExpr struct {
		X  Expr
		Op lexer.Tok
	}

	GroupExpr struct {
		lparen lexer.Position
		X      Expr
		rparen lexer.Position
	}
)

func (e BinaryExpr) exprNode()               {}
func (e UnaryExpr) exprNode()                {}
func (e GroupExpr) exprNode()                {}
func (e GroupExpr) Start() lexer.Position { return e.lparen }
func (e GroupExpr) End() lexer.Position   { return e.rparen }
