package controllers

import (
	"context"
	"fmt"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	utilrand "k8s.io/apimachinery/pkg/util/rand"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	modeljobsv1alpha1 "github.com/kleveross/klever-model-registry/pkg/apis/modeljob/v1alpha1"
	test "github.com/kleveross/klever-model-registry/testutil"
)

var _ = Describe("Modeljob Controller", func() {
	initGlobalVar()
	PresetAnalyzeImageConfig = test.InitPresetModelImageConfigMap()
	modeljobName := "test-modeljob"

	req := ctrl.Request{
		NamespacedName: types.NamespacedName{
			Name:      modeljobName,
			Namespace: "default",
		},
	}

	const timeout = time.Second * 5
	const interval = time.Second * 1

	BeforeEach(func() {
		modeljobName = fmt.Sprintf("modeljob-test-%s", utilrand.String(8))
		// failed test runs that don't clean up leave resources behind.
		modeljob := &modeljobsv1alpha1.ModelJob{}
		err := k8sClient.DeleteAllOf(context.Background(), modeljob, client.InNamespace("default"))
		Expect(err).ToNot(HaveOccurred())

		job := &batchv1.Job{}
		err = k8sClient.DeleteAllOf(context.Background(), job, client.InNamespace("default"))
		Expect(err).ToNot(HaveOccurred())
	})

	AfterEach(func() {
		// Add any teardown steps that needs to be executed after each test
	})

	Context("ModelJob", func() {
		It("Should be created correctly", func() {
			key := types.NamespacedName{
				Name:      modeljobName,
				Namespace: "default",
			}

			toCreate := &modeljobsv1alpha1.ModelJob{
				ObjectMeta: metav1.ObjectMeta{
					Name:      key.Name,
					Namespace: key.Namespace,
				},
				Spec: modeljobsv1alpha1.ModelJobSpec{
					Model: "release/savedmodel:v1",
					ModelJobSource: modeljobsv1alpha1.ModelJobSource{
						Extraction: &modeljobsv1alpha1.ExtractionSource{
							Format: modeljobsv1alpha1.FormatSavedModel,
						},
					},
				},
			}
			Expect(k8sClient.Create(context.Background(), toCreate)).Should(Succeed())

			By("Expecting to reconcile a job successfully, the result is get job successfully")
			Eventually(func() error {
				getedModelJob := &modeljobsv1alpha1.ModelJob{}
				err := k8sClient.Get(context.Background(), key, getedModelJob)
				if err != nil {
					return err
				}

				reconciler.Reconcile(req)
				job := &batchv1.Job{}
				return k8sClient.Get(context.Background(), key, job)
			}, timeout, interval).Should(Succeed())

			By("Expecting to reconcile a job successfully, the modejob status update to active")
			Eventually(func() error {
				job := &batchv1.Job{}
				err := k8sClient.Get(context.Background(), key, job)
				Expect(err).To(BeNil())

				job.Status.StartTime = &metav1.Time{
					Time: time.Now(),
				}
				job.Status.Active = 1
				err = k8sClient.Status().Update(context.Background(), job)
				if err != nil {
					return err
				}

				_, err = reconciler.Reconcile(req)
				if err != nil {
					return err
				}

				job = &batchv1.Job{}
				return k8sClient.Get(context.Background(), key, job)
			}, timeout, interval).Should(Succeed())

			By("Expecting to reconcile a job successfully, the modejob status update to success")
			Eventually(func() error {
				job := &batchv1.Job{}
				err := k8sClient.Get(context.Background(), key, job)
				Expect(err).To(BeNil())

				job.Status.StartTime = &metav1.Time{
					Time: time.Now(),
				}
				job.Status.Active = 0
				job.Status.Succeeded = 1
				err = k8sClient.Status().Update(context.Background(), job)
				if err != nil {
					return err
				}

				_, err = reconciler.Reconcile(req)
				if err != nil {
					return err
				}

				job = &batchv1.Job{}
				return k8sClient.Get(context.Background(), key, job)
			}, timeout, interval).Should(Succeed())

			By("Create pod for modejob -> job -> pod")
			job := &batchv1.Job{}
			err := k8sClient.Get(context.Background(), key, job)
			Expect(err).To(BeNil())
			err = test.CreateFailedPodForJob(k8sClient, job)
			Expect(err).To(BeNil())

			By("Update job for modejob -> job failed")
			job.Status.Active = 0
			job.Status.Succeeded = 0
			job.Status.Failed = 1
			k8sClient.Status().Update(context.Background(), job)
			Expect(err).To(BeNil())

			By("Expecting to reconcile a job successfully, the modejob status update to failed, failed reason is ormb login error")
			Eventually(func() error {
				pod := &corev1.Pod{}
				err := k8sClient.Get(context.Background(), key, pod)
				if err != nil {
					return err
				}

				pod.Status.ContainerStatuses = []corev1.ContainerStatus{
					{
						Name: "status",
						State: corev1.ContainerState{
							Terminated: &corev1.ContainerStateTerminated{
								ExitCode: ErrORMBLogin,
							},
						},
					},
				}

				err = k8sClient.Status().Update(context.Background(), pod)
				if err != nil {
					return err
				}

				_, err = reconciler.Reconcile(req)
				if err != nil {
					return err
				}

				job = &batchv1.Job{}
				return k8sClient.Get(context.Background(), key, job)
			}, timeout, interval).Should(Succeed())

			By("Expecting to delete successfully")
			Eventually(func() error {
				getedModelJob := &modeljobsv1alpha1.ModelJob{}
				k8sClient.Get(context.Background(), key, getedModelJob)
				return k8sClient.Delete(context.Background(), getedModelJob)
			}, timeout, interval).Should(Succeed())

			By("Expecting to delete finish")
			Eventually(func() error {
				getedModelJob := &modeljobsv1alpha1.ModelJob{}
				return k8sClient.Get(context.Background(), key, getedModelJob)
			}, timeout, interval).ShouldNot(Succeed())
		})
	})
})
