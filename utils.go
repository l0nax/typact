package typact

import (
	"reflect"

	"go.l0nax.org/typact/internal/features"
)

// zeroValue returns the zero value of T.
func zeroValue[T any]() (t T) {
	return t
}

// isScalarCopyable returns whether the given Kind is a scalar type
// which can be copied by value.
func isScalarCopyable(kind reflect.Kind) bool {
	switch kind {
	case reflect.Bool,
		reflect.Uintptr,
		reflect.Func,
		reflect.UnsafePointer,
		reflect.Invalid:

		return true
	case reflect.String:
		// WARN: With arenas, a string cannot be simply copied.
		return !features.GoArenaAvail
	}

	return false
}
