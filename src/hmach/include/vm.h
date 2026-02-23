#ifndef HMACH_VM_H
#define HMACH_VM_H

#include "chunk.h"
#include "value.h"

namespace hmach {

enum class InterpretResult {
    INTERPRET_OK,
    INTERPRET_COMPILE_ERROR,
    INTERPRET_RUNTIME_ERROR
};

class VM {
public:
    VM();
    ~VM();

    InterpretResult interpret(Chunk* chunk);

    void push(Value value);
    Value pop();

private:
    Chunk* chunk;
    uint8_t* ip;

    static const int STACK_MAX = 256;
    Value stack[STACK_MAX];
    Value* stackTop;

    InterpretResult run();
    uint8_t readByte();
    Value readConstant();
};

} // namespace hmach

#endif // HMACH_VM_H
