#pragma once

#include "shared.hpp"

namespace Types {
    enum ObjectType {
        ARRAY,
        STRUCT,
        STRING,
    };

    // Force 8 bytes alignment so that we can allocate Object pointers on a
    // raw byte array without casing problems
    struct alignas(8) Object {
        ObjectType type;
        void * gc{nullptr}; //TODO Embed gc with type in 8 bytes
    };
}