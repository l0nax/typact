package typact

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"unsafe"
)

// Option represents an optional value.
// Every [Option] is either [Some] and contains a value, or [None],
// in which case it does not contain a value.
//
// It is based on the std::option::Option type from
// Rust (https://doc.rust-lang.org/std/option/enum.Option.html).
type Option[T any] struct {
	val  T
	some bool
}

// None returns [Option] with no value.
func None[T any]() Option[T] {
	return Option[T]{
		some: false,
	}
}

// Some returns [Option] with val.
//
//gcassert:inline
func Some[T any](val T) Option[T] {
	return Option[T]{
		some: true,
		val:  val,
	}
}

// TryWrap executes fn and wraps the value in an [Option].
func TryWrap[T any](fn func() (T, bool)) Option[T] {
	return Wrap(fn())
}

// Wrap wraps val and some into an [Option].
//
// To eagerly evaluate and wrap a value, use [TryWrap].
func Wrap[T any](val T, some bool) Option[T] {
	return Option[T]{
		some: some,
		val:  val,
	}
}

// IsSome returns true if o contains a value.
//
//gcassert:inline
func (o Option[T]) IsSome() bool {
	return o.some
}

// IsNone returns true if o contains no value.
//
//gcassert:inline
func (o Option[T]) IsNone() bool {
	return !o.IsSome()
}

// Deconstruct returns the value and whether it is present.
//
//gcassert:inline
func (o Option[T]) Deconstruct() (T, bool) {
	return o.val, o.some
}

// UnsafeUnwrap returns the value without checking whether
// the value is present.
//
// WARN: Only use this method as a last resort!
func (o Option[T]) UnsafeUnwrap() T {
	return o.val
}

// Unwrap returns the value or panics if it is
// not present.
//
//gcassert:inline
func (o Option[T]) Unwrap() T {
	if o.some {
		return o.val
	}

	panic("option does not contain a value")
}

// UnwrapOr returns the value of o, if present.
// Otherwise the provided value is returned
//
//gcassert:inline
func (o Option[T]) UnwrapOr(value T) T {
	if o.IsSome() {
		return o.val
	}

	return value
}

// UnwrapOrZero returns the value of o, if present
// Otherwise the zero value of T is returned
func (o Option[T]) UnwrapOrZero() T {
	if o.IsSome() {
		return o.val
	}

	return zeroValue[T]()
}

// UnwrapAsRef unwraps o and returns the reference to the value.
// If the value is not present, this method will panic.
func (o *Option[T]) UnwrapAsRef() *T {
	if o.IsSome() {
		return &o.val
	}

	panic("option does not contain a value")
}

// UnwrapOrElse returns the value of o, if present.
// Otherwise it executes fn and returns the value.
func (o Option[T]) UnwrapOrElse(fn func() T) T {
	if o.IsSome() {
		return o.val
	}

	return fn()
}

// Inspect executes fn if o contains a value.
// It returns o.
func (o Option[T]) Inspect(fn func(T)) Option[T] {
	if o.IsSome() {
		fn(o.UnsafeUnwrap())
	}

	return o
}

// Inserts val into o and returns the reference to the inserted value.
// If the option already contains a value, the old value is dropped.
//
// See [GetOrInsert] which doesn't update the value if the option already contains [Some].
func (o *Option[T]) Insert(val T) *T {
	*o = Some(val)

	return o.UnwrapAsRef()
}

// GetOrInsert inserts val in o if o is [None].
// It subsequently returns a pointer to the value.
//
// See [Insert], which updates the value even if the
// option already contains [Some].
func (o *Option[T]) GetOrInsert(val T) *T {
	if o.IsNone() {
		*o = Some(val)
	}

	return o.UnwrapAsRef()
}

// GetOrInsertWith inserts a value returned by fn into o, if o is [None].
// It subsequently returns a pointer to the value.
func (o *Option[T]) GetOrInsertWith(fn func() T) *T {
	if o.IsNone() {
		*o = Some(fn())
	}

	return o.UnwrapAsRef()
}

// Map maps Option[T] to Option[T] by calling fn on the value held by o, if [Some].
// It returns [Some] with the new value returned by fn.
// Otherwise [None] will be returned.
func (o Option[T]) Map(fn func(T) T) Option[T] {
	if o.IsSome() {
		return Some(fn(o.UnsafeUnwrap()))
	}

	return o
}

// MapOr returns the provided default result (if [None]),
// or applies fn to the contained value (if [Some]).
// Otherwise the provided (fallback) value is returned.
func (o Option[T]) MapOr(fn func(T) T, value T) T {
	if o.IsSome() {
		return fn(o.UnsafeUnwrap())
	}

	return value
}

// MapOrElse applies the function mapFn to the value held by o if it exists,
// and returns the result. If o does not hold a value, it applies valueFn and returns its result.
//
// This allows conditional transformation of the Option's value or generation of a default value when none is present.
func (o Option[T]) MapOrElse(mapFn func(T) T, valueFn func() T) T {
	if o.IsSome() {
		return mapFn(o.UnsafeUnwrap())
	}

	return valueFn()
}

// Filter returns [None] if o is [None], otherwise calls fn with the value of o and returns:
//   - [None] if fn returns false
//   - [Some] if fn returns true
func (o Option[T]) Filter(fn func(T) bool) Option[T] {
	if o.IsSome() && fn(o.UnsafeUnwrap()) {
		return o
	}

	return None[T]()
}

// Replace replaces o with [Some] of val and returns the old
// value of o.
func (o *Option[T]) Replace(val T) Option[T] {
	old := *o
	*o = Some(val)

	return old
}

// And returns [None] if o is [None], otherwise opt is returned.
func (o Option[T]) And(opt Option[T]) Option[T] {
	if o.IsSome() {
		return opt
	}

	return None[T]()
}

// AndThen returns None if o is none, otherwise fn is called and the return
// value is wrapped and returned.
func (o Option[T]) AndThen(fn func(T) Option[T]) Option[T] {
	if o.IsSome() {
		return fn(o.val)
	}

	return None[T]()
}

// Or returns the option if it contains a value, otherwise returns value.
func (o Option[T]) Or(value Option[T]) Option[T] {
	if o.IsSome() {
		return o
	}

	return value
}

// OrElse returns o, if [Some].
// Otherwise the return value of valueFn is returned.
func (o Option[T]) OrElse(valueFn func() Option[T]) Option[T] {
	if o.IsSome() {
		return o
	}

	return valueFn()
}

// Scan implements the [sql.Scanner] interface.
func (o *Option[T]) Scan(src any) error {
	// reset first
	o.some = false

	if src == nil {
		// only allocate in slow path.
		// this overrides any previously defined value in the field.
		o.val = zeroValue[T]()

		return nil
	}

	av, err := driver.DefaultParameterConverter.ConvertValue(src)
	if err != nil {
		// only allocate in slow path
		// this overrides any previously defined value in the field.
		o.val = zeroValue[T]()

		return errors.New("unable to scan Option[T]")
	}

	if v, ok := av.(T); ok {
		o.some = true
		o.val = v
	} else {
		// explicitly copy src to prevent heap escape.
		tmp := src
		return fmt.Errorf("got unexpected type %T", tmp)
	}

	return nil
}

// UnmarshalJSON implements the [json.Marshaler] interface.
// If value is not present, 'null' be encoded.
func (o Option[T]) MarshalJSON() ([]byte, error) {
	if o.some {
		return json.Marshal(o.val)
	}

	return json.Marshal(nil)
}

// UnmarshalJSON implements the [json.Unmarshaler] interface.
func (o *Option[T]) UnmarshalJSON(data []byte) error {
	// reset first
	o.some = false

	if bytes.Equal(data, []byte("null")) {
		// only allocate in slow path
		o.val = zeroValue[T]()

		return nil
	}

	err := json.Unmarshal(data, &o.val)
	if err != nil {
		// only allocate in slow path
		o.val = zeroValue[T]()

		return err
	}

	o.some = true

	return nil
}

// MarshalText implements the [encoding.TextMarshaler] interface.
// It returns the JSON representation.
func (o Option[T]) MarshalText() ([]byte, error) {
	return json.Marshal(o)
}

// UnmarshalText implements the [encoding.TextUnmarshaler] interface.
// It expects JSON as input.
func (o *Option[T]) UnmarshalText(data []byte) error {
	return json.Unmarshal(data, o)
}

type copy[T any] interface {
	Copy() T
}

func (o Option[T]) Copy() Option[T] {
	if o.IsNone() {
		return None[T]()
	}

	refVal := reflect.ValueOf(o.val)

	copyType := reflect.TypeFor[copy[T]]()
	if refVal.Type().Implements(copyType) {
		vv := refVal.Interface().(copy[T]).Copy()
		return Some(vv)
	}

	if IsScalar(refVal.Kind()) {
		vv := copyScalarValue(refVal).Interface().(T)
		return Some(vv)
	}

	panic("X")

	return None[T]()
}

func (o Option[T]) IS_ValueOf() bool {
	if o.IsNone() {
		return false
	}

	val := reflect.ValueOf(o.val)
	return IsScalar(val.Kind())
}

func (o Option[T]) IS_TypeFrom() bool {
	if o.IsNone() {
		return false
	}

	tt := reflect.TypeFor[T]()
	return IsScalar(tt.Kind())
}

func IsScalar(k reflect.Kind) bool {
	switch k {
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
		// If arena is not enabled, string can be copied as scalar safely
		// as it's immutable by design.
		return true
		//return !arenaIsEnabled
	}

	return false
}

func copyScalarValue(src reflect.Value) reflect.Value {
	if src.CanInterface() {
		return src
	}

	dst := newScalarValue(src)
	return dst.Convert(src.Type())
}

func newScalarValue(src reflect.Value) reflect.Value {
	// src is an unexported field value. Copy its value.
	switch src.Kind() {
	case reflect.Bool:
		return reflect.ValueOf(src.Bool())

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

	case reflect.Float32:
		return reflect.ValueOf(float32(src.Float()))
	case reflect.Float64:
		return reflect.ValueOf(src.Float())

	case reflect.Complex64:
		return reflect.ValueOf(complex64(src.Complex()))
	case reflect.Complex128:
		return reflect.ValueOf(src.Complex())

	case reflect.String:
		return reflect.ValueOf(src.String())
	case reflect.Func:
		panic("func")
	case reflect.UnsafePointer:
		return reflect.ValueOf(unsafe.Pointer(src.Pointer()))
	}

	panic(fmt.Errorf("go-clone: <bug> impossible type `%v` when cloning private field", src.Type()))
}
