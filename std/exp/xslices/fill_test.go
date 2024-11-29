package xslices

import (
	"fmt"
	"reflect"
	"testing"
)

var fillBenchmarkStages = []int{
	1, 2, 5, 10, 20, 100, 1_000,
}

type basicStruct struct {
	a int
	b string
	c bool
	d float64
}

func TestFill(t *testing.T) {
	type testData struct {
		name string

		data []int
	}

	data := []testData{
		{
			name: "empty",
			data: []int{},
		},
		{
			name: "one",
			data: []int{1},
		},
		{
			name: "even_entries",
			data: []int{1, 1, 1, 1},
		},
		{
			name: "uneven_entries",
			data: []int{1, 1, 1, 1, 1},
		},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			dest := make([]int, len(d.data))
			Fill(dest, 1)

			if !reflect.DeepEqual(dest, d.data) {
				t.Errorf("expected %v, got %v", d.data, dest)
			}
		})
	}
}

func BenchmarkFill(b *testing.B) {
	for _, stage := range fillBenchmarkStages {
		// ------ NATIVE
		b.Run(fmt.Sprintf("type=Native/Int/%d", stage), func(b *testing.B) {
			benchmarkIntRange(b, stage)
		})

		b.Run(fmt.Sprintf("type=Native/SmallString/%d", stage), func(b *testing.B) {
			benchmarkStringRange(b, stage)
		})

		b.Run(fmt.Sprintf("type=Native/Struct/%d", stage), func(b *testing.B) {
			benchmarkStructRange(b, stage)
		})

		b.Run(fmt.Sprintf("type=Native/Ptr/%d", stage), func(b *testing.B) {
			benchmarkStructPtrRange(b, stage)
		})

		// ------ Fill(...)
		b.Run(fmt.Sprintf("type=Fill/Int/%d", stage), func(b *testing.B) {
			benchmarkIntFill(b, stage)
		})

		b.Run(fmt.Sprintf("type=Fill/SmallString/%d", stage), func(b *testing.B) {
			benchmarkStringFill(b, stage)
		})

		b.Run(fmt.Sprintf("type=Fill/Struct/%d", stage), func(b *testing.B) {
			benchmarkStructFill(b, stage)
		})

		b.Run(fmt.Sprintf("type=Fill/Ptr/%d", stage), func(b *testing.B) {
			benchmarkStructPtrFill(b, stage)
		})
	}
}

func benchmarkIntFill(b *testing.B, stage int) {
	dest := make([]int, stage)

	b.ReportAllocs()

	for i := range b.N {
		Fill(dest, i)
	}
}

func benchmarkStringFill(b *testing.B, stage int) {
	const strVal = "Hello World"

	dest := make([]string, stage)

	b.ResetTimer()

	for range b.N {
		Fill(dest, strVal)
	}
}

func benchmarkIntRange(b *testing.B, stage int) {
	dest := make([]int, stage)

	b.ReportAllocs()

	for i := range b.N {
		for j := range stage {
			dest[j] = i
		}
	}
}

func benchmarkStringRange(b *testing.B, stage int) {
	const strVal = "Hello World"

	dest := make([]string, stage)

	b.ReportAllocs()

	for range b.N {
		for j := range stage {
			dest[j] = strVal
		}
	}
}

func benchmarkStructFill(b *testing.B, stage int) {
	data := basicStruct{
		a: 1,
		b: "Hello World",
		c: true,
		d: 3.14,
	}
	dest := make([]basicStruct, stage)

	b.ReportAllocs()

	for range b.N {
		Fill(dest, data)
	}
}

func benchmarkStructRange(b *testing.B, stage int) {
	data := basicStruct{
		a: 1,
		b: "Hello World",
		c: true,
		d: 3.14,
	}
	dest := make([]basicStruct, stage)

	b.ReportAllocs()

	for range b.N {
		for j := range stage {
			dest[j] = data
		}
	}
}

func benchmarkStructPtrFill(b *testing.B, stage int) {
	data := &basicStruct{
		a: 1,
		b: "Hello World",
		c: true,
		d: 3.14,
	}
	dest := make([]*basicStruct, stage)

	b.ReportAllocs()

	for range b.N {
		Fill(dest, data)
	}
}

func benchmarkStructPtrRange(b *testing.B, stage int) {
	data := &basicStruct{
		a: 1,
		b: "Hello World",
		c: true,
		d: 3.14,
	}
	dest := make([]*basicStruct, stage)

	b.ReportAllocs()

	for range b.N {
		for j := range stage {
			dest[j] = data
		}
	}
}
