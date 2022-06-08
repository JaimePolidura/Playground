#include <stdio.h>
#include <stdbool.h>
#include <malloc.h>


struct Stream {
    char * stream_elements;
    int n_elements;
    size_t size;
};

typedef char * stream_element_t;
typedef struct Stream stream_t;

stream_t * make_stream(void *, int, size_t);
stream_t * map(stream_t *, stream_element_t (* mapper)(stream_element_t));
stream_t * filter(stream_t, bool (* filter)(stream_element_t));
stream_element_t reduce(stream_t, void * (* reducer)(stream_element_t *, stream_element_t));

int multiplyByTwo(int num){
    return num * 2;
}

int main(){
    return 1;
}

stream_t * map(stream_t * stream, stream_element_t (* mapper)(stream_element_t)) {
    int n_elements = stream->n_elements;
    size_t size = stream->size;
    char * stream_elements_mapped = malloc(size * n_elements);

    for(int i = 0; i < n_elements; i++){
        stream_element_t * elementToMap = (stream_element_t *) (stream->stream_elements + (i * size));
        void * element_mapped = mapper(* elementToMap);
        char * element_mapped_p = (char *) &element_mapped;

        for(int j = 0; j < size; j++){
            *(stream_elements_mapped + (i * size) + j) = *(element_mapped_p + j);
        }
    }

    stream_t * stream_result = malloc(sizeof(stream_t));
    stream_result->n_elements = n_elements;
    stream_result->stream_elements = (char *) stream_elements_mapped;
    stream_result->size = size;

    return stream_result;
}

stream_t * make_stream(void * data, int n_elements, size_t size){
    stream_t * stream = malloc(size);
    stream->n_elements = n_elements;
    stream->size = size;
    stream->stream_elements = (char *) data;

    return stream;
}

stream_t * filter(stream_t stream, bool (* filter)(stream_element_t)) {
    return NULL;
}

stream_element_t reduce(stream_t stream, void * (* reducer)(stream_element_t *, stream_element_t)) {
    return NULL;
}







