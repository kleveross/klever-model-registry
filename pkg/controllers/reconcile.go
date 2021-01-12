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
			return ctrl.Result{Requeue: true, RequeueAfter: 15 * time.Second}, nil
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

	return getModelJobMesageByPods(&pods)
}

func getModelJobMesageByPods(pods *corev1.PodList) (string, error) {
	if len(pods.Items) == 0 {
		return errContainerCreating, nil
	}

	// For pod condition.
	for _, c := range pods.Items[0].Status.Conditions {
		if c.Type == corev1.PodScheduled && c.Status == corev1.ConditionFalse {
			return c.Message, nil
		}
	}

	// deal common error
	dealCommonErr := func(containerStatus *corev1.ContainerStatus) string {
		if containerStatus.LastTerminationState.Waiting != nil {
			errMsg := convertPodReasonToMessage(containerStatus.LastTerminationState.Waiting.Reason)
			if errMsg != "" {
				return errMsg
			}
		}

		if containerStatus.State.Waiting != nil {
			errMsg := convertPodReasonToMessage(containerStatus.State.Waiting.Reason)
			if errMsg != "" {
				return errMsg
			}
		}

		return ""
	}

	// For init container.
	if len(pods.Items[0].Status.InitContainerStatuses) == 0 {
		return errContainerCreating, nil
	}
	for _, cs := range pods.Items[0].Status.InitContainerStatuses {
		errMsg := dealCommonErr(&cs)
		if errMsg != "" {
			return errMsg, nil
		}

		if cs.State.Terminated != nil && cs.State.Terminated.ExitCode != 0 {
			return errORMBPull, nil
		}
	}

	// For task container.
	if len(pods.Items[0].Status.ContainerStatuses) == 0 {
		return errContainerCreating, nil
	}
	for _, cs := range pods.Items[0].Status.ContainerStatuses {
		errMsg := dealCommonErr(&cs)
		if errMsg != "" {
			return errMsg, nil
		}

		if cs.State.Terminated != nil {
			switch cs.State.Terminated.ExitCode {
			case ErrORMBLogin:
				return errORMBLogin, nil
			case ErrORMBPullModel:
				return errORMBPull, nil
			case ErrORMBExportModel:
				return errORMBExport, nil
			case ErrRunTask:
				return errRunTask, nil
			case ErrORMBSaveModel:
				return errORMBSave, nil
			case ErrORMBPushModel:
				return errORMBPush, nil
			default:
				return fmt.Sprintf("unknow error, err code: %v", cs.State.Terminated.ExitCode), nil
			}
		}
	}
	return "", nil
}

func convertPodReasonToMessage(reason string) string {
	switch reason {
	case ReasonImagePull, ReasonImagePullBackOff:
		return errPullImage
	case ReasonOOMKilled:
		return errContainerOutOfMemory
	case ReasonContainerCreating:
		return errContainerCreating
	case ReasonCrashLoopBackOff:
		return errContainerCrashed
	case ReasonInvalidImageName:
		return errContainerImageInvalid
	}

	return ""
}
