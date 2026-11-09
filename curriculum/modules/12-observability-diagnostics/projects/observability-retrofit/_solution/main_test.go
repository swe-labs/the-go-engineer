package main

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
)

func TestOrderService_MetricsAndTracing(t *testing.T) {
	reg := prometheus.NewRegistry()
	metrics := NewMetricsRegistry(reg)

	exp := tracetest.NewInMemoryExporter()
	tp := sdktrace.NewTracerProvider(sdktrace.WithSyncer(exp))
	tracer := tp.Tracer("test-tracer")

	var logBuf bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&logBuf, nil))

	svc := NewOrderService(metrics, tracer, logger)

	ctx := context.Background()

	// 1. Create Order Success
	order, err := svc.CreateOrder(ctx, CreateOrderRequest{
		Customer:    "Alice Corp",
		AmountCents: 4500,
	})
	if err != nil {
		t.Fatalf("unexpected error creating order: %v", err)
	}
	if order.ID == "" || order.Customer != "Alice Corp" || order.AmountCents != 4500 {
		t.Fatalf("unexpected order data: %+v", order)
	}

	// Verify trace spans recorded
	spans := exp.GetSpans()
	if len(spans) == 0 {
		t.Fatalf("expected at least 1 span recorded, got 0")
	}
	createSpan := spans[0]
	if createSpan.Name != "OrderService.CreateOrder" {
		t.Errorf("expected span name 'OrderService.CreateOrder', got %q", createSpan.Name)
	}

	// Verify span attributes
	hasIDAttr := false
	for _, attr := range createSpan.Attributes {
		if string(attr.Key) == "order.id" && attr.Value.AsString() == order.ID {
			hasIDAttr = true
			break
		}
	}
	if !hasIDAttr {
		t.Errorf("expected span to contain attribute order.id=%s", order.ID)
	}

	// Verify structured logging output contains order_id
	logStr := logBuf.String()
	if !strings.Contains(logStr, order.ID) || !strings.Contains(logStr, "Alice Corp") {
		t.Errorf("expected log to contain order ID and customer name, got:\n%s", logStr)
	}

	// 2. Query /metrics via HTTP handler
	handler := NewHTTPHandler(svc, reg)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK from /metrics, got %d", rec.Code)
	}
	metricsBody := rec.Body.String()

	if !strings.Contains(metricsBody, `order_requests_total{operation="create",status="success"} 1`) {
		t.Errorf("metrics missing create counter increment:\n%s", metricsBody)
	}
	if !strings.Contains(metricsBody, `order_duration_seconds_count{operation="create"} 1`) {
		t.Errorf("metrics missing duration histogram sample:\n%s", metricsBody)
	}
	if !strings.Contains(metricsBody, `active_orders_count 1`) {
		t.Errorf("metrics missing active_orders_count 1:\n%s", metricsBody)
	}

	// 3. Cancel order and verify active count decreases
	err = svc.CancelOrder(ctx, order.ID)
	if err != nil {
		t.Fatalf("unexpected error cancelling order: %v", err)
	}

	recCancelMetrics := httptest.NewRecorder()
	handler.ServeHTTP(recCancelMetrics, httptest.NewRequest(http.MethodGet, "/metrics", nil))
	cancelMetricsBody := recCancelMetrics.Body.String()

	if !strings.Contains(cancelMetricsBody, `active_orders_count 0`) {
		t.Errorf("expected active_orders_count to be 0 after cancel, got:\n%s", cancelMetricsBody)
	}
	if !strings.Contains(cancelMetricsBody, `order_requests_total{operation="cancel",status="success"} 1`) {
		t.Errorf("metrics missing cancel counter increment:\n%s", cancelMetricsBody)
	}
}

func TestOrderService_InvalidOrder_ErrorMetrics(t *testing.T) {
	reg := prometheus.NewRegistry()
	metrics := NewMetricsRegistry(reg)
	svc := NewOrderService(metrics, nil, nil)

	ctx := context.Background()

	_, err := svc.CreateOrder(ctx, CreateOrderRequest{
		Customer:    "",
		AmountCents: -10,
	})
	if err == nil {
		t.Fatalf("expected error on invalid order, got nil")
	}

	handler := NewHTTPHandler(svc, reg)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/metrics", nil))

	if !strings.Contains(rec.Body.String(), `order_requests_total{operation="create",status="error"} 1`) {
		t.Errorf("expected error counter increment in metrics:\n%s", rec.Body.String())
	}
}

func TestHTTPHandler_Endpoints(t *testing.T) {
	reg := prometheus.NewRegistry()
	metrics := NewMetricsRegistry(reg)
	svc := NewOrderService(metrics, nil, nil)
	handler := NewHTTPHandler(svc, reg)

	// POST /orders
	payload := `{"customer":"Bob Corp","amount_cents":1200}`
	req := httptest.NewRequest(http.MethodPost, "/orders", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected status 201 Created, got %d", rec.Code)
	}

	var created Order
	if err := json.Unmarshal(rec.Body.Bytes(), &created); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if created.ID == "" || created.Customer != "Bob Corp" {
		t.Fatalf("invalid created order: %+v", created)
	}

	// GET /orders/{id}
	reqGet := httptest.NewRequest(http.MethodGet, "/orders/"+created.ID, nil)
	recGet := httptest.NewRecorder()
	handler.ServeHTTP(recGet, reqGet)

	if recGet.Code != http.StatusOK {
		t.Fatalf("expected status 200 OK, got %d", recGet.Code)
	}

	// POST /orders/{id}/cancel
	reqCancel := httptest.NewRequest(http.MethodPost, "/orders/"+created.ID+"/cancel", nil)
	recCancel := httptest.NewRecorder()
	handler.ServeHTTP(recCancel, reqCancel)

	if recCancel.Code != http.StatusNoContent {
		t.Fatalf("expected status 204 No Content, got %d", recCancel.Code)
	}
}
