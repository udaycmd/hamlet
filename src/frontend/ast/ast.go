// Copyright 2025 Uday Tiwari. All rights reserved.
// Use of this source code is governed by MIT
// license that can be found in the LICENSE file.

package ast

type Node interface {
	NodeVal() string
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
