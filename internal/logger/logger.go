package logger

import (
	"github.com/dongfg/dogecli/internal/constants"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Init initializes the logger
func Init() {
	logrus.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})
	if viper.GetBool(constants.VerboseMode) {
		logrus.SetLevel(logrus.DebugLevel)
	} else {
		logrus.SetLevel(logrus.InfoLevel)
	}
}
