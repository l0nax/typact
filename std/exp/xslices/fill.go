package xslices

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
	// allowing the runtime to GC them faster
	clear(slice[1:])

	for i := 1; i<len(slice); i *= 2 {
		copy(slice[i:], slice[:i])
	}
}
