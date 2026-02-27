#include "chunk.h"
#include "opcode.h"
#include "value.h"
#include <cstdio>

namespace hmach {

int dis_simple_inst(OpCode op, int offset) {
	std::printf("%s\n", to_string(op));
	return ++offset;
}

void Chunk::write(OpCode op) { m_code.emplace_back(op); }

void Chunk::write_const(Value v) { m_consts.emplace_back(v); }

void Chunk::dis(const char* name) const {
	std::printf("=== %s ===\n", name);

	for (int offset = 0; offset < static_cast<int>(m_code.size());) {
		std::printf("%05d ", offset);
		auto op = m_code[offset];

		switch (op) {
		case hmach::OpCode::Return:
			offset = dis_simple_inst(op, offset);
			break;
		default:
			offset += 1;
		}
	}
}

} // namespace hmach
