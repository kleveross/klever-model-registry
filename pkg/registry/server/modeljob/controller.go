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
package modeljob

import (
	"encoding/json"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	modeljobsv1alpha1 "github.com/kleveross/klever-model-registry/pkg/apis/modeljob/v1alpha1"
	"github.com/kleveross/klever-model-registry/pkg/common"
	"github.com/kleveross/klever-model-registry/pkg/registry/client"
	"github.com/kleveross/klever-model-registry/pkg/registry/errors"
	"github.com/kleveross/klever-model-registry/pkg/registry/server/basecontroller"
)

// ModeljobController for modeljob APIs impl
type ModeljobController struct {
	basecontroller.BaseController
}

func (m *ModeljobController) Create() {
	modeljob := &modeljobsv1alpha1.ModelJob{}

	err := json.Unmarshal(m.Ctx.Input.RequestBody, modeljob)
	if err != nil {
		m.RendorError(errors.BadRequestCode, err.Error())
		return
	}

	err = ExchangeModelJobNameAndID(&modeljob.ObjectMeta)
	if err != nil {
		m.RendorError(errors.GeneralCode, err.Error())
		return
	}

	_, err = client.KubeModelJobClient.ModeljobsV1alpha1().
		ModelJobs(common.DefaultModelJobNamespace).Create(modeljob)
	if err != nil {
		m.RendorError(errors.GeneralCode, err.Error())
		return
	}

	m.RendorCreateSuccess()
	return
}

func (m *ModeljobController) Get() {
	modeljobID := m.Ctx.Input.Param(":modeljob_id")

	modeljob, err := client.KubeModelJobClient.ModeljobsV1alpha1().
		ModelJobs(common.DefaultModelJobNamespace).Get(modeljobID, metav1.GetOptions{})
	if err != nil {
		m.RendorError(errors.GeneralCode, err.Error())
	}

	m.RendorSuccessData(modeljob)
	return
}

func (m *ModeljobController) Delete() {
	modeljobID := m.Ctx.Input.Param(":modeljob_id")

	err := client.KubeModelJobClient.ModeljobsV1alpha1().
		ModelJobs(common.DefaultModelJobNamespace).Delete(modeljobID, &metav1.DeleteOptions{})
	if err != nil {
		m.RendorError(errors.GeneralCode, err.Error())
		return
	}

	m.RendorDeleteSuccess()
	return
}

func (m *ModeljobController) List() {
	modeljobs, err := client.KubeModelJobClient.ModeljobsV1alpha1().
		ModelJobs(common.DefaultModelJobNamespace).List(metav1.ListOptions{})
	if err != nil {
		m.RendorError(errors.GeneralCode, err.Error())
		return
	}

	m.RendorSuccessData(modeljobs)
	return
}
