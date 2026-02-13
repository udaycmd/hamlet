// Copyright 2025 Uday Tiwari. All rights reserved.
// Use of this source code is governed by MIT
// license that can be found in the LICENSE file.

package hmach

type Opcode byte

const (
	Push Opcode = iota + 32
	Pop
	Duplicate
	Swap
	Negate
	Not
	Equal
	NotEqual
	Less
	Greater
	LessOrEqual
	GreaterOrEqual
	Load
	Store
	LoadGlobal
	StoreGlobal
	Jump
	Call
	Return
	BuiltinCall
	ArrayLen
)

var opcodes = [...]string{
	Push:           "push",
	Pop:            "pop",
	Duplicate:      "dup",
	Swap:           "swp",
	Negate:         "neg",
	Not:            "not",
	Equal:          "eq",
	NotEqual:       "neq",
	Less:           "lss",
	Greater:        "grt",
	LessOrEqual:    "loe",
	GreaterOrEqual: "goe",
	Load:           "ldl",
	Store:          "stl",
	LoadGlobal:     "ldg",
	StoreGlobal:    "stg",
	Jump:           "jmp",
	Call:           "call",
	Return:         "ret",
	BuiltinCall:    "bcall",
	ArrayLen:       "alen",
}

func (op Opcode) String() string {
	if int(op) >= len(opcodes) || opcodes[op] == "" {
		panic("unknown opcode encountered")
	}

	return opcodes[op]
}
