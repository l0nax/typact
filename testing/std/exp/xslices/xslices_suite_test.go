package xslices_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestXslices(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "xslices Suite")
}
