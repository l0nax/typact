package typact

import (
	"fmt"
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

	// this ONLY includes direct scalars, i.e. type aliases are included
	// A custom type of, e.g., string will return false here.
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
	refKind := refVal.Kind()

	if isScalarCopyable(refKind) {
		return Some(o.val)
	}

	return slowClone[T](o.val, refVal, refKind)
}

// slowClone will clone the underlying value of raw and return it as [Option[T]].
func slowClone[T any](raw T, val reflect.Value, kind reflect.Kind) Option[T] {
	switch {
	case isScalarCopyable(kind):
		return Some(raw)
	case kind == reflect.String:
		// XXX: Special case: if [isScalarCopyable] returns false
		// for a string, it means that we need to explicitly create a copy.
		srcStr := any(raw).(string)
		cpy := strings.Clone(srcStr)

		return Some(any(cpy).(T))
	case kind == reflect.Slice:
		return Some(cloneSlice[T](raw, val))
	case kind == reflect.Ptr:
		return Some(clonePtr[T](val))
	}

	panic("unable to clone value: type does not implement std.Cloner interface")
}

// slowReflectClone works like [slowClone] but it does not operate on the raw value.
// But instead uses [reflect] only.
func slowReflectClone(value reflect.Value, kind reflect.Kind) reflect.Value {
	switch {
	case isScalarCopyable(kind):
		return cloneScalarValue(value)
	}

	panic(
		fmt.Errorf(
			"unable to clone value: type <%v> not supported to clone with reflect",
			kind,
		),
	)
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

// clonePtr returns a deep copy of the val pointer.
func clonePtr[T any](val reflect.Value) T {
	if val.IsNil() {
		return zeroValue[T]()
	}

	// TODO: Handle opaque pointers like
	// `elliptic.Curve`, which is `*elliptic.CurveParam` or `elliptic.p256Curve`;

	elem := val.Elem()
	elemType := elem.Type()
	elemKind := elem.Kind()

	dst := reflect.New(elemType)

	switch elemKind {
	case reflect.Struct:
		panic("pointer to struct MUST implement std.Cloner[T] interface!")
	case reflect.Array:
		panic("TODO: arrays currently not supported")
	}

	cloned := slowReflectClone(elem, elemKind)
	dst.Elem().Set(cloned)

	return dst.Interface().(T)
}

// isScalar reports whether the type T is a scalar type.
// The check is done without allocating on the heap.
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

// cloneScalarValue returns a clone of a scalar value val.
func cloneScalarValue(val reflect.Value) reflect.Value {
	if val.CanInterface() {
		return val
	}

	// TODO: Benchmark
	// dst := newScalarValue(src)

	dst := reflect.New(val.Type())
	return dst.Convert(val.Type())
}

func newScalarValue(src reflect.Value) reflect.Value {
	switch src.Kind() {
	case reflect.Bool:
		return reflect.ValueOf(src.Bool())

	// Numbers
	case reflect.Int:
		return reflect.ValueOf(int(src.Int()))
	case reflect.Int8:
		return reflect.ValueOf(int8(src.Int()))
	case reflect.Int16:
		return reflect.ValueOf(int16(src.Int()))
	case reflect.Int32:
		return reflect.ValueOf(int32(src.Int()))
	case reflect.Int64:
		return reflect.ValueOf(src.Int())

	// Positive Numbers Only
	case reflect.Uint:
		return reflect.ValueOf(uint(src.Uint()))
	case reflect.Uint8:
		return reflect.ValueOf(uint8(src.Uint()))
	case reflect.Uint16:
		return reflect.ValueOf(uint16(src.Uint()))
	case reflect.Uint32:
		return reflect.ValueOf(uint32(src.Uint()))
	case reflect.Uint64:
		return reflect.ValueOf(src.Uint())
	case reflect.Uintptr:
		return reflect.ValueOf(uintptr(src.Uint()))

	// Non Integer Numbers
	case reflect.Float32:
		return reflect.ValueOf(float32(src.Float()))
	case reflect.Float64:
		return reflect.ValueOf(src.Float())

	// Z
	case reflect.Complex64:
		return reflect.ValueOf(complex64(src.Complex()))
	case reflect.Complex128:
		return reflect.ValueOf(src.Complex())

	case reflect.String:
		// TODO(l0nax): Validate in case of arenas
		return reflect.ValueOf(src.String())

	case reflect.Func:
		t := src.Type()

		if src.IsNil() {
			return reflect.Zero(t)
		}

		// XXX: Very special, very rare case: if the RO flag is set, we COULD
		// force unset it. But unless someone reports this SPECIFIC usecase
		// I will not implement it, because it may break with any Go update
		// and is very unsafe.
		panic(
			"unable to clone func(...): RO flag may be set: please report this to the repository!",
		)
	case reflect.UnsafePointer:
		return reflect.ValueOf(unsafe.Pointer(src.Pointer()))
	}

	// I have no idea how this can be triggered, but just in case
	panic(fmt.Errorf("BUG: unable to clone type '%v'", src.Type()))
}
