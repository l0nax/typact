package typact

import "reflect"

func reflectIsScalar(kind reflect.Kind) bool {
	switch kind {
	case reflect.Bool,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr,
		reflect.Float32, reflect.Float64,
		reflect.Complex64, reflect.Complex128,
		reflect.Func,
		reflect.UnsafePointer,
		reflect.Invalid:
		return true

	case reflect.String:
		// XXX: If arena is not enabled, string can be copied safely
		// as it's immutable by design.
		return true
	}

	return false
}

func (o Option[T]) Copy() Option[T] {
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
