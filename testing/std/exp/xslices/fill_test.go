package xslices

import (
	"math/bits"
	"reflect"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.l0nax.org/typact/std/exp/xslices"
)

func getSlice() []byte {
	return make([]byte, 73437)
}

func getIntSlice() []uint64 {
	return make([]uint64, 73437)
}

func TestFillSlice(t *testing.T) {
	t.Parallel()

	t.Run("nil slice", func(t *testing.T) {
		var src []int
		xslices.Fill(src, 0)
	})

	t.Run("no entry", func(t *testing.T) {
		src := []int{}
		xslices.Fill(src, 0)

		if len(src) != 0 {
			t.Errorf("Expected len 0, got %d", len(src))
		}
	})

	t.Run("two entries", func(t *testing.T) {
		src := []int{5, 100}
		xslices.Fill(src, 1024)

		exp := []int{1024, 1024}
		if !reflect.DeepEqual(src, exp) {
			t.Errorf("Expected %v, got %v", exp, src)
		}
	})

	t.Run("normal slice", func(t *testing.T) {
		src := []int{5, 100, 12315, 345634, 34534}
		xslices.Fill(src, 1)

		expect := []int{1, 1, 1, 1, 1}
		if !reflect.DeepEqual(src, expect) {
			t.Errorf("Expected %v, got %v", expect, src)
		}
	})

	t.Run("append", func(t *testing.T) {
		src := make([]int, 0, 10)
		src = append(src, 5, 100, 12315, 345634, 34534)
		xslices.Fill(src, 1)

		expect := []int{1, 1, 1, 1, 1}
		if !reflect.DeepEqual(src, expect) {
			t.Errorf("Expected %v, got %v", expect, src)
		}
	})
}

func BenchmarkFillSlice_Index(b *testing.B) {
	b.Run("byte", func(b *testing.B) {
		bigSlice := getSlice()

		b.ResetTimer()
		b.ReportAllocs()

		for range b.N {
			for i := 0; i < len(bigSlice); i++ {
				bigSlice[i] = 65
			}
		}
	})

	b.Run("uint64", func(b *testing.B) {
		bigSlice := getIntSlice()

		b.ResetTimer()
		b.ReportAllocs()

		for range b.N {
			for i := 0; i < len(bigSlice); i++ {
				bigSlice[i] = 65
			}
		}
	})
}

func BenchmarkFillSlice_Range(b *testing.B) {
	b.Run("byte", func(b *testing.B) {
		bigSlice := getSlice()

		b.ResetTimer()
		b.ReportAllocs()

		for range b.N {
			for i := range bigSlice {
				bigSlice[i] = 66
			}
		}
	})

	b.Run("uint64", func(b *testing.B) {
		bigSlice := getIntSlice()

		b.ResetTimer()
		b.ReportAllocs()

		for range b.N {
			for i := range bigSlice {
				bigSlice[i] = 66
			}
		}
	})
}

func BenchmarkFill(b *testing.B) {
	b.Run("byte", func(b *testing.B) {
		data := make([]byte, 73437)

		b.ResetTimer()
		b.ReportAllocs()

		for range b.N {
			xslices.Fill(data, 66)
		}
	})

	b.Run("uint64", func(b *testing.B) {
		data := make([]uint64, 73437)

		b.ResetTimer()
		b.ReportAllocs()

		for range b.N {
			xslices.Fill(data, 66)
		}
	})
}

var _ = Describe("Fill", func() {
	Describe("FillValues", func() {
		It("should correctly fill a perfect slice", func() {
			pattern := []int{1, 2, 3, 4}
			slice := make([]int, 4*8)
			expect := repeat(pattern, 8)

			xslices.FillValues(slice, pattern...)
			Expect(slice).To(Equal(expect))
		})

		It("should fill if slice has an uneven length", func() {
			pattern := []int{1, 2, 3, 4}
			slice := make([]int, 4*8+1)

			expect := repeat(pattern, 8)
			expect = append(expect, 1)

			xslices.FillValues(slice, pattern...)
			Expect(slice).To(Equal(expect))
		})
	})
})

// repeat returns a new slice that repeats the provided slice the given number of times.
// The result has length and capacity len(x) * count.
// The result is never nil.
// Repeat panics if count is negative or if the result of (len(x) * count)
// overflows.
//
// NOTE: This has been copied from the GoLang source code!
// TODO: Replace me with [slices.Repeat] once go 1.23 has been released!
func repeat[S ~[]E, E any](x S, count int) S {
	if count < 0 {
		panic("cannot be negative")
	}

	const maxInt = ^uint(0) >> 1
	if hi, lo := bits.Mul(uint(len(x)), uint(count)); hi > 0 || lo > maxInt {
		panic("the result of (len(x) * count) overflows")
	}

	newslice := make(S, len(x)*count)
	n := copy(newslice, x)
	for n < len(newslice) {
		n += copy(newslice[n:], newslice[:n])
	}
	return newslice
}
