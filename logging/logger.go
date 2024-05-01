package logging 

import "github.com/sirupsen/logrus"

var (
  Logger *logrus.Logger
)

func SetupLogger() {
  Logger = logrus.New()
  Logger.Formatter = &logrus.JSONFormatter{}  // Log in JSON format
  Logger.Level = logrus.DebugLevel             // Set the logging level
}
