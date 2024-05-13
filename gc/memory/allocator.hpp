#pragma once

#include "types/string.hpp"
#include "types/array.hpp"
#include "types/struct.hpp"

namespace Memory {
    class Allocator {
        virtual Types::StringObject * allocString(char * data) = 0;

        virtual Types::StructObject * allocStruct(int nFields) = 0;

        virtual Types::ArrayObject * allocArray(int nElements) = 0;
    };
}