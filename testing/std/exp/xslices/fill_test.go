package xslices

import (
	"reflect"
	"testing"

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
