package es.jaime.list;

import lombok.AllArgsConstructor;
import lombok.Getter;
import lombok.experimental.Accessors;

import java.util.*;
import java.util.function.Predicate;

public class LinkedList<E> implements Iterable<E> {
    @Getter private final int maxCapacity;
    private Node<E> head;
    private Node<E> tail;
    @Getter private int size;

    public LinkedList(int maxCapacity) {
        this.maxCapacity = maxCapacity;
    }

    public LinkedList(Collection<? extends E> collection) {
        this.maxCapacity = -1;

        collection.forEach(this::add);
    }

    public LinkedList() {
        this.maxCapacity = -1;
    }

    public E getFirstElement(){
        return head.element;
    }

    public E getLastElement(){
        return tail.element;
    }

    public boolean isEmpty () {
        return head == null;
    }

    public void clear(){
        head = null;
        tail = null;
        size = 0;
    }

    public E addFirst (E element) {
        checkIfNotNullAndNotReachedMaxCap(element);

        Node<E> toAdd = new Node<>(head, null, element);

        if(isEmpty()){
            head = toAdd;
            tail = toAdd;
        }else{
            head.prev = toAdd;
            head = toAdd;
        }

        size++;

        return toAdd.element;
    }

    public E add (E element){
        checkIfNotNullAndNotReachedMaxCap(element);

        if(isEmpty()){
            return addFirst(element);
        }else{
            Node<E> toAdd = new Node<>(null, tail, element);
            tail.next = toAdd;
            this.tail = toAdd;

            size++;

            return element;
        }
    }

    public E add (int index, E element) {
        checkIfNotNullAndNotReachedMaxCap(element);
        checkIfIndexInsideCollectionToAdd(index);

        if (index + 1 == size){
            return add(element);
        }else if(index == 0 || isEmpty()){
            return addFirst(element);
        }else{
            Node<E> alreadyInIndex = this.getNodeAt(index);
            Node<E> alreadyInIndexPrev = alreadyInIndex.prev;
            Node<E> toAdd = new Node<>(alreadyInIndex, alreadyInIndexPrev, element);

            alreadyInIndexPrev.next = toAdd;
            alreadyInIndex.prev = toAdd;

            size++;

            return toAdd.element;
        }
    }

    public E get(int index) {
        checkIfIndexInsideCollection(index);

        if(index + 1 == size){
            return tail.element;
        }else if (index == 0) {
            return head.element;
        }

        return getNodeAt(index).element;
    }

    public boolean removeFirst(){
        if(isEmpty()) return false;

        Node<E> nextToHead = head.next;
        this.head = nextToHead;
        nextToHead.prev = null;
        size--;

        return true;
    }

    public boolean removeLast(){
        if(isEmpty()) return false;

        Node<E> prevToTail = tail.prev;
        this.tail = prevToTail;
        prevToTail.next = null;
        size--;

        return true;
    }

    public boolean removeElement (E element){
        Objects.requireNonNull(element);

        Node<E> nodeToRemove = getNodeOfElement(element);
        if(nodeToRemove == null || isEmpty()){
            return false;
        }else if (nodeToRemove == head) {
            return removeFirst();
        }else if (nodeToRemove == tail){
            return removeLast();
        }

        unlinkNode(nodeToRemove);
        return true;
    }

    public boolean removeAt (int index) {
        checkIfIndexInsideCollection(index);

        if(isEmpty()){
            return false;
        }else if(index == 0){
            return removeFirst();
        }else if(index == size - 1){
            return removeLast();
        }else{
            Node<E> nodeToRemove = getNodeAt(index);
            unlinkNode(nodeToRemove);
        }

        return true;
    }

    public int removeAllElements (E element) {
        Objects.requireNonNull(head);

        Node<E> iterate = head;
        int itemsRemoved = 0;

        while (iterate != null){
            if(iterate.element == element){
                if(iterate == head){
                    removeFirst();
                }else if (iterate == tail) {
                    removeLast();
                }else{
                   unlinkNode(iterate);
                }
                itemsRemoved++;
            }
            iterate = iterate.next;
        }

        return itemsRemoved;
    }

    public LinkedList<E> removeRange(int beginIndex, int endIndex) {
        checkIfRangeInsideCollection(beginIndex, endIndex);

        if(endIndex - beginIndex == 0) {
            clear();
            return this;
        }

        Node<E> nodetToStart = getNodeAt(beginIndex);
        Node<E> nodetToStartPrev = nodetToStart.prev;
        Node<E> nodeToEnd = getNodeAt(endIndex);
        Node<E> nodeToEndNext = nodeToEnd.next;

        nodetToStartPrev.next = nodeToEndNext;
        nodeToEndNext.prev = nodetToStartPrev;
        size = size - (endIndex - beginIndex);

        return this;
    }

    public E modify (int index, E element) {
        Objects.requireNonNull(element);

        Node<E> nodeToModify = getNodeAt(index);

        if(nodeToModify == null){
            return null;
        }else{
            return nodeToModify.element = element;
        }
    }

    public int indexOf (E element) {
        int index = 0;

        for (Node<E> iterate = head; iterate != null; iterate = iterate.next) {
            if (iterate.element.equals(element))
                return index;
            index = index + 1;
        }

        return -1;
    }

    public ArrayList<Integer> indexOfElements(E element){
        Node<E> iterate = head;
        int index = 0;
        ArrayList<Integer> toReturn = new ArrayList<>();

        while (iterate != null){
            if(iterate.element == element){
                toReturn.add(index);
            }
            index++;
            iterate = iterate.next;
        }

        return toReturn;
    }

    public boolean contains (E element) {
        return indexOf(element) != -1;
    }

    public LinkedList<E> retainAll (Collection<? extends E> collection){
        Node<E> iterate = head;

        while (iterate != null) {
            if (!collection.contains(iterate.element)) {
                removeAllElements(iterate.element);
            }
            iterate = iterate.next;
        }

        return this;
    }

    public LinkedList<E> addAll(Collection<? extends E> collection){
        collection.forEach(this::add);

        return this;
    }

    public LinkedList<E> addAll(int index, Collection<? extends E> collection){
        checkIfIndexInsideCollectionToAdd(index);

        for (E element : collection) {
            add(index, element);
            index++;
        }

        return this;
    }

    public LinkedList<E> reverse () {
        Node<E> iterate = tail;
        clear();

        while (iterate != null){
            add(iterate.element);
            iterate = iterate.prev;
        }

        return this;
    }

    public LinkedList<E> subList (int beginIndex, int finalIndex) {
        checkIfRangeInsideCollection(beginIndex, finalIndex);

        LinkedList<E> subList = new LinkedList<>();
        Node<E> nodeIte = getNodeAt(beginIndex);

        for(int i = beginIndex; i < finalIndex + 1; i++){
            subList.add(nodeIte.element);
            nodeIte = nodeIte.next;
        }

        return subList;
    }

    public boolean removeIf (Predicate<? super E> condition) {
        Objects.requireNonNull(condition);

        boolean removed = false;

        Node<E> actualNode = head;

        int index = 0;
        while (actualNode != null){
            if(condition.test(actualNode.element)){
                removeAt(index);
                removed = true;
            }else{
                index++;
            }
            actualNode = actualNode.next;
        }

        return removed;
    }

    public ArrayList<E> toList () {
        ArrayList<E> list = new ArrayList<>();

        Node<E> iterate = head;
        while (iterate != null){
            list.add(iterate.element);
            iterate = iterate.next;
        }

        return list;
    }

    public void sort (Comparator<? super E> comparator) {
        Objects.requireNonNull(comparator);
        ArrayList<E> editableList = toList();
        clear();

        while (editableList.size() != 0) {
            E fixedElement = editableList.get(0);
            ArrayList<E> subListOfLargestThanFixed = new ArrayList<>();
            int numberOfZeros = 0;

            for (E elementToCompare : editableList) {
                int resultToCompare = comparator.compare(elementToCompare, fixedElement);

                if (resultToCompare == -1)
                    subListOfLargestThanFixed.add(elementToCompare);
                else if (resultToCompare == 0)
                    numberOfZeros++;
            }

            if (subListOfLargestThanFixed.size() == 0) {
                for (int k = 0; k < numberOfZeros; k++) {
                    add(fixedElement);
                    editableList.remove(fixedElement);
                }
            } else {
                E largestElementOfSubList = getTheLargestElementFromAnArrayList(editableList, comparator);

                add(largestElementOfSubList);
                editableList.removeIf(ele -> ele.equals(largestElementOfSubList));
            }
        }
    }

    private E getTheLargestElementFromAnArrayList (ArrayList<E> list, Comparator<? super E> comparator) {
        E largest = list.get(0);

        for (E element : list) {
            if(comparator.compare(largest, element) == -1){
                largest = element;
            }
        }

        return largest;
    }

    private void unlinkNode (Node<E> nodeToUnlink) {
        Node<E> nodeUnlinkPrev = nodeToUnlink.prev;
        Node<E> nodeUnlinkNext = nodeToUnlink.next;

        nodeUnlinkPrev.next = nodeUnlinkNext;
        nodeUnlinkNext.prev = nodeUnlinkPrev;

        size--;
    }
    
    private Node<E> getNodeAt (int index) {
        int middlePoint = (int) Math.round((double) (size) / 2);

        if(index + 1 < middlePoint){
            return getNodeIterateToward(index);
        }else{
            return getNodeIterateBackward(index);
        }
    }

    private Node<E> getNodeIterateBackward (int index) {
        int middlePoint = (int) Math.round((double) (size) / 2);
        Node<E> ite = this.tail;

        for(int i = size; i >= middlePoint; i--){
            if(index + 1 == i){
                return ite;
            }
            ite = ite.prev;
        }

        return null;
    }

    private Node<E> getNodeIterateToward (int index) {
        int middlePoint = (int) Math.round((double) (size) / 2);
        Node<E> ite = this.head;

        for(int i = 1; i <= middlePoint; i++){
            if(index + 1 == i){
                return ite;
            }
            ite = ite.next;
        }

        return null;
    }

    private Node<E> getNodeOfElement(E element) {
        Node<E> ite = head;

        while (ite != null) {
            if (ite.element == element) {
                return ite;
            }
            ite = ite.next;
        }

        return null;
    }

    private void checkIfNotNullAndNotReachedMaxCap (E element) {
        Objects.requireNonNull(element);

        if(maxCapacity != -1 && size + 1 >= maxCapacity){
             throw new RuntimeException("The linked list has reached the max capacity");
        }
    }

    private void checkIfRangeInsideCollection(int beginIndex, int finalIndex) {
        if(beginIndex > finalIndex) throw new IllegalArgumentException("Final index cannot be bigger than beginIndex");
        if(beginIndex < 0 || finalIndex >= size) throw new ArrayIndexOutOfBoundsException("The range of begin index and finalIndex must be on the list range");
    }

    private void checkIfIndexInsideCollectionToAdd(int index) {
        if(index < 0 || index > size) throw new ArrayIndexOutOfBoundsException("The index must be inside the list");
    }

    private void checkIfIndexInsideCollection(int index) {
        if(index < 0 || index > size) throw new ArrayIndexOutOfBoundsException("The index must be inside the list");
    }

    @Override
    public String toString() {
        StringBuilder result = new StringBuilder();
        Node<E> currentNode = head;
        
        while (currentNode != null) {
            result.append(currentNode.element)
                    .append(currentNode != tail ? " <-> " : "");

            currentNode = currentNode.next;
        }

        return result.toString();
    }

    @Override
    public Iterator<E> iterator() {
        return new MyIterator<>();
    }

    private class MyIterator<DoubleLinkedList> implements Iterator<E>{
        private int index = -1;

        @Override
        public boolean hasNext() {
            return getNodeAt(index + 1).next != null;
        }

        @Override
        public E next() {
            return getNodeAt(++index).element;
        }

        @Override
        public void remove () {
            removeAt(index);
        }
    }

    @AllArgsConstructor
    private static class Node<E> {
        Node<E> next;
        Node<E> prev;
        E element;

        @Override
        public String toString () {
            return String.valueOf(element);
        }
    }
}
