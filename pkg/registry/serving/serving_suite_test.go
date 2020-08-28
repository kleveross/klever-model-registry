package serving_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestServing(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Serving Suite")
}
