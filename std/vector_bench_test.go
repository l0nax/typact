package std_test

import (
	"testing"

	"go.l0nax.org/typact/std"
)

func BenchmarkVector_AppendVector(b *testing.B) {
	b.ReportAllocs()

	genVec := func(n int) std.Vector[int] {
		tmp := make([]int, n)
		for i := 0; i<n; i++ {
			tmp[i] = i
		}

		return tmp
	}

	b.Run("10", func(b *testing.B) {
		other := std.VectorFromSlice(genVec(10))
		for i := 0; i < b.N; i++ {
			vec := std.VectorFromSlice(genVec(10))
			vec.AppendVector(other)
			_ = vec
		}

	})

	b.Run("512", func(b *testing.B) {
		other := std.VectorFromSlice(genVec(512))
		for i := 0; i < b.N; i++ {
			vec := std.VectorFromSlice(genVec(512))
			vec.AppendVector(other)
			_ = vec
		}
	})

	b.Run("1024", func(b *testing.B) {
		other := std.VectorFromSlice(genVec(1024))
		for i := 0; i < b.N; i++ {
			vec := std.VectorFromSlice(genVec(1024))
			vec.AppendVector(other)
			_ = vec
		}
	})

	b.Run("10240", func(b *testing.B) {
		other := std.VectorFromSlice(genVec(10240))
		for i := 0; i < b.N; i++ {
			vec := std.VectorFromSlice(genVec(10240))
			vec.AppendVector(other)
			_ = vec
		}
	})
}

func BenchmarkVector_Clone(b *testing.B) {
	b.ReportAllocs()

	genVec := func(n int) std.Vector[int] {
		tmp := make([]int, n)
		for i := 0; i<n; i++ {
			tmp[i] = i
		}

		return tmp
	}

	b.Run("10", func(b *testing.B) {
		vec := std.VectorFromSlice(genVec(10))
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			tmp := vec.Clone()
			_ = tmp
		}
	})

	b.Run("512", func(b *testing.B) {
		vec := std.VectorFromSlice(genVec(512))
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			tmp := vec.Clone()
			_ = tmp
		}
	})

	b.Run("1024", func(b *testing.B) {
		vec := std.VectorFromSlice(genVec(1024))
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			tmp := vec.Clone()
			_ = tmp
		}
	})
}
