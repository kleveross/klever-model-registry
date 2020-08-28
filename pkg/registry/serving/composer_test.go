package serving_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	seldonv1 "github.com/seldonio/seldon-core/operator/apis/machinelearning.seldon.io/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/kleveross/klever-model-registry/pkg/registry/serving"
)

var sdepSingleGraph *seldonv1.SeldonDeployment
var sdepDoubleGraph *seldonv1.SeldonDeployment

var _ = Describe("Composer", func() {

	It("Should compose single graph successfully", func() {
		err := serving.Compose(sdepSingleGraph)
		Expect(err).To(BeNil())

		Expect(len(sdepSingleGraph.Spec.Predictors)).Should(Equal(1))
		Expect(len(sdepSingleGraph.Spec.Predictors[0].ComponentSpecs[0].Spec.Containers[0].Ports)).Should(Equal(1))
		Expect(sdepSingleGraph.Spec.Predictors[0].ComponentSpecs[0].Spec.Containers[0].Ports[0].ContainerPort).Should(Equal(int32(8000)))

		Expect(len(sdepSingleGraph.Spec.Predictors[0].ComponentSpecs[0].Spec.Containers[0].VolumeMounts)).Should(Equal(1))
		Expect(sdepSingleGraph.Spec.Predictors[0].ComponentSpecs[0].Spec.Containers[0].VolumeMounts[0].MountPath).Should(Equal("/mnt/sdep-name"))

		Expect(len(sdepSingleGraph.Spec.Predictors[0].ComponentSpecs[0].Spec.InitContainers[0].VolumeMounts)).Should(Equal(1))

	})

	It("Should compose double graph successfully", func() {
		err := serving.Compose(sdepDoubleGraph)
		Expect(err).To(BeNil())

		Expect(len(sdepDoubleGraph.Spec.Predictors)).Should(Equal(2))
		Expect(len(sdepDoubleGraph.Spec.Predictors[0].ComponentSpecs[0].Spec.Containers[0].Ports)).Should(Equal(1))
		Expect(len(sdepDoubleGraph.Spec.Predictors[1].ComponentSpecs[0].Spec.Containers[0].Ports)).Should(Equal(1))
		Expect(sdepDoubleGraph.Spec.Predictors[0].ComponentSpecs[0].Spec.Containers[0].Ports[0].ContainerPort).Should(Equal(int32(8000)))
		Expect(sdepDoubleGraph.Spec.Predictors[1].ComponentSpecs[0].Spec.Containers[0].Ports[0].ContainerPort).Should(Equal(int32(8000)))

		Expect(len(sdepDoubleGraph.Spec.Predictors[0].ComponentSpecs[0].Spec.InitContainers[0].VolumeMounts)).Should(Equal(1))
		Expect(len(sdepDoubleGraph.Spec.Predictors[1].ComponentSpecs[0].Spec.InitContainers[0].VolumeMounts)).Should(Equal(1))
		Expect(sdepDoubleGraph.Spec.Predictors[0].ComponentSpecs[0].Spec.Containers[0].VolumeMounts[0].MountPath).Should(Equal("/mnt/sdep-name"))
		Expect(sdepDoubleGraph.Spec.Predictors[1].ComponentSpecs[0].Spec.Containers[0].VolumeMounts[0].MountPath).Should(Equal("/mnt/sdep-name"))
	})
})

var _ = BeforeEach(func() {
	sdepSingleGraph = &seldonv1.SeldonDeployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "sdep-name",
			Namespace: "default",
		},
		Spec: seldonv1.SeldonDeploymentSpec{
			Name: "deployment-name",
			Predictors: []seldonv1.PredictorSpec{
				{
					Graph: seldonv1.PredictiveUnit{
						Name:               "graph1",
						ModelURI:           "harbor-harbor-core.kleveross-system/release/savedmodel:v1",
						ServiceAccountName: "default",
						Endpoint: &seldonv1.Endpoint{
							Type: seldonv1.REST,
						},
						Parameters: []seldonv1.Parameter{
							{
								Name:  "cpu",
								Value: "1",
							},
							{
								Name:  "mem",
								Value: "2Gi",
							},
							{
								Name:  "format",
								Value: "SavedModel",
							},
						},
					},
				},
			},
		},
	}

	sdepDoubleGraph = &seldonv1.SeldonDeployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "sdep-name",
			Namespace: "default",
		},
		Spec: seldonv1.SeldonDeploymentSpec{
			Name: "deployment-name",
			Predictors: []seldonv1.PredictorSpec{
				{
					Graph: seldonv1.PredictiveUnit{
						Name:               "graph1",
						ModelURI:           "harbor-harbor-core.kleveross-system/release/savedmodel:v1",
						ServiceAccountName: "default",
						Endpoint: &seldonv1.Endpoint{
							Type: seldonv1.REST,
						},
						Parameters: []seldonv1.Parameter{
							{
								Name:  "cpu",
								Value: "1",
							},
							{
								Name:  "mem",
								Value: "2Gi",
							},
							{
								Name:  "format",
								Value: "SavedModel",
							},
						},
					},
				},
				{
					Graph: seldonv1.PredictiveUnit{
						Name:               "graph2",
						ModelURI:           "harbor-harbor-core.kleveross-system/release/pmml:v1",
						ServiceAccountName: "default",
						Endpoint: &seldonv1.Endpoint{
							Type: seldonv1.REST,
						},
						Parameters: []seldonv1.Parameter{
							{
								Name:  "cpu",
								Value: "1",
							},
							{
								Name:  "mem",
								Value: "2Gi",
							},
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
})
