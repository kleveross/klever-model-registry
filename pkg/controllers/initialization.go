package controllers

import (
	"github.com/spf13/viper"

	modeljobsv1alpha1 "github.com/kleveross/klever-model-registry/pkg/apis/modeljob/v1alpha1"
	"github.com/kleveross/klever-model-registry/pkg/common"
)

func initGlobalVar() {
	ModelFormatToFrameworkMapping = map[modeljobsv1alpha1.Format]modeljobsv1alpha1.Framework{
		modeljobsv1alpha1.FormatSavedModel:  modeljobsv1alpha1.FrameworkTensorflow,
		modeljobsv1alpha1.FormatONNX:        modeljobsv1alpha1.FrameworkONNX,
		modeljobsv1alpha1.FormatH5:          modeljobsv1alpha1.FrameworkKeras,
		modeljobsv1alpha1.FormatPMML:        modeljobsv1alpha1.FrameworkPMML,
		modeljobsv1alpha1.FormatCaffeModel:  modeljobsv1alpha1.FrameworkCaffe,
		modeljobsv1alpha1.FormatNetDef:      modeljobsv1alpha1.FrameworkCaffe2,
		modeljobsv1alpha1.FormatMXNETParams: modeljobsv1alpha1.FrameworkMXNet,
		modeljobsv1alpha1.FormatTorchScript: modeljobsv1alpha1.FrameworkPyTorch,
		modeljobsv1alpha1.FormatGraphDef:    modeljobsv1alpha1.FrameworkTensorflow,
		modeljobsv1alpha1.FormatTensorRT:    modeljobsv1alpha1.FrameworkTensorRT,
	}

	common.ORMBDomain = viper.GetString(common.ORMBDomainEnvKey)
	common.ORMBUserName = viper.GetString(common.ORMBUsernameEnvkey)
	common.ORMBPassword = viper.GetString(common.ORMBPasswordEnvKey)
}

func Initialization() error {
	viper.AutomaticEnv()

	initGlobalVar()

	return nil
}
