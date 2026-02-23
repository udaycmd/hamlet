#pragma once

#include <cstdint>

namespace hmach {

enum class Opcode : uint8_t {
  Halt = 32,
  Push,
  Pop,
  Duplicate,
  Swap,
  Negate,
  Not,
  Equal,
  NotEqual,
  Less,
  Greater,
  LessOrEqual,
  GreaterOrEqual,
  Load,
  Store,
  LoadGlobal,
  StoreGlobal,
  Jump,
  Call,
  Return,
  BuiltinCall,
  ArrayLen,
};

} // namespace hmach
