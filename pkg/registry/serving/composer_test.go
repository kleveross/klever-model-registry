package serving

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	seldonv1 "github.com/seldonio/seldon-core/operator/apis/machinelearning.seldon.io/v1"
	"github.com/spf13/viper"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/kleveross/klever-model-registry/pkg/common"
)

var sdepSingleGraph *seldonv1.SeldonDeployment
var sdepDoubleGraph *seldonv1.SeldonDeployment
var sdepCustomImageGraph *seldonv1.SeldonDeployment

var _ = Describe("Composer", func() {

	It("Should compose single graph successfully", func() {
		err := Compose(sdepSingleGraph)
		Expect(err).To(BeNil())

		Expect(sdepSingleGraph.Spec.Predictors[0].Name).Should(Equal(sdepSingleGraph.Spec.Predictors[0].Graph.Name))

		Expect(len(sdepSingleGraph.Spec.Predictors)).Should(Equal(1))
		Expect(len(sdepSingleGraph.Spec.Predictors[0].ComponentSpecs[0].Spec.Containers[0].Ports)).Should(Equal(2))
		Expect(sdepSingleGraph.Spec.Predictors[0].ComponentSpecs[0].Spec.Containers[0].Ports[0].ContainerPort).Should(Equal(int32(8000)))

		Expect(len(sdepSingleGraph.Spec.Predictors[0].ComponentSpecs[0].Spec.Containers[0].VolumeMounts)).Should(Equal(1))
		Expect(sdepSingleGraph.Spec.Predictors[0].ComponentSpecs[0].Spec.Containers[0].VolumeMounts[0].MountPath).Should(Equal("/mnt/sdep-name"))

		Expect(len(sdepSingleGraph.Spec.Predictors[0].ComponentSpecs[0].Spec.InitContainers[0].VolumeMounts)).Should(Equal(1))

	})

	It("Should compose double graph successfully", func() {
		err := Compose(sdepDoubleGraph)
		Expect(err).To(BeNil())

		Expect(sdepDoubleGraph.Spec.Predictors[0].Name).Should(Equal(sdepDoubleGraph.Spec.Predictors[0].Graph.Name))

		Expect(len(sdepDoubleGraph.Spec.Predictors)).Should(Equal(2))
		Expect(len(sdepDoubleGraph.Spec.Predictors[0].ComponentSpecs[0].Spec.Containers[0].Ports)).Should(Equal(2))
		Expect(len(sdepDoubleGraph.Spec.Predictors[1].ComponentSpecs[0].Spec.Containers[0].Ports)).Should(Equal(1))
		Expect(sdepDoubleGraph.Spec.Predictors[0].ComponentSpecs[0].Spec.Containers[0].Ports[0].ContainerPort).Should(Equal(int32(8000)))
		Expect(sdepDoubleGraph.Spec.Predictors[1].ComponentSpecs[0].Spec.Containers[0].Ports[0].ContainerPort).Should(Equal(int32(8000)))

		Expect(len(sdepDoubleGraph.Spec.Predictors[0].ComponentSpecs[0].Spec.InitContainers[0].VolumeMounts)).Should(Equal(1))
		Expect(len(sdepDoubleGraph.Spec.Predictors[1].ComponentSpecs[0].Spec.InitContainers[0].VolumeMounts)).Should(Equal(1))
		Expect(sdepDoubleGraph.Spec.Predictors[0].ComponentSpecs[0].Spec.Containers[0].VolumeMounts[0].MountPath).Should(Equal("/mnt/sdep-name"))
		Expect(sdepDoubleGraph.Spec.Predictors[1].ComponentSpecs[0].Spec.Containers[0].VolumeMounts[0].MountPath).Should(Equal("/mnt/sdep-name"))
	})

	It("Should compose custom image graph successfully", func() {
		err := Compose(sdepCustomImageGraph)
		Expect(err).To(BeNil())

		Expect(sdepCustomImageGraph.Spec.Predictors[0].Name).Should(Equal(sdepCustomImageGraph.Spec.Predictors[0].Graph.Name))

		Expect(len(sdepCustomImageGraph.Spec.Predictors)).Should(Equal(1))
		Expect(len(sdepCustomImageGraph.Spec.Predictors[0].ComponentSpecs[0].Spec.Containers[0].Ports)).Should(Equal(1))
		Expect(sdepCustomImageGraph.Spec.Predictors[0].ComponentSpecs[0].Spec.Containers[0].Ports[0].ContainerPort).Should(Equal(int32(8000)))

		Expect(len(sdepCustomImageGraph.Spec.Predictors[0].ComponentSpecs[0].Spec.Containers[0].VolumeMounts)).Should(Equal(1))
		Expect(sdepCustomImageGraph.Spec.Predictors[0].ComponentSpecs[0].Spec.Containers[0].VolumeMounts[0].MountPath).Should(Equal("/model_store/sdep-name"))

		Expect(len(sdepCustomImageGraph.Spec.Predictors[0].ComponentSpecs[0].Spec.InitContainers[0].VolumeMounts)).Should(Equal(1))
	})
})

var _ = BeforeEach(func() {
	viper.Set("MODEL_INITIALIZER_CPU", "1")
	viper.Set("MODEL_INITIALIZER_MEM", "1Gi")

	genTestResource := func() corev1.ResourceList {
		resourcesList := make(corev1.ResourceList)
		cpuQuantity, _ := resource.ParseQuantity("1")
		resourcesList[corev1.ResourceCPU] = cpuQuantity
		memQuantity, _ := resource.ParseQuantity("1Gi")
		resourcesList[corev1.ResourceMemory] = memQuantity

		return resourcesList
	}

	sdepSingleGraph = &seldonv1.SeldonDeployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "sdep-name",
			Namespace: "default",
		},
		Spec: seldonv1.SeldonDeploymentSpec{
			Name: "deployment-name",
			Predictors: []seldonv1.PredictorSpec{
				{
					ComponentSpecs: []*seldonv1.SeldonPodSpec{
						{
							Metadata: metav1.ObjectMeta{
								Name: "graph1",
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
						Name:               "graph1",
						ModelURI:           "harbor-harbor-core.kleveross-system/release/savedmodel:v1",
						ServiceAccountName: "default",
						Parameters: []seldonv1.Parameter{
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

	sdepCustomImageGraph = &seldonv1.SeldonDeployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "sdep-name",
			Namespace: "default",
		},
		Spec: seldonv1.SeldonDeploymentSpec{
			Name: "deployment-name",
			Predictors: []seldonv1.PredictorSpec{
				{
					ComponentSpecs: []*seldonv1.SeldonPodSpec{
						{
							Metadata: metav1.ObjectMeta{
								Name: "graph1",
							},
							Spec: corev1.PodSpec{
								Containers: []corev1.Container{
									{
										Image: "lever-dev.cargo.io/release/tritonserver:v0.2.0",
										Env: []corev1.EnvVar{
											{
												Name:  "MODEL_STORE",
												Value: "/model_store",
											},
										},
										Ports: []corev1.ContainerPort{
											{
												Name:          "http",
												ContainerPort: 8000,
												Protocol:      corev1.ProtocolTCP,
											},
										},
										Command: []string{"/entrypoint.sh"},
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
						Name:               "graph1",
						ModelURI:           "harbor-harbor-core.kleveross-system/release/savedmodel:v1",
						ServiceAccountName: "default",
						Parameters: []seldonv1.Parameter{
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
					ComponentSpecs: []*seldonv1.SeldonPodSpec{
						{
							Metadata: metav1.ObjectMeta{
								Name: "graph1",
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
						Name:               "graph1",
						ModelURI:           "harbor-harbor-core.kleveross-system/release/savedmodel:v1",
						ServiceAccountName: "default",
						Parameters: []seldonv1.Parameter{
							{
								Name:  "format",
								Value: "SavedModel",
							},
						},
					},
				},
				{
					ComponentSpecs: []*seldonv1.SeldonPodSpec{
						{
							Metadata: metav1.ObjectMeta{
								Name: "graph2",
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
						Name:               "graph2",
						ModelURI:           "harbor-harbor-core.kleveross-system/release/pmml:v1",
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
})

func Test_rewriteModelURI(t *testing.T) {
	common.ORMBDomain = "domain"

	type args struct {
		uri string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "rewrite successfully",
			args: args{
				uri: "repo/savedmodel:v1",
			},
			want: "domain/repo/savedmodel:v1",
		},
		{
			name: "rewrite failed",
			args: args{
				uri: "test/repo/savedmodel:v1",
			},
			want: "test/repo/savedmodel:v1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := rewriteModelURI(tt.args.uri); got != tt.want {
				t.Errorf("rewriteModelURI() = %v, want %v", got, tt.want)
			}
		})
	}
}
