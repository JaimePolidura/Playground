#include "chunk.h"

#include <stdlib.h>

static void resize_chunk(struct chunk * chunk_to_resize, int new_capacity);

struct chunk * alloc_chunk() {
    struct chunk * allocated_chunk = malloc(sizeof(struct chunk));
    allocated_chunk->capacity = 0;
    allocated_chunk->in_use = 0;
    allocated_chunk->code = NULL;
    alloc_lox_array(&allocated_chunk->constants);

    return allocated_chunk;
}

int add_constant_to_chunk(struct chunk * chunk_to_write, lox_value_t constant) {
    write_lox_array(&chunk_to_write->constants, constant);
    return chunk_to_write->constants.in_use - 1;
}

void write_chunk(struct chunk * chunk_to_write, uint8_t byte) {
    if(chunk_to_write->in_use + 1 > chunk_to_write->capacity) {
        resize_chunk(chunk_to_write, GROW_ARRAY_CAPACITY(chunk_to_write->capacity));
    }

    chunk_to_write->code[chunk_to_write->in_use++] = byte;
}

void free_chunk(struct chunk * chunk_to_free) {
    free(chunk_to_free->code);
    free(&chunk_to_free->constants);
}

static void resize_chunk(struct chunk * chunk_to_resize, int new_capacity) {
    const int old_capacity = chunk_to_resize->capacity;
    chunk_to_resize->capacity = new_capacity;
    chunk_to_resize->code = reallocate_array(chunk_to_resize->code, old_capacity, new_capacity);
}