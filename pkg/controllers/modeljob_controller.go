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

package controllers

import (
	"context"

	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	modeljobsv1alpha1 "github.com/caicloud/temp-model-registry/pkg/api/v1alpha1"
)

const (
	// ControllerName is the name of the controller.
	ControllerName = "model-registry-controller"
)

// ModelJobReconciler reconciles a ModelJob object
type ModelJobReconciler struct {
	client.Client
	record.EventRecorder
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=modeljobs.caicloud.io,resources=modeljobs,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=modeljobs.caicloud.io,resources=modeljobs/status,verbs=get;update;patch

func (r *ModelJobReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	_ = context.Background()
	_ = r.Log.WithValues("modeljob", req.NamespacedName)

	// Fetch the ModelJob instance
	instance := &modeljobsv1alpha1.ModelJob{}
	err := r.Get(context.TODO(), req.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Object not found, return.  Created objects are automatically garbage collected.
			// For additional cleanup logic use finalizers.
			return reconcile.Result{}, nil
		}

		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	return r.reconcile(instance)
}

func (r *ModelJobReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&modeljobsv1alpha1.ModelJob{}).
		Owns(&corev1.Pod{}).
		Complete(r)
}
