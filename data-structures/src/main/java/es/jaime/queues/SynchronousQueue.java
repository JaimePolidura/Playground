package es.jaime.queues;

import lombok.SneakyThrows;

import java.util.concurrent.locks.Condition;
import java.util.concurrent.locks.Lock;
import java.util.concurrent.locks.ReentrantLock;

public final class SynchronousQueue<T> implements Queue<T>{
    private final Lock lock;
    private final Condition itemDequeued;
    private final Condition itemQueued;
    private T item;

    public SynchronousQueue() {
        this.lock = new ReentrantLock();
        this.itemDequeued = this.lock.newCondition();
        this.itemQueued = this.lock.newCondition();
    }

    @Override
    @SneakyThrows
    public void enqueue(T itemToEnqueue) {
        try {
            this.lock.lock();

            this.item = itemToEnqueue;

            this.itemQueued.signalAll();

            while (item != null)
                this.itemDequeued.await();
        }finally {
            this.lock.unlock();
        }
    }

    @Override
    @SneakyThrows
    public T dequeue() {
        try{
            this.lock.lock();

            while (item == null)
                this.itemQueued.await();

            T item = this.item;

            this.item = null;

            this.itemDequeued.signalAll();

            return item;
        }finally {
            this.lock.unlock();
        }
    }

    @Override
    public int size() {
        return this.item != null ? 1 : 0;
    }
}
