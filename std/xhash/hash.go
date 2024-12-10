package xhash

import (
	"hash"
	"reflect"
)

// hashableImpl holds the [reflect.Type] of [Hashable].
var hashableImpl = reflect.TypeOf((Hashable)(nil))

// Hashable is implemented by custom types to allow them to be hashed.
//
// ## Prefix collisions
//
// When implementing this interface, the implementee must ensure that the data passed
// to any call of [Hasher] is prefix free.
// This means that values which are not equal should cause two different sequences of
// values to be written, and neither of the two sequences should be a prefix of the other.
//
// In simple words: if a struct hash two string fields and [Hasher.WriteString] is called
// successively, the caller must ensure that the resulting hash is unique.
// I.e., the string tuples ("ab", "c") and ("a", "bc") must result in a different hash.
type Hashable interface {
	// Hash uses h to hash its value.
	Hash(h Hasher)
}

// Hasher is the hashing implementation.
// For the default implementation use [NewHasher].
type Hasher interface {
	hash.Hash64

	// WriteByte writes the single byte b into the hasher.
	WriteByte(b byte) error

	// WriteString writes the given string into the hasher.
	//
	// WARN: This method does not ensure that the input is prefix-free!
	WriteString(s string) (int, error)

	// WriteInt writes the given integer into the hasher.
	WriteInt(n int)

	// WriteUint64 writes the given integer into the hasher.
	WriteUint64(n uint64)

	// WriteFloat64 writes the given float into the hasher.
	WriteFloat64(n float64)

	// WriteInterface writes the given interface into the hasher.
	// If the type implements [Hashable] it will be hashed recursively.
	//
	// This method ensures that the hashed data is prefix-free!
	// If v is a struct, its fields will be hashed recursively – even non-exported ones.
	WriteInterface(v interface{})
}
