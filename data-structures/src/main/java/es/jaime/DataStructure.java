package es.jaime;

public interface DataStructure {
    void clear();

    int size();

    default boolean isEmpty() {
        return this.size() == 0;
    }
}
