package integration

import (
	httpexpect "github.com/gavv/httpexpect/v2"
	. "github.com/onsi/ginkgo"
	// . "github.com/onsi/gomega"
)

var _ = Describe("Model Registry", func() {
	e := httpexpect.New(GinkgoT(), ModelRegistryHost)
	Context("ModelJobs", func() {
		It("Should get the ModelJobs successfully", func() {
			e.GET("/api/v1alpha1/namespaces/{namespace}/modeljobs/",
				"default").Expect().Status(200)
		})
	})
	Context("Servings", func() {
		It("Should get the Servings successfully", func() {
			e.GET("/api/v1alpha1/namespaces/{namespace}/servings/",
				"default").Expect().Status(200)
		})
	})
	Context("Models", func() {
		It("Should push the model successfully", func() {
			e.POST("/api/v1alpha1/projects/{projectName}/models/{modelName}/versions/{versionName}/upload",
				"library", "tensorflow", "test").WithMultipart().WithFile("model", "./models/model.zip").Expect().Status(200)
		})
	})
})
