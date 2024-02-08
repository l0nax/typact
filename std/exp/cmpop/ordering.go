package cmpop

// NOTE: There are two reasons why I added this to the package:
// 1. I wanted compare operations to look more natural
// 2. (the most crucial) reading comparsion code would be mutch easier to understand and less error prone.

// Ordering represents the ordering of a comparsion operation
// between two values.
//
// It is basically a helper type/ syntactic suggar.
type Ordering = int

const (
	// Less represents the case when a compared value is less than another.
	Less Ordering = -1
	// Equal represents the case when a compared value is equal to another.
	Equal Ordering = 0
	// Greater represents the case when a compared value is greater than another.
	Greater Ordering = 1
)
