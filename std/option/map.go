package option

import "go.l0nax.org/typact"

// Map maps Option[T] to Option[K] if the value of src is present.
// Otherwise [typact.None] will be returned.
func Map[T any, K any](src typact.Option[T], mapFn func(T) K) typact.Option[K] {
	if src.IsSome() {
		return typact.Some(mapFn(src.UnsafeUnwrap()))
	}

	return typact.None[K]()
}

// MapOr maps src with mapFn if the value is present.
// Otherwise the provided value is returned.
func MapOr[T any, K any](src typact.Option[T], mapFn func(T) K, value K) K {
	if src.IsSome() {
		return mapFn(src.UnsafeUnwrap())
	}

	return value
}

// MapOrElse applies the function mapFn to the value held by src if it exists,
// and returns the result. If src does not hold a value, it applies valueFn and returns its result.
//
// This allows conditional transformation of the Option's value or generation of a default value when none is present.
func MapOrElse[T any, K any](src typact.Option[T], mapFn func(T) K, valueFn func() K) K {
	if src.IsSome() {
		return mapFn(src.UnsafeUnwrap())
	}

	return valueFn()
}

