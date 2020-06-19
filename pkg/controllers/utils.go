package controllers

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	modeljobsv1alpha1 "github.com/caicloud/temp-model-registry/pkg/api/v1alpha1"
)

func init() {
	viper.AutomaticEnv()
}

func getPodName(modeljobName string) string {
	return fmt.Sprintf("modeljob-%v", modeljobName)
}

func getFrameworkByFormat(format modeljobsv1alpha1.Format) modeljobsv1alpha1.Framework {
	return ModelFormatToFrameworkMapping[format]
}

func generatePod(modeljob *modeljobsv1alpha1.ModelJob) (*corev1.Pod, error) {
	dstFramework := getFrameworkByFormat(modeljobsv1alpha1.Format(modeljob.Spec.Conversion.MMdnn.From))
	dstFormat := modeljob.Spec.Conversion.MMdnn.From
	dstTag := ""
	if modeljob.Spec.DesiredTag != nil {
		dstTag = *modeljob.Spec.DesiredTag
		dstFramework = getFrameworkByFormat(modeljobsv1alpha1.Format(modeljob.Spec.Conversion.MMdnn.To))
		dstFormat = modeljob.Spec.Conversion.MMdnn.To
	} else {
		dstTag = "empty"
	}

	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: modeljob.Namespace,
			Name:      getPodName(modeljob.Name),
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name: "executor",
					// TODO(simon): set different image for different job.
					Image:      "harbor.caicloud.com/test/extract:v0.2",
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
							Name:  modeljobsv1alpha1.ORMBDomainEnvKey,
							Value: viper.GetString(modeljobsv1alpha1.ORMBDomainEnvKey),
						},
						corev1.EnvVar{
							Name:  modeljobsv1alpha1.ORMBUsernameEnvkey,
							Value: viper.GetString(modeljobsv1alpha1.ORMBUsernameEnvkey),
						},
						corev1.EnvVar{
							Name:  modeljobsv1alpha1.ORMBPasswordEnvKey,
							Value: viper.GetString(modeljobsv1alpha1.ORMBPasswordEnvKey),
						},
					},
				},
			},
			RestartPolicy: corev1.RestartPolicyNever,
		},
	}
	return pod, nil
}
