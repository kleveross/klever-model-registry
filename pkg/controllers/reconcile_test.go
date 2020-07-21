package controllers

import (
	"testing"

	corev1 "k8s.io/api/core/v1"
)

func Test_getModelJobFailedMesage(t *testing.T) {
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
			want: "",
		},
		{
			name: "ormb login err",
			args: args{
				pods: &corev1.PodList{
					Items: []corev1.Pod{
						{
							Status: corev1.PodStatus{
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
			want: "ormb login error",
		},
		{
			name: "ormb pull model error",
			args: args{
				pods: &corev1.PodList{
					Items: []corev1.Pod{
						{
							Status: corev1.PodStatus{
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
			want: "ormb pull model error",
		},
		{
			name: "ormb export model error",
			args: args{
				pods: &corev1.PodList{
					Items: []corev1.Pod{
						{
							Status: corev1.PodStatus{
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
			want: "ormb export model error",
		},
		{
			name: "run task error",
			args: args{
				pods: &corev1.PodList{
					Items: []corev1.Pod{
						{
							Status: corev1.PodStatus{
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
			want: "run task error",
		},
		{
			name: "ormb save model error",
			args: args{
				pods: &corev1.PodList{
					Items: []corev1.Pod{
						{
							Status: corev1.PodStatus{
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
			want: "ormb save model error",
		},
		{
			name: "ormb push model error",
			args: args{
				pods: &corev1.PodList{
					Items: []corev1.Pod{
						{
							Status: corev1.PodStatus{
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
			want: "ormb push model error",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getModelJobFailedMesage(tt.args.pods); got != tt.want {
				t.Errorf("getModelJobFailedMesage() = %v, want %v", got, tt.want)
			}
		})
	}
}
