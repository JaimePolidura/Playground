#include "chunk_disassembler.h"

#include <stdio.h>

static int constant_instruction(const char * name, const struct chunk * chunk, int offset);
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
        case OP_RETURN: return simple_instruction("RETURN", offset);
        case OP_NEGATE: return simple_instruction("NEGATE", offset);
        case OP_CONSTANT: return constant_instruction("CONSTANT", chunk, offset);
        case OP_ADD: return simple_instruction("ADD", offset);
        case OP_SUB: return simple_instruction("SUB", offset);
        case OP_MUL: return simple_instruction("MUL", offset);
        case OP_DIV: return simple_instruction("DIV", offset);
        default:
            printf("Unknown opcode %d\n", instruction);
            return offset + 1;
    }
}

static int simple_instruction(const char * name, const int offset) {
    printf("%s\n", name);
    return offset + 1;
}

static int constant_instruction(const char * name, const struct chunk * chunk, int offset) {
    const uint8_t constant = chunk->code[offset + 1];
    printf("%s '", name);
    print_value(chunk->constants.values[constant]);
    printf("'\n");

    return offset + 2;
}

void print_value(lox_value_t value) {
    printf("%g", value);
}