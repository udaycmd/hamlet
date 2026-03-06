// Copyright 2026 Uday Tiwari. All rights reserved.
// Use of this source code is governed by MIT
// license that can be found in the LICENSE file.

package blackbird

type Scope uint8

const (
	ScopeBuiltin Scope = iota + 1
	ScopeGlobal
	ScopeLocal
	ScopeUpValue
)

type Symbol struct {
	Name    string
	ScopeID Scope
	Index   int
}

// SymbolTable implements a scoped symbol table.
type SymbolTable struct {
	surrounding *SymbolTable       // surrounding scope's symbol table
	table       map[string]*Symbol // actual store
	maxDecls    int                // maximum number of declaration at current scope
	numDecls    int                // number of declaration at current scope
	upvalues    []*Symbol          // upvalues in current scope
	builtins    []*Symbol          // builtins in current scope
	block       bool               // is current scope is a non-functional block
}

func NewSymtab() *SymbolTable {
	return &SymbolTable{
		table: make(map[string]*Symbol),
	}
}

func (st *SymbolTable) nextIndex() int {
	// on top level st.block can't be 'true'
	if st.block {
		return st.surrounding.nextIndex() + st.numDecls
	}

	return st.numDecls
}

func (st *SymbolTable) updateMaxDecls(n int) {
	st.maxDecls = max(n, st.maxDecls)
	if st.block {
		st.surrounding.updateMaxDecls(n)
	}
}

// add a new symbol in the current scope
func (st *SymbolTable) Insert(name string) *Symbol {
	s := &Symbol{Name: name, Index: st.nextIndex()}
	st.numDecls++

	if st.surrounding == nil {
		s.ScopeID = ScopeGlobal
	} else {
		s.ScopeID = ScopeLocal
	}

	st.table[name] = s
	st.updateMaxDecls(s.Index + 1)
	return s
}

// LookUp resolves a symbol with a given name
func (st *SymbolTable) LookUp(name string) (*Symbol, int, bool) {
	s, ok := st.table[name]
	if ok {
		return s, 0, true
	}

	if st.surrounding == nil {
		return nil, 0, false
	}

	s, depth, ok := st.surrounding.LookUp(name)
	if !ok {
		return nil, 0, false
	}
	depth++

	// defined in parent table and if it's not in global/builtin scope
	// then its a upvalue
	if !st.block && s.ScopeID != ScopeGlobal && s.ScopeID != ScopeBuiltin {
		return st.UpValue(s), depth, true
	}

	return s, depth, true
}

// InitScope creates a new symbol table for a new scope
func (st *SymbolTable) InitScope(block bool) *SymbolTable {
	return &SymbolTable{
		table:       make(map[string]*Symbol),
		surrounding: st,
		block:       block,
	}
}

// ParentScope return the SymbolTable of parent scope of the current scope
func (st *SymbolTable) ParentScope(skipBlock bool) *SymbolTable {
	if skipBlock && st.block {
		return st.surrounding.ParentScope(skipBlock)
	}

	return st.surrounding
}

func (st *SymbolTable) UpValue(s0 *Symbol) *Symbol {
	st.upvalues = append(st.upvalues, s0)
	s1 := &Symbol{
		Name:    s0.Name,
		Index:   len(st.upvalues) - 1,
		ScopeID: ScopeUpValue,
	}

	st.table[s0.Name] = s1
	return s1
}

func (st *SymbolTable) AddBuiltin(index int, name string) *Symbol {
	if st.surrounding != nil {
		return st.surrounding.AddBuiltin(index, name)
	}

	s := &Symbol{
		Name:    name,
		Index:   index,
		ScopeID: ScopeBuiltin,
	}

	st.table[name] = s
	st.builtins = append(st.builtins, s)
	return s
}
