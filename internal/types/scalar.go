package types

// IsScalar reports whether the type T is a scalar type.
// The check is done without allocating on the heap.
func IsScalar[T any]() bool {
	var v T

	switch any(v).(type) {
	case int, int8, int16, int32, int64,
		uint, uint8, uint16, uint32, uint64,
		float32, float64,
		complex64, complex128,
		bool,
		string:

		return true
	}

	return false
}

// ZeroValue returns the zero value of T.
func ZeroValue[T any]() (t T) {
	return t
}
