package utils

import (
	"fmt"
	"log"
	"os"
	"time"
)

type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
	FATAL
)

var Logger = NewLoggerUtil(INFO)

type LoggerUtil struct {
	level      LogLevel
	LoggerUtil *log.Logger
}

func NewLoggerUtil(level LogLevel) *LoggerUtil {
	return &LoggerUtil{
		level:      level,
		LoggerUtil: log.New(os.Stdout, "", 0),
	}
}

func (l *LoggerUtil) log(level LogLevel, label string, format string, args ...any) {
	if level < l.level {
		return
	}

	msg := fmt.Sprintf(format, args...)

	l.LoggerUtil.Printf(
		"[%s] [%s] %s",
		time.Now().Format("2006-01-02 15:04:05"),
		label,
		msg,
	)

	if level == FATAL {
		os.Exit(1)
	}
}

func (l *LoggerUtil) Debug(format string, args ...any) {
	l.log(DEBUG, "DEBUG", format, args...)
}

func (l *LoggerUtil) Info(format string, args ...any) {
	l.log(INFO, "INFO", format, args...)
}

func (l *LoggerUtil) Warn(format string, args ...any) {
	l.log(WARN, "WARN", format, args...)
}

func (l *LoggerUtil) Error(format string, args ...any) {
	l.log(ERROR, "ERROR", format, args...)
}

func (l *LoggerUtil) Fatal(format string, args ...any) {
	l.log(FATAL, "FATAL", format, args...)
}
