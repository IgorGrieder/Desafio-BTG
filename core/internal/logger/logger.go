package logger

import (
	"log/slog"
	"os"
)

var Logger *slog.Logger

// Init initializes the global structured logger with JSON output
func Init(env string) {
	var level slog.Level

	switch env {
	case "production", "prod":
		level = slog.LevelInfo
	case "development", "dev":
		level = slog.LevelDebug
	default:
		level = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{
		Level:     level,
		AddSource: true, // Add source file and line number
	}

	handler := slog.NewJSONHandler(os.Stdout, opts)
	Logger = slog.New(handler)

	// Set as default logger
	slog.SetDefault(Logger)
}

// Info logs an info level message
func Info(msg string, args ...any) {
	Logger.Info(msg, args...)
}

// Debug logs a debug level message
func Debug(msg string, args ...any) {
	Logger.Debug(msg, args...)
}

// Warn logs a warning level message
func Warn(msg string, args ...any) {
	Logger.Warn(msg, args...)
}

// Error logs an error level message
func Error(msg string, args ...any) {
	Logger.Error(msg, args...)
}

// Fatal logs an error and exits
func Fatal(msg string, args ...any) {
	Logger.Error(msg, args...)
	os.Exit(1)
}
