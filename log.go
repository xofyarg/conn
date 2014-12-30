package main

import (
	"fmt"
	"log"
)

const (
	LogNone = iota
	LogError
	LogWarn
	LogInfo
	LogDebug
)

var _logLevel = LogNone

var banner = []string{
	"NONE",
	"ERROR",
	"WARN",
	"INFO",
	"DEBUG",
}

func writeLog(l int, v ...interface{}) {
	if l > _logLevel {
		return
	}

	if len(v) == 1 {
		log.Printf("[%s] %v", banner[l], v[0])
	} else {
		f := fmt.Sprintf("[%s] %s", banner[l], v[0])
		log.Printf(f, v[1:]...)
	}
}

func LogLevel(l int) {
	if l <= len(banner) {
		_logLevel = l
	} else {
		Warn("invalid level")
	}
}

func Debug(v ...interface{}) {
	writeLog(LogDebug, v...)
}

func Warn(v ...interface{}) {
	writeLog(LogWarn, v...)
}

func Info(v ...interface{}) {
	writeLog(LogInfo, v...)
}

func Error(v ...interface{}) {
	writeLog(LogError, v...)
}
