package serving

import (
	"fmt"
	"strings"

	seldonv1 "github.com/seldonio/seldon-core/operator/apis/machinelearning.seldon.io/v1"
	"github.com/spf13/viper"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	modeljobsv1alpha1 "github.com/kleveross/klever-model-registry/pkg/apis/modeljob/v1alpha1"
	"github.com/kleveross/klever-model-registry/pkg/common"
)

const (
	// modelSharedMountName is a shared dir for initContainer and userContainer,
	// the model from harbor by ormb pull will store in the mount point.
	modelSharedMountName = "models-mnt"

	// envTRTServingImage is the preset image for trtserver.
	envTRTServingImage = "TRT_SERVING_IMAGE"

	// envPMMLServingImage is the preset image for pmml.
	envPMMLServingImage = "PMML_SERVING_IMAGE"

	// envModelInitializerImage is the preset image for model initializer.
	envModelInitializerImage = "MODEL_INITIAIZER_IMAGE"
)

func Compose(sdep *seldonv1.SeldonDeployment) error {
	sdep.Spec.Name = sdep.ObjectMeta.Name
	sdep.Spec.Transport = seldonv1.TransportRest

	modelMountPath := getModelMountPath(sdep.Name)

	for i, p := range sdep.Spec.Predictors {
		if sdep.Spec.Predictors[i].Annotations == nil {
			sdep.Spec.Predictors[i].Annotations = make(map[string]string)
		}

		// use no-engine mode
		sdep.Spec.Predictors[i].Annotations[seldonv1.ANNOTATION_NO_ENGINE] = "true"
		graphName, err := getModelTag(p.Graph.ModelURI)
		if err != nil {
			return err
		}
		sdep.Spec.Predictors[i].Name = graphName
		sdep.Spec.Predictors[i].Graph.Name = graphName

		resources, err := getRuntimeResource(&sdep.Spec.Predictors[i].Graph)
		if err != nil {
			return err
		}

		modelFormat := getModelFormat(&sdep.Spec.Predictors[i].Graph)
		probe := getProbe(modelFormat, sdep.Name)

		image := getUserContainerImage(&sdep.Spec.Predictors[i].Graph)
		// compose user containers
		container := corev1.Container{
			Name:  graphName,
			Image: image,
			Env: []corev1.EnvVar{
				{
					Name:  "MODEL_STORE",
					Value: modelMountPath,
				},
				{
					Name:  "SERVING_NAME",
					Value: sdep.Name,
				},
			},
			VolumeMounts: []corev1.VolumeMount{
				{
					Name:      modelSharedMountName,
					MountPath: modelMountPath,
				},
			},
			ReadinessProbe: probe,
			LivenessProbe:  probe,
			// not support graph now, so there are one container only.
			Resources: *resources,
		}

		sdep.Spec.Predictors[i].ComponentSpecs = []*seldonv1.SeldonPodSpec{
			&seldonv1.SeldonPodSpec{
				Metadata: metav1.ObjectMeta{
					Name: graphName,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{container},
					Volumes: []corev1.Volume{
						{
							Name: modelSharedMountName,
							VolumeSource: corev1.VolumeSource{
								EmptyDir: &corev1.EmptyDirVolumeSource{},
							},
						},
					},
				},
			},
		}
	}

	composeInitContainer(sdep)

	return nil
}

// getProbe generate readiness and liveiness
func getProbe(format, servingName string) *corev1.Probe {
	path := fmt.Sprintf("/api/status/%v", servingName)
	if format == string(modeljobsv1alpha1.FormatPMML) {
		path = fmt.Sprintf("/openscoring/model/%v", servingName)
	}

	return &corev1.Probe{
		Handler: corev1.Handler{
			HTTPGet: &corev1.HTTPGetAction{
				Path: path,
				Port: intstr.IntOrString{
					Type:   intstr.Int,
					IntVal: 8000,
				},
				Scheme: corev1.URISchemeHTTP,
			},
		},
	}
}

func getModelFormat(pu *seldonv1.PredictiveUnit) string {
	for _, p := range pu.Parameters {
		if p.Name == "format" {
			return p.Value
		}
	}

	return ""
}

func getUserContainerImage(pu *seldonv1.PredictiveUnit) string {
	format := getModelFormat(pu)
	if format == string(modeljobsv1alpha1.FormatPMML) {
		return viper.GetString(envPMMLServingImage)
	}

	return viper.GetString(envTRTServingImage)
}

func getRuntimeResource(pu *seldonv1.PredictiveUnit) (*corev1.ResourceRequirements, error) {
	cpu := ""
	mem := ""
	for _, p := range pu.Parameters {
		if p.Name == "cpu" {
			cpu = p.Value
		}
		if p.Name == "mem" {
			mem = p.Value
		}
	}

	resourcesList := make(corev1.ResourceList)
	cpuQuantity, err := resource.ParseQuantity(cpu)
	if err != nil {
		return nil, err
	}
	resourcesList[corev1.ResourceCPU] = cpuQuantity

	memQuantity, err := resource.ParseQuantity(mem)
	if err != nil {
		return nil, err
	}
	resourcesList[corev1.ResourceMemory] = memQuantity
	// TODO(@simon-cj): Support GPU resource

	return &corev1.ResourceRequirements{
		Limits:   resourcesList,
		Requests: resourcesList,
	}, nil
}

func getModelTag(modelUri string) (string, error) {
	modelURISlice := strings.Split(modelUri, ":")
	if len(modelURISlice) < 2 {
		return "", fmt.Errorf("modelUri's format is error")
	}

	return modelURISlice[len(modelURISlice)-1], nil
}

func getModelMountPath(servingName string) string {
	return fmt.Sprintf("/mnt/%v", servingName)
}

func composeInitContainer(sdep *seldonv1.SeldonDeployment) error {
	modelMountPath := getModelMountPath(sdep.Name)

	for _, p := range sdep.Spec.Predictors {
		// simple model serving, the number of ComponentSpecs is 1
		if len(p.ComponentSpecs) != 1 {
			return fmt.Errorf("too many or too less componentspecs")
		}

		container := corev1.Container{
			Name:  "model-initializer",
			Image: viper.GetString(envModelInitializerImage),
			Args:  []string{p.Graph.ModelURI, modelMountPath},
			// Get username and password from environment
			// Here AWS_SECRET_ACCESS_KEY and AWS_ACCESS_KEY_ID are used
			// because Seldon Core does not support renaming the environment variable name.
			// it is used in ormb-storage-initializer.
			// please refenence https://github.com/kleveross/ormb/blob/master/cmd/ormb-storage-initializer/cmd/pull-and-export.go#L47
			Env: []corev1.EnvVar{
				{
					Name:  "AWS_ACCESS_KEY_ID",
					Value: common.ORMBUserName,
				},
				{
					Name:  "AWS_SECRET_ACCESS_KEY",
					Value: common.ORMBPassword,
				},
				{
					Name:  "ROOTPATH",
					Value: "/mnt",
				},
			},
			VolumeMounts: []corev1.VolumeMount{
				{
					Name:      modelSharedMountName,
					MountPath: modelMountPath,
				},
			},
		}
		p.ComponentSpecs[0].Spec.InitContainers = []corev1.Container{container}
	}
	return nil
}
