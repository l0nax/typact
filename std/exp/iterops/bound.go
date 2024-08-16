package iterops

import "go.l0nax.org/typact"

// NOTE: This file contains an implementation/ port of
// std::ops::Bound from Rust, because I like the concept.

// BoundType describes the type of the range.
type BoundType int8

const (
	// BoundIncluded describes that the key should be included.
	//
	// Math equivalent to: `[key, key]`
	BoundIncluded BoundType = iota + 1
	// BoundExcluded describes that the key should be excluded.
	//
	// Math equivalent to: `(key, key]`
	BoundExcluded
	// BoundUnbounded describes no bound.
	BoundUnbounded
)

// Bound describes a range of keys in an iterator.
type Bound[T any] struct {
	key       typact.Option[T]
	boundType BoundType
}

// Included returns an [BoundIncluded] [Bound].
func Included[T any](key T) Bound[T] {
	return Bound[T]{
		key:       typact.Some(key),
		boundType: BoundIncluded,
	}
}

// Excluded returns an [BoundExcluded] [Bound].
func Excluded[T any](key T) Bound[T] {
	return Bound[T]{
		key:       typact.Some(key),
		boundType: BoundExcluded,
	}
}

// Unbounded returns an [BoundUnbounded] [Bound].
func Unbounded[T any]() Bound[T] {
	return Bound[T]{
		boundType: BoundUnbounded,
	}
}

// Key returns the key of b.
// It is only available if BoundType is either
// [BoundIncluded] or [BoundExcluded].
func (b Bound[T]) Key() typact.Option[T] {
	return b.key
}

// BoundType returns the [BoundType] of b.
func (b Bound[T]) BoundType() BoundType {
	return b.boundType
}
