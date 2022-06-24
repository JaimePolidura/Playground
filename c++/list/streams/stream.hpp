#include <cstdlib>
#include <iostream>

template<typename T>
class Stream {
private:
    T * pointer;
    size_t size;

public:
    Stream(T * pointer, size_t size): pointer{pointer}, size{size} {}
};

template<typename T>
static Stream<T> streamOf (Iterator<T> * iterator){
    size_t size = iterator->size();
    T * newMemPointer = new T[size];

    T& actual = iterator->next();
    int actualCount = 0;
    while (iterator->hasNext()) {
        actual = iterator->next();
        newMemPointer[actualCount++] = actual;
    }

    return Stream<T>{newMemPointer, size};
}