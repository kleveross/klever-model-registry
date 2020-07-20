package modeljob

import (
	"testing"

	"github.com/kleveross/klever-model-registry/pkg/common"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestExchangeModelJobNameAndID(t *testing.T) {
	type args struct {
		objectMeta *metav1.ObjectMeta
	}
	case1ModelJobName := "modelJobName"
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "case-1",
			args: args{
				objectMeta: &metav1.ObjectMeta{
					Name: case1ModelJobName,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ExchangeModelJobNameAndID(tt.args.objectMeta); (err != nil) != tt.wantErr ||
				tt.args.objectMeta.Labels[common.ResourceNameLabelKey] != case1ModelJobName {
				t.Errorf("ExchangeModelJobNameAndID() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
