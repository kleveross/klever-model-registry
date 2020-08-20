package main

import (
	"github.com/caicloud/nirvana"
	"github.com/caicloud/nirvana/config"
	"github.com/caicloud/nirvana/log"
	"github.com/caicloud/nirvana/plugins/reqlog"
	"github.com/spf13/viper"
	"sigs.k8s.io/controller-runtime/pkg/manager/signals"

	"github.com/kleveross/klever-model-registry/pkg/registry/apis"
	"github.com/kleveross/klever-model-registry/pkg/registry/apis/v1alpha1/descriptors"
	"github.com/kleveross/klever-model-registry/pkg/registry/client"
	"github.com/kleveross/klever-model-registry/pkg/registry/filters"
	"github.com/kleveross/klever-model-registry/pkg/registry/modifiers"
)

const (
	// kleverModelRegistryPort is model registry default port, default 8080
	kleverModelRegistryPort = "klever_model_registry_port"
)

func init() {
	viper.AutomaticEnv()
	viper.SetDefault(kleverModelRegistryPort, "8080")
}

func main() {
	// Start nirvana
	option := &config.Option{
		Port: uint16(viper.GetInt(kleverModelRegistryPort)),
	}
	cmd := config.NewNirvanaCommand(option)
	cmd.EnablePlugin(
		&reqlog.Option{
			DoubleLog:  true,
			SourceAddr: true,
			RequestID:  true,
		},
	)

	serverConfig := nirvana.NewConfig()
	serverConfig.Configure(
		nirvana.Logger(log.DefaultLogger()),
		nirvana.Filter(filters.Filters()...),
		nirvana.Modifier(modifiers.Modifiers()...),
		nirvana.Descriptor(apis.AllDescriptor...),
	)

	// Set nirvana command hooks.
	cmd.SetHook(&config.NirvanaCommandHookFunc{
		PreServeFunc: func(c *nirvana.Config, server nirvana.Server) error {
			if err := client.InitClient(signals.SetupSignalHandler()); err != nil {
				log.Fatal(err)
			}
			descriptors.InitModelJobController()
			descriptors.InitLogController()
			descriptors.InitEventController()

			return nil
		},
	})

	if err := cmd.ExecuteWithConfig(serverConfig); err != nil {
		log.Fatal(err)
	}
	return
}
