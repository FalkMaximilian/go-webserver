// logger/logger.go
package logger

import (
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

var log = logrus.New()

func init() {
	// Get the log level from the environment variable
	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		logLevel = "debug" // Default log level
	}

	// Parse the log level string into a logrus log level
	level, err := logrus.ParseLevel(strings.ToLower(logLevel))
	if err != nil {
		log.Fatal(err)
	}

	// Set the log level
	log.SetLevel(level)
}

// Exported functions for logging
func Info(args ...interface{}) {
	log.Info(args...)
}

func Debug(args ...interface{}) {
	log.Debug(args...)
}

func Warn(args ...interface{}) {
	log.Warn(args...)
}

func Error(args ...interface{}) {
	log.Error(args...)
}

func Fatal(args ...interface{}) {
	log.Fatal(args...)
}
