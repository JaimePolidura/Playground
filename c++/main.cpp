#include <iostream>

struct Pair {
    char * key;
    int value;
};

typedef struct Pair pair_t;

class Jugador {
    char * nombre;
    double dinero;

    public: Jugador(char * nombre, double dinero){
        this->nombre = nombre;
        this->dinero = dinero;
    }

    Jugador(const Jugador& jugador_to_copy): dinero{jugador_to_copy.dinero}, nombre{new char[10]} {
    }

    double getDinero() const{
        return this->dinero;
    }

    char * getNombre() const{
        return this->nombre;
    }

    Jugador& operator=(const Jugador& other){
        if(this == &other) return * this;

        const char * newNombre = new char[16];
        double newDinero = other.dinero;

        return * this;
    }

    ~Jugador(){
        printf("F\n");
    }
};

void printRef(int& data){
    printf("%i\n", data);
}

//int main() {
//    auto the_answer { 42 };
//
//    auto* jugador = new Jugador("jhaime", 1);
//
//    int p{};
//    printf("%i\n", p);
//
//    int original = 100;
//    printRef(original);
//
//    throw std::runtime_error("error");
//}