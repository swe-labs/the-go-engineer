package main

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"strings"
	"testing"
	"time"
)

func TestLogRequest_TextHandler(t *testing.T) {
	tests := []struct {
		name   string
		method string
		path   string
		status int
		dur    time.Duration
	}{
		{"GET users", "GET", "/api/users", 200, 42 * time.Millisecond},
		{"POST order", "POST", "/api/orders", 201, 150 * time.Millisecond},
		{"DELETE error", "DELETE", "/api/orders/1", 500, 5200 * time.Millisecond},
		{"GET not found", "GET", "/api/unknown", 404, 1 * time.Millisecond},
		{"PUT update", "PUT", "/api/users/1", 200, 30 * time.Millisecond},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			logger := slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{Level: slog.LevelInfo}))
			LogRequest(logger, tt.method, tt.path, tt.status, tt.dur)
			output := buf.String()
			for _, key := range []string{"method", "path", "status", "duration_ms"} {
				if !strings.Contains(output, key) {
					t.Errorf("output missing key %q\noutput: %s", key, output)
				}
			}
		})
	}
}

func TestLogRequest_JSONHandler(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&buf, &slog.HandlerOptions{Level: slog.LevelInfo}))
	LogRequest(logger, "GET", "/api/test", 200, 50*time.Millisecond)

	var result map[string]any
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("failed to parse JSON log output: %v", err)
	}
	if result["method"] != "GET" {
		t.Errorf("expected method=GET, got %v", result["method"])
	}
	if result["path"] != "/api/test" {
		t.Errorf("expected path=/api/test, got %v", result["path"])
	}
}

func TestLogAtLevel(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&buf, &slog.HandlerOptions{Level: slog.LevelDebug}))

	levels := []struct {
		level slog.Level
		msg   string
	}{
		{slog.LevelDebug, "debug message"},
		{slog.LevelInfo, "info message"},
		{slog.LevelWarn, "warn message"},
		{slog.LevelError, "error message"},
	}

	for _, l := range levels {
		buf.Reset()
		LogAtLevel(logger, l.level, l.msg)

		var result map[string]any
		if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
			t.Fatalf("failed to parse JSON for level %v: %v", l.level, err)
		}
		if result["msg"] != l.msg {
			t.Errorf("expected msg=%q, got %v", l.msg, result["msg"])
		}
	}
}
