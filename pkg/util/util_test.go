package util_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/kleveross/klever-model-registry/pkg/util"
)

var _ = Describe("Util", func() {
	It("Should util successfully", func() {
		prefix := "name-prefix"
		name := util.RandomNameWithPrefix(prefix)
		Expect(len(name)).To(Equal(len(prefix) + 24))
	})
})
