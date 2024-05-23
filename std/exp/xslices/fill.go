package xslices

import "go.l0nax.org/typact/internal/types"

// Fill fills the slice with the provided value.
// Fill zeros the elements in the slice before overriding them.
func Fill[S ~[]E, E any](slice S, value E) {
	if len(slice) == 0 {
		return
	}
	if len(slice) == 1 {
		slice[0] = value
		return
	}

	// preload value
	slice[0] = value

	// clear the values which will be overriden
	// allowing the runtime to GC them faster and to
	// prevent memory leaks.
	if !types.IsScalar[E]() {
		clear(slice[1:])
	}

	// the bigger the slice, the faster the copy.
	// The cost for calling copy is amortized over time.
	for i := 1; i < len(slice); i *= 2 {
		copy(slice[i:], slice[:i])
	}
}
