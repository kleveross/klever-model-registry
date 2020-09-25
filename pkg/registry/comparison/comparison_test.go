package comparison

import (
	"context"
	"reflect"
	"testing"

	ormbmodel "github.com/kleveross/ormb/pkg/model"

	"github.com/kleveross/klever-model-registry/pkg/registry/harbor"
	"github.com/kleveross/klever-model-registry/pkg/registry/paging"
)

func TestComposeComparison(t *testing.T) {
	tests := []struct {
		name        string
		comparison  Comparison
		expected    []*ormbmodel.Model
		expectedErr bool
	}{
		{
			name: "compare savedmodel with tensorrt",
			comparison: Comparison{
				Models: []ComparisonModel{
					{
						Name:    "tensorrt",
						Project: "release",
						Tag:     "v1",
					},
					{
						Name:    "savedmodel",
						Project: "release",
						Tag:     "v1",
					},
				},
			},
			expected: []*ormbmodel.Model{
				{
					Path: "/release/tensorrt:v1",
					Metadata: &ormbmodel.Metadata{
						Author: "Klever",
						Format: "TensorRT",
					},
				},
				{
					Path: "/release/savedmodel:v1",
					Metadata: &ormbmodel.Metadata{
						Author:    "Klever",
						Format:    "SavedModel",
						Framework: "TensorFlow",
					},
				},
			},
			expectedErr: false,
		},
	}

	for _, tt := range tests {
		proxy := harbor.NewFakeProxy()
		actual, err := composeComparison(tt.comparison.Models, proxy)
		if len(actual) != len(tt.expected) {
			t.Errorf("composeComparison() actual = %v, expected %v", actual, tt.expected)
		}
		for k, v := range actual {
			if !reflect.DeepEqual(v.Path, tt.expected[k].Path) &&
				!reflect.DeepEqual(v.Metadata, tt.expected[k].Metadata) {
				t.Errorf("composeComparison() actual = %v, expected %v", v, tt.expected[k])
			}
		}
		if (err != nil) != tt.expectedErr {
			t.Errorf("composeComparison() error = %v, expectedErr %v", err, tt.expectedErr)
			return
		}
	}
}

func Test_toORMBModelList(t *testing.T) {

	limit := 1

	type args struct {
		items []*ormbmodel.Model
		opt   *paging.ListOption
	}
	tests := []struct {
		name string
		args args
		want *ORMBModelList
	}{
		{
			name: "toORMBModelList successfully",
			args: args{
				items: []*ormbmodel.Model{
					&ormbmodel.Model{
						Path: "/test1",
					},
					&ormbmodel.Model{
						Path: "/test2",
					},
				},
				opt: &paging.ListOption{
					Start: 1,
					Limit: &limit,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := toORMBModelList(tt.args.items, tt.args.opt)
			if len(got.Items) != 1 {
				t.Errorf("len(got.Items) = %v, want 1", len(got.Items))
			}

			if got.ListMeta.TotalItems != 2 {
				t.Errorf("got.ListMeta.TotalItems = %v, want 2", got.ListMeta.TotalItems)
			}
		})
	}
}

func Test_generateCSVFile(t *testing.T) {
	type args struct {
		ctx   context.Context
		metas []*ormbmodel.Model
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "createCSVFile successfully",
			args: args{
				ctx: context.Background(),
				metas: []*ormbmodel.Model{
					&ormbmodel.Model{
						Metadata: &ormbmodel.Metadata{
							Signature: &ormbmodel.Signature{
								Inputs:  []ormbmodel.Tensor{},
								Outputs: []ormbmodel.Tensor{},
							},
						},
						Path: "/",
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := generateCSVFile(tt.args.ctx, "testFile", tt.args.metas); (err != nil) != tt.wantErr {
				t.Errorf("createCSVFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
