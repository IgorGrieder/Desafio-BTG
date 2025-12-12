package logger

import (
	"go.uber.org/zap"
)

var Logger *zap.Logger

// Init initializes the global structured logger with JSON output using Zap
func Init(env string) {
	var config zap.Config

	switch env {
	case "production", "prod":
		config = zap.NewProductionConfig()
	case "development", "dev":
		config = zap.NewDevelopmentConfig()
		config.Encoding = "json"
	default:
		config = zap.NewProductionConfig()
	}

	// Enable caller (file and line number)
	config.DisableCaller = false
	config.DisableStacktrace = false

	var err error
	Logger, err = config.Build()
	if err != nil {
		panic("Failed to initialize zap logger: " + err.Error())
	}
}

// Info logs an info level message
func Info(msg string, fields ...zap.Field) {
	Logger.Info(msg, fields...)
}

// Debug logs a debug level message
func Debug(msg string, fields ...zap.Field) {
	Logger.Debug(msg, fields...)
}

// Warn logs a warning level message
func Warn(msg string, fields ...zap.Field) {
	Logger.Warn(msg, fields...)
}

// Error logs an error level message
func Error(msg string, fields ...zap.Field) {
	Logger.Error(msg, fields...)
}

// Fatal logs an error and exits
func Fatal(msg string, fields ...zap.Field) {
	Logger.Fatal(msg, fields...)
}

// Sync flushes any buffered log entries
func Sync() {
	_ = Logger.Sync()
}
