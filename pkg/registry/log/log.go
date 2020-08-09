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
package log

import (
	"strconv"

	"github.com/kleveross/klever-model-registry/pkg/registry/errors"
	"github.com/kleveross/klever-model-registry/pkg/registry/resource/container"
	"github.com/kleveross/klever-model-registry/pkg/registry/resource/logs"
	"k8s.io/client-go/kubernetes"
)

type LogController struct {
	kubeMainClient kubernetes.Interface
}

func New(kubeMainClient kubernetes.Interface) *LogController {
	return &LogController{
		kubeMainClient: kubeMainClient,
	}
}

// GetPodLogs is get pod log, now only get one container log.
// If need, we can support containerID param in the furture.
func (l LogController) GetPodLogs(namespace, podID, containerID, refTimestamp string, refLineNum int,
	usePreviousLogs bool, offsetFrom, offsetTo, logFilePosition string) (*logs.LogDetails, error) {
	if refTimestamp == "" {
		refTimestamp = logs.NewestTimestamp
	}

	offsetStart, err1 := strconv.Atoi(offsetFrom)
	offsetEnd, err2 := strconv.Atoi(offsetTo)

	logSelector := logs.DefaultSelection
	if err1 == nil && err2 == nil {
		logSelector = &logs.Selection{
			ReferencePoint: logs.LogLineId{
				LogTimestamp: logs.LogTimestamp(refTimestamp),
				LineNum:      refLineNum,
			},
			OffsetFrom:      offsetStart,
			OffsetTo:        offsetEnd,
			LogFilePosition: logFilePosition,
		}
	}

	logs, err := container.GetLogDetails(l.kubeMainClient, namespace, podID,
		containerID, logSelector, usePreviousLogs)
	if err != nil {
		return nil, errors.RenderError(err)
	}

	return logs, nil
}
