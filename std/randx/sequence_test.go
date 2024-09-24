package randx

import "testing"

func BenchmarkRuneSequence(b *testing.B) {
	b.ReportAllocs()

	for range b.N {
		_, err := RuneSequence(64, AlphaNum)
		if err != nil {
			b.Fatal(err)
		}
	}
}
