// Copyright 2026 Uday Tiwari. All rights reserved.
// Use of this source code is governed by MIT
// license that can be found in the LICENSE file.

package blackbird

import (
	"fmt"

	"github.com/udaycmd/hamlet/src/frontend/parser"
	"github.com/udaycmd/hamlet/src/frontend/token"
)

type CompileError struct {
	Sm   *token.SourceManager
	Node parser.Node
	Err  error
}

func (e *CompileError) Error() string {
	return fmt.Sprintf("CompileError at %s: %s", e.Sm.SrcPos(e.Node.Start()), e.Err.Error())
}

type Compiler struct {
	file   *token.SourceHandle
	parent *Compiler
}
