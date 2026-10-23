package main

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type contextKey string

const correlationIDKey = contextKey("correlation_id")

// CorrelationIDFromContext extracts the correlation ID or returns empty string.
func CorrelationIDFromContext(ctx context.Context) string {
	if val, ok := ctx.Value(correlationIDKey).(string); ok {
		return val
	}
	return ""
}

// ContextWithCorrelationID injects a correlation ID into context.
func ContextWithCorrelationID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, correlationIDKey, id)
}

// MetricsRegistry encapsulates Prometheus metrics for Opslane observability.
type MetricsRegistry struct {
	Registry                 *prometheus.Registry
	OrdersProcessedTotal     *prometheus.CounterVec
	OrderProcessingDuration  *prometheus.HistogramVec
	ActiveWorkersGauge       prometheus.Gauge
}

func NewMetricsRegistry() *MetricsRegistry {
	reg := prometheus.NewRegistry()

	ordersCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "opslane_orders_processed_total",
			Help: "Total count of processed orders by tenant and status.",
		},
		[]string{"tenant_id", "status"},
	)

	orderDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "opslane_order_processing_duration_seconds",
			Help:    "Histogram of order processing execution durations in seconds.",
			Buckets: []float64{0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1.0, 2.5},
		},
		[]string{"tenant_id", "status"},
	)

	activeWorkers := prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "opslane_active_workers",
			Help: "Current number of actively executing background worker goroutines.",
		},
	)

	reg.MustRegister(ordersCounter)
	reg.MustRegister(orderDuration)
	reg.MustRegister(activeWorkers)

	return &MetricsRegistry{
		Registry:                reg,
		OrdersProcessedTotal:    ordersCounter,
		OrderProcessingDuration: orderDuration,
		ActiveWorkersGauge:      activeWorkers,
	}
}

// ObservabilityServer combines metrics and structured logging.
type ObservabilityServer struct {
	metrics *MetricsRegistry
	logger  *slog.Logger
}

func NewObservabilityServer(metrics *MetricsRegistry, logger *slog.Logger) *ObservabilityServer {
	return &ObservabilityServer{
		metrics: metrics,
		logger:  logger,
	}
}

func (s *ObservabilityServer) Routes() http.Handler {
	mux := http.NewServeMux()

	// Metrics scrape endpoint
	mux.Handle("GET /metrics", promhttp.HandlerFor(s.metrics.Registry, promhttp.HandlerOpts{}))

	// Business endpoint instrumented with logs and metrics
	mux.HandleFunc("POST /api/v1/orders/process", s.handleProcessOrder)

	return s.correlationMiddleware(mux)
}

func (s *ObservabilityServer) correlationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		corrID := r.Header.Get("X-Correlation-ID")
		if corrID == "" {
			b := make([]byte, 8)
			_, _ = rand.Read(b)
			corrID = hex.EncodeToString(b)
		}

		ctx := ContextWithCorrelationID(r.Context(), corrID)
		w.Header().Set("X-Correlation-ID", corrID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (s *ObservabilityServer) handleProcessOrder(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	corrID := CorrelationIDFromContext(r.Context())
	tenantID := r.Header.Get("X-Tenant-ID")
	if tenantID == "" {
		tenantID = "default"
	}

	// Simulate work
	time.Sleep(10 * time.Millisecond)

	duration := time.Since(start).Seconds()

	// 1. Record Prometheus metrics
	s.metrics.OrdersProcessedTotal.WithLabelValues(tenantID, "success").Inc()
	s.metrics.OrderProcessingDuration.WithLabelValues(tenantID, "success").Observe(duration)

	// 2. Structured slog log with correlation ID
	s.logger.Info("Order processed successfully",
		"correlation_id", corrID,
		"tenant_id", tenantID,
		"duration_sec", duration,
	)

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"status":"success"}`))
}

func main() {
	var logBuf bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&logBuf, nil))
	metrics := NewMetricsRegistry()
	server := NewObservabilityServer(metrics, logger)

	fmt.Println("Opslane Observability server listening on :8080...")
	_ = http.ListenAndServe(":8080", server.Routes())
}
