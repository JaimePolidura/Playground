#include "vm.h"

#include "../chunk/chunk_disassembler.h"

struct vm current_vm;

static double check_number();
static bool check_boolean();
static interpret_result run();
static void print_stack();
static void runtime_errpr(char * format, ...);
static lox_value_t values_equal(lox_value_t a, lox_value_t b);

interpret_result interpret_vm(struct chunk * chunk) {
    current_vm.chunk = chunk;
    current_vm.pc = chunk->code;

    return run();
}

static interpret_result run() {
#define READ_BYTE() (*current_vm.pc++)
#define READ_CONSTANT() (current_vm.chunk->constants.values[READ_BYTE()])
#define BINARY_OP(op) \
    do { \
        double b = check_number(pop_stack_vm()); \
        double a = check_number(pop_stack_vm()); \
        push_stack_vm(FROM_NUMBER(a op b)); \
    }while(false);

#define COMPARATION_OP(op) \
    do { \
        double b = check_number(pop_stack_vm()); \
        double a = check_number(pop_stack_vm()); \
        push_stack_vm(FROM_BOOL(a op b)); \
    }while(false);

    for(;;) {
#ifdef  DEBUG_TRACE_EXECUTION
        disassemble_chunk_instruction(current_vm.chunk, current_vm.stack - current_vm.esp);
        print_stack();
#endif
        switch (READ_BYTE()) {
            case OP_RETURN: print_value(pop_stack_vm()); break;
            case OP_CONSTANT: push_stack_vm(READ_CONSTANT()); break;
            case OP_NEGATE: push_stack_vm(FROM_NUMBER(-check_number())); break;
            case OP_ADD: BINARY_OP(+) break;
            case OP_SUB: BINARY_OP(-) break;
            case OP_MUL: BINARY_OP(*) break;
            case OP_DIV: BINARY_OP(/) break;
            case OP_GREATER: COMPARATION_OP(>) break;
            case OP_LESS: COMPARATION_OP(<) break;
            case OP_FALSE: push_stack_vm(FROM_BOOL(false)); break;
            case OP_TRUE: push_stack_vm(FROM_BOOL(true)); break;
            case OP_NIL: push_stack_vm(FROM_NIL); break;
            case OP_NOT: push_stack_vm(FROM_BOOL(!check_boolean())); break;
            case OP_EQUAL: push_stack_vm(values_equal(pop_stack_vm(), pop_stack_vm())); break;
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

static double check_number() {
    lox_value_t value = pop_stack_vm();
    if(IS_NUMBER(value)) {
        return TO_NUMBER(value);
    } else {
        runtime_errpr("Operand must be a number.");
        exit(1);
    }
}

static bool check_boolean() {
    lox_value_t value = pop_stack_vm();
    if(IS_BOOL(value)) {
        return TO_BOOL(value);
    } else {
        runtime_errpr("Operand must be a boolean.");
        exit(1);
    }
}

static lox_value_t values_equal(lox_value_t a, lox_value_t b) {
    if(a.type != b.type) {
        return FROM_BOOL(false);
    }

    switch (a.type) {
        case VAL_NIL: return FROM_BOOL(true);
        case VAL_NUMBER: return FROM_BOOL(a.as.number == b.as.number);
        case VAL_BOOL: return FROM_BOOL(a.as.boolean == b.as.boolean);
        case VAL_OBJ: return FROM_BOOL(strcmp(TO_STRING_CHARS(a), TO_STRING_CHARS(b)) == 0);
        default:
            runtime_errpr("Operator '==' not supported for that type");
            return FROM_BOOL(false); //Unreachable, runtime_error executes exit()
    }
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

static void runtime_errpr(char * format, ...) {
    va_list args;
    va_start(args, format);
    vfprintf(stderr, format, args);
    va_end(args);
    fputs("\n", stderr);

    size_t instruction = ((uint8_t *) current_vm.esp) - current_vm.chunk->code - 1;
    int line = current_vm.chunk->lines[instruction];
    fprintf(stderr, "[line %d] in script\n", line);
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