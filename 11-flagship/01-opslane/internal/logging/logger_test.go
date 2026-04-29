// Copyright (c) 2026 Rasel Hossen
// See LICENSE for usage terms.

package logging

import (
	"bytes"
	"context"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGenerateCorrelationID(t *testing.T) {
	t.Parallel()

	id1 := GenerateCorrelationID()
	id2 := GenerateCorrelationID()

	if len(id1) != 32 {
		t.Fatalf("expected 32-char hex string, got %d chars: %q", len(id1), id1)
	}
	if id1 == id2 {
		t.Fatal("two generated IDs should not be equal")
	}
}

func TestCorrelationIDContextRoundTrip(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	if got := CorrelationID(ctx); got != "" {
		t.Fatalf("expected empty correlation ID, got %q", got)
	}

	ctx = WithCorrelationID(ctx, "test-abc-123")
	if got := CorrelationID(ctx); got != "test-abc-123" {
		t.Fatalf("expected %q, got %q", "test-abc-123", got)
	}
}

func TestCorrelationIDFromRequestExtractsHeader(t *testing.T) {
	t.Parallel()

	r := httptest.NewRequest("GET", "/", nil)
	r.Header.Set("X-Correlation-ID", "from-client")

	got := CorrelationIDFromRequest(r)
	if got != "from-client" {
		t.Fatalf("expected %q, got %q", "from-client", got)
	}
}

func TestCorrelationIDFromRequestGeneratesWhenMissing(t *testing.T) {
	t.Parallel()

	r := httptest.NewRequest("GET", "/", nil)
	got := CorrelationIDFromRequest(r)

	if got == "" {
		t.Fatal("expected generated correlation ID, got empty string")
	}
	if len(got) != 32 {
		t.Fatalf("expected 32-char hex ID, got %d chars: %q", len(got), got)
	}
}

func TestContextAttrsIncludesCorrelationID(t *testing.T) {
	t.Parallel()

	ctx := WithCorrelationID(context.Background(), "ctx-id-42")
	attrs := ContextAttrs(ctx)

	if len(attrs) != 1 {
		t.Fatalf("expected 1 attr, got %d", len(attrs))
	}
	if attrs[0].Key != "correlation_id" || attrs[0].Value.String() != "ctx-id-42" {
		t.Fatalf("unexpected attr: %v", attrs[0])
	}
}

func TestContextAttrsEmptyWithoutCorrelationID(t *testing.T) {
	t.Parallel()

	attrs := ContextAttrs(context.Background())
	if len(attrs) != 0 {
		t.Fatalf("expected 0 attrs, got %d", len(attrs))
	}
}

func TestNewLoggerJSON(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	logger := New(&buf, slog.LevelInfo, "json")
	logger.Info("test message", slog.String("key", "value"))

	output := buf.String()
	if !strings.Contains(output, `"msg":"test message"`) {
		t.Fatalf("expected JSON log output, got %q", output)
	}
}

func TestNewLoggerText(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	logger := New(&buf, slog.LevelInfo, "text")
	logger.Info("test message")

	output := buf.String()
	if !strings.Contains(output, "test message") {
		t.Fatalf("expected text log output, got %q", output)
	}
}

func TestWithContextEnrichesLogger(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	logger := New(&buf, slog.LevelInfo, "json")

	ctx := WithCorrelationID(context.Background(), "enrich-id-99")
	enriched := WithContext(ctx, logger)
	enriched.Info("enriched")

	output := buf.String()
	if !strings.Contains(output, "enrich-id-99") {
		t.Fatalf("expected correlation_id in output, got %q", output)
	}
}

func TestRequestLoggerGeneratesCorrelationID(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	logger := New(&buf, slog.LevelInfo, "json")

	handler := RequestLogger(logger)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify context has correlation ID.
		id := CorrelationID(r.Context())
		if id == "" {
			t.Error("expected correlation ID in context")
		}
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/test", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	// Check response header.
	if cid := rec.Header().Get("X-Correlation-ID"); cid == "" {
		t.Error("expected X-Correlation-ID response header")
	}

	// Check log output.
	output := buf.String()
	if !strings.Contains(output, "correlation_id") {
		t.Fatalf("expected correlation_id in log, got %q", output)
	}
	if !strings.Contains(output, `"status":200`) {
		t.Fatalf("expected status 200 in log, got %q", output)
	}
}

func TestRequestLoggerPreservesClientCorrelationID(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	logger := New(&buf, slog.LevelInfo, "json")

	handler := RequestLogger(logger)(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusAccepted)
	}))

	req := httptest.NewRequest("POST", "/submit", nil)
	req.Header.Set("X-Correlation-ID", "client-provided-id")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if got := rec.Header().Get("X-Correlation-ID"); got != "client-provided-id" {
		t.Fatalf("expected %q, got %q", "client-provided-id", got)
	}

	output := buf.String()
	if !strings.Contains(output, "client-provided-id") {
		t.Fatalf("expected client ID in log, got %q", output)
	}
}

func TestRequestLoggerCapturesStatus(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	logger := New(&buf, slog.LevelInfo, "json")

	handler := RequestLogger(logger)(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))

	req := httptest.NewRequest("GET", "/missing", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	output := buf.String()
	if !strings.Contains(output, `"status":404`) {
		t.Fatalf("expected status 404 in log, got %q", output)
	}
}
