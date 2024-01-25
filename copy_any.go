package typact

import "reflect"

func (o Option[T]) CopyAny() Option[T] {
	if o.IsNone() {
		return None[T]()
	}

	anyVal := any(o.val)
	switch anyVal.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64,
		string, bool, complex64, complex128, uintptr:
		return Some(o.val)
	default:
	}

	vv, ok := anyVal.(copy[T])
	if ok {
		return Some(vv.Copy())
	}

	return None[T]()
}

func isScalar(vv any) bool {
	switch vv.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64,
		string, bool, complex64, complex128, uintptr:
		return true
	default:
		return false
	}
}

func (o Option[T]) NewCopy() Option[T] {
	if o.IsNone() {
		return None[T]()
	}

	refVal := reflect.ValueOf(o.val)
	if IsScalar(refVal.Kind()) {
		return Some(o.val)
	}

	// go currently does not support type assertion on generic types.
	// We try to reduce the overhead – which is especially true
	// for scalar types like string, int, ... – by moving only converting
	// it if there is no other way.
	// For scalar types a conversion to any causes to escape to the heap AND
	// create a new allocation.
	anyVal := any(o.val)

	vv, ok := anyVal.(copy[T])
	if ok {
		return Some(vv.Copy())
	}

	panic("unknown")
}
