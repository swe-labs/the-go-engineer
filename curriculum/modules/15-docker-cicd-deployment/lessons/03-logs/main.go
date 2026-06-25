package main

import (
	"io"
	"log/slog"
	"os"
	"time"
)

func newJSONLogger(w io.Writer, minLevel slog.Level) *slog.Logger {
	return slog.New(slog.NewJSONHandler(w, &slog.HandlerOptions{
		Level: minLevel,
	}))
}

func logRequest(logger *slog.Logger, method, path string, status int, duration time.Duration) {
	logger.Info("request",
		"method", method,
		"path", path,
		"status", status,
		"duration_ms", duration.Milliseconds(),
	)
}

func main() {
	logger := newJSONLogger(os.Stdout, slog.LevelInfo)
	logRequest(logger, "GET", "/api/users", 200, 15*time.Millisecond)
	logRequest(logger, "POST", "/api/orders", 201, 42*time.Millisecond)
	logger.Debug("this message is suppressed at Info level")
}
