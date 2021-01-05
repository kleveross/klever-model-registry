package controllers

import (
	"testing"

	corev1 "k8s.io/api/core/v1"
)

func Test_getModelJobMesageByPods(t *testing.T) {
	type args struct {
		pods *corev1.PodList
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "no pods",
			args: args{
				pods: &corev1.PodList{
					Items: []corev1.Pod{},
				},
			},
			want: "waiting pod to create",
		},
		{
			name: "ormb login err",
			args: args{
				pods: &corev1.PodList{
					Items: []corev1.Pod{
						{
							Status: corev1.PodStatus{
								InitContainerStatuses: []corev1.ContainerStatus{
									{
										State: corev1.ContainerState{
											Terminated: &corev1.ContainerStateTerminated{
												ExitCode: 0,
											},
										},
									},
								},
								ContainerStatuses: []corev1.ContainerStatus{
									{
										State: corev1.ContainerState{
											Terminated: &corev1.ContainerStateTerminated{
												ExitCode: ErrORMBLogin,
											},
										},
									},
								},
							},
						},
					},
				},
			},
			want: "failed to login model registry",
		},
		{
			name: "ormb pull model error",
			args: args{
				pods: &corev1.PodList{
					Items: []corev1.Pod{
						{
							Status: corev1.PodStatus{
								InitContainerStatuses: []corev1.ContainerStatus{
									{
										State: corev1.ContainerState{
											Terminated: &corev1.ContainerStateTerminated{
												ExitCode: 0,
											},
										},
									},
								},
								ContainerStatuses: []corev1.ContainerStatus{
									{
										State: corev1.ContainerState{
											Terminated: &corev1.ContainerStateTerminated{
												ExitCode: ErrORMBPullModel,
											},
										},
									},
								},
							},
						},
					},
				},
			},
			want: "failed to pull model from model registry",
		},
		{
			name: "ormb export model error",
			args: args{
				pods: &corev1.PodList{
					Items: []corev1.Pod{
						{
							Status: corev1.PodStatus{
								InitContainerStatuses: []corev1.ContainerStatus{
									{
										State: corev1.ContainerState{
											Terminated: &corev1.ContainerStateTerminated{
												ExitCode: 0,
											},
										},
									},
								},
								ContainerStatuses: []corev1.ContainerStatus{
									{
										State: corev1.ContainerState{
											Terminated: &corev1.ContainerStateTerminated{
												ExitCode: ErrORMBExportModel,
											},
										},
									},
								},
							},
						},
					},
				},
			},
			want: "failed to export model to local",
		},
		{
			name: "run task error",
			args: args{
				pods: &corev1.PodList{
					Items: []corev1.Pod{
						{
							Status: corev1.PodStatus{
								InitContainerStatuses: []corev1.ContainerStatus{
									{
										State: corev1.ContainerState{
											Terminated: &corev1.ContainerStateTerminated{
												ExitCode: 0,
											},
										},
									},
								},
								ContainerStatuses: []corev1.ContainerStatus{
									{
										State: corev1.ContainerState{
											Terminated: &corev1.ContainerStateTerminated{
												ExitCode: ErrRunTask,
											},
										},
									},
								},
							},
						},
					},
				},
			},
			want: "failed to run extract/convert task",
		},
		{
			name: "ormb save model error",
			args: args{
				pods: &corev1.PodList{
					Items: []corev1.Pod{
						{
							Status: corev1.PodStatus{
								InitContainerStatuses: []corev1.ContainerStatus{
									{
										State: corev1.ContainerState{
											Terminated: &corev1.ContainerStateTerminated{
												ExitCode: 0,
											},
										},
									},
								},
								ContainerStatuses: []corev1.ContainerStatus{
									{
										State: corev1.ContainerState{
											Terminated: &corev1.ContainerStateTerminated{
												ExitCode: ErrORMBSaveModel,
											},
										},
									},
								},
							},
						},
					},
				},
			},
			want: "failed to save model to localhost",
		},
		{
			name: "ormb push model error",
			args: args{
				pods: &corev1.PodList{
					Items: []corev1.Pod{
						{
							Status: corev1.PodStatus{
								InitContainerStatuses: []corev1.ContainerStatus{
									{
										State: corev1.ContainerState{
											Terminated: &corev1.ContainerStateTerminated{
												ExitCode: 0,
											},
										},
									},
								},
								ContainerStatuses: []corev1.ContainerStatus{
									{
										State: corev1.ContainerState{
											Terminated: &corev1.ContainerStateTerminated{
												ExitCode: ErrORMBPushModel,
											},
										},
									},
								},
							},
						},
					},
				},
			},
			want: "failed to push model to model registry",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := getModelJobMesageByPods(tt.args.pods); got != tt.want {
				t.Errorf("getModelJobMesage() = %v, want %v", got, tt.want)
			}
		})
	}
}
