#pragma once

#include "types/types.hpp"
#include "shared.hpp"

#define AS_STRUCT(object) reinterpret_cast<Types::StructObject *>(object)

namespace Types {
    struct StructObject {
        Object object;
        int nFields;
        Types::Object * fields[];

        static int size(int nFields) {
            return sizeof(Types::StructObject) + (nFields * sizeof(Types::Object *));
        }
    };
}