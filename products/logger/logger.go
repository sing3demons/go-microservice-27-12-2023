package logger

import (
	"io"
	"log"
	"os"

	"github.com/sirupsen/logrus"
)

type Logger struct {
	*logrus.Logger
}

func NewLogger() *Logger {
	logLevel, err := logrus.ParseLevel(os.Getenv("LOG_LEVEL"))
	if err != nil {
		logLevel = logrus.InfoLevel
	}

	logger := logrus.New()
	logger.Formatter = &logrus.JSONFormatter{}
	logrus.SetLevel(logLevel)
	log.SetOutput(logger.Writer())
	logger.SetOutput(io.MultiWriter(os.Stdout))
	return &Logger{logger}
}
