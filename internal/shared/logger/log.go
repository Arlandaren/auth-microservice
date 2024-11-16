package logger

import (
	log "github.com/sirupsen/logrus"
	"service/internal/shared/config"
)

type DefaultFieldsHook struct {
	Service string
	Env     string
}

func Init() {
	log.SetFormatter(&log.JSONFormatter{})

	log.SetLevel(log.InfoLevel)

	log.AddHook(&DefaultFieldsHook{
		Service: "auth_service",
		Env:     config.GetEnvironment(),
	})
}

func (h *DefaultFieldsHook) Levels() []log.Level {
	return log.AllLevels
}

func (h *DefaultFieldsHook) Fire(entry *log.Entry) error {
	entry.Data["service"] = h.Service
	entry.Data["env"] = h.Env
	return nil
}

func Info(args ...interface{}) {
	log.Info(args...)
}

func Infof(format string, args ...interface{}) {
	log.Infof(format, args...)
}

func Error(args ...interface{}) {
	log.Error(args...)
}

func Errorf(format string, args ...interface{}) {
	log.Errorf(format, args...)
}

func Fatal(args ...interface{}) {
	log.Fatal(args...)
}

func Fatalf(format string, args ...interface{}) {
	log.Fatalf(format, args...)
}
