package blackbird

import (
	"testing"
)

func TestDefine(t *testing.T) {
	expected := map[string]*Symbol{
		"a": {Name: "a", ScopeID: ScopeGlobal, Index: 0},
		"b": {Name: "b", ScopeID: ScopeGlobal, Index: 1},
		"c": {Name: "c", ScopeID: ScopeLocal, Index: 0},
		"d": {Name: "d", ScopeID: ScopeLocal, Index: 1},
		"e": {Name: "e", ScopeID: ScopeLocal, Index: 0},
		"f": {Name: "f", ScopeID: ScopeLocal, Index: 1},
	}

	global := NewSymtab()

	a := global.Insert("a")
	if a.Name != expected["a"].Name || a.ScopeID != expected["a"].ScopeID || a.Index != expected["a"].Index {
		t.Errorf("expected a=%+v, got=%+v", expected["a"], a)
	}

	b := global.Insert("b")
	if b.Name != expected["b"].Name || b.ScopeID != expected["b"].ScopeID || b.Index != expected["b"].Index {
		t.Errorf("expected b=%+v, got=%+v", expected["b"], b)
	}

	firstLocal := global.InitScope(false)
	c := firstLocal.Insert("c")
	if c.Name != expected["c"].Name || c.ScopeID != expected["c"].ScopeID || c.Index != expected["c"].Index {
		t.Errorf("expected c=%+v, got=%+v", expected["c"], c)
	}

	d := firstLocal.Insert("d")
	if d.Name != expected["d"].Name || d.ScopeID != expected["d"].ScopeID || d.Index != expected["d"].Index {
		t.Errorf("expected d=%+v, got=%+v", expected["d"], d)
	}

	secondLocal := firstLocal.InitScope(false)
	e := secondLocal.Insert("e")
	if e.Name != expected["e"].Name || e.ScopeID != expected["e"].ScopeID || e.Index != expected["e"].Index {
		t.Errorf("expected e=%+v, got=%+v", expected["e"], e)
	}

	f := secondLocal.Insert("f")
	if f.Name != expected["f"].Name || f.ScopeID != expected["f"].ScopeID || f.Index != expected["f"].Index {
		t.Errorf("expected f=%+v, got=%+v", expected["f"], f)
	}
}

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

	for _, tt := range expected {
		sym, _, ok := global.LookUp(tt.name, false)
		if !ok {
			t.Errorf("name %s not resolvable", tt.name)
			continue
		}
		if sym.Name != tt.sym.Name || sym.ScopeID != tt.sym.ScopeID || sym.Index != tt.sym.Index {
			t.Errorf("expected %s to resolve to %+v, got %+v", tt.name, tt.sym, sym)
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

	for _, tt := range expected {
		sym, _, ok := local.LookUp(tt.name, false)
		if !ok {
			t.Errorf("name %s not resolvable", tt.name)
			continue
		}
		if sym.Name != tt.sym.Name || sym.ScopeID != tt.sym.ScopeID || sym.Index != tt.sym.Index {
			t.Errorf("expected %s to resolve to %+v, got %+v", tt.name, tt.sym, sym)
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

	for _, tt := range tests {
		for _, sym := range tt.expectedSymbols {
			result, _, ok := tt.table.LookUp(sym.name, false)
			if !ok {
				t.Errorf("name %s not resolvable", sym.name)
				continue
			}
			if result.Name != sym.sym.Name || result.ScopeID != sym.sym.ScopeID || result.Index != sym.sym.Index {
				t.Errorf("expected %s to resolve to %+v, got %+v", sym.name, sym.sym, result)
			}
		}
	}
}

func TestResolveBuiltins(t *testing.T) {
	global := NewSymtab()
	global.AddBuiltin(0, "len")
	global.AddBuiltin(1, "first")

	local := global.InitScope(false)
	local.Insert("a")

	nestedLocal := local.InitScope(false)

	expected := []struct {
		name string
		sym  *Symbol
	}{
		{"len", &Symbol{Name: "len", ScopeID: ScopeBuiltin, Index: 0}},
		{"first", &Symbol{Name: "first", ScopeID: ScopeBuiltin, Index: 1}},
	}

	for _, tt := range expected {
		sym, _, ok := global.LookUp(tt.name, false)
		if !ok {
			t.Errorf("name %s not resolvable", tt.name)
			continue
		}
		if sym.Name != tt.sym.Name || sym.ScopeID != tt.sym.ScopeID || sym.Index != tt.sym.Index {
			t.Errorf("expected %s to resolve to %+v, got %+v", tt.name, tt.sym, sym)
		}

		sym, _, ok = local.LookUp(tt.name, false)
		if !ok {
			t.Errorf("name %s not resolvable", tt.name)
			continue
		}
		if sym.Name != tt.sym.Name || sym.ScopeID != tt.sym.ScopeID || sym.Index != tt.sym.Index {
			t.Errorf("expected %s to resolve to %+v, got %+v", tt.name, tt.sym, sym)
		}

		sym, _, ok = nestedLocal.LookUp(tt.name, false)
		if !ok {
			t.Errorf("name %s not resolvable", tt.name)
			continue
		}
		if sym.Name != tt.sym.Name || sym.ScopeID != tt.sym.ScopeID || sym.Index != tt.sym.Index {
			t.Errorf("expected %s to resolve to %+v, got %+v", tt.name, tt.sym, sym)
		}
	}
}

func TestResolveUnresolvableFree(t *testing.T) {
	global := NewSymtab()
	global.Insert("a")

	local := global.InitScope(false)
	local.Insert("c")

	_, _, ok := local.LookUp("b", false)
	if ok {
		t.Errorf("expected 'b' to not be resolvable")
	}
}

func TestBlocks(t *testing.T) {
	global := NewSymtab()
	global.Insert("a")

	block1 := global.InitScope(true)
	b := block1.Insert("b")
	if b.Index != 1 {
		t.Errorf("expected b index to be 1, got %d", b.Index)
	}

	block2 := block1.InitScope(true)
	c := block2.Insert("c")
	if c.Index != 2 {
		t.Errorf("expected c index to be 2, got %d", c.Index)
	}

	sym, _, ok := block2.LookUp("b", false)
	if !ok || sym.Name != "b" || sym.ScopeID != ScopeLocal || sym.Index != 1 {
		t.Errorf("expected 'b' to resolve to local in block scope, got ok: %v, sym: %+v", ok, sym)
	}
}
