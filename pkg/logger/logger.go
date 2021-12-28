package logger

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/timmilesdw/backup-controller/pkg/config"
)

func InitLog(level string) error {
	l, err := logrus.ParseLevel(level)
	if err != nil {
		return fmt.Errorf("parse logging level: %s", err)
	}
	logrus.SetLevel(l)
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	logrus.SetReportCaller(true)
	logrus.SetOutput(os.Stdout)
	return nil
}

func UpdateLogLevel(l config.System) {
	lvl, err := logrus.ParseLevel(l.LogLevel)
	if err != nil {
		logrus.Fatal(err)
	}
	logrus.SetLevel(lvl)
}
