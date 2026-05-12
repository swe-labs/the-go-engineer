// Copyright (c) 2026 Rasel Hossen
// See LICENSE for usage terms.

package logging

import (
	"context"
	"io"
	"log/slog"
)

// New (Function): creates a structured slog logger with json or text format
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

// WithContext (Function): returns a logger enriched with correlation and auth fields from context
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

// Info (Function): logs an info message enriched with context fields
func Info(ctx context.Context, logger *slog.Logger, msg string, args ...any) {
	WithContext(ctx, logger).InfoContext(ctx, msg, args...)
}

// Error (Function): logs an error message enriched with context fields
func Error(ctx context.Context, logger *slog.Logger, msg string, args ...any) {
	WithContext(ctx, logger).ErrorContext(ctx, msg, args...)
}

// Debug (Function): logs a debug message enriched with context fields
func Debug(ctx context.Context, logger *slog.Logger, msg string, args ...any) {
	WithContext(ctx, logger).DebugContext(ctx, msg, args...)
}

// Warn (Function): logs a warning message enriched with context fields
func Warn(ctx context.Context, logger *slog.Logger, msg string, args ...any) {
	WithContext(ctx, logger).WarnContext(ctx, msg, args...)
}
