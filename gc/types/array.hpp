#pragma once

#include "types/types.hpp"
#include "shared.hpp"

#define AS_ARRAY(object) reinterpret_cast<Types::ArrayObject *>(object)

namespace Types {
    struct ArrayObject {
        Object object;
        int nElements;
        Types::Object * elements[];

        static int size(int nElements) {
            return sizeof(Types::ArrayObject) + (nElements * sizeof(Types::Object *));
        }
    };
}
