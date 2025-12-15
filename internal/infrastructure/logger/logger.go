package logger

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"runtime"
	"time"
)

type LogLevel string

const (
	DEBUG LogLevel = "DEBUG"
	INFO  LogLevel = "INFO"
	WARN  LogLevel = "WARN"
	ERROR LogLevel = "ERROR"
)

var (
	infoLogger  *log.Logger
	errorLogger *log.Logger
	debugLogger *log.Logger
	warnLogger  *log.Logger
	logLevel    LogLevel
)

type LogEntry struct {
	Level     string                 `json:"level"`
	Message   string                 `json:"message"`
	Timestamp string                 `json:"timestamp"`
	RequestID string                 `json:"request_id,omitempty"`
	Fields    map[string]interface{} `json:"fields,omitempty"`
	File      string                 `json:"file,omitempty"`
	Line      int                    `json:"line,omitempty"`
}

func init() {
	infoLogger = log.New(os.Stdout, "", 0)
	errorLogger = log.New(os.Stderr, "", 0)
	debugLogger = log.New(os.Stdout, "", 0)
	warnLogger = log.New(os.Stdout, "", 0)

	// Set log level from environment
	level := os.Getenv("LOG_LEVEL")
	switch level {
	case "DEBUG":
		logLevel = DEBUG
	case "WARN":
		logLevel = WARN
	case "ERROR":
		logLevel = ERROR
	default:
		logLevel = INFO
	}
}

func Info(message string, args ...interface{}) {
	if logLevel > INFO {
		return
	}
	logStructured(INFO, message, nil, args...)
}

func InfoWithContext(ctx context.Context, message string, fields map[string]interface{}) {
	if logLevel > INFO {
		return
	}
	logStructured(INFO, message, ctx, fields)
}

func Error(message string, args ...interface{}) {
	logStructured(ERROR, message, nil, args...)
}

func ErrorWithContext(ctx context.Context, message string, fields map[string]interface{}) {
	logStructured(ERROR, message, ctx, fields)
}

func Debug(message string, args ...interface{}) {
	if logLevel > DEBUG {
		return
	}
	logStructured(DEBUG, message, nil, args...)
}

func DebugWithContext(ctx context.Context, message string, fields map[string]interface{}) {
	if logLevel > DEBUG {
		return
	}
	logStructured(DEBUG, message, ctx, fields)
}

func Warn(message string, args ...interface{}) {
	if logLevel > WARN {
		return
	}
	logStructured(WARN, message, nil, args...)
}

func WarnWithContext(ctx context.Context, message string, fields map[string]interface{}) {
	if logLevel > WARN {
		return
	}
	logStructured(WARN, message, ctx, fields)
}

func logStructured(level LogLevel, message string, ctx context.Context, args ...interface{}) {
	entry := LogEntry{
		Level:     string(level),
		Message:   message,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}

	// Get request ID from context
	if ctx != nil {
		if reqID := getRequestIDFromContext(ctx); reqID != "" {
			entry.RequestID = reqID
		}
	}

	// Parse fields from args
	if len(args) > 0 {
		if fields, ok := args[0].(map[string]interface{}); ok {
			entry.Fields = fields
		} else {
			// Legacy support: parse key-value pairs
			fields := make(map[string]interface{})
			for i := 0; i < len(args)-1; i += 2 {
				if key, ok := args[i].(string); ok {
					fields[key] = args[i+1]
				}
			}
			if len(fields) > 0 {
				entry.Fields = fields
			}
		}
	}

	// Add file and line info for errors and debug
	if level == ERROR || level == DEBUG {
		_, file, line, ok := runtime.Caller(2)
		if ok {
			entry.File = file
			entry.Line = line
		}
	}

	// JSON output for structured logging
	if isJSONLogging() {
		jsonBytes, _ := json.Marshal(entry)
		switch level {
		case ERROR:
			errorLogger.Println(string(jsonBytes))
		case WARN:
			warnLogger.Println(string(jsonBytes))
		case DEBUG:
			debugLogger.Println(string(jsonBytes))
		default:
			infoLogger.Println(string(jsonBytes))
		}
	} else {
		// Human-readable output for development
		logMessage := fmt.Sprintf("%s [%s] %s", entry.Timestamp, entry.Level, entry.Message)
		if entry.RequestID != "" {
			logMessage += fmt.Sprintf(" request_id=%s", entry.RequestID)
		}
		if entry.Fields != nil {
			for k, v := range entry.Fields {
				logMessage += fmt.Sprintf(" %s=%v", k, v)
			}
		}
		if entry.File != "" {
			logMessage += fmt.Sprintf(" (%s:%d)", entry.File, entry.Line)
		}

		switch level {
		case ERROR:
			errorLogger.Println(logMessage)
		case WARN:
			warnLogger.Println(logMessage)
		case DEBUG:
			debugLogger.Println(logMessage)
		default:
			infoLogger.Println(logMessage)
		}
	}
}

func isJSONLogging() bool {
	return os.Getenv("LOG_FORMAT") == "json"
}

func getRequestIDFromContext(ctx context.Context) string {
	type contextKey string
	const requestIDKey contextKey = "request_id"

	if reqID, ok := ctx.Value(requestIDKey).(string); ok {
		return reqID
	}
	return ""
}

// Legacy functions for backward compatibility
func formatArgs(args []interface{}) string {
	format := ""
	for i := 0; i < len(args); i += 2 {
		if i+1 < len(args) {
			format += " %s=%v"
		}
	}
	return format
}
