package randx

import (
	"crypto/rand"
	"fmt"
	"math/bits"
)

// is32bit tells whether the system is a 32 bit system or not.
const is32bit = ^uint(0)>>32 == 0

// IntN returns a uniform random value in [0, max).
// It panics if max <= 0.
func IntN(max int) (int, error) {
	num, err := Uint64N(uint64(max))
	if err != nil {
		return 0, err
	}

	return int(num), nil
}

// Uint64N returns a uniform random value in [0, max).
// It panics if max <= 0.
func Uint64N(max uint64) (uint64, error) {
	if max <= 0 {
		panic("Argument max must be > 0")
	}

	num, err := Uint64()
	if err != nil {
		return 0, nil
	}

	if max&(max-1) == 0 { // n is power of two, can mask
		return num & (max - 1), nil
	}

	// NOTE: This base on the code in math/rand/v2.
	// Copyright 2009 The Go Authors
	//
	// See https://arxiv.org/abs/1805.10941
	//
	// Suppose we have a uint64 x uniform in the range [0,2⁶⁴)
	// and want to reduce it to the range [0,n) preserving exact uniformity.
	// We can simulate a scaling arbitrary precision x * (n/2⁶⁴) by
	// the high bits of a double-width multiply of x*n, meaning (x*n)/2⁶⁴.
	// Since there are 2⁶⁴ possible inputs x and only n possible outputs,
	// the output is necessarily biased if n does not divide 2⁶⁴.
	// In general (x*n)/2⁶⁴ = k for x*n in [k*2⁶⁴,(k+1)*2⁶⁴).
	// There are either floor(2⁶⁴/n) or ceil(2⁶⁴/n) possible products
	// in that range, depending on k.
	// But suppose we reject the sample and try again when
	// x*n is in [k*2⁶⁴, k*2⁶⁴+(2⁶⁴%n)), meaning rejecting fewer than n possible
	// outcomes out of the 2⁶⁴.
	// Now there are exactly floor(2⁶⁴/n) possible ways to produce
	// each output value k, so we've restored uniformity.
	// To get valid uint64 math, 2⁶⁴ % n = (2⁶⁴ - n) % n = -n % n,
	// so the direct implementation of this algorithm would be:
	//
	//	hi, lo := bits.Mul64(r.Uint64(), n)
	//	thresh := -n % n
	//	for lo < thresh {
	//		hi, lo = bits.Mul64(r.Uint64(), n)
	//	}
	//
	// That still leaves an expensive 64-bit division that we would rather avoid.
	// We know that thresh < n, and n is usually much less than 2⁶⁴, so we can
	// avoid the last four lines unless lo < n.
	//
	// See also:
	// https://lemire.me/blog/2016/06/27/a-fast-alternative-to-the-modulo-reduction
	// https://lemire.me/blog/2016/06/30/fast-random-shuffling
	hi, lo := bits.Mul64(num, max)
	if lo < max {
		thresh := -max % max
		for lo < thresh {
			num, err = Uint64()
			if err != nil {
				return 0, nil
			}

			hi, lo = bits.Mul64(num, max)
		}
	}

	return hi, nil
}

// Int returns a non-negative uniform random value as int.
func Int() (int, error) {
	num, err := Uint64()
	if err != nil {
		return 0, err
	}

	return int(uint(num) << 1 >> 1), nil
}

// Uint64 returns a uniform random value as uint64.
func Uint64() (uint64, error) {
	var data [8]byte

	_, err := rand.Read(data[:])
	if err != nil {
		return 0, fmt.Errorf("unable to read random data: %w", err)
	}

	return uint64(data[0]) |
		uint64(data[1])<<8 |
		uint64(data[2])<<16 |
		uint64(data[3])<<24 |
		uint64(data[4])<<32 |
		uint64(data[5])<<40 |
		uint64(data[6])<<48 |
		uint64(data[7])<<56, nil
}

// Uint32N returns a uniform random value in [0, max).
// It panics if max <= 0.
func Uint32N(max uint32) (uint32, error) {
	num, err := Uint64N(uint64(max))
	if err != nil {
		return 0, err
	}

	return uint32(num), nil
}

// Uint32 returns a uniform random value as uint32.
func Uint32() (uint32, error) {
	num, err := Uint64()
	if err != nil {
		return 0, err
	}

	return uint32(num >> 32), nil
}
