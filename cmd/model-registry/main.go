package main

import (
	"github.com/caicloud/nirvana"
	"github.com/caicloud/nirvana/config"
	"github.com/caicloud/nirvana/log"
	"github.com/caicloud/nirvana/plugins/reqlog"
	"github.com/caicloud/nirvana/service"
	"github.com/spf13/viper"
	"sigs.k8s.io/controller-runtime/pkg/manager/signals"

	"github.com/kleveross/klever-model-registry/pkg/common"
	"github.com/kleveross/klever-model-registry/pkg/registry/apis"
	"github.com/kleveross/klever-model-registry/pkg/registry/apis/v1alpha1/descriptors"
	"github.com/kleveross/klever-model-registry/pkg/registry/client"
	registryconfig "github.com/kleveross/klever-model-registry/pkg/registry/config"
	"github.com/kleveross/klever-model-registry/pkg/registry/filters"
	"github.com/kleveross/klever-model-registry/pkg/registry/modifiers"
)

func main() {
	viper.AutomaticEnv()

	// Start nirvana
	option := &config.Option{
		Port: 8080,
	}
	customOption := registryconfig.New()
	cmd := config.NewNirvanaCommand(option)
	cmd.AddOption("ORMB", customOption)
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
	)

	// Set nirvana command hooks.
	cmd.SetHook(&config.NirvanaCommandHookFunc{
		PreConfigureFunc: func(c *nirvana.Config) error {
			log.Info("Harbor domain: ", customOption.Domain,
				", Harbor Username: ", customOption.Username)
			common.ORMBDomain = customOption.Domain
			common.ORMBPassword = customOption.Password
			common.ORMBUserName = customOption.Username
			c.Configure(
				nirvana.Descriptor(apis.AllDescriptors(
					customOption.Domain,
					customOption.Username,
					customOption.Password)...),
			)

			err := service.RegisterConsumer(service.NewSimpleSerializer("application/vnd.oci.image.manifest.v1+json"))
			if err != nil {
				return err
			}
			return nil
		},
		PreServeFunc: func(c *nirvana.Config, server nirvana.Server) error {
			if err := client.InitClient(customOption.KubeConfig,
				customOption.Domain, customOption.Username,
				customOption.Password, signals.SetupSignalHandler()); err != nil {
				return err
			}

			descriptors.InitModelJobController()
			descriptors.InitLogController()
			descriptors.InitEventController()
			descriptors.InitServingController()

			return nil
		},
	})

	if err := cmd.ExecuteWithConfig(serverConfig); err != nil {
		log.Fatal(err)
	}
	return
}
