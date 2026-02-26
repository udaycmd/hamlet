#include "chunk.h"
#include "opcode.h"

int main() {
	hmach::Chunk chunk;
	chunk.write(hmach::OpCode::Return);
	chunk.dis("test_chunk");
}
