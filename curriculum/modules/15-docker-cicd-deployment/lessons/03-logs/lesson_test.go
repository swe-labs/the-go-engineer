package main

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"strings"
	"testing"
	"time"
)

func TestNewJSONLogger(t *testing.T) {
	var buf bytes.Buffer
	logger := newJSONLogger(&buf, slog.LevelInfo)

	logger.Info("test message", "key", "value")

	var result map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}

	if result["msg"] != "test message" {
		t.Errorf("got msg=%v, want 'test message'", result["msg"])
	}
	if result["level"] != "INFO" {
		t.Errorf("got level=%v, want 'INFO'", result["level"])
	}
	if result["key"] != "value" {
		t.Errorf("got key=%v, want 'value'", result["key"])
	}
}

func TestLogRequest(t *testing.T) {
	var buf bytes.Buffer
	logger := newJSONLogger(&buf, slog.LevelInfo)
	logRequest(logger, "POST", "/api/test", 201, 30*time.Millisecond)

	var result map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}

	if result["method"] != "POST" {
		t.Errorf("got method=%v, want 'POST'", result["method"])
	}
	if result["status"] != float64(201) {
		t.Errorf("got status=%v, want 201", result["status"])
	}
}

func TestDebugSuppressed(t *testing.T) {
	var buf bytes.Buffer
	logger := newJSONLogger(&buf, slog.LevelInfo)
	logger.Debug("should not appear")

	if buf.Len() > 0 {
		t.Errorf("expected no output for debug at Info level, got: %s", buf.String())
	}
}

func TestDebugEnabled(t *testing.T) {
	var buf bytes.Buffer
	logger := newJSONLogger(&buf, slog.LevelDebug)
	logger.Debug("debug message")

	if !strings.Contains(buf.String(), "debug message") {
		t.Errorf("expected debug message in output, got: %s", buf.String())
	}
}
