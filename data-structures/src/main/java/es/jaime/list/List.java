package es.jaime.list;

public interface List<T> {
    boolean add(T item);

    T get(int index);

    int size();

    void clear();

    boolean contains(T item);

    boolean isEmpty();
}
