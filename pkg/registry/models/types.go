package models

import (
	"github.com/kleveross/ormb/pkg/model"
)

type Model struct {
	ModelName   string         `json:"modelName,omitempty"`
	ProjectName string         `json:"projectName,omitempty"`
	VersionName string         `json:"versionName,omitempty"`
	Description string         `json:"description,omitempty"`
	Format      string         `json:"format,omitempty"`
	FrameWork   string         `json:"framework,omitempty"`
	Inputs      []model.Tensor `json:"inputs,omitempty"`
	Outputs     []model.Tensor `json:"outputs,omitempty"`
}
