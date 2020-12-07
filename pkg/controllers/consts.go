package controllers

import (
	modeljobsv1alpha1 "github.com/kleveross/klever-model-registry/pkg/apis/modeljob/v1alpha1"
)

var (
	// ModelFormatToFrameworkMapping is the map for model's format to model's framework.
	ModelFormatToFrameworkMapping map[modeljobsv1alpha1.Format]modeljobsv1alpha1.Framework

	// KleverModelRegistryAddressEnvKey is the address of klever-model-registry
	KleverModelRegistryAddressEnvKey = "KLEVER_MODEL_REGISTRY_ADDRESS"

	// DefaultSchedulerName is default scheduler name
	DefaultSchedulerName = "default-scheduler"

	// SchedulerNameEnvKey is the env key for scheduler name, it is set in Deployment
	SchedulerNameEnvKey = "SCHEDULER_NAME"

	// ModelInitializerCPUEnvKey defines the cpu env key for model initializer container.
	ModelInitializerCPUEnvKey = "MODEL_INITIALIZER_CPU"
	// ModelInitializerMEMEnvKey defines the mem env key for model initializer container.
	ModelInitializerMEMEnvKey = "MODEL_INITIALIZER_MEM"
	// ModelJobTaskCPUEnvKey defines the cpu env key for modeljob task container.
	ModelJobTaskCPUEnvKey = "MODELJOB_TASK_CPU"
	// ModelJobTaskMEMEnvKey defines the mem env key for modeljob task container.
	ModelJobTaskMEMEnvKey = "MODELJOB_TASK_MEM"
)

var presetImage = map[string]string{
	"caffemodel-extract":  "CAFFE_EXTRACT_IMAGE",
	"netdef-extract":      "NETDEF_EXTRACT_IMAGE",
	"graphdef-extract":    "GRAPHDEF_EXTRACT_IMAGE",
	"h5-extract":          "H5_EXTRACT_IMAGE",
	"mxnetparams-extract": "MXNETPARAMS_EXTRACT_IMAGE",
	"onnx-extract":        "ONNX_EXTRACT_IMAGE",
	"savedmodel-extract":  "SAVEDMODEL_EXTRACT_IMAGE",
	"torchscript-extract": "TORCHSCRIPT_EXTRACT_IMAGE",
	"pmml-extract":        "PMML_EXTRACT_IMAGE",
	"caffemodel-convert":  "CAFFE_CONVERSION_IMAGE",
	"mxnetparams-convert": "MXNET_CONVERSION_IMAGE",
	"h5-convert":          "H5_CONVERSION_IMAGE",
	"netdef-convert":      "NETDEF_CONVERSION_IMAGE",
	"initializer":         "ORMB_INITIALIZER_IMAGE",
}
