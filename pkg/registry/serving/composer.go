package serving

import (
	"fmt"
	"strings"

	"github.com/caicloud/nirvana/log"
	seldonv1 "github.com/seldonio/seldon-core/operator/apis/machinelearning.seldon.io/v1"
	"github.com/spf13/viper"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
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

	// envModelStorePath is for custome image.
	// if the image is not empty, must set MODEL_STORE in each SeldonPodSpec's env.
	envModelStorePath = "MODEL_STORE"

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

// validateComponentSpecs validate basic infomation in CRD, now, we are not support multi graph,
// so the length for ComponentSpecs and Containers must be equal 1.
func validateComponentSpecs(p *seldonv1.PredictorSpec) error {
	if len(p.ComponentSpecs) != 1 {
		log.Warningf("Component's length must be equal 1 for simple model serving, the actual ComponentSpecs is %v", p.ComponentSpecs)
		return fmt.Errorf("Component's length must be equal 1 for simple model serving")
	}

	if len(p.ComponentSpecs[0].Spec.Containers) != 1 {
		log.Warningf("Component's length must be equal 1 for simple model serving, the actual Containers is %v", p.ComponentSpecs[0].Spec.Containers)
		return fmt.Errorf("Container's length must be equal 1 for simple model serving")
	}

	return nil
}

func Compose(sdep *seldonv1.SeldonDeployment) error {
	sdep.Spec.Name = sdep.ObjectMeta.Name

	for _, p := range sdep.Spec.Predictors {
		if !seldonv1.IsPrepack(&p.Graph) {
			if err := validateComponentSpecs(&p); err != nil {
				return err
			}

			componentSpecMap := getComponentsMap(p.ComponentSpecs)

			// Compose user containers
			if p.ComponentSpecs[0].Spec.Containers[0].Image != "" {
				// For custome image
				err := composeCustomesUserContainer(sdep, &p.Graph, componentSpecMap)
				if err != nil {
					return err
				}
			} else {
				// For default image
				err := composeDefaultUserContainer(sdep, &p.Graph, componentSpecMap)
				if err != nil {
					return err
				}
			}

			// Conpose init container for pod
			composeInitContainer(sdep, &p)
		}

		// Setup no-engine mode
		setupNoEngineMode(&p)
	}

	return nil
}

// getComponentsMap convert []*seldonv1.SeldonPodSpec as map struct, the key is SeldonPodSpec.Metadata.Name, the value is SeldonPodSpec.
func getComponentsMap(componentSpecs []*seldonv1.SeldonPodSpec) map[string]*seldonv1.SeldonPodSpec {
	componentSpecMap := make(map[string]*seldonv1.SeldonPodSpec)
	for _, cs := range componentSpecs {
		componentSpecMap[cs.Metadata.Name] = cs
	}

	return componentSpecMap
}

func setupNoEngineMode(pu *seldonv1.PredictorSpec) {
	if pu.Annotations == nil {
		pu.Annotations = make(map[string]string)
	}
	// use no-engine mode
	pu.Annotations[seldonv1.ANNOTATION_NO_ENGINE] = "true"

	pu.Name = pu.Graph.Name
}

// composeCustomeUserContainer compose user container for custome image
func composeCustomesUserContainer(sdep *seldonv1.SeldonDeployment, pu *seldonv1.PredictiveUnit, componentSpecMap map[string]*seldonv1.SeldonPodSpec) error {
	// Find the Graph node's PodSpec and Container
	var seldonPodSpec *seldonv1.SeldonPodSpec
	var ok bool
	if seldonPodSpec, ok = componentSpecMap[pu.Name]; !ok {
		return fmt.Errorf("can't find ComponentSpec for graph %v", pu.Name)
	}
	if len(seldonPodSpec.Spec.Containers) == 0 {
		return fmt.Errorf("container length is zero")
	}
	container := &seldonPodSpec.Spec.Containers[0]
	container.Name = pu.Name

	// Set default env.
	if len(container.Env) == 0 {
		container.Env = []corev1.EnvVar{}
	}
	container.Env = append(container.Env, []corev1.EnvVar{
		{
			Name:  "SERVING_NAME",
			Value: sdep.Name,
		},
	}...)

	// Set Probe
	probe := &corev1.Probe{
		Handler: corev1.Handler{
			Exec: &corev1.ExecAction{
				Command: []string{"ls", "/"},
			},
			HTTPGet:   nil,
			TCPSocket: nil,
		},
		TimeoutSeconds:   5,
		FailureThreshold: 5,
	}
	container.ReadinessProbe = probe
	container.LivenessProbe = probe

	// Set VolumeMounts
	if len(container.VolumeMounts) == 0 {
		container.VolumeMounts = []corev1.VolumeMount{}
	}
	modelMountPath := getModelMountPath(container, sdep.Name)
	container.VolumeMounts = append(container.VolumeMounts, []corev1.VolumeMount{
		{
			Name:      modelSharedMountName,
			MountPath: modelMountPath,
		},
	}...)

	for idx := range pu.Children {
		err := composeCustomesUserContainer(sdep, &pu.Children[idx], componentSpecMap)
		if err != nil {
			return err
		}
	}

	return nil
}

// composeDefaultUserContainer compose user container for default image
func composeDefaultUserContainer(sdep *seldonv1.SeldonDeployment, pu *seldonv1.PredictiveUnit, componentSpecMap map[string]*seldonv1.SeldonPodSpec) error {
	// Find the Graph node's PodSpec and Container
	var seldonPodSpec *seldonv1.SeldonPodSpec
	var ok bool
	if seldonPodSpec, ok = componentSpecMap[pu.Name]; !ok {
		return fmt.Errorf("can't find ComponentSpec for graph %v", pu.Name)
	}
	if len(seldonPodSpec.Spec.Containers) == 0 {
		return fmt.Errorf("container length is zero")
	}
	container := &seldonPodSpec.Spec.Containers[0]

	container.Name = pu.Name
	modelFormat := getModelFormat(pu)
	image := getUserContainerImage(modelFormat)
	container.Image = image

	// Must set ports, otherwise it will can not traffic diversion in unique port(default: 8000) for multi deployment.
	// please refer https://github.com/SeldonIO/seldon-core/blob/master/operator/apis/machinelearning.seldon.io/v1/seldondeployment_webhook.go#L142-L145
	ports := getDefaultUserContainerPorts(modelFormat)
	container.Ports = ports

	// Must set probe, otherwise the default probe by seldon's webhook will cause error.
	probe := getDefaultProbe(modelFormat, sdep.Name)
	container.ReadinessProbe = probe
	container.LivenessProbe = probe

	// Set default env.
	if len(container.Env) == 0 {
		container.Env = []corev1.EnvVar{}
	}
	container.Env = append(container.Env, []corev1.EnvVar{
		{
			Name:  envModelStorePath,
			Value: modelStorePath,
		},
		{
			Name:  "SERVING_NAME",
			Value: sdep.Name,
		},
	}...)
	// If not set GPU resource, must set env key is equal "NVIDIA_VISIBLE_DEVICES" and value is empty string.
	if getGPUAmount(container.Resources) == 0 {
		container.Env = append(container.Env, corev1.EnvVar{
			Name: envNvidiaVisibleDevices,
		})
	}

	if len(container.VolumeMounts) == 0 {
		container.VolumeMounts = []corev1.VolumeMount{}
	}
	modelMountPath := getModelMountPath(container, sdep.Name)
	container.VolumeMounts = append(container.VolumeMounts, []corev1.VolumeMount{
		{
			Name:      modelSharedMountName,
			MountPath: modelMountPath,
		},
	}...)

	for idx := range pu.Children {
		err := composeDefaultUserContainer(sdep, &pu.Children[idx], componentSpecMap)
		if err != nil {
			return err
		}
	}
	return nil
}

func composeSeldonPodSpec(pu *seldonv1.PredictiveUnit, componentSpecMap map[string]*seldonv1.SeldonPodSpec) error {
	var seldonPodSpec *seldonv1.SeldonPodSpec
	var ok bool
	if seldonPodSpec, ok = componentSpecMap[pu.Name]; !ok {
		return fmt.Errorf("can't find ComponentSpec for graph %v", pu.Name)
	}

	modelTag, err := getModelTag(pu.ModelURI)
	if err != nil {
		return err
	}
	podName := modelTag
	seldonPodSpec.Metadata.Name = podName

	composeSchedulerName(seldonPodSpec)

	if len(seldonPodSpec.Spec.Volumes) == 0 {
		seldonPodSpec.Spec.Volumes = []corev1.Volume{}
	}
	seldonPodSpec.Spec.Volumes = append(seldonPodSpec.Spec.Volumes, []corev1.Volume{
		{
			Name: modelSharedMountName,
			VolumeSource: corev1.VolumeSource{
				EmptyDir: &corev1.EmptyDirVolumeSource{},
			},
		},
	}...)

	for idx := range pu.Children {
		err := composeSeldonPodSpec(&pu.Children[idx], componentSpecMap)
		if err != nil {
			return err
		}
	}

	return nil
}

// getGPUAmount returns the number of GPUs in the resources.
func getGPUAmount(resource corev1.ResourceRequirements) int64 {
	for k, v := range resource.Limits {
		if strings.Contains(strings.ToLower(k.String()), "gpu") {
			return v.Value()
		}
	}

	return 0
}

// getDefaultUserContainerPorts get container ports for default image.
func getDefaultUserContainerPorts(format string) []corev1.ContainerPort {
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

// getDefaultProbe generate readiness and liveiness.
func getDefaultProbe(format, servingName string) *corev1.Probe {
	path := fmt.Sprintf("/v2/models/%v", servingName)
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
		TimeoutSeconds:   5,
		FailureThreshold: 5,
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
	switch format {
	// Group1 for PMML image
	case string(modeljobsv1alpha1.FormatPMML):
		return viper.GetString(envPMMLServingImage)

	// Group2 for mlserver image
	case string(modeljobsv1alpha1.FormatSKLearn):
		fallthrough
	case string(modeljobsv1alpha1.FormatXGBoost):
		return viper.GetString(envMLServerImage)

	// Group3 for TRT server image
	default:
		return viper.GetString(envTRTServingImage)
	}
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
// For default image, it is /mnt/servingName
// For custom image, it is /<modelStorePath>/servingName, modelStorePath is come from env of container.
func getModelMountPath(container *corev1.Container, servingName string) string {
	for _, env := range container.Env {
		// for custom image, the mount path is pass by frontend.
		if env.Name == envModelStorePath && env.Value != "" {
			return fmt.Sprintf("%v/%v", env.Value, servingName)
		}
	}

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

func composeInitContainer(sdep *seldonv1.SeldonDeployment, pu *seldonv1.PredictorSpec) error {
	p := pu.ComponentSpecs[0]

	if len(p.Spec.Containers) == 0 {
		return fmt.Errorf("there are no container in SeldonPodSpec")
	}
	userContainer := &p.Spec.Containers[0]
	modelMountPath := getModelMountPath(userContainer, sdep.Name)

	initContainer := &corev1.Container{
		// mimics the behavior of seldon model initializer for it will disable the default init container injection
		// in case of a pre-packaged server implementation was selected: https://github.com/SeldonIO/seldon-core/blob/0bd83773228a18e7f376270f4b85cbef69395b8f/operator/controllers/model_initializer_injector.go#L142
		// the default name will be generated here:
		// https://github.com/SeldonIO/seldon-core/blob/0ef45fd234a674fc9b6c8d034cd2e42b4c9ebd05/operator/controllers/model_initializer_injector.go#L118
		Name:  pu.Name + "-" + "model-initializer",
		Image: viper.GetString(envModelInitializerImage),
		Args:  []string{pu.Graph.ModelURI, modelMountPath},
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
				Value: modelStorePath,
			},
		},
		VolumeMounts: []corev1.VolumeMount{
			{
				Name:      modelSharedMountName,
				MountPath: modelMountPath,
			},
		},
	}

	err := composeModelInitailzerContainerResource(initContainer)
	if err != nil {
		return err
	}

	p.Spec.InitContainers = []corev1.Container{*initContainer}

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
