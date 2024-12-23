package pred

// IsZero returns whether val holds the zero value.
func IsZero[T comparable](val T) bool {
	var zero T
	return val == zero
}

// IsNotZero returns whether val does not hold the zero value.
// It's the opposite of [IsZero].
func IsNotZero[T comparable](val T) bool {
	var zero T
	return val != zero
}
