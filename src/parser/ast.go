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

	ProcStmt struct {
		Proc     token.Position
		ProcName *Ident
		Params   *IdentList
		Body     *BlockStmt
	}

	Ident struct {
		Name string
		Pos  token.Position
	}

	BranchStmt struct {
		Kind  token.Tok
		Pos   token.Position
		Label *Ident
	}

	ReturnStmt struct {
		Ret token.Position
		e   Expr
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

func (s *ProcStmt) stmtNode()             {}
func (s *ProcStmt) Start() token.Position { return s.Proc }
func (s *ProcStmt) End() token.Position   { return s.Body.End() }

func (s *Ident) stmtNode()             {}
func (s *Ident) Start() token.Position { return s.Pos }
func (s *Ident) End() token.Position   { return token.Position(int(s.Pos) + len(s.Name)) }

func (s *BranchStmt) stmtNode()             {}
func (s *BranchStmt) Start() token.Position { return s.Pos }
func (s *BranchStmt) End() token.Position {
	if s.Label != nil {
		return s.Label.End()
	}

	return token.Position(int(s.Pos) + len(s.Kind.String()))
}

func (s *ReturnStmt) stmtNode()             {}
func (s *ReturnStmt) Start() token.Position { return s.Ret }
func (s *ReturnStmt) End() token.Position   { return s.e.End() }

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

	CallExpr struct {
		Proc     Expr
		LParen   token.Position
		Args     []Expr
		Ellipsis token.Position
		RParen   token.Position
	}

	StringLit struct {
		Val     string
		Literal string
		Pos     token.Position
	}
)

func (e *BadExpr) exprNode()             {}
func (e *BadExpr) Start() token.Position { return e.From }
func (e *BadExpr) End() token.Position   { return e.To }

func (e *ImportExpr) exprNode()             {}
func (e *ImportExpr) Start() token.Position { return e.Pos }
func (e *ImportExpr) End() token.Position   { return token.Position(int(e.Pos) + len(e.ModuleName) + 10) }
func (e *ImportExpr) String() string        { return `import("` + e.ModuleName + `")` }

func (e *CallExpr) exprNode()             {}
func (e *CallExpr) Start() token.Position { return e.Proc.Start() }
func (e *CallExpr) End() token.Position   { return e.RParen + 1 }

func (e *StringLit) exprNode()             {}
func (e *StringLit) Start() token.Position { return e.Pos }
func (e *StringLit) End() token.Position   { return token.Position(int(e.Pos) + len(e.Literal)) }
