package comparison

import (
	"reflect"
	"testing"

	ormbmodel "github.com/kleveross/ormb/pkg/model"

	"github.com/kleveross/klever-model-registry/pkg/registry/harbor"
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
