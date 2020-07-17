package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/caicloud/nirvana/config"
	"github.com/caicloud/nirvana/log"
	"github.com/caicloud/nirvana/plugins/reqlog"
	"github.com/spf13/viper"

	"github.com/kleveross/klever-model-registry/pkg/registry/apis"
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
	// Start nirvana
	option := &config.Option{
		Port: uint16(viper.GetInt(kleverModelRegistryPort)),
	}
	nirvana := config.NewNirvanaCommand(option)
	nirvana.EnablePlugin(
		&reqlog.Option{
			DoubleLog:  true,
			SourceAddr: true,
			RequestID:  true,
		},
	)

	if err := nirvana.Execute(
		apis.AllDescriptor...); err != nil {
		log.Fatal(err)
	}

	return
}
