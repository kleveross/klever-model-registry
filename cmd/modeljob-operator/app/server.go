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
package app

import (
	"bufio"
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	apiextensionsv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	"k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/sets"
	yamlutil "k8s.io/apimachinery/pkg/util/yaml"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"k8s.io/client-go/rest"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/yaml"

	"github.com/kleveross/klever-model-registry/cmd/modeljob-operator/app/options"
	modeljobsv1alpha1 "github.com/kleveross/klever-model-registry/pkg/apis/modeljob/v1alpha1"
	"github.com/kleveross/klever-model-registry/pkg/controllers"
	"github.com/kleveross/klever-model-registry/pkg/version"
)

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

const (
	// CustomResourceDefinitionKind is the kind of resource customresourcedefinition.
	CustomResourceDefinitionKind = "CustomResourceDefinition"
	// CustomResourceDefinitionPath is the path of crd yaml.
	CustomResourceDefinitionPath = "crds"
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))

	utilruntime.Must(modeljobsv1alpha1.AddToScheme(scheme))
	// +kubebuilder:scaffold:scheme
}

func Run(opt *options.ServerOption) error {
	if opt.PrintVersion {
		version.PrintVersion()
		os.Exit(1)
	}

	config := ctrl.GetConfigOrDie()
	mgr, err := ctrl.NewManager(config, ctrl.Options{
		Scheme:             scheme,
		MetricsBindAddress: opt.MetricsAddr,
		Port:               9443,
		LeaderElection:     opt.EnableLeaderElection,
		LeaderElectionID:   "9dabac17.kleveross.io",
	})
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		return err
	}

	if err := ensureCRD(config); err != nil {
		setupLog.Error(err, "unable to start install modeljob crd")
		return err
	}

	if err := controllers.Initialization(); err != nil {
		setupLog.Error(err, "init error")
		return err
	}

	if err = (&controllers.ModelJobReconciler{
		Client:        mgr.GetClient(),
		Log:           ctrl.Log.WithName(controllers.ControllerName).WithName("ModelJob"),
		EventRecorder: mgr.GetEventRecorderFor(controllers.ControllerName),
		Scheme:        mgr.GetScheme(),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "ModelJob")
		return err
	}
	// +kubebuilder:scaffold:builder

	setupLog.Info("starting manager")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		setupLog.Error(err, "problem running manager")
		return err
	}

	return nil
}

func renderCRDs(crdpath string) ([]*apiextensionsv1beta1.CustomResourceDefinition, error) {
	var (
		info  os.FileInfo
		files []os.FileInfo
		dir   string
	)

	info, err := os.Stat(crdpath)
	if err != nil {
		return nil, err
	}

	if !info.IsDir() {
		dir, files = filepath.Dir(crdpath), []os.FileInfo{info}
	} else {
		if files, err = ioutil.ReadDir(crdpath); err != nil {
			return nil, err
		}
		dir = crdpath
	}

	return readCRDs(dir, files)
}

// readCRDs reads the CRDs from files and Unmarshals them into structs
func readCRDs(basePath string, files []os.FileInfo) ([]*apiextensionsv1beta1.CustomResourceDefinition, error) {
	crds := make([]*apiextensionsv1beta1.CustomResourceDefinition, 0, len(files))
	crdExts := sets.NewString(".yaml", ".yml")
	for _, file := range files {
		// Only parse allowed file types
		if !crdExts.Has(filepath.Ext(file.Name())) {
			continue
		}

		// Unmarshal CRDs from file into structs
		docs, err := readDocuments(filepath.Join(basePath, file.Name()))
		if err != nil {
			return nil, err
		}

		for _, doc := range docs {
			crd := &apiextensionsv1beta1.CustomResourceDefinition{}
			if err = yaml.Unmarshal(doc, crd); err != nil {
				return nil, err
			}

			// Check that it is actually a CRD
			if crd.Kind != CustomResourceDefinitionKind || crd.Spec.Names.Kind == "" || crd.Spec.Group == "" {
				continue
			}
			crds = append(crds, crd)
		}
	}
	return crds, nil
}

// readDocuments reads documents from file
func readDocuments(fp string) ([][]byte, error) {
	b, err := ioutil.ReadFile(fp)
	if err != nil {
		return nil, err
	}

	// A file may contain multiple documents separated by `---`, eg:
	//   a: 1
	//   ---
	//   a: 2
	// will be converted to: ["a: 1", "a: 2"]
	docs := make([][]byte, 0, 1)
	reader := yamlutil.NewYAMLReader(bufio.NewReader(bytes.NewReader(b)))
	for {
		doc, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}

			return nil, err
		}
		docs = append(docs, doc)
	}
	return docs, nil
}

func ensureCRD(config *rest.Config) error {
	client, err := clientset.NewForConfig(config)
	if err != nil {
		return err
	}
	crds, err := renderCRDs(CustomResourceDefinitionPath)
	if err != nil {
		return err
	}
	for _, crd := range crds {
		_, err = client.ApiextensionsV1beta1().CustomResourceDefinitions().Create(crd)
		if err != nil {
			setupLog.Error(err, "unable to ensure crds")
			if !errors.IsAlreadyExists(err) {
				return err
			}
		}
	}

	return nil
}
