package immutable

import "go.l0nax.org/typact"

// FromSlice returns a new [List] from s.
// After this function has been called, s should not be used anymore.
func FromSlice[T any](s []T) List[T] {
	return List[T]{
		value: s,
	}
}

// List is an immutable list of elements.
type List[T any] struct {
	value []T
}

// Get returns the value at the given index.
// If the index is below zero or greater than len(list) - 1, [typact.None] will be returned.
func (l List[T]) Get(index int) typact.Option[T] {
	if index < 0 || index >= len(l.value) {
		return typact.None[T]()
	}

	return typact.Some(l.value[index])
}

// Len returns the number of elements in the list.
func (l List[T]) Len() int {
	return len(l.value)
}

// Iter is an iterator over the elements in the list.
func (l List[T]) Iter(yield func(index int, value T) bool) {
	for i, v := range l.value {
		if !yield(i, v) {
			break
		}
	}
}
