// Copyright 2026 Uday Tiwari. All rights reserved.
// Use of this source code is governed by MIT
// license that can be found in the LICENSE file.

package blackbird

import (
	"fmt"
	"io"
	"reflect"
	"strings"

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
	symtab *SymbolTable

	// tracing
	tracing     bool
	traceW      io.Writer // trace output
	traceIndent int       // trace indentation
}

func NewCompiler(file *token.SourceHandle, table *SymbolTable, traceW io.Writer) *Compiler {
	if table == nil {
		table = NewSymtab()
	}

	return &Compiler{
		file:    file,
		symtab:  table,
		tracing: traceW != nil,
		traceW:  traceW,
	}
}

func (c *Compiler) tracePrint(stringer ...any) {
	fmt.Fprint(c.traceW, strings.Repeat(" ", 2*c.traceIndent))
	fmt.Fprintln(c.traceW, stringer...)
}

func trace(c *Compiler, msg string) *Compiler {
	c.tracePrint(msg)
	c.traceIndent++
	return c
}

func untrace(c *Compiler) {
	c.traceIndent--
}

func (c *Compiler) errorf(node parser.Node, format string, args ...any) error {
	return &CompileError{
		Sm:   c.file.Manager,
		Node: node,
		Err:  fmt.Errorf(format, args...),
	}
}

func (c *Compiler) Compile(node parser.Node) error {
	if c.tracing {
		if node != nil {
			defer untrace(trace(c, fmt.Sprintf("compile{%s}", reflect.TypeOf(node).Elem().Name())))
		} else {
			defer untrace(trace(c, "nil"))
		}
	}

	switch node := node.(type) {
	case *parser.File:
		for i := range node.Statements {
			if err := c.Compile(node.Statements[i]); err != nil {
				return err
			}
		}
	case *parser.ExprStmt:
		if err := c.Compile(node.Get()); err != nil {
			return err
		}
	case *parser.BlockStmt:
		if len(node.Stmts) == 0 {
			return nil
		}

		c.symtab = c.symtab.InitScope(true)
		defer func() {
			c.symtab = c.symtab.ParentScope(false)
		}()

		for i := range node.Stmts {
			if err := c.Compile(node.Stmts[i]); err != nil {
				return err
			}
		}
	case *parser.ProcStmt:
		for i := range node.Params.List {
			c.symtab.Insert(node.Params.List[i].Name)
		}

		if err := c.Compile(node.Body); err != nil {
			return err
		}
	case *parser.ReturnStmt:
		if c.symtab.ParentScope(true) == nil {
			return c.errorf(node, "return statement outside a procedure block")
		}
	}

	return nil
}
