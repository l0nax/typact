package xslices

import (
	"testing"

	"go.l0nax.org/typact/std/exp/xslices"
)


func getSlice() []byte {
	return make([]byte, 73437)
}

func getIntSlice() []uint64 {
	return make([]uint64, 73437)
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
