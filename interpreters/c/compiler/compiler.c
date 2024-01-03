#include "compiler.h"

struct parser {
    struct token current;
    struct token previous;
    bool has_error;
};

struct compiler {
    struct scanner scanner;
    struct parser parser;
    struct chunk chunk;
};

static void report_error(struct compiler * compiler, struct token token, const char * message);
static void advance(struct compiler * compiler);
static void init_parser(struct parser * parser);
static void consume(struct compiler * compiler, tokenType_t expected_token_type, const char * error_message);
static void emit_bytecode(struct compiler * compiler, uint8_t bytecode);
static void emit_bytecodes(struct compiler * compiler, uint8_t bytecodeA, uint8_t bytecodeB);
static struct compiler * alloc_compiler();

bool compile(char * source_code, struct chunk * output_chunk) {
    struct compiler * compiler = alloc_compiler();

    //TODO

    emit_bytecode(compiler, OP_RETURN);

    bool is_success = !compiler->parser.has_error;
    free(compiler);
    return is_success;
}

static void advance(struct compiler * compiler) {
    compiler->parser.previous = compiler->parser.current;

    struct token token = next_token_scanner(&compiler->scanner);
    compiler->parser.current = token;

    if(token.type == TOKEN_ERROR) {
        report_error(compiler, token, "");
    }
}

static void emit_bytecodes(struct compiler * compiler, uint8_t bytecodeA, uint8_t bytecodeB) {
    emit_bytecode(compiler, bytecodeA);
    emit_bytecode(compiler, bytecodeB);
}

static void emit_bytecode(struct compiler * compiler, uint8_t bytecode) {
    write_chunk(&compiler->chunk, bytecode, compiler->parser.previous.line);
}

static void consume(struct compiler * compiler, tokenType_t expected_token_type, const char * error_message) {
    if(compiler->parser.current.type == expected_token_type) {
        advance(compiler);
        return;
    }

    report_error(compiler, compiler->parser.current, error_message);
}

static void report_error(struct compiler * compiler, struct token token, const char * message) {
    fprintf(stderr, "[line %d] Error", token.line);
    if (token.type == TOKEN_EOF) {
        fprintf(stderr, " at end");
    } else if (token.type == TOKEN_ERROR) {
        // Nothing.
    } else {
        fprintf(stderr, " at '%.*s'", token.length, token.start);
    }
    fprintf(stderr, ": %s\n", message);
    compiler->parser.has_error = true;
}

static void init_parser(struct parser * parser) {
    parser->has_error = false;
}

static struct compiler * alloc_compiler() {
    struct compiler * compiler = malloc(sizeof(struct compiler));
    init_scanner(&compiler->scanner);
    init_parser(&compiler->parser);
    init_chunk(&compiler->chunk);
    return compiler;
}