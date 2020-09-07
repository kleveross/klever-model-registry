package integration

import (
	"encoding/json"

	httpexpect "github.com/gavv/httpexpect/v2"
	"github.com/kleveross/klever-model-registry/pkg/registry/models"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
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
			project := "library"
			model := "tensorflow"
			version := "test"
			modelContent := models.Model{
				ProjectName: project,
				ModelName:   model,
				VersionName: version,
				Format:      "SavedModel",
			}
			bytes, err := json.Marshal(modelContent)
			Expect(err).Should(BeNil())

			e.POST("/api/v1alpha1/projects/{projectName}/models/{modelName}/versions/{versionName}/upload",
				project, model, version).
				WithHeaders(map[string]string{
					"X-Tenant": "test",
					"X-User":   "test",
				}).
				WithMultipart().
				WithFile("file", "./models/model.zip").WithFormField("model", string(bytes)).
				Expect().Status(201)
		})
	})
})
