package code

import "errors"

type Streamable[T comparable] interface {
	Stream() * Stream[T]
}

type Stream[T comparable] struct {
	iterator Iterator[T]
}

func (stream * Stream[I]) Foreach(consumer func(it I)) * Stream[I] {
	for stream.iterator.HastNext() {
		consumer(stream.iterator.Next())
	}

	return stream
}

func (stream * Stream[I]) Filter(predicate func(it I) bool) * Stream[I] {
	list := &LinkedList[I]{}

	for stream.iterator.HastNext() {
		it := stream.iterator.Next()
		if predicate(it) {
			list.Add(it)
		}
	}

	return list.Stream()
}

func (stream * Stream[I]) ToList() * LinkedList[I] {
	list := &LinkedList[I]{}

	for stream.iterator.HastNext() {
		list.Add(stream.iterator.Next())
	}

	return list
}

func (stream * Stream[I]) FindAny(predicate func(it I) bool) (I, error) {
	var defaultIfNotFound I

	for stream.iterator.HastNext() {
		it := stream.iterator.Next()

		if predicate(it) {
			return it, nil
		}
	}

	return defaultIfNotFound, errors.New("not found")
}