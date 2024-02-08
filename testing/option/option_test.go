package option_test

import (
	"encoding/json"
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.l0nax.org/typact"
)

var _ = Describe("Option", func() {
	Describe("UnwrapAsRef", func() {
		It("with string", func() {
			val := typact.Some("this is my option")
			(*val.UnwrapAsRef()) = "new value"

			Expect(val.Unwrap()).To(BeEquivalentTo("new value"))
		})

		Context("with string slice", func() {
			It("should alter in place", func() {
				val := typact.Some([]string{"hello", "world"})
				(*val.UnwrapAsRef())[1] = "my new"

				Expect(val.Unwrap()).To(BeEquivalentTo([]string{
					"hello",
					"my new",
				}))
			})

			It("should alter the refence", func() {
				val := typact.Some([]string{"hello", "world"})
				ref := val.UnwrapAsRef()
				(*ref) = append((*ref), "from home")

				Expect(val.Unwrap()).To(BeEquivalentTo([]string{
					"hello",
					"world",
					"from home",
				}))
			})
		})

		Context("with map", func() {
			It("should alter in place", func() {
				mm := map[string]string{
					"hello": "world",
					"foo":   "bar",
				}
				val := typact.Some(mm)
				(*val.UnwrapAsRef())["foo"] = "my new"

				Expect(val.Unwrap()).To(BeEquivalentTo(map[string]string{
					"hello": "world",
					"foo":   "my new",
				}))
			})

			It("should alter the refence", func() {
				mm := map[string]string{
					"hello": "world",
					"foo":   "bar",
				}

				val := typact.Some(mm)
				ref := val.UnwrapAsRef()
				delete(*ref, "foo")

				Expect(val.Unwrap()).To(BeEquivalentTo(map[string]string{
					"hello": "world",
				}))
			})
		})
	})

	Describe("IsSomeAnd", func() {
		It("should return true if fn returns true and is Some", func() {
			val := typact.Some("foo bar")
			ret := val.IsSomeAnd(func(s string) bool {
				return s == "foo bar"
			})

			Expect(ret).To(BeTrue())
		})

		It("should return false and not call fn if is None", func() {
			val := typact.None[string]()
			ret := val.IsSomeAnd(func(s string) bool {
				Fail("fn should not be called")

				return true
			})

			Expect(ret).To(BeFalse())
		})

		It("should return false if fn returns false and is Some", func() {
			val := typact.Some("foo bar")
			ret := val.IsSomeAnd(func(s string) bool {
				return false
			})

			Expect(ret).To(BeFalse())
		})
	})

	Describe("Or", func() {
		It("when first is Some and second is Some", func() {
			val := typact.Some("this is my option")
			other := typact.Some("other option")

			newVal := val.Or(other)

			Expect(newVal.Unwrap()).To(BeEquivalentTo("this is my option"))
		})

		It("when first is Some and second is None", func() {
			val := typact.Some("this is my option")
			other := typact.None[string]()

			newVal := val.Or(other)

			Expect(newVal.Unwrap()).To(BeEquivalentTo("this is my option"))
		})

		It("when first is None and second is Some", func() {
			val := typact.None[string]()
			other := typact.Some("other option")

			newVal := val.Or(other)

			Expect(newVal.Unwrap()).To(BeEquivalentTo("other option"))
		})

		It("when first is None and second is None", func() {
			val := typact.None[string]()
			other := typact.None[string]()

			newVal := val.Or(other)

			Expect(newVal.IsNone()).To(BeTrue())
		})
	})

	Describe("Unwrap", func() {
		It("should return the value", func() {
			x := typact.Some(5)
			Expect(x.Unwrap()).To(Equal(5))
		})

		It("should panic", func() {
			x := typact.None[int]()
			Expect(
				func() {
					x.Unwrap()
				}).
				To(Panic())
		})
	})

	Describe("UnwrapOrZero", func() {
		It("should return the value", func() {
			x := typact.Some(5)
			Expect(x.UnwrapOrZero()).To(Equal(5))
		})

		It("should return the zero value", func() {
			x := typact.None[int]()
			Expect(x.UnwrapOrZero()).To(Equal(0))
		})
	})

	Describe("UnwrapOrElse", func() {
		getTwo := func() int {
			return 2
		}

		It("should return the value", func() {
			x := typact.Some(5)
			Expect(x.UnwrapOrElse(getTwo)).To(Equal(5))
		})

		It("should return the zero value", func() {
			x := typact.None[int]()
			Expect(x.UnwrapOrElse(getTwo)).To(Equal(2))
		})
	})

	Describe("UnwrapAsRef", func() {
		It("should return a pointer to the value", func() {
			x := typact.Some(5)

			ref := x.UnwrapAsRef()
			*ref = 10
			Expect(x.Unwrap()).To(Equal(10))
		})

		It("should panic", func() {
			x := typact.None[int]()
			Expect(
				func() {
					x.UnwrapAsRef()
				}).
				To(Panic())
		})
	})

	Describe("Insert", func() {
		It("should change the value once we change it via the pointer", func() {
			opt := typact.None[int]()
			val := opt.Insert(5)

			Expect(*val).To(BeEquivalentTo(5))
			Expect(opt.Unwrap()).To(BeEquivalentTo(5))

			*val = 3
			Expect(opt.Unwrap()).To(BeEquivalentTo(3), "option value must have changed")
		})
	})

	Describe("GetOrInsert", func() {
		It("should change the value if its None", func() {
			opt := typact.None[int]()
			val := opt.GetOrInsert(5)

			Expect(*val).To(BeEquivalentTo(5))
			Expect(opt.Unwrap()).To(BeEquivalentTo(5))

			*val = 3
			Expect(opt.Unwrap()).To(BeEquivalentTo(3), "option value must have changed")
		})

		It("should NOT change the value if its None", func() {
			opt := typact.Some[int](10)
			val := opt.GetOrInsert(5)

			Expect(*val).To(BeEquivalentTo(10))
			Expect(opt.Unwrap()).To(BeEquivalentTo(10))

			*val = 3
			Expect(opt.Unwrap()).To(BeEquivalentTo(3), "option value must have changed")
		})
	})

	Describe("Replace", func() {
		It("should replace the value if Some", func() {
			opt := typact.Some(5)
			old := opt.Replace(10)

			Expect(opt.Unwrap()).To(BeEquivalentTo(10))
			Expect(old.Unwrap()).To(BeEquivalentTo(5))
		})

		It("should replace the value if None", func() {
			opt := typact.None[int]()
			old := opt.Replace(10)

			Expect(opt.Unwrap()).To(BeEquivalentTo(10))
			Expect(old.IsNone()).To(BeTrue())
		})
	})

	Describe("AndThen", func() {
		It(
			"should apply the function and return the new option when the option is present",
			func() {
				option := typact.Some(5)
				fn := func(val int) typact.Option[int] {
					return typact.Some(val * 2)
				}

				newOption := option.AndThen(fn)
				Expect(newOption.IsSome()).To(BeTrue())
				Expect(newOption.UnsafeUnwrap()).To(Equal(10))
			},
		)

		It("should return None when the option is not present", func() {
			option := typact.None[int]()
			fn := func(val int) typact.Option[int] {
				return typact.Some(val * 2)
			}

			newOption := option.AndThen(fn)
			Expect(newOption.IsSome()).To(BeFalse())
		})
	})

	Describe("And", func() {
		It("should return None when the original option is None", func() {
			original := typact.None[string]()
			opt := typact.Some[string]("value")

			result := original.And(opt)

			Expect(result.IsSome()).To(BeFalse())
		})

		It("should return the provided option when the original option is Some", func() {
			original := typact.Some[string]("original")
			opt := typact.Some[string]("value")

			result := original.And(opt)

			Expect(result.IsSome()).To(BeTrue())
			Expect(result.Unwrap()).To(Equal("value"))
		})
	})

	Describe("Filter", func() {
		isEven := func(n int) bool { return n%2 == 0 }

		It("should return None when the original option is None", func() {
			original := typact.None[int]()
			result := original.Filter(isEven)

			Expect(result.IsSome()).To(BeFalse())
		})

		It(
			"should return None when the original option is Some but does not satisfy the filter function",
			func() {
				original := typact.Some[int](3)
				result := original.Filter(isEven)

				Expect(result.IsSome()).To(BeFalse())
			},
		)

		It("should return the original option when it satisfies the filter function", func() {
			original := typact.Some[int](4)
			result := original.Filter(isEven)

			Expect(result.IsSome()).To(BeTrue())
			Expect(result.UnsafeUnwrap()).To(Equal(4))
		})
	})

	Describe("Map", func() {
		It("should return transformed value if present", func() {
			opt := typact.Some("foo")
			result := opt.Map(func(str string) string { return "bar" })

			Expect(result.UnsafeUnwrap()).To(Equal("bar"))
		})

		It("should return None if no value is present", func() {
			opt := typact.None[string]()
			result := opt.Map(func(str string) string { return "bar" })

			Expect(result.IsSome()).To(BeFalse())
		})
	})

	Describe("MapOr", func() {
		It("should return transformed value if present", func() {
			opt := typact.Some("foo")
			result := opt.MapOr(func(str string) string { return "bar" }, "alt")

			Expect(result).To(Equal("bar"))
		})

		It("should return provided value if no value is present", func() {
			opt := typact.None[string]()
			result := opt.MapOr(func(str string) string { return "bar" }, "alt")

			Expect(result).To(Equal("alt"))
		})
	})

	Describe("MapOrElse", func() {
		It("should return transformed value if present", func() {
			opt := typact.Some("foo")
			result := opt.MapOrElse(
				func(str string) string { return "bar" },
				func() string { return "alternative" },
			)

			Expect(result).To(Equal("bar"))
		})

		It("should return result of valueFn if no value is present", func() {
			opt := typact.None[string]()
			result := opt.MapOrElse(
				func(str string) string { return "bar" },
				func() string { return "alternative" },
			)

			Expect(result).To(Equal("alternative"))
		})
	})

	Describe("OrElse with AndThen", func() {
		type ScalarWrapper struct {
			Number uint32
		}

		type OptionalHolder struct {
			Wrapper typact.Option[*ScalarWrapper]
		}

		It("should return the changed value from the base", func() {
			ref := &OptionalHolder{
				Wrapper: typact.Some(&ScalarWrapper{
					Number: 21324534,
				}),
			}

			ref.Wrapper = ref.Wrapper.
				OrElse(func() typact.Option[*ScalarWrapper] {
					return typact.Some(&ScalarWrapper{
						Number: 456456,
					})
				}).
				AndThen(func(sw *ScalarWrapper) typact.Option[*ScalarWrapper] {
					Expect(sw.Number).To(BeEquivalentTo(21324534))
					sw.Number = 89722
					return typact.Some(sw)
				})

			Expect(ref.Wrapper.Unwrap().Number).To(BeEquivalentTo(89722))
		})

		It("should return the changed value from the alternative", func() {
			ref := &OptionalHolder{
				Wrapper: typact.None[*ScalarWrapper](),
			}

			ref.Wrapper = ref.Wrapper.
				OrElse(func() typact.Option[*ScalarWrapper] {
					return typact.Some(&ScalarWrapper{
						Number: 456456,
					})
				}).
				AndThen(func(sw *ScalarWrapper) typact.Option[*ScalarWrapper] {
					Expect(sw.Number).To(BeEquivalentTo(456456))
					sw.Number = 89722
					return typact.Some(sw)
				})

			Expect(ref.Wrapper.Unwrap().Number).To(BeEquivalentTo(89722))
		})
	})

	Describe("JSON", func() {
		type MyData struct {
			Str typact.Option[string] `json:"str"`
			Num typact.Option[int]    `json:"num"`
		}

		Context("Unmarshal", func() {
			It("should be able to handle null", func() {
				const raw = `{"str":null,"num":125}`

				var data MyData

				err := json.Unmarshal([]byte(raw), &data)
				Expect(err).ToNot(HaveOccurred())

				Expect(data.Str.IsNone()).To(BeTrue())
				Expect(data.Num.Unwrap()).To(BeEquivalentTo(125))
			})

			It("should handle missing fields", func() {
				const raw = `{"num":125}`

				var data MyData

				err := json.Unmarshal([]byte(raw), &data)
				Expect(err).ToNot(HaveOccurred())

				Expect(data.Str.IsNone()).To(BeTrue())
				Expect(data.Num.Unwrap()).To(BeEquivalentTo(125))
			})

			It("should handle empty fields", func() {
				const raw = `{"str":"","num":125}`

				var data MyData

				err := json.Unmarshal([]byte(raw), &data)
				Expect(err).ToNot(HaveOccurred())

				Expect(data.Str.Unwrap()).To(BeEquivalentTo(""))
				Expect(data.Num.Unwrap()).To(BeEquivalentTo(125))
			})

			It("should handle invalid JSON", func() {
				const raw = `{"num":null,"str":}`

				var data MyData

				err := json.Unmarshal([]byte(raw), &data)
				Expect(err).To(HaveOccurred())
			})
		})

		Context("Marshal", func() {
			It("should correctly encode None", func() {
				data := MyData{
					Str: typact.None[string](),
					Num: typact.Some(125),
				}

				b, err := json.Marshal(data)
				Expect(err).ToNot(HaveOccurred())
				Expect(string(b)).To(BeEquivalentTo(`{"str":null,"num":125}`))
			})
		})
	})

	Describe("Database Scan/Value", func() {
		DescribeTable("Scanning on string should work",
			func(inputValue any, wantErr bool, isSome bool) {
				var vv typact.Option[string]

				err := vv.Scan(inputValue)
				if wantErr {
					Expect(err).To(HaveOccurred())
				} else {
					Expect(err).ToNot(HaveOccurred())

					if isSome {
						Expect(vv.IsSome()).To(BeTrue())
						Expect(vv.Unwrap()).To(BeEquivalentTo(inputValue))
					} else {
						Expect(vv.IsNone()).To(BeTrue())
					}
				}
			},
			Entry("normal string input", "foo bar", false, true),
			Entry("null input", nil, false, false),
			Entry("byte slice as input", []byte("hello world"), true, false),
			Entry("complete other type", 555, true, false),
		)

		Context("Scan on unsupported type", func() {
			It("should error on non-null input", func() {
				var vv typact.Option[func()]

				err := vv.Scan("foo")
				Expect(err).To(HaveOccurred())
				Expect(vv.IsNone()).To(BeTrue())
			})

			// NOTE: This is the same behavior as the [database/sql.Null] type!
			It("should not error on null input", func() {
				var vv typact.Option[func()]

				err := vv.Scan(nil)
				Expect(err).ToNot(HaveOccurred())
				Expect(vv.IsNone()).To(BeTrue())
			})

			It("should handle errors of custom scanner", func() {
				var vv typact.Option[errScanner]

				err := vv.Scan("foo")
				Expect(err).To(HaveOccurred())
				Expect(err).To(MatchError(ContainSubstring("never scan")))
				Expect(vv.IsNone()).To(BeTrue())
			})

			It("should call the custom scanner implementation", func() {
				var vv typact.Option[customScanner]

				err := vv.Scan("foo")
				Expect(err).ToNot(HaveOccurred())
				Expect(vv.IsSome()).To(BeTrue())
				Expect(vv.Unwrap()).To(Equal(customScanner{Data: "custom: foo"}))
			})

			It("should override existing values", func() {
				vv := typact.Some(customScanner{Data: "start"})

				err := vv.Scan("foo")
				Expect(err).ToNot(HaveOccurred())
				Expect(vv.IsSome()).To(BeTrue())
				Expect(vv.Unwrap()).To(Equal(customScanner{Data: "custom: foo"}))
			})
		})
	})
})

type errScanner struct {
	Data string
}

func (s *errScanner) Scan(val any) error {
	return fmt.Errorf("never scan")
}

type customScanner struct {
	Data string
}

func (c *customScanner) Scan(val any) error {
	switch v := val.(type) {
	case string:
		c.Data = "custom: " + v
		return nil
	}

	return fmt.Errorf("unsupported type %T", val)
}
