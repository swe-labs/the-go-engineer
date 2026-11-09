package main

import (
	"fmt"
	"log/slog"
	"time"
)

// BlindLog prints request info without structured context.
func BlindLog(path string, status int, dur time.Duration) {
	fmt.Printf("request handled: %s %d %v\n", path, status, dur)
}

// ObservableLog logs a request with structured context.
func ObservableLog(path string, status int, dur time.Duration) {
	slog.Info("request handled",
		"path", path,
		"status_code", status,
		"duration_ms", dur.Milliseconds(),
	)
}

// LogRequest writes a structured log entry with request-scoped fields.
func LogRequest(logger *slog.Logger, handler, method, path string, statusCode int, latencyMs int64) {
	logger.Info("request completed",
		"handler", handler,
		"method", method,
		"path", path,
		"status_code", statusCode,
		"latency_ms", latencyMs,
	)
}

func main() {
	fmt.Println("=== BlindLog (unstructured) ===")
	BlindLog("/api/orders", 500, 5200*time.Millisecond)

	fmt.Println("\n=== ObservableLog (structured) ===")
	ObservableLog("/api/orders", 500, 5200*time.Millisecond)

	fmt.Println("\n=== LogRequest ===")
	LogRequest(slog.Default(), "OrderHandler", "POST", "/api/orders", 500, 5200)
}
