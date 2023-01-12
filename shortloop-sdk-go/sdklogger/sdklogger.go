package sdklogger

import (
	"log"
	"os"
)

type SDKLogger struct {
	loggingEnabled bool
	logLevel       LogLevel
	infoLogger     *log.Logger
	errorLogger    *log.Logger
}

type LogLevel string

const (
	INFO  LogLevel = "INFO"
	ERROR LogLevel = "ERROR"
)

func GetLogLevel(logLevel string) LogLevel {
	switch logLevel {
	case "INFO":
		return INFO
	case "ERROR":
		return ERROR
	default:
		return ERROR
	}
}

func (logger *SDKLogger) Info(v ...interface{}) {
	if logger.loggingEnabled && logger.logLevel == INFO && logger.infoLogger != nil {
		logger.infoLogger.Println(v...)
	}
}

func (logger *SDKLogger) Error(v ...interface{}) {
	if logger.loggingEnabled && (logger.logLevel == INFO || logger.logLevel == ERROR) && logger.errorLogger != nil {
		logger.errorLogger.Println(v...)
	}
}

func (logger *SDKLogger) InfoF(format string, v ...interface{}) {
	if logger.loggingEnabled && logger.logLevel == INFO && logger.infoLogger != nil {
		logger.infoLogger.Printf(format, v...)
	}
}

func (logger *SDKLogger) ErrorF(format string, v ...interface{}) {
	if logger.loggingEnabled && (logger.logLevel == INFO || logger.logLevel == ERROR) && logger.errorLogger != nil {
		logger.errorLogger.Printf(format, v...)
	}
}

func (logger *SDKLogger) SetLoggingEnabled(loggingEnabled bool) {
	logger.loggingEnabled = loggingEnabled
}

func (logger *SDKLogger) SetLogLevel(logLevel LogLevel) {
	logger.logLevel = logLevel
}

var Logger *SDKLogger

func init() {
	Logger = &SDKLogger{
		loggingEnabled: false,
		logLevel:       ERROR,
		infoLogger:     log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime),
		errorLogger:    log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime),
	}
}
