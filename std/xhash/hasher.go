package xhash

import (
	"hash/maphash"
	"math"
	"reflect"
)

// NewHasher returns the default [Hasher] implementation with a random seed.
//
// NOTE: Because of [Hash flooding] the seed is generated with each
// call to this function.
//
// [Hash flooding]: https://en.wikipedia.org/wiki/Collision_attack#Hash_flooding
func NewHasher() Hasher {
	dh := &defaultHasher{}
	dh.hh.SetSeed(maphash.MakeSeed())

	return dh
}

// defaultHasher is the default [Hasher] implementation using [hash/maphash].
type defaultHasher struct {
	hh maphash.Hash
}

// WriteFloat64 implements Hasher.
func (d *defaultHasher) WriteFloat64(n float64) {
	if n == 0 {
		d.hh.WriteByte(0)
		return
	}

	d.WriteUint64(math.Float64bits(n))
}

// WriteInt implements Hasher.
func (d *defaultHasher) WriteInt(n int) {
	d.WriteUint64(uint64(n))
}

// WriteUint64 implements Hasher.
func (d *defaultHasher) WriteUint64(n uint64) {
	data := [8]byte{
		byte(n >> 56),
		byte(n >> 48),
		byte(n >> 40),
		byte(n >> 32),
		byte(n >> 24),
		byte(n >> 16),
		byte(n >> 8),
		byte(n),
	}

	_, _ = d.hh.Write(data[:])
}

func (d *defaultHasher) Write(p []byte) (n int, err error) {
	return d.hh.Write(p)
}

// Sum appends the current hash to b and returns the resulting slice.
// It does not change the underlying hash state.
func (d *defaultHasher) Sum(b []byte) []byte {
	return d.hh.Sum(b)
}

// Reset resets the Hash to its initial state.
func (d *defaultHasher) Reset() {
	d.hh.Reset()
}

// Size returns the number of bytes Sum will return.
func (d *defaultHasher) Size() int {
	return d.hh.Size()
}

// BlockSize returns the hash's underlying block size.
// The Write method must be able to accept any amount
// of data, but it may operate more efficiently if all writes
// are a multiple of the block size.
func (d *defaultHasher) BlockSize() int {
	return d.hh.BlockSize()
}

func (d *defaultHasher) Sum64() uint64 {
	return d.hh.Sum64()
}

func (d *defaultHasher) WriteByte(b byte) error {
	return d.hh.WriteByte(b)
}

func (d *defaultHasher) WriteString(s string) (int, error) {
	return d.hh.WriteString(s)
}

// WriteInterface implements Hasher.
func (d *defaultHasher) WriteInterface(v interface{}) {
	d.reflectWrite(reflect.ValueOf(v))
}

func (d *defaultHasher) reflectWrite(val reflect.Value) {
	// we write the type name first to ensure prefix-freedom
	// but instead of a (type name) string we use the address since
	// it will be identical for the same type
	typ := val.Type()
	d.hh.WriteString(typ.String())

	valKind := val.Kind()

	if val.IsValid() {
		impl, ok := val.Interface().(Hashable)
		if ok {
			impl.Hash(d)
			return
		}
	}

	switch valKind {
	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
		d.WriteUint64(uint64(val.Int()))

	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint, reflect.Uintptr:
		d.WriteUint64(uint64(val.Uint()))

	case reflect.Array, reflect.Slice:
		for i := range val.Len() {
			// prevent hashing to the same value
			// [2]string{"foo", ""} and [2]string{"", "foo"}.
			d.WriteUint64(uint64(i))

			d.reflectWrite(val.Index(i))
		}

	case reflect.String:
		d.WriteUint64(uint64(val.Len()))
		d.WriteString(val.String())

	case reflect.Struct:
		for i := range typ.NumField() {
			// ensure prefix-freedom
			d.WriteUint64(uint64(i))

			// skip all non-exported fields
			fldTyp := typ.Field(i)
			if !fldTyp.IsExported() {
				continue
			}

			d.reflectWrite(val.Field(i))
		}

	case reflect.Complex64, reflect.Complex128:
		c := val.Complex()
		d.WriteFloat64(real(c))
		d.WriteFloat64(imag(c))

	case reflect.Float32, reflect.Float64:
		d.WriteFloat64(val.Float())

	case reflect.Bool:
		if val.Bool() {
			d.WriteUint64(1)
		} else {
			d.WriteUint64(0)
		}

	case reflect.Interface:
		if !val.CanAddr() {
			// we cannot hash it. This may be the case if the field is a interface
			// with value nil.
			return
		}

		d.WriteUint64(uint64(val.UnsafeAddr()))

		if !val.IsNil() {
			d.reflectWrite(val.Elem())
		}

	case reflect.Pointer:
		// write the address to ensure prefix-freedom, even if the value is nil
		d.WriteUint64(uint64(uintptr(val.UnsafePointer())))

		if val.IsNil() {
			return
		}

		d.reflectWrite(val.Elem())

	case reflect.Chan, reflect.Func, reflect.UnsafePointer:
		d.WriteUint64(uint64(uintptr(val.UnsafePointer())))

	case reflect.Map:
		mapLen := val.Len()
		keys := getMapKeys(mapLen)

		for i, iter := 0, val.MapRange(); i < mapLen && iter.Next(); i++ {
			keyVal := iter.Key()

			(*keys)[i] = mapKeyEntry{
				key: keyVal.String(),
				val: keyVal,
			}
		}

		keys.Sort()

		for i, key := range *keys {
			// to be prefix-free, we write the key index
			d.WriteUint64(uint64(i))

			d.hh.WriteString(key.key)

			d.reflectWrite(val.MapIndex(key.val))
		}

	default:
		panic("xhash.defaultHasher: type " + val.Type().String() + " not supported")
	}
}
