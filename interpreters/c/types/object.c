#include "object.h"

static void add_heap_object(struct object * object);
static struct string_object * alloc_string_object(char * string_ptr, int length);
static uint32_t hash_string(const char * string_ptr, int length);

extern struct vm current_vm;

struct string_object * from_chars_to_string_object(char * chars, int length) {
    return alloc_string_object(chars, length);
}

struct string_object * copy_chars_to_string_object(const char * chars, int length) {
    char * string_ptr = malloc(sizeof(char) * length + 1);
    memcpy(string_ptr, chars, length);
    string_ptr[length] = '\0';

    return alloc_string_object(string_ptr, length);
}

static void add_heap_object(struct object * object) {
    object->next = current_vm.heap;
    current_vm.heap = object;
}

static struct string_object * alloc_string_object(char * string_ptr, int length) {
    uint32_t string_hash = hash_string(string_ptr, length);

    struct string_object * interned_string = get_key_by_hash(&current_vm.strings, string_hash);
    bool string_already_interned = interned_string != NULL;

    if(!string_already_interned){
        struct string_object * string = malloc(sizeof(struct string_object));
        string->object.type = OBJ_STRING;
        string->hash = string_hash;
        string->chars = string_ptr;
        string->length = length;

        put_hash_table(&current_vm.strings, string, FROM_NIL());
        add_heap_object((struct object *) string);

        return string;
    } else {
        return interned_string;
    }
}

static uint32_t hash_string(const char * string_ptr, int length) {
    uint32_t hash = 2166136261u;
    for (int i = 0; i < length; i++) {
        hash ^= string_ptr[i];
        hash *= 16777619;
    }

    return hash;
}