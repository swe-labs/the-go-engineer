// Copyright (c) 2026 Rasel Hossen
// See LICENSE for usage terms.

package handlers

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/rasel9t6/the-go-engineer/11-flagship/01-opslane/internal/auth"
	"github.com/rasel9t6/the-go-engineer/11-flagship/01-opslane/internal/db"
	"github.com/rasel9t6/the-go-engineer/11-flagship/01-opslane/internal/middleware"
)

const healthDatabaseTimeout = 2 * time.Second

type Application struct {
	Logger      *slog.Logger
	Store       *db.Store
	Tokens      *auth.TokenManager
	ServiceName string
	Environment string
}

func (app *Application) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /", app.handleIndex)
	mux.HandleFunc("GET /health", app.handleHealth)
	mux.Handle("GET /me", auth.RequireAuth(app.Tokens)(http.HandlerFunc(app.handleMe)))

	handler := middleware.SecureHeaders(
		middleware.LogRequest(app.Logger)(
			middleware.RecoverPanic(app.Logger)(mux),
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

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}
