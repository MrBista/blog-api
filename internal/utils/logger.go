package utils

import "github.com/sirupsen/logrus"

var Logger *logrus.Logger

func InitLogger() {
	if Logger == nil {
		Logger = logrus.New()
		Logger.SetLevel(logrus.TraceLevel)
		Logger.SetFormatter(&logrus.JSONFormatter{})
	}
}
