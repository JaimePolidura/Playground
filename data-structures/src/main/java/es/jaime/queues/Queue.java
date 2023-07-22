package es.jaime.queues;

public interface Queue<T> {
    void enqueue(T element);
    T dequeue();
    int size();

    default boolean isEmpty() {
        return this.size() == 0;
    }
}

