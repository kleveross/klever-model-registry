package integration

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	httpexpect "github.com/gavv/httpexpect/v2"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	modeljobsv1alpha1 "github.com/kleveross/klever-model-registry/pkg/apis/modeljob/v1alpha1"
	"github.com/kleveross/klever-model-registry/pkg/registry/models"
)

var _ = Describe("Model Registry", func() {
	const timeout = time.Second * 5
	const interval = time.Second * 1

	e := httpexpect.New(GinkgoT(), ModelRegistryHost)
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
			// Upload model will create ModelJob automatically, now get ModelJobList.
			Eventually(func() error {
				length := e.GET("/api/v1alpha1/namespaces/{namespace}/modeljobs/",
					"default").Expect().JSON().Object().Value("items").Array().Length()
				if int(length.Raw()) != 1 {
					return fmt.Errorf("Not found any modeljob")
				}
				return nil
			}, timeout, interval).Should(Succeed())

		})
		It("Should list the model successfully", func() {
			artifact := e.GET("/api/v2.0/projects/{project_name}/repositories/{repository_name}/artifacts/{version}",
				project, model, version).Expect().Status(http.StatusOK).JSON().Object()
			// Validate that the artifact is a SavedModel.
			// It is blocked since goharbor/harbor-helm does not support 2.1 now.
			// artifact.Value("extra_attrs").Object().Value("format").Equal("SavedModel")
			artifact.Value("type").Equal("MODEL")
		})

		Context("ModelJobs", func() {
			It("Should get the ModelJobs successfully", func() {
				e.GET("/api/v1alpha1/namespaces/{namespace}/modeljobs/",
					"default").Expect().Status(http.StatusOK)
			})

			It("Should create the ModelJobs successfully", func() {
				job := modeljobsv1alpha1.ModelJob{
					Spec: modeljobsv1alpha1.ModelJobSpec{
						// Use ModelRegistryHost/<project>/<model>:<version>.
						Model: fmt.Sprintf("%s/%s/%s:%s", ModelRegistryHost[7:], project, model, version),
						ModelJobSource: modeljobsv1alpha1.ModelJobSource{
							Extraction: &modeljobsv1alpha1.ExtractionSource{
								Format: modeljobsv1alpha1.FormatSavedModel,
							},
						},
					},
				}
				e.POST("/api/v1alpha1/namespaces/{namespace}/modeljobs", "default").WithJSON(job).Expect().Status(http.StatusCreated)
			})
		})
	})
})
