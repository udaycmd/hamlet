// Copyright 2026 Uday Tiwari. All rights reserved.
// Use of this source code is governed by MIT
// license that can be found in the LICENSE file.

package parser

import (
	"strings"

	"github.com/udaycmd/hamlet/src/token"
)

type Stringer interface {
	String() string
}

func parenthesize(op token.Tok, expr ...Expr) string {
	sb := strings.Builder{}

	sb.WriteByte('(')
	sb.WriteString(op.String())

	for _ = range expr {
		sb.WriteByte(' ')
		sb.WriteString("TODO")
	}

	sb.WriteByte(')')
	return sb.String()
}

func (expr BinaryExpr) String() string {
	return parenthesize(expr.Op, expr.Lhs, expr.Rhs)
}
