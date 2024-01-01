#include "memory.h"

#include <stdlib.h>

void * reallocate_array(void * ptr, int old_size, int new_size) {
    if(new_size == 0) {
        free(ptr);
        return NULL;
    }

    void * result = realloc(ptr, new_size);
    if (result == NULL) {
        perror("Out of memory");
        exit(1);
    }

    return result;
}
