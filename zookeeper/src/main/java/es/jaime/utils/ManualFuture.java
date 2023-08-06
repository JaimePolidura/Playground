package es.jaime.utils;

import java.util.concurrent.*;

public final class ManualFuture<T> implements Future<T> {
    private final CountDownLatch latch;

    private State state;
    private T value;

    public ManualFuture(T value) {
        this.latch = new CountDownLatch(1);
        this.state = State.PENDING;
        this.value = value;
    }

    public void complete(T value) {
        this.value = value;
        this.state = State.DONE;

        this.latch.countDown();
    }

    public static <T> ManualFuture<T> ofDefault(T defaultValue) {
        return new ManualFuture<>(defaultValue);
    }

    @Override
    public boolean cancel(boolean mayInterruptIfRunning) {
        this.state = State.CANCELLED;

        return true;
    }

    @Override
    public boolean isCancelled() {
        return this.state == State.CANCELLED;
    }

    @Override
    public boolean isDone() {
        return this.state == State.DONE;
    }

    @Override
    public T get() throws InterruptedException {
        if (state == State.DONE || state == State.CANCELLED) {
            return value;
        }

        latch.await();

        state = State.DONE;

        return value;
    }

    @Override
    public T get(long timeout, TimeUnit unit) throws InterruptedException, TimeoutException {
        if (state == State.DONE || state == State.CANCELLED) {
            return value;
        }

        boolean timedOut = latch.await(timeout, unit);

        if (!timedOut) {
            state = State.DONE;
            return value;
        } else {
            throw new TimeoutException();
        }
    }

    private enum State {
        PENDING, CANCELLED, DONE
    }
}
