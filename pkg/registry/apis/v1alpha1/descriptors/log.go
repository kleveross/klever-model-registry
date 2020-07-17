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

	"github.com/kleveross/klever-model-registry/pkg/registry/log"
)

func init() {
	register(logAPI)
}

var logAPI = definition.Descriptor{
	Description: "APIs for log",
	Children: []definition.Descriptor{
		{
			Path:        "/namespaces/{namespace}/pods/{podID}/log",
			Definitions: []definition.Definition{getPodLogs},
		},
	},
}

var getPodLogs = definition.Definition{
	Method:      definition.Get,
	Summary:     "Get pod logs",
	Description: "Get pod logs",
	Parameters: []definition.Parameter{
		definition.PathParameterFor("namespace", "modeljob namespace"),
		definition.PathParameterFor("podID", "pod id"),
		definition.QueryParameterFor("containerid", "container id"),
		definition.QueryParameterFor("reftimestamp", "ref timestamp"),
		definition.QueryParameterFor("reflineNum", "ref line num"),
		definition.QueryParameterFor("userPreviousLogs", "us previous logs"),
		definition.QueryParameterFor("offsetFrom", "offset from"),
		definition.QueryParameterFor("offsetTo", "offset to"),
		definition.QueryParameterFor("logFilePosition", "log file position"),
	},
	Results: []definition.Result{
		definition.DataResultFor("pod logs"),
		definition.ErrorResult(),
	},
	Function: func(ctx context.Context, namespace, podID, containerID, refTimestamp string,
		refLineNum int, usePreviousLogs bool, offsetFrom,
		offsetTo, logFilePosition string) (interface{}, error) {
		return log.GetPodLogs(namespace, podID, containerID, refTimestamp,
			refLineNum, usePreviousLogs, offsetFrom, offsetTo, logFilePosition)
	},
}
