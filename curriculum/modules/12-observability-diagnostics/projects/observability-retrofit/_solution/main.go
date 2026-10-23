package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

var (
	ErrOrderNotFound = errors.New("order not found")
	ErrInvalidOrder  = errors.New("invalid order payload")
)

type Order struct {
	ID          string    `json:"id"`
	Customer    string    `json:"customer"`
	AmountCents int       `json:"amount_cents"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
}

type CreateOrderRequest struct {
	Customer    string `json:"customer"`
	AmountCents int    `json:"amount_cents"`
}

type MetricsRegistry struct {
	RequestsTotal   *prometheus.CounterVec
	RequestDuration *prometheus.HistogramVec
	ActiveOrders    prometheus.Gauge
}

func NewMetricsRegistry(reg prometheus.Registerer) *MetricsRegistry {
	if reg == nil {
		reg = prometheus.DefaultRegisterer
	}
	m := &MetricsRegistry{
		RequestsTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "order_requests_total",
				Help: "Total count of order operations",
			},
			[]string{"operation", "status"},
		),
		RequestDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "order_duration_seconds",
				Help:    "Execution duration of order operations in seconds",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"operation"},
		),
		ActiveOrders: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Name: "active_orders_count",
				Help: "Current number of active orders",
			},
		),
	}

	_ = reg.Register(m.RequestsTotal)
	_ = reg.Register(m.RequestDuration)
	_ = reg.Register(m.ActiveOrders)

	return m
}

func GenerateID() string {
	b := make([]byte, 8)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

type OrderService struct {
	mu      sync.RWMutex
	orders  map[string]Order
	metrics *MetricsRegistry
	tracer  trace.Tracer
	logger  *slog.Logger
}

func NewOrderService(metrics *MetricsRegistry, tracer trace.Tracer, logger *slog.Logger) *OrderService {
	if tracer == nil {
		tracer = otel.GetTracerProvider().Tracer("order-service")
	}
	if logger == nil {
		logger = slog.Default()
	}
	return &OrderService{
		orders:  make(map[string]Order),
		metrics: metrics,
		tracer:  tracer,
		logger:  logger,
	}
}

func (s *OrderService) CreateOrder(ctx context.Context, req CreateOrderRequest) (Order, error) {
	ctx, span := s.tracer.Start(ctx, "OrderService.CreateOrder")
	defer span.End()

	timer := prometheus.NewTimer(s.metrics.RequestDuration.WithLabelValues("create"))
	defer timer.ObserveDuration()

	if strings.TrimSpace(req.Customer) == "" || req.AmountCents <= 0 {
		s.metrics.RequestsTotal.WithLabelValues("create", "error").Inc()
		span.RecordError(ErrInvalidOrder)
		s.logger.WarnContext(ctx, "invalid order request",
			slog.String("customer", req.Customer),
			slog.Int("amount_cents", req.AmountCents),
		)
		return Order{}, ErrInvalidOrder
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	order := Order{
		ID:          GenerateID(),
		Customer:    strings.TrimSpace(req.Customer),
		AmountCents: req.AmountCents,
		Status:      "CREATED",
		CreatedAt:   time.Now().UTC(),
	}

	s.orders[order.ID] = order
	s.metrics.ActiveOrders.Inc()
	s.metrics.RequestsTotal.WithLabelValues("create", "success").Inc()

	span.SetAttributes(
		attribute.String("order.id", order.ID),
		attribute.String("order.customer", order.Customer),
		attribute.Int("order.amount_cents", order.AmountCents),
	)

	s.logger.InfoContext(ctx, "order created",
		slog.String("order_id", order.ID),
		slog.String("customer", order.Customer),
		slog.Int("amount_cents", order.AmountCents),
	)

	return order, nil
}

func (s *OrderService) GetOrder(ctx context.Context, id string) (Order, error) {
	ctx, span := s.tracer.Start(ctx, "OrderService.GetOrder")
	defer span.End()

	timer := prometheus.NewTimer(s.metrics.RequestDuration.WithLabelValues("get"))
	defer timer.ObserveDuration()

	span.SetAttributes(attribute.String("order.id", id))

	s.mu.RLock()
	defer s.mu.RUnlock()

	order, exists := s.orders[id]
	if !exists {
		s.metrics.RequestsTotal.WithLabelValues("get", "error").Inc()
		s.logger.WarnContext(ctx, "order not found", slog.String("order_id", id))
		return Order{}, ErrOrderNotFound
	}

	s.metrics.RequestsTotal.WithLabelValues("get", "success").Inc()
	return order, nil
}

func (s *OrderService) CancelOrder(ctx context.Context, id string) error {
	ctx, span := s.tracer.Start(ctx, "OrderService.CancelOrder")
	defer span.End()

	timer := prometheus.NewTimer(s.metrics.RequestDuration.WithLabelValues("cancel"))
	defer timer.ObserveDuration()

	span.SetAttributes(attribute.String("order.id", id))

	s.mu.Lock()
	defer s.mu.Unlock()

	order, exists := s.orders[id]
	if !exists {
		s.metrics.RequestsTotal.WithLabelValues("cancel", "error").Inc()
		s.logger.WarnContext(ctx, "cancel failed: order not found", slog.String("order_id", id))
		return ErrOrderNotFound
	}

	if order.Status != "CANCELLED" {
		order.Status = "CANCELLED"
		s.orders[id] = order
		s.metrics.ActiveOrders.Dec()
	}

	s.metrics.RequestsTotal.WithLabelValues("cancel", "success").Inc()
	s.logger.InfoContext(ctx, "order cancelled", slog.String("order_id", id))
	return nil
}

func NewHTTPHandler(svc *OrderService, gatherer prometheus.Gatherer) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /orders", func(w http.ResponseWriter, r *http.Request) {
		var req CreateOrderRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid json", http.StatusBadRequest)
			return
		}
		order, err := svc.CreateOrder(r.Context(), req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(order)
	})

	mux.HandleFunc("GET /orders/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		order, err := svc.GetOrder(r.Context(), id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(order)
	})

	mux.HandleFunc("POST /orders/{id}/cancel", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		err := svc.CancelOrder(r.Context(), id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	})

	if gatherer == nil {
		mux.Handle("GET /metrics", promhttp.Handler())
	} else {
		mux.Handle("GET /metrics", promhttp.HandlerFor(gatherer, promhttp.HandlerOpts{}))
	}

	return mux
}

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	metrics := NewMetricsRegistry(nil)
	svc := NewOrderService(metrics, nil, logger)
	handler := NewHTTPHandler(svc, nil)

	server := &http.Server{
		Addr:         ":8080",
		Handler:      handler,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	go func() {
		slog.Info("server running with observability retrofit", "addr", server.Addr)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("server error", "error", err)
		}
	}()

	<-shutdown
	slog.Info("shutting down cleanly")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = server.Shutdown(ctx)
}
