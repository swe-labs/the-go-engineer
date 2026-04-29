package handlers

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/rasel9t6/the-go-engineer/11-flagship/01-opslane/internal/auth"
	"github.com/rasel9t6/the-go-engineer/11-flagship/01-opslane/internal/db"
	"github.com/rasel9t6/the-go-engineer/11-flagship/01-opslane/internal/models"
	paymentflow "github.com/rasel9t6/the-go-engineer/11-flagship/01-opslane/internal/payment"
	"github.com/rasel9t6/the-go-engineer/11-flagship/01-opslane/internal/services"
)

func TestProtectedMeRouteReturnsAuthenticatedTenantIdentity(t *testing.T) {
	t.Parallel()

	tokens := newHandlerTestTokenManager(t)
	token, err := tokens.Issue(auth.Identity{
		UserID:   42,
		TenantID: 7,
		Email:    "admin@example.com",
		Role:     models.UserRoleAdmin,
	})
	if err != nil {
		t.Fatalf("Issue returned error: %v", err)
	}

	app := &Application{
		Logger:      slog.Default(),
		Tokens:      tokens,
		ServiceName: "opslane",
		Environment: "test",
	}

	req := httptest.NewRequest(http.MethodGet, "/me", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	res := httptest.NewRecorder()

	app.Routes().ServeHTTP(res, req)

	if res.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusOK)
	}

	var payload map[string]any
	if err := json.NewDecoder(res.Body).Decode(&payload); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if payload["tenant_id"] != float64(7) {
		t.Fatalf("tenant_id = %v, want 7", payload["tenant_id"])
	}
}

func TestProtectedMeRouteRejectsAnonymousRequest(t *testing.T) {
	t.Parallel()

	app := &Application{
		Logger:      slog.Default(),
		Tokens:      newHandlerTestTokenManager(t),
		ServiceName: "opslane",
		Environment: "test",
	}

	res := httptest.NewRecorder()
	app.Routes().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/me", nil))

	if res.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusUnauthorized)
	}
}

func TestCreateTenantReturnsCreatedTenant(t *testing.T) {
	t.Parallel()

	store := &fakeStore{}
	app := newTestApplication(t, store)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/tenants", bytes.NewBufferString(`{
		"name": "Acme Inc",
		"slug": "Acme"
	}`))
	res := httptest.NewRecorder()

	app.Routes().ServeHTTP(res, req)

	if res.Code != http.StatusCreated {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusCreated)
	}

	if store.createdTenant.Slug != "acme" {
		t.Fatalf("tenant slug = %q, want acme", store.createdTenant.Slug)
	}
}

func TestCreateTenantReturnsConflictForDuplicateSlug(t *testing.T) {
	t.Parallel()

	store := &fakeStore{createTenantErr: db.ErrDuplicateValue}
	app := newTestApplication(t, store)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/tenants", bytes.NewBufferString(`{
		"name": "Acme Inc",
		"slug": "acme"
	}`))
	res := httptest.NewRecorder()

	app.Routes().ServeHTTP(res, req)

	if res.Code != http.StatusConflict {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusConflict)
	}

	var payload map[string]map[string]string
	if err := json.NewDecoder(res.Body).Decode(&payload); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if payload["error"]["code"] != "duplicate_slug" {
		t.Fatalf("error code = %q, want duplicate_slug", payload["error"]["code"])
	}
}

func TestCreateUserHashesPasswordBeforePersisting(t *testing.T) {
	t.Parallel()

	store := &fakeStore{}
	app := newTestApplication(t, store)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewBufferString(`{
		"tenant_id": 7,
		"email": "Admin@Example.com",
		"display_name": "Admin User",
		"password": "CorrectHorse7Battery",
		"role": "admin"
	}`))
	res := httptest.NewRecorder()

	app.Routes().ServeHTTP(res, req)

	if res.Code != http.StatusCreated {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusCreated)
	}

	if store.createdUser.Email != "admin@example.com" {
		t.Fatalf("email = %q, want normalized email", store.createdUser.Email)
	}

	if store.createdUser.PasswordHash == "CorrectHorse7Battery" {
		t.Fatal("password must be hashed before persistence")
	}
}

func TestCreateUserReturnsNotFoundForUnknownTenant(t *testing.T) {
	t.Parallel()

	store := &fakeStore{createUserErr: db.ErrInvalidReference}
	app := newTestApplication(t, store)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewBufferString(`{
		"tenant_id": 999,
		"email": "admin@example.com",
		"display_name": "Admin User",
		"password": "CorrectHorse7Battery",
		"role": "admin"
	}`))
	res := httptest.NewRecorder()

	app.Routes().ServeHTTP(res, req)

	if res.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusNotFound)
	}

	var payload map[string]map[string]string
	if err := json.NewDecoder(res.Body).Decode(&payload); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if payload["error"]["code"] != "tenant_not_found" {
		t.Fatalf("error code = %q, want tenant_not_found", payload["error"]["code"])
	}
}

func TestCreateUserReturnsConflictForDuplicateEmail(t *testing.T) {
	t.Parallel()

	store := &fakeStore{createUserErr: db.ErrDuplicateValue}
	app := newTestApplication(t, store)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewBufferString(`{
		"tenant_id": 7,
		"email": "admin@example.com",
		"display_name": "Admin User",
		"password": "CorrectHorse7Battery",
		"role": "admin"
	}`))
	res := httptest.NewRecorder()

	app.Routes().ServeHTTP(res, req)

	if res.Code != http.StatusConflict {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusConflict)
	}

	var payload map[string]map[string]string
	if err := json.NewDecoder(res.Body).Decode(&payload); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if payload["error"]["code"] != "duplicate_email" {
		t.Fatalf("error code = %q, want duplicate_email", payload["error"]["code"])
	}
}

func TestCreateUserValidatesRequiredFieldsBeforePasswordHashing(t *testing.T) {
	t.Parallel()

	store := &fakeStore{}
	app := newTestApplication(t, store)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewBufferString(`{
		"tenant_id": 0,
		"email": "",
		"display_name": "",
		"password": "short",
		"role": "admin"
	}`))
	res := httptest.NewRecorder()

	app.Routes().ServeHTTP(res, req)

	if res.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusBadRequest)
	}

	var payload map[string]map[string]string
	if err := json.NewDecoder(res.Body).Decode(&payload); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if payload["error"]["code"] != "invalid_user" {
		t.Fatalf("error code = %q, want invalid_user", payload["error"]["code"])
	}
}

func TestLoginReturnsAccessToken(t *testing.T) {
	t.Parallel()

	passwordHash, err := auth.HashPassword("CorrectHorse7Battery")
	if err != nil {
		t.Fatalf("HashPassword returned error: %v", err)
	}

	store := &fakeStore{
		userByEmail: models.User{
			ID:           42,
			TenantID:     7,
			Email:        "admin@example.com",
			DisplayName:  "Admin User",
			PasswordHash: passwordHash,
			Role:         models.UserRoleAdmin,
		},
	}
	app := newTestApplication(t, store)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewBufferString(`{
		"tenant_id": 7,
		"email": "admin@example.com",
		"password": "CorrectHorse7Battery"
	}`))
	res := httptest.NewRecorder()

	app.Routes().ServeHTTP(res, req)

	if res.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusOK)
	}

	var payload map[string]any
	if err := json.NewDecoder(res.Body).Decode(&payload); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if payload["access_token"] == "" {
		t.Fatal("access_token should not be empty")
	}
}

func TestCreateOrderUsesAuthenticatedTenantAndUser(t *testing.T) {
	t.Parallel()

	store := &fakeStore{}
	app := newTestApplication(t, store)
	token := issueHandlerTestToken(t, app.Tokens, auth.Identity{
		UserID:   42,
		TenantID: 7,
		Email:    "admin@example.com",
		Role:     models.UserRoleAdmin,
	})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/orders", bytes.NewBufferString(`{
		"total_cents": 2500,
		"currency": "usd",
		"idempotency_key": "checkout-123"
	}`))
	req.Header.Set("Authorization", "Bearer "+token)
	res := httptest.NewRecorder()

	app.Routes().ServeHTTP(res, req)

	if res.Code != http.StatusCreated {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusCreated)
	}

	if store.createdOrder.TenantID != 7 || store.createdOrder.UserID != 42 {
		t.Fatalf("order identity = tenant %d user %d, want tenant 7 user 42", store.createdOrder.TenantID, store.createdOrder.UserID)
	}

	if store.createdOrder.Currency != "USD" {
		t.Fatalf("currency = %q, want USD", store.createdOrder.Currency)
	}
}

func TestCreateOrderReturnsExistingOrderForIdempotentRetry(t *testing.T) {
	t.Parallel()

	store := &fakeStore{
		orderByIdempotencyKey: models.Order{
			ID:             101,
			TenantID:       7,
			UserID:         42,
			Status:         models.OrderStatusPending,
			TotalCents:     2500,
			Currency:       "USD",
			IdempotencyKey: "checkout-123",
			CreatedAt:      time.Now().UTC(),
			UpdatedAt:      time.Now().UTC(),
		},
	}
	app := newTestApplication(t, store)
	token := issueHandlerTestToken(t, app.Tokens, auth.Identity{
		UserID:   42,
		TenantID: 7,
		Email:    "admin@example.com",
		Role:     models.UserRoleAdmin,
	})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/orders", bytes.NewBufferString(`{
		"total_cents": 2500,
		"currency": "usd",
		"idempotency_key": "checkout-123"
	}`))
	req.Header.Set("Authorization", "Bearer "+token)
	res := httptest.NewRecorder()

	app.Routes().ServeHTTP(res, req)

	if res.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusOK)
	}

	var payload models.Order
	if err := json.NewDecoder(res.Body).Decode(&payload); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if payload.ID != 101 {
		t.Fatalf("order id = %d, want 101", payload.ID)
	}
	if store.createOrderCalls != 0 {
		t.Fatalf("createOrderCalls = %d, want 0", store.createOrderCalls)
	}
}

func TestListOrdersRejectsAnonymousRequest(t *testing.T) {
	t.Parallel()

	app := newTestApplication(t, &fakeStore{})
	res := httptest.NewRecorder()

	app.Routes().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/api/v1/orders", nil))

	if res.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusUnauthorized)
	}

	var payload map[string]map[string]string
	if err := json.NewDecoder(res.Body).Decode(&payload); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if payload["error"]["code"] != "unauthorized" {
		t.Fatalf("error code = %q, want unauthorized", payload["error"]["code"])
	}
}

func TestHealthRouteBypassesRateLimit(t *testing.T) {
	t.Parallel()

	app := newTestApplication(t, &fakeStore{})

	for i := 0; i < apiRateLimitMaxRequests+5; i++ {
		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/health", nil)

		app.Routes().ServeHTTP(res, req)

		if res.Code != http.StatusOK {
			t.Fatalf("iteration %d status = %d, want %d", i, res.Code, http.StatusOK)
		}
	}
}

func TestCreatePaymentUsesAuthenticatedTenant(t *testing.T) {
	t.Parallel()

	store := &fakeStore{orderByID: paymentTestOrder()}
	app := newTestApplication(t, store)
	token := issueHandlerTestToken(t, app.Tokens, auth.Identity{
		UserID:   42,
		TenantID: 7,
		Email:    "admin@example.com",
		Role:     models.UserRoleAdmin,
	})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/payments", bytes.NewBufferString(`{
		"order_id": 101,
		"provider_reference": "pay_123",
		"amount_cents": 2500
	}`))
	req.Header.Set("Authorization", "Bearer "+token)
	res := httptest.NewRecorder()

	app.Routes().ServeHTTP(res, req)

	if res.Code != http.StatusCreated {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusCreated)
	}

	if store.createdPayment.TenantID != 7 {
		t.Fatalf("payment tenant = %d, want 7", store.createdPayment.TenantID)
	}

	if store.updatedPayment.Status != models.PaymentStatusSettled {
		t.Fatalf("payment status = %q, want settled", store.updatedPayment.Status)
	}
}

func TestCreatePaymentReturnsNotFoundForInvalidOrder(t *testing.T) {
	t.Parallel()

	store := &fakeStore{}
	app := newTestApplication(t, store)
	token := issueHandlerTestToken(t, app.Tokens, auth.Identity{
		UserID:   42,
		TenantID: 7,
		Email:    "admin@example.com",
		Role:     models.UserRoleAdmin,
	})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/payments", bytes.NewBufferString(`{
		"order_id": 404,
		"provider_reference": "pay_missing",
		"amount_cents": 2500
	}`))
	req.Header.Set("Authorization", "Bearer "+token)
	res := httptest.NewRecorder()

	app.Routes().ServeHTTP(res, req)

	if res.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusNotFound)
	}

	var payload map[string]map[string]string
	if err := json.NewDecoder(res.Body).Decode(&payload); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if payload["error"]["code"] != "order_not_found" {
		t.Fatalf("error code = %q, want order_not_found", payload["error"]["code"])
	}
}

func TestCreatePaymentReturnsExistingPaymentForDuplicateProviderReference(t *testing.T) {
	t.Parallel()

	store := &fakeStore{
		orderByID: paymentTestOrder(),
		paymentByReference: models.Payment{
			ID:                501,
			TenantID:          7,
			OrderID:           101,
			Status:            models.PaymentStatusSettled,
			ProviderReference: "pay_123",
			AmountCents:       2500,
			CreatedAt:         time.Now().UTC(),
			UpdatedAt:         time.Now().UTC(),
		},
	}
	app := newTestApplication(t, store)
	token := issueHandlerTestToken(t, app.Tokens, auth.Identity{
		UserID:   42,
		TenantID: 7,
		Email:    "admin@example.com",
		Role:     models.UserRoleAdmin,
	})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/payments", bytes.NewBufferString(`{
		"order_id": 101,
		"provider_reference": "pay_123",
		"amount_cents": 2500
	}`))
	req.Header.Set("Authorization", "Bearer "+token)
	res := httptest.NewRecorder()

	app.Routes().ServeHTTP(res, req)

	if res.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusOK)
	}

	var payload models.Payment
	if err := json.NewDecoder(res.Body).Decode(&payload); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if payload.ID != 501 {
		t.Fatalf("payment id = %d, want existing payment 501", payload.ID)
	}
}

func TestListPaymentsByOrderUsesAuthenticatedTenant(t *testing.T) {
	t.Parallel()

	store := &fakeStore{
		createdPayment: models.Payment{
			ID:       501,
			TenantID: 7,
			OrderID:  101,
			Status:   models.PaymentStatusPending,
		},
	}
	app := newTestApplication(t, store)
	token := issueHandlerTestToken(t, app.Tokens, auth.Identity{
		UserID:   42,
		TenantID: 7,
		Email:    "admin@example.com",
		Role:     models.UserRoleAdmin,
	})
	req := httptest.NewRequest(http.MethodGet, "/api/v1/orders/101/payments", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	res := httptest.NewRecorder()

	app.Routes().ServeHTTP(res, req)

	if res.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusOK)
	}

	if store.listPaymentsTenantID != 7 || store.listPaymentsOrderID != 101 {
		t.Fatalf("list payments scope = tenant %d order %d, want tenant 7 order 101", store.listPaymentsTenantID, store.listPaymentsOrderID)
	}
}

func newHandlerTestTokenManager(t *testing.T) *auth.TokenManager {
	t.Helper()

	tokens, err := auth.NewTokenManager(
		"handler-test-secret-with-at-least-thirty-two-characters",
		"opslane-test",
		time.Hour,
	)
	if err != nil {
		t.Fatalf("NewTokenManager returned error: %v", err)
	}

	return tokens
}

func newTestApplication(t *testing.T, store Store) *Application {
	t.Helper()

	orders := services.NewOrderService(store, services.NewNoopInventoryCoordinator())
	return &Application{
		Logger:      slog.Default(),
		Store:       store,
		Orders:      orders,
		Payments:    services.NewPaymentService(store, store, orders, paymentflow.NewSimulatedGateway(), services.PaymentServiceOptions{GatewayTimeout: time.Second, MaxAttempts: 1}),
		Tokens:      newHandlerTestTokenManager(t),
		ServiceName: "opslane",
		Environment: "test",
	}
}

func issueHandlerTestToken(t *testing.T, tokens *auth.TokenManager, identity auth.Identity) string {
	t.Helper()

	token, err := tokens.Issue(identity)
	if err != nil {
		t.Fatalf("Issue returned error: %v", err)
	}

	return token
}

type fakeStore struct {
	createdTenant          models.Tenant
	createTenantErr        error
	createdUser            models.User
	createUserErr          error
	userByEmail            models.User
	orderByID              models.Order
	orderByIDErr           error
	orderByIdempotencyKey  models.Order
	orderByIdempotencyErr  error
	createdOrder           models.Order
	createOrderCalls       int
	createOrderErr         error
	updatedOrder           models.Order
	updateOrderStatusErr   error
	createdPayment         models.Payment
	createPaymentErr       error
	paymentByReference     models.Payment
	paymentByReferenceErr  error
	updatedPayment         models.Payment
	updatePaymentStatusErr error
	listPaymentsTenantID   int64
	listPaymentsOrderID    int64
}

func (s *fakeStore) Ping(context.Context) error {
	return nil
}

func (s *fakeStore) CreateTenant(_ context.Context, tenant *models.Tenant) error {
	if s.createTenantErr != nil {
		return s.createTenantErr
	}

	tenant.ID = 7
	tenant.CreatedAt = time.Now().UTC()
	s.createdTenant = *tenant
	return nil
}

func (s *fakeStore) CreateUser(_ context.Context, user *models.User) error {
	if s.createUserErr != nil {
		return s.createUserErr
	}

	user.ID = 42
	user.CreatedAt = time.Now().UTC()
	s.createdUser = *user
	return nil
}

func (s *fakeStore) GetUserByEmail(context.Context, int64, string) (models.User, error) {
	return s.userByEmail, nil
}

func (s *fakeStore) CreateOrder(_ context.Context, order *models.Order) error {
	s.createOrderCalls++
	if s.createOrderErr != nil {
		return s.createOrderErr
	}

	order.ID = 101
	order.CreatedAt = time.Now().UTC()
	order.UpdatedAt = order.CreatedAt
	s.createdOrder = *order
	return nil
}

func (s *fakeStore) GetOrderByID(context.Context, int64, int64) (models.Order, error) {
	if s.orderByIDErr != nil {
		return models.Order{}, s.orderByIDErr
	}
	if s.orderByID.ID == 0 {
		return models.Order{}, sql.ErrNoRows
	}

	return s.orderByID, nil
}

func (s *fakeStore) GetOrderByIdempotencyKey(context.Context, int64, string) (models.Order, error) {
	if s.orderByIdempotencyErr != nil {
		return models.Order{}, s.orderByIdempotencyErr
	}
	if s.orderByIdempotencyKey.ID == 0 {
		return models.Order{}, sql.ErrNoRows
	}

	return s.orderByIdempotencyKey, nil
}

func (s *fakeStore) UpdateOrderStatus(_ context.Context, _ int64, _ int64, status models.OrderStatus) (models.Order, error) {
	if s.updateOrderStatusErr != nil {
		return models.Order{}, s.updateOrderStatusErr
	}

	s.updatedOrder = s.orderByID
	s.updatedOrder.Status = status
	s.updatedOrder.UpdatedAt = time.Now().UTC()
	s.orderByID = s.updatedOrder
	return s.updatedOrder, nil
}

func (s *fakeStore) ListOrdersByTenant(context.Context, int64) ([]models.Order, error) {
	return []models.Order{s.createdOrder}, nil
}

func (s *fakeStore) CreatePayment(_ context.Context, payment *models.Payment) error {
	if s.createPaymentErr != nil {
		return s.createPaymentErr
	}

	payment.ID = 501
	payment.CreatedAt = time.Now().UTC()
	payment.UpdatedAt = payment.CreatedAt
	s.createdPayment = *payment
	return nil
}

func (s *fakeStore) GetPaymentByProviderReference(context.Context, int64, string) (models.Payment, error) {
	if s.paymentByReferenceErr != nil {
		return models.Payment{}, s.paymentByReferenceErr
	}
	if s.paymentByReference.ID == 0 {
		return models.Payment{}, sql.ErrNoRows
	}

	return s.paymentByReference, nil
}

func (s *fakeStore) UpdatePaymentStatus(_ context.Context, tenantID int64, providerReference string, status models.PaymentStatus, failureReason string) (models.Payment, error) {
	if s.updatePaymentStatusErr != nil {
		return models.Payment{}, s.updatePaymentStatusErr
	}

	payment := s.createdPayment
	if payment.ID == 0 || payment.TenantID != tenantID || payment.ProviderReference != providerReference {
		payment = s.paymentByReference
	}
	if payment.ID == 0 {
		return models.Payment{}, sql.ErrNoRows
	}

	payment.Status = status
	payment.FailureReason = failureReason
	payment.UpdatedAt = time.Now().UTC()
	s.updatedPayment = payment
	s.paymentByReference = payment
	return payment, nil
}

func (s *fakeStore) ListPaymentsByOrder(_ context.Context, tenantID, orderID int64) ([]models.Payment, error) {
	s.listPaymentsTenantID = tenantID
	s.listPaymentsOrderID = orderID
	return []models.Payment{s.createdPayment}, nil
}

func paymentTestOrder() models.Order {
	now := time.Now().UTC()
	return models.Order{
		ID:             101,
		TenantID:       7,
		UserID:         42,
		Status:         models.OrderStatusPending,
		TotalCents:     2500,
		Currency:       "USD",
		IdempotencyKey: "checkout-123",
		CreatedAt:      now,
		UpdatedAt:      now,
	}
}
