package code

type Iterator[T comparable] interface {
	Next() T
	HastNext() bool
}

type Iterable[T comparable] interface {
	Iterate() Iterator[T]
}
