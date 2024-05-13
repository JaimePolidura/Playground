#pragma once

#include "types/types.hpp"
#include "shared.hpp"

#define AS_STRING(object) reinterpret_cast<Types::StringObject *>(object)

namespace Types {
    struct StringObject {
        Object object;
        int nChars;
        char contents[]; //Aligned by allocator in 8 byte boundary

        static int size(char * ptr) {
            return sizeof(Types::StringObject) + std::ceil(strlen(ptr) / sizeof(Types::Object));
        }
    };
}