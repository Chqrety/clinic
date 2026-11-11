package bootstrap

import (
	"log"
	"os"
)

type Logger struct {
	*log.Logger
}

func InitLogger(appName, logLevel string) *Logger {
	stdLogger := log.New(os.Stdout, fmtPrefix(appName), log.LstdFlags|log.Lmsgprefix)
	return &Logger{Logger: stdLogger}
}

func fmtPrefix(name string) string {
	return "[" + name + "] "
}

