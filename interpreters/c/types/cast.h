#pragma once

#include "../shared.h"
#include "object.h"

//If value is type number, the returned char * will be allocated in the heap
char * cast_to_string(lox_value_t value);