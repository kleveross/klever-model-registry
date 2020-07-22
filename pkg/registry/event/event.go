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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/reference"

	modeljobscheme "github.com/kleveross/klever-model-registry/pkg/clientset/clientset/versioned/scheme"
	"github.com/kleveross/klever-model-registry/pkg/registry/client"
)

func GetModelJobEvents(namespace, modeljobID string) (*corev1.EventList, error) {
	modeljob, err := client.KubeModelJobClient.KleverossV1alpha1().
		ModelJobs(namespace).Get(modeljobID, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	var events *corev1.EventList
	if ref, err := reference.GetReference(scheme.Scheme, modeljob); err != nil {
		log.Errorf("get modeljob reference err: %v", err)
	} else {
		events, err = client.KubeMainClient.CoreV1().Events(namespace).Search(modeljobscheme.Scheme, ref)
		if err != nil {
			log.Errorf("search modeljob event err: %v", err)
		}
	}

	return events, nil
}
