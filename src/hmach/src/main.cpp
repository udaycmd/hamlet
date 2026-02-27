#include "chunk.h"
#include "opcode.h"

int main() {
	hmach::Chunk c;
	c.write(hmach::OpCode::Return);
	c.dis("test");
}
