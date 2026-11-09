package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func generateSampleBatch(size int) [][]byte {
	records := make([][]byte, size)
	services := []string{"auth-service", "billing-service", "inventory-service", "notification-service"}
	levels := []string{"INFO", "WARN", "ERROR"}

	for i := 0; i < size; i++ {
		svc := services[i%len(services)]
		lvl := levels[i%len(levels)]
		records[i] = []byte(`{"timestamp":"2026-10-23T10:00:00Z","level":"` + lvl + `","service":"` + svc + `","message":"processed transaction batch item"}`)
	}
	return records
}

func TestProcessBatch_Correctness(t *testing.T) {
	processor := NewEventProcessor()

	batch := [][]byte{
		[]byte(`{"timestamp":"2026-10-23T10:00:00Z","level":"INFO","service":"auth","message":"login ok"}`),
		[]byte(`{"timestamp":"2026-10-23T10:00:01Z","level":"ERROR","service":"auth","message":"database timeout"}`),
		[]byte(`{"timestamp":"2026-10-23T10:00:02Z","level":"INFO","service":"billing","message":"charge ok"}`),
		[]byte(`{"timestamp":"2026-10-23T10:00:03Z","level":"ERROR","service":"billing","message":"declined"}`),
		[]byte(`{"timestamp":"2026-10-23T10:00:04Z","level":"FATAL","service":"billing","message":"gateway unavailable"}`),
	}

	report, err := processor.ProcessBatch(batch)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if report.TotalEvents != 5 {
		t.Errorf("expected 5 total events, got %d", report.TotalEvents)
	}
	if report.TotalErrors != 3 {
		t.Errorf("expected 3 total errors, got %d", report.TotalErrors)
	}
	if len(report.Services) != 2 {
		t.Fatalf("expected 2 services, got %d", len(report.Services))
	}

	// Verify deterministic sorting: "auth" before "billing"
	if report.Services[0].Service != "auth" || report.Services[0].Total != 2 || report.Services[0].Errors != 1 {
		t.Errorf("unexpected auth stat: %+v", report.Services[0])
	}
	if report.Services[1].Service != "billing" || report.Services[1].Total != 3 || report.Services[1].Errors != 2 {
		t.Errorf("unexpected billing stat: %+v", report.Services[1])
	}

	// Verify digest summary
	expectedDigest := "TOTAL=5;ERRORS=3;SERVICES=auth:2,billing:3"
	if report.DigestSummary != expectedDigest {
		t.Errorf("expected digest %q, got %q", expectedDigest, report.DigestSummary)
	}

	// Verify JSON formatter
	var out bytes.Buffer
	if err := FormatReportJSON(&out, report); err != nil {
		t.Fatalf("failed to format json: %v", err)
	}
	if !strings.Contains(out.String(), `"total_events": 5`) {
		t.Errorf("missing total_events in json:\n%s", out.String())
	}
}

func TestProcessBatch_EdgeCases(t *testing.T) {
	processor := NewEventProcessor()

	// Empty batch
	_, err := processor.ProcessBatch(nil)
	if err != ErrEmptyBatch {
		t.Errorf("expected ErrEmptyBatch, got %v", err)
	}

	// Corrupt JSON
	_, err = processor.ProcessBatch([][]byte{[]byte(`{corrupt json`)})
	if err == nil || !strings.Contains(err.Error(), "corrupt log record") {
		t.Errorf("expected corrupt data error, got %v", err)
	}
}

func TestPprofServer_Endpoints(t *testing.T) {
	// Standard net/http/pprof handler mounts on DefaultServeMux
	req := httptest.NewRequest(http.MethodGet, "/debug/pprof/", nil)
	rec := httptest.NewRecorder()
	http.DefaultServeMux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK from /debug/pprof/, got %d", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "Types of profiles available") {
		t.Errorf("unexpected pprof index response: %s", rec.Body.String())
	}
}

func BenchmarkProcessBatch(b *testing.B) {
	processor := NewEventProcessor()
	batch := generateSampleBatch(100)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_, err := processor.ProcessBatch(batch)
		if err != nil {
			b.Fatalf("benchmark failed: %v", err)
		}
	}
}
