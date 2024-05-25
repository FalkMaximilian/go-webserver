// logger/logger.go
package logger

import (
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

var Log = logrus.New()

func init() {

	// Get the log level from the environment variable
	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		logLevel = "debug" // Default log level
	}

	// Parse the log level string into a logrus log level
	level, err := logrus.ParseLevel(strings.ToLower(logLevel))
	if err != nil {
		Log.Fatal(err)
	}

	// Set the log level
	Log.SetLevel(level)
	Log.SetOutput(os.Stdout)
	Log.SetReportCaller(true)
	Log.SetFormatter(&logrus.JSONFormatter{})
}
