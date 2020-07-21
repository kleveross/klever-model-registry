package controllers

import (
	"context"
	"fmt"

	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	modeljobsv1alpha1 "github.com/kleveross/klever-model-registry/pkg/apis/modeljob/v1alpha1"
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

	err := r.reconcileJob(modeljob)
	if err != nil {
		return ctrl.Result{Requeue: true}, nil
	}

	// Update modeljob's status.
	if !equality.Semantic.DeepEqual(modeljob.Status, oldModelJob.Status) {
		if err := r.Status().Update(context.Background(), modeljob); err != nil {
			return reconcile.Result{}, err
		}
	}

	return reconcile.Result{}, nil
}

func (r *ModelJobReconciler) reconcileJob(modeljob *modeljobsv1alpha1.ModelJob) error {
	var err error

	job := &batchv1.Job{}
	err = r.Get(context.TODO(), types.NamespacedName{Namespace: modeljob.Namespace, Name: modeljob.Name}, job)
	if err != nil {
		if errors.IsNotFound(err) {
			job, err := generateJobResource(modeljob)
			if err != nil {
				modeljob.Status.Phase = modeljobsv1alpha1.ModelJobFailed
				r.Log.Error(err, "New job failed")
				r.Event(modeljob, "Error", "Failed", fmt.Sprintf("New jod failed"))
				return nil
			}

			if err := controllerutil.SetControllerReference(modeljob, job, r.Scheme); err != nil {
				modeljob.Status.Phase = modeljobsv1alpha1.ModelJobFailed
				r.Log.Error(err, "Set job ownreference failed")
				r.Event(modeljob, "Error", "Failed", fmt.Sprintf("Set job ownreference failed"))
				return nil
			}

			if err := r.Create(context.TODO(), job); err != nil {
				if errors.IsAlreadyExists(err) {
					return nil
				}
				modeljob.Status.Phase = modeljobsv1alpha1.ModelJobFailed
				r.Log.Error(err, "Create job failed")
				r.Event(modeljob, "Error", "Failed", fmt.Sprintf("Create job failed"))
				return nil
			}

			modeljob.Status.Phase = modeljobsv1alpha1.ModelJobPending
		}
	}

	r.updateModelJobStatus(job, modeljob)

	return nil
}

func (r *ModelJobReconciler) updateModelJobStatus(job *batchv1.Job, modeljob *modeljobsv1alpha1.ModelJob) {

	if job == nil || modeljob == nil {
		return
	}

	if job.Status.StartTime == nil {
		modeljob.Status.Phase = modeljobsv1alpha1.ModelJobPending
		return
	}

	if job.Status.Active != 0 {
		modeljob.Status.Phase = modeljobsv1alpha1.ModelJobRunning
		return
	}

	if job.Status.Succeeded != 0 {
		modeljob.Status.Phase = modeljobsv1alpha1.ModelJobSucceeded
		return
	}

	if job.Status.Failed != 0 {
		modeljob.Status.Phase = modeljobsv1alpha1.ModelJobFailed

		pods := corev1.PodList{}
		err := r.List(context.TODO(), &pods,
			&client.MatchingLabels{"job-name": job.Name},
			&client.MatchingFields{"metadata.namespace": job.ObjectMeta.Namespace})
		if err != nil {
			r.Log.Error(err, fmt.Sprintf("Get pod for modeljob %v", modeljob.Name))
			return
		}

		modeljob.Status.Message = getModelJobFailedMesage(&pods)
		return
	}

	return
}

func getModelJobFailedMesage(pods *corev1.PodList) string {
	if len(pods.Items) == 0 {
		return ""
	}

	cs := pods.Items[0].Status.ContainerStatuses
	if cs != nil && cs[0].State.Terminated != nil {
		switch cs[0].State.Terminated.ExitCode {
		case ErrORMBLogin:
			return "ormb login error"
		case ErrORMBPullModel:
			return "ormb pull model error"
		case ErrORMBExportModel:
			return "ormb export model error"
		case ErrRunTask:
			return "run task error"
		case ErrORMBSaveModel:
			return "ormb save model error"
		case ErrORMBPushModel:
			return "ormb push model error"
		default:
			return fmt.Sprintf("unknow error, error code: %v", cs[0].State.Terminated.ExitCode)
		}
	}
	return ""
}
