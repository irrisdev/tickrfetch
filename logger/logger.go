package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

var log *logrus.Logger

func init() {
	log = logrus.New()
	log.SetOutput(os.Stdout)
	log.SetLevel(logrus.InfoLevel)
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
}

// Info logs an info message
func Info(args ...interface{}) {
	log.Info(args...)
}

// Warn logs a warning message
func Warn(args ...interface{}) {
	log.Warn(args...)
}

// Fatal logs a fatal message
func Fatal(args ...interface{}) {
	log.Fatal(args...)
}

// Error logs an error message
func Error(args ...interface{}) {
	log.Error(args...)
}

// WithField logs a message with an additional field
func WithField(key string, value interface{}) *logrus.Entry {
	return log.WithField(key, value)
}

// WithFields logs a message with additional fields
func WithFields(fields logrus.Fields) *logrus.Entry {
	return log.WithFields(fields)
}
