#pragma once

namespace Types {
    enum ObjectType {
        ARRAY,
        STRUCT,
        STRING,
        NUMBER
    };

    // Force 8 bytes alignment so that we can allocate Object pointers on a
    // raw byte array without casing problems
    struct alignas(8) Object {
        ObjectType type;
    };
}