package controllers

import (
	"context"
	"github.com/spf13/viper"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

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

	kubeconfigPath := viper.GetString("kubeconfig")

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err != nil {
		return err
	}

	kubeClient, err := kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}

	PresetAnalyzeImageConfig, err = kubeClient.CoreV1().ConfigMaps("kleveross-system").Get(context.TODO(), "modeljob-image-config", metav1.GetOptions{})
	if err != nil {
		return err
	}

	return nil
}
