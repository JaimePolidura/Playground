#include "chunk_disassembler.h"

#include <stdio.h>

static int simple_instruction(const char * name, int offset);

void disassemble_chunk(const struct chunk * chunk, char * name) {
    printf("== %s ==\n", name);

    for(int offset = 0; offset < chunk->in_use;) {
        offset = disassemble_chunk_instruction(chunk, offset);
    }
}

int disassemble_chunk_instruction(const struct chunk * chunk, const int offset) {
    const uint8_t instruction = chunk->code[offset];
    switch (instruction) {
        case OP_RETURN:
            return simple_instruction("RETURN", offset);
        default:
            printf("Unknown opcode %d\n", instruction);
            return offset + 1;
    }
}

static int simple_instruction(const char * name, const int offset) {
    printf("%s\n", name);
    return offset + 1;
}