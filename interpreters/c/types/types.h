#pragma once

#include "../shared.h"
#include "../memory/memory.h"

struct object;

typedef enum {
  VAL_BOOL,
  VAL_NIL,
  VAL_NUMBER,
  VAL_OBJ,
} lox_value_type;

typedef struct {
  lox_value_type type;
  union {
    bool boolean;
    double number;
    struct object * object;
  } as;
} lox_value_t;

#define FROM_NUMBER(value) ((lox_value_t){VAL_NUMBER, {.number = value}})
#define FROM_BOOL(value) ((lox_value_t){VAL_BOOL, {.boolean = value}})
#define FROM_NIL ((lox_value_t){VAL_NIL, {.number = 0}})

#define TO_NUMBER(value) ((value).as.number)
#define TO_BOOL(value) ((value).as.boolean)

#define IS_BOOL(value) ((value).type == VAL_BOOL)
#define IS_NIL(value) ((value).type == VAL_NIL)
#define IS_NUMBER(value) ((value).type == VAL_NUMBER)

struct lox_array {
  lox_value_t * values;
  int capacity;
  int in_use;
};

void alloc_lox_array(struct lox_array * array);
void write_lox_array(struct lox_array * array, lox_value_t value);
void free_lox_array(struct lox_array * array);