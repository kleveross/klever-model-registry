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
	"fmt"
	"time"

	"github.com/kleveross/ormb/pkg/oras"
	"github.com/kleveross/ormb/pkg/ormb"
	seldonv1 "github.com/seldonio/seldon-core/operator/client/machinelearning.seldon.io/v1/clientset/versioned"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"

	kleverossv1alpha1 "github.com/kleveross/klever-model-registry/pkg/clientset/clientset/versioned"
	kleverossinformers "github.com/kleveross/klever-model-registry/pkg/clientset/informers/externalversions"
	"github.com/kleveross/klever-model-registry/pkg/clientset/informers/externalversions/modeljob/v1alpha1"
)

var (
	// kubeMainClient is the client for k8s builtin resource.
	kubeMainClient kubernetes.Interface

	// kubeKleverOssClient is the client of klever CRD.
	kubeKleverOssClient kleverossv1alpha1.Interface
	// kleverOssModelJobInformer is the informer of kleveross.
	kleverOssModelJobInformer v1alpha1.ModelJobInformer

	// kubeSeldonClient is the client of seldon development.
	kubeSeldonClient seldonv1.Interface

	// ormbClient is interact with harbor.
	ormbClient ormb.Interface
)

func GetKubeMainClient() kubernetes.Interface {
	return kubeMainClient
}

func GetKubeKleverOssClient() kleverossv1alpha1.Interface {
	return kubeKleverOssClient
}

func GetKubeKleverOssModelJobInformer() v1alpha1.ModelJobInformer {
	return kleverOssModelJobInformer
}

func GetKubeSeldonClient() seldonv1.Interface {
	return kubeSeldonClient
}

func GetORMBClient() ormb.Interface {
	return ormbClient
}

// InitClient initializes the client.
func InitClient(kubeconfigPath, domain, username, password string,
	stopCh <-chan struct{}) error {
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err != nil {
		return err
	}

	// init k8s main client
	kubeMainClient, err = kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}

	// init kleveross client
	kubeKleverOssClient, err = kleverossv1alpha1.NewForConfig(config)
	if err != nil {
		return err
	}

	factory := kleverossinformers.NewSharedInformerFactory(kubeKleverOssClient, 30*time.Second)

	go factory.Start(stopCh)

	kleverOssModelJobInformer = factory.Kleveross().V1alpha1().ModelJobs()
	if !cache.WaitForCacheSync(stopCh, kleverOssModelJobInformer.Informer().HasSynced) {
		return fmt.Errorf("failed to wait for modeljob synced")
	}

	// init seldon core client
	kubeSeldonClient, err = seldonv1.NewForConfig(config)
	if err != nil {
		return err
	}

	// init ormb client
	ormbClient, err = ormb.New(oras.ClientOptPlainHTTP(true))
	if err != nil {
		return err
	}
	err = ormbClient.Login(domain, username, password, true)
	if err != nil {
		return err
	}

	return nil
}
