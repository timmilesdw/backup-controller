package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"github.com/timmilesdw/backup-controller/pkg/backupper"
	"github.com/timmilesdw/backup-controller/pkg/config"
	"github.com/timmilesdw/backup-controller/pkg/logger"
	"github.com/timmilesdw/backup-controller/pkg/metrics"
)

func init() {
	err := logger.InitLog("debug")
	if err != nil {
		panic(err)
	}
}

func main() {
	logrus.Info("Hello from backup-controller!")
	cfgPath, err := config.ParseFlags()
	if err != nil {
		logrus.Fatal(err)
	}
	cfg, err := config.NewConfig(cfgPath)
	if err != nil {
		if verr, ok := err.(validator.ValidationErrors); ok {
			for _, err := range verr {
				if err.Field() == "APIVersion" {
					logrus.Fatalf("Unknown APIVersion: %s", err.Value())
				}
				logrus.Fatal(err)
			}
		}
		logrus.Fatal(err)
	}
	logger.UpdateLogLevel(cfg.Spec.System)
	backupper := backupper.Backupper{
		ConfigSpec: cfg.Spec,
	}
	backupper.ConfigureCron()
	ms := metrics.RegisterMetrics(cfg.Spec)
	go ms.Start()
	logrus.Infof("Started metrics server on port %s, route %s", ms.Port, ms.Route)

	backupper.Start()

	quitChannel := make(chan os.Signal, 1)
	signal.Notify(quitChannel, syscall.SIGINT, syscall.SIGTERM)
	<-quitChannel
}
