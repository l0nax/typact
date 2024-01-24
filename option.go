package typact

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
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
