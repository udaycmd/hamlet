// Copyright 2025 Uday Tiwari. All rights reserved.
// Use of this source code is governed by MIT
// license that can be found in the LICENSE file.

package parser

import (
	"testing"

	"github.com/udaycmd/hamlet/src/lexer"
	"github.com/udaycmd/hamlet/src/token"
)

var testSrcManager = token.NewSourceManager()

func TestFuncDecl(t *testing.T) {
	input := `
proc main(a, b, c) {
}
`
	testFile := testSrcManager.AddFile("testfile", -1, len(input))
	p := NewParser(testFile, []byte(input), 10, lexer.NoAsi, nil)

	f, err := p.Parse()
	if err != nil {
		t.Fatalf("parser error: %v", err)
	}

	if len(f.Statements) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(f.Statements))
	}

	procStmt, ok := f.Statements[0].(*ProcStmt)
	if !ok {
		t.Fatalf("expected ProcStmt, got %T", f.Statements[0])
	}

	if procStmt.ProcName.Name != "main" {
		t.Errorf("expected function name 'main', got '%s'", procStmt.ProcName.Name)
	}

	if len(procStmt.Params.List) != 3 {
		t.Errorf("expected 3 parameters, got %d", len(procStmt.Params.List))
	}
}
