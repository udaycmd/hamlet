package blackbird

import (
	"testing"

	"github.com/udaycmd/hamlet/src/utils"
)

func TestResolveGlobal(t *testing.T) {
	global := NewSymtab()
	global.Insert("a")
	global.Insert("b")

	expected := []struct {
		name string
		sym  *Symbol
	}{
		{"a", &Symbol{Name: "a", ScopeID: ScopeGlobal, Index: 0}},
		{"b", &Symbol{Name: "b", ScopeID: ScopeGlobal, Index: 1}},
	}

	for _, tc := range expected {
		sym, _, ok := global.LookUp(tc.name)
		if !ok {
			utils.Fail(t, "name %s not resolvable in global scope", tc.name)
		}
		if sym.Name != tc.sym.Name || sym.ScopeID != tc.sym.ScopeID || sym.Index != tc.sym.Index {
			utils.Fail(t, "expected %s to resolve to %+v, got %+v", tc.name, tc.sym, sym)
		}
	}
}

func TestResolveLocal(t *testing.T) {
	global := NewSymtab()
	global.Insert("a")
	global.Insert("b")

	local := global.InitScope(false)
	local.Insert("c")
	local.Insert("d")

	expected := []struct {
		name string
		sym  *Symbol
	}{
		{"a", &Symbol{Name: "a", ScopeID: ScopeGlobal, Index: 0}},
		{"b", &Symbol{Name: "b", ScopeID: ScopeGlobal, Index: 1}},
		{"c", &Symbol{Name: "c", ScopeID: ScopeLocal, Index: 0}},
		{"d", &Symbol{Name: "d", ScopeID: ScopeLocal, Index: 1}},
	}

	for _, tc := range expected {
		sym, _, ok := local.LookUp(tc.name)
		if !ok {
			utils.Fail(t, "name %s not resolvable in local scope", tc.name)
		}
		if sym.Name != tc.sym.Name || sym.ScopeID != tc.sym.ScopeID || sym.Index != tc.sym.Index {
			utils.Fail(t, "expected %s to resolve to %+v, got %+v", tc.name, tc.sym, sym)
		}
	}
}

func TestResolveNestedLocal(t *testing.T) {
	global := NewSymtab()
	global.Insert("a")
	global.Insert("b")

	firstLocal := global.InitScope(false)
	firstLocal.Insert("c")
	firstLocal.Insert("d")

	secondLocal := firstLocal.InitScope(false)
	secondLocal.Insert("e")
	secondLocal.Insert("f")

	tests := []struct {
		table           *SymbolTable
		expectedSymbols []struct {
			name string
			sym  *Symbol
		}
	}{
		{
			firstLocal,
			[]struct {
				name string
				sym  *Symbol
			}{
				{"a", &Symbol{Name: "a", ScopeID: ScopeGlobal, Index: 0}},
				{"b", &Symbol{Name: "b", ScopeID: ScopeGlobal, Index: 1}},
				{"c", &Symbol{Name: "c", ScopeID: ScopeLocal, Index: 0}},
				{"d", &Symbol{Name: "d", ScopeID: ScopeLocal, Index: 1}},
			},
		},
		{
			secondLocal,
			[]struct {
				name string
				sym  *Symbol
			}{
				{"a", &Symbol{Name: "a", ScopeID: ScopeGlobal, Index: 0}},
				{"b", &Symbol{Name: "b", ScopeID: ScopeGlobal, Index: 1}},
				{"c", &Symbol{Name: "c", ScopeID: ScopeUpValue, Index: 0}},
				{"d", &Symbol{Name: "d", ScopeID: ScopeUpValue, Index: 1}},
				{"e", &Symbol{Name: "e", ScopeID: ScopeLocal, Index: 0}},
				{"f", &Symbol{Name: "f", ScopeID: ScopeLocal, Index: 1}},
			},
		},
	}

	for _, tc := range tests {
		for _, sym := range tc.expectedSymbols {
			result, _, ok := tc.table.LookUp(sym.name)
			if !ok {
				utils.Fail(t, "name %s not resolvable", sym.name)
			}
			if result.Name != sym.sym.Name || result.ScopeID != sym.sym.ScopeID || result.Index != sym.sym.Index {
				utils.Fail(t, "expected %s to resolve to %+v, got %+v", sym.name, sym.sym, result)
			}
		}
	}
}

func TestResolveBuiltins(t *testing.T) {
	global := NewSymtab()
	global.AddBuiltin(0, "builtin1")
	global.AddBuiltin(1, "builtin2")

	local := global.InitScope(false)
	local.Insert("a")

	nestedLocal := local.InitScope(false)

	expected := []struct {
		name string
		sym  *Symbol
	}{
		{"builtin1", &Symbol{Name: "builtin1", ScopeID: ScopeBuiltin, Index: 0}},
		{"builtin2", &Symbol{Name: "builtin2", ScopeID: ScopeBuiltin, Index: 1}},
	}

	for _, tc := range expected {
		// check in global scope
		sym, _, ok := global.LookUp(tc.name)
		if !ok {
			utils.Fail(t, "builtin name %s not resolvable in global scope", tc.name)
		}

		if sym.Name != tc.sym.Name || sym.ScopeID != tc.sym.ScopeID || sym.Index != tc.sym.Index {
			utils.Fail(t, "expected %s to resolve to %+v, got %+v", tc.name, tc.sym, sym)
		}

		// check in local scope
		sym, _, ok = local.LookUp(tc.name)
		if !ok {
			utils.Fail(t, "builtin name %s not resolvable in local scope", tc.name)
		}

		if sym.Name != tc.sym.Name || sym.ScopeID != tc.sym.ScopeID || sym.Index != tc.sym.Index {
			utils.Fail(t, "expected %s to resolve to %+v, got %+v", tc.name, tc.sym, sym)
		}

		sym, _, ok = nestedLocal.LookUp(tc.name)
		if !ok {
			utils.Fail(t, "builtin name %s not resolvable in nested local scope", tc.name)
		}

		if sym.Name != tc.sym.Name || sym.ScopeID != tc.sym.ScopeID || sym.Index != tc.sym.Index {
			utils.Fail(t, "expected %s to resolve to %+v, got %+v", tc.name, tc.sym, sym)
		}
	}
}

func TestResolveUnresolvableSymbol(t *testing.T) {
	global := NewSymtab()
	global.Insert("a")

	local := global.InitScope(false) // functional scope
	local.Insert("c")

	_, _, ok := local.LookUp("b")
	if ok {
		utils.Fail(t, "expected 'b' unresolvable")
	}
}

func TestBlocks(t *testing.T) {
	global := NewSymtab()
	global.Insert("a")

	block1 := global.InitScope(true)
	b := block1.Insert("b")
	if b.Index != 1 {
		utils.Fail(t, "expected b index to be 1, got %d", b.Index)
	}

	block2 := block1.InitScope(true)
	c := block2.Insert("c")
	if c.Index != 2 {
		utils.Fail(t, "expected c index to be 2, got %d", c.Index)
	}

	sym, depth, ok := block2.LookUp("b")
	if !ok || sym.Name != "b" || depth != 1 || sym.ScopeID != ScopeLocal || sym.Index != 1 {
		utils.Fail(t, "expected 'b' to resolve to local in block scope, got ok: %v, sym: %+v", ok, sym)
	}
}
