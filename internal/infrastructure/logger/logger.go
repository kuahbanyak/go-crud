package logger

import (
	"log"
	"os"
)

var (
	infoLogger  *log.Logger
	errorLogger *log.Logger
	debugLogger *log.Logger
)

func init() {
	infoLogger = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	errorLogger = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
	debugLogger = log.New(os.Stdout, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func Info(message string, args ...interface{}) {
	if len(args) > 0 {
		infoLogger.Printf(message+formatArgs(args), args...)
	} else {
		infoLogger.Println(message)
	}
}

func Error(message string, args ...interface{}) {
	if len(args) > 0 {
		errorLogger.Printf(message+formatArgs(args), args...)
	} else {
		errorLogger.Println(message)
	}
}

func Debug(message string, args ...interface{}) {
	if len(args) > 0 {
		debugLogger.Printf(message+formatArgs(args), args...)
	} else {
		debugLogger.Println(message)
	}
}

func formatArgs(args []interface{}) string {
	format := ""
	for i := 0; i < len(args); i += 2 {
		if i+1 < len(args) {
			format += " %s=%v"
		}
	}
	return format
}
