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

	modeljobsv1alpha1 "github.com/kleveross/klever-model-registry/pkg/apis/modeljob/v1alpha1"
	"github.com/kleveross/klever-model-registry/pkg/common"
	"github.com/kleveross/klever-model-registry/pkg/registry/client"
)

func Create(modeljob *modeljobsv1alpha1.ModelJob) (*modeljobsv1alpha1.ModelJob, error) {
	err := ExchangeModelJobNameAndID(&modeljob.ObjectMeta)
	if err != nil {
		return nil, err
	}

	result, err := client.KubeModelJobClient.ModeljobsV1alpha1().
		ModelJobs(common.DefaultModelJobNamespace).Create(modeljob)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func Get(modeljobID string) (*modeljobsv1alpha1.ModelJob, error) {
	modeljob, err := client.KubeModelJobClient.ModeljobsV1alpha1().
		ModelJobs(common.DefaultModelJobNamespace).Get(modeljobID, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	return modeljob, err
}

func Delete(modeljobID string) error {
	err := client.KubeModelJobClient.ModeljobsV1alpha1().
		ModelJobs(common.DefaultModelJobNamespace).Delete(modeljobID, &metav1.DeleteOptions{})
	if err != nil {
		return err
	}

	return nil
}

func List() (*modeljobsv1alpha1.ModelJobList, error) {
	modeljobs, err := client.KubeModelJobClient.ModeljobsV1alpha1().
		ModelJobs(common.DefaultModelJobNamespace).List(metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	return modeljobs, nil
}
