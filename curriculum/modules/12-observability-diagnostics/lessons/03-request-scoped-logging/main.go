package main

import (
	"context"
	"log/slog"
	"os"
)

type contextKey string

const loggerKey contextKey = "logger"

// AttachLogger stores a slog.Logger in context.
func AttachLogger(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}

// LoggerFromContext retrieves a slog.Logger from context or returns slog.Default().
func LoggerFromContext(ctx context.Context) *slog.Logger {
	if logger, ok := ctx.Value(loggerKey).(*slog.Logger); ok {
		return logger
	}
	return slog.Default()
}

// LogWithCtx logs a message using the logger from context.
func LogWithCtx(ctx context.Context, level slog.Level, msg string, attrs ...slog.Attr) {
	LoggerFromContext(ctx).LogAttrs(ctx, level, msg, attrs...)
}

// HandleRequest simulates a request handler with request-scoped logging.
func HandleRequest(ctx context.Context, requestID string) {
	requestLogger := LoggerFromContext(ctx).With("request_id", requestID)
	ctx = AttachLogger(ctx, requestLogger)
	LogWithCtx(ctx, slog.LevelInfo, "processing request")
	processInner(ctx)
}

func processInner(ctx context.Context) {
	LogWithCtx(ctx, slog.LevelInfo, "inner processing step")
}

func main() {
	baseLogger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	ctx := AttachLogger(context.Background(), baseLogger)
	HandleRequest(ctx, "req-001")
	HandleRequest(ctx, "req-002")
}
