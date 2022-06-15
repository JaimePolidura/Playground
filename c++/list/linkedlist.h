#pragma once

#include <iostream>

using String = std::string;

template<typename T>
class Node {
    public: Node * next;
    public: Node * back;
    public: T& value;

    public: Node(Node * next, Node * back, T& value):
            next {next}, back {back}, value {value}
    {}
};

template<typename T>
class Linkedlist {
private:
    Node<T> * first;
    Node<T> * last;
    int size;

public:
    Linkedlist():
        first {nullptr}, last {nullptr}, size{0}
    {}

    ~Linkedlist(){
        Node<T> * actual = this->first;

        while (actual != nullptr){
            Node<T> * nextToActual = actual->next;
            delete actual;
            actual = nextToActual;
        }
    }

    T& operator[](int index) const{
        return getNode(index)->value;
    }

    Linkedlist * add(T& value);
    bool remove(int index);
    bool isEmpty();
    int getSize();
    void clear();
    T& findBy(bool (* predicate)(T& value));
    T& get(int requiredIndex);
    int indexOf(T& value);

    private: Node<T> * getNode(int requiredIndex){
        auto actualIndex = -1;

        for(Node<T> * actualNode = this->first; actualNode != nullptr; actualNode = actualNode->next){
            actualIndex++;
            if(actualIndex == requiredIndex)
                return actualNode;
        }

        throw std::logic_error("index not found");
    }
};

template<typename T>
Linkedlist<T> * Linkedlist<T>::add(T& value){
    if(this->size == 0){
        auto * newNode = new Node<T>(nullptr, nullptr, value);
        this->first = newNode;
        this->last = newNode;
        this->size = 1;

    }else{
        auto * newNode = new Node<T>(nullptr, this->last, value);
        this->last->next = newNode;
        this->last = newNode;
        this->size = this->size + 1;
    }

    return this;
}

template<typename T>
bool Linkedlist<T>::remove(int index){
    Node<T> * nodeToDelete = this->getNode(index);

    if(nodeToDelete == nullptr) return false;

    if(this->size == 1){ //Remove the only element
        this->first = nullptr;
        this->last = nullptr;

    }else if(index == 0){ //First element
        Node<T> * nextNodeToFirst = this->first->next;
        nextNodeToFirst->back = nullptr;
        this->first = nextNodeToFirst;
    }else if(index == size - 1){ //Last element
        Node<T> * backNodeToLast = this->first->back;
        backNodeToLast->next = nullptr;
        this->last = backNodeToLast;
    }else { //Node between other nodes
        Node<T> * nextNodeToNodeToRemove = nodeToDelete->next;
        Node<T> * backNodeToNodeToRemove = nodeToDelete->back;
        nextNodeToNodeToRemove->back = backNodeToNodeToRemove;
        backNodeToNodeToRemove->next = nextNodeToNodeToRemove;
    }

    free(nodeToDelete);
    this->size = this->size - 1;

    return true;
}

template<typename T>
bool Linkedlist<T>::isEmpty(){
    return this->size == 0;
}

template<typename T>
int Linkedlist<T>::getSize(){
    return this->size;
}

template<typename T>
void Linkedlist<T>::clear(){
    this->size = 0;
    Node<T> * actualNode = this->first;

    while (actualNode != nullptr){
        Node<T> * nextNodeAux = actualNode->next;
        delete actualNode;
        actualNode = nextNodeAux->next;
    }
}

template<typename T>
T& Linkedlist<T>::findBy(bool (* predicate)(T& value)){
    for(Node<T> * actualNode = this->first; actualNode != nullptr; actualNode = actualNode->next){
        if(predicate(actualNode->value))
            return actualNode->value;
    }

    throw std::logic_error("index not found");
}

template<typename T>
T& Linkedlist<T>::get(int requiredIndex){
    auto actualIndex = -1;

    if(requiredIndex < 0 || requiredIndex + 1 > size)
        throw std::out_of_range("Item in list out of bounds");

    for(Node<T> * actualNode = this->first; actualNode != nullptr; actualNode = actualNode->next){
        actualIndex++;
        if(actualIndex == requiredIndex)
            return actualNode->value;
    }

    throw std::logic_error("index not found");
}

template<typename T>
int Linkedlist<T>::indexOf(T& value){
    int actualIndex = -1;

    for(Node<T> * actualNode = this->first; actualNode != nullptr; actualNode = actualNode->next){
        actualIndex++;
        if(actualNode->value.compare(value) == 0)
            break;
    }

    return actualIndex;
}