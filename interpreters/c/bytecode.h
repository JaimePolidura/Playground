#pragma once

#include "shared.h"

typedef enum {
    OP_CONSTANT,
    OP_RETURN,
    OP_NEGATE,
    OP_EOF,
    OP_ADD,
    OP_SUB,
    OP_MUL,
    OP_DIV
} op_code;

