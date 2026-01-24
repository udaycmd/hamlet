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
	File struct {
		SrcFile    *token.SourceHandle
		Statements []Stmt
	}
)

func (F *File) Start() token.Position { return token.Position(F.SrcFile.Base) }
func (F *File) End() token.Position   { return token.Position(F.SrcFile.Base + F.SrcFile.Len) }
