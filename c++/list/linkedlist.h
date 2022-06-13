#pragma once

#include <iostream>

using String = std::string;

template<typename T>
class Node {
    public: Node * next;
    public: Node * back;
    public: T value;

    public: Node(Node * next, Node * back, T value):
            next {next}, back {back}, value {value}
    {}

    ~Node(){
        delete value;
    }
};

template<typename T>
class Linkedlist {
    private: Node<T> * first;
    private: Node<T> * last;
    private: int size;

    public: Linkedlist():
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

    public: T operator[](int index) const{
        return getNode(index)->value;
    }

    public: Linkedlist * add(const T value){
        if(this->size == 0){
            Node<T> * newNode = new Node<T>(nullptr, nullptr, value);
            this->first = newNode;
            this->last = newNode;
            this->size = 1;

        }else{
            Node<T> * newNode = new Node<T>(nullptr, this->last, value);
            this->last->next = newNode;
            this->last = newNode;
            this->size = this->size + 1;
        }

        return this;
    }

    public: bool remove(int index){
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

        delete nodeToDelete;
        this->size = this->size - 1;

        return true;
    }

    public: bool isEmpty(){
        return this->size == 0;
    }

    public: int getSize(){
        return this->size;
    }

    public: void clear(){
        this->size = 0;
        Node<T> * actualNode = this->first;

        while (actualNode != nullptr){
            Node<T> * nextNodeAux = actualNode->next;
            delete actualNode;
            actualNode = nextNodeAux->next;
        }
    }

    public: T findBy(bool (* predicate)(T value)){
        for(Node<T> * actualNode = this->first; actualNode != nullptr; actualNode = actualNode->next){
            if(predicate(actualNode->value))
                return actualNode->value;
        }

        return nullptr;
    }

    public: T get(int requiredIndex){
        auto actualIndex = -1;

        if(requiredIndex < 0 || requiredIndex + 1 >= size)
            throw std::out_of_range("Item in list out of bounds");

        for(Node<T> * actualNode = this->first; actualNode != nullptr; actualNode = actualNode->next){
            actualIndex++;
            if(actualIndex == requiredIndex)
                return actualNode->value;
        }
    }

    public: int indexOf(T * value){
        int actualIndex = -1;

        for(Node<T> * actualNode = this->first; actualNode != nullptr; actualNode = actualNode->next){
            actualIndex++;
            if(actualNode->value.compare(value) == 0)
                break;
        }

        return actualIndex;
    }

    private: Node<T> * getNode(int requiredIndex){
        auto actualIndex = -1;

        for(Node<T> * actualNode = this->first; actualNode != nullptr; actualNode = actualNode->next){
            actualIndex++;
            if(actualIndex == requiredIndex)
                return actualNode;
        }

        return nullptr;
    }
};