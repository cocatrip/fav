package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

var log = logrus.New()

func init() {
	// log.SetReportCaller(true)

	log.SetOutput(os.Stdout)

	log.SetLevel(logrus.DebugLevel)
}

func GetLogger() *logrus.Logger {
	return log
}
