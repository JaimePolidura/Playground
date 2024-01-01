#pragma once

#include "shared.h"
#include "memory.h"

struct chunk {
    uint8_t * code;
    int capacity; //nº of code allocated
    int in_use; //nº of code used
} chunk;

void write_chunk(struct chunk * chunk, uint8_t byte);
struct chunk * alloc_chunk();
void free_chunk(struct chunk * chunk);