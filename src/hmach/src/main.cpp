#include <iostream>
#include "vm.h"
#include "chunk.h"
#include "opcode.h"

int main() {
    hmach::VM vm;
    hmach::Chunk chunk;

    int constant = chunk.addConstant(1.2);
    chunk.writeChunk(static_cast<uint8_t>(hmach::Opcode::Push), 123);
    chunk.writeChunk(constant, 123);

    chunk.writeChunk(static_cast<uint8_t>(hmach::Opcode::Halt), 123);

    hmach::InterpretResult result = vm.interpret(&chunk);
    
    if (result == hmach::InterpretResult::INTERPRET_OK) {
        std::cout << "Executed successfully.\n";
    } else {
        std::cout << "Execution failed.\n";
    }

    return 0;
}
