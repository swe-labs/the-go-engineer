// Copyright (c) 2026 Rasel Hossen
// See LICENSE for usage terms.

package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/rasel9t6/the-go-engineer/11-flagship/01-opslane/internal/auth"
	"github.com/rasel9t6/the-go-engineer/11-flagship/01-opslane/internal/db"
	"github.com/rasel9t6/the-go-engineer/11-flagship/01-opslane/internal/models"
	paymentflow "github.com/rasel9t6/the-go-engineer/11-flagship/01-opslane/internal/payment"
	"github.com/rasel9t6/the-go-engineer/11-flagship/01-opslane/internal/services"
)

const maxJSONBodySize = 1 << 20

type createTenantRequest struct {
	Name string `json:"name"`
	Slug string `json:"slug"`
}

type createUserRequest struct {
	TenantID    int64           `json:"tenant_id"`
	Email       string          `json:"email"`
	DisplayName string          `json:"display_name"`
	Password    string          `json:"password"`
	Role        models.UserRole `json:"role"`
}

type loginRequest struct {
	TenantID int64  `json:"tenant_id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type createOrderRequest struct {
	TotalCents     int64  `json:"total_cents"`
	Currency       string `json:"currency"`
	IdempotencyKey string `json:"idempotency_key"`
}

type createPaymentRequest struct {
	OrderID           int64                `json:"order_id"`
	Status            models.PaymentStatus `json:"status"`
	ProviderReference string               `json:"provider_reference"`
	AmountCents       int64                `json:"amount_cents"`
	FailureReason     string               `json:"failure_reason"`
}

func (app *Application) handleCreateTenant(w http.ResponseWriter, r *http.Request) {
	var req createTenantRequest
	if err := readJSON(w, r, &req); err != nil {
		app.writeError(w, r, http.StatusBadRequest, "invalid_json", err.Error())
		return
	}

	req.Name = strings.TrimSpace(req.Name)
	req.Slug = strings.ToLower(strings.TrimSpace(req.Slug))
	if req.Name == "" || req.Slug == "" {
		app.writeError(w, r, http.StatusBadRequest, "invalid_tenant", "name and slug are required")
		return
	}

	tenant := models.Tenant{
		Name: req.Name,
		Slug: req.Slug,
	}
	if err := app.Store.CreateTenant(r.Context(), &tenant); err != nil {
		if errors.Is(err, db.ErrDuplicateValue) {
			app.writeError(w, r, http.StatusConflict, "duplicate_slug", "tenant slug already exists")
			return
		}
		app.writeError(w, r, http.StatusInternalServerError, "tenant_create_failed", "failed to create tenant")
		return
	}

	writeJSON(w, http.StatusCreated, tenant)
}

func (app *Application) handleCreateUser(w http.ResponseWriter, r *http.Request) {
	var req createUserRequest
	if err := readJSON(w, r, &req); err != nil {
		app.writeError(w, r, http.StatusBadRequest, "invalid_json", err.Error())
		return
	}

	role := req.Role
	if role == "" {
		role = models.UserRoleMember
	}
	if !isAllowedUserRole(role) {
		app.writeError(w, r, http.StatusBadRequest, "invalid_role", "role must be admin, member, or billing")
		return
	}

	email := strings.ToLower(strings.TrimSpace(req.Email))
	displayName := strings.TrimSpace(req.DisplayName)
	if req.TenantID <= 0 || email == "" || displayName == "" {
		app.writeError(w, r, http.StatusBadRequest, "invalid_user", "tenant_id, email, and display_name are required")
		return
	}

	passwordHash, err := auth.HashPassword(req.Password)
	if err != nil {
		app.writeError(w, r, http.StatusBadRequest, "weak_password", "password does not meet policy")
		return
	}

	user := models.User{
		TenantID:     req.TenantID,
		Email:        email,
		DisplayName:  displayName,
		PasswordHash: passwordHash,
		Role:         role,
	}

	if err := app.Store.CreateUser(r.Context(), &user); err != nil {
		if errors.Is(err, db.ErrInvalidReference) {
			app.writeError(w, r, http.StatusNotFound, "tenant_not_found", "tenant does not exist")
			return
		}
		if errors.Is(err, db.ErrDuplicateValue) {
			app.writeError(w, r, http.StatusConflict, "duplicate_email", "email already exists for this tenant")
			return
		}
		app.writeError(w, r, http.StatusInternalServerError, "user_create_failed", "failed to create user")
		return
	}

	writeJSON(w, http.StatusCreated, user)
}

func (app *Application) handleLogin(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := readJSON(w, r, &req); err != nil {
		app.writeError(w, r, http.StatusBadRequest, "invalid_json", err.Error())
		return
	}

	service := auth.NewService(app.Store, app.Tokens)
	result, err := service.Login(r.Context(), auth.LoginRequest{
		TenantID: req.TenantID,
		Email:    strings.ToLower(strings.TrimSpace(req.Email)),
		Password: req.Password,
	})
	if err != nil {
		app.writeError(w, r, http.StatusUnauthorized, "invalid_credentials", "invalid credentials")
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"access_token": result.Token,
		"token_type":   "Bearer",
		"expires_at":   result.Identity.ExpiresAt,
		"identity":     result.Identity,
	})
}

func (app *Application) requireAPIAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		identity, err := auth.IdentityFromRequest(app.Tokens, r)
		if err != nil {
			app.writeError(w, r, http.StatusUnauthorized, "unauthorized", "missing or invalid bearer token")
			return
		}

		next.ServeHTTP(w, r.WithContext(auth.WithIdentity(r.Context(), identity)))
	})
}

func (app *Application) handleListOrders(w http.ResponseWriter, r *http.Request) {
	identity, err := auth.RequireIdentity(r.Context())
	if err != nil {
		app.writeError(w, r, http.StatusInternalServerError, "missing_identity", "missing authenticated identity")
		return
	}

	orders, err := app.Store.ListOrdersByTenant(r.Context(), identity.TenantID)
	if err != nil {
		app.writeError(w, r, http.StatusInternalServerError, "orders_list_failed", "failed to list orders")
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{"data": orders})
}

func (app *Application) handleCreateOrder(w http.ResponseWriter, r *http.Request) {
	identity, err := auth.RequireIdentity(r.Context())
	if err != nil {
		app.writeError(w, r, http.StatusInternalServerError, "missing_identity", "missing authenticated identity")
		return
	}

	var req createOrderRequest
	if err := readJSON(w, r, &req); err != nil {
		app.writeError(w, r, http.StatusBadRequest, "invalid_json", err.Error())
		return
	}

	if app.Orders == nil {
		app.writeError(w, r, http.StatusInternalServerError, "order_service_unavailable", "order service is not configured")
		return
	}

	result, err := app.Orders.CreateOrder(r.Context(), services.CreateOrderInput{
		TenantID:       identity.TenantID,
		UserID:         identity.UserID,
		TotalCents:     req.TotalCents,
		Currency:       req.Currency,
		IdempotencyKey: req.IdempotencyKey,
	})
	if err != nil {
		switch {
		case errors.Is(err, services.ErrInvalidOrder):
			app.writeError(w, r, http.StatusBadRequest, "invalid_order", "total_cents, currency, and idempotency_key are required")
		case errors.Is(err, services.ErrInventoryUnavailable):
			app.writeError(w, r, http.StatusConflict, "inventory_unavailable", "inventory could not be reserved for this order")
		default:
			app.writeError(w, r, http.StatusInternalServerError, "order_create_failed", "failed to create order")
		}
		return
	}

	statusCode := http.StatusCreated
	if !result.Created {
		statusCode = http.StatusOK
	}

	writeJSON(w, statusCode, result.Order)
}

func (app *Application) handleCreatePayment(w http.ResponseWriter, r *http.Request) {
	identity, err := auth.RequireIdentity(r.Context())
	if err != nil {
		app.writeError(w, r, http.StatusInternalServerError, "missing_identity", "missing authenticated identity")
		return
	}

	var req createPaymentRequest
	if err := readJSON(w, r, &req); err != nil {
		app.writeError(w, r, http.StatusBadRequest, "invalid_json", err.Error())
		return
	}

	if req.OrderID <= 0 || req.AmountCents <= 0 || strings.TrimSpace(req.ProviderReference) == "" {
		app.writeError(w, r, http.StatusBadRequest, "invalid_payment", "order_id, amount_cents, and provider_reference are required")
		return
	}

	if app.Payments == nil {
		app.writeError(w, r, http.StatusInternalServerError, "payment_service_unavailable", "payment service is not configured")
		return
	}

	result, err := app.Payments.ProcessPayment(r.Context(), paymentflow.Job{
		TenantID:          identity.TenantID,
		OrderID:           req.OrderID,
		ProviderReference: strings.TrimSpace(req.ProviderReference),
		AmountCents:       req.AmountCents,
	})
	if err != nil {
		switch {
		case errors.Is(err, services.ErrInvalidPayment):
			app.writeError(w, r, http.StatusBadRequest, "invalid_payment", "payment does not match the tenant-scoped order")
		case errors.Is(err, services.ErrOrderNotFound):
			app.writeError(w, r, http.StatusNotFound, "order_not_found", "order does not exist for this tenant")
		case errors.Is(err, paymentflow.ErrGatewayTimeout), errors.Is(err, paymentflow.ErrGatewayUnavailable):
			writeJSON(w, http.StatusAccepted, result.Payment)
		default:
			app.writeError(w, r, http.StatusInternalServerError, "payment_create_failed", "failed to create payment")
		}
		return
	}

	statusCode := http.StatusCreated
	if !result.Created {
		statusCode = http.StatusOK
	}

	writeJSON(w, statusCode, result.Payment)
}

func (app *Application) handleListPaymentsByOrder(w http.ResponseWriter, r *http.Request) {
	identity, err := auth.RequireIdentity(r.Context())
	if err != nil {
		app.writeError(w, r, http.StatusInternalServerError, "missing_identity", "missing authenticated identity")
		return
	}

	orderID, err := strconv.ParseInt(r.PathValue("orderID"), 10, 64)
	if err != nil || orderID <= 0 {
		app.writeError(w, r, http.StatusBadRequest, "invalid_order_id", "order id must be a positive integer")
		return
	}

	payments, err := app.Store.ListPaymentsByOrder(r.Context(), identity.TenantID, orderID)
	if err != nil {
		app.writeError(w, r, http.StatusInternalServerError, "payments_list_failed", "failed to list payments")
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{"data": payments})
}

func readJSON(w http.ResponseWriter, r *http.Request, dst any) error {
	r.Body = http.MaxBytesReader(w, r.Body, maxJSONBodySize)
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(dst); err != nil {
		return fmt.Errorf("decode json body: %w", err)
	}

	if err := decoder.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		return fmt.Errorf("body must contain a single json object")
	}

	return nil
}

func isAllowedUserRole(role models.UserRole) bool {
	switch role {
	case models.UserRoleAdmin, models.UserRoleMember, models.UserRoleBilling:
		return true
	default:
		return false
	}
}

func isAllowedPaymentStatus(status models.PaymentStatus) bool {
	switch status {
	case models.PaymentStatusPending, models.PaymentStatusAuthorized, models.PaymentStatusSettled, models.PaymentStatusFailed, models.PaymentStatusRefunded:
		return true
	default:
		return false
	}
}

func (app *Application) writeError(w http.ResponseWriter, _ *http.Request, status int, code, message string) {
	writeJSON(w, status, map[string]any{
		"error": map[string]string{
			"code":    code,
			"message": message,
		},
	})
}
