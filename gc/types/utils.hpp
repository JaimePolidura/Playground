#pragma once

#include "array.hpp"
#include "struct.hpp"
#include "string.hpp"
#include "types.hpp"

namespace Types {
    std::size_t sizeofObject(Types::Object * object);

    void traverseObjectDeep(Types::Object * object, std::function<bool(Types::Object *)> callback);

    void copy(Types::Object * dst, Types::Object * src);
}
