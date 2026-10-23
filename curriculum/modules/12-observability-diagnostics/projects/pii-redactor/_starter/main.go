package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"regexp"
	"sync/atomic"
)

var (
	emailRegex = regexp.MustCompile(`(?i)[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}`)
	ccRegex    = regexp.MustCompile(`\b(?:\d{4}[ -]?){3}\d{4}\b`)
	ssnRegex   = regexp.MustCompile(`\b\d{3}-\d{2}-\d{4}\b`)
)

type RedactionStats struct {
	LogsProcessed atomic.Int64
	PIIRedacted   atomic.Int64
}

// PIIRedactorHandler wraps an slog.Handler to sanitize sensitive attributes.
type PIIRedactorHandler struct {
	next  slog.Handler
	stats *RedactionStats
}

// NewPIIRedactorHandler creates a privacy-preserving slog.Handler.
func NewPIIRedactorHandler(next slog.Handler, stats *RedactionStats) *PIIRedactorHandler {
	if stats == nil {
		stats = &RedactionStats{}
	}
	return &PIIRedactorHandler{
		next:  next,
		stats: stats,
	}
}

func (h *PIIRedactorHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.next.Enabled(ctx, level)
}

func (h *PIIRedactorHandler) Handle(ctx context.Context, r slog.Record) error {
	// TODO: Increment h.stats.LogsProcessed
	// TODO: Create sanitized slog.Record with masked message
	// TODO: Iterate over attributes, sanitize each, add to sanitized record
	// TODO: Delegate to h.next.Handle(ctx, sanitized)
	return h.next.Handle(ctx, r)
}

func (h *PIIRedactorHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	sanitizedAttrs := make([]slog.Attr, len(attrs))
	for i, a := range attrs {
		sanitizedAttrs[i] = h.sanitizeAttr(a)
	}
	return &PIIRedactorHandler{
		next:  h.next.WithAttrs(sanitizedAttrs),
		stats: h.stats,
	}
}

func (h *PIIRedactorHandler) WithGroup(name string) slog.Handler {
	return &PIIRedactorHandler{
		next:  h.next.WithGroup(name),
		stats: h.stats,
	}
}

func (h *PIIRedactorHandler) sanitizeAttr(a slog.Attr) slog.Attr {
	// TODO: Check if a.Key is in sensitive list (password, token, secret, authorization, api_key)
	// TODO: If sensitive, replace value with "[REDACTED_SECRET]" and increment PIIRedacted
	// TODO: If a.Value is string, call maskString
	// TODO: If a.Value is group, recursively sanitize child attributes
	return a
}

func (h *PIIRedactorHandler) maskString(text string) string {
	// TODO: Mask credit cards (preserve last 4 digits)
	// TODO: Mask SSN (preserve last 4 digits)
	// TODO: Mask email address (j***e@domain.com)
	// TODO: Increment h.stats.PIIRedacted for each match
	return text
}

func main() {
	stats := &RedactionStats{}
	baseHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
	redactor := NewPIIRedactorHandler(baseHandler, stats)
	logger := slog.New(redactor)

	logger.Info("User registration initiated",
		"email", "alice.smith@corporate.org",
		"credit_card", "4111-2222-3333-4589",
		"ssn", "000-12-3456",
		"password", "verysecretpassword123",
	)

	fmt.Printf("\nRedaction Stats -> Logs Processed: %d, PII Redacted: %d\n",
		stats.LogsProcessed.Load(),
		stats.PIIRedacted.Load(),
	)
}
