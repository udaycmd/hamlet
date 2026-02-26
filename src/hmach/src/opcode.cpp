#include "opcode.h"

namespace hmach {

const char* to_string(hmach::OpCode op) {
	switch (op) {
	case hmach::OpCode::Return:
		return "ret";
	default:
		return "<not_impl>";
	}
}

} // namespace hmach
