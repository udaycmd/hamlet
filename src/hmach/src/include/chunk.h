#pragma once

#include "opcode.h"
#include <vector>

namespace hmach {

class Chunk {
  public:
	Chunk() : m_code() {}
	void write(OpCode op);
	void dis(const char* name);

  private:
	std::vector<OpCode> m_code;
};

} // namespace hmach
