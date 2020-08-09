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

	"github.com/kleveross/klever-model-registry/pkg/registry/client"
	"github.com/kleveross/klever-model-registry/pkg/registry/log"
	"github.com/kleveross/klever-model-registry/pkg/registry/resource/logs"
)

var (
	logController *log.LogController
)

func init() {
	register(logAPI)
}

func InitLogController() {
	logController = log.New(client.GetKubeMainClient())
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
		definition.Parameter{
			Source:      definition.Query,
			Name:        "containerID",
			Description: "container id",
			Optional:    true,
		},
		definition.Parameter{
			Source:      definition.Query,
			Name:        "reftimestamp",
			Description: "ref timestamp",
			Optional:    true,
		},
		definition.Parameter{
			Source:      definition.Query,
			Name:        "reflineNum",
			Description: "ref line num",
			Optional:    true,
		},
		definition.Parameter{
			Source:      definition.Query,
			Name:        "userPreviousLogs",
			Description: "user previous logs",
			Optional:    true,
		},
		definition.Parameter{
			Source:      definition.Query,
			Name:        "offsetFrom",
			Description: "offset from",
			Optional:    true,
		},
		definition.Parameter{
			Source:      definition.Query,
			Name:        "offsetTo",
			Description: "offset to",
			Optional:    true,
		},
		definition.Parameter{
			Source:      definition.Query,
			Name:        "logFilePosition",
			Description: "log file position",
			Optional:    true,
		},
	},
	Results: []definition.Result{
		definition.DataResultFor("pod logs"),
		definition.ErrorResult(),
	},
	Function: func(ctx context.Context, namespace, podID, containerID, refTimestamp string,
		refLineNum int, usePreviousLogs bool, offsetFrom,
		offsetTo, logFilePosition string) (*logs.LogDetails, error) {
		return logController.GetPodLogs(namespace, podID, containerID, refTimestamp,
			refLineNum, usePreviousLogs, offsetFrom, offsetTo, logFilePosition)
	},
}
