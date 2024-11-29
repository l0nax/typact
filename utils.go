package typact

import (
	"reflect"
	"unsafe"

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

// string2Bytes converts the given string to a byte slice without memory allocation.
//
// Note it may break if string and/or slice header will change in future go versions.
func string2Bytes(s string) (b []byte) {
	return unsafe.Slice(unsafe.StringData(s), len(s))
}

// bytes2String converts a byte slice to a string in a performant way.
func bytes2String(bs []byte) string {
	return unsafe.String(unsafe.SliceData(bs), len(bs))
}
