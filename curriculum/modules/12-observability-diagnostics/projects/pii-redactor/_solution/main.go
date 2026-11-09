package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"regexp"
	"strings"
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
	h.stats.LogsProcessed.Add(1)

	sanitized := slog.NewRecord(r.Time, r.Level, h.maskString(r.Message), r.PC)
	r.Attrs(func(a slog.Attr) bool {
		sanitized.AddAttrs(h.sanitizeAttr(a))
		return true
	})

	return h.next.Handle(ctx, sanitized)
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
	// Sensitive key names are redacted outright
	keyLower := strings.ToLower(a.Key)
	switch keyLower {
	case "password", "token", "secret", "authorization", "api_key", "apikey":
		h.stats.PIIRedacted.Add(1)
		return slog.String(a.Key, "[REDACTED_SECRET]")
	}

	switch a.Value.Kind() {
	case slog.KindString:
		masked := h.maskString(a.Value.String())
		return slog.String(a.Key, masked)
	case slog.KindGroup:
		attrs := a.Value.Group()
		sanitizedGroup := make([]slog.Attr, len(attrs))
		for i, sub := range attrs {
			sanitizedGroup[i] = h.sanitizeAttr(sub)
		}
		return slog.Group(a.Key, anySliceToAttrs(sanitizedGroup)...)
	default:
		return a
	}
}

func anySliceToAttrs(attrs []slog.Attr) []any {
	res := make([]any, len(attrs))
	for i, a := range attrs {
		res[i] = a
	}
	return res
}

func (h *PIIRedactorHandler) maskString(text string) string {
	res := text

	// Mask Credit Cards (keep last 4 digits)
	res = ccRegex.ReplaceAllStringFunc(res, func(cc string) string {
		h.stats.PIIRedacted.Add(1)
		digitsOnly := strings.ReplaceAll(strings.ReplaceAll(cc, "-", ""), " ", "")
		if len(digitsOnly) >= 4 {
			return "****-****-****-" + digitsOnly[len(digitsOnly)-4:]
		}
		return "[REDACTED_CC]"
	})

	// Mask SSN
	if ssnRegex.MatchString(res) {
		res = ssnRegex.ReplaceAllStringFunc(res, func(ssn string) string {
			h.stats.PIIRedacted.Add(1)
			return "***-**-" + ssn[len(ssn)-4:]
		})
	}

	// Mask Email Addresses
	if emailRegex.MatchString(res) {
		res = emailRegex.ReplaceAllStringFunc(res, func(email string) string {
			h.stats.PIIRedacted.Add(1)
			parts := strings.Split(email, "@")
			if len(parts) == 2 && len(parts[0]) > 2 {
				return string(parts[0][0]) + "***" + string(parts[0][len(parts[0])-1]) + "@" + parts[1]
			}
			return "[REDACTED_EMAIL]"
		})
	}

	return res
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
