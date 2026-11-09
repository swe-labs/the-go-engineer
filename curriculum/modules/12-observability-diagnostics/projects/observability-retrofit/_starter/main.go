package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
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

// TODO: Task 2 - Add Prometheus Metrics Registry
// Define and register:
//   - order_requests_total (CounterVec: operation, status)
//   - order_duration_seconds (HistogramVec: operation)
//   - active_orders_count (Gauge)
type MetricsRegistry struct {
	RequestsTotal   *prometheus.CounterVec
	RequestDuration *prometheus.HistogramVec
	ActiveOrders    prometheus.Gauge
}

func NewMetricsRegistry(reg prometheus.Registerer) *MetricsRegistry {
	// TODO: Initialize Prometheus metrics with appropriate labels and register with reg
	return &MetricsRegistry{}
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
	// TODO: Task 3 - Start OpenTelemetry span: OrderService.CreateOrder
	// TODO: Task 2 - Start Prometheus duration timer for operation="create"

	if strings.TrimSpace(req.Customer) == "" || req.AmountCents <= 0 {
		// TODO: Task 1 - Replace log.Println with slog.WarnContext
		// TODO: Task 2 - Increment order_requests_total counter with status="error"
		log.Println("invalid order payload:", req.Customer, req.AmountCents)
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

	// TODO: Task 1 - Log structured order creation with slog.InfoContext
	// TODO: Task 2 - Increment active_orders_count and order_requests_total (status="success")
	// TODO: Task 3 - Set span attributes (order.id, order.customer, order.amount_cents)
	log.Println("order created:", order.ID)

	return order, nil
}

func (s *OrderService) GetOrder(ctx context.Context, id string) (Order, error) {
	// TODO: Task 3 - Start OpenTelemetry span: OrderService.GetOrder
	// TODO: Task 2 - Start Prometheus duration timer for operation="get"

	s.mu.RLock()
	defer s.mu.RUnlock()

	order, exists := s.orders[id]
	if !exists {
		// TODO: Task 1 - Log warn with slog
		// TODO: Task 2 - Increment get error counter
		return Order{}, ErrOrderNotFound
	}

	// TODO: Task 2 - Increment get success counter
	return order, nil
}

func (s *OrderService) CancelOrder(ctx context.Context, id string) error {
	// TODO: Task 3 - Start OpenTelemetry span: OrderService.CancelOrder
	// TODO: Task 2 - Start Prometheus duration timer for operation="cancel"

	s.mu.Lock()
	defer s.mu.Unlock()

	order, exists := s.orders[id]
	if !exists {
		return ErrOrderNotFound
	}

	if order.Status != "CANCELLED" {
		order.Status = "CANCELLED"
		s.orders[id] = order
		// TODO: Task 2 - Decrement active_orders_count gauge
	}

	// TODO: Task 2 - Increment cancel success counter
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

	// TODO: Task 2 - Expose /metrics using promhttp.HandlerFor(gatherer, ...)

	return mux
}

func main() {
	metrics := NewMetricsRegistry(nil)
	svc := NewOrderService(metrics, nil, nil)
	handler := NewHTTPHandler(svc, nil)

	server := &http.Server{
		Addr:    ":8080",
		Handler: handler,
	}

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Println("legacy server running on :8080 without full observability")
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("server error: %v", err)
		}
	}()

	<-shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = server.Shutdown(ctx)
}
