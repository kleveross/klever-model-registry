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

	"github.com/kleveross/klever-model-registry/pkg/registry/event"
)

func init() {
	register(eventAPI)
}

var eventAPI = definition.Descriptor{
	Description: "APIs for event",
	Children: []definition.Descriptor{
		{
			Path:        "/namespaces/{namespace}/modeljobs/{modeljobID}/events",
			Definitions: []definition.Definition{getModelJobEvents},
		},
	},
}

var getModelJobEvents = definition.Definition{
	Method:      definition.Get,
	Summary:     "Get modeljob events",
	Description: "Get modeljob events",
	Parameters: []definition.Parameter{
		definition.PathParameterFor("namespace", "modeljob namespace"),
		definition.PathParameterFor("modeljobID", "modeljob id"),
	},
	Results: []definition.Result{
		definition.DataResultFor("modeljob events"),
		definition.ErrorResult(),
	},
	Function: func(ctx context.Context, namespace, modeljobID string) (interface{}, error) {
		return event.GetModelJobEvents(namespace, modeljobID)
	},
}
