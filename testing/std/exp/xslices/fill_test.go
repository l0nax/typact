package xslices

import (
	"testing"

	"go.l0nax.org/typact/std/exp/xslices"
)

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
