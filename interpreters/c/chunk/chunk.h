#pragma once

#include "../shared.h"
#include "../memory/memory.h"
#include "../types/types.h"

struct chunk {
    struct lox_array constants;
    uint8_t * code;
    int capacity; //nº of code allocated
    int in_use; //nº of code used
};

void write_chunk(struct chunk * chunk_to_write, uint8_t byte);
int add_constant_to_chunk(struct chunk * chunk_to_write, lox_number_t constant);
struct chunk * alloc_chunk();
void free_chunk(struct chunk * chunk_to_free);