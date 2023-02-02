#include <iostream>
#include <filesystem>

#include "./list/linkedlist.h"
#include "./list/streams/stream.hpp"

void printItAsInt(int& element) {
    printf("%i ", element);
}

int multiplyByTwo(int& element) {
    return element * 2;
}

bool isEven(int& num) {
    return num % 2 == 0;
}


int main(){
    auto linkedlistTestStreams = new Linkedlist<int>();
    linkedlistTestStreams->add(1)->add(2)->add(3)->add(3)->add(4)->add(5)
                ->add(6)->add(7)->add(8)->add(9)->add(10);

    auto size = linkedlistTestStreams->getSize();

    Iterator<int> * iterator = linkedlistTestStreams->iterator();
    Stream<int> stream = streamOf(iterator)
            .filter(isEven) // 2 4 6 8 10
            .map(multiplyByTwo) // 4 8 12 16 20
            .forEach(printItAsInt);

    return 0;
}