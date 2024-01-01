#include "types.h"

#include <stdlib.h>

void alloc_lox_array(struct lox_array * array) {
    array->values = NULL;
    array->capacity = 0;
    array->in_use = 0;
}

void write_lox_array(struct lox_array * array, lox_value_t value) {
    if(array->in_use + 1 > array->capacity) {
        const int old_capacity = array->capacity;
        const int new_capacity = GROW_ARRAY_CAPACITY(array->capacity);
        array->capacity = new_capacity;
        array->values = reallocate_array(array->values, old_capacity, new_capacity);
    }

    array->values[array->in_use++] = value;
}

void free_lox_array(struct lox_array * array){
    free(array->values);
}