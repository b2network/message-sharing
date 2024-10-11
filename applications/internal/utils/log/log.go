package log

import (
	log "github.com/sirupsen/logrus"
	"os"
)

type Logger struct {
	Name string
}

func NewLogger(name string, logLevel uint32) *Logger {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)

	log.SetLevel(log.Level(logLevel))
	return &Logger{
		Name: name,
	}
}

func (l *Logger) Debug(args ...interface{}) {
	log.WithField("name", l.Name).Debug(args...)
}

func (l *Logger) Warn(args ...interface{}) {
	log.WithField("name", l.Name).Warn(args...)
}

func (l *Logger) Info(args ...interface{}) {
	log.WithField("name", l.Name).Info(args...)
}

func (l *Logger) Error(args ...interface{}) {
	log.WithField("name", l.Name).Error(args...)
}

func (l *Logger) Panic(args ...interface{}) {
	log.WithField("name", l.Name).Panic(args...)
}

func (l *Logger) Debugf(format string, args ...interface{}) {
	log.WithField("name", l.Name).Debugf(format, args...)
}

func (l *Logger) Warnf(format string, args ...interface{}) {
	log.WithField("name", l.Name).Warnf(format, args...)
}

func (l *Logger) Infof(format string, args ...interface{}) {
	log.WithField("name", l.Name).Infof(format, args...)
}

func (l *Logger) Errorf(format string, args ...interface{}) {
	log.WithField("name", l.Name).Errorf(format, args...)
}

func (l *Logger) Panicf(format string, args ...interface{}) {
	log.WithField("name", l.Name).Panicf(format, args...)
}
