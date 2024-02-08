package pred_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestPred(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Pred Suite")
}
