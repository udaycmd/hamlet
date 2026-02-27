#pragma once

#include "value.h"
#include <stack>

namespace hmach {

struct hVM {
	std::stack<hObject> stk;
};

} // namespace hmach
