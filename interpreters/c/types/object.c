#include "object.h"

static void add_heap_object(struct object * object);
static struct string_object * alloc_string_object(char * string_ptr, int length);
static uint32_t hash_string(const char * string_ptr, int length);

extern struct vm current_vm;

struct string_object * from_chars_to_string_object(char * chars, int length) {
    struct string_object * string =  alloc_string_object(chars, length);
    add_heap_object((struct object *) string);

    return string;
}

struct string_object * copy_chars_to_string_object(const char * chars, int length) {
    char * string_ptr = malloc(sizeof(char) * length + 1);
    memcpy(string_ptr, chars, length);
    string_ptr[length] = '\0';

    struct string_object * string =  alloc_string_object(string_ptr, length);
    add_heap_object((struct object *) string);

    return string;
}

static void add_heap_object(struct object * object) {
    object->next = current_vm.heap;
    current_vm.heap = object;
}

static struct string_object * alloc_string_object(char * string_ptr, int length) {
    struct string_object * string = malloc(sizeof(struct string_object));
    string->hash = hash_string(string_ptr, length);
    string->object.type = OBJ_STRING;
    string->chars = string_ptr;
    string->length = length;

    return string;
}

static uint32_t hash_string(const char * string_ptr, int length) {
    uint32_t hash = 2166136261u;
    for (int i = 0; i < length; i++) {
        hash ^= string_ptr[i];
        hash *= 16777619;
    }

    return hash;
}