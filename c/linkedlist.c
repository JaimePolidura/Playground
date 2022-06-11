#include <stdio.h>
#include <malloc.h>
#include <string.h>

struct node {
    int value;
    struct node *next;
    struct node *back;
};
typedef struct node node_t;

struct linkedlist {
    int size;
    struct node *last;
};
typedef struct linkedlist linkedlist_t;

void create_linkedlist(linkedlist_t *linkedlist, int initialValue);
void add_new_node_linkedlist(linkedlist_t *linkedlist, int newValue);
node_t* get_node(linkedlist_t *linkedlist, int index);
int delete_node(linkedlist_t *linkedlist, int index);

enum boolean {TRUE, FALSE};

int main(){
    linkedlist_t list;

    create_linkedlist(&list, 1);
    printf("The size of the linkeslist: %i\n", list.size);
    printf("The last node: %i\n\n", list.last->value);

    add_new_node_linkedlist(&list, 2);
    printf("The size of the linkeslist: %i\n", list.size);
    printf("The last node: %i\n\n", list.last->value);

    node_t *node1 = get_node(&list, 0);
    printf("Node foudn in %i: %i\n", 0, node1->value);
    node_t *node2 = get_node(&list, 3);
    printf("Node foudn in %i: %i\n", 3, node2->value);

    return 0;
}

int delete_node(linkedlist_t *linkedlist, int index) {
    node_t *nodeToDelete = get_node(linkedlist, index);

    if(nodeToDelete == NULL)
        return -1;

    if(nodeToDelete->next == NULL && nodeToDelete->back == NULL){ //delete first
        linkedlist->last = NULL;

    }else if(nodeToDelete->next == NULL){
        nodeToDelete->back->next = NULL;
        linkedlist->last = nodeToDelete->back;

    }else if(nodeToDelete->back == NULL){
        nodeToDelete->next->back = NULL;

    }else{
        node_t *nextNodeToDelete = nodeToDelete->next;
        node_t *backNodeToDelete = nodeToDelete->back;

        nextNodeToDelete->back = backNodeToDelete;
        backNodeToDelete->next = nextNodeToDelete;
    }

    linkedlist->size = linkedlist->size - 1;
    free(nodeToDelete);

    return 1;
}

node_t* get_node(linkedlist_t *linkedlist, int index){
    if(index < 0 || index > linkedlist->size)
        return NULL;

    node_t *actual = linkedlist->last;
    int actualIndex = 0;
    while (actual != NULL){
        if(index == actualIndex)
            return actual;

        actual = actual->next;
        actualIndex++;
    }

    return NULL;
}

void add_new_node_linkedlist(linkedlist_t *linkedlist, int newValue){
    node_t *oldLastNode = linkedlist->last;

    node_t *newNode = malloc(sizeof(node_t));
    newNode->value = newValue;
    newNode->back = oldLastNode;
    newNode->next = NULL;

    oldLastNode->next = newNode;

    linkedlist->last = newNode;
    linkedlist->size = linkedlist->size + 1;
}

void create_linkedlist(linkedlist_t *linkedlist, int initialValue){
    node_t *newNode = malloc(sizeof(node_t));
    newNode->back = NULL;
    newNode->next = NULL;
    newNode->value = initialValue;

    linkedlist->last = newNode;
    linkedlist->size = 1;
}