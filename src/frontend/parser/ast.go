// Copyright 2025 Uday Tiwari. All rights reserved.
// Use of this source code is governed by MIT
// license that can be found in the LICENSE file.

package parser

import "github.com/udaycmd/hamlet/src/frontend/token"

// --- AST Nodes ---

type Node interface {
	Start() token.Position
	End() token.Position
}

type (
	// Expression is an ast node in hamlet
	// that is evaluated to determine its value
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

	ExprStmt struct {
		e Expr
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

	DeclStmt struct {
		Decl  token.Position
		Ident *Ident
		Equal token.Position
		Val   Expr
	}

	ConstStmt struct {
		Const token.Position
		Ident *Ident
		Equal token.Position
		Val   Expr
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

	ExportStmt struct {
		Pos token.Position
		e   Expr
	}

	IfStmt struct {
		IfPos token.Position
		Cond  Expr
		Body  *BlockStmt
		Else  Stmt
	}

	ForStmt struct {
		ForPos token.Position
		Cond   Expr
		Body   *BlockStmt
	}

	ForInStmt struct {
		ForPos   token.Position
		Index    *Ident
		Val      *Ident
		Iterable Expr
		Body     *BlockStmt
	}

	AssignStmt struct {
		Lhs    []Expr
		EqType token.Tok
		Eq     token.Position
		Rhs    []Expr
	}
)

func (s *BadStmt) stmtNode()             {}
func (s *BadStmt) Start() token.Position { return s.From }
func (s *BadStmt) End() token.Position   { return s.To }

func (s *ExprStmt) stmtNode()             {}
func (s *ExprStmt) Start() token.Position { return s.Start() }
func (s *ExprStmt) End() token.Position   { return s.End() }

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

func (s *DeclStmt) stmtNode()             {}
func (s *DeclStmt) Start() token.Position { return s.Decl }
func (s *DeclStmt) End() token.Position   { return s.Val.End() }

func (s *ConstStmt) stmtNode()             {}
func (s *ConstStmt) Start() token.Position { return s.Const }
func (s *ConstStmt) End() token.Position   { return s.Val.End() }

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

func (s *ExportStmt) stmtNode()             {}
func (s *ExportStmt) Start() token.Position { return s.Pos }
func (s *ExportStmt) End() token.Position   { return s.e.End() }

func (s *IfStmt) stmtNode()             {}
func (s *IfStmt) Start() token.Position { return s.IfPos }
func (s *IfStmt) End() token.Position {
	if s.Else != nil {
		return s.Else.End()
	}

	return s.Body.End()
}

func (s *ForStmt) stmtNode()             {}
func (s *ForStmt) Start() token.Position { return s.ForPos }
func (s *ForStmt) End() token.Position   { return s.Body.End() }

func (s *ForInStmt) stmtNode()             {}
func (s *ForInStmt) Start() token.Position { return s.ForPos }
func (s *ForInStmt) End() token.Position   { return s.Body.End() }

func (s *AssignStmt) stmtNode()             {}
func (s *AssignStmt) Start() token.Position { return s.Lhs[0].Start() }
func (s *AssignStmt) End() token.Position   { return s.Rhs[len(s.Rhs)-1].End() }

// --- Expressions ---

type (
	BadExpr struct {
		From token.Position
		To   token.Position
	}

	Ident struct {
		Name string
		Pos  token.Position
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

	EmptyLit struct {
		Pos token.Position
	}

	IntLit struct {
		Val     int64
		Literal string
		Pos     token.Position
	}

	FloatLit struct {
		Val     float64
		Literal string
		Pos     token.Position
	}

	CharLit struct {
		Val     rune
		Literal string
		Pos     token.Position
	}

	StringLit struct {
		Val     string
		Literal string
		Pos     token.Position
	}

	BoolLit struct {
		Val     bool
		Literal string
		Pos     token.Position
	}

	ProcLit struct {
		Proc   token.Position
		Params *IdentList
		Body   *BlockStmt
	}

	KvLit struct {
		KeyPos token.Position
		Key    string
		Colon  token.Position
		Value  Expr
	}

	MapLit struct {
		LBrace token.Position
		Kvs    []*KvLit
		RBrace token.Position
	}

	ArrayLit struct {
		LBrack token.Position
		Items  []Expr
		RBrack token.Position
	}

	BinaryExpr struct {
		Lhs   Expr
		Rhs   Expr
		Op    token.Tok
		OpPos token.Position
	}

	UnaryExpr struct {
		Op    token.Tok
		OpPos token.Position
		X     Expr
	}

	GroupedExpr struct {
		Lparen token.Position
		X      Expr
		Rparen token.Position
	}

	IndexExpr struct {
		X     Expr
		LBrac token.Position
		Index Expr
		Rbrac token.Position
	}

	SliceExpr struct {
		X     Expr
		LBrac token.Position
		Lo    Expr
		Hi    Expr
		Rbrac token.Position
	}

	ReceiverExpr struct {
		X  Expr
		Id *Ident
	}

	TernaryExpr struct {
		X     Expr
		Ques  token.Position
		True  Expr
		Colon token.Position
		False Expr
	}

	RangeExpr struct {
		Lhs    Expr
		Ranger token.Position
		Rhs    Expr
	}
)

func (e *BadExpr) exprNode()             {}
func (e *BadExpr) Start() token.Position { return e.From }
func (e *BadExpr) End() token.Position   { return e.To }

func (s *Ident) exprNode()             {}
func (s *Ident) Start() token.Position { return s.Pos }
func (s *Ident) End() token.Position   { return token.Position(int(s.Pos) + len(s.Name)) }

func (e *ImportExpr) exprNode()             {}
func (e *ImportExpr) Start() token.Position { return e.Pos }
func (e *ImportExpr) End() token.Position   { return token.Position(int(e.Pos) + len(e.ModuleName) + 10) }
func (e *ImportExpr) String() string        { return `import("` + e.ModuleName + `")` }

func (e *CallExpr) exprNode()             {}
func (e *CallExpr) Start() token.Position { return e.Proc.Start() }
func (e *CallExpr) End() token.Position   { return e.RParen + 1 }

func (e *EmptyLit) exprNode()             {}
func (e *EmptyLit) Start() token.Position { return e.Pos }
func (e *EmptyLit) End() token.Position   { return e.Pos + token.Position(5) }

func (e *IntLit) exprNode()             {}
func (e *IntLit) Start() token.Position { return e.Pos }
func (e *IntLit) End() token.Position   { return e.Pos + token.Position(len(e.Literal)) }

func (e *FloatLit) exprNode()             {}
func (e *FloatLit) Start() token.Position { return e.Pos }
func (e *FloatLit) End() token.Position   { return e.Pos + token.Position(len(e.Literal)) }

func (e *CharLit) exprNode()             {}
func (e *CharLit) Start() token.Position { return e.Pos }
func (e *CharLit) End() token.Position   { return e.Pos + token.Position(len(e.Literal)) }

func (e *StringLit) exprNode()             {}
func (e *StringLit) Start() token.Position { return e.Pos }
func (e *StringLit) End() token.Position   { return e.Pos + token.Position(len(e.Literal)) }

func (e *BoolLit) exprNode()             {}
func (e *BoolLit) Start() token.Position { return e.Pos }
func (e *BoolLit) End() token.Position   { return e.Pos + token.Position(len(e.Literal)) }

func (e *ProcLit) exprNode()             {}
func (e *ProcLit) Start() token.Position { return e.Proc }
func (e *ProcLit) End() token.Position   { return e.Proc + e.Body.End() }

func (e *KvLit) exprNode()             {}
func (e *KvLit) Start() token.Position { return e.KeyPos }
func (e *KvLit) End() token.Position   { return e.Value.End() }

func (e *MapLit) exprNode()             {}
func (e *MapLit) Start() token.Position { return e.LBrace }
func (e *MapLit) End() token.Position   { return e.RBrace + 1 }

func (e *ArrayLit) exprNode()             {}
func (e *ArrayLit) Start() token.Position { return e.LBrack }
func (e *ArrayLit) End() token.Position   { return e.RBrack + 1 }

func (e *BinaryExpr) exprNode()             {}
func (e *BinaryExpr) Start() token.Position { return e.Lhs.Start() }
func (e *BinaryExpr) End() token.Position   { return e.Rhs.End() }

func (e *UnaryExpr) exprNode()             {}
func (e *UnaryExpr) Start() token.Position { return e.OpPos }
func (e *UnaryExpr) End() token.Position   { return e.X.End() }

func (e *GroupedExpr) exprNode()             {}
func (e *GroupedExpr) Start() token.Position { return e.Lparen }
func (e *GroupedExpr) End() token.Position   { return e.Rparen }

func (e *IndexExpr) exprNode()             {}
func (e *IndexExpr) Start() token.Position { return e.X.Start() }
func (e *IndexExpr) End() token.Position   { return e.Rbrac }

func (e *SliceExpr) exprNode()             {}
func (e *SliceExpr) Start() token.Position { return e.X.Start() }
func (e *SliceExpr) End() token.Position   { return e.Rbrac }

func (e *ReceiverExpr) exprNode()             {}
func (e *ReceiverExpr) Start() token.Position { return e.X.Start() }
func (e *ReceiverExpr) End() token.Position   { return e.Id.Pos + token.Position(len(e.Id.Name)) }

func (e *TernaryExpr) exprNode()             {}
func (e *TernaryExpr) Start() token.Position { return e.X.Start() }
func (e *TernaryExpr) End() token.Position   { return e.False.End() }

func (e *RangeExpr) exprNode()             {}
func (e *RangeExpr) Start() token.Position { return e.Lhs.Start() }
func (e *RangeExpr) End() token.Position   { return e.Rhs.End() }
