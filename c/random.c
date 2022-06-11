#include <stdio.h>
#include <stdlib.h>

int main(){
    int * p = malloc(sizeof(int));
    * p = 2;

    printf("El valor de p %i\n", p);
    printf("El valor de p %i\n", *p);
    printf("El valor de p %i\n", &p);

    return 0;
}