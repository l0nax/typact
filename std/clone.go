package std

// Cloner implements the Clone method, which allows one to deeply clone
// T.
type Cloner[T any] interface {
	// Clone returns a deep copy of T.
	Clone() T
}
