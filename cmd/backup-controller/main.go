package main

import (
	"os"
	"os/signal"
	"syscall"

	"git.ufanet.ru/hw-k8s/software/backup-controller/pkg/backupper"
	"git.ufanet.ru/hw-k8s/software/backup-controller/pkg/config"
	"git.ufanet.ru/hw-k8s/software/backup-controller/pkg/logger"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
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
	backupper := backupper.Backupper{
		ConfigSpec: cfg.Spec,
	}
	backupper.Start()

	quitChannel := make(chan os.Signal, 1)
	signal.Notify(quitChannel, syscall.SIGINT, syscall.SIGTERM)
	<-quitChannel
}
