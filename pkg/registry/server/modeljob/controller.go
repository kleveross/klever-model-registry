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

	modeljobsv1alpha1 "github.com/caicloud/temp-model-registry/pkg/apis/modeljob/v1alpha1"
	"github.com/caicloud/temp-model-registry/pkg/registry/client"
	"github.com/caicloud/temp-model-registry/pkg/registry/common"
	"github.com/caicloud/temp-model-registry/pkg/registry/errors"
	"github.com/caicloud/temp-model-registry/pkg/registry/server/basecontroller"
)

// modeljobController for modeljob APIs impl
type modeljobController struct {
	basecontroller.BaseController
}

func (m *modeljobController) Create() {
	modeljob := &modeljobsv1alpha1.ModelJob{}

	err := json.Unmarshal(m.Ctx.Input.RequestBody, modeljob)
	if err != nil {
		m.RendorError(errors.BadRequestCode, err.Error())
		return
	}

	_, err = client.KubeModelJobClient.ModeljobsV1alpha1().
		ModelJobs(common.KubeSystemNamespace).Create(modeljob)
	if err != nil {
		m.RendorError(errors.GeneralCode, err.Error())
		return
	}

	m.RendorCreateSuccess()
	return
}

func (m *modeljobController) Get() {
	modeljobID := m.Ctx.Input.Param(":modeljob_id")

	modeljob, err := client.KubeModelJobClient.ModeljobsV1alpha1().
		ModelJobs(common.KubeSystemNamespace).Get(modeljobID, metav1.GetOptions{})
	if err != nil {
		m.RendorError(errors.GeneralCode, err.Error())
	}

	m.RendorSuccessData(modeljob)
	return
}

func (m *modeljobController) Delete() {
	modeljobID := m.Ctx.Input.Param(":modeljob_id")

	err := client.KubeModelJobClient.ModeljobsV1alpha1().
		ModelJobs(common.KubeSystemNamespace).Delete(modeljobID, &metav1.DeleteOptions{})
	if err != nil {
		m.RendorError(errors.GeneralCode, err.Error())
		return
	}

	m.RendorDeleteSuccess()
	return
}

func (m *modeljobController) List() {
	modeljobs, err := client.KubeModelJobClient.ModeljobsV1alpha1().ModelJobs(common.KubeSystemNamespace).List(metav1.ListOptions{})
	if err != nil {
		m.RendorError(errors.GeneralCode, err.Error())
		return
	}

	m.RendorSuccessData(modeljobs)
	return
}
