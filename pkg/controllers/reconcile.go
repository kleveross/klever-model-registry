package controllers

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	modeljobsv1alpha1 "github.com/caicloud/temp-model-registry/pkg/api/v1alpha1"
)

func (r *ModelJobReconciler) reconcile(modeljob *modeljobsv1alpha1.ModelJob) (ctrl.Result, error) {

	// Update state if is deleting
	if !modeljob.ObjectMeta.DeletionTimestamp.IsZero() && modeljob.Status.Phase != modeljobsv1alpha1.ModelJobDeleting {
		modeljob.Status.Phase = modeljobsv1alpha1.ModelJobDeleting
		err := r.Status().Update(context.Background(), modeljob)
		if err != nil {
			return ctrl.Result{}, err
		}
		return ctrl.Result{Requeue: true}, nil
	}

	// Get a local copy of modeljob's instance.
	oldModelJob := modeljob.DeepCopy()

	r.reconcileWorkerPod(modeljob)

	// Update modeljob's status.
	if !equality.Semantic.DeepEqual(modeljob.Status, oldModelJob.Status) {
		if err := r.Status().Update(context.Background(), modeljob); err != nil {
			return reconcile.Result{}, err
		}
	}

	return reconcile.Result{}, nil
}

func (r *ModelJobReconciler) reconcileWorkerPod(modeljob *modeljobsv1alpha1.ModelJob) error {
	var err error

	pod := &corev1.Pod{}
	err = r.Get(context.TODO(), types.NamespacedName{Namespace: modeljob.Namespace, Name: getPodName(modeljob.Name)}, pod)
	if err != nil {
		if errors.IsNotFound(err) {
			pod, err := generatePod(modeljob)
			if err != nil {
				modeljob.Status.Phase = modeljobsv1alpha1.ModelJobFailed
				r.Log.Error(err, "New pod failed")
				r.Event(modeljob, "Error", "Failed", fmt.Sprintf("New pod failed"))
				return nil
			}

			if err := controllerutil.SetControllerReference(modeljob, pod, r.Scheme); err != nil {
				modeljob.Status.Phase = modeljobsv1alpha1.ModelJobFailed
				r.Log.Error(err, "Set pod ownreference failed")
				r.Event(modeljob, "Error", "Failed", fmt.Sprintf("Set pod ownreference failed"))
				return nil
			}

			if err := r.Create(context.TODO(), pod); err != nil {
				if errors.IsAlreadyExists(err) {
					return nil
				}
				modeljob.Status.Phase = modeljobsv1alpha1.ModelJobFailed
				r.Log.Error(err, "Create pod failed")
				r.Event(modeljob, "Error", "Failed", fmt.Sprintf("Create pod failed"))
			}

			modeljob.Status.Phase = modeljobsv1alpha1.ModelJobPending
		}
	}

	if pod == nil {
		r.Log.Info("pod is nil")
		return nil
	}
	switch pod.Status.Phase {
	case corev1.PodRunning:
		if modeljob.Status.Phase == modeljobsv1alpha1.ModelJobPending {
			modeljob.Status.Phase = modeljobsv1alpha1.ModelJobRunning
		}
	case corev1.PodSucceeded:
		modeljob.Status.Phase = modeljobsv1alpha1.ModelJobSucceeded
	case corev1.PodFailed, corev1.PodUnknown:
		modeljob.Status.Phase = modeljobsv1alpha1.ModelJobFailed
		modeljob.Status.Message = pod.Status.Message

	}
	return nil
}
