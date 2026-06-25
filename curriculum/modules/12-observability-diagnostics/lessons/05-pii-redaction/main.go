package main

import (
	"context"
	"log/slog"
	"os"
	"strings"
)

// RedactHandler wraps an slog.Handler and redacts values for sensitive keys.
type RedactHandler struct {
	inner  slog.Handler
	redact map[string]bool
}

// NewRedactHandler creates a handler that redacts specified keys.
func NewRedactHandler(inner slog.Handler, sensitiveKeys ...string) *RedactHandler {
	r := make(map[string]bool, len(sensitiveKeys))
	for _, k := range sensitiveKeys {
		r[strings.ToLower(k)] = true
	}
	return &RedactHandler{inner: inner, redact: r}
}

func (h *RedactHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.inner.Enabled(ctx, level)
}

func (h *RedactHandler) Handle(ctx context.Context, rec slog.Record) error {
	var attrs []slog.Attr
	rec.Attrs(func(a slog.Attr) bool {
		if h.redact[strings.ToLower(a.Key)] {
			a.Value = slog.StringValue("[REDACTED]")
		}
		attrs = append(attrs, a)
		return true
	})
	newRec := slog.NewRecord(rec.Time, rec.Level, rec.Message, rec.PC)
	for _, a := range attrs {
		newRec.AddAttrs(a)
	}
	return h.inner.Handle(ctx, newRec)
}

func (h *RedactHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &RedactHandler{inner: h.inner.WithAttrs(attrs), redact: h.redact}
}

func (h *RedactHandler) WithGroup(name string) slog.Handler {
	return &RedactHandler{inner: h.inner.WithGroup(name), redact: h.redact}
}

// LogUser logs user data with PII automatically redacted.
func LogUser(logger *slog.Logger, userID, email, name, role string) {
	logger.Info("user operation",
		"user_id", userID,
		"email", email,
		"name", name,
		"role", role,
	)
}

func main() {
	baseHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
	redactHandler := NewRedactHandler(baseHandler, "email", "password", "ssn", "credit_card")
	logger := slog.New(redactHandler)

	LogUser(logger, "u-42", "alice@example.com", "Alice Smith", "admin")

	slog.Info("note: email value above should appear as [REDACTED]")
}
