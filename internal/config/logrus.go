package config

import (
	"github.com/sirupsen/logrus"
	"os"
	"strconv"
)

func NewLogrus() *logrus.Logger {
	log := logrus.New()

	logLevelInt, err := strconv.Atoi(os.Getenv("APP_LOG_LEVEL"))
	if err != nil {
		panic("failed to convert APP_LOG_LEVEL to int: " + err.Error())
	}

	log.SetLevel(logrus.Level(logLevelInt))
	log.SetFormatter(&logrus.JSONFormatter{})

	return log
}
