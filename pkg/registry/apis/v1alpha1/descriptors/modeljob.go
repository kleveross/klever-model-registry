/*
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package descriptors

import (
	"context"

	"github.com/caicloud/nirvana/definition"

	modeljobsv1alpha1 "github.com/kleveross/klever-model-registry/pkg/apis/modeljob/v1alpha1"
	"github.com/kleveross/klever-model-registry/pkg/registry/modeljob"
)

func init() {
	register(modeljobAPI)
}

var modeljobAPI = definition.Descriptor{
	Description: "APIs for modeljob",
	Children: []definition.Descriptor{
		{
			Path:        "/modeljobs",
			Definitions: []definition.Definition{createModelJob, listModelJob},
		},
		{
			Path:        "/modeljobs/{modeljobID}",
			Definitions: []definition.Definition{deleteModelJob, getModelJob},
		},
	},
}

var createModelJob = definition.Definition{
	Method:      definition.Create,
	Summary:     "Create modeljob",
	Description: "Create modeljob",
	Parameters: []definition.Parameter{
		{
			Source:      definition.Body,
			Name:        "job",
			Description: "modeljob body",
		},
	},
	Results: []definition.Result{
		definition.ErrorResult(),
	},
	Function: func(ctx context.Context, job *modeljobsv1alpha1.ModelJob) error {
		_, err := modeljob.Create(job)
		return err
	},
}

var listModelJob = definition.Definition{
	Method:      definition.List,
	Summary:     "List modeljob",
	Description: "List modeljob",
	Parameters:  []definition.Parameter{},
	Results:     definition.DataErrorResults("modeljob list"),
	Function: func(ctx context.Context) (*modeljobsv1alpha1.ModelJobList, error) {
		return modeljob.List()
	},
}

var getModelJob = definition.Definition{
	Method:      definition.Get,
	Summary:     "Get modeljob",
	Description: "Get modeljob",
	Parameters: []definition.Parameter{
		definition.PathParameterFor("modeljobID", "modeljob id"),
	},
	Results: definition.DataErrorResults("modeljob"),
	Function: func(ctx context.Context, modeljobID string) (*modeljobsv1alpha1.ModelJob, error) {
		return modeljob.Get(modeljobID)
	},
}

var deleteModelJob = definition.Definition{
	Method:      definition.Delete,
	Summary:     "Delete modeljob",
	Description: "Delete modeljob",

	Parameters: []definition.Parameter{
		definition.PathParameterFor("modeljobID", "modeljob id"),
	},
	Results: []definition.Result{
		definition.ErrorResult(),
	},
	Function: func(ctx context.Context, modeljobID string) error {
		return modeljob.Delete(modeljobID)
	},
}
