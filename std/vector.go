package std

// NOTE: This implementation was added due to performance considerations
// with the Clone method of Option for slices of custom types.
//
// This implementation is based on the [container/vector.Vector] implementation
// from the standard library.

// TODO: Once iterators are added to GoLang (https://github.com/golang/go/issues/61897)
// it would be nice to support them.

// Vector is a slice of elements of type T.
// It behaves like a slice and is similar to [container/vector.Vector].
type Vector[T any] []T

// VectorFromSlice returns a vector from the provided slice.
// The slice should not be used afterwards.
func VectorFromSlice[T any](s []T) Vector[T] {
	return s
}

// AppendVector appends other to v.
// After calling this method, other should not be used anymore.
func (v *Vector[T]) AppendVector(other Vector[T]) {
	*v = append(*v, other...)
}

// Copy returns a shallow clone of v.
func (v Vector[T]) Copy() Vector[T] {
	if v == nil {
		return nil
	}

	return append(Vector[T](nil), v...)
}

// Clone returns a deep clone of v.
func (v Vector[T]) Clone() Vector[T] {
	return append(Vector[T](nil), v...)
}

func _() {
	x := VectorFromSlice([]int{0, 1, 2, 3, 4, 5})
	_ = x[0]
}
