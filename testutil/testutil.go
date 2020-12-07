package testutil

import (
	"context"
	"os"

	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func InitPresetModelImage() {
	os.Setenv("SAVEDMODEL_EXTRACT_IMAGE", "demo.goharbor.com/release/savedmodel-extract:v0.2.0")
	os.Setenv("H5_CONVERSION_IMAGE", "demo.goharbor.com/release/h5_to_savedmodel:v0.2.0")
	os.Setenv("ORMB_INITIALIZER_IMAGE", "demo.goharbor.com/release/klever-ormb-storage-initializer:v0.0.8")
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
