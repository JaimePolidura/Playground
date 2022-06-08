#include <stdio.h>

int length(char * string){
    int n = 0;

    while(*string++ != '\0')
        n++;

    return n;
}

int main(){
    int variable = 12;
    printf("%i\n", &variable);

    printf("\n%i\n", length("hola adios"));

    return 0;
}