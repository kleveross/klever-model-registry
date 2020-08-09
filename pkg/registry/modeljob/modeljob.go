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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"

	modeljobsv1alpha1 "github.com/kleveross/klever-model-registry/pkg/apis/modeljob/v1alpha1"
	clientset "github.com/kleveross/klever-model-registry/pkg/clientset/clientset/versioned"
	"github.com/kleveross/klever-model-registry/pkg/clientset/informers/externalversions/modeljob/v1alpha1"
)

type ModelJobController struct {
	kleverossClient  clientset.Interface
	modeljobInformer v1alpha1.ModelJobInformer
}

func New(kleverossClient clientset.Interface, modeljobInformer v1alpha1.ModelJobInformer) *ModelJobController {
	return &ModelJobController{
		kleverossClient:  kleverossClient,
		modeljobInformer: modeljobInformer,
	}
}

func (m ModelJobController) Create(namespace string, modeljob *modeljobsv1alpha1.ModelJob) (*modeljobsv1alpha1.ModelJob, error) {
	err := ExchangeModelJobNameAndID(&modeljob.ObjectMeta)
	if err != nil {
		return nil, err
	}

	result, err := m.kleverossClient.KleverossV1alpha1().ModelJobs(namespace).Create(modeljob)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (m ModelJobController) Get(namespace, modeljobID string) (*modeljobsv1alpha1.ModelJob, error) {
	modeljob, err := m.modeljobInformer.Lister().ModelJobs(namespace).Get(modeljobID)
	if err != nil {
		return nil, err
	}

	return modeljob, err
}

func (m ModelJobController) Delete(namespace, modeljobID string) error {
	err := m.kleverossClient.KleverossV1alpha1().ModelJobs(namespace).Delete(modeljobID, &metav1.DeleteOptions{})
	if err != nil {
		return err
	}

	return nil
}

func (m ModelJobController) List(namespace string) ([]*modeljobsv1alpha1.ModelJob, error) {
	modeljobs, err := m.modeljobInformer.Lister().ModelJobs(namespace).List(labels.Everything())
	if err != nil {
		return nil, err
	}

	return modeljobs, nil
}
