#ifndef HMACH_VALUE_H
#define HMACH_VALUE_H

#include <variant>
#include <vector>
#include <string>

namespace hmach {

// Assuming basic value types for now: double for numbers, boolean, etc.
using Value = std::variant<std::monostate, double, bool, std::string>;

// Vector to store a pool of values/constants
using ValueArray = std::vector<Value>;

} // namespace hmach

#endif // HMACH_VALUE_H
