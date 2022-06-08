#include <stdio.h>

void swap(int *x, int *y);

int main() {
    int x = 5;
    int y = 10;

    swap(&x, &y);

    printf("&x -> %i\n", &x);
    printf("&y -> %i\n", &y);
    printf("x -> %i\n", x);
    printf("y -> %i", y);

    return 0;
}

void swap(int *x, int *y) {
//    int aux = *x;
//    *x = *y;
//    *y = aux;

    printf("%i\n", x);
    printf("%i\n", y);
    printf("*x -> %i\n", *x);
    printf("*y -> %i\n", *y);
}
