// Copyright 2025 Uday Tiwari. All rights reserved.
// Use of this source code is governed by MIT
// license that can be found in the LICENSE file.

package parser

import (
	"reflect"
	"testing"

	"github.com/udaycmd/hamlet/src/frontend/token"
)

// --- Test Helpers ---

func fail(t *testing.T, format string, args ...any) {
	t.Helper()
	t.Errorf(format, args...)
}

func expectEqual(t *testing.T, x, y any, format string, args ...any) {
	t.Helper()
	if !reflect.DeepEqual(x, y) {
		fail(t, format, args...)
	}
}

func mustParse(t *testing.T, input string) *File {
	t.Helper()
	test := token.NewSourceManager()
	testfile := test.AddFile("testfile", -1, len(input))
	root, err := NewParser(testfile, []byte(input), 10, 0, nil).Parse()
	if err != nil {
		fail(t, "unexpected parse error: %s", err.Error())
	}

	return root
}

func expectParseError(t *testing.T, input string) error {
	t.Helper()
	test := token.NewSourceManager()
	testfile := test.AddFile("testfile", -1, len(input))
	_, err := NewParser(testfile, []byte(input), 10, 0, nil).Parse()
	return err
}

// --- ParseError Tests ---

func TestParseError(t *testing.T) {
	err := &ParseError{
		position: token.SrcPos{Line: 1, Column: 10, FileName: "test.ham"},
		msg:      "unexpected token",
	}
	expected := "test.ham(1, 10): unexpected token"
	expectEqual(t, expected, err.Error(), "ParseError.Error()")
}

func TestParseErrorNoFile(t *testing.T) {
	err := &ParseError{
		position: token.SrcPos{},
		msg:      "some error",
	}
	expected := "some error"
	expectEqual(t, expected, err.Error(), "ParseError.Error() with no file")
}

func TestParseErrors(t *testing.T) {
	var errs ParseErrors
	errs.Extend(token.SrcPos{Line: 3, Column: 5, FileName: "a.ham"}, "error 3")
	errs.Extend(token.SrcPos{Line: 1, Column: 1, FileName: "a.ham"}, "error 1")
	errs.Extend(token.SrcPos{Line: 2, Column: 3, FileName: "a.ham"}, "error 2")

	expectEqual(t, 3, errs.Len(), "ParseErrors.Len()")

	errs.Sort()
	expectEqual(t, "error 1", errs[0].msg, "sorted first error")
	expectEqual(t, "error 2", errs[1].msg, "sorted second error")
	expectEqual(t, "error 3", errs[2].msg, "sorted third error")

	if errs.Error() == "" {
		fail(t, "ParseErrors.Error() returned empty string")
	}
}

func TestParseErrorsEmpty(t *testing.T) {
	var errs ParseErrors
	expectEqual(t, "empty", errs.Error(), "empty ParseErrors.Error()")
	expectEqual(t, nil, errs.GetParseErrors(), "empty ParseErrors.err()")
}

// --- Literal Parsing Tests ---

func TestParseIntLit(t *testing.T) {
	testCases := []struct {
		input    string
		expected int64
	}{
		{"0", 0},
		{"123", 123},
		{"0xFF", 255},
		{"42", 42},
	}

	for _, tc := range testCases {
		file := mustParse(t, tc.input)
		expectEqual(t, 1, len(file.Statements), "statement count for: %s", tc.input)

		exprStmt, ok := file.Statements[0].(*ExprStmt)
		if !ok {
			fail(t, "expected ExprStmt for input: %s", tc.input)
		}

		intLit, ok := exprStmt.e.(*IntLit)
		if !ok {
			fail(t, "expected IntLit for input: %s", tc.input)
		}

		expectEqual(t, tc.expected, intLit.Val, "int value for: %s", tc.input)
	}
}

func TestParseFloatLit(t *testing.T) {
	testCases := []struct {
		input    string
		expected float64
	}{
		{"0.0", 0.0},
		{"3.14", 3.14},
		{"1e10", 1e10},
	}

	for _, tc := range testCases {
		file := mustParse(t, tc.input)
		expectEqual(t, 1, len(file.Statements), "statement count for: %s", tc.input)

		exprStmt, ok := file.Statements[0].(*ExprStmt)
		if !ok {
			fail(t, "expected ExprStmt for input: %s", tc.input)
		}

		floatLit, ok := exprStmt.e.(*FloatLit)
		if !ok {
			fail(t, "expected FloatLit for input: %s", tc.input)
		}

		expectEqual(t, tc.expected, floatLit.Val, "float value for: %s", tc.input)
	}
}

func TestParseStringLit(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{`"hello"`, "hello"},
		{`""`, ""},
	}

	for _, tc := range testCases {
		file := mustParse(t, tc.input)
		expectEqual(t, 1, len(file.Statements), "statement count for: %s", tc.input)

		exprStmt, ok := file.Statements[0].(*ExprStmt)
		if !ok {
			t.Fatalf("expected ExprStmt for input: %s", tc.input)
		}

		strLit, ok := exprStmt.e.(*StringLit)
		if !ok {
			fail(t, "expected StringLit for input: %s", tc.input)
		}

		expectEqual(t, tc.expected, strLit.Val, "string value for: %s", tc.input)
	}
}

func TestParseBoolLit(t *testing.T) {
	for _, tc := range []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
	} {
		file := mustParse(t, tc.input)
		exprStmt := file.Statements[0].(*ExprStmt)
		expectEqual(t, 1, len(file.Statements), "statement count for: %s", tc.input)

		boolLit := exprStmt.e.(*BoolLit)
		expectEqual(t, tc.expected, boolLit.Val, "bool value for: %s", tc.input)
	}
}

func TestParseCharLit(t *testing.T) {
	for _, tc := range []struct {
		input    string
		expected rune
	}{
		{"'a'", 'a'},
		{"'Z'", 'Z'},
	} {
		file := mustParse(t, tc.input)
		exprStmt := file.Statements[0].(*ExprStmt)
		expectEqual(t, 1, len(file.Statements), "statement count for: %s", tc.input)

		charLit := exprStmt.e.(*CharLit)
		expectEqual(t, tc.expected, charLit.Val, "char value for: %s", tc.input)
	}
}

func TestParseEmptyLit(t *testing.T) {
	file := mustParse(t, "empty")
	exprStmt := file.Statements[0].(*ExprStmt)
	_, ok := exprStmt.e.(*EmptyLit)
	if !ok {
		fail(t, "expected EmptyLit")
	}
}

// --- Identifier Tests ---

func TestParseIdent(t *testing.T) {
	for _, name := range []string{"foo", "bar", "_test", "camelCase"} {
		file := mustParse(t, name)
		exprStmt := file.Statements[0].(*ExprStmt)
		ident := exprStmt.e.(*Ident)
		expectEqual(t, name, ident.Name, "identifier name")
	}
}

// --- Expression Tests ---

func TestParseUnaryExpr(t *testing.T) {
	for _, tc := range []struct {
		input string
		op    token.Tok
	}{
		{"-5", token.MINUS},
		{"+5", token.PLUS},
		{"!true", token.BANG},
		{"^42", token.XOR},
	} {
		file := mustParse(t, tc.input)
		exprStmt := file.Statements[0].(*ExprStmt)
		expectEqual(t, 1, len(file.Statements), "statement count for: %s", tc.input)

		unaryExpr := exprStmt.e.(*UnaryExpr)
		expectEqual(t, tc.op, unaryExpr.Op, "unary operator for: %s", tc.input)
	}
}

func TestParseBinaryExpr(t *testing.T) {
	for _, tc := range []struct {
		input string
		op    token.Tok
	}{
		{"1 + 2", token.PLUS},
		{"1 - 2", token.MINUS},
		{"1 * 2", token.STAR},
		{"1 / 2", token.SLASH},
		{"1 == 2", token.EQUALS},
		{"1 != 2", token.BANG_EQ},
		{"1 < 2", token.LESS},
		{"1 > 2", token.GREATER},
		{"1 && 2", token.AND},
		{"1 || 2", token.OR},
	} {
		file := mustParse(t, tc.input)
		exprStmt := file.Statements[0].(*ExprStmt)
		expectEqual(t, 1, len(file.Statements), "statement count for: %s", tc.input)

		binaryExpr := exprStmt.e.(*BinaryExpr)
		expectEqual(t, tc.op, binaryExpr.Op, "binary operator for: %s", tc.input)
	}
}

func TestParseGroupedExpr(t *testing.T) {
	file := mustParse(t, "(1 + 2)")
	exprStmt := file.Statements[0].(*ExprStmt)
	grouped := exprStmt.e.(*GroupedExpr)
	_, ok := grouped.X.(*BinaryExpr)
	if !ok {
		fail(t, "expected BinaryExpr inside GroupedExpr")
	}
}

func TestParseTernaryExpr(t *testing.T) {
	file := mustParse(t, "a ? b : c")
	exprStmt := file.Statements[0].(*ExprStmt)

	ternary := exprStmt.e.(*TernaryExpr)
	cond := ternary.X.(*Ident)
	expectEqual(t, "a", cond.Name, "condition ident")

	trueExpr := ternary.True.(*Ident)
	expectEqual(t, "b", trueExpr.Name, "true branch ident")

	falseExpr := ternary.False.(*Ident)
	expectEqual(t, "c", falseExpr.Name, "false branch ident")
}

// --- Array and Map Literal Tests ---

func TestParseArrayLit(t *testing.T) {
	file := mustParse(t, "[1, 2, 3]")
	exprStmt := file.Statements[0].(*ExprStmt)
	arrayLit := exprStmt.e.(*ArrayLit)
	expectEqual(t, 3, len(arrayLit.Items), "array items count")
}

func TestParseEmptyArray(t *testing.T) {
	file := mustParse(t, "[]")
	exprStmt := file.Statements[0].(*ExprStmt)
	arrayLit := exprStmt.e.(*ArrayLit)
	expectEqual(t, 0, len(arrayLit.Items), "empty array")
}

func TestParseMapLit(t *testing.T) {
	file := mustParse(t, `{foo: 1, bar: 2}`)
	exprStmt := file.Statements[0].(*ExprStmt)
	mapLit := exprStmt.e.(*MapLit)
	expectEqual(t, 2, len(mapLit.Kvs), "map entries count")
	expectEqual(t, "foo", mapLit.Kvs[0].Key, "first key")
	expectEqual(t, "bar", mapLit.Kvs[1].Key, "second key")
}

func TestParseEmptyMap(t *testing.T) {
	file := mustParse(t, "{}")
	exprStmt := file.Statements[0].(*ExprStmt)
	mapLit := exprStmt.e.(*MapLit)
	expectEqual(t, 0, len(mapLit.Kvs), "empty map")
}

// --- Call Expression Tests ---

func TestParseCallExpr(t *testing.T) {
	file := mustParse(t, "foo(1, 2, 3)")
	exprStmt := file.Statements[0].(*ExprStmt)
	callExpr := exprStmt.e.(*CallExpr)

	ident := callExpr.Proc.(*Ident)
	expectEqual(t, "foo", ident.Name, "function name")
	expectEqual(t, 3, len(callExpr.Args), "argument count")
}

func TestParseCallExprNoArgs(t *testing.T) {
	file := mustParse(t, "foo()")
	exprStmt := file.Statements[0].(*ExprStmt)
	callExpr := exprStmt.e.(*CallExpr)
	expectEqual(t, 0, len(callExpr.Args), "no arguments")
}

// --- Index and Slice Expression Tests ---

func TestParseIndexExpr(t *testing.T) {
	file := mustParse(t, "arr[0]")
	exprStmt := file.Statements[0].(*ExprStmt)
	indexExpr := exprStmt.e.(*IndexExpr)

	ident := indexExpr.X.(*Ident)
	expectEqual(t, "arr", ident.Name, "indexed expression")
}

func TestParseSliceExpr(t *testing.T) {
	file := mustParse(t, "arr[1:3]")
	exprStmt := file.Statements[0].(*ExprStmt)
	sliceExpr := exprStmt.e.(*SliceExpr)

	ident := sliceExpr.X.(*Ident)
	expectEqual(t, "arr", ident.Name, "sliced expression")
}

func TestParseSliceExprNoOffset(t *testing.T) {
	file := mustParse(t, "arr[:]")
	exprStmt := file.Statements[0].(*ExprStmt)
	sliceExpr := exprStmt.e.(*SliceExpr)

	expectEqual(t, nil, sliceExpr.Lo, "no low slice offset")
	expectEqual(t, nil, sliceExpr.Hi, "no high slice offset")
}

// --- Receiver Expression Tests ---

func TestParseReceiverExpr(t *testing.T) {
	file := mustParse(t, "obj.method")
	exprStmt := file.Statements[0].(*ExprStmt)
	receiverExpr := exprStmt.e.(*ReceiverExpr)

	ident := receiverExpr.X.(*Ident)
	expectEqual(t, "obj", ident.Name, "receiver object")
	expectEqual(t, "method", receiverExpr.Id.Name, "method name")
}

// --- Statement Tests ---

func TestParseIfStmt(t *testing.T) {
	file := mustParse(t, "if x { y }")
	ifStmt := file.Statements[0].(*IfStmt)

	cond := ifStmt.Cond.(*Ident)
	expectEqual(t, "x", cond.Name, "condition")

	if ifStmt.Else != nil {
		fail(t, "expected no else clause")
	}
}

func TestParseIfElseStmt(t *testing.T) {
	file := mustParse(t, "if x { y } else { z }")
	ifStmt := file.Statements[0].(*IfStmt)

	if ifStmt.Else == nil {
		fail(t, "expected else clause")
	}

	_, ok := ifStmt.Else.(*BlockStmt)
	if !ok {
		fail(t, "expected BlockStmt as else clause")
	}
}

func TestParseIfElseIfStmt(t *testing.T) {
	file := mustParse(t, "if x { a } else if y { b }")
	ifStmt := file.Statements[0].(*IfStmt)

	elseIf := ifStmt.Else.(*IfStmt)
	cond := elseIf.Cond.(*Ident)
	expectEqual(t, "y", cond.Name, "else-if condition")
}

func TestParseForStmt(t *testing.T) {
	file := mustParse(t, "for { break }")
	forStmt := file.Statements[0].(*ForStmt)

	if forStmt.Cond != nil {
		fail(t, "expected no condition for infinite loop")
	}
}

func TestParseForCondStmt(t *testing.T) {
	file := mustParse(t, "for x < 10 { x }")
	forStmt := file.Statements[0].(*ForStmt)

	if forStmt.Cond == nil {
		fail(t, "expected condition in for loop")
	}

	_, ok := forStmt.Cond.(*BinaryExpr)
	if !ok {
		fail(t, "expected BinaryExpr as condition")
	}
}

func TestParseForInStmt(t *testing.T) {
	file := mustParse(t, "for x in arr { x }")
	forInStmt := file.Statements[0].(*ForInStmt)
	expectEqual(t, "x", forInStmt.Val.Name, "loop variable")
}

func TestParseForInWithIndexStmt(t *testing.T) {
	file := mustParse(t, "for i, v in arr { v }")
	forInStmt := file.Statements[0].(*ForInStmt)
	expectEqual(t, "i", forInStmt.Index.Name, "index variable")
	expectEqual(t, "v", forInStmt.Val.Name, "value variable")
}

func TestParseForInWithUntilSeqWithoutIndex(t *testing.T) {
	file := mustParse(t, "for v in 0 until 5 { v }")
	forInStmt := file.Statements[0].(*ForInStmt)
	expectEqual(t, "v", forInStmt.Val.Name, "value variable")
	if forInStmt.Index != nil {
		fail(t, "expected forIn index to be nil")
	}

	_, ok := forInStmt.Iterable.(*RangeExpr)
	if !ok {
		fail(t, "expected rangeExpr in forIn.Iterable")
	}
}

func TestParseForInWithUntilSeqWithIndex(t *testing.T) {
	file := mustParse(t, "for i, v in 0 until 5 { v }")
	forInStmt := file.Statements[0].(*ForInStmt)
	expectEqual(t, "i", forInStmt.Index.Name, "index variable")
	expectEqual(t, "v", forInStmt.Val.Name, "value variable")

	_, ok := forInStmt.Iterable.(*RangeExpr)
	if !ok {
		fail(t, "expected rangeExpr in forIn.Iterable")
	}
}

func TestParseProcStmt(t *testing.T) {
	file := mustParse(t, "proc foo(a, b) { return a + b }")
	procStmt := file.Statements[0].(*ProcStmt)

	expectEqual(t, "foo", procStmt.ProcName.Name, "proc name")
	expectEqual(t, 2, len(procStmt.Params.List), "param count")
	expectEqual(t, "a", procStmt.Params.List[0].Name, "first param")
	expectEqual(t, "b", procStmt.Params.List[1].Name, "second param")
}

func TestParseProcStmtVarArgs(t *testing.T) {
	file := mustParse(t, "proc foo(...args) { args }")
	procStmt := file.Statements[0].(*ProcStmt)

	expectEqual(t, true, procStmt.Params.VarArgs, "varargs")
	expectEqual(t, 1, len(procStmt.Params.List), "param count")
	expectEqual(t, "args", procStmt.Params.List[0].Name, "varargs param")
}

func TestParseDeclStmt(t *testing.T) {
	file := mustParse(t, "decl x = 5")
	declStmt := file.Statements[0].(*DeclStmt)

	expectEqual(t, "x", declStmt.Ident.Name, "decl ident")

	intLit := declStmt.Val.(*IntLit)
	expectEqual(t, int64(5), intLit.Val, "decl value")
}

func TestParseDeclStmtNoInit(t *testing.T) {
	file := mustParse(t, "decl x")
	declStmt := file.Statements[0].(*DeclStmt)

	expectEqual(t, "x", declStmt.Ident.Name, "decl ident")

	if declStmt.Val != nil {
		fail(t, "expected nil value for uninitialized decl")
	}
}

func TestParseConstStmt(t *testing.T) {
	file := mustParse(t, "const PI = 3.14")
	constStmt := file.Statements[0].(*ConstStmt)

	expectEqual(t, "PI", constStmt.Ident.Name, "const ident")

	floatLit := constStmt.Val.(*FloatLit)
	expectEqual(t, 3.14, floatLit.Val, "const value")
}

func TestParseUninitializedConstStmt(t *testing.T) {
	err := expectParseError(t, "const PI")
	expectEqual(t, "testfile(1, 9): expected '=', found newline", err.Error(), "Uninitialized const error")
}

func TestParseReturnStmt(t *testing.T) {
	file := mustParse(t, "return 42")
	returnStmt := file.Statements[0].(*ReturnStmt)

	intLit := returnStmt.e.(*IntLit)
	expectEqual(t, int64(42), intLit.Val, "return value")
}

func TestParseReturnStmtNoValue(t *testing.T) {
	file := mustParse(t, "return")
	returnStmt := file.Statements[0].(*ReturnStmt)

	if returnStmt.e != nil {
		t.Error("expected nil return value")
	}
}

func TestParseAssignStmt(t *testing.T) {
	file := mustParse(t, "x = 5")
	assignStmt := file.Statements[0].(*AssignStmt)

	expectEqual(t, 1, len(assignStmt.Lhs), "lhs count")
	expectEqual(t, 1, len(assignStmt.Rhs), "rhs count")
	expectEqual(t, token.ASSIGN, assignStmt.EqType, "assign type")
}

func TestParseMultiAssignStmt(t *testing.T) {
	file := mustParse(t, "x, y = 1, 2")
	assignStmt := file.Statements[0].(*AssignStmt)

	expectEqual(t, 2, len(assignStmt.Lhs), "lhs count")
	expectEqual(t, 2, len(assignStmt.Rhs), "rhs count")
}

func TestParseCompoundAssignStmt(t *testing.T) {
	for _, tc := range []struct {
		input string
		op    token.Tok
	}{
		{"x += 1", token.PLUS_EQ},
		{"x -= 1", token.MINUS_EQ},
		{"x *= 2", token.STAR_EQ},
		{"x /= 2", token.SLASH_EQ},
		{"x ^= 1", token.XOR_EQ},
	} {
		file := mustParse(t, tc.input)
		assignStmt := file.Statements[0].(*AssignStmt)
		expectEqual(t, tc.op, assignStmt.EqType, "compound assign type for: %s", tc.input)
	}
}

func TestParseExportStmt(t *testing.T) {
	file := mustParse(t, "export x")
	exportStmt := file.Statements[0].(*ExportStmt)

	ident := exportStmt.e.(*Ident)
	expectEqual(t, "x", ident.Name, "export ident")
}

func TestParseBranchStmt(t *testing.T) {
	for _, input := range []string{"break", "continue"} {
		file := mustParse(t, input)
		branchStmt := file.Statements[0].(*BranchStmt)
		if branchStmt.Pos == 0 {
			fail(t, "expected valid position for: %s", input)
		}
	}
}

func TestParseImportExpr(t *testing.T) {
	file := mustParse(t, `import("std")`)
	exprStmt := file.Statements[0].(*ExprStmt)
	importExpr := exprStmt.e.(*ImportExpr)
	expectEqual(t, "std", importExpr.ModuleName, "module name")
}

func TestParseProcLit(t *testing.T) {
	file := mustParse(t, "x = proc(x) { x }")
	assignStmt := file.Statements[0].(*AssignStmt)
	procLit := assignStmt.Rhs[0].(*ProcLit)

	expectEqual(t, 1, len(procLit.Params.List), "param count")
	expectEqual(t, "x", procLit.Params.List[0].Name, "param name")
}

// --- Operator Precedence Tests ---

func TestOperatorPrecedence(t *testing.T) {
	file := mustParse(t, "1 + 2 * 3")
	exprStmt := file.Statements[0].(*ExprStmt)

	// Should be parsed as 1 + (2 * 3)
	addExpr := exprStmt.e.(*BinaryExpr)
	expectEqual(t, token.PLUS, addExpr.Op, "outer operator should be +")

	mulExpr := addExpr.Rhs.(*BinaryExpr)
	expectEqual(t, token.STAR, mulExpr.Op, "inner operator should be *")
}

// --- Edge Case Tests ---

func TestParseMultipleStatements(t *testing.T) {
	file := mustParse(t, "x = 1\ny = 2\nz = 3")
	expectEqual(t, 3, len(file.Statements), "statement count")
}

func TestParseNestedExpressions(t *testing.T) {
	file := mustParse(t, "((1 + 2) * 3)")
	exprStmt := file.Statements[0].(*ExprStmt)

	grouped := exprStmt.e.(*GroupedExpr)
	mulExpr, ok := grouped.X.(*BinaryExpr)
	if !ok {
		fail(t, "expected BinaryExpr inside GroupedExpr")
	}

	// Should be parsed as (1 + 2) * 3
	expectEqual(t, token.STAR, mulExpr.Op, "outer operator should be *")

	addExpr := mulExpr.Lhs.(*GroupedExpr).X.(*BinaryExpr)
	expectEqual(t, token.PLUS, addExpr.Op, "inner operator should be +")
}
