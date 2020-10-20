package modeljob

import (
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	modeljobsv1alpha1 "github.com/kleveross/klever-model-registry/pkg/apis/modeljob/v1alpha1"
	"github.com/kleveross/klever-model-registry/pkg/util"
)

// GenerateExtractionModelJob will generate ModelJob by base information.
func GenerateExtractionModelJob(domain, project, modelName, versionName, format string) *modeljobsv1alpha1.ModelJob {
	modeljob := modeljobsv1alpha1.ModelJob{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "kleveross.io/v1alpha1",
			Kind:       "ModelJob",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      util.RandomNameWithPrefix(fmt.Sprintf("modeljob-%v-%v-%v", project, modelName, versionName)),
			Namespace: "default",
		},
		Spec: modeljobsv1alpha1.ModelJobSpec{
			Model: fmt.Sprintf("%v/%v/%v:%v", domain, project, modelName, versionName),
			ModelJobSource: modeljobsv1alpha1.ModelJobSource{
				Extraction: &modeljobsv1alpha1.ExtractionSource{
					Format: modeljobsv1alpha1.Format(format),
				},
			},
		},
	}

	return &modeljob
}

var formatNeedNotExtraction = map[string]bool{
	string(modeljobsv1alpha1.FormatTensorRT): true,
	string(modeljobsv1alpha1.FormatMLflow):   true,
}

// IsExtractModel return bool which is represent whether extract model or not.
// For `TensorRT` format, since extract MUST have GPU, but GPU resource is precious, so not extract.
// For `Others` format, not extract.
func IsExtractModel(format string) bool {
	if _, ok := formatNeedNotExtraction[format]; ok {
		return false
	}

	return true
}
