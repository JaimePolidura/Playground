package es.jaime.queues;

import lombok.AllArgsConstructor;
import lombok.Getter;
import lombok.SneakyThrows;

import java.util.EmptyStackException;
import java.util.concurrent.atomic.AtomicInteger;
import java.util.concurrent.atomic.AtomicReference;

/**
 * Dequeue -> Head  (-> next) Tail <- to enqueue
 *
 * See the art of multicore page 251 Unbounded lock-free queue
 */
public final class LockFreeQueue<T> implements Queue<T> {
    private final AtomicReference<Node<T>> head; //Pointing to element to dequeue
    private final AtomicReference<Node<T>> tail; //Pointing to last element added
    private final AtomicInteger size;

    public LockFreeQueue() {
        Node<T> sentinelHeadNode = new Node<>(null, null);
        sentinelHeadNode.next = new AtomicReference<>(null);

        this.head = new AtomicReference<>(sentinelHeadNode);
        this.tail = new AtomicReference<>(null);
        this.size = new AtomicInteger(0);
    }

    @Override
    public void enqueue(T element) {
        Node<T> toEnqueue = new Node<T>(element, this.tail);

        while (true) {
            Node<T> last = tail.get();
            Node<T> next = last.next.get();

            if(last == tail.get()) {
                if(next == null){
                    if(last.next.compareAndSet(null, toEnqueue)){
                        tail.compareAndSet(last, toEnqueue);
                        size.getAndIncrement();
                        return;
                    }
                }else{
                    tail.compareAndSet(last, next);
                }
            }
        }
    }

    @Override
    @SneakyThrows
    public T dequeue() {
        while (true) {
            Node<T> first = this.head.get();
            Node<T> last = this.tail.get();
            Node<T> next = first.next.get();

            if(next == null)
                throw new EmptyStackException();

            if(first == this.head.get()){
                if(first == this.tail.get()){
                    this.tail.compareAndSet(last, next);
                    continue;
                }

                if(this.head.compareAndSet(first, next)){
                    this.size.decrementAndGet();
                    return first.data;
                }
            }
        }
    }

    @Override
    public int size() {
        return this.size.get();
    }

    @AllArgsConstructor
    private static class Node<T> {
        @Getter private T data;
        @Getter private AtomicReference<Node<T>> next;
    }
}
