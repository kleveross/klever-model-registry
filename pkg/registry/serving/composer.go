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

	// envMLServerLImage is the preset image for mlserver.
	envMLServerImage = "MLSERVER_IMAGE"

	// envModelInitializerImage is the preset image for model initializer.
	envModelInitializerImage = "MODEL_INITIALIZER_IMAGE"

	// envModelInitializerCPU is the cpu config for model-initializer container.
	envModelInitializerCPU = "MODEL_INITIALIZER_CPU"

	// envModelInitializerMem is the cpu config for model-initializer container.
	envModelInitializerMem = "MODEL_INITIALIZER_MEM"

	// envNvidiaVisibleDevices set the env empty string, the container will not use GPU.
	envNvidiaVisibleDevices = "NVIDIA_VISIBLE_DEVICES"

	// envSchedulerName will set podSpec's SchedulerName.
	envSchedulerName = "SCHEDULER_NAME"

	// defaultInferenceHTTPPort is default port for http.
	defaultInferenceHTTPPort = 8000

	// defaultInferenceGRPCPort is default port for grpc.
	defaultInferenceGRPCPort = 8001

	// defaultMLServerHTTPPort is default port for http.
	defaultMLServerHTTPPort = 8080

	// defaultMLServerGRPCPort is default port for grpc.
	defaultMLServerGRPCPort = 8081

	// modelStorePath is trtserver param --model-repository path.
	modelStorePath = "/mnt"
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
		modelTag, err := getModelTag(p.Graph.ModelURI)
		if err != nil {
			return err
		}
		podName := modelTag
		containerName := sdep.Spec.Predictors[i].Graph.Name
		sdep.Spec.Predictors[i].Name = sdep.Spec.Predictors[i].Graph.Name

		resources, gpuFlag, err := getRuntimeResource(&sdep.Spec.Predictors[i].Graph)
		if err != nil {
			return err
		}

		modelFormat := getModelFormat(&sdep.Spec.Predictors[i].Graph)
		// Must set probe, otherwise the default probe by seldon's webhook will cause error.
		probe := getProbe(modelFormat, sdep.Name)

		ports := getUserContainerPorts(modelFormat)

		image := getUserContainerImage(modelFormat)
		// compose user containers
		container := corev1.Container{
			Name:            containerName,
			Image:           image,
			ImagePullPolicy: corev1.PullAlways,
			Env: []corev1.EnvVar{
				{
					Name:  "MODEL_STORE",
					Value: modelStorePath,
				},
				{
					Name:  "SERVING_NAME",
					Value: sdep.Name,
				},
			},
			// Must set ports, otherwise it will can not traffic diversion in unique port(default: 8000) for multi deployment.
			// please refer https://github.com/SeldonIO/seldon-core/blob/master/operator/apis/machinelearning.seldon.io/v1/seldondeployment_webhook.go#L142-L145
			Ports: ports,
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

		// If not set GPU resource, must set env key is equal "NVIDIA_VISIBLE_DEVICES" and value is empty string.
		if !gpuFlag {
			container.Env = append(container.Env, corev1.EnvVar{
				Name: envNvidiaVisibleDevices,
			})
		}

		seldPodSpec := &seldonv1.SeldonPodSpec{
			Metadata: metav1.ObjectMeta{
				Name: podName,
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
		}
		composeSchedulerName(seldPodSpec)
		sdep.Spec.Predictors[i].ComponentSpecs = []*seldonv1.SeldonPodSpec{seldPodSpec}
	}

	composeInitContainer(sdep)

	return nil
}

func getUserContainerPorts(format string) []corev1.ContainerPort {
	if format == string(modeljobsv1alpha1.FormatSKLearn) || format == string(modeljobsv1alpha1.FormatXGBoost) {
		ports := []corev1.ContainerPort{
			{
				Name:          "http",
				Protocol:      corev1.ProtocolTCP,
				ContainerPort: defaultMLServerHTTPPort,
			},
			{
				Name:          "grpc",
				Protocol:      corev1.ProtocolTCP,
				ContainerPort: defaultMLServerGRPCPort,
			},
		}
		return ports
	}

	ports := []corev1.ContainerPort{
		{
			Name:          "http",
			Protocol:      corev1.ProtocolTCP,
			ContainerPort: defaultInferenceHTTPPort,
		},
	}

	if format != string(modeljobsv1alpha1.FormatPMML) {
		ports = append(ports, corev1.ContainerPort{
			Name:          "grpc",
			Protocol:      corev1.ProtocolTCP,
			ContainerPort: defaultInferenceGRPCPort,
		})
	}

	return ports
}

// getProbe generate readiness and liveiness.
func getProbe(format, servingName string) *corev1.Probe {
	path := fmt.Sprintf("/api/status/%v", servingName)
	port := defaultInferenceHTTPPort
	if format == string(modeljobsv1alpha1.FormatPMML) {
		path = fmt.Sprintf("/openscoring/model/%v", servingName)
	} else if format == string(modeljobsv1alpha1.FormatSKLearn) || format == string(modeljobsv1alpha1.FormatXGBoost) {
		path = fmt.Sprintf("/v2/models/%v/ready", servingName)
		port = defaultMLServerHTTPPort
	}

	return &corev1.Probe{
		Handler: corev1.Handler{
			HTTPGet: &corev1.HTTPGetAction{
				Path: path,
				Port: intstr.IntOrString{
					Type:   intstr.Int,
					IntVal: int32(port),
				},
				Scheme: corev1.URISchemeHTTP,
			},
		},
	}
}

// getModelFormat get model format from Graph.Parameters, eg:
// "parameters": [
// 	{
// 		"name": "format",
// 		"value": "SavedModel"
// 	},
// 	...
// ]
func getModelFormat(pu *seldonv1.PredictiveUnit) string {
	for _, p := range pu.Parameters {
		if p.Name == "format" {
			return p.Value
		}
	}

	return ""
}

// getUserContainerImage get image by different model format.
func getUserContainerImage(format string) string {
	if format == string(modeljobsv1alpha1.FormatPMML) {
		return viper.GetString(envPMMLServingImage)
	}
	if format == string(modeljobsv1alpha1.FormatSKLearn) || format == string(modeljobsv1alpha1.FormatXGBoost) {
		return viper.GetString(envMLServerImage)
	}

	return viper.GetString(envTRTServingImage)
}

// getRuntimeResource get resource from Graph.Parameters, eg:
// "parameters": [
// 	{
// 		"name": "cpu",
// 		"value": "1"
// 	},
// 	{
// 		"name": "mem",
// 		"value": "2Gi"
// 	},
// 	...
// ]
func getRuntimeResource(pu *seldonv1.PredictiveUnit) (*corev1.ResourceRequirements, bool, error) {
	cpu := ""
	mem := ""
	gpuType := ""
	gpuNum := ""
	gpuFlag := false
	for _, p := range pu.Parameters {
		if p.Name == "cpu" {
			cpu = p.Value
		}
		if p.Name == "mem" {
			mem = p.Value
		}
		if p.Name == "gpuType" {
			gpuType = p.Value
		}
		if p.Name == "gpuNum" {
			gpuNum = p.Value
		}
	}

	resourcesList := make(corev1.ResourceList)
	cpuQuantity, err := resource.ParseQuantity(cpu)
	if err != nil {
		return nil, gpuFlag, err
	}
	resourcesList[corev1.ResourceCPU] = cpuQuantity

	memQuantity, err := resource.ParseQuantity(mem)
	if err != nil {
		return nil, gpuFlag, err
	}
	resourcesList[corev1.ResourceMemory] = memQuantity

	// Support gpu scheduling, the detail please refer https://kubernetes.io/zh/docs/tasks/manage-gpus/scheduling-gpus/
	if gpuType != "" && gpuNum != "" {
		gpuNumQuantity, err := resource.ParseQuantity(gpuNum)
		if err != nil {
			return nil, gpuFlag, err
		}
		resourcesList[corev1.ResourceName(gpuType)] = gpuNumQuantity
		gpuFlag = true
	}

	return &corev1.ResourceRequirements{
		Limits:   resourcesList,
		Requests: resourcesList,
	}, gpuFlag, nil
}

// getModelTag gets model tag, eg: harbor.demo.io/release/savedmodel:v1, it will return `v1`.
func getModelTag(modelUri string) (string, error) {
	modelURISlice := strings.Split(modelUri, ":")
	if len(modelURISlice) < 2 {
		return "", fmt.Errorf("modelUri's format is error")
	}

	return modelURISlice[len(modelURISlice)-1], nil
}

// getModelMountPath will generate model mount path in container,
// ormb-storage-initializer will pull and export model to this path.
func getModelMountPath(servingName string) string {
	return fmt.Sprintf("%v/%v", modelStorePath, servingName)
}

// composeModelInitailzerContainerResource get the default resource config.
func composeModelInitailzerContainerResource(container *corev1.Container) error {
	cpu := viper.GetString(envModelInitializerCPU)
	mem := viper.GetString(envModelInitializerMem)
	if cpu != "" && mem != "" {
		resourcesList := make(corev1.ResourceList)
		cpuQuantity, err := resource.ParseQuantity(cpu)
		if err != nil {
			return err
		}
		resourcesList[corev1.ResourceCPU] = cpuQuantity

		memQuantity, err := resource.ParseQuantity(mem)
		if err != nil {
			return err
		}
		resourcesList[corev1.ResourceMemory] = memQuantity

		container.Resources = corev1.ResourceRequirements{
			Limits:   resourcesList,
			Requests: resourcesList,
		}
		return nil
	}

	return nil
}

func composeInitContainer(sdep *seldonv1.SeldonDeployment) error {
	modelMountPath := getModelMountPath(sdep.Name)

	for _, p := range sdep.Spec.Predictors {
		// simple model serving, the number of ComponentSpecs is 1
		if len(p.ComponentSpecs) != 1 {
			return fmt.Errorf("too many or too less componentspecs")
		}

		container := &corev1.Container{
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

		err := composeModelInitailzerContainerResource(container)
		if err != nil {
			return err
		}

		p.ComponentSpecs[0].Spec.InitContainers = []corev1.Container{*container}
	}

	return nil
}

// composeSchedulerName set container for inference task.
func composeSchedulerName(seldonPodSpec *seldonv1.SeldonPodSpec) {
	schedulerName := viper.GetString(envSchedulerName)
	if schedulerName == "" {
		return
	}
	seldonPodSpec.Spec.SchedulerName = schedulerName
}
