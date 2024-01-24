package option_test

import (
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
		It("should apply the function and return the new option when the option is present", func() {
			option := typact.Some(5)
			fn := func(val int) typact.Option[int] {
				return typact.Some(val * 2)
			}

			newOption := option.AndThen(fn)
			Expect(newOption.IsSome()).To(BeTrue())
			Expect(newOption.UnsafeUnwrap()).To(Equal(10))
		})

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
		var isEven = func(n int) bool { return n%2 == 0 }

		It("should return None when the original option is None", func() {
			original := typact.None[int]()
			result := original.Filter(isEven)

			Expect(result.IsSome()).To(BeFalse())
		})

		It("should return None when the original option is Some but does not satisfy the filter function", func() {
			original := typact.Some[int](3)
			result := original.Filter(isEven)

			Expect(result.IsSome()).To(BeFalse())
		})

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
})
