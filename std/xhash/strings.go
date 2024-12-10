package xhash

import (
	"cmp"
	"reflect"
	"slices"
	"sync"

	"go.l0nax.org/typact/std/exp/cmpop"
)

var stringsPools = &sync.Pool{
	New: func() any {
		return new(mapKeySlice)
	},
}

type mapKeyEntry struct {
	key string
	val reflect.Value
}

type mapKeySlice []mapKeyEntry

// getMapKeys returns a non-nil pointer to a slice with length n.
func getMapKeys(n int) *mapKeySlice {
	s := stringsPools.Get().(*mapKeySlice)
	if cap(*s) < n {
		*s = nil // release the memory, if any, to reduce GC pressure
		*s = make([]mapKeyEntry, n)
	}

	*s = (*s)[:n]

	return s
}

func putStrings(s *mapKeySlice) {
	if cap(*s) > 1<<10 {
		*s = nil // avoid pinning arbitrarily large amounts of memory
	}
	stringsPools.Put(s)
}

// Sort sorts the string slice according to RFC 8785, section 3.2.3.
func (ss *mapKeySlice) Sort() {
	slices.SortFunc(*ss, func(a, b mapKeyEntry) cmpop.Ordering {
		return cmp.Compare(a.key, b.key)
	})
}
