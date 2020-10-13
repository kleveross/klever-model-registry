package integration

import (
	"fmt"
	"net/http"
	"time"

	httpexpect "github.com/gavv/httpexpect/v2"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	// . "github.com/onsi/gomega"
	seldonv1 "github.com/seldonio/seldon-core/operator/apis/machinelearning.seldon.io/v1"
)

var _ = Describe("Model Registry", func() {
	const timeout = time.Second * 5
	const interval = time.Second * 1
	project := "library"
	model := "tensorflow"
	version := "test"

	genTestResource := func() corev1.ResourceList {
		resourcesList := make(corev1.ResourceList)
		cpuQuantity, _ := resource.ParseQuantity("1")
		resourcesList[corev1.ResourceCPU] = cpuQuantity
		memQuantity, _ := resource.ParseQuantity("1Gi")
		resourcesList[corev1.ResourceMemory] = memQuantity

		return resourcesList
	}

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
							ComponentSpecs: []*seldonv1.SeldonPodSpec{
								{
									Metadata: metav1.ObjectMeta{
										Name: "test",
									},
									Spec: corev1.PodSpec{
										Containers: []corev1.Container{
											{
												Resources: corev1.ResourceRequirements{
													Limits:   genTestResource(),
													Requests: genTestResource(),
												},
											},
										},
									},
								},
							},
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
			Eventually(func() int {
				// TODO: better attersions
				return int(e.GET("/api/v1alpha1/namespaces/{namespace}/servings/",
					"default").Expect().Status(http.StatusOK).
					JSON().Object().Value("items").Array().Length().Raw())
			}, timeout, interval).Should(Equal(1))
		})
		It("Should update the Serving with HPA successfully", func() {
			reqBody := []byte(`{
				"spec": {
					"predictors": [
						{
							"componentSpecs": [
								{
									"hpaSpec": {
										"minReplicas": 2,
										"maxReplicas": 4
									}
								}
							]
						}
					]
				}
			}`)
			Eventually(
				func() int {
					return e.PUT("/api/v1alpha1/namespaces/{namespace}/servings/{servingID}", "default", "test").
						WithHeader("Content-Type", "application/json").WithBytes(reqBody).Expect().Raw().StatusCode
				}, timeout, interval).Should(Equal(http.StatusOK))
		})
		It("Should get the updated Serving with correct HPA", func() {
			Eventually(func() bool {
				_, ok := e.GET("/api/v1alpha1/namespaces/{namespace}/servings/{servingID}", "default", "test").
					Expect().Status(http.StatusOK).JSON().Object().Value("spec").Object().Value("predictors").
					Array().Element(0).Object().Value("componentSpecs").Array().Element(0).Object().Raw()["hpaSpec"]
				return ok
			}, timeout, interval).Should(BeTrue())
		})
	})
})
