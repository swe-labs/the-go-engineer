package main

import (
	"bytes"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestObservability_MetricIncrementOnOrder(t *testing.T) {
	// Core practice task requirement:
	// Add custom Prometheus metric to track order processing duration,
	// instrument the order handler, and write a test verifying that the metric increments correctly.
	var logBuf bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&logBuf, nil))
	metrics := NewMetricsRegistry()
	server := NewObservabilityServer(metrics, logger)
	router := server.Routes()

	// 1. Process an order for tenant 'tenant-globex'
	reqProcess := httptest.NewRequest(http.MethodPost, "/api/v1/orders/process", nil)
	reqProcess.Header.Set("X-Tenant-ID", "tenant-globex")
	rrProcess := httptest.NewRecorder()
	router.ServeHTTP(rrProcess, reqProcess)

	if rrProcess.Code != http.StatusOK {
		t.Fatalf("expected 200 OK from order processing, got %d", rrProcess.Code)
	}

	// 2. Scrape the /metrics endpoint
	reqMetrics := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	rrMetrics := httptest.NewRecorder()
	router.ServeHTTP(rrMetrics, reqMetrics)

	if rrMetrics.Code != http.StatusOK {
		t.Fatalf("expected 200 OK from /metrics, got %d", rrMetrics.Code)
	}

	metricsOutput := rrMetrics.Body.String()

	// 3. Verify Counter metric incremented with matching label
	expectedCounter := `opslane_orders_processed_total{status="success",tenant_id="tenant-globex"} 1`
	if !strings.Contains(metricsOutput, expectedCounter) {
		t.Errorf("expected metrics output to contain counter %q, got:\n%s", expectedCounter, metricsOutput)
	}

	// 4. Verify Histogram metric duration recorded
	expectedHistogram := `opslane_order_processing_duration_seconds_count{status="success",tenant_id="tenant-globex"} 1`
	if !strings.Contains(metricsOutput, expectedHistogram) {
		t.Errorf("expected metrics output to contain histogram count %q", expectedHistogram)
	}
}

func TestObservability_CorrelationIDLogging(t *testing.T) {
	var logBuf bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&logBuf, nil))
	metrics := NewMetricsRegistry()
	server := NewObservabilityServer(metrics, logger)
	router := server.Routes()

	customCorrID := "corr-test-trace-9999"
	req := httptest.NewRequest(http.MethodPost, "/api/v1/orders/process", nil)
	req.Header.Set("X-Correlation-ID", customCorrID)
	req.Header.Set("X-Tenant-ID", "tenant-acme")
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	// Verify Header echoed
	if headerVal := rr.Header().Get("X-Correlation-ID"); headerVal != customCorrID {
		t.Errorf("expected X-Correlation-ID header %q, got %q", customCorrID, headerVal)
	}

	// Verify structured log entry contains correlation ID
	logOutput := logBuf.String()
	if !strings.Contains(logOutput, customCorrID) {
		t.Errorf("expected log output to include correlation ID %q, got: %s", customCorrID, logOutput)
	}
	if !strings.Contains(logOutput, "tenant-acme") {
		t.Errorf("expected log output to include tenant ID 'tenant-acme', got: %s", logOutput)
	}
}
