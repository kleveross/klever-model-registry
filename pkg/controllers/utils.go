package controllers

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	modeljobsv1alpha1 "github.com/kleveross/klever-model-registry/pkg/apis/modeljob/v1alpha1"
	"github.com/kleveross/klever-model-registry/pkg/common"
)

func getFrameworkByFormat(format modeljobsv1alpha1.Format) modeljobsv1alpha1.Framework {
	return ModelFormatToFrameworkMapping[format]
}

func generateJobResource(modeljob *modeljobsv1alpha1.ModelJob) (*batchv1.Job, error) {
	var dstFormat modeljobsv1alpha1.Format
	var dstFramework modeljobsv1alpha1.Framework
	dstTag := ""
	image := ""
	if modeljob.Spec.Conversion != nil {
		dstTag = *modeljob.Spec.DesiredTag
		dstFormat = modeljob.Spec.Conversion.MMdnn.To
		dstFramework = getFrameworkByFormat(dstFormat)
		image = PresetAnalyzeImageConfig.Data[strings.ToLower(string(modeljob.Spec.Conversion.MMdnn.From))+"-convert"]
	} else if modeljob.Spec.Extraction != nil {
		dstTag = "empty"
		dstFormat = modeljob.Spec.Extraction.Format
		dstFramework = getFrameworkByFormat(dstFormat)
		image = PresetAnalyzeImageConfig.Data[strings.ToLower(string(dstFormat))+"-extract"]
	} else {
		return nil, fmt.Errorf("%v", "not support source")
	}

	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: modeljob.Namespace,
			Name:      modeljob.Name,
		},
		Spec: batchv1.JobSpec{
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:       "executor",
							Image:      image,
							WorkingDir: "/models",
							Command:    []string{"sh"},
							Args: []string{
								"-c",
								fmt.Sprintf("/scripts/run.sh"),
							},
							Env: []corev1.EnvVar{
								corev1.EnvVar{
									Name:  modeljobsv1alpha1.FrameworkEnvKey,
									Value: string(dstFramework),
								},
								corev1.EnvVar{
									Name:  modeljobsv1alpha1.FormatEnvKey,
									Value: string(dstFormat),
								},
								corev1.EnvVar{
									Name:  modeljobsv1alpha1.SourceModelTagEnvKey,
									Value: modeljob.Spec.Model,
								},
								corev1.EnvVar{
									Name:  modeljobsv1alpha1.DestinationModelTagEnvKey,
									Value: dstTag,
								},
								corev1.EnvVar{
									Name:  modeljobsv1alpha1.SourceModelPathEnvKey,
									Value: modeljobsv1alpha1.SourceModelPath,
								},
								corev1.EnvVar{
									Name:  modeljobsv1alpha1.DestinationModelPathEnvKey,
									Value: modeljobsv1alpha1.DestinationModelPath,
								},
								corev1.EnvVar{
									Name:  modeljobsv1alpha1.ExtractorEnvKey,
									Value: strings.ToLower(string(dstFormat)),
								},
								corev1.EnvVar{
									Name:  common.ORMBDomainEnvKey,
									Value: viper.GetString(common.ORMBDomainEnvKey),
								},
								corev1.EnvVar{
									Name:  common.ORMBUsernameEnvkey,
									Value: viper.GetString(common.ORMBUsernameEnvkey),
								},
								corev1.EnvVar{
									Name:  common.ORMBPasswordEnvKey,
									Value: viper.GetString(common.ORMBPasswordEnvKey),
								},
							},
						},
					},
					RestartPolicy: corev1.RestartPolicyNever,
				},
			},
		},
	}

	return job, nil
}
