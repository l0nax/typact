package option_test

import (
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.l0nax.org/typact"
)

var _ = Describe("Clone", func() {
	It("should clone a None value", func() {
		vv := typact.None[string]()
		cpy := vv.Clone()

		Expect(cpy.IsNone()).To(BeTrue())
	})

	Describe("primitive types", func() {
		It("should clone string", func() {
			vv := typact.Some("foo bar")
			cpy := vv.Clone().Unwrap()

			// change vv
			xx := vv.UnwrapAsRef()
			*xx = "changed"

			Expect(cpy).To(BeEquivalentTo("foo bar"))
		})

		It("should clone integer", func() {
			vv := typact.Some(42)
			cpy := vv.Clone().Unwrap()

			// change vv
			xx := vv.UnwrapAsRef()
			*xx = 654654

			Expect(cpy).To(BeEquivalentTo(42))
		})

		It("should clone bool", func() {
			vv := typact.Some(true)
			cpy := vv.Clone().Unwrap()

			// change vv
			xx := vv.UnwrapAsRef()
			*xx = false

			Expect(cpy).To(BeEquivalentTo(true))
		})

		It("should clone float ptr", func() {
			vv := typact.Some(toPtr(3.14))
			cpy := vv.Clone().Unwrap()

			// change vv
			xx := vv.UnwrapAsRef()
			*xx = toPtr(7.0)

			Expect(cpy).ToNot(BeNil())
			Expect(*cpy).To(BeEquivalentTo(3.14))
		})

		It("should clone a nil ptr", func() {
			vv := typact.Some((*string)(nil))
			cpy := vv.Clone().Unwrap()

			// change vv
			xx := vv.UnwrapAsRef()
			*xx = toPtr("foo bar")

			Expect(cpy).To(BeNil())
		})

		It("should clone a string slice", func() {
			vv := typact.Some([]string{"foo", "bar"})
			cpy := vv.Clone().Unwrap()

			// change vv
			vv.Unwrap()[0] = "hello"

			Expect(cpy).To(BeEquivalentTo([]string{"foo", "bar"}))
		})
	})

	Describe("wrapped primitive types", func() {
		It("should use the custom Clone method", func() {
			vv := typact.Some(myStrWithClone("foo bar"))
			cpy := vv.Clone().Unwrap()

			Expect(cpy).To(BeEquivalentTo("cpy: foo bar"))
		})

		It("should NOT use the custom Clone method on pointer", func() {
			data := myStrWithClone("foo bar")
			vv := typact.Some(&data)
			cpy := vv.Clone().Unwrap()

			// NOTE: It is not directly possible to call the custom Clone[T] method
			// and since [myStrWithClone] implements the [std.Clone] interface only
			// for the non-pointer variant, it should fallback to our implementation.
			//
			// If someone wants to support this specific usecase, please create an issue.
			Expect(*cpy).To(BeEquivalentTo("foo bar"))
		})

		It("should clone an alias to primitive type", func() {
			vv := typact.Some(myStrAlias("foo bar"))
			cpy := vv.Clone().Unwrap()

			// change vv
			xx := vv.UnwrapAsRef()
			*xx = "hello world"

			Expect(cpy).To(BeEquivalentTo(myStrAlias("foo bar")))
		})

		It("should clone a wrapped primitive type without Cloner", func() {
			vv := typact.Some(mySimpleStr("foo bar"))
			cpy := vv.Clone().Unwrap()

			// change vv
			xx := vv.UnwrapAsRef()
			*xx = "hello world"

			Expect(cpy).To(BeEquivalentTo(mySimpleStr("foo bar")))
		})

		It("should clone a string slice alias", func() {
			vv := typact.Some(myStrSliceAlias([]string{"foo", "bar"}))
			cpy := vv.Clone().Unwrap()

			// change vv
			vv.Unwrap()[0] = "hello"

			Expect(cpy).To(BeEquivalentTo(myStrSliceAlias([]string{"foo", "bar"})))
		})

		It("should call Clone on a wrapped string slice", func() {
			vv := typact.Some(myStrSlice([]string{"foo", "bar"}))
			cpy := vv.Clone().Unwrap()

			Expect(cpy).To(BeEquivalentTo(myStrSlice([]string{"cpy: foo", "cpy: bar"})))
		})
	})

	Describe("clone custom structs", func() {
		It("should clone the custom method", func() {
			tt := time.Now()

			vv := typact.Some(myStruct{
				CreatedAt: tt,
				Data:      "foo bar",
			})
			cpy := vv.Clone().Unwrap()

			// change vv
			vv.UnwrapAsRef().CreatedAt = tt.AddDate(0, 0, 1)

			Expect(cpy).To(BeEquivalentTo(myStruct{
				CreatedAt: tt,
				Data:      "foo bar",
			}))
		})

		It("should fallback to the pointer receiver", func() {
			vv := typact.Some(MyData{
				Data: "foo bar",
			})
			cpy := vv.Clone().Unwrap()

			// change vv
			vv.UnwrapAsRef().Data = "hello world"

			Expect(cpy).To(BeEquivalentTo(MyData{
				Data: "cpy: foo bar",
			}))
		})

		/*
			// XXX: This will/ must panic!
					It("should call the custom Clone method on ptr ref", func() {
						tt := time.Now()

						vv := typact.Some(&myStruct{
							CreatedAt: tt,
							Data:      "foo bar",
						})
						cpy := vv.Clone().Unwrap()

						// change vv
						vv.Unwrap().CreatedAt = tt.AddDate(0, 0, 1)

						Expect(*cpy).To(BeEquivalentTo(myStruct{
							CreatedAt: tt,
							Data:      "foo bar",
						}))
					})
		*/
	})
})

type myStrSliceAlias = []string

type myStrSlice []string

func (m myStrSlice) Clone() myStrSlice {
	cpy := make(myStrSlice, len(m))
	for i := range m {
		cpy[i] = "cpy: " + m[i]
	}

	return cpy
}

type mySimpleStr string

type myStrAlias = string

type myStrWithClone string

func (m myStrWithClone) Clone() myStrWithClone {
	return "cpy: " + m
}

func toPtr[T any](val T) *T {
	return &val
}

type myStruct struct {
	CreatedAt time.Time
	Data      string
}

func (m myStruct) Clone() myStruct {
	return myStruct{
		CreatedAt: m.CreatedAt,
		Data:      m.Data,
	}
}

type MyData struct {
	Data string
}

func (m *MyData) Clone() *MyData {
	return &MyData{
		Data: "cpy: " + m.Data,
	}
}
