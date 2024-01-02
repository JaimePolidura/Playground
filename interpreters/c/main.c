#include <stdio.h>
#include <stdlib.h>

#include "shared.h"
#include "chunk/chunk.h"
#include "chunk/chunk_disassembler.h"
#include "vm/vm.h"

void debug_simple_calculation();
void prod();

int main(int argc, char* args[]) {
    debug_simple_calculation();
    return 0;
}

void debug_simple_calculation() {
    start_vm();
    struct chunk * chunk = alloc_chunk();

    // -((1.2 + 3.4) / 5.6)
    write_chunk(chunk, OP_CONSTANT, 1);
    write_chunk(chunk, add_constant_to_chunk(chunk, 1.2), 1);
    write_chunk(chunk, OP_CONSTANT, 1);
    write_chunk(chunk, add_constant_to_chunk(chunk, 3.4), 1);
    write_chunk(chunk, OP_ADD, 1);
    write_chunk(chunk, OP_CONSTANT, 1);
    write_chunk(chunk, add_constant_to_chunk(chunk, 5.6), 1);
    write_chunk(chunk, OP_DIV, 1);
    write_chunk(chunk, OP_NEGATE, 1);
    write_chunk(chunk, OP_RETURN, 1);
    write_chunk(chunk, OP_EOF, 1); //Expect -0.81

    interpret(chunk);

    stop_vm();
    free_chunk(chunk);
}

void prod() {
    start_vm();



    stop_vm();
}