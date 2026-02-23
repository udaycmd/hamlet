#ifndef HMACH_CHUNK_H
#define HMACH_CHUNK_H

#include <vector>
#include <cstdint>
#include "value.h"

namespace hmach {

class Chunk {
public:
    Chunk() = default;
    ~Chunk() = default;

    void writeChunk(uint8_t byte, int line);
    int addConstant(Value value);

    std::vector<uint8_t>& getCode() { return code; }
    ValueArray& getConstants() { return constants; }

private:
    std::vector<uint8_t> code;
    ValueArray constants;
    std::vector<int> lines; // Line information for debugging
};

} // namespace hmach

#endif // HMACH_CHUNK_H
