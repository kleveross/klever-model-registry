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
package event

import (
	"fmt"

	log "github.com/astaxie/beego/logs"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/reference"

	modeljobscheme "github.com/caicloud/temp-model-registry/pkg/clientset/clientset/versioned/scheme"
	"github.com/caicloud/temp-model-registry/pkg/registry/client"
	"github.com/caicloud/temp-model-registry/pkg/registry/errors"
	"github.com/caicloud/temp-model-registry/pkg/registry/server/basecontroller"
)

// EventController for event APIs
type EventController struct {
	basecontroller.BaseController
}

func (m *EventController) GetModelJobEvents() {
	namespace := m.Ctx.Input.Param(":namespace")
	modeljobID := m.Ctx.Input.Param(":modeljob_id")

	modeljob, err := client.KubeModelJobClient.ModeljobsV1alpha1().ModelJobs(namespace).Get(modeljobID, metav1.GetOptions{})
	if err != nil {
		m.RendorError(errors.GeneralCode, err.Error())
		return
	}

	var events *corev1.EventList
	if ref, err := reference.GetReference(scheme.Scheme, modeljob); err != nil {
		log.Error(fmt.Errorf("get modeljob reference err: %v", err))
	} else {
		events, err = client.KubeMainClient.CoreV1().Events(namespace).Search(modeljobscheme.Scheme, ref)
		if err != nil {
			log.Error(fmt.Errorf("search modeljob event err: %v", err))
		}
	}

	m.RendorSuccessData(events)
	return
}
