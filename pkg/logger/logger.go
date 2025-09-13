package logger

import (
	"log"
	"os"
)

// Levels
const (
	DEBUG = "DEBUG"
	INFO  = "INFO"
	WARN  = "WARN"
	ERROR = "ERROR"
)

var (
	debugLogger = log.New(os.Stdout, "["+DEBUG+"] ", log.Ldate|log.Ltime|log.Lshortfile)
	infoLogger  = log.New(os.Stdout, "["+INFO+"] ", log.Ldate|log.Ltime|log.Lshortfile)
	warnLogger  = log.New(os.Stdout, "["+WARN+"] ", log.Ldate|log.Ltime|log.Lshortfile)
	errorLogger = log.New(os.Stderr, "["+ERROR+"] ", log.Ldate|log.Ltime|log.Lshortfile)
)

func Debug(format string, v ...interface{}) { debugLogger.Printf(format, v...) }
func Info(format string, v ...interface{})  { infoLogger.Printf(format, v...) }
func Warn(format string, v ...interface{})  { warnLogger.Printf(format, v...) }
func Error(format string, v ...interface{}) { errorLogger.Printf(format, v...) }
