#include "chunk.h"

#include <stdlib.h>

static void resize_chunk(struct chunk * chunk, int new_capacity);

struct chunk * alloc_chunk() {
    struct chunk * allocated_chunk = malloc(sizeof(struct chunk));
    allocated_chunk->capacity = 0;
    allocated_chunk->in_use = 0;
    allocated_chunk->code = NULL;

    return allocated_chunk;
}

void write_chunk(struct chunk * chunk, uint8_t byte) {
    if(chunk->in_use + 1 > chunk->capacity) {
        resize_chunk(chunk, GROW_ARRAY_CAPACITY(chunk->capacity));
    }

    chunk->code[chunk->in_use++] = byte;
}

void free_chunk(struct chunk * chunk) {
    free(chunk->code);
}

static void resize_chunk(struct chunk * chunk, int new_capacity) {
    const int old_capacity = chunk->capacity;
    chunk->capacity = new_capacity;
    chunk->code = reallocate_array(chunk->code, old_capacity, new_capacity);
}