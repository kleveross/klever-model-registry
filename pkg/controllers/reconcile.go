package controllers

import (
	"context"
	"fmt"

	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
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
				r.recordStatus(modeljob, modeljobsv1alpha1.ModelJobFailed, corev1.EventTypeWarning, ModelJobReasonFailed, "failed to generate job", err)
				return nil
			}

			if err := controllerutil.SetControllerReference(modeljob, job, r.Scheme); err != nil {
				r.recordStatus(modeljob, modeljobsv1alpha1.ModelJobFailed, corev1.EventTypeWarning, ModelJobReasonFailed, "failed to set job ownreference failed", err)
				return nil
			}

			if err := r.Create(context.TODO(), job); err != nil {
				if errors.IsAlreadyExists(err) {
					return nil
				}
				r.recordStatus(modeljob, modeljobsv1alpha1.ModelJobFailed, corev1.EventTypeWarning, ModelJobReasonFailed, "failed to create job", err)
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
		r.recordStatus(modeljob, modeljobsv1alpha1.ModelJobPending, corev1.EventTypeNormal, ModelJobReasonPending, "modeljob pending", nil)
		return
	}

	if job.Status.Active != 0 {
		r.recordStatus(modeljob, modeljobsv1alpha1.ModelJobRunning, corev1.EventTypeNormal, ModelJobReasonStartRunning, "modelJob start running", nil)
		return
	}

	if job.Status.Succeeded != 0 {
		r.recordStatus(modeljob, modeljobsv1alpha1.ModelJobSucceeded, corev1.EventTypeNormal, ModelJobReasonSucceded, "modelJob run successfully", nil)
		return
	}

	if job.Status.Failed != 0 {
		modeljob.Status.Phase = modeljobsv1alpha1.ModelJobFailed

		pods := corev1.PodList{}
		opt := client.ListOptions{
			LabelSelector: labels.SelectorFromSet(labels.Set(map[string]string{"job-name": modeljob.Name})),
			Namespace:     job.Namespace,
		}
		err := r.List(context.TODO(), &pods, &opt)
		if err != nil {
			r.recordStatus(modeljob, "", corev1.EventTypeWarning, "", "failed to get pod for modeljob", err)
			return
		}

		modeljob.Status.Message = getModelJobFailedMesage(&pods)
		return
	}

	return
}

func (r *ModelJobReconciler) recordStatus(modeljob *modeljobsv1alpha1.ModelJob, phase modeljobsv1alpha1.ModelJobPhase,
	eventType, reason, message string, err error) {
	if phase != "" {
		modeljob.Status.Phase = phase

	}
	modeljob.Status.Message = message

	if err != nil {
		message := fmt.Sprintf("%v, err: %v", message, err)
		r.Log.Error(err, message, "modelJobName", modeljob.Name)
	} else if eventType == corev1.EventTypeWarning {
		r.Log.Info(message, "modelJobName", modeljob.Name)
	}

	r.Event(modeljob, eventType, reason, message)
}

func getModelJobFailedMesage(pods *corev1.PodList) string {
	if len(pods.Items) == 0 {
		return ""
	}

	cs := pods.Items[0].Status.InitContainerStatuses
	if cs != nil && cs[0].State.Terminated != nil {
		return "run model-initializer error"
	}

	cs = pods.Items[0].Status.ContainerStatuses
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
