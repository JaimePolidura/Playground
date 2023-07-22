package es.jaime.list;

import lombok.AllArgsConstructor;
import lombok.Getter;

import java.util.Comparator;

// head -- tail
public final class SkipList<T> implements List<T> {
    private static final int MAX_LEVEL = 2;

    private final Comparator<T> comparator;
    private Node<T> head;
    private Node<T> tail;
    private int size;

    public SkipList(Comparator<T> comparator) {
        this.comparator = comparator;
        this.head = Node.of(null);
        this.tail = Node.of(null);

        for(int actualLevel = MAX_LEVEL; actualLevel >= 0; actualLevel--) {
            this.head.levels[actualLevel].next = this.tail;
            this.tail.levels[actualLevel].back = this.head;
        }
    }

    @Override
    public boolean add(T item) {
        if(this.size == 0){
            addFirst(item);
            this.size = this.size + 1;
            return true;
        }

        Node<T>[][] nodesPredSuc = new Node[MAX_LEVEL + 1][2];
        for(int actualLevel = 0; actualLevel <= MAX_LEVEL; actualLevel++){
            Node<T> predecesor = this.head;
            Node<T> successor = this.head.levels[0].next;

            do {
                predecesor = predecesor.levels[0].next;
                successor = successor.levels[0].next;

                nodesPredSuc[actualLevel][0] = predecesor;
                nodesPredSuc[actualLevel][1] = successor;

                int comparationWithPredecesor = this.comparator.compare(item, predecesor.data);

                if(predecesor.data.equals(item) || (successor != null && successor.data != null && successor.data.equals(item)))
                    return false;
                if(successor == null || successor == this.tail) {
                    nodesPredSuc[actualLevel][0] = comparationWithPredecesor > 0 ? predecesor : predecesor.levels[0].back;
                    nodesPredSuc[actualLevel][1] = comparationWithPredecesor > 0 ? successor : predecesor;
                    break;
                }

                int comparationWithSuccessor = this.comparator.compare(item, successor.data);

                if(comparationWithPredecesor < 0 && comparationWithSuccessor > 0)
                    break;

            } while (predecesor.levels[actualLevel] != null);
        }

        Node<T> newNode = Node.of(item);

        for(int actualLevel = 0; actualLevel <= MAX_LEVEL; actualLevel++){
            Node<T> predecesor = nodesPredSuc[actualLevel][0];
            Node<T> successor = nodesPredSuc[actualLevel][1];
            boolean randomTest = this.randomTest50();

            if(actualLevel == 0 || randomTest){
                predecesor.levels[actualLevel].next = newNode;
                successor.levels[actualLevel].back = newNode;
                newNode.levels[actualLevel].back = predecesor;
                newNode.levels[actualLevel].next = successor;
            }

            if(actualLevel != 0 && !randomTest)
                break;
        }

        this.size = this.size + 1;

        return true;
    }

    private void addFirst(T item) {
        Node<T> newNode = Node.of(item);
        boolean lastLevelAdded = true;

        for(int actualLevel = 0; actualLevel <= MAX_LEVEL && lastLevelAdded; actualLevel++) {
            if(actualLevel == 0 || randomTest50()){
                this.head.levels[actualLevel].next = newNode;
                this.tail.levels[actualLevel].back = newNode;
                newNode.levels[actualLevel].next = this.tail;
                newNode.levels[actualLevel].back = this.head;
            }else {
                lastLevelAdded = false;
            }
        }
    }

    @Override
    public T get(int requiredIndex) {
        Node<T> actualNode = this.head.levels[0].next;
        int actualIndex = 0;

        if(actualNode == null)
            return null;

        while (actualNode != null) {
            if(actualIndex == requiredIndex)
                return actualNode.data;

            actualNode = actualNode.levels[0].next;
            actualIndex = actualIndex + 1;
        }

        return null;
    }

    @Override
    public boolean contains(T item) {
        if(this.size == 0)
            return false;
        if(this.size == 1)
            return this.head.levels[0].next.data.equals(item);

        int actualLevelIndex = MAX_LEVEL;
        Node<T> actualNode = this.head.levels[MAX_LEVEL].next;

        while (actualLevelIndex != 0) {
            if(actualNode != this.tail && actualNode.data.equals(item))
                return true;

            int comparationWithItem = actualNode != this.tail ? this.comparator.compare(item, actualNode.data) : Integer.MIN_VALUE;
            actualLevelIndex = actualLevelIndex - 1;

            if(comparationWithItem > 0){ //item bigger than actual -> go right (next)
                if(actualLevelIndex < 0)
                    break;

                do {
                    if(actualNode.levels[actualLevelIndex].next != null)
                        actualNode = actualNode.levels[actualLevelIndex].next;
                    else
                        actualLevelIndex--;
                }while (actualLevelIndex >= 0 && actualNode.levels[actualLevelIndex].next == null);

            }else if(comparationWithItem < 0) { //item smaller than actual <- go left (back)
                if(actualLevelIndex < 0)
                    break;

                do {
                    if(actualNode.levels[actualLevelIndex].back != null)
                        actualNode = actualNode.levels[actualLevelIndex].back;
                    else
                        actualLevelIndex--;
                }while (actualLevelIndex >= 0 && actualNode.levels[actualLevelIndex].back == null);

            }else if(comparationWithItem == 0) {
                return true;
            }
        }

        if(actualNode.data != null && actualNode.data.equals(item))
            return true;

        if(actualNode == this.head || actualNode == this.tail){
            return actualNode == this.head ?
                    iterateRightUntilFound(item, actualNode.levels[0].next) :
                    iterateLeftUntilFound(item,  actualNode.levels[0].back);
        }else {
            return this.comparator.compare(item, actualNode.data) > 0 ?
                    iterateRightUntilFound(item, actualNode) :
                    iterateLeftUntilFound(item, actualNode);
        }
    }

    private boolean iterateLeftUntilFound(T item, Node<T> actualNode) {
        while (actualNode.data != null && actualNode != this.head && this.comparator.compare(item, actualNode.data) <= 0){
            if(actualNode.data.equals(item))
                return true;

            actualNode = actualNode.levels[0].back;
        }

        return false;
    }

    private boolean iterateRightUntilFound(T item, Node<T> actualNode) {
        while (actualNode.data != null && actualNode != this.tail && this.comparator.compare(item, actualNode.data) >= 0){
            if(actualNode.data.equals(item))
                return true;

            actualNode = actualNode.levels[0].next;
        }

        return false;
    }

    @Override
    public void clear() {
        for(int i = 0; i < MAX_LEVEL; i++){
            this.head.levels[i].next = null;
            this.head.levels[i].back = null;

            this.tail.levels[i].next = null;
            this.tail.levels[i].back = null;
        }

        this.size = 0;
    }

    @Override
    public int size() {
        return this.size;
    }

    @Override
    public boolean isEmpty() {
        return this.size == 0;
    }

    private boolean randomTest50() {
        return Math.random() >= 0.5;
    }

    @AllArgsConstructor
    private static class Node<T> {
        @Getter private final T data;
        @Getter private final Level<T>[] levels;

        public static <T> Node<T> of(T data) {
            return new Node<>(data, initializeLevelArray());
        }

        private static <T> Level<T>[] initializeLevelArray() {
            Level<T>[] levelArray = new Level[MAX_LEVEL + 1];

            for(int i = 0; i <= MAX_LEVEL; i++){
                levelArray[i] = new Level<>();
                levelArray[i] = new Level<>();
            }

            return levelArray;
        }
    }

    private static class Level<T> {
        @Getter private Node<T> next;
        @Getter private Node<T> back;
    }
}
