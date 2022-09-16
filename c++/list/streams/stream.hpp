#include <cstdlib>
#include <iostream>
#include <functional>

template<typename T>
class Stream {
private:
    T * pointer;
    int size;

public:
    Stream(T * pointer, const int& size): pointer{pointer}, size{size} {}

    Stream<T> filter (bool (* predicate)(T& element)) {
        int lastIndexOfMatch = 0;

        for(int i = 0; i < this->size; i++){
            T * actualValue = this->pointer + i;
            bool matches = predicate(* actualValue);

            if(matches){
                * (pointer + lastIndexOfMatch) = * actualValue;
                lastIndexOfMatch++;
            }
        }

        return Stream{pointer, lastIndexOfMatch};
    }

    template<typename O>
    Stream<O> map(O (* mapper)(T& element)) {
        O * result = new O[this->size];

        for(int i = 0; i < this->size; i++){
            T * actualValueToMap = this->pointer + i;
            O mappedValue = mapper(* actualValueToMap);

            * (result + i) = mappedValue;
        }

        delete[] this->pointer;

        return Stream<O>{result, this->size};
    }
    
    T& get(int position) {
        if(position < 0 || position > this->size - 1)
            throw std::out_of_range("Out of range error");

        return this->pointer[position];
    }

    int count() {
       return this->size;
    }

    Stream<T> forEach(void (* consumer)(T& element)) {
        for(int i = 0; i < this->size; i++)
           consumer(this->pointer[i]);
    }
};

template<typename T>
static Stream<T> streamOf (Iterator<T> * iterator){
    size_t size = iterator->size();
    T * newMemPointer = new T[size];

    int actualCount = 0;
    while (iterator->hasNext()) {
        const T& actual = iterator->next();
        newMemPointer[actualCount++] = actual;
    }

    return Stream<T>{newMemPointer, size};
}