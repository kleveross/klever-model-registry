package descriptors

import (
	"context"

	"github.com/caicloud/nirvana/definition"

	"github.com/kleveross/klever-model-registry/pkg/registry/client"
	"github.com/kleveross/klever-model-registry/pkg/registry/serving"
	seldonv1 "github.com/seldonio/seldon-core/operator/apis/machinelearning.seldon.io/v1"
)

var servingController *serving.ServingController

func init() {
	register(servingAPI)
}

func InitServingController() {
	servingController = serving.New(client.GetKubeSeldonClient())
}

var servingAPI = definition.Descriptor{
	Description: "APIs for serving",
	Children: []definition.Descriptor{
		{
			Path:        "/namespaces/{namespace}/servings",
			Definitions: []definition.Definition{createServing, listServing},
		},
		{
			Path:        "/namespaces/{namespace}/servings/{servingID}",
			Definitions: []definition.Definition{deleteServing, getServing},
		},
	},
}

var createServing = definition.Definition{
	Method:      definition.Create,
	Summary:     "Create serving",
	Description: "Create serving",
	Parameters: []definition.Parameter{
		definition.PathParameterFor("namespace", "namespace"),
		{
			Source:      definition.Body,
			Name:        "serving",
			Description: "serving body",
		},
	},
	Results: []definition.Result{
		definition.ErrorResult(),
	},
	Function: func(ctx context.Context, namespace string, sdep *seldonv1.SeldonDeployment) error {
		return servingController.Create(namespace, sdep)
	},
}

var listServing = definition.Definition{
	Method:      definition.List,
	Summary:     "List definition",
	Description: "List definition",
	Parameters: []definition.Parameter{
		definition.PathParameterFor("namespace", "namespace"),
	},
	Results: definition.DataErrorResults("definition list"),
	Function: func(ctx context.Context, namespace string) (*seldonv1.SeldonDeploymentList, error) {
		return servingController.List(namespace)
	},
}

var getServing = definition.Definition{
	Method:      definition.Get,
	Summary:     "Get definition",
	Description: "Get definition",
	Parameters: []definition.Parameter{
		definition.PathParameterFor("namespace", "namespace"),
		definition.PathParameterFor("servingID", "definition id"),
	},
	Results: definition.DataErrorResults("definition"),
	Function: func(ctx context.Context, namespace, servingID string) (*seldonv1.SeldonDeployment, error) {
		return servingController.Get(namespace, servingID)
	},
}

var deleteServing = definition.Definition{
	Method:      definition.Delete,
	Summary:     "Delete definition",
	Description: "Delete definition",

	Parameters: []definition.Parameter{
		definition.PathParameterFor("namespace", "namespace"),
		definition.PathParameterFor("servingID", "definition id"),
	},
	Results: []definition.Result{
		definition.ErrorResult(),
	},
	Function: func(ctx context.Context, namespace, servingID string) error {
		return servingController.Delete(namespace, servingID)
	},
}
