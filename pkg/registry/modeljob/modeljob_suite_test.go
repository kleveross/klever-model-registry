package modeljob_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestModeljob(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Modeljob Suite")
}
