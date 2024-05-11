#pragma once

#include "shared.hpp"
#include "types/types.hpp"

namespace VM {
    struct Package {
        std::map<std::string, Types::Object *> globals;
    };
}