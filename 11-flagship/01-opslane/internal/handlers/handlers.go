// Copyright (c) 2026 Rasel Hossen
// See LICENSE for usage terms.

package handlers

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/netip"
	"sync/atomic"
	"time"

	"github.com/swe-labs/the-go-engineer/11-flagship/01-opslane/internal/auth"
	"github.com/swe-labs/the-go-engineer/11-flagship/01-opslane/internal/logging"
	"github.com/swe-labs/the-go-engineer/11-flagship/01-opslane/internal/metrics"
	"github.com/swe-labs/the-go-engineer/11-flagship/01-opslane/internal/middleware"
	"github.com/swe-labs/the-go-engineer/11-flagship/01-opslane/internal/models"
	"github.com/swe-labs/the-go-engineer/11-flagship/01-opslane/internal/otel"
	paymentflow "github.com/swe-labs/the-go-engineer/11-flagship/01-opslane/internal/payment"
	"github.com/swe-labs/the-go-engineer/11-flagship/01-opslane/internal/ratelimit"
	"github.com/swe-labs/the-go-engineer/11-flagship/01-opslane/internal/services"
)

const healthDatabaseTimeout = 2 * time.Second
const apiRateLimitWindow = time.Minute
const apiRateLimitMaxRequests = 120

// Application holds the dependencies for the HTTP server. It acts as the
// central composition root for all incoming web requests.
type Application struct {
	Logger            *slog.Logger
	Metrics           *metrics.AppMetrics
	Store             Store
	Orders            OrderWorkflow
	Payments          PaymentWorkflow
	Tokens            *auth.TokenManager
	ServiceName       string
	Environment       string
	TrustedProxyCIDRs []netip.Prefix
	IsDraining        *atomic.Bool
	RateLimiter       *ratelimit.Limiter
	Tracer            *otel.Tracer
}

// OrderWorkflow defines the subset of the Order Service required by HTTP handlers.
// Using narrow interfaces here makes testing handlers much easier.
type OrderWorkflow interface {
	CreateOrder(ctx context.Context, input services.CreateOrderInput) (services.CreateOrderResult, error)
}

// PaymentWorkflow defines the subset of the Payment Service required by handlers.
type PaymentWorkflow interface {
	ProcessPayment(ctx context.Context, job paymentflow.Job) (services.ProcessPaymentResult, error)
}

// Store defines the data access methods required by the HTTP handlers.
// By depending on this interface rather than the concrete db.Store, the handlers
// remain decoupled from PostgreSQL.
type Store interface {
	Ping(ctx context.Context) error
	CreateTenant(ctx context.Context, tenant *models.Tenant) error
	CreateUser(ctx context.Context, user *models.User) error
	GetUserByEmail(ctx context.Context, tenantID int64, email string) (models.User, error)
	CreateOrder(ctx context.Context, order *models.Order) error
	GetOrderByID(ctx context.Context, tenantID, orderID int64) (models.Order, error)
	GetOrderByIdempotencyKey(ctx context.Context, tenantID int64, idempotencyKey string) (models.Order, error)
	UpdateOrderStatus(ctx context.Context, tenantID, orderID int64, status models.OrderStatus) (models.Order, error)
	ListOrdersByTenant(ctx context.Context, tenantID int64) ([]models.Order, error)
	CreatePayment(ctx context.Context, payment *models.Payment) error
	GetPaymentByProviderReference(ctx context.Context, tenantID int64, providerReference string) (models.Payment, error)
	UpdatePaymentStatus(ctx context.Context, tenantID int64, providerReference string, status models.PaymentStatus, failureReason string) (models.Payment, error)
	ListPaymentsByOrder(ctx context.Context, tenantID, orderID int64) ([]models.Payment, error)
}

// Routes configures all the HTTP endpoints for the application, attaching
// necessary middleware (CORS, Rate Limiting, Authentication) to specific paths.
func (app *Application) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /", app.handleIndex)
	mux.HandleFunc("GET /health", app.handleHealth)
	mux.HandleFunc("GET /readyz", app.handleHealth)
	mux.HandleFunc("GET /livez", app.handleLiveness)
	mux.Handle("GET /metrics", metrics.PrometheusHandler(app.Metrics))
	mux.Handle("GET /me", auth.RequireAuth(app.Tokens)(http.HandlerFunc(app.handleMe)))
	mux.HandleFunc("POST /api/v1/tenants", app.handleCreateTenant)
	mux.HandleFunc("POST /api/v1/users", app.handleCreateUser)
	mux.HandleFunc("POST /api/v1/auth/login", app.handleLogin)

	protected := app.requireAPIAuth
	mux.Handle("GET /api/v1/me", protected(http.HandlerFunc(app.handleMe)))
	mux.Handle("GET /api/v1/orders", protected(http.HandlerFunc(app.handleListOrders)))
	mux.Handle("POST /api/v1/orders", protected(http.HandlerFunc(app.handleCreateOrder)))
	mux.Handle("POST /api/v1/payments", protected(http.HandlerFunc(app.handleCreatePayment)))
	mux.Handle("GET /api/v1/orders/{orderID}/payments", protected(http.HandlerFunc(app.handleListPaymentsByOrder)))

	rootHandler := middleware.RecoverPanic(app.Logger)(mux)
	if app.Tracer != nil && app.Tracer.Enabled() {
		rootHandler = otel.HTTPMiddleware(app.Tracer)(rootHandler)
	}
	baseHandler := logging.RequestLogger(app.Logger)(
		metrics.HTTPMetrics(app.Metrics)(rootHandler),
	)
	var rateLimitedHandler http.Handler
	if app.RateLimiter != nil {
		rateLimitedHandler = ratelimit.PerIPMiddleware(app.RateLimiter, app.TrustedProxyCIDRs, app.Logger)(baseHandler)
	} else {
		rateLimitedHandler = middleware.RateLimit(apiRateLimitMaxRequests, apiRateLimitWindow, app.TrustedProxyCIDRs)(baseHandler)
	}
	httpSurface := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/health", "/readyz", "/livez":
			baseHandler.ServeHTTP(w, r)
			return
		}

		rateLimitedHandler.ServeHTTP(w, r)
	})

	handler := middleware.SecureHeaders(
		middleware.CORS(
			httpSurface,
		),
	)

	return handler
}

func (app *Application) handleIndex(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{
		"service": app.ServiceName,
		"env":     app.Environment,
		"message": "Opslane PostgreSQL foundation is running",
	})
}

func (app *Application) handleHealth(w http.ResponseWriter, r *http.Request) {
	statusCode := http.StatusOK
	databaseStatus := "ok"
	serviceStatus := "ok"

	pingCtx, cancelPing := context.WithTimeout(r.Context(), healthDatabaseTimeout)
	defer cancelPing()

	if err := app.Store.Ping(pingCtx); err != nil {
		statusCode = http.StatusServiceUnavailable
		databaseStatus = "degraded"
		serviceStatus = "degraded"
	}

	if app.IsDraining != nil && app.IsDraining.Load() {
		statusCode = http.StatusServiceUnavailable
		serviceStatus = "draining"
	}

	writeJSON(w, statusCode, map[string]string{
		"service":  app.ServiceName,
		"env":      app.Environment,
		"status":   serviceStatus,
		"database": databaseStatus,
	})
}

func (app *Application) handleMe(w http.ResponseWriter, r *http.Request) {
	identity, err := auth.RequireIdentity(r.Context())
	if err != nil {
		http.Error(w, "missing authenticated identity", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"user_id":    identity.UserID,
		"tenant_id":  identity.TenantID,
		"email":      identity.Email,
		"role":       identity.Role,
		"expires_at": identity.ExpiresAt,
	})
}

func (app *Application) handleLiveness(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{
		"service": app.ServiceName,
		"env":     app.Environment,
		"status":  "alive",
	})
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}
