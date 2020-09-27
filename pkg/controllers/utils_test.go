package controllers

import (
	"os"
	"reflect"
	"testing"

	"github.com/spf13/viper"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	modeljobsv1alpha1 "github.com/kleveross/klever-model-registry/pkg/apis/modeljob/v1alpha1"
	"github.com/kleveross/klever-model-registry/pkg/common"
	test "github.com/kleveross/klever-model-registry/testutil"
)

func Test_getFrameworkByFormat(t *testing.T) {
	initGlobalVar()

	type args struct {
		format modeljobsv1alpha1.Format
	}
	tests := []struct {
		name string
		args args
		want modeljobsv1alpha1.Framework
	}{
		{
			name: "caffe model",
			args: args{
				format: modeljobsv1alpha1.FormatCaffeModel,
			},
			want: modeljobsv1alpha1.FrameworkCaffe,
		},
		{
			name: "savedmodel",
			args: args{
				format: modeljobsv1alpha1.FormatSavedModel,
			},
			want: modeljobsv1alpha1.FrameworkTensorflow,
		},
		{
			name: "onnx",
			args: args{
				format: modeljobsv1alpha1.FormatONNX,
			},
			want: modeljobsv1alpha1.FrameworkONNX,
		},
		{
			name: "h5",
			args: args{
				format: modeljobsv1alpha1.FormatH5,
			},
			want: modeljobsv1alpha1.FrameworkKeras,
		},
		{
			name: "pmml",
			args: args{
				format: modeljobsv1alpha1.FormatPMML,
			},
			want: modeljobsv1alpha1.FrameworkPMML,
		},
		{
			name: "netdef",
			args: args{
				format: modeljobsv1alpha1.FormatNetDef,
			},
			want: modeljobsv1alpha1.FrameworkCaffe2,
		},
		{
			name: "mxnet",
			args: args{
				format: modeljobsv1alpha1.FormatMXNETParams,
			},
			want: modeljobsv1alpha1.FrameworkMXNet,
		},
		{
			name: "torchscript",
			args: args{
				format: modeljobsv1alpha1.FormatTorchScript,
			},
			want: modeljobsv1alpha1.FrameworkPyTorch,
		},
		{
			name: "graphdef",
			args: args{
				format: modeljobsv1alpha1.FormatGraphDef,
			},
			want: modeljobsv1alpha1.FrameworkTensorflow,
		},
		{
			name: "tensorrt",
			args: args{
				format: modeljobsv1alpha1.FormatTensorRT,
			},
			want: modeljobsv1alpha1.FrameworkTensorRT,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getFrameworkByFormat(tt.args.format); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getFrameworkByFormat() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_generateJobResource(t *testing.T) {
	initGlobalVar()
	PresetAnalyzeImageConfig = test.InitPresetModelImageConfigMap()
	conversionDesiredTag := "release/savedmodel:v2"

	type args struct {
		modeljob *modeljobsv1alpha1.ModelJob
	}
	tests := []struct {
		name    string
		args    args
		want    *batchv1.Job
		wantErr bool
	}{
		{
			name: "Extraction successfully",
			args: args{
				modeljob: &modeljobsv1alpha1.ModelJob{
					Spec: modeljobsv1alpha1.ModelJobSpec{
						Model: "release/savedmodel:v1",
						ModelJobSource: modeljobsv1alpha1.ModelJobSource{
							Extraction: &modeljobsv1alpha1.ExtractionSource{
								Format: modeljobsv1alpha1.FormatSavedModel,
							},
						},
					},
				},
			},
			want: &batchv1.Job{
				Spec: batchv1.JobSpec{
					Template: corev1.PodTemplateSpec{
						Spec: corev1.PodSpec{
							Containers: []corev1.Container{
								corev1.Container{
									Env: []corev1.EnvVar{
										corev1.EnvVar{
											Name:  modeljobsv1alpha1.FrameworkEnvKey,
											Value: string(modeljobsv1alpha1.FrameworkTensorflow),
										},
									},
								},
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Conversition Successfully",
			args: args{
				modeljob: &modeljobsv1alpha1.ModelJob{
					Spec: modeljobsv1alpha1.ModelJobSpec{
						Model:      "release/savedmodel:v1",
						DesiredTag: &conversionDesiredTag,
						ModelJobSource: modeljobsv1alpha1.ModelJobSource{
							Conversion: &modeljobsv1alpha1.ConversionSource{
								MMdnn: &modeljobsv1alpha1.MMdnnSpec{
									ConversionBaseSpec: modeljobsv1alpha1.ConversionBaseSpec{
										From: modeljobsv1alpha1.FormatH5,
										To:   modeljobsv1alpha1.FormatSavedModel,
									},
								},
							},
						},
					},
				},
			},
			want: &batchv1.Job{
				Spec: batchv1.JobSpec{
					Template: corev1.PodTemplateSpec{
						Spec: corev1.PodSpec{
							Containers: []corev1.Container{
								corev1.Container{
									Env: []corev1.EnvVar{
										corev1.EnvVar{
											Name:  modeljobsv1alpha1.FrameworkEnvKey,
											Value: string(modeljobsv1alpha1.FrameworkTensorflow),
										},
									},
								},
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Conversition Failed, Desired tag is nil",
			args: args{
				modeljob: &modeljobsv1alpha1.ModelJob{
					Spec: modeljobsv1alpha1.ModelJobSpec{
						Model: "model tag",
						ModelJobSource: modeljobsv1alpha1.ModelJobSource{
							Conversion: &modeljobsv1alpha1.ConversionSource{
								MMdnn: &modeljobsv1alpha1.MMdnnSpec{
									ConversionBaseSpec: modeljobsv1alpha1.ConversionBaseSpec{
										From: modeljobsv1alpha1.FormatH5,
										To:   modeljobsv1alpha1.FormatSavedModel,
									},
								},
							},
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Error, None ModelJobSource",
			args: args{
				modeljob: &modeljobsv1alpha1.ModelJob{
					Spec: modeljobsv1alpha1.ModelJobSpec{
						ModelJobSource: modeljobsv1alpha1.ModelJobSource{},
					},
				},
			},
			want: &batchv1.Job{
				Spec: batchv1.JobSpec{
					Template: corev1.PodTemplateSpec{
						Spec: corev1.PodSpec{
							Containers: []corev1.Container{
								corev1.Container{
									Env: []corev1.EnvVar{},
								},
							},
						},
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := generateJobResource(tt.args.modeljob)
			if (err != nil) != tt.wantErr {
				t.Errorf("generateJobResource() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err != nil {
				return
			}

			succ := false
			for _, gotEnv := range got.Spec.Template.Spec.Containers[0].Env {
				for _, wantEnv := range tt.want.Spec.Template.Spec.Containers[0].Env {
					if gotEnv.Name == wantEnv.Name && gotEnv.Value == wantEnv.Value {
						succ = true
						break
					}
				}
			}
			if !succ {
				t.Errorf("generateJobResource() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_generateInitContainers(t *testing.T) {
	viper.AutomaticEnv()
	os.Setenv(common.ORMBDomainEnvKey, "demo.goharbo.com")
	os.Setenv(common.ORMBUsernameEnvkey, "ormbtest")
	os.Setenv(common.ORMBPasswordEnvKey, "ORMBtest12345")
	initGlobalVar()
	PresetAnalyzeImageConfig = test.InitPresetModelImageConfigMap()

	type args struct {
		modeljob *modeljobsv1alpha1.ModelJob
	}
	tests := []struct {
		name    string
		args    args
		want    []corev1.Container
		wantErr bool
	}{
		{
			name: "generateInitContainers successfully",
			args: args{
				modeljob: &modeljobsv1alpha1.ModelJob{
					ObjectMeta: metav1.ObjectMeta{
						Name: "test-modeljob",
					},
					Spec: modeljobsv1alpha1.ModelJobSpec{
						Model: "demo.goharbor.com/release/testmodel:v1",
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := generateInitContainers(tt.args.modeljob)
			if (err != nil) != tt.wantErr {
				t.Errorf("generateInitContainers() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got[0].Args[0] != tt.args.modeljob.Spec.Model {
				t.Errorf("generateInitContainers() = %v, want %v", got, tt.want)
			}
		})
	}
}
