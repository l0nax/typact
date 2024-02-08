package pred_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"go.l0nax.org/typact/std/exp/pred"
)

var _ = Describe("Zero", func() {
	Describe("IsZero", func() {
		It("should return true for zero values", func() {
			var (
				str    string
				num    int
				bl     bool
				strPtr *string
			)

			Expect(pred.IsZero(str)).To(BeTrue())
			Expect(pred.IsZero(num)).To(BeTrue())
			Expect(pred.IsZero(bl)).To(BeTrue())
			Expect(pred.IsZero(strPtr)).To(BeTrue())
		})

		It("should return false for non-zero values", func() {
			var (
				str    = "Hello, World!"
				num    = 12
				bl     = true
				strPtr = &str
			)

			Expect(pred.IsZero(str)).To(BeFalse())
			Expect(pred.IsZero(num)).To(BeFalse())
			Expect(pred.IsZero(bl)).To(BeFalse())
			Expect(pred.IsZero(strPtr)).To(BeFalse())
		})
	})
})
