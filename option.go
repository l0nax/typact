package typact

import (
	"bytes"
	"database/sql"
	"database/sql/driver"
	"encoding"
	"encoding/json"
	"fmt"
	"math"
	"strconv"

	"go.l0nax.org/typact/internal/types"
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

// IsZero returns whether o is [None].
//
// NOTE: This method has only been added for compatibility with
// unmarshaling/ marshaling libraries, e.g. YAML, TOML, JSON
// in conjunction with the "omitzero" tag.
//
// Deprecated: Use [Option.IsNone] instead!
func (o Option[T]) IsZero() bool {
	return !o.some
}

// IsSome returns true if o contains a value.
//
//gcassert:inline
func (o Option[T]) IsSome() bool {
	return o.some
}

// IsSomeAnd returns true if o contains a value and fn(o) returns true.
//
//gcassert:inline
func (o Option[T]) IsSomeAnd(fn func(T) bool) bool {
	return o.some && fn(o.UnsafeUnwrap())
}

// IsNone returns true if o contains no value.
//
//gcassert:inline
func (o Option[T]) IsNone() bool {
	return !o.some
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

// Expect returns the contained value of o, if it is present.
// Otherwise it panics with msg.
//
//gcassert:inline
func (o Option[T]) Expect(msg string) T {
	if o.some {
		return o.UnsafeUnwrap()
	}

	panic(msg)
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
	if o.some {
		return o.val
	}

	return value
}

// UnwrapOrZero returns the value of o, if present
// Otherwise the zero value of T is returned
func (o Option[T]) UnwrapOrZero() T {
	if o.some {
		return o.val
	}

	return types.ZeroValue[T]()
}

// UnwrapAsRef unwraps o and returns the reference to the value.
// If the value is not present, this method will panic.
func (o *Option[T]) UnwrapAsRef() *T {
	if o.some {
		return &o.val
	}

	panic("option does not contain a value")
}

// UnwrapOrElse returns the value of o, if present.
// Otherwise it executes fn and returns the value.
func (o Option[T]) UnwrapOrElse(fn func() T) T {
	if o.some {
		return o.val
	}

	return fn()
}

// Inspect executes fn if o contains a value.
// It returns o.
func (o Option[T]) Inspect(fn func(T)) Option[T] {
	if o.some {
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
	if !o.some {
		*o = Some(val)
	}

	return o.UnwrapAsRef()
}

// GetOrInsertWith inserts a value returned by fn into o, if o is [None].
// It subsequently returns a pointer to the value.
func (o *Option[T]) GetOrInsertWith(fn func() T) *T {
	if !o.some {
		*o = Some(fn())
	}

	return o.UnwrapAsRef()
}

// Map maps Option[T] to Option[T] by calling fn on the value held by o, if [Some].
// It returns [Some] with the new value returned by fn.
// Otherwise [None] will be returned.
func (o Option[T]) Map(fn func(T) T) Option[T] {
	if o.some {
		return Some(fn(o.UnsafeUnwrap()))
	}

	return o
}

// MapOr returns the provided default result (if [None]),
// or applies fn to the contained value (if [Some]).
// Otherwise the provided (fallback) value is returned.
func (o Option[T]) MapOr(fn func(T) T, value T) T {
	if o.some {
		return fn(o.UnsafeUnwrap())
	}

	return value
}

// MapOrElse applies the function mapFn to the value held by o if it exists,
// and returns the result. If o does not hold a value, it applies valueFn and returns its result.
//
// This allows conditional transformation of the Option's value or generation of a default value when none is present.
func (o Option[T]) MapOrElse(mapFn func(T) T, valueFn func() T) T {
	if o.some {
		return mapFn(o.UnsafeUnwrap())
	}

	return valueFn()
}

// Filter returns [None] if o is [None], otherwise calls fn with the value of o and returns:
//   - [None] if fn returns false
//   - [Some] if fn returns true
func (o Option[T]) Filter(fn func(T) bool) Option[T] {
	if o.some && fn(o.UnsafeUnwrap()) {
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
	if o.some {
		return opt
	}

	return None[T]()
}

// AndThen returns None if o is none, otherwise fn is called and the return
// value is wrapped and returned.
func (o Option[T]) AndThen(fn func(T) Option[T]) Option[T] {
	if o.some {
		return fn(o.val)
	}

	return None[T]()
}

// Or returns the option if it contains a value, otherwise returns value.
func (o Option[T]) Or(value Option[T]) Option[T] {
	if o.some {
		return o
	}

	return value
}

// OrElse returns o, if [Some].
// Otherwise the return value of valueFn is returned.
func (o Option[T]) OrElse(valueFn func() Option[T]) Option[T] {
	if o.some {
		return o
	}

	return valueFn()
}

// Take takes the value of o and returns it, leaving [None] in its place.
//
// Experimental: This method is considered experimental and may change or be removed in the future.
func (o *Option[T]) Take() Option[T] {
	if !o.some {
		return None[T]()
	}

	vv := Some(o.UnsafeUnwrap())
	*o = None[T]()

	return vv
}

// Value implements the [driver.Valuer] interface.
// It returns NULL if o is [None], otherwise it
// returns the value of o.
//
// If T implements the [driver.Valuer] interface, the method
// will be called instead.
func (o Option[T]) Value() (driver.Value, error) {
	if !o.some {
		return nil, nil
	}

	if implementsDriverValuer[T]() {
		return any(o.val).(driver.Valuer).Value()
	}

	return driver.DefaultParameterConverter.ConvertValue(o.val)
}

// Scan implements the [sql.Scanner] interface.
//
// If *T implements [sql.Scanner], the custom method will be called.
func (o *Option[T]) Scan(src any) error {
	// reset first
	o.some = false

	if src == nil {
		// only allocate in slow path.
		// this overrides any previously defined value in the field.
		o.val = types.ZeroValue[T]()

		return nil
	}

	if implementsSQLScanner[T]() {
		// TODO(l0nax): Add tests to check if override works!
		// we first ensure to set o.val to the zero value, just in case
		o.val = types.ZeroValue[T]()

		scanner := any(&o.val).(sql.Scanner)

		err := scanner.Scan(src)
		o.some = err == nil

		return err
	}

	av, err := driver.DefaultParameterConverter.ConvertValue(src)
	if err != nil {
		// only allocate in slow path
		// this overrides any previously defined value in the field.
		o.val = types.ZeroValue[T]()

		// TODO(l0nax): Wrap the returned error and return it!
		return err
	}

	v, ok := av.(T)
	if !ok {
		// explicitly copy src to prevent heap escape.
		tmp := src
		return fmt.Errorf("got unexpected type %T", tmp)
	}

	o.some = true
	o.val = v

	return nil
}

// implementsSqlScanner returns true if the pointer of T implements [sql.Scanner].
func implementsSQLScanner[T any]() bool {
	var zero *T

	_, ok := any(zero).(sql.Scanner)
	return ok
}

// implementsDriverValuer returns true if T implements [driver.Valuer].
func implementsDriverValuer[T any]() bool {
	var zero T

	_, ok := any(&zero).(driver.Valuer)
	return ok
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
		o.val = types.ZeroValue[T]()

		return nil
	}

	err := json.Unmarshal(data, &o.val)
	if err != nil {
		// only allocate in slow path
		o.val = types.ZeroValue[T]()

		// TODO(l0nax): Wrap the returned error and return it!
		return err
	}

	o.some = true

	return nil
}

// MarshalText implements the [encoding.TextMarshaler] interface.
//
// Please not that for scalar types it is advised to define the "omitempty" tag!
func (o Option[T]) MarshalText() ([]byte, error) {
	if !o.some {
		return nil, nil
	}

	if !types.IsScalar[T]() {
		enc, ok := any(o.val).(encoding.TextMarshaler)
		if !ok {
			return nil, fmt.Errorf("type %T does not implement encoding.TextMarshaler", o.val)
		}

		return enc.MarshalText()
	}

	// it's a scalar type
	zz := any(o.val)

	switch val := zz.(type) {
	case string:
		return string2Bytes(val), nil

	case []byte:
		return val, nil

	case int:
		raw := strconv.FormatInt(int64(val), 10)
		return string2Bytes(raw), nil

	case uint:
		raw := strconv.FormatUint(uint64(val), 10)
		return string2Bytes(raw), nil

	case int8:
		raw := strconv.FormatInt(int64(val), 10)
		return string2Bytes(raw), nil

	case uint8:
		raw := strconv.FormatUint(uint64(val), 10)
		return string2Bytes(raw), nil

	case int16:
		raw := strconv.FormatInt(int64(val), 10)
		return string2Bytes(raw), nil

	case uint16:
		raw := strconv.FormatUint(uint64(val), 10)
		return string2Bytes(raw), nil

	case int32:
		raw := strconv.FormatInt(int64(val), 10)
		return string2Bytes(raw), nil

	case uint32:
		raw := strconv.FormatUint(uint64(val), 10)
		return string2Bytes(raw), nil

	case int64:
		raw := strconv.FormatInt(int64(val), 10)
		return string2Bytes(raw), nil

	case uint64:
		raw := strconv.FormatUint(uint64(val), 10)
		return string2Bytes(raw), nil

	case float32:
		raw := strconv.FormatFloat(float64(val), 'f', -1, 32)
		return string2Bytes(raw), nil

	case float64:
		raw := strconv.FormatFloat(float64(val), 'f', -1, 64)
		return string2Bytes(raw), nil

	case bool:
		if val {
			return []byte("true"), nil
		}

		return []byte("false"), nil
	}

	return nil, fmt.Errorf("type %T does not implement encoding.TextMarshaler", o.val)
}

// UnmarshalText implements the [encoding.TextUnmarshaler] interface.
//
// NOTE: when using scalar types, it is advised use the "omitzero" tag!
func (o *Option[T]) UnmarshalText(data []byte) error {
	o.some = false

	if !types.IsScalar[T]() {
		enc, ok := any(o.val).(encoding.TextUnmarshaler)
		if !ok {
			// only allocate in slow path.
			// this overrides any previously defined value in the field.
			o.val = types.ZeroValue[T]()

			return fmt.Errorf("type %T does not implement encoding.TextUnmarshaler", o.val)
		}

		if err := enc.UnmarshalText(data); err != nil {
			// only allocate in slow path.
			// this overrides any previously defined value in the field.
			o.val = types.ZeroValue[T]()

			return err
		}

		o.some = true

		return nil
	}

	err := unmarshalText(&o.val, data)
	if err != nil {
		// only allocate in slow path.
		// this overrides any previously defined value in the field.
		o.val = types.ZeroValue[T]()

		return fmt.Errorf("error unmarshaling data: %w", err)
	}

	o.some = true

	return nil
}

func unmarshalText(dest interface{}, data []byte) error {
	switch val := any(dest).(type) {
	case *string:
		*val = string(data)

	case *int:
		num, err := strconv.ParseInt(bytes2String(data), 10, 64)
		if err != nil {
			return err
		}

		if num < math.MinInt || num > math.MaxInt {
			panic("under/overflow")
		}

		*val = int(num)

	case *int8:
		num, err := strconv.ParseInt(bytes2String(data), 10, 8)
		if err != nil {
			return err
		}

		*val = int8(num)

	case *int16:
		num, err := strconv.ParseInt(bytes2String(data), 10, 16)
		if err != nil {
			return err
		}

		*val = int16(num)

	case *int32:
		num, err := strconv.ParseInt(bytes2String(data), 10, 32)
		if err != nil {
			return err
		}

		*val = int32(num)

	case *int64:
		num, err := strconv.ParseInt(bytes2String(data), 10, 64)
		if err != nil {
			return err
		}

		*val = int64(num)

	case *uint8:
		num, err := strconv.ParseUint(bytes2String(data), 10, 8)
		if err != nil {
			return err
		}

		*val = uint8(num)

	case *uint16:
		num, err := strconv.ParseUint(bytes2String(data), 10, 16)
		if err != nil {
			return err
		}

		*val = uint16(num)

	case *uint32:
		num, err := strconv.ParseUint(bytes2String(data), 10, 32)
		if err != nil {
			return err
		}

		*val = uint32(num)

	case *uint64:
		num, err := strconv.ParseUint(bytes2String(data), 10, 64)
		if err != nil {
			return err
		}

		*val = uint64(num)

	case *float32:
		num, err := strconv.ParseFloat(bytes2String(data), 32)
		if err != nil {
			return err
		}

		*val = float32(num)

	case *float64:
		num, err := strconv.ParseFloat(bytes2String(data), 64)
		if err != nil {
			return err
		}

		*val = float64(num)

	case *[]byte:
		*val = []byte(data)

	case *bool:
		switch bytes2String(data) {
		case "true":
			*val = true

		case "false":
			*val = false

		default:
			return fmt.Errorf("invalid boolean value: %q", data)
		}

	default:
		return fmt.Errorf("BUG: scalar type %T is not supported", dest)
	}

	return nil
}
