#include "opcode.h"

namespace hmach {

const char* to_string(OpCode op) {
	switch (op) {
	case OpCode::Halt:
		return "hlt";
	case OpCode::Push:
		return "push";
	case OpCode::Pop:
		return "pop";
	case OpCode::Duplicate:
		return "dup";
	case OpCode::Swap:
		return "swp";
	case OpCode::Negate:
		return "neg";
	case OpCode::Not:
		return "not";
	case OpCode::Equal:
		return "eql";
	case OpCode::NotEqual:
		return "neql";
	case OpCode::Less:
		return "lss";
	case OpCode::Greater:
		return "grt";
	case OpCode::LessOrEqual:
		return "leq";
	case OpCode::GreaterOrEqual:
		return "geq";
	case OpCode::Load:
		return "load";
	case OpCode::Store:
		return "str";
	case OpCode::LoadGlobal:
		return "gload";
	case OpCode::StoreGlobal:
		return "gstr";
	case OpCode::Jump:
		return "jmp";
	case OpCode::Call:
		return "call";
	case OpCode::Return:
		return "ret";
	case OpCode::BuiltinCall:
		return "bcall";
	case OpCode::ArrayLen:
		return "alen";
	case OpCode::Const:
		return "ldc";
	}
}

} // namespace hmach
