package models

import (
	"github.com/kleveross/ormb/pkg/model"
)

type Model struct {
	ModelName   string         `json:"modelName"`
	VersionName string         `json:"versionName"`
	Description string         `json:"description"`
	Format      string         `json:"format"`
	FrameWork   string         `json:"framework"`
	Inputs      []model.Tensor `json:"inputs"`
	Outputs     []model.Tensor `json:"outputs"`
}
