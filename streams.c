#include <stdio.h>
#include <stdbool.h>
#include <malloc.h>

struct Stream {
    void * stream_elements;
    int n_elements;
    size_t size;
};

typedef struct Stream stream_t;
typedef void * stream_element_t;

stream_t * make_stream(void *, int, size_t);
stream_t * map(stream_t *, stream_element_t * (* mapper)(stream_element_t));
stream_t * filter(stream_t, bool (* filter)(stream_element_t));
stream_element_t reduce(stream_t, void * (* reducer)(stream_element_t *, stream_element_t));

int multiplyByTwo(int num){
    return num * 2;
}

int main(){
    int * array = malloc(sizeof(int));
    * array = 0;
    * (array + 1) = 1;
    * (array + 2) = 2;

    stream_t * stream = make_stream(array, 3, sizeof(int));
    stream_t * stream_map_result = map(stream, (stream_element_t *(*)(stream_element_t)) multiplyByTwo);
    void * data_map_result = stream_map_result->stream_elements;

    for(int i = 0; i < 3; i++){
        int * data_p = (int *) data_map_result + (i * 2);

        printf("b: %i\n", *data_p);
    }

    return 1;
}

stream_t * map(stream_t * stream, stream_element_t * (* mapper)(stream_element_t)) {
    int n_elements = stream->n_elements;
    size_t size = stream->size;
    stream_element_t * stream_elements_mapped = malloc(size * n_elements);

    for(int i = 0; i < n_elements; i++){
        stream_element_t * elementToMap = stream->stream_elements + (i * size);
        stream_element_t * element_mapped = mapper(*elementToMap);

        *(stream_elements_mapped + i) = element_mapped;
    }

    stream_t * stream_result = malloc(sizeof(stream_t));
    stream_result->n_elements = n_elements;
    stream_result->stream_elements = stream_elements_mapped;
    stream_result->size = size;

    return stream_result;
}

stream_t * make_stream(void * data, int n_elements, size_t size){
    stream_t * stream = malloc(size);
    stream->n_elements = n_elements;
    stream->size = size;
    stream->stream_elements = data;

    return stream;
}

stream_t * filter(stream_t stream, bool (* filter)(stream_element_t)) {
    return NULL;
}

stream_element_t reduce(stream_t stream, void * (* reducer)(stream_element_t *, stream_element_t)) {
    return NULL;
}







