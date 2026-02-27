#pragma once

#include <cstdint>

namespace hmach {

using real_type = double;
using int_type = std::int64_t;

enum OType : uint8_t {
	hInt,
	hReal,
	hStr,
};

union Value {
	real_type real_num;
	int_type  int_num;
};

struct hObject {
	OType t;
	Value v;
};

} // namespace hmach
