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

	"github.com/kleveross/klever-model-registry/pkg/registry/client"
	"github.com/kleveross/klever-model-registry/pkg/registry/errors"
	"github.com/kleveross/klever-model-registry/pkg/registry/resource/container"
	"github.com/kleveross/klever-model-registry/pkg/registry/resource/logs"
	"github.com/kleveross/klever-model-registry/pkg/registry/server/basecontroller"
)

// LogController for log APIs implement
type LogController struct {
	basecontroller.BaseController
}

// Get is get pod log, now only get one container log.
// If need, we can support containerID param in the furture.
func (l *LogController) Get() {
	namespace := l.Ctx.Input.Param(":namespace")
	podID := l.Ctx.Input.Param(":pod_id")
	containerID := l.Ctx.Input.Query("container")

	refTimestamp := l.Ctx.Input.Query("reference_timestamp")
	if refTimestamp == "" {
		refTimestamp = logs.NewestTimestamp
	}

	refLineNum, err := strconv.Atoi(l.Ctx.Input.Query("reference_line_num"))
	if err != nil {
		refLineNum = 0
	}
	usePreviousLogs := l.Ctx.Input.Query("previous") == "true"
	offsetFrom, err1 := strconv.Atoi(l.Ctx.Input.Query("offset_from"))
	offsetTo, err2 := strconv.Atoi(l.Ctx.Input.Query("offset_to"))
	logFilePosition := l.Ctx.Input.Query("log_file_position")

	logSelector := logs.DefaultSelection
	if err1 == nil && err2 == nil {
		logSelector = &logs.Selection{
			ReferencePoint: logs.LogLineId{
				LogTimestamp: logs.LogTimestamp(refTimestamp),
				LineNum:      refLineNum,
			},
			OffsetFrom:      offsetFrom,
			OffsetTo:        offsetTo,
			LogFilePosition: logFilePosition,
		}
	}

	result, err := container.GetLogDetails(client.KubeMainClient, namespace, podID, containerID, logSelector, usePreviousLogs)
	if err != nil {
		l.RendorError(errors.GeneralCode, err.Error())
		return
	}

	l.RendorSuccessData(result)
}
