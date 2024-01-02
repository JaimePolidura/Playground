#include <stdio.h>
#include <stdlib.h>

#include "shared.h"
#include "chunk/chunk.h"
#include "chunk/chunk_disassembler.h"
#include "vm/vm.h"

void debug();
void prod();

int main(int argc, char* args[]) {
    debug();
    return 0;
}

void debug() {
    start_vm();
    struct chunk * chunk = alloc_chunk();
    write_chunk(chunk, OP_CONSTANT, 1);
    write_chunk(chunk, add_constant_to_chunk(chunk, 10), 1); //add_constant_to_chunk returns offset
    write_chunk(chunk, OP_RETURN, 1);

    interpret(chunk);

    stop_vm();
    free_chunk(chunk);
}

void prod() {
    start_vm();



    stop_vm();
}