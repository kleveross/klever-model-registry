package integration

import (
	httpexpect "github.com/gavv/httpexpect/v2"
	. "github.com/onsi/ginkgo"
	// . "github.com/onsi/gomega"
)

var _ = Describe("Model Registry", func() {
	e := httpexpect.New(GinkgoT(), ModelRegistryHost)
	Context("ModelJobs", func() {
		It("Should get the ModelJob successfully", func() {
			e.GET("/api/v1alpha1/namespaces/{namespace}/modeljobs/",
				"default").Expect().Status(200)
		})
	})
})
