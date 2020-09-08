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
	"sort"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"

	modeljobsv1alpha1 "github.com/kleveross/klever-model-registry/pkg/apis/modeljob/v1alpha1"
	clientset "github.com/kleveross/klever-model-registry/pkg/clientset/clientset/versioned"
	"github.com/kleveross/klever-model-registry/pkg/clientset/informers/externalversions/modeljob/v1alpha1"
	"github.com/kleveross/klever-model-registry/pkg/registry/errors"
	"github.com/kleveross/klever-model-registry/pkg/registry/paging"
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
		return nil, errors.RenderError(err)
	}

	result, err := m.kleverossClient.KleverossV1alpha1().ModelJobs(namespace).Create(modeljob)
	if err != nil {
		return nil, errors.RenderError(err)
	}

	return result, nil
}

func (m ModelJobController) Get(namespace, modeljobID string) (*modeljobsv1alpha1.ModelJob, error) {
	modeljob, err := m.modeljobInformer.Lister().ModelJobs(namespace).Get(modeljobID)
	if err != nil {
		return nil, errors.RenderError(err)
	}

	return modeljob, nil
}

func (m ModelJobController) Delete(namespace, modeljobID string) error {
	err := m.kleverossClient.KleverossV1alpha1().ModelJobs(namespace).Delete(modeljobID, &metav1.DeleteOptions{})
	if err != nil {
		return errors.RenderError(err)
	}

	return nil
}

// ModelJobList is the response of List Interface.
type ModelJobList struct {
	ListMeta paging.ListMeta               `json:"metadata"`
	Items    []*modeljobsv1alpha1.ModelJob `json:"items"`
}

// toModelJobList is convert to ModelJobList struct.
func toModelJobList(items []*modeljobsv1alpha1.ModelJob, opt *paging.ListOption) *ModelJobList {
	datas := paging.Page(items, opt)
	servingList := &ModelJobList{
		ListMeta: paging.ListMeta{
			TotalItems: datas.TotalItems,
		},
		Items: []*modeljobsv1alpha1.ModelJob{},
	}

	for _, d := range datas.Items {
		servingList.Items = append(servingList.Items, d.(*modeljobsv1alpha1.ModelJob))
	}
	return servingList
}

func (m ModelJobController) List(namespace string, opt *paging.ListOption) (*ModelJobList, error) {
	modeljobs, err := m.modeljobInformer.Lister().ModelJobs(namespace).List(labels.Everything())
	if err != nil {
		return nil, errors.RenderError(err)
	}

	sort.SliceStable(modeljobs, func(i, j int) bool {
		return modeljobs[i].Name < modeljobs[j].Name
	})
	return toModelJobList(modeljobs, opt), nil
}
