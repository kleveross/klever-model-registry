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
)
