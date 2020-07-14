package modeljob

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/caicloud/temp-model-registry/pkg/common"
	"github.com/caicloud/temp-model-registry/pkg/util"
)

// ExchangeModelJobNameAndID set crd name in labels field, and set uuid in name field
func ExchangeModelJobNameAndID(objectMeta *metav1.ObjectMeta) error {
	id := util.RandomNameWithPrefix("modeljob")

	if objectMeta.Labels == nil {
		objectMeta.Labels = map[string]string{}
	}
	objectMeta.Labels[common.ResourceNameLabelKey] = objectMeta.Name
	objectMeta.Name = id
	return nil
}
