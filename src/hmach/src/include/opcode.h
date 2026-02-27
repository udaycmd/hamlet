#pragma once

#include <cstdint>

namespace hmach {

enum class OpCode : uint8_t {
	Halt,
	Push,
	Pop,
	Duplicate,
	Swap,
	Negate,
	Not,
	Equal,
	NotEqual,
	Less,
	Greater,
	LessOrEqual,
	GreaterOrEqual,
	Load,
	Store,
	LoadGlobal,
	StoreGlobal,
	Jump,
	Call,
	Return,
	BuiltinCall,
	ArrayLen,
	Const
};

const char* to_string(OpCode op);

} // namespace hmach
