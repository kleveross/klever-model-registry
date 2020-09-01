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
		if modeljob.Spec.DesiredTag == nil {
			return nil, fmt.Errorf("modeljob desired tag is nil")
		}
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

	initContainers, err := generateInitContainers(modeljob)
	if err != nil {
		return nil, err
	}

	backoffLimit := int32(0)
	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: modeljob.Namespace,
			Name:      modeljob.Name,
		},
		Spec: batchv1.JobSpec{
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					InitContainers: initContainers,
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
							ImagePullPolicy: corev1.PullAlways,
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
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "model",
									MountPath: modeljobsv1alpha1.SourceModelPath,
								},
							},
						},
					},
					Volumes: []corev1.Volume{
						{
							Name: "model",
							VolumeSource: corev1.VolumeSource{
								EmptyDir: &corev1.EmptyDirVolumeSource{},
							},
						},
					},
					RestartPolicy: corev1.RestartPolicyNever,
				},
			},
			BackoffLimit: &backoffLimit,
		},
	}

	return job, nil
}

// generateInitContainers will pull model from harbor and export the model to /models/input path
func generateInitContainers(modeljob *modeljobsv1alpha1.ModelJob) ([]corev1.Container, error) {
	ormbDomain := viper.GetString(common.ORMBDomainEnvKey)
	ormbUsername := viper.GetString(common.ORMBUsernameEnvkey)
	ormbPassword := viper.GetString(common.ORMBPasswordEnvKey)
	if ormbDomain == "" || ormbUsername == "" || ormbPassword == "" {
		return nil, nil
	}

	image, ok := PresetAnalyzeImageConfig.Data["ormb-storage-initializer"]
	if !ok {
		return nil, fmt.Errorf("failed get ormb-storage-initializer image")
	}

	initContainers := []corev1.Container{
		{
			Name:  "model-initializer",
			Image: image,
			// Set --relayout=false, only pull and export model, not move any file
			// please refenrence https://github.com/kleveross/ormb/blob/master/cmd/ormb-storage-initializer/cmd/pull-and-export.go
			Args:       []string{modeljob.Spec.Model, modeljobsv1alpha1.SourceModelPath, "--relayout=false"},
			WorkingDir: "/models",
			Env: []corev1.EnvVar{
				corev1.EnvVar{
					Name:  "AWS_ACCESS_KEY_ID",
					Value: ormbUsername,
				},
				corev1.EnvVar{
					Name:  "AWS_SECRET_ACCESS_KEY",
					Value: ormbPassword,
				},
			},
			VolumeMounts: []corev1.VolumeMount{
				{
					Name:      "model",
					MountPath: modeljobsv1alpha1.SourceModelPath,
				},
			},
			ImagePullPolicy: corev1.PullAlways,
		},
	}

	return initContainers, nil
}
