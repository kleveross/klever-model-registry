package integration

import (
	"fmt"
	"net/http"
	"time"

	httpexpect "github.com/gavv/httpexpect/v2"
	. "github.com/onsi/ginkgo"

	// . "github.com/onsi/gomega"
	seldonv1 "github.com/seldonio/seldon-core/operator/apis/machinelearning.seldon.io/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe("Model Registry", func() {
	const timeout = time.Second * 5
	const interval = time.Second * 1
	project := "library"
	model := "tensorflow"
	version := "test"

	e := httpexpect.New(GinkgoT(), ModelRegistryHost)
	Context("Servings", func() {
		name := "test"
		It("Should create the Serving successfully", func() {
			seldonCoreDeploy := &seldonv1.SeldonDeployment{
				ObjectMeta: metav1.ObjectMeta{
					Name: name,
				},
				Spec: seldonv1.SeldonDeploymentSpec{
					Predictors: []seldonv1.PredictorSpec{
						{
							Graph: seldonv1.PredictiveUnit{
								Name:     "test",
								ModelURI: fmt.Sprintf("%s/%s/%s:%s", ModelRegistryHost[7:], project, model, version),
								Parameters: []seldonv1.Parameter{
									{
										Name:  "cpu",
										Value: "1",
									},
									{
										Name:  "mem",
										Value: "100Mi",
									},
								},
							},
						},
					},
				},
			}
			e.POST("/api/v1alpha1/namespaces/{namespace}/servings", "default").WithJSON(seldonCoreDeploy).Expect().Status(http.StatusCreated)
		})
		It("Should get the Servings successfully", func() {
			e.GET("/api/v1alpha1/namespaces/{namespace}/servings/",
				"default").Expect().Status(http.StatusOK).
				JSON().Object().Value("items").Array().Length().Equal(1)
		})
	})
})
