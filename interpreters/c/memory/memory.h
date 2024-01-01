#pragma once

#include "../shared.h"

#define GROW_ARRAY_CAPACITY(capacity) (capacity < 8 ? 8 : capacity << 2)

void * reallocate_array(void * ptr, int old_size, int new_size);