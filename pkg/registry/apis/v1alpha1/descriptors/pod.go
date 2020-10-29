package descriptors

import (
	"context"

	"github.com/caicloud/nirvana/definition"
	corev1 "k8s.io/api/core/v1"

	"github.com/kleveross/klever-model-registry/pkg/registry/client"
	"github.com/kleveross/klever-model-registry/pkg/registry/pod"
)

var (
	podController *pod.PodController
)

func init() {
	register(podAPI)
}

// InitPodController inits the pod controller
func InitPodController() {
	podController = pod.New(client.GetKubeMainClient())
}

var podAPI = definition.Descriptor{
	Description: "APIs for pod",
	Children: []definition.Descriptor{
		{
			Path:        "/namespaces/{namespace}/pods",
			Definitions: []definition.Definition{getPods},
		},
	},
}

var getPods = definition.Definition{
	Method:      definition.Get,
	Summary:     "Get pods",
	Description: "Get pods",
	Parameters: []definition.Parameter{
		definition.PathParameterFor("namespace", "namespace"),
		definition.QueryParameterFor("filterBy", "filter pod, eg: job-name=\"123\""),
	},
	Results: []definition.Result{
		definition.DataResultFor("pods"),
		definition.ErrorResult(),
	},
	Function: func(ctx context.Context, namespace, filterBy string) (*corev1.PodList, error) {
		return podController.GetPods(namespace, filterBy)
	},
}
