#pragma once

#include "shared.hpp"
#include "types/types.hpp"

namespace VM {
    struct Global {
        Types::Object * value;
    };

    struct Package {
        std::map<std::string, VM::Global> globals;
    };
}