#ifndef STREAMS_H
#define STREAMS_H

#include <stdio.h>
#include <stdbool.h>

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

#endif