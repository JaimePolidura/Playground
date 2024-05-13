#pragma once

#include "shared.hpp"

namespace Types {
    enum ObjectType {
        ARRAY,
        STRUCT,
        STRING,
    };

    struct Object {
        ObjectType type;
        void * gc{nullptr}; //TODO Embed gc with type in 8 bytes
    };
}