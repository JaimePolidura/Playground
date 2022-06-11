#include <stdio.h>
#include <malloc.h>

char * substring(char *, int, int);

struct hola {
    long hola;
    char adios;
};

int main(){
    printf("%c", substring("hola soy jaime", 1, 3));

    return 1;
}


//Funciona xd
char * substring(char * source, int from, int to){
    char * substring = (char *) malloc(sizeof(char) * (to - from));
    source += from;
    int count = from;

    while(*source++ != '\0' && count < to){
        *(substring + count) = *source;
        count++;
    }

    return substring;
}