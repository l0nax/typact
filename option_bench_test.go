package typact_test

import (
	"testing"

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
