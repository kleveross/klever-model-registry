package controllers

import (
	"context"
	"fmt"
	"time"

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
		return ctrl.Result{Requeue: true, RequeueAfter: 10 * time.Second}, nil
	}

	// Get a local copy of modeljob's instance.
	oldModelJob := modeljob.DeepCopy()

	result, err := r.reconcileJob(modeljob)
	if err != nil || result.Requeue {
		return result, nil
	}

	// Update modeljob's status.
	if !equality.Semantic.DeepEqual(modeljob.Status, oldModelJob.Status) {
		if err := r.Status().Update(context.Background(), modeljob); err != nil {
			return reconcile.Result{Requeue: true, RequeueAfter: 10 * time.Second}, err
		}
	}

	return reconcile.Result{}, nil
}

func (r *ModelJobReconciler) reconcileJob(modeljob *modeljobsv1alpha1.ModelJob) (ctrl.Result, error) {
	var err error

	job := &batchv1.Job{}
	err = r.Get(context.TODO(), types.NamespacedName{Namespace: modeljob.Namespace, Name: modeljob.Name}, job)
	if err != nil {
		if errors.IsNotFound(err) {
			job, err := generateJobResource(modeljob)
			if err != nil {
				r.recordStatus(modeljob, modeljobsv1alpha1.ModelJobFailed, corev1.EventTypeWarning, ModelJobReasonFailed, "failed to generate job", err)
				return ctrl.Result{}, nil
			}

			if err := controllerutil.SetControllerReference(modeljob, job, r.Scheme); err != nil {
				r.recordStatus(modeljob, modeljobsv1alpha1.ModelJobFailed, corev1.EventTypeWarning, ModelJobReasonFailed, "failed to set job ownreference failed", err)
				return ctrl.Result{}, nil
			}

			if err := r.Create(context.TODO(), job); err != nil {
				if errors.IsAlreadyExists(err) {
					return ctrl.Result{}, nil
				}
				r.recordStatus(modeljob, modeljobsv1alpha1.ModelJobFailed, corev1.EventTypeWarning, ModelJobReasonFailed, "failed to create job", err)
				return ctrl.Result{}, nil
			}

			modeljob.Status.Phase = modeljobsv1alpha1.ModelJobPending
		}
	}

	return r.updateModelJobStatus(job, modeljob)
}

func (r *ModelJobReconciler) updateModelJobStatus(job *batchv1.Job, modeljob *modeljobsv1alpha1.ModelJob) (ctrl.Result, error) {

	if job == nil || modeljob == nil {
		return ctrl.Result{}, nil
	}

	if job.Status.Active != 0 {
		message, err := r.getModelJobMesage(modeljob)
		if err != nil || message != "" {
			r.recordStatus(modeljob, modeljobsv1alpha1.ModelJobPending, corev1.EventTypeWarning, ModelJobReasonPending, message, err)
			return ctrl.Result{}, nil
		}

		r.recordStatus(modeljob, modeljobsv1alpha1.ModelJobRunning, corev1.EventTypeNormal, ModelJobReasonStartRunning, "modelJob start running", nil)
		return ctrl.Result{Requeue: true, RequeueAfter: 15 * time.Second}, nil
	}

	if job.Status.Succeeded != 0 {
		r.recordStatus(modeljob, modeljobsv1alpha1.ModelJobSucceeded, corev1.EventTypeNormal, ModelJobReasonSucceded, "modelJob run successfully", nil)
		return ctrl.Result{}, nil
	}

	if job.Status.Failed != 0 {
		message, err := r.getModelJobMesage(modeljob)
		r.recordStatus(modeljob, modeljobsv1alpha1.ModelJobFailed, corev1.EventTypeWarning, ModelJobReasonFailed, message, err)
		return ctrl.Result{}, nil
	}

	return ctrl.Result{}, nil
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

func (r *ModelJobReconciler) getModelJobMesage(modeljob *modeljobsv1alpha1.ModelJob) (string, error) {
	pods := corev1.PodList{}
	opt := client.ListOptions{
		LabelSelector: labels.SelectorFromSet(labels.Set(map[string]string{"job-name": modeljob.Name})),
		Namespace:     modeljob.Namespace,
	}
	err := r.List(context.TODO(), &pods, &opt)
	if err != nil {
		return "", err
	}

	if len(pods.Items) == 0 {
		return "waiting pod to creat", nil
	}

	// For pod condition.
	for _, c := range pods.Items[0].Status.Conditions {
		if c.Type == corev1.PodScheduled && c.Status == corev1.ConditionFalse {
			return c.Message, nil
		}
	}

	// deal common error
	// 1. pull image error.
	// 2. container is out of memory.
	// 3. container is crashoff.
	dealCommonErr := func(containerStatus *corev1.ContainerStatus) string {
		if containerStatus.LastTerminationState.Waiting != nil {
			switch containerStatus.State.Waiting.Reason {
			case ReasonImagePull, ReasonImagePullBackOff:
				return "failed to pull image"
			case ReasonOOMKilled:
				return "container is out of memory"
			case ReasonContainerCreating:
				return "container is creating"
			}
		}

		if containerStatus.State.Waiting != nil {
			switch containerStatus.State.Waiting.Reason {
			case ReasonImagePull, ReasonImagePullBackOff:
				return "failed to pull image"
			case ReasonCrashLoopBackOff:
				return "container is crash"
			case ReasonContainerCreating:
				return "container is creating"
			}
		}

		return ""
	}

	// For init container.
	if len(pods.Items[0].Status.InitContainerStatuses) == 0 {
		return "init container is creating", nil
	}
	for _, cs := range pods.Items[0].Status.InitContainerStatuses {
		errMsg := dealCommonErr(&cs)
		if errMsg != "" {
			return errMsg, nil
		}

		if cs.State.Terminated != nil && cs.State.Terminated.ExitCode != 0 {
			return "failed to pull model", nil
		}
	}

	// For task container.
	if len(pods.Items[0].Status.ContainerStatuses) == 0 {
		return "container is creating", nil
	}
	for _, cs := range pods.Items[0].Status.ContainerStatuses {
		errMsg := dealCommonErr(&cs)
		if errMsg != "" {
			return errMsg, nil
		}

		if cs.State.Terminated != nil {
			switch cs.State.Terminated.ExitCode {
			case ErrORMBLogin:
				return "failed to login model registry", nil
			case ErrORMBPullModel:
				return "failed to pull model from model registry", nil
			case ErrORMBExportModel:
				return "failed to export model to local", nil
			case ErrRunTask:
				return "failed to run extract/convert task", nil
			case ErrORMBSaveModel:
				return "failed to save model to localhost", nil
			case ErrORMBPushModel:
				return "failed to push model to model registry", nil
			default:
				return fmt.Sprintf("unknow error, err code: %v", cs.State.Terminated.ExitCode), nil
			}
		}
	}
	return "", nil
}
