package main

import (
	"log/slog"
	"os"
	"time"
)

// NewTextLogger creates a structured text logger.
func NewTextLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
}

// NewJSONLogger creates a structured JSON logger.
func NewJSONLogger() *slog.Logger {
	return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
}

// LogAtLevel logs a message at the specified level with structured attributes.
func LogAtLevel(logger *slog.Logger, level slog.Level, msg string, attrs ...slog.Attr) {
	logger.LogAttrs(nil, level, msg, attrs...)
}

// LogRequest logs a structured HTTP request entry.
func LogRequest(logger *slog.Logger, method, path string, status int, dur time.Duration) {
	logger.Info("http request",
		"method", method,
		"path", path,
		"status", status,
		"duration_ms", dur.Milliseconds(),
	)
}

func main() {
	textLogger := NewTextLogger()
	jsonLogger := NewJSONLogger()

	slog.Info("=== Text logger ===")
	LogRequest(textLogger, "GET", "/api/users", 200, 42*time.Millisecond)
	LogAtLevel(textLogger, slog.LevelWarn, "slow query detected", slog.Int("query_ms", 3200))

	slog.Info("\n=== JSON logger ===")
	LogRequest(jsonLogger, "POST", "/api/orders", 201, 150*time.Millisecond)
	LogAtLevel(jsonLogger, slog.LevelError, "database connection failed", slog.String("db", "postgres"), slog.Int("retry", 3))
}
