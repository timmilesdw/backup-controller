package main

import (
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"github.com/timmilesdw/backup-controller/pkg/backupper"
	"github.com/timmilesdw/backup-controller/pkg/config"
	"github.com/timmilesdw/backup-controller/pkg/exporters"
	"github.com/timmilesdw/backup-controller/pkg/logger"
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
				logrus.Fatal(err)
			}
		}
		logrus.Fatal(err)
	}
	logger.UpdateLogLevel(cfg.Logger)
	exporters.PopulateExporters(cfg.Backupper.Databases)
	exporters.PopulateStorers(cfg.Backupper.Storages)
	backupper := backupper.Backupper{
		ConfigSpec: cfg.Backupper,
	}
	backupper.ConfigureCron()
	// ms := metrics.RegisterMetrics(cfg.Spec)
	// go ms.Start()
	// logrus.Infof("Started metrics server on port %s, route %s", ms.Port, ms.Route)

	// go backupper.Start()
	// go server.StartServer(cfg.Spec)
	// quitChannel := make(chan os.Signal, 1)
	// signal.Notify(quitChannel, syscall.SIGINT, syscall.SIGTERM)
	// <-quitChannel
}
