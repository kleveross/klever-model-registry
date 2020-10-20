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

// getORMBDomain is get domain for ModelJob task.
// For model extraction, when it is completed, it should push model to harbor directly
// so that not create ModelJob repeatedly, but for model conversion, when convert complete,
// it should push to klever-model-registry, so that it can extract model automatically.
func getORMBDomain(isConvert bool) string {
	ormbDomain := viper.GetString(common.ORMBDomainEnvKey)
	if isConvert {
		ormbDomain = viper.GetString(KleverModelRegistryAddressEnvKey)
	}

	return ormbDomain
}

// replaceModelRefDomain will replace modelRef domain.
// The real domain is depend on getORMBDomain.
func replaceModelRefDomain(inputModelRef, ormbDomain string) (string, error) {
	refSlice := strings.Split(inputModelRef, "/")

	modelRef := ""
	if len(refSlice) == 2 {
		// The form like release/savedmodel:v1.0, do not have default domain.
		modelRef = strings.Join([]string{ormbDomain, refSlice[0], refSlice[1]}, "/")
	} else if len(refSlice) == 3 {
		// The form like harbor.io/release/savedmodel:v1.0, it have default domain.
		modelRef = strings.Join([]string{ormbDomain, refSlice[1], refSlice[2]}, "/")
	} else {
		return "", fmt.Errorf("The model ref's format is error")
	}

	return modelRef, nil
}

func generateJobResource(modeljob *modeljobsv1alpha1.ModelJob) (*batchv1.Job, error) {
	var dstFormat modeljobsv1alpha1.Format
	var dstFramework modeljobsv1alpha1.Framework
	var image string
	var srcModelRef string
	var dstModelRef string
	var ormbDomain string
	var err error

	if modeljob.Spec.Conversion != nil {
		if modeljob.Spec.DesiredTag == nil {
			return nil, fmt.Errorf("modeljob desired tag is nil")
		}
		ormbDomain = getORMBDomain(true)
		dstModelRef, err = replaceModelRefDomain(*modeljob.Spec.DesiredTag, ormbDomain)
		if err != nil {
			return nil, err
		}
		dstFormat = modeljob.Spec.Conversion.MMdnn.To
		dstFramework = getFrameworkByFormat(dstFormat)
		image = PresetAnalyzeImageConfig.Data[strings.ToLower(string(modeljob.Spec.Conversion.MMdnn.From))+"-convert"]
	} else if modeljob.Spec.Extraction != nil {
		ormbDomain = getORMBDomain(false)
		dstModelRef = "empty"
		dstFormat = modeljob.Spec.Extraction.Format
		dstFramework = getFrameworkByFormat(dstFormat)
		image = PresetAnalyzeImageConfig.Data[strings.ToLower(string(dstFormat))+"-extract"]
	} else {
		return nil, fmt.Errorf("%v", "not support source")
	}

	srcModelRef, err = replaceModelRefDomain(modeljob.Spec.Model, ormbDomain)
	if err != nil {
		return nil, err
	}

	initContainers, err := generateInitContainers(modeljob)
	if err != nil {
		return nil, err
	}

	schedulerName := getSchedulerName()
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
							ImagePullPolicy: corev1.PullIfNotPresent,
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
									Value: srcModelRef,
								},
								corev1.EnvVar{
									Name:  modeljobsv1alpha1.DestinationModelTagEnvKey,
									Value: dstModelRef,
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
									Value: ormbDomain,
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
					SchedulerName: schedulerName,
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

func getSchedulerName() string {
	schedulerName := viper.GetString(SchedulerNameEnvKey)
	if schedulerName == "" {
		schedulerName = DefaultSchedulerName
	}

	return schedulerName
}
