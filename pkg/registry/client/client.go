// +build !test
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

package client

import (
	seldonv1 "github.com/seldonio/seldon-core/operator/client/machinelearning.seldon.io/v1/clientset/versioned"
	"github.com/spf13/viper"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	modeljobv1alpha1 "github.com/kleveross/klever-model-registry/pkg/clientset/clientset/versioned"
)

var (
	KubeMainClient     kubernetes.Interface
	KubeModelJobClient modeljobv1alpha1.Interface
	KubelSeldonClient  seldonv1.Interface
)

func InitClient() error {
	kubeconfigPath := viper.GetString("kubeconfig")

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err != nil {
		return err
	}

	KubeMainClient, err = kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}

	KubeModelJobClient, err = modeljobv1alpha1.NewForConfig(config)
	if err != nil {
		return err
	}

	KubelSeldonClient, err = seldonv1.NewForConfig(config)
	if err != nil {
		return err
	}

	return nil
}
