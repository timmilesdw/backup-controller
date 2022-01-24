package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"github.com/timmilesdw/backup-controller/pkg/backupper"
	"github.com/timmilesdw/backup-controller/pkg/config"
	"github.com/timmilesdw/backup-controller/pkg/exporters"
	"github.com/timmilesdw/backup-controller/pkg/logger"
	"github.com/timmilesdw/backup-controller/pkg/metrics"
	"github.com/timmilesdw/backup-controller/pkg/server"
)

func init() {
	err := logger.InitLog("debug")
	if err != nil {
		panic(err)
	}
}

func main() {
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
	if *cfg.Metrics.Enabled {
		ms := metrics.RegisterMetrics(cfg.Metrics)
		go ms.Start()
		logrus.Infof("Started metrics server on port %s, route %s", ms.Port, ms.Route)
	}
	if *cfg.Backupper.Enabled {
		exporters.PopulateExporters(cfg.Backupper.Exporters)
		exporters.PopulateStorers(cfg.Backupper.Storers)
		backupper := backupper.Backupper{
			ConfigSpec: cfg.Backupper,
		}
		backupper.ConfigureCron()
		go backupper.Start()
	}
	if *cfg.UI.Enabled {
		go server.StartServer(cfg.UI)
	}
	quitChannel := make(chan os.Signal, 1)
	signal.Notify(quitChannel, syscall.SIGINT, syscall.SIGTERM)
	<-quitChannel
}
