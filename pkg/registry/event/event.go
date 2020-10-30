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
	"github.com/caicloud/nirvana/log"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/reference"

	kleverossv1alpha1 "github.com/kleveross/klever-model-registry/pkg/clientset/clientset/versioned"
	modeljobscheme "github.com/kleveross/klever-model-registry/pkg/clientset/clientset/versioned/scheme"
	"github.com/kleveross/klever-model-registry/pkg/clientset/informers/externalversions/modeljob/v1alpha1"
	"github.com/kleveross/klever-model-registry/pkg/registry/errors"
	"github.com/kleveross/klever-model-registry/pkg/registry/paging"
)

type EventController struct {
	kubeMainClient kubernetes.Interface

	kleverossClient  kleverossv1alpha1.Interface
	modeljobInformer v1alpha1.ModelJobInformer
}

func New(kubeMainClient kubernetes.Interface, kleverossClient kleverossv1alpha1.Interface,
	modeljobInformer v1alpha1.ModelJobInformer) *EventController {
	return &EventController{
		kubeMainClient:   kubeMainClient,
		kleverossClient:  kleverossClient,
		modeljobInformer: modeljobInformer,
	}
}

// EventList is the response of List Interface.
type EventList struct {
	ListMeta paging.ListMeta `json:"metadata"`
	Items    []corev1.Event  `json:"items"`
}

// toEventList is convert to ModelJobList struct.
func toEventList(items []corev1.Event, opt *paging.ListOption) *EventList {
	datas := paging.Page(items, opt)
	eventList := &EventList{
		ListMeta: paging.ListMeta{
			TotalItems: datas.TotalItems,
		},
		Items: []corev1.Event{},
	}

	for _, d := range datas.Items {
		eventList.Items = append(eventList.Items, d.(corev1.Event))
	}
	return eventList
}

func (e EventController) GetModelJobEvents(namespace, modeljobID string, opt *paging.ListOption) (*EventList, error) {
	modeljob, err := e.modeljobInformer.Lister().ModelJobs(namespace).Get(modeljobID)
	if err != nil {
		return nil, errors.RenderError(err)
	}

	var events *corev1.EventList
	if ref, err := reference.GetReference(scheme.Scheme, modeljob); err != nil {
		log.Errorf("failed to get modeljob reference, err: %v", err)
	} else {
		events, err = e.kubeMainClient.CoreV1().Events(namespace).Search(modeljobscheme.Scheme, ref)
		if err != nil {
			log.Errorf("failed to search modeljob event, err: %v", err)
		}
	}

	if events != nil {
		return toEventList(events.Items, opt), nil
	}

	return &EventList{}, nil
}
