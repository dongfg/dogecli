package logger

import (
	"github.com/sirupsen/logrus"
)

// Init initializes the logger
func Init() {
	logrus.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})
	logrus.SetLevel(logrus.InfoLevel)
}
