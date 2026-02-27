#pragma once

#include "opcode.h"
#include "value.h"
#include <vector>

namespace hmach {

class Chunk {
  public:
	Chunk() : m_code(), m_consts() {}
	void write(OpCode op);
	void write_const(Value v);
	void dis(const char* name) const;

  private:
	std::vector<OpCode> m_code;
	std::vector<Value>	m_consts;
};

} // namespace hmach
