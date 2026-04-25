package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/rasel9t6/the-go-engineer/11-flagship/01-opslane/internal/auth"
	"github.com/rasel9t6/the-go-engineer/11-flagship/01-opslane/internal/models"
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
