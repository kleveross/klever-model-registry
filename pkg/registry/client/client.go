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
	"github.com/spf13/viper"
	"k8s.io/client-go/tools/clientcmd"

	clientset "github.com/caicloud/temp-model-registry/pkg/clientset/clientset/versioned"
)

var (
	KubeCRDClient *clientset.Clientset
)

func init() {
	kubeconfigPath := viper.GetString("kubeconfig")

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err != nil {
		panic(err)
	}

	KubeCRDClient, err = clientset.NewForConfig(config)
	if err != nil {
		panic(err)
	}
}
