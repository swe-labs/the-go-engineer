package main

import (
	"bytes"
	"context"
	"log/slog"
	"strings"
	"testing"
)

func TestRedactHandler_RedactsSensitiveKeys(t *testing.T) {
	tests := []struct {
		name       string
		keys       []string
		logArgs    []slog.Attr
		wantRedact bool
	}{
		{
			name:       "redacts email",
			keys:       []string{"email"},
			logArgs:    []slog.Attr{slog.String("email", "alice@example.com")},
			wantRedact: true,
		},
		{
			name:       "redacts password",
			keys:       []string{"password"},
			logArgs:    []slog.Attr{slog.String("password", "s3cret!")},
			wantRedact: true,
		},
		{
			name:       "redacts ssn",
			keys:       []string{"ssn"},
			logArgs:    []slog.Attr{slog.String("ssn", "123-45-6789")},
			wantRedact: true,
		},
		{
			name:       "multiple sensitive fields",
			keys:       []string{"email", "password", "ssn", "credit_card"},
			logArgs:    []slog.Attr{slog.String("email", "bob@test.com"), slog.String("password", "pass123"), slog.String("ssn", "987-65-4321")},
			wantRedact: true,
		},
		{
			name:       "no sensitive fields in log",
			keys:       []string{"email", "password"},
			logArgs:    []slog.Attr{slog.String("user_id", "u-1"), slog.String("role", "admin")},
			wantRedact: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			baseHandler := slog.NewTextHandler(&buf, &slog.HandlerOptions{Level: slog.LevelInfo})
			handler := NewRedactHandler(baseHandler, tt.keys...)
			logger := slog.New(handler)
			logger.LogAttrs(context.Background(), slog.LevelInfo, "test", tt.logArgs...)
			output := buf.String()

			if tt.wantRedact && !strings.Contains(output, "[REDACTED]") {
				t.Errorf("expected [REDACTED] for sensitive keys\noutput: %s", output)
			}
			if !tt.wantRedact && strings.Contains(output, "[REDACTED]") {
				t.Errorf("expected no [REDACTED] when no sensitive keys present\noutput: %s", output)
			}
		})
	}
}

func TestRedactHandler_PassesThroughNonSensitiveKeys(t *testing.T) {
	var buf bytes.Buffer
	baseHandler := slog.NewTextHandler(&buf, &slog.HandlerOptions{Level: slog.LevelInfo})
	handler := NewRedactHandler(baseHandler, "email", "password")
	logger := slog.New(handler)
	logger.Info("test", "user_id", "u-42", "action", "login", "email", "test@test.com")
	output := buf.String()

	if !strings.Contains(output, "u-42") {
		t.Errorf("expected non-sensitive value u-42 to pass through\noutput: %s", output)
	}
	if !strings.Contains(output, "login") {
		t.Errorf("expected non-sensitive value login to pass through\noutput: %s", output)
	}
	if !strings.Contains(output, "[REDACTED]") {
		t.Errorf("expected email to be redacted\noutput: %s", output)
	}
}
