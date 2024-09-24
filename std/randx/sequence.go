package randx

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

var rander = rand.Reader // random function

var (
	// AlphaNum contains runes [abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789].
	AlphaNum = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	// Alpha contains runes [abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ].
	Alpha = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	// AlphaLowerNum contains runes [abcdefghijklmnopqrstuvwxyz0123456789].
	AlphaLowerNum = []rune("abcdefghijklmnopqrstuvwxyz0123456789")
	// AlphaUpperNum contains runes [ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789].
	AlphaUpperNum = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	// AlphaLower contains runes [abcdefghijklmnopqrstuvwxyz].
	AlphaLower = []rune("abcdefghijklmnopqrstuvwxyz")
	// AlphaUpper contains runes [ABCDEFGHIJKLMNOPQRSTUVWXYZ].
	AlphaUpper = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	// Numeric contains runes [0123456789].
	Numeric = []rune("0123456789")
)

// RuneSequence returns a cryptographically secure random
// sequence using the defined allowed runes.
func RuneSequence(l int, allowedRunes []rune) (seq []rune, err error) {
	maxLen := uint64(len(allowedRunes))
	seq = make([]rune, l)

	for i := 0; i < l; i++ {
		r, err := Uint64N(maxLen)
		if err != nil {
			return seq, fmt.Errorf("unable to genrate random number: %w", err)
		}

		rn := allowedRunes[r]
		seq[i] = rn
	}

	return seq, nil
}

// MustString returns a cryptographically secure random
// string sequence using the defined runes.
//
// Panics on error.
func MustString(l int, allowedRunes []rune) string {
	seq, err := RuneSequence(l, allowedRunes)
	if err != nil {
		panic(err)
	}

	return string(seq)
}

// MustNumeric returns a cryptographically secure random number in the range of [0, num).
func MustNumeric(num int) int {
	n, err := rand.Int(rand.Reader, big.NewInt(int64(num)))
	if err != nil {
		panic(err)
	}

	return int(n.Int64())
}
