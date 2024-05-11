#pragma once

#include "types/types.hpp"
#include "shared.hpp"

namespace Types {
    struct StructObject {
        Object object;
        int n_fields;
        Object * fields;
    };
}