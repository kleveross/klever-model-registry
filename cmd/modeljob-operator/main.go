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

package main

import (
	"flag"
	"os"

	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	modeljobsv1alpha1 "github.com/caicloud/temp-model-registry/pkg/api/v1alpha1"
	"github.com/caicloud/temp-model-registry/pkg/controllers"
	// +kubebuilder:scaffold:imports
)

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))

	utilruntime.Must(modeljobsv1alpha1.AddToScheme(scheme))
	// +kubebuilder:scaffold:scheme

	controllers.ModelFormatToFrameworkMapping = map[modeljobsv1alpha1.Format]modeljobsv1alpha1.Framework{
		modeljobsv1alpha1.FormatSavedModel:  modeljobsv1alpha1.FrameworkTensorflow,
		modeljobsv1alpha1.FormatONNX:        modeljobsv1alpha1.FrameworkONNX,
		modeljobsv1alpha1.FormatH5:          modeljobsv1alpha1.FrameworkKeras,
		modeljobsv1alpha1.FormatPMML:        modeljobsv1alpha1.FrameworkPMML,
		modeljobsv1alpha1.FormatCaffeModel:  modeljobsv1alpha1.FrameworkCaffe,
		modeljobsv1alpha1.FormatNetDef:      modeljobsv1alpha1.FrameworkCaffe2,
		modeljobsv1alpha1.FormatMXNETParams: modeljobsv1alpha1.FrameworkMXNet,
		modeljobsv1alpha1.FormatTouchScript: modeljobsv1alpha1.FrameworkPyTorch,
		modeljobsv1alpha1.FormatGraphDef:    modeljobsv1alpha1.FrameworkTensorflow,
		modeljobsv1alpha1.FormatTensorRT:    modeljobsv1alpha1.FrameworkTensorRT,
	}

}

func main() {
	var metricsAddr string
	var enableLeaderElection bool
	flag.StringVar(&metricsAddr, "metrics-addr", ":8080", "The address the metric endpoint binds to.")
	flag.BoolVar(&enableLeaderElection, "enable-leader-election", false,
		"Enable leader election for controller manager. "+
			"Enabling this will ensure there is only one active controller manager.")
	flag.Parse()

	ctrl.SetLogger(zap.New(zap.UseDevMode(true)))

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:             scheme,
		MetricsBindAddress: metricsAddr,
		Port:               9443,
		LeaderElection:     enableLeaderElection,
		LeaderElectionID:   "9dabac17.caicloud.io",
	})
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	if err = (&controllers.ModelJobReconciler{
		Client:        mgr.GetClient(),
		Log:           ctrl.Log.WithName(controllers.ControllerName).WithName("ModelJob"),
		EventRecorder: mgr.GetEventRecorderFor(controllers.ControllerName),
		Scheme:        mgr.GetScheme(),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "ModelJob")
		os.Exit(1)
	}
	// +kubebuilder:scaffold:builder

	setupLog.Info("starting manager")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		setupLog.Error(err, "problem running manager")
		os.Exit(1)
	}
}
