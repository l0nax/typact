package randx

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.l0nax.org/typact/std/randx"
)

func TestRunePatterns(t *testing.T) {
	for k, v := range []struct {
		runes       []rune
		shouldMatch string
	}{
		{randx.Alpha, "[a-zA-Z]{52}"},
		{randx.AlphaLower, "[a-z]{26}"},
		{randx.AlphaUpper, "[A-Z]{26}"},
		{randx.AlphaNum, "[a-zA-Z0-9]{62}"},
		{randx.AlphaLowerNum, "[a-z0-9]{36}"},
		{randx.AlphaUpperNum, "[A-Z0-9]{36}"},
		{randx.Numeric, "[0-9]{10}"},
	} {
		valid, err := regexp.Match(v.shouldMatch, []byte(string(v.runes)))
		assert.Nil(t, err, "Case %d", k)
		assert.True(t, valid, "Case %d", k)
	}
}

func TestRuneSequenceMatchesPattern(t *testing.T) {
	for k, v := range []struct {
		runes       []rune
		shouldMatch string
		length      int
	}{
		{randx.Alpha, "[a-zA-Z]+", 25},
		{randx.AlphaLower, "[a-z]+", 46},
		{randx.AlphaUpper, "[A-Z]+", 21},
		{randx.AlphaNum, "[a-zA-Z0-9]+", 123},
		{randx.AlphaLowerNum, "[a-z0-9]+", 41},
		{randx.AlphaUpperNum, "[A-Z0-9]+", 94914},
		{randx.Numeric, "[0-9]+", 94914},
	} {
		seq, err := randx.RuneSequence(v.length, v.runes)
		assert.Nil(t, err, "case %d", k)
		assert.Equal(t, v.length, len(seq), "case %d", k)

		valid, err := regexp.Match(v.shouldMatch, []byte(string(seq)))
		assert.Nil(t, err, "case %d", k)
		assert.True(t, valid, "case %d\nrunes %s\nresult %s", k, v.runes, string(seq))
	}
}

func TestRuneSequenceIsPseudoUnique(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}

	times := 100
	runes := []rune("ab")
	length := 32
	s := make(map[string]bool)

	for i := 0; i < times; i++ {
		k, err := randx.RuneSequence(length, runes)
		assert.Nil(t, err)

		ks := string(k)

		_, ok := s[ks]
		assert.False(t, ok)
		if ok {
			return
		}

		s[ks] = true
	}
}

func BenchmarkTestInt64(b *testing.B) {
	length := 25
	pattern := []rune("abcdefghijklmnopqrstuvwxyz")

	for i := 0; i < b.N; i++ {
		randx.RuneSequence(length, pattern)
	}
}
