// Copyright 2026 Uday Tiwari. All rights reserved.
// Use of this source code is governed by MIT
// license that can be found in the LICENSE file.

package token

import "fmt"

type Position struct {
	Line   int
	Column int
	Offset int
}

var (
	InvalidPos = Position{Line: 0, Column: 0, Offset: 0}
)

func (p Position) String() string {
	return fmt.Sprintf("%d:%d", p.Line, p.Column) // line:column
}
