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
)

type (
	File struct {
		SrcFile    *token.SourceHandle
		Statements []Stmt
	}

	IdentList struct {
		LParen  token.Position
		VarArgs bool
		List    []*Ident
		RParen  token.Position
	}
)

func (f *File) Start() token.Position { return token.Position(f.SrcFile.Base) }
func (f *File) End() token.Position   { return token.Position(f.SrcFile.Base + f.SrcFile.Len) }

func (il *IdentList) Start() token.Position { return il.LParen }
func (il *IdentList) End() token.Position   { return il.RParen + 1 }

// --- Statements ---

type (
	BadStmt struct {
		From token.Position
		To   token.Position
	}

	EmptyStmt struct {
		Semicolon  token.Position
		IsImplicit bool
	}

	BlockStmt struct {
		Stmts  []Stmt
		LBrace token.Position
		RBrace token.Position
	}
)

func (s *BadStmt) stmtNode()             {}
func (s *BadStmt) Start() token.Position { return s.From }
func (s *BadStmt) End() token.Position   { return s.To }

func (s *EmptyStmt) stmtNode()             {}
func (s *EmptyStmt) Start() token.Position { return s.Semicolon }
func (s *EmptyStmt) End() token.Position {
	if s.IsImplicit {
		return s.Semicolon
	}

	return s.Semicolon + 1
}

func (s *BlockStmt) stmtNode()             {}
func (s *BlockStmt) Start() token.Position { return s.LBrace }
func (s *BlockStmt) End() token.Position   { return s.RBrace + 1 }

// --- Expressions ---

type (
	BadExpr struct {
		From token.Position
		To   token.Position
	}

	ImportExpr struct {
		ModuleName string
		Pos        token.Position
	}

	FuncMeta struct {
		FnName *Ident
		Params *IdentList
	}

	FuncDecl struct {
		Meta *FuncMeta
		Body *BlockStmt
	}

	Ident struct {
		Name string
		Pos  token.Position
	}
)

func (e *BadExpr) exprNode()             {}
func (e *BadExpr) Start() token.Position { return e.From }
func (e *BadExpr) End() token.Position   { return e.To }

func (e *ImportExpr) exprNode()             {}
func (e *ImportExpr) Start() token.Position { return e.Pos }
func (e *ImportExpr) End() token.Position   { return token.Position(int(e.Pos) + len(e.ModuleName) + 10) }
func (e *ImportExpr) String() string        { return `import("` + e.ModuleName + `")` }

func (e *FuncMeta) exprNode()             {}
func (e *FuncMeta) Start() token.Position { return e.FnName.Pos }
func (e *FuncMeta) End() token.Position   { return e.Params.End() }

func (e *FuncDecl) exprNode()             {}
func (e *FuncDecl) Start() token.Position { return e.Meta.Start() }
func (e *FuncDecl) End() token.Position   { return e.Body.End() }

func (e *Ident) exprNode()             {}
func (e *Ident) Start() token.Position { return e.Pos }
func (e *Ident) End() token.Position   { return token.Position(int(e.Pos) + len(e.Name)) }
