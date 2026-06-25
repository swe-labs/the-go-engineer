package main

import (
	"strings"
	"testing"
)

func TestLogLevelString(t *testing.T) {
	tests := []struct {
		level LogLevel
		want  string
	}{
		{DEBUG, "DEBUG"},
		{INFO, "INFO"},
		{WARN, "WARN"},
		{ERROR, "ERROR"},
		{FATAL, "FATAL"},
	}

	for _, tc := range tests {
		t.Run(tc.want, func(t *testing.T) {
			if tc.level.String() != tc.want {
				t.Errorf("LogLevel(%d).String() = %q, want %q", tc.level, tc.level.String(), tc.want)
			}
		})
	}
}

func TestRedact(t *testing.T) {
	logger := NewSecureLogger(DEBUG)

	fields := map[string]interface{}{
		"user_id":      "u123",
		"password":     "s3cret!",
		"api_key":      "sk-live-abc123",
		"email":        "alice@example.com",
		"access_token": "eyJhbGciOiJIUzI1NiIs...",
	}

	result := logger.redact(fields)

	tests := []struct {
		key  string
		want string
	}{
		{key: "user_id", want: "u123"},
		{key: "password", want: "[REDACTED]"},
		{key: "api_key", want: "[REDACTED]"},
		{key: "email", want: "alice@example.com"},
		{key: "access_token", want: "[REDACTED]"},
	}

	for _, tc := range tests {
		t.Run(tc.key, func(t *testing.T) {
			got, ok := result[tc.key]
			if !ok {
				t.Errorf("key %q not found in result", tc.key)
				return
			}
			if got != tc.want {
				t.Errorf("result[%q] = %v, want %q", tc.key, got, tc.want)
			}
		})
	}
}

func TestRedactCaseInsensitive(t *testing.T) {
	logger := NewSecureLogger(DEBUG)

	fields := map[string]interface{}{
		"Password": "s3cret!",
		"SECRET":   "my-secret",
		"Token":    "abc123",
	}

	result := logger.redact(fields)

	for k, v := range result {
		if v != "[REDACTED]" {
			t.Errorf("expected %q to be redacted, got %v", k, v)
		}
	}
}

func TestLogLevelFiltering(t *testing.T) {
	logger := NewSecureLogger(WARN)

	// These should not produce output (not captured by test, just checking no panic)
	logger.Debug("debug message", nil)
	logger.Info("info message", nil)
	logger.Warn("warn message", nil)
	logger.Error("error message", nil)
}

func TestAuditTrail(t *testing.T) {
	trail := &AuditTrail{}
	logger := NewSecureLogger(INFO)
	logger.AuditEnabled = true
	logger.auditCallback = trail.Append

	logger.Warn("test warning", map[string]interface{}{"key": "value"})

	if len(trail.Entries) != 1 {
		t.Errorf("expected 1 audit entry, got %d", len(trail.Entries))
	}
}

func TestAuditTrailSummary(t *testing.T) {
	trail := &AuditTrail{}
	trail.Append(LogEntry{Message: "test"})

	summary := trail.Summary()
	if !strings.Contains(summary, "1 entries") {
		t.Errorf("unexpected summary: %q", summary)
	}
}
