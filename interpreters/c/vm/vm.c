#include "vm.h"

#include "../chunk/chunk_disassembler.h"

struct vm current_vm;

static interpret_result run();
static void print_stack();

interpret_result interpret(struct chunk * chunk) {
    current_vm.chunk = chunk;
    current_vm.pc = chunk->code;

    return run();
}

static interpret_result run() {
#define READ_BYTE() (*current_vm.pc++)
#define READ_CONSTANT() (current_vm.chunk->constants.values[READ_BYTE()])
#define BINARY_OP(op) \
    do { \
        double b = pop_stack_vm(); \
        double a = pop_stack_vm(); \
        push_stack_vm(a op b); \
    }while(false); \

    for(;;) {
#ifdef  DEBUG_TRACE_EXECUTION
        disassemble_chunk_instruction(current_vm.chunk, current_vm.stack - current_vm.esp);
        print_stack();
#endif
        uint8_t instruction;
        switch (instruction = READ_BYTE()) {
            case OP_RETURN: print_value(pop_stack_vm()); break;
            case OP_CONSTANT: push_stack_vm(READ_CONSTANT()); break;
            case OP_NEGATE: push_stack_vm(-pop_stack_vm()); break;
            case OP_ADD: BINARY_OP(+); break;
            case OP_SUB: BINARY_OP(-); break;
            case OP_MUL: BINARY_OP(*); break;
            case OP_DIV: BINARY_OP(/); break;
            case OP_EOF: return INTERPRET_OK;
            default:
                perror("Unhandled bytecode op\n");
                return INTERPRET_RUNTIME_ERROR;
        }
    }

#undef READ_CONSTANT
#undef BINARY_OP
#undef READ_BYTE
}

//TODO Check overflow
void push_stack_vm(lox_value_t value) {
    *current_vm.esp = value;
    current_vm.esp++;
}

//TODO Check underflow
lox_value_t pop_stack_vm() {
    auto val = *--current_vm.esp;
    return val;
}

void start_vm() {
    current_vm.esp = current_vm.stack; //Reset stack
}

void stop_vm() {
}

static void print_stack() {
    printf("\t");
    for(lox_value_t * value = current_vm.stack; value < current_vm.esp; value++)  {
        printf("[");
        print_value(*value);
        printf("]");
    }
    printf("\n");
}