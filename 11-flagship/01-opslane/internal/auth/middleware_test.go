package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/rasel9t6/the-go-engineer/11-flagship/01-opslane/internal/models"
)

func TestRequireAuthAddsIdentityToRequestContext(t *testing.T) {
	t.Parallel()

	tokens := newTestTokenManager(t, time.Now().UTC(), time.Hour)
	token, err := tokens.Issue(Identity{
		UserID:   42,
		TenantID: 7,
		Email:    "admin@example.com",
		Role:     models.UserRoleAdmin,
	})
	if err != nil {
		t.Fatalf("Issue returned error: %v", err)
	}

	var got Identity
	handler := RequireAuth(tokens)(http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
		identity, ok := IdentityFromContext(r.Context())
		if !ok {
			t.Fatal("identity missing from request context")
		}
		got = identity
	}))

	req := httptest.NewRequest(http.MethodGet, "/me", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	res := httptest.NewRecorder()

	handler.ServeHTTP(res, req)

	if res.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusOK)
	}

	if got.UserID != 42 || got.TenantID != 7 {
		t.Fatalf("identity = %+v, want user 42 tenant 7", got)
	}
}

func TestRequireAuthAcceptsCaseInsensitiveBearerScheme(t *testing.T) {
	t.Parallel()

	tokens := newTestTokenManager(t, time.Now().UTC(), time.Hour)
	token, err := tokens.Issue(Identity{
		UserID:   42,
		TenantID: 7,
		Email:    "admin@example.com",
		Role:     models.UserRoleAdmin,
	})
	if err != nil {
		t.Fatalf("Issue returned error: %v", err)
	}

	handler := RequireAuth(tokens)(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}))
	req := httptest.NewRequest(http.MethodGet, "/me", nil)
	req.Header.Set("Authorization", "bearer "+token)
	res := httptest.NewRecorder()

	handler.ServeHTTP(res, req)

	if res.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusOK)
	}
}

func TestRequireAuthRejectsMissingBearerToken(t *testing.T) {
	t.Parallel()

	tokens := newTestTokenManager(t, time.Now().UTC(), time.Hour)
	handler := RequireAuth(tokens)(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		t.Fatal("next handler should not run without auth")
	}))

	res := httptest.NewRecorder()
	handler.ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/me", nil))

	if res.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusUnauthorized)
	}
}
