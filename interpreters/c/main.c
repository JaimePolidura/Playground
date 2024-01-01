#include <stdio.h>
#include <stdlib.h>

#include "shared.h"
#include "chunk/chunk.h"
#include "chunk/chunk_disassembler.h"

int main(int argc, char* args[]) {
    struct chunk * chunk = alloc_chunk();
    write_chunk(chunk, 0x00);
    disassemble_chunk(chunk, "Chunk #1");
    free_chunk(chunk);


    printf("Hello, World!\n");
    return 0;
}
