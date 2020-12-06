package controllers

import (
	corev1 "k8s.io/api/core/v1"

	modeljobsv1alpha1 "github.com/kleveross/klever-model-registry/pkg/apis/modeljob/v1alpha1"
)

var (
	// ModelFormatToFrameworkMapping is the map for model's format to model's framework.
	ModelFormatToFrameworkMapping map[modeljobsv1alpha1.Format]modeljobsv1alpha1.Framework
	// PresetAnalyzeImageConfig is the preset image of analyze for some model's format.
	PresetAnalyzeImageConfig *corev1.ConfigMap

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
