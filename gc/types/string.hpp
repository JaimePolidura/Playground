#pragma once

#include "types/types.hpp"
#include "shared.hpp"

namespace Types {
    struct StringObject {
        Object object;
        int size;
        std::byte * contents;
    };
}