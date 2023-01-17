package log

import (
	"time"

	"github.com/sirupsen/logrus"
	"watchmen/config"
)

func SetupLogger(cfg config.Logger) {
	logLevel, err := logrus.ParseLevel(cfg.Level)
	if err != nil {
		logLevel = logrus.ErrorLevel
	}

	logrus.SetLevel(logLevel)
	logrus.SetReportCaller(true)
	logrus.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.RFC3339,
	})
}
