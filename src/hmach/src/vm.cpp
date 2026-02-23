#include "vm.h"
#include "opcode.h"
#include <iostream>

namespace hmach {

VM::VM() {
    stackTop = stack;
}

VM::~VM() {
}

void VM::push(Value value) {
    *stackTop = value;
    stackTop++;
}

Value VM::pop() {
    stackTop--;
    return *stackTop;
}

InterpretResult VM::interpret(Chunk* chunk) {
    this->chunk = chunk;
    this->ip = chunk->getCode().data();
    return run();
}

uint8_t VM::readByte() {
    return *ip++;
}

Value VM::readConstant() {
    return chunk->getConstants()[readByte()];
}

InterpretResult VM::run() {
    for (;;) {
        uint8_t instruction;
        switch (instruction = readByte()) {
            case static_cast<uint8_t>(Opcode::Push): {
                Value constant = readConstant();
                push(constant);
                break;
            }
            case static_cast<uint8_t>(Opcode::Pop): {
                pop();
                break;
            }
            case static_cast<uint8_t>(Opcode::Halt): {
                return InterpretResult::INTERPRET_OK;
            }
            // Add other opcodes here as needed
            default:
                std::cerr << "Unknown opcode: " << (int)instruction << std::endl;
                return InterpretResult::INTERPRET_RUNTIME_ERROR;
        }
    }
}

} // namespace hmach
