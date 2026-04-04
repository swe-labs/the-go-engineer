// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"sync"
)

// ============================================================================
// Section 23: Structured Logging — Custom slog.Handler
// Level: Advanced
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - The slog.Handler interface: the 4 methods you must implement
//   - Building a pretty-print terminal handler for development
//   - Thread-safety requirements: the mutex pattern for concurrent writes
//   - How to fan out to multiple handlers simultaneously
//
// ENGINEERING DEPTH:
//   slog.Handler is the extension point for the entire logging ecosystem.
//   Every third-party logging backend (Datadog, Sentry, OpenTelemetry) plugs
//   in via this interface. By implementing it yourself you understand exactly
//   why "structured logging" means what it means.
//
//   THE INTERFACE:
//     type Handler interface {
//         Enabled(ctx context.Context, level Level) bool
//         Handle(ctx context.Context, r Record) error
//         WithAttrs(attrs []Attr) Handler
//         WithGroup(name string) Handler
//     }
//
//   Enabled() is a fast pre-check. If it returns false, slog skips all attribute
//   evaluation — this is why debug logging has near-zero cost when disabled.
//
// RUN: go run ./23-structured-logging/3-custom-handler
// ============================================================================

// ============================================================================
// PrettyHandler — colorised terminal output for local development
// ============================================================================
// In production you use JSONHandler. In development, this handler renders logs
// as human-readable colored lines:
//   10:32:11 INFO  server started addr=:8080 env=dev
//   10:32:12 WARN  high memory  mb=3814

const (
	colorReset  = "\033[0m"
	colorGray   = "\033[90m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorRed    = "\033[31m"
	colorBlue   = "\033[34m"
)

type PrettyHandler struct {
	mu    sync.Mutex // Protects w — multiple goroutines may log concurrently
	out   *os.File
	level slog.Level
	attrs []slog.Attr // Pre-loaded attrs from .With()
	group string      // Active group prefix from .WithGroup()
}

func NewPrettyHandler(out *os.File, level slog.Level) *PrettyHandler {
	return &PrettyHandler{out: out, level: level}
}

// Enabled is the fast path. slog calls this before evaluating any attributes.
// Return false to skip the record entirely at zero cost.
func (h *PrettyHandler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.level
}

// Handle formats and writes a single log record.
// The Record contains: time, level, message, and the attribute list.
func (h *PrettyHandler) Handle(_ context.Context, r slog.Record) error {
	var buf bytes.Buffer

	// Timestamp (muted)
	buf.WriteString(colorGray)
	buf.WriteString(r.Time.Format("15:04:05"))
	buf.WriteString(colorReset)
	buf.WriteByte(' ')

	// Level (colored)
	buf.WriteString(levelColor(r.Level))
	buf.WriteString(fmt.Sprintf("%-5s", r.Level.String()))
	buf.WriteString(colorReset)
	buf.WriteByte(' ')

	// Message
	buf.WriteString(r.Message)

	// Pre-loaded attrs from .With() calls
	for _, a := range h.attrs {
		buf.WriteByte(' ')
		writeAttr(&buf, h.group, a)
	}

	// Per-record attrs
	r.Attrs(func(a slog.Attr) bool {
		buf.WriteByte(' ')
		writeAttr(&buf, h.group, a)
		return true
	})

	buf.WriteByte('\n')

	h.mu.Lock()
	defer h.mu.Unlock()
	_, err := h.out.Write(buf.Bytes())
	return err
}

func writeAttr(buf *bytes.Buffer, group string, a slog.Attr) {
	key := a.Key
	if group != "" {
		key = group + "." + key
	}
	buf.WriteString(colorBlue)
	buf.WriteString(key)
	buf.WriteString(colorReset)
	buf.WriteByte('=')
	buf.WriteString(fmt.Sprintf("%v", a.Value.Any()))
}

func levelColor(l slog.Level) string {
	switch {
	case l >= slog.LevelError:
		return colorRed
	case l >= slog.LevelWarn:
		return colorYellow
	case l >= slog.LevelInfo:
		return colorGreen
	default:
		return colorGray
	}
}

// WithAttrs returns a new handler with the given attrs pre-loaded.
// IMPORTANT: never mutate h.attrs. Always allocate a new slice.
func (h *PrettyHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	newAttrs := make([]slog.Attr, len(h.attrs)+len(attrs))
	copy(newAttrs, h.attrs)
	copy(newAttrs[len(h.attrs):], attrs)
	return &PrettyHandler{out: h.out, level: h.level, attrs: newAttrs, group: h.group}
}

// WithGroup returns a new handler that prefixes all subsequent attr keys.
func (h *PrettyHandler) WithGroup(name string) slog.Handler {
	group := name
	if h.group != "" {
		group = h.group + "." + name
	}
	return &PrettyHandler{out: h.out, level: h.level, attrs: h.attrs, group: group}
}

// ============================================================================
// MultiHandler — fan-out to N handlers simultaneously
// ============================================================================
// Send DEBUG to a file, INFO+ to stdout, ERROR+ to Sentry — all from one logger.

type MultiHandler struct {
	handlers []slog.Handler
}

func NewMultiHandler(handlers ...slog.Handler) *MultiHandler {
	return &MultiHandler{handlers: handlers}
}

func (m *MultiHandler) Enabled(ctx context.Context, level slog.Level) bool {
	for _, h := range m.handlers {
		if h.Enabled(ctx, level) {
			return true
		}
	}
	return false
}

func (m *MultiHandler) Handle(ctx context.Context, r slog.Record) error {
	for _, h := range m.handlers {
		if h.Enabled(ctx, r.Level) {
			if err := h.Handle(ctx, r.Clone()); err != nil {
				return err
			}
		}
	}
	return nil
}

func (m *MultiHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	handlers := make([]slog.Handler, len(m.handlers))
	for i, h := range m.handlers {
		handlers[i] = h.WithAttrs(attrs)
	}
	return &MultiHandler{handlers: handlers}
}

func (m *MultiHandler) WithGroup(name string) slog.Handler {
	handlers := make([]slog.Handler, len(m.handlers))
	for i, h := range m.handlers {
		handlers[i] = h.WithGroup(name)
	}
	return &MultiHandler{handlers: handlers}
}

// ============================================================================
// ErrorOnlyHandler — captures only errors for alert pipelines
// ============================================================================

type ErrorOnlyHandler struct {
	Alerts []map[string]any // In production: send to PagerDuty / Sentry
	slog.Handler
}

func (e *ErrorOnlyHandler) Handle(ctx context.Context, r slog.Record) error {
	if r.Level >= slog.LevelError {
		alert := map[string]any{
			"msg":   r.Message,
			"level": r.Level.String(),
			"time":  r.Time,
		}
		r.Attrs(func(a slog.Attr) bool {
			alert[a.Key] = a.Value.Any()
			return true
		})
		e.Alerts = append(e.Alerts, alert)
		// In production: sentry.CaptureMessage(r.Message) or pagerduty.Trigger(...)
	}
	return nil
}

func main() {
	// =========================================================================
	// Demo 1: PrettyHandler for local development
	// =========================================================================
	pretty := slog.New(NewPrettyHandler(os.Stdout, slog.LevelDebug))
	pretty.Debug("cache miss", slog.String("key", "user:42"))
	pretty.Info("server started", slog.String("addr", ":8080"))
	pretty.Warn("high memory", slog.Int("mb", 3814))
	pretty.Error("db timeout", slog.Duration("elapsed", 0))

	// Using .With() and .WithGroup()
	reqLog := pretty.With(slog.String("request_id", "req_001"))
	reqLog.Info("request completed",
		slog.Group("response", slog.Int("status", 200), slog.Duration("latency", 0)))

	fmt.Println()

	// =========================================================================
	// Demo 2: MultiHandler — stdout + error capture
	// =========================================================================
	errCapture := &ErrorOnlyHandler{Handler: slog.DiscardHandler}
	jsonHandler := slog.NewJSONHandler(os.Stdout, nil)
	multi := slog.New(NewMultiHandler(jsonHandler, errCapture))

	multi.Info("order placed")
	multi.Error("payment failed", slog.String("reason", "insufficient funds"))
	multi.Info("retry succeeded")

	fmt.Printf("\nCaptured %d alert(s):\n", len(errCapture.Alerts))
	for _, a := range errCapture.Alerts {
		b, _ := json.Marshal(a)
		fmt.Println(" ", string(b))
	}

	// KEY TAKEAWAY:
	// - Implement 4 methods: Enabled, Handle, WithAttrs, WithGroup
	// - Enabled() is the hot path — return false quickly when below min level
	// - Never mutate h.attrs in WithAttrs — always allocate a new slice
	// - MultiHandler fans records out to N backends simultaneously
	// - This is how every slog-compatible library (Datadog, Sentry) is built
}
