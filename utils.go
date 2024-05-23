package typact

import (
	"reflect"

	"go.l0nax.org/typact/internal/features"
)

// isScalarCopyable returns whether the given Kind is a scalar type
// which can be copied by value.
func isScalarCopyable(kind reflect.Kind) bool {
	switch kind {
	case reflect.Bool,
		reflect.Uintptr,
		reflect.Func,
		reflect.UnsafePointer,
		reflect.Invalid,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64,
		reflect.Complex64, reflect.Complex128:

		return true
	case reflect.String:
		// WARN: With arenas, a string cannot be simply copied.
		return !features.GoArenaAvail
	}

	return false
}
