#pragma once

#include "types/types.hpp"
#include "shared.hpp"

namespace Types {
    struct ArrayObject {
        Object object;
        ObjectType contentType;
        int size;
        Types::Object * content;
    };
}
