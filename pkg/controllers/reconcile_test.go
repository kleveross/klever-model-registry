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
			want: errContainerCreating,
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
			want: errORMBLogin,
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
			want: errORMBPull,
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
			want: errORMBExport,
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
			want: errRunTask,
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
			want: errORMBSave,
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
			want: errORMBPush,
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
