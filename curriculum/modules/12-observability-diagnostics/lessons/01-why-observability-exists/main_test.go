package main

import (
	"bytes"
	"log/slog"
	"strings"
	"testing"
	"time"
)

func TestLogRequest_ContainsAllKeys(t *testing.T) {
	tests := []struct {
		name       string
		handler    string
		method     string
		path       string
		statusCode int
		latencyMs  int64
		wantKeys   []string
	}{
		{
			name:       "successful GET",
			handler:    "UserHandler",
			method:     "GET",
			path:       "/api/users",
			statusCode: 200,
			latencyMs:  42,
			wantKeys:   []string{"handler", "method", "path", "status_code", "latency_ms"},
		},
		{
			name:       "server error POST",
			handler:    "OrderHandler",
			method:     "POST",
			path:       "/api/orders",
			statusCode: 500,
			latencyMs:  5200,
			wantKeys:   []string{"handler", "method", "path", "status_code", "latency_ms"},
		},
		{
			name:       "not found GET",
			handler:    "NotFoundHandler",
			method:     "GET",
			path:       "/api/unknown",
			statusCode: 404,
			latencyMs:  1,
			wantKeys:   []string{"handler", "method", "path", "status_code", "latency_ms"},
		},
		{
			name:       "redirect GET",
			handler:    "RedirectHandler",
			method:     "GET",
			path:       "/api/old",
			statusCode: 301,
			latencyMs:  5,
			wantKeys:   []string{"handler", "method", "path", "status_code", "latency_ms"},
		},
		{
			name:       "bad request POST",
			handler:    "ValidationHandler",
			method:     "POST",
			path:       "/api/validate",
			statusCode: 400,
			latencyMs:  15,
			wantKeys:   []string{"handler", "method", "path", "status_code", "latency_ms"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			logger := slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{Level: slog.LevelInfo}))
			LogRequest(logger, tt.handler, tt.method, tt.path, tt.statusCode, tt.latencyMs)
			output := buf.String()
			for _, key := range tt.wantKeys {
				if !strings.Contains(output, key) {
					t.Errorf("LogRequest output missing key %q\noutput: %s", key, output)
				}
			}
			if !strings.Contains(output, tt.handler) {
				t.Errorf("LogRequest output missing handler value %q\noutput: %s", tt.handler, output)
			}
		})
	}
}

func TestObservableLog_Structured(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{Level: slog.LevelInfo}))
	slog.SetDefault(logger)
	ObservableLog("/test", 200, 100*time.Millisecond)
	output := buf.String()
	for _, key := range []string{"path", "status_code", "duration_ms"} {
		if !strings.Contains(output, key) {
			t.Errorf("ObservableLog output missing key %q\noutput: %s", key, output)
		}
	}
}
