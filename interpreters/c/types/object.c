#include "object.h"

struct string_object * chars_to_string_object(const char * chars, int length) {
    char * string_ptr = malloc(sizeof(char) * length + 1);
    memcpy(string_ptr, chars, length);
    string_ptr[length] = '\0';

    struct string_object * string = malloc(sizeof(struct string_object));
    string->length = length + 1;
    string->chars = string_ptr;
    string->object.type = OBJ_STRING;

    return string;
}