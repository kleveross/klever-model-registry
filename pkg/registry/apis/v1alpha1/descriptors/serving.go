package descriptors

import (
	"context"

	"github.com/caicloud/nirvana/definition"

	"github.com/kleveross/klever-model-registry/pkg/registry/client"
	"github.com/kleveross/klever-model-registry/pkg/registry/paging"
	"github.com/kleveross/klever-model-registry/pkg/registry/serving"
	seldonv1 "github.com/seldonio/seldon-core/operator/apis/machinelearning.seldon.io/v1"
)

var servingController *serving.ServingController

func init() {
	register(servingAPI)
}

// InitServingController inits the seldon serving controller
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
		definition.BodyParameterFor("serving body"),
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
	Summary:     "List Serving",
	Description: "List Serving",
	Parameters: []definition.Parameter{
		definition.PathParameterFor("namespace", "namespace"),
		paging.PageDefinitionParameter(),
	},
	Results: definition.DataErrorResults("definition list"),
	Function: func(ctx context.Context, namespace string, opt *paging.ListOption) (*serving.ServingList, error) {
		return servingController.List(namespace, opt)
	},
}

var updateServing = definition.Definition{
	Method:      definition.Update,
	Summary:     "Updates serving",
	Description: "Updates serving",
	Parameters: []definition.Parameter{
		definition.PathParameterFor("namespace", "namespace"),
		definition.PathParameterFor("servingID", "serving id"),
		definition.BodyParameterFor("serving body"),
	},
	Results: definition.DataErrorResults("serving updated"),
	Function: func(ctx context.Context, namespace string, servingID string, body []byte) (*seldonv1.SeldonDeployment, error) {
		return servingController.Update(namespace, servingID, body)
	},
}

var getServing = definition.Definition{
	Method:      definition.Get,
	Summary:     "Get Serving",
	Description: "Get Serving",
	Parameters: []definition.Parameter{
		definition.PathParameterFor("namespace", "namespace"),
		definition.PathParameterFor("servingID", "serving id"),
	},
	Results: definition.DataErrorResults("serving"),
	Function: func(ctx context.Context, namespace, servingID string) (*seldonv1.SeldonDeployment, error) {
		return servingController.Get(namespace, servingID)
	},
}

var deleteServing = definition.Definition{
	Method:      definition.Delete,
	Summary:     "Delete Serving",
	Description: "Delete Serving",
	Parameters: []definition.Parameter{
		definition.PathParameterFor("namespace", "namespace"),
		definition.PathParameterFor("servingID", "serving id"),
	},
	Results: []definition.Result{
		definition.ErrorResult(),
	},
	Function: func(ctx context.Context, namespace, servingID string) error {
		return servingController.Delete(namespace, servingID)
	},
}
