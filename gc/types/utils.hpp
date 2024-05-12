#pragma once

#include "array.hpp"
#include "struct.hpp"
#include "string.hpp"
#include "types.hpp"

namespace Types {
    std::size_t sizeofObject(ObjectType type);

    void traverseObjectDeep(Types::Object * object, std::function<bool(Types::Object *)> callback);
}
