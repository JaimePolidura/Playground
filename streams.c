#include <stdio.h>
#include <stdbool.h>
#include <malloc.h>
#include <math.h>

#pragma clang diagnostic ignored "-Wint-to-void-pointer-cast"

struct Stream {
    char * stream_elements;
    int n_elements;
    size_t size;
};

typedef char * stream_element_t;
typedef struct Stream stream_t;

typedef stream_element_t (* mapper_t)(stream_element_t);
typedef bool (* predicate_t)(void *);
typedef stream_element_t (* reducer_t)(stream_element_t, stream_element_t);
typedef void (* consumer_t)(stream_element_t);

stream_t * make_stream(void *, int, size_t);
stream_t * map(stream_t *, mapper_t);
stream_t * filter(stream_t *, predicate_t);
stream_t * foreach(stream_t *, consumer_t);
stream_element_t reduce(stream_t *, stream_element_t, reducer_t);

static int count_nmatches(char *elements, int n_elements, size_t size, bool (*predicate)(void *));

stream_t * foreach(stream_t * stream, consumer_t consumer){
    int n_elements = stream->n_elements;
    size_t size = stream->size;

    for(int i = 0; i < n_elements; i++){
        stream_element_t * elementToMap = (stream_element_t *) (stream->stream_elements + (i * size));

        consumer(* elementToMap);
    }

    return stream;
}

stream_t * map(stream_t * stream, mapper_t mapper) {
    int n_elements = stream->n_elements;
    size_t size = stream->size;
    char * stream_elements_mapped = malloc(size * n_elements);

    for(int i = 0; i < n_elements; i++){
        stream_element_t * elementToMap = (stream_element_t *) (stream->stream_elements + (i * size));
        void * element_mapped = mapper(* elementToMap);
        char * element_mapped_p = (char *) &element_mapped;

        for(int j = 0; j < size; j++)
            *(stream_elements_mapped + (i * size) + j) = *(element_mapped_p + j);
    }

    stream_t * stream_result = malloc(sizeof(stream_t));
    stream_result->n_elements = n_elements;
    stream_result->stream_elements = stream_elements_mapped;
    stream_result->size = size;

    return stream_result;
}

stream_t * filter(stream_t * stream, predicate_t predicate) {
    int n_elements = stream->n_elements;
    size_t size = stream->size;
    int n_matches = count_nmatches(stream->stream_elements, n_elements, size, predicate);
    char * stream_elements_filtered = malloc(size * n_matches);

    for(int i = 0; i < n_elements; i++){
        stream_element_t * stream_element = (stream_element_t *) (stream->stream_elements + (i * size));
        bool matches = predicate(* stream_element);

        if(matches == true){
            char * element_filtered_p = (char *) stream_element;

            for(int j = 0; j < size; j++)
                *(stream_elements_filtered + (i * size) + j) = *(element_filtered_p + j);
        }
    }

    return make_stream(stream_elements_filtered, n_matches, size);
}

static int count_nmatches(char * elements, int n_elements, size_t size, predicate_t predicate) {
    int count = 0;

    for(int i = 0; i < n_elements; i++){
        void * element = (void *) *(elements + (size * i));
        if(predicate(element) == true)
            count++;
    }

    return count;
}

stream_element_t reduce(stream_t * stream, stream_element_t initial, reducer_t reducer){
    int n_elements = stream->n_elements;
    size_t size = stream->size;
    stream_element_t actumulator = initial;

    for(int i = 0; i < n_elements; i++){
        stream_element_t * stream_element = (stream_element_t *) (stream->stream_elements + (i * size));
        actumulator = reducer(actumulator, *stream_element);
    }

    return actumulator;
}

stream_t * make_stream(void * data, int n_elements, size_t size){
    stream_t * stream = malloc(sizeof(stream_t));
    stream->n_elements = n_elements;
    stream->size = size;
    stream->stream_elements = (char *) data;

    return stream;
}