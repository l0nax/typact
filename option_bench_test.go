package typact_test

import (
	"os"
	"testing"
	"time"

	"go.l0nax.org/typact"
)

type SmallStruct struct {
	Data    string
	Raw     []byte
	Age     int
	Boolean bool
}

type OptionalHolder struct {
	Wrapper typact.Option[*ScalarWrapper]
}

type ScalarWrapper struct {
	Number uint32
}

func BenchmarkOption_GetOrInsert(b *testing.B) {
	b.ReportAllocs()

	b.Run("None", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			opt := typact.None[int]()

			x := opt.GetOrInsert(5)
			_ = x
		}
	})

	b.Run("Some", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			opt := typact.Some(5)

			x := opt.GetOrInsert(5)
			_ = x
		}
	})
}

func BenchmarkOption_Insert(b *testing.B) {
	b.ReportAllocs()

	// NOTE: We can reuse opt since Insert will always override
	// the value

	b.Run("None", func(b *testing.B) {
		opt := typact.None[int]()

		for i := 0; i < b.N; i++ {
			x := opt.Insert(5)
			_ = x
		}
	})

	b.Run("Some", func(b *testing.B) {
		opt := typact.Some(5)

		for i := 0; i < b.N; i++ {
			x := opt.Insert(10)
			_ = x
		}
	})
}

// BenchmarkOrElseAndThenSome benchmarks a real-world usage of
// OrElse(..).AndThen(...)
func BenchmarkOption_OrElseAndThenSome(b *testing.B) {
	ref := &OptionalHolder{
		Wrapper: typact.Some(&ScalarWrapper{
			Number: 21324534,
		}),
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		ref.Wrapper = ref.Wrapper.
			OrElse(func() typact.Option[*ScalarWrapper] {
				return typact.Some(&ScalarWrapper{
					Number: 456456,
				})
			}).
			AndThen(func(sw *ScalarWrapper) typact.Option[*ScalarWrapper] {
				sw.Number = 76575
				return typact.Some(sw)
			})
	}
}

func BenchmarkOption_OrElseAndThenNone(b *testing.B) {
	ref := &OptionalHolder{
		Wrapper: typact.None[*ScalarWrapper](),
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		ref.Wrapper = ref.Wrapper.
			OrElse(func() typact.Option[*ScalarWrapper] {
				return typact.Some(&ScalarWrapper{
					Number: 456456,
				})
			}).
			AndThen(func(sw *ScalarWrapper) typact.Option[*ScalarWrapper] {
				sw.Number = 76575
				return typact.Some(sw)
			})
	}
}

func getSomeStr() typact.Option[string] {
	if os.Getenv("NOT_EXISTING") == "" {
		return typact.Some("Foo")
	}

	return typact.None[string]()
}

func getNoneStr() typact.Option[string] {
	if os.Getenv("NOT_EXISTING") != "" {
		return typact.Some("Foo")
	}

	return typact.None[string]()
}

func BenchmarkOption_Unwrap(b *testing.B) {
	b.ReportAllocs()

	vv := getSomeStr()
	for i := 0; i < b.N; i++ {
		str := vv.Unwrap()
		_ = str
	}
}

func BenchmarkOption_Expect(b *testing.B) {
	b.ReportAllocs()

	vv := getSomeStr()
	for i := 0; i < b.N; i++ {
		str := vv.Expect("my string")
		_ = str
	}
}

func BenchmarkOption_IsSomeAnd(b *testing.B) {
	b.ReportAllocs()

	b.Run("Some", func(b *testing.B) {
		vv := getSomeStr()

		for i := 0; i < b.N; i++ {
			ok := vv.IsSomeAnd(func(str string) bool {
				return str == "Foo"
			})

			_ = ok
		}
	})

	b.Run("None", func(b *testing.B) {
		vv := getNoneStr()

		for i := 0; i < b.N; i++ {
			ok := vv.IsSomeAnd(func(str string) bool {
				return str == "Foo"
			})

			_ = ok
		}
	})

	b.Run("Some_Slice", func(b *testing.B) {
		var vv typact.Option[[]string]
		if os.Getenv("NOT_EXISTING") == "" {
			vv = typact.Some([]string{
				"Hello", "World", "Foo", "Bar",
				"Test", "Something", "Home",
			})
		}

		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			isLongEnough := vv.IsSomeAnd(func(strs []string) bool {
				n := 0
				for i := range strs {
					n += len(strs[i])
				}

				return n > 100
			})
			_ = isLongEnough
		}
	})
}

func BenchmarkOption_UnwrapOr(b *testing.B) {
	b.ReportAllocs()

	b.Run("Native", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			var str string
			if str == "" {
				str = "My String"
			}

			_ = str
		}
	})

	b.Run("Some", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			str := typact.Some("Foo").UnwrapOr("Bar")
			_ = str
		}
	})

	b.Run("None", func(b *testing.B) {
		str := typact.None[string]().UnwrapOr("Bar")
		_ = str
	})
}

func BenchmarkOption_UnwrapOrSome(b *testing.B) {
	v := typact.Some(&SmallStruct{
		Data:    "dflgködfslgkäsöfkgäösdkägöksdfäölgksdäfölgkä",
		Age:     321321654,
		Boolean: true,
		Raw:     []byte("sldfkädafgkäadfkgäölsfdkgäölsdkfgäö"),
	})

	optionUnwrapOr(b, v)
}

func BenchmarkOption_UnwrapOrNone(b *testing.B) {
	v := typact.None[*SmallStruct]()

	optionUnwrapOr(b, v)
}

func optionUnwrapOr(b *testing.B, o typact.Option[*SmallStruct]) {
	alter := &SmallStruct{
		Data: "This is my alternative!!",
		Raw:  []byte("aösädllfdgä#ödsfl#gäsldf#"),
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		x := o.UnwrapOr(alter)
		_ = x
	}
}

func BenchmarkOptionUnwrapOrElseSome(b *testing.B) {
	v := typact.Some(&SmallStruct{
		Data:    "dflgködfslgkäsöfkgäösdkägöksdfäölgksdäfölgkä",
		Age:     321321654,
		Boolean: true,
		Raw:     []byte("sldfkädafgkäadfkgäölsfdkgäölsdkfgäö"),
	})

	optionUnwrapOrElse(b, v)
}

func BenchmarkOption_UnwrapOrElseNone(b *testing.B) {
	v := typact.None[*SmallStruct]()

	optionUnwrapOrElse(b, v)
}

func optionUnwrapOrElse(b *testing.B, o typact.Option[*SmallStruct]) {
	alterFn := func() *SmallStruct {
		return &SmallStruct{
			Data: "This is my alternative!!",
			Raw:  []byte("aösädllfdgä#ödsfl#gäsldf#"),
		}
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		x := o.UnwrapOrElse(alterFn)
		_ = x
	}
}

func BenchmarkOption_UnwrapAsRef(b *testing.B) {
	str := typact.Some("this is my struct")
	slice := typact.Some([]string{"hello", "world"})

	b.ResetTimer()

	b.Run("Scalar", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			ref := str.UnwrapAsRef()
			_ = ref
		}
	})

	b.Run("Slice", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			ref := slice.UnwrapAsRef()
			_ = ref
		}
	})
}

func BenchmarkOption_AndThen(b *testing.B) {
	opt := typact.Some(5)
	fn := func(val int) typact.Option[int] {
		return typact.Some(val * 2)
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		ret := opt.AndThen(fn)
		_ = ret
	}
}

func BenchmarkOption_AndThenNone(b *testing.B) {
	opt := typact.None[int]()
	fn := func(val int) typact.Option[int] {
		return typact.Some(val * 2)
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		ret := opt.AndThen(fn)
		_ = ret
	}
}

type myData struct {
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (m *myData) Clone() *myData {
	return &myData{
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

type myStruct struct {
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (m myStruct) Clone() myStruct {
	return myStruct{
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

type (
	MySimpleDataPtrList   []*myData
	MySimpleDataList      []myData
	MySimpleStructPtrList []*myStruct
	MySimpleStructList    []myStruct
)

func BenchmarkOption_Clone(b *testing.B) {
	b.ReportAllocs()

	b.Run("None", func(b *testing.B) {
		val := typact.None[string]()

		for i := 0; i < b.N; i++ {
			tmp := val.Clone()
			_ = tmp
		}
	})

	b.Run("String", func(b *testing.B) {
		val := typact.Some("Foo Bar")

		for i := 0; i < b.N; i++ {
			tmp := val.Clone()
			_ = tmp
		}
	})

	b.Run("Int64", func(b *testing.B) {
		val := typact.Some(int64(123123))

		for i := 0; i < b.N; i++ {
			tmp := val.Clone()
			_ = tmp
		}
	})

	b.Run("CustomStructPointer", func(b *testing.B) {
		val := typact.Some(&myData{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		})

		for i := 0; i < b.N; i++ {
			tmp := val.Clone()
			_ = tmp
		}
	})

	b.Run("CustomStructFallback", func(b *testing.B) {
		// NOTE: The myData struct implements the std.Cloner interface with a pointer
		// receiver.
		val := typact.Some(myData{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		})

		for i := 0; i < b.N; i++ {
			tmp := val.Clone()
			_ = tmp
		}
	})

	b.Run("ScalarSlice", func(b *testing.B) {
		val := typact.Some([]string{
			"Foo", "Bar",
			"Hello", "World",
		})

		for i := 0; i < b.N; i++ {
			tmp := val.Clone()
			_ = tmp
		}
	})

	// This benchmarks
	//   3. T is []*E; E implements [std.Cloner] with a pointer receiver
	b.Run("PtrSliceWrapper_PtrRecv", func(b *testing.B) {
		val := typact.Some(MySimpleDataPtrList{
			{
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			{
				CreatedAt: time.Now().Add(time.Hour),
				UpdatedAt: time.Now().Add(time.Hour),
			},
		})

		b.ReportAllocs()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			tmp := val.Clone()
			_ = tmp
		}
	})

	// This benchmarks
	//   1. T is []E; E implements [std.Cloner] with a pointer receiver
	b.Run("SliceWrapper_PtrRecv", func(b *testing.B) {
		val := typact.Some(MySimpleDataList{
			{
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			{
				CreatedAt: time.Now().Add(time.Hour),
				UpdatedAt: time.Now().Add(time.Hour),
			},
		})

		b.ReportAllocs()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			tmp := val.Clone()
			_ = tmp
		}
	})

	// This benchmarks
	//   4. T is []*E; E implements [std.Cloner] with a normal receiver
	b.Run("PtrSliceWrapper_NormalRecv", func(b *testing.B) {
		val := typact.Some(MySimpleStructPtrList{
			{
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			{
				CreatedAt: time.Now().Add(time.Hour),
				UpdatedAt: time.Now().Add(time.Hour),
			},
		})

		b.ReportAllocs()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			tmp := val.Clone()
			_ = tmp
		}
	})

	// This benchmarks
	//   2. T is []E; E implements [std.Cloner] with a normal receiver
	b.Run("SliceWrapper_NormalRecv", func(b *testing.B) {
		val := typact.Some(MySimpleStructList{
			{
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			{
				CreatedAt: time.Now().Add(time.Hour),
				UpdatedAt: time.Now().Add(time.Hour),
			},
		})

		b.ReportAllocs()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			tmp := val.Clone()
			_ = tmp
		}
	})
}

func BenchmarkOption_MarshalText(b *testing.B) {
	b.Run("String", func(b *testing.B) {
		b.ReportAllocs()

		val := typact.Some("Hello")

		for range b.N {
			vv, err := val.MarshalText()
			if err != nil {
				b.Fatal(err)
			}

			_ = vv
		}
	})

	b.Run("Int", func(b *testing.B) {
		benchmarkMarshalText(int(46546), b)
	})

	b.Run("Int8", func(b *testing.B) {
		benchmarkMarshalText(int8(46), b)
	})

	b.Run("Int16", func(b *testing.B) {
		benchmarkMarshalText(int16(446), b)
	})

	b.Run("Int32", func(b *testing.B) {
		benchmarkMarshalText(int32(446), b)
	})

	b.Run("Int64", func(b *testing.B) {
		benchmarkMarshalText(int64(46546), b)
	})

	b.Run("Uint", func(b *testing.B) {
		benchmarkMarshalText(uint(46546), b)
	})

	b.Run("Uint8", func(b *testing.B) {
		benchmarkMarshalText(uint8(46), b)
	})

	b.Run("Uint16", func(b *testing.B) {
		benchmarkMarshalText(uint16(46546), b)
	})

	b.Run("Uint32", func(b *testing.B) {
		benchmarkMarshalText(uint32(46546), b)
	})

	b.Run("Uint64", func(b *testing.B) {
		benchmarkMarshalText(uint64(46546), b)
	})

	b.Run("Float32", func(b *testing.B) {
		benchmarkMarshalText(float32(46546.34), b)
	})

	b.Run("Float64", func(b *testing.B) {
		benchmarkMarshalText(float64(46546.345), b)
	})

	b.Run("Bool", func(b *testing.B) {
		benchmarkMarshalText(true, b)
	})
}

func benchmarkMarshalText[T any](vv T, b *testing.B) {
	b.ReportAllocs()

	val := typact.Some(vv)

	for range b.N {
		vv, err := val.MarshalText()
		if err != nil {
			b.Fatal(err)
		}

		_ = vv
	}
}

func BenchmarkOption_UnmarshalText(b *testing.B) {
	b.Run("String", func(b *testing.B) {
		benchmarkUnmarshalText[string]([]byte("hello"), b)
	})

	//	b.Run("Int", func(b *testing.B) {
	//		benchmarkMarshalText(int(46546), b)
	//	})
	//
	//	b.Run("Int8", func(b *testing.B) {
	//		benchmarkMarshalText(int8(46), b)
	//	})
	//
	//	b.Run("Int16", func(b *testing.B) {
	//		benchmarkMarshalText(int16(446), b)
	//	})
	//
	//	b.Run("Int32", func(b *testing.B) {
	//		benchmarkMarshalText(int32(446), b)
	//	})
	//
	//	b.Run("Int64", func(b *testing.B) {
	//		benchmarkMarshalText(int64(46546), b)
	//	})
	//
	//	b.Run("Uint", func(b *testing.B) {
	//		benchmarkMarshalText(uint(46546), b)
	//	})
	//
	//	b.Run("Uint8", func(b *testing.B) {
	//		benchmarkMarshalText(uint8(46), b)
	//	})
	//
	//	b.Run("Uint16", func(b *testing.B) {
	//		benchmarkMarshalText(uint16(46546), b)
	//	})
	//
	//	b.Run("Uint32", func(b *testing.B) {
	//		benchmarkMarshalText(uint32(46546), b)
	//	})
	//
	//	b.Run("Uint64", func(b *testing.B) {
	//		benchmarkMarshalText(uint64(46546), b)
	//	})
	//
	//	b.Run("Float32", func(b *testing.B) {
	//		benchmarkMarshalText(float32(46546.34), b)
	//	})
	//
	//	b.Run("Float64", func(b *testing.B) {
	//		benchmarkMarshalText(float64(46546.345), b)
	//	})
	//
	//	b.Run("Bool", func(b *testing.B) {
	//		benchmarkMarshalText(true, b)
	//	})
}

func benchmarkUnmarshalText[T any](value []byte, b *testing.B) {
	b.ReportAllocs()

	val := typact.None[T]()

	for range b.N {
		err := val.UnmarshalText(value)
		if err != nil {
			b.Fatal(err)
		}

		_ = val
	}
}
