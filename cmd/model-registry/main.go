package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/astaxie/beego"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/caicloud/temp-model-registry/pkg/registry/server"
)

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
	viper.AutomaticEnv()

	beego.BConfig.CopyRequestBody = true

	server.RegisterRoutes()

	closing := make(chan struct{})
	done := make(chan struct{})
	go gracefulShutdown(closing, done)

	beego.Run()

	return
}
