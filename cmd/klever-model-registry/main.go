package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/caicloud/nirvana"
	"github.com/caicloud/nirvana/config"
	"github.com/caicloud/nirvana/log"
	"github.com/caicloud/nirvana/plugins/reqlog"
	"github.com/spf13/viper"

	"github.com/kleveross/klever-model-registry/pkg/registry/apis"
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

func gracefulShutdown(closing, done chan struct{}) {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	log.Infof("capture system signal %s, to close \"closing\" channel", <-signals)
	close(closing)
	select {
	case <-done:
		log.Infof("Goroutines exited normally")
	case <-time.After(time.Second * 3):
		log.Infof("Timeout waiting goroutines to exit")
	}
	os.Exit(0)
}

func main() {
	err := client.InitClient()
	if err != nil {
		log.Fatal(err)
	}

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

	if err := cmd.ExecuteWithConfig(serverConfig); err != nil {
		log.Fatal(err)
	}
	return
}
