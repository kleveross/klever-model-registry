package comparison

import (
	ormbmodel "github.com/kleveross/ormb/pkg/model"

	"github.com/kleveross/klever-model-registry/pkg/registry/paging"
)

type Comparison struct {
	Models []ComparisonModel `json:"models"`
}

type ComparisonModel struct {
	Name    string `json:"name"`
	Project string `json:"project"`
	Tag     string `json:"tag"`
}

type ORMBModelList struct {
	ListMeta paging.ListMeta    `json:"metadata"`
	Items    []*ormbmodel.Model `json:"items"`
}
