// Copyright (c) 2026 Rasel Hossen
// See LICENSE for usage terms.

package logging

import (
	"context"
	"io"
	"log/slog"
)

// New creates a structured logger. The format parameter selects the
// slog handler: "json" produces machine-readable output suitable for
// production log aggregation; "text" produces human-friendly output
// for local development.
//
// This factory centralises logger construction so the application has
// one place to control output format, level filtering, and handler
// options.
func New(w io.Writer, level slog.Level, format string) *slog.Logger {
	opts := &slog.HandlerOptions{Level: level}

	var handler slog.Handler
	switch format {
	case "json":
		handler = slog.NewJSONHandler(w, opts)
	default:
		handler = slog.NewTextHandler(w, opts)
	}

	return slog.New(handler)
}

// WithContext returns a new logger enriched with fields from the context.
// It reads the correlation ID from the logging context and, if present,
// any auth identity that the caller has attached.
//
// This bridges context-based request identity into structured log output
// without requiring every call site to extract the fields manually.
func WithContext(ctx context.Context, logger *slog.Logger) *slog.Logger {
	attrs := ContextAttrs(ctx)
	if len(attrs) == 0 {
		return logger
	}

	args := make([]any, 0, len(attrs))
	for _, a := range attrs {
		args = append(args, a)
	}
	return logger.With(args...)
}

// Info logs an info message enriched with context fields.
func Info(ctx context.Context, logger *slog.Logger, msg string, args ...any) {
	WithContext(ctx, logger).InfoContext(ctx, msg, args...)
}

// Error logs an error message enriched with context fields.
func Error(ctx context.Context, logger *slog.Logger, msg string, args ...any) {
	WithContext(ctx, logger).ErrorContext(ctx, msg, args...)
}

// Debug logs a debug message enriched with context fields.
func Debug(ctx context.Context, logger *slog.Logger, msg string, args ...any) {
	WithContext(ctx, logger).DebugContext(ctx, msg, args...)
}

// Warn logs a warning message enriched with context fields.
func Warn(ctx context.Context, logger *slog.Logger, msg string, args ...any) {
	WithContext(ctx, logger).WarnContext(ctx, msg, args...)
}
