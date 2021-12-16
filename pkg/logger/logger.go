package logger

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
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
	logrus.SetOutput(os.Stdout)
	return nil
}
