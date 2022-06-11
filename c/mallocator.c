#include <stdio.h>
#include <malloc.h>

struct Node {
    struct node* next;
    unsigned int size;
    unsigned short empty;
};
typedef struct Node node_t;
void * allocatem(size_t size);
void freem(void *);

static void init_node(node_t *);

node_t * first;

int main(){
    node_t first_node;
    first = &first_node;
    init_node(first);

    char * pointer = allocatem(sizeof("hola"));

    return 0;
}

void * allocatem(size_t size){
    node_t * actual = first;
    while (actual->next != NULL)
        actual = (node_t *) actual->next;

    node_t * next = (actual + size);
    init_node(next);

    actual->size = size;
    actual->empty = 0;
    actual->next = (struct node *) next;

    return actual + 1;
}

void freem(void * node_to_remove_pointer){
    size_t size = sizeof(node_t);

    node_t * node_to_remove = (node_t *) (node_to_remove_pointer - size);
    node_t * prev_node_to_node_to_remove = node_to_remove - size;
    node_t * next_node_to_node_to_remove = (node_t *) node_to_remove->next;

    if(prev_node_to_node_to_remove != NULL && next_node_to_node_to_remove != NULL){ //Between two nodes
        prev_node_to_node_to_remove->empty = 0;
        prev_node_to_node_to_remove->next = (struct node *) next_node_to_node_to_remove;
        prev_node_to_node_to_remove->size = node_to_remove->size;

    }else if(prev_node_to_node_to_remove != NULL){ //Last node
        prev_node_to_node_to_remove->next = NULL;
        prev_node_to_node_to_remove->size = -1;
        prev_node_to_node_to_remove->empty = 1;

    }else if(next_node_to_node_to_remove != NULL){ //First node
        first = next_node_to_node_to_remove;
        first->empty = 0;
        first->size = next_node_to_node_to_remove->size;

    }else { //Only one node
        first = NULL;
    }
}

static void init_node(node_t * node){
    node->empty = 1;
    node->size = 0;
    node->next = NULL;
}