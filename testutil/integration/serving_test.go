package integration

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/caicloud/nirvana/log"
	httpexpect "github.com/gavv/httpexpect/v2"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"github.com/kleveross/klever-model-registry/pkg/util"
	seldonv1 "github.com/seldonio/seldon-core/operator/apis/machinelearning.seldon.io/v1"
)

var _ = Describe("Model Registry", func() {
	const timeout = time.Second * 5
	const interval = time.Second * 1
	var replica int32 = 1
	project := "library"
	model := "pmml"
	version := "v1"
	ns := "default"
	servingName := util.RandomNameWithPrefix("serving-test")

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
		It("Should create the Serving successfully", func() {
			minReplicas := int32(1)
			hpaSpec := seldonv1.SeldonHpaSpec{
				MinReplicas: &minReplicas,
				MaxReplicas: 2,
			}
			seldonCoreDeploy := &seldonv1.SeldonDeployment{
				ObjectMeta: metav1.ObjectMeta{
					Name: servingName,
					Namespace: ns,
				},
				Spec: seldonv1.SeldonDeploymentSpec{
					Predictors: []seldonv1.PredictorSpec{
						{
							ComponentSpecs: []*seldonv1.SeldonPodSpec{
								{
									Metadata: metav1.ObjectMeta{
										Name: servingName,
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
									HpaSpec: &hpaSpec,
									Replicas: &replica,
								},
							},
							Graph: seldonv1.PredictiveUnit{
								Name:     servingName,
								ModelURI: fmt.Sprintf("harbor-harbor-core.harbor-system/%s/%s:%s", project, model, version),
								ServiceAccountName: "default",
								Parameters: []seldonv1.Parameter{
									{
										Name:  "format",
										Value: "PMML",
									},
								},
							},
						},
					},
				},
			}
			e.POST("/api/v1alpha1/namespaces/{namespace}/servings", ns).WithJSON(seldonCoreDeploy).
				Expect().Status(http.StatusCreated)
		})

		It("The list should contains the Serving", func() {
			Eventually(func() []byte {
				str, err := json.Marshal(e.GET("/api/v1alpha1/namespaces/{namespace}/servings/", ns).
					Expect().Status(http.StatusOK).
					JSON().Object().Raw())
				if err != nil{
					log.Error(err)
				}
				return str
			}, timeout, interval).Should(ContainSubstring(servingName))
		})

		It("Should get the Serving successfully", func() {
			body := e.GET("/api/v1alpha1/namespaces/{namespace}/servings/{servingID}", ns, servingName).
				Expect().Status(http.StatusOK).Body().Raw()
			sdep := seldonv1.SeldonDeployment{}
			json.Unmarshal([]byte(body),&sdep)
			Expect(sdep.Spec.Predictors[0].Graph.ModelURI).Should(ContainSubstring(project+"/"+model+":"+version))
			Expect(sdep.ObjectMeta.Name).Should(Equal(servingName))
			Expect(sdep.Spec.Predictors[0].ComponentSpecs[0].Spec.Containers[0].Resources.Requests.Cpu().String()).Should(Equal("1"))
			Expect(sdep.Spec.Predictors[0].ComponentSpecs[0].Spec.Containers[0].Resources.Requests.Memory().String()).Should(Equal("1Gi"))
			Expect(sdep.Spec.Predictors[0].ComponentSpecs[0].Spec.Containers[0].Resources.Limits.Cpu().String()).Should(Equal("1"))
			Expect(sdep.Spec.Predictors[0].ComponentSpecs[0].Spec.Containers[0].Resources.Limits.Memory().String()).Should(Equal("1Gi"))
		})

		It("should delete the serving successfully", func() {
			e.DELETE("/api/v1alpha1/namespaces/{namespace}/servings/{servingID}", ns, servingName).
				Expect().Status(http.StatusNoContent)
		})

		It("Shouldn't get the deleted serving", func() {
			Eventually(func() int {
				return e.GET("/api/v1alpha1/namespaces/{namespace}/servings/{servingID}", ns, servingName).
					Expect().Raw().StatusCode
			}, timeout, interval).Should(Equal(http.StatusNotFound))

		})

	})
})
