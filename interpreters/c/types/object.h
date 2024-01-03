#pragma once

#include "../shared.h"

#include "types.h"

typedef enum {
    OBJ_STRING,
} object_type_t;

struct object {
    object_type_t type;
};

struct string_object {
    struct object object;
    int length;
    char * chars;
};

struct string_object * chars_to_string_object(const char * chars, int length);

#define FROM_OBJECT(value) ((lox_value_t){VAL_OBJ, {.object = (struct object*) value}})

#define TO_OBJECT(value) ((value).as.object)
#define TO_STRING(value) ((struct string_object *)TO_OBJECT(value))
#define TO_STRING_CHARS(value) (((struct string_object *) value.as.object)->chars)

#define IS_STRING(value) is_object_type(value, OBJ_STRING)
#define IS_OBJECT(value) ((value).type == VAL_OBJ)

#define OBJECT_TYPE(value) (TO_OBJECT(value)->type)

static inline bool is_object_type(lox_value_t value, object_type_t type) {
    return IS_OBJECT(value) && TO_OBJECT(value)->type == type;
}
