package main

import (
	"bytes"
	"context"
	"log/slog"
	"strings"
	"testing"
)

func TestWithAndGetCorrelationID(t *testing.T) {
	tests := []struct {
		name string
		cid  string
	}{
		{"standard id", "corr-abc-123"},
		{"empty id", ""},
		{"uuid style", "550e8400-e29b-41d4-a716-446655440000"},
		{"short id", "x"},
		{"long id", "a-very-long-correlation-id-that-spans-multiple-services-and-requests"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := WithCorrelationID(context.Background(), tt.cid)
			got := GetCorrelationID(ctx)
			if got != tt.cid {
				t.Errorf("GetCorrelationID = %q, want %q", got, tt.cid)
			}
		})
	}
}

func TestCallServices_ContainsCorrelationID(t *testing.T) {
	tests := []struct {
		name     string
		cid      string
		services []string
	}{
		{
			name:     "three services",
			cid:      "corr-xyz",
			services: []string{"auth", "db", "cache"},
		},
		{
			name:     "single service",
			cid:      "corr-001",
			services: []string{"api-gateway"},
		},
		{
			name:     "five services",
			cid:      "corr-999",
			services: []string{"edge", "auth", "orders", "payment", "notification"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			logger := slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{Level: slog.LevelInfo}))
			slog.SetDefault(logger)
			ctx := WithCorrelationID(context.Background(), tt.cid)
			chain := CallServices(ctx, tt.services)
			if len(chain) != len(tt.services) {
				t.Errorf("chain length = %d, want %d", len(chain), len(tt.services))
			}
			output := buf.String()
			if !strings.Contains(output, tt.cid) {
				t.Errorf("output missing correlation_id %q\noutput: %s", tt.cid, output)
			}
			for _, svc := range tt.services {
				if !strings.Contains(output, svc) {
					t.Errorf("output missing service %q\noutput: %s", svc, output)
				}
			}
		})
	}
}

func TestGetCorrelationID_EmptyWhenMissing(t *testing.T) {
	cid := GetCorrelationID(context.Background())
	if cid != "" {
		t.Errorf("expected empty string, got %q", cid)
	}
}
