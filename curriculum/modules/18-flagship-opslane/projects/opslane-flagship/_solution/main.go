package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	_ "modernc.org/sqlite"
)

// Common errors.
var (
	ErrUnauthorized     = errors.New("unauthorized")
	ErrForbidden        = errors.New("forbidden")
	ErrNotFound         = errors.New("resource not found")
	ErrTenantMismatch   = errors.New("cross-tenant access prohibited")
	ErrInvalidState     = errors.New("invalid state transition")
	ErrDuplicateIdemKey = errors.New("idempotency key conflict")
	ErrQueueFull        = errors.New("worker queue full, backpressure applied")
)

// Role definitions.
type Role string

const (
	RoleAdmin    Role = "admin"
	RoleOperator Role = "operator"
	RoleViewer   Role = "viewer"
)

func (r Role) CanWrite() bool {
	return r == RoleAdmin || r == RoleOperator
}

// OrderStatus state machine.
type OrderStatus string

const (
	StatusPending    OrderStatus = "PENDING"
	StatusProcessing OrderStatus = "PROCESSING"
	StatusCompleted  OrderStatus = "COMPLETED"
	StatusFailed     OrderStatus = "FAILED"
	StatusCancelled  OrderStatus = "CANCELLED"
)

var allowedTransitions = map[OrderStatus][]OrderStatus{
	StatusPending:    {StatusProcessing, StatusCancelled},
	StatusProcessing: {StatusCompleted, StatusFailed},
	StatusFailed:     {StatusProcessing, StatusCancelled},
	StatusCompleted:  {},
	StatusCancelled:  {},
}

func isValidTransition(from, to OrderStatus) bool {
	for _, s := range allowedTransitions[from] {
		if s == to {
			return true
		}
	}
	return false
}

// TenantContext keys.
type contextKey string

const (
	tenantCtxKey contextKey = "tenant_id"
	roleCtxKey   contextKey = "role"
	reqIDCtxKey  contextKey = "request_id"
)

// Metrics counters.
type Metrics struct {
	TotalRequests   atomic.Uint64
	ActiveRequests  atomic.Int64
	OrdersCreated   atomic.Uint64
	PaymentsSuccess atomic.Uint64
	PaymentsFailed  atomic.Uint64
	QueueDropped    atomic.Uint64
}

// Config encapsulates server runtime settings.
type Config struct {
	Port           int
	Environment    string
	WorkerCount    int
	QueueCapacity  int
	RequestTimeout time.Duration
}

func DefaultConfig() *Config {
	return &Config{
		Port:           8080,
		Environment:    "development",
		WorkerCount:    4,
		QueueCapacity:  64,
		RequestTimeout: 10 * time.Second,
	}
}

// Order model.
type Order struct {
	ID          string      `json:"id"`
	TenantID    string      `json:"tenant_id"`
	Customer    string      `json:"customer"`
	AmountCents int64       `json:"amount_cents"`
	Status      OrderStatus `json:"status"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
}

// Payment model.
type Payment struct {
	ID             string    `json:"id"`
	TenantID       string    `json:"tenant_id"`
	OrderID        string    `json:"order_id"`
	IdempotencyKey string    `json:"idempotency_key"`
	AmountCents    int64     `json:"amount_cents"`
	Status         string    `json:"status"`
	CreatedAt      time.Time `json:"created_at"`
}

// WorkerPool manages background asynchronous processing.
type WorkerPool struct {
	tasks      chan func(ctx context.Context)
	wg         sync.WaitGroup
	workers    int
	isShutdown atomic.Bool
}

func NewWorkerPool(workers, queueCap int) *WorkerPool {
	wp := &WorkerPool{
		tasks:   make(chan func(context.Context), queueCap),
		workers: workers,
	}
	for i := 0; i < workers; i++ {
		wp.wg.Add(1)
		go func(workerID int) {
			defer wp.wg.Done()
			for task := range wp.tasks {
				task(context.Background())
			}
		}(i)
	}
	return wp
}

func (wp *WorkerPool) Submit(task func(ctx context.Context)) error {
	if wp.isShutdown.Load() {
		return errors.New("worker pool shutting down")
	}
	select {
	case wp.tasks <- task:
		return nil
	default:
		return ErrQueueFull
	}
}

func (wp *WorkerPool) Stop() {
	if wp.isShutdown.CompareAndSwap(false, true) {
		close(wp.tasks)
		wp.wg.Wait()
	}
}

// Cache with singleflight double-checking.
type Cache struct {
	mu    sync.RWMutex
	items map[string]cachedItem
}

type cachedItem struct {
	data      any
	expiresAt time.Time
}

func NewCache() *Cache {
	return &Cache{items: make(map[string]cachedItem)}
}

func (c *Cache) GetOrCompute(key string, ttl time.Duration, compute func() (any, error)) (any, error) {
	c.mu.RLock()
	item, found := c.items[key]
	c.mu.RUnlock()

	if found && time.Now().Before(item.expiresAt) {
		return item.data, nil
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	// Double-check after acquiring write lock
	if item, found := c.items[key]; found && time.Now().Before(item.expiresAt) {
		return item.data, nil
	}

	data, err := compute()
	if err != nil {
		return nil, err
	}

	c.items[key] = cachedItem{
		data:      data,
		expiresAt: time.Now().Add(ttl),
	}
	return data, nil
}

func (c *Cache) Invalidate(key string) {
	c.mu.Lock()
	delete(c.items, key)
	c.mu.Unlock()
}

// OpslaneServer is the unified production service.
type OpslaneServer struct {
	cfg     *Config
	db      *sql.DB
	pool    *WorkerPool
	cache   *Cache
	metrics *Metrics
	logger  *slog.Logger
	isReady atomic.Bool
}

func NewOpslaneServer(cfg *Config, db *sql.DB, logger *slog.Logger) (*OpslaneServer, error) {
	if logger == nil {
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	if err := initSchema(db); err != nil {
		return nil, fmt.Errorf("init schema: %w", err)
	}

	srv := &OpslaneServer{
		cfg:     cfg,
		db:      db,
		pool:    NewWorkerPool(cfg.WorkerCount, cfg.QueueCapacity),
		cache:   NewCache(),
		metrics: &Metrics{},
		logger:  logger,
	}
	srv.isReady.Store(true)
	return srv, nil
}

func initSchema(db *sql.DB) error {
	schema := `
	CREATE TABLE IF NOT EXISTS orders (
		id TEXT PRIMARY KEY,
		tenant_id TEXT NOT NULL,
		customer TEXT NOT NULL,
		amount_cents INTEGER NOT NULL,
		status TEXT NOT NULL,
		created_at TIMESTAMP NOT NULL,
		updated_at TIMESTAMP NOT NULL
	);
	CREATE INDEX IF NOT EXISTS idx_orders_tenant ON orders(tenant_id);

	CREATE TABLE IF NOT EXISTS payments (
		id TEXT PRIMARY KEY,
		tenant_id TEXT NOT NULL,
		order_id TEXT NOT NULL,
		idempotency_key TEXT UNIQUE NOT NULL,
		amount_cents INTEGER NOT NULL,
		status TEXT NOT NULL,
		created_at TIMESTAMP NOT NULL
	);
	CREATE INDEX IF NOT EXISTS idx_payments_tenant ON payments(tenant_id);

	CREATE TABLE IF NOT EXISTS audit_logs (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		tenant_id TEXT NOT NULL,
		action TEXT NOT NULL,
		entity_id TEXT NOT NULL,
		created_at TIMESTAMP NOT NULL
	);
	`
	_, err := db.Exec(schema)
	return err
}

// Core Domain Logic with Multi-Tenancy & Idempotency.

var (
	orderSeq   atomic.Uint64
	paymentSeq atomic.Uint64
)

func (s *OpslaneServer) CreateOrder(ctx context.Context, tenantID, customer string, amountCents int64) (*Order, error) {
	if tenantID == "" {
		return nil, ErrTenantMismatch
	}
	if customer == "" || amountCents <= 0 {
		return nil, errors.New("invalid order parameters: customer required and amount > 0")
	}

	now := time.Now().UTC()
	seq := orderSeq.Add(1)
	order := &Order{
		ID:          fmt.Sprintf("ord_%d_%d", now.UnixNano(), seq),
		TenantID:    tenantID,
		Customer:    customer,
		AmountCents: amountCents,
		Status:      StatusPending,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	query := `INSERT INTO orders (id, tenant_id, customer, amount_cents, status, created_at, updated_at)
	          VALUES (?, ?, ?, ?, ?, ?, ?)`
	_, err := s.db.ExecContext(ctx, query, order.ID, order.TenantID, order.Customer, order.AmountCents, string(order.Status), order.CreatedAt, order.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("insert order: %w", err)
	}

	s.metrics.OrdersCreated.Add(1)
	s.cache.Invalidate(fmt.Sprintf("%s:orders:all", tenantID))

	// Async event logging via worker pool
	_ = s.pool.Submit(func(c context.Context) {
		_, _ = s.db.Exec(`INSERT INTO audit_logs (tenant_id, action, entity_id, created_at) VALUES (?, ?, ?, ?)`,
			tenantID, "order.created", order.ID, time.Now().UTC())
	})

	return order, nil
}

func (s *OpslaneServer) GetOrder(ctx context.Context, tenantID, orderID string) (*Order, error) {
	cacheKey := fmt.Sprintf("%s:order:%s", tenantID, orderID)
	val, err := s.cache.GetOrCompute(cacheKey, 30*time.Second, func() (any, error) {
		query := `SELECT id, tenant_id, customer, amount_cents, status, created_at, updated_at
		          FROM orders WHERE id = ? AND tenant_id = ?`
		row := s.db.QueryRowContext(ctx, query, orderID, tenantID)

		var ord Order
		var statusStr string
		if err := row.Scan(&ord.ID, &ord.TenantID, &ord.Customer, &ord.AmountCents, &statusStr, &ord.CreatedAt, &ord.UpdatedAt); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, ErrNotFound
			}
			return nil, err
		}
		ord.Status = OrderStatus(statusStr)
		return &ord, nil
	})
	if err != nil {
		return nil, err
	}
	return val.(*Order), nil
}

func (s *OpslaneServer) ListOrders(ctx context.Context, tenantID string) ([]*Order, error) {
	cacheKey := fmt.Sprintf("%s:orders:all", tenantID)
	val, err := s.cache.GetOrCompute(cacheKey, 15*time.Second, func() (any, error) {
		query := `SELECT id, tenant_id, customer, amount_cents, status, created_at, updated_at
		          FROM orders WHERE tenant_id = ? ORDER BY created_at DESC LIMIT 100`
		rows, err := s.db.QueryContext(ctx, query, tenantID)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		var results []*Order
		for rows.Next() {
			var ord Order
			var statusStr string
			if err := rows.Scan(&ord.ID, &ord.TenantID, &ord.Customer, &ord.AmountCents, &statusStr, &ord.CreatedAt, &ord.UpdatedAt); err != nil {
				return nil, err
			}
			ord.Status = OrderStatus(statusStr)
			results = append(results, &ord)
		}
		return results, rows.Err()
	})
	if err != nil {
		return nil, err
	}
	return val.([]*Order), nil
}

func (s *OpslaneServer) ProcessPayment(ctx context.Context, tenantID, orderID, idemKey string) (*Payment, error) {
	if idemKey == "" {
		return nil, errors.New("idempotency_key is required")
	}

	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback()

	// 1. Check existing payment with idempotency key
	var existing Payment
	err = tx.QueryRowContext(ctx, `SELECT id, tenant_id, order_id, idempotency_key, amount_cents, status, created_at
	                               FROM payments WHERE idempotency_key = ?`, idemKey).
		Scan(&existing.ID, &existing.TenantID, &existing.OrderID, &existing.IdempotencyKey, &existing.AmountCents, &existing.Status, &existing.CreatedAt)

	if err == nil {
		// Found idempotent replay
		if existing.TenantID != tenantID {
			return nil, ErrTenantMismatch
		}
		return &existing, nil
	} else if !errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("query payment by idem key: %w", err)
	}

	// 2. Fetch order strictly within tenant boundary
	var ord Order
	var statusStr string
	err = tx.QueryRowContext(ctx, `SELECT id, tenant_id, customer, amount_cents, status, created_at, updated_at
	                               FROM orders WHERE id = ? AND tenant_id = ?`, orderID, tenantID).
		Scan(&ord.ID, &ord.TenantID, &ord.Customer, &ord.AmountCents, &statusStr, &ord.CreatedAt, &ord.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	ord.Status = OrderStatus(statusStr)

	// Invariant validation: must be transitionable to COMPLETED
	if ord.Status == StatusCompleted {
		return nil, errors.New("order is already completed")
	}
	if !isValidTransition(ord.Status, StatusProcessing) && ord.Status != StatusProcessing {
		return nil, fmt.Errorf("%w: cannot process payment for status %s", ErrInvalidState, ord.Status)
	}

	// 3. Insert payment record
	now := time.Now().UTC()
	pSeq := paymentSeq.Add(1)
	payment := &Payment{
		ID:             fmt.Sprintf("pay_%d_%d", now.UnixNano(), pSeq),
		TenantID:       tenantID,
		OrderID:        orderID,
		IdempotencyKey: idemKey,
		AmountCents:    ord.AmountCents,
		Status:         "SUCCESS",
		CreatedAt:      now,
	}

	_, err = tx.ExecContext(ctx, `INSERT INTO payments (id, tenant_id, order_id, idempotency_key, amount_cents, status, created_at)
	                             VALUES (?, ?, ?, ?, ?, ?, ?)`,
		payment.ID, payment.TenantID, payment.OrderID, payment.IdempotencyKey, payment.AmountCents, payment.Status, payment.CreatedAt)
	if err != nil {
		s.metrics.PaymentsFailed.Add(1)
		return nil, fmt.Errorf("insert payment: %w", err)
	}

	// 4. Update order status to COMPLETED
	_, err = tx.ExecContext(ctx, `UPDATE orders SET status = ?, updated_at = ? WHERE id = ? AND tenant_id = ?`,
		string(StatusCompleted), now, orderID, tenantID)
	if err != nil {
		s.metrics.PaymentsFailed.Add(1)
		return nil, fmt.Errorf("update order status: %w", err)
	}

	if err := tx.Commit(); err != nil {
		s.metrics.PaymentsFailed.Add(1)
		return nil, fmt.Errorf("commit payment tx: %w", err)
	}

	s.metrics.PaymentsSuccess.Add(1)
	s.cache.Invalidate(fmt.Sprintf("%s:order:%s", tenantID, orderID))
	s.cache.Invalidate(fmt.Sprintf("%s:orders:all", tenantID))

	// Submit background fulfillment event
	_ = s.pool.Submit(func(c context.Context) {
		_, _ = s.db.Exec(`INSERT INTO audit_logs (tenant_id, action, entity_id, created_at) VALUES (?, ?, ?, ?)`,
			tenantID, "payment.completed", payment.ID, time.Now().UTC())
	})

	return payment, nil
}

// HTTP Middleware & Handlers.

func (s *OpslaneServer) Routes() http.Handler {
	mux := http.NewServeMux()

	// Probes & Observability
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{"status": "UP"})
	})

	mux.HandleFunc("GET /readyz", func(w http.ResponseWriter, r *http.Request) {
		if !s.isReady.Load() {
			http.Error(w, `{"status":"NOT_READY"}`, http.StatusServiceUnavailable)
			return
		}
		if err := s.db.Ping(); err != nil {
			http.Error(w, `{"status":"DB_DOWN"}`, http.StatusServiceUnavailable)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{"status": "READY"})
	})

	mux.HandleFunc("GET /metrics", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		stats := map[string]any{
			"total_requests":   s.metrics.TotalRequests.Load(),
			"active_requests":  s.metrics.ActiveRequests.Load(),
			"orders_created":   s.metrics.OrdersCreated.Load(),
			"payments_success": s.metrics.PaymentsSuccess.Load(),
			"payments_failed":  s.metrics.PaymentsFailed.Load(),
			"queue_dropped":    s.metrics.QueueDropped.Load(),
		}
		_ = json.NewEncoder(w).Encode(stats)
	})

	// Authenticated API V1
	mux.HandleFunc("POST /api/v1/orders", s.handleCreateOrder)
	mux.HandleFunc("GET /api/v1/orders", s.handleListOrders)
	mux.HandleFunc("GET /api/v1/orders/{id}", s.handleGetOrder)
	mux.HandleFunc("POST /api/v1/orders/{id}/pay", s.handlePayOrder)

	return s.withMiddleware(mux)
}

func (s *OpslaneServer) withMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.metrics.TotalRequests.Add(1)
		s.metrics.ActiveRequests.Add(1)
		defer s.metrics.ActiveRequests.Add(-1)

		start := time.Now()
		reqID := r.Header.Get("X-Request-ID")
		if reqID == "" {
			reqID = fmt.Sprintf("req_%d", time.Now().UnixNano())
		}
		w.Header().Set("X-Request-ID", reqID)

		// Tenant & Role extraction from headers (simulating verified token claims)
		tenantID := r.Header.Get("X-Tenant-ID")
		roleStr := r.Header.Get("X-Role")
		if roleStr == "" {
			roleStr = string(RoleViewer)
		}

		ctx := context.WithValue(r.Context(), reqIDCtxKey, reqID)
		ctx = context.WithValue(ctx, tenantCtxKey, tenantID)
		ctx = context.WithValue(ctx, roleCtxKey, Role(roleStr))

		defer func() {
			if rec := recover(); rec != nil {
				s.logger.Error("panic recovered in HTTP handler", "panic", rec, "request_id", reqID)
				http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, r.WithContext(ctx))

		s.logger.Info("http request",
			"method", r.Method,
			"path", r.URL.Path,
			"tenant_id", tenantID,
			"request_id", reqID,
			"duration_ms", time.Since(start).Milliseconds(),
		)
	})
}

func (s *OpslaneServer) handleCreateOrder(w http.ResponseWriter, r *http.Request) {
	tenantID, _ := r.Context().Value(tenantCtxKey).(string)
	role, _ := r.Context().Value(roleCtxKey).(Role)

	if tenantID == "" {
		http.Error(w, `{"error":"missing X-Tenant-ID header"}`, http.StatusUnauthorized)
		return
	}
	if !role.CanWrite() {
		http.Error(w, `{"error":"insufficient permissions for order creation"}`, http.StatusForbidden)
		return
	}

	var req struct {
		Customer    string `json:"customer"`
		AmountCents int64  `json:"amount_cents"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid request json"}`, http.StatusBadRequest)
		return
	}

	order, err := s.CreateOrder(r.Context(), tenantID, req.Customer, req.AmountCents)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":%q}`, err.Error()), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(order)
}

func (s *OpslaneServer) handleListOrders(w http.ResponseWriter, r *http.Request) {
	tenantID, _ := r.Context().Value(tenantCtxKey).(string)
	if tenantID == "" {
		http.Error(w, `{"error":"missing X-Tenant-ID header"}`, http.StatusUnauthorized)
		return
	}

	orders, err := s.ListOrders(r.Context(), tenantID)
	if err != nil {
		http.Error(w, `{"error":"failed to retrieve orders"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(orders)
}

func (s *OpslaneServer) handleGetOrder(w http.ResponseWriter, r *http.Request) {
	tenantID, _ := r.Context().Value(tenantCtxKey).(string)
	orderID := r.PathValue("id")

	if tenantID == "" {
		http.Error(w, `{"error":"missing X-Tenant-ID header"}`, http.StatusUnauthorized)
		return
	}

	order, err := s.GetOrder(r.Context(), tenantID, orderID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			http.Error(w, `{"error":"order not found"}`, http.StatusNotFound)
			return
		}
		http.Error(w, fmt.Sprintf(`{"error":%q}`, err.Error()), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(order)
}

func (s *OpslaneServer) handlePayOrder(w http.ResponseWriter, r *http.Request) {
	tenantID, _ := r.Context().Value(tenantCtxKey).(string)
	role, _ := r.Context().Value(roleCtxKey).(Role)
	orderID := r.PathValue("id")

	if tenantID == "" {
		http.Error(w, `{"error":"missing X-Tenant-ID header"}`, http.StatusUnauthorized)
		return
	}
	if !role.CanWrite() {
		http.Error(w, `{"error":"insufficient permissions for payment"}`, http.StatusForbidden)
		return
	}

	var req struct {
		IdempotencyKey string `json:"idempotency_key"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || strings.TrimSpace(req.IdempotencyKey) == "" {
		http.Error(w, `{"error":"valid idempotency_key is required"}`, http.StatusBadRequest)
		return
	}

	payment, err := s.ProcessPayment(r.Context(), tenantID, orderID, req.IdempotencyKey)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			http.Error(w, `{"error":"order not found"}`, http.StatusNotFound)
			return
		}
		if errors.Is(err, ErrTenantMismatch) {
			http.Error(w, `{"error":"cross-tenant access prohibited"}`, http.StatusForbidden)
			return
		}
		http.Error(w, fmt.Sprintf(`{"error":%q}`, err.Error()), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(payment)
}

func (s *OpslaneServer) Close() {
	s.isReady.Store(false)
	s.pool.Stop()
	_ = s.db.Close()
}

func main() {
	cfg := DefaultConfig()
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))

	db, err := sql.Open("sqlite", "file:opslane_prod.db?cache=shared&mode=rwc")
	if err != nil {
		logger.Error("failed to open database", "err", err)
		os.Exit(1)
	}

	srv, err := NewOpslaneServer(cfg, db, logger)
	if err != nil {
		logger.Error("failed to initialize opslane server", "err", err)
		os.Exit(1)
	}
	defer srv.Close()

	httpServer := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Port),
		Handler:      srv.Routes(),
		ReadTimeout:  cfg.RequestTimeout,
		WriteTimeout: cfg.RequestTimeout,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		logger.Info("starting Opslane Flagship server", "port", cfg.Port, "env", cfg.Environment)
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("server fatal error", "err", err)
		}
	}()

	<-stop
	logger.Info("graceful shutdown initiated")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		logger.Error("forced server shutdown", "err", err)
	}
	logger.Info("Opslane Flagship stopped cleanly")
}
