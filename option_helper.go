package typact

import (
	"reflect"
	"strings"
	"unsafe"

	"go.l0nax.org/typact/std"
)

// Clone returns a deep copy of T, if o contains a value.
// Otherwise [None] is returned.
//
// The below special types are handled by the method by simply copying
// the value:
//
//   - Scalar types: all number-like types are copied by value.
//   - string: Copied by value as string is immutable by design.
//   - func: Copied by value as func is an opaque pointer at runtime.
//   - unsafe.Pointer: Copied by value as we don't know what's in it.
//   - chan: A new empty chan is created as we cannot read data inside the old chan.
//
// WARN: If T is not part of the special types above AND not DOES NOT
// implement [std.Cloner], this method will panic!
//
// BUG: Currently unsupported are arrays of any type.
//
// Unstable: This method is unstable and not guarded by the SemVer promise!
func (o Option[T]) Clone() Option[T] {
	if o.IsNone() {
		return None[T]()
	}

	if isScalar[T]() {
		return Some(o.val)
	}

	// NOTE: Converting to any should be last restort because if we use it
	// with a scalar type, it will create a new allocation just for
	// the conversion.
	//
	// With [implementsCloner] we can check if T implements [std.Cloner]
	// without allocating.
	// Additionally, we skip all the reflect checks, which is nice too.
	if implementsCloner[T]() {
		cloner := any(o.val).(std.Cloner[T])
		return Some(cloner.Clone())
	}

	// NOTE: After extensive benchmarking, I found out that
	// ValueOf is the fastest way to accessing the information.
	//
	// Sadly, Go currently does not support type assertion on
	// generic types.
	// We might switch to type assertion, once available.
	refVal := reflect.ValueOf(o.val)
	kind := refVal.Kind()

	switch {
	case isScalarCopyable(kind):
		return Some(o.val)
	case kind == reflect.String:
		// XXX: Special case: if [isScalarCopyable] returns false
		// for a string, it means that we need to explicitly create a copy.
		srcStr := any(o.val).(string)
		cpy := strings.Clone(srcStr)

		return Some(any(cpy).(T))
	case kind == reflect.Slice:
		return Some(cloneSlice[T](o.UnsafeUnwrap(), refVal))
	}

	panic("unable to clone value: type does not implement std.Cloner interface")
}

// maxByteSize is a large enough value to cheat Go compiler
// when converting unsafe address to []byte.
// It's not actually used in runtime.
//
// The value 2^30 is the max value AFAIK to make Go compiler happy on all archs.
const maxByteSize = 1 << 30

// cloneSlice returns a deep copy of the val slice.
//
// NOTE: T represents the slice, not the type of a slice element!
func cloneSlice[T any](raw T, val reflect.Value) T {
	if val.IsNil() {
		return zeroValue[T]()
	}

	valType := val.Type()
	elems := val.Len()
	vCap := val.Cap()

	// for scalar slices, we can copy the underlying values directly
	// => fast path.
	if isScalarCopyable(valType.Elem().Kind()) {
		ret := reflect.MakeSlice(valType, elems, vCap)
		src := unsafe.Pointer(val.Pointer())
		dst := unsafe.Pointer(ret.Pointer())
		sz := int(valType.Elem().Size())

		l := elems * sz
		cc := vCap * sz

		copy((*[maxByteSize]byte)(dst)[:l:cc], (*[maxByteSize]byte)(src)[:l:cc])

		return ret.Interface().(T)
	}

	// fast path check if T implements std.Cloner.
	// We can convert to any here because we know that the element of T
	// is not a scalar thus we are not creating an unnecessary allocation.
	vv, ok := any(raw).(std.Cloner[T])
	if ok {
		return vv.Clone()
	}

	ret := reflect.MakeSlice(valType, elems, vCap)

	// The caller did not implement a helper type, thus
	// we have to manually clone the elements one by one.
	for i := 0; i < elems; i++ {
		elem := val.Index(i)
		vv := elem.Interface().(std.Cloner[T])
		ret.Index(i).Set(reflect.ValueOf(vv.Clone()))
	}

	return ret.Interface().(T)
}

func isScalar[T any]() bool {
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

// implementsCloner returns true if T implements the [std.Cloner] interface.
func implementsCloner[T any]() bool {
	var v T

	_, ok := any(v).(std.Cloner[T])
	return ok
}
