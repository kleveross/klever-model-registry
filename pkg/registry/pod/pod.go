package pod

import (
	"context"
	"strings"

	"github.com/caicloud/nirvana/log"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type PodController struct {
	kubeMainClient kubernetes.Interface
}

func New(kubeMainClient kubernetes.Interface) *PodController {
	return &PodController{
		kubeMainClient: kubeMainClient,
	}
}

func (e PodController) GetPods(namespace string, filterBy string) (*corev1.PodList, error) {
	seletor := metav1.ListOptions{}
	if len(filterBy) > 0 {
		pairs := strings.Split(filterBy, ",")
		for _, item := range pairs {
			data := strings.Split(item, ":")
			seletor.LabelSelector += strings.Join(data, "=")
		}
	}

	pods, err := e.kubeMainClient.CoreV1().Pods(namespace).List(context.TODO(), seletor)
	if err != nil {
		log.Errorf("failed to list pods, err: %v", err)
		return nil, err
	}

	return pods, nil
}
