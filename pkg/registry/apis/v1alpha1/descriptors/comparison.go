package descriptors

import (
	"context"

	"github.com/caicloud/nirvana/definition"

	"github.com/kleveross/klever-model-registry/pkg/registry/comparison"
	"github.com/kleveross/klever-model-registry/pkg/registry/paging"
)

func init() {
	register(comparisonAPI)
}

var comparisonAPI = definition.Descriptor{
	Description: "APIs for model comparison",
	Children: []definition.Descriptor{
		{
			Path:        "/comparisons",
			Definitions: []definition.Definition{generateComparison},
		},
		{
			Path:        "/comparativedocument",
			Definitions: []definition.Definition{downloadComparison},
		},
	},
}

var generateComparison = definition.Definition{
	Method:      definition.List,
	Description: "Generate Comparison",
	Parameters: []definition.Parameter{
		definition.BodyParameterFor("Comparison Body"),
		paging.PageDefinitionParameter(),
	},
	Results: definition.DataErrorResults("generate comparison"),
	Function: func(ctx context.Context, models comparison.Comparison, opt *paging.ListOption) (*comparison.ORMBModelList, error) {
		return comparison.Generator(ctx, models, opt)
	},
}

var downloadComparison = definition.Definition{
	Method:      definition.Get,
	Description: "Download Comparison",
	Results:     []definition.Result{definition.ErrorResult()},
	Parameters: []definition.Parameter{
		definition.BodyParameterFor("Comparison Body"),
	},
	Function: func(ctx context.Context, models comparison.Comparison) error {
		return comparison.DownloadCSVFile(ctx, models)
	},
}
