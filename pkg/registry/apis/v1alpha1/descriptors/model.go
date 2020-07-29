package descriptors

import (
	"context"

	"github.com/caicloud/nirvana/definition"
	"github.com/kleveross/klever-model-registry/pkg/registry/models"
)

func init() {
	register(modelAPI)
}

var modelAPI = definition.Descriptor{
	Description: "APIs for model",
	Children: []definition.Descriptor{
		{
			Path:        "/models/{modelName}/versions/{versionName}/upload",
			Definitions: []definition.Definition{uploadModel},
		},
		{
			Path:        "/models/{modelName}/versions/{versionName}/download",
			Definitions: []definition.Definition{downloadModel},
		},
	},
}

var uploadModel = definition.Definition{
	Method:      definition.Create,
	Summary:     "Upload model",
	Description: "Upload model",
	Parameters: []definition.Parameter{
		definition.HeaderParameterFor("X-Tenant", "Tenant name"),
		definition.HeaderParameterFor("X-User", "User name"),
		definition.PathParameterFor("modelName", "model name"),
		definition.PathParameterFor("versionName", "version name"),
	},
	Results: []definition.Result{
		definition.ErrorResult(),
	},
	Function: func(ctx context.Context, tenant, user, modelName, versionName string) error {
		return models.UploadFile(ctx, tenant, user, modelName, versionName)
	},
}

var downloadModel = definition.Definition{
	Method:      definition.Create,
	Summary:     "Download model",
	Description: "Download model",
	Parameters: []definition.Parameter{
		definition.HeaderParameterFor("X-Tenant", "Tenant name"),
		definition.HeaderParameterFor("X-User", "User name"),
		definition.PathParameterFor("modelName", "model name"),
		definition.PathParameterFor("versionName", "version name"),
		{
			Source:      definition.Body,
			Name:        "model",
			Description: "model body",
		},
	},
	Results: []definition.Result{
		definition.ErrorResult(),
	},
	Function: func(ctx context.Context, tenant, user, modelName, versionName string, model *models.Model) error {
		return models.DownloadFile(ctx, tenant, user, modelName, versionName, model)
	},
}
