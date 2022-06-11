//
// Created by polid on 10/06/2022.
//

#include "linkedlist.h"
#include <iostream>

using String = std::string;

class Node {
    public: Node * next;
    public: Node * back;
    public: String value;

    public: Node(Node * next, Node * back, String value):
            next {next}, back {back}, value {value}
    {}
};

class Linkedlist {
    private: Node * first;
    private: Node * last;
    private: int size;

    public: Linkedlist():
        first {nullptr}, last {nullptr}, size{0}
    {}

    public: Linkedlist * add(const String value){
        if(this->size == 0){
            Node * newNode = new Node(nullptr, nullptr, value);
            this->first = newNode;
            this->last = newNode;
            this->size = 1;

        }else{
            Node * newNode = new Node(nullptr, this->last, value);
            this->last->next = newNode;
            this->size = ++this->size;
        }

        return this;
    }

    public: bool remove(int index){
        Node * nodeToDelete = this->getNode(index);

        if(nodeToDelete == nullptr) return false;

        if(this->size == 1){ //Remove the only element
            this->first = nullptr;
            this->last = nullptr;

        }else if(index == 0){ //First element
            Node * nextNodeToFirst = this->first->next;
            nextNodeToFirst->back = nullptr;
            this->first = nextNodeToFirst;
        }else if(index == size - 1){ //Last element
            Node * backNodeToLast = this->first->back;
            backNodeToLast->next = nullptr;
            this->last = backNodeToLast;
        }else { //Node between other nodes
            Node * nextNodeToNodeToRemove = nodeToDelete->next;
            Node * backNodeToNodeToRemove = nodeToDelete->back;
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
        Node * actualNode = this->first;

        while (actualNode != nullptr){
            Node * nextNodeAux = actualNode->next;
            delete actualNode;
            actualNode = nextNodeAux->next;
        }
    }

    public: String findBy(bool (* predicate)(String value)){
        for(Node * actualNode = this->first; actualNode != nullptr; actualNode = actualNode->next){
            if(predicate(actualNode->value))
                return actualNode->value;
        }

        return nullptr;
    }

    public: String get(int requiredIndex){
        if(requiredIndex < 0 || requiredIndex + 1 >= size)
            return nullptr;

        auto actualIndex = -1;

        for(Node * actualNode = this->first; actualNode != nullptr; actualNode = actualNode->next){
            actualIndex++;
            if(actualIndex == requiredIndex)
                return actualNode->value;
        }

        return nullptr;
    }

    public: int indexOf(String value){
        int actualIndex = -1;

        for(Node * actualNode = this->first; actualNode != nullptr; actualNode = actualNode->next){
            actualIndex++;
            if(actualNode->value.compare(value))
                break;
        }

        return actualIndex;
    }

    private: Node * getNode(int requiredIndex){
        auto actualIndex = 0;

        for(Node * actualNode = this->first; actualNode != nullptr; actualNode = actualNode->next){
            actualIndex++;
            if(actualIndex == requiredIndex)
                return actualNode;
        }

        return nullptr;
    }
};

int main(){
    Linkedlist * list = new Linkedlist();

    list->add("jaime")
            ->add("pedro")
            ->add("paula")
            ->add("javier");

    auto text = list->get(2);


    printf("%s", text.c_str());
}