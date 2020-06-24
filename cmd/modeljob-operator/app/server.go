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
	"os"

	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/caicloud/temp-model-registry/cmd/modeljob-operator/app/options"
	modeljobsv1alpha1 "github.com/caicloud/temp-model-registry/pkg/api/v1alpha1"
	"github.com/caicloud/temp-model-registry/pkg/controllers"
	"github.com/caicloud/temp-model-registry/pkg/version"
)

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
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

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:             scheme,
		MetricsBindAddress: opt.MetricsAddr,
		Port:               9443,
		LeaderElection:     opt.EnableLeaderElection,
		LeaderElectionID:   "9dabac17.caicloud.io",
	})
	if err != nil {
		setupLog.Error(err, "unable to start manager")
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
