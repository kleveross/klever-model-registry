package integration

import (
	"encoding/json"
	"net/http"

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
				"default").Expect().Status(http.StatusOK)
		})
	})
	Context("Servings", func() {
		It("Should get the Servings successfully", func() {
			e.GET("/api/v1alpha1/namespaces/{namespace}/servings/",
				"default").Expect().Status(http.StatusOK)
		})
	})
	Context("Models", func() {
		project := "library"
		model := "tensorflow"
		version := "test"

		It("Should push the model successfully", func() {
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
				Expect().Status(http.StatusCreated)
		})
		It("Should list the model successfully", func() {
			artifact := e.GET("/api/v2.0/projects/{project_name}/repositories/{repository_name}/artifacts/{version}",
				project, model, version).Expect().Status(http.StatusOK).JSON().Object()
			// Validate that the artifact is a SavedModel.
			// It is blocked since goharbor/harbor-helm does not support 2.1 now.
			// artifact.Value("extra_attrs").Object().Value("format").Equal("SavedModel")
			artifact.Value("type").Equal("MODEL")
		})
	})
})
