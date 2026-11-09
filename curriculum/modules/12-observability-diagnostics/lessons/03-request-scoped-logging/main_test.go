package main

import (
	"bytes"
	"context"
	"log/slog"
	"strings"
	"testing"
)

func TestAttachLoggerAndLogWithCtx(t *testing.T) {
	tests := []struct {
		name      string
		requestID string
		wantKey   string
	}{
		{"first request", "req-001", "req-001"},
		{"second request", "req-002", "req-002"},
		{"error request", "req-err-99", "req-err-99"},
		{"empty request", "", ""},
		{"long request id", "req-abc-def-ghi-123456", "req-abc-def-ghi-123456"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			logger := slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{Level: slog.LevelInfo}))
			ctx := AttachLogger(context.Background(), logger)
			HandleRequest(ctx, tt.requestID)
			output := buf.String()
			if tt.wantKey != "" && !strings.Contains(output, tt.wantKey) {
				t.Errorf("output missing request_id %q\noutput: %s", tt.wantKey, output)
			}
			if !strings.Contains(output, "request_id") {
				t.Errorf("output missing request_id key\noutput: %s", output)
			}
		})
	}
}

func TestLoggerFromContext_ReturnsDefaultWhenMissing(t *testing.T) {
	logger := LoggerFromContext(context.Background())
	if logger == nil {
		t.Error("LoggerFromContext should never return nil")
	}
}

func TestLogWithCtx_UsesContextLogger(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{Level: slog.LevelInfo}))
	ctx := AttachLogger(context.Background(), logger)
	LogWithCtx(ctx, slog.LevelInfo, "test message")
	output := buf.String()
	if !strings.Contains(output, "test message") {
		t.Errorf("output missing message text\noutput: %s", output)
	}
}
