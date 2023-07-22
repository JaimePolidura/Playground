package es.jaime.queues;

import lombok.AllArgsConstructor;
import lombok.Getter;
import lombok.SneakyThrows;

import java.util.concurrent.locks.Condition;
import java.util.concurrent.locks.Lock;
import java.util.concurrent.locks.ReentrantLock;

public final class BlockingQueue<T> implements Queue<T>{
    private final Lock lock;
    private final Condition notEmptySignal;
    private Node<T> head; //Pointing to element to dequeue
    private Node<T> tail; //Pointing to last element added
    private int size;

    public BlockingQueue() {
        this.lock = new ReentrantLock(true);
        this.notEmptySignal = this.lock.newCondition();
    }

    @Override
    public void enqueue(T element) {
        this.lock.lock();

        Node<T> toEnqueue = new Node<T>(element, null, null);

        if(this.isEmpty()){
            this.head = this.tail = toEnqueue;
        }else if(this.size == 1){
            this.head.prev = toEnqueue;
            toEnqueue.next = this.head;
            this.tail = toEnqueue;
        }else{
            this.tail.prev = toEnqueue;
            toEnqueue.next = this.tail;
            this.tail = toEnqueue;
        }

        this.size = size + 1;

        this.notEmptySignal.signalAll();

        this.lock.unlock();
    }

    @Override
    @SneakyThrows
    public T dequeue() {
        this.lock.lock();

        while (this.size == 0)
            this.notEmptySignal.await();

        Node<T> node = this.head;
        this.head = this.head.prev;
        this.size = size - 1;

        this.lock.unlock();

        return node.data;
    }

    @Override
    public int size() {
        return this.size;
    }

    @AllArgsConstructor
    private static class Node<T> {
        @Getter private T data;
        @Getter private Node<T> next;
        @Getter private Node<T> prev;
    }
}
