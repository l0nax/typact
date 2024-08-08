package immutable

import "go.l0nax.org/typact"

// FromSlice returns a new [List] from s.
// After this function has been called, s should not be used anymore.
func FromSlice[S ~[]E, E any](s S) List[S, E] {
	return List[S, E]{
		value: s,
	}
}

// List is an immutable list of elements.
type List[S ~[]E, E any] struct {
	value S
}

// Get returns the value at the given index.
// If the index is below zero or greater than len(list) - 1, [typact.None] will be returned.
func (l List[S, E]) Get(index int) typact.Option[E] {
	if index < 0 || index >= len(l.value) {
		return typact.None[E]()
	}

	return typact.Some(l.value[index])
}

// Len returns the number of elements in the list.
func (l List[S, E]) Len() int {
	return len(l.value)
}

// Iter is an iterator over the elements in the list.
func (l List[S, E]) Iter(yield func(index int, value E) bool) {
	for i, v := range l.value {
		if !yield(i, v) {
			break
		}
	}
}
