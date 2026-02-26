#pragma once

#include "opcode.h"
#include <vector>

namespace hmach {

struct Chunk {
	std::vector<OpCode> m_code;

	Chunk() : m_code() {}
	void write(OpCode op);
	void dis(const char* name);
};

} // namespace hmach
