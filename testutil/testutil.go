package testutil

import (
	"context"

	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func InitPresetModelImageConfigMap() *corev1.ConfigMap {
	presetImageConfigMap := &corev1.ConfigMap{}
	presetImageConfigMap.Data = map[string]string{
		"savedmodel-extract": "cargo.dev.caicloud.xyz/release/savedmodel:v0.2",
	}

	return presetImageConfigMap
}

func CreateFailedPodForJob(c client.Client, job *batchv1.Job) error {
	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      job.Name,
			Namespace: job.Namespace,
			Labels: map[string]string{
				"job-name": job.Name,
			},
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:  "containet-name",
					Image: "image",
				},
			},
			RestartPolicy: corev1.RestartPolicyNever,
		},
	}
	return c.Create(context.Background(), pod)
}
