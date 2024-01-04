#pragma once

#include "../chunk/chunk.h"
#include "../shared.h"
#include "../types/cast.h"
#include "../table/table.h"

#define STACK_MAX 256

struct vm {
    struct chunk * chunk;
    uint8_t * pc; // Actual instruction
    lox_value_t stack[STACK_MAX];
    lox_value_t * esp; // Top of the stack
    struct object * heap; // Linkedlist of heap allocated objects
    struct hash_table strings;
};

typedef enum {
    INTERPRET_OK,
    INTERPRET_COMPILE_ERROR,
    INTERPRET_RUNTIME_ERROR,
} interpret_result;

interpret_result interpret_vm(struct chunk* chunk);

void push_stack_vm(lox_value_t value);
lox_value_t pop_stack_vm();

void start_vm();
void stop_vm();