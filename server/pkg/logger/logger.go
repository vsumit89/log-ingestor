package logger

import (
	"net/http"
)

var logInstance ILogger

type ILogger interface {
	Info(message string, args ...interface{})
	Warn(message string, args ...interface{})
	Error(message string, args ...interface{})
	GetHTTPMiddleWare() func(next http.Handler) http.Handler
	Fatal(message string, args ...interface{})
}

func GetInstance() ILogger {
	if logInstance != nil {
		return logInstance
	}
	logInstance = NewSlogAdapter()
	return logInstance
}
