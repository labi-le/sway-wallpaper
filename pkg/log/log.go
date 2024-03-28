package log

import (
	"github.com/sirupsen/logrus"
)

var logger = MustLogger()

func MustLogger() *logrus.Logger {
	l := logrus.New()
	l.SetFormatter(&logrus.TextFormatter{
		DisableColors: false,
		FullTimestamp: true,
	})

	l.SetReportCaller(true)

	return l
}

func Fatal(args ...any) {
	logger.Fatal(args...)
}

func Fatalf(format string, v ...any) {
	logger.Fatalf(format, v...)
}

func Error(args ...any) {
	logger.Error(args...)
}

func Errorf(format string, v ...any) {
	logger.Errorf(format, v...)
}

func Info(args ...any) {
	logger.Info(args...)
}

func Infof(format string, v ...any) {
	logger.Infof(format, v...)
}

func Debug(args ...any) {
	logger.Debug(args...)
}

func Debugf(format string, v ...any) {
	logger.Debugf(format, v...)
}

func Warn(args ...any) {
	logger.Warn(args...)
}

func Warnf(format string, v ...any) {
	logger.Warnf(format, v...)
}
