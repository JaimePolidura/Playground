package es.jaime.trees;


import java.util.ArrayList;

public class MinHeap<E extends Comparable<E>> {
    private final ArrayList<E> elements;

    public MinHeap () {
        this.elements = new ArrayList<>();
    }

    public boolean isEmpty () {
        return elements.size() == 0;
    }

    public void add(E element) {
        elements.add(element);
        if(elements.size() == 1){
            return;
        }

        addRecursive(element, elements.size() - 1);
    }

    private void addRecursive(E element, int index) {
        int indexParent = getIndexOfParent(index);

        if((elements.size() - 1) <= indexParent){
            if(element.compareTo(elements.get(0)) < 0){
                swapElements(0, index);
            }
        }else{
            E parent = elements.get(indexParent);
            if(element.compareTo(parent) < 0){
                swapElements(index, indexParent);
                addRecursive(element, indexParent);
            }
        }
    }

    private int getIndexLefttChild(int index) {
        return 2 * index + 1;
    }

    private int getIndexRightChild(int index) {
        return 2 * index + 2;
    }

    private int getIndexOfParent (int index) {
        return (index - 1) / 2;
    }

    private void swapElements (int a, int b) {
        E aElement = elements.get(a);
        E bElement = elements.get(b);

        elements.set(a, bElement);
        elements.set(b, aElement);
    }

    @Override
    public String toString () {
        return elements.toString();
    }
}
