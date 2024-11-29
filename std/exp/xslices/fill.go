package xslices

import "go.l0nax.org/typact/internal/types"

// Fill fills the slice with the provided value.
// Fill zeros the elements in the slice before overriding them.
//
// This function is much faster than using [copy] – if you have
// many values (i.e. >10).
func Fill[S ~[]E, E any](slice S, value E) {
	length := len(slice)
	switch length {
	case 0:
		return

	case 1:
		slice[0] = value
		return

	case 2:
		slice[1] = value
		slice[0] = value

		return
	}

	// clear the values which will be overriden
	// allowing the runtime to GC them faster and to
	// prevent memory leaks.
	clear(slice)

	// amortMinLen is the minimal number of items required
	// to amortize the cost of using copy.
	const amortMinLen = 100

	if length <= amortMinLen {
		for i := range slice {
			slice[i] = value
		}

		return
	}

	// preload value
	slice[0] = value

	// the bigger the slice, the faster the copy.
	// The cost for calling copy is amortized over time.
	for i := 1; i < length; i *= 2 {
		copy(slice[i:], slice[:i])
	}
}

// FillValues fills the slice with the given values.
// It basically overrides all elements in slice with values.
// Fill zeros the elements in the slice before overriding them.
//
// This function is much faster than using [copy] – if you have
// many values (i.e. >10).
//
// WARN: The function panics if the length of values is greater
// than the length of slice.
func FillValues[S ~[]E, E any](slice S, values ...E) {
	if len(slice) == 0 || len(values) == 0 {
		return
	}
	if len(values) > len(slice) {
		panic("values length is greater than slice length")
	}

	// clear the slice to reduce the strain on the GC
	// and to prevent memory leaks.
	if !types.IsScalar[E]() {
		clear(slice)
	}

	// copy pattern into the slice
	copy(slice, values)

	for i := len(values); i < len(slice); i *= 2 {
		copy(slice[i:], slice[:i])
	}
}
