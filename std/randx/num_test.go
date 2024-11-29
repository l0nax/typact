package randx

import "testing"

var AlwaysFalse = false

func keep[T int | uint | int32 | uint32 | int64 | uint64](x T) T {
	if AlwaysFalse {
		return -x
	}
	return x
}

func BenchmarkUint64(b *testing.B) {
	b.ReportAllocs()

	for range b.N {
		_, err := Uint64()
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkUint64N(b *testing.B) {
	b.ReportAllocs()

	b.Run("N=1000", func(b *testing.B) {
		b.ReportAllocs()

		arg := keep(uint64(1_000))

		for range b.N {
			_, err := Uint64N(arg)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("N=1e8", func(b *testing.B) {
		b.ReportAllocs()

		arg := keep(uint64(1e8))

		for range b.N {
			_, err := Uint64N(arg)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("N=1e9", func(b *testing.B) {
		b.ReportAllocs()

		arg := keep(uint64(1e9))

		for range b.N {
			_, err := Uint64N(arg)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("N=1e18", func(b *testing.B) {
		b.ReportAllocs()

		arg := keep(uint64(1e18))

		for range b.N {
			_, err := Uint64N(arg)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}
