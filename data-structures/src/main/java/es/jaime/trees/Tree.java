package es.jaime.trees;

import es.jaime.DataStructure;

public interface Tree<T extends Comparable<T>> extends DataStructure {
    boolean insert(T data);

    boolean search(T data);

    boolean delete(T data);
}
