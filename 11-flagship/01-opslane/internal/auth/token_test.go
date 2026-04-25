package auth

import (
	"encoding/base64"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/rasel9t6/the-go-engineer/11-flagship/01-opslane/internal/models"
)

const testTokenSecret = "test-secret-with-at-least-thirty-two-characters"

func TestTokenManagerIssuesAndVerifiesTenantIdentity(t *testing.T) {
	t.Parallel()

	baseTime := time.Date(2026, 4, 25, 12, 0, 0, 0, time.UTC)
	tokens := newTestTokenManager(t, baseTime, time.Hour)

	token, err := tokens.Issue(Identity{
		UserID:   42,
		TenantID: 7,
		Email:    "admin@example.com",
		Role:     models.UserRoleAdmin,
	})
	if err != nil {
		t.Fatalf("Issue returned error: %v", err)
	}

	identity, err := tokens.Verify(token)
	if err != nil {
		t.Fatalf("Verify returned error: %v", err)
	}

	if identity.UserID != 42 || identity.TenantID != 7 {
		t.Fatalf("identity = %+v, want user 42 tenant 7", identity)
	}

	if identity.Role != models.UserRoleAdmin {
		t.Fatalf("role = %q, want admin", identity.Role)
	}

	if !identity.ExpiresAt.Equal(baseTime.Add(time.Hour)) {
		t.Fatalf("expires at = %v, want %v", identity.ExpiresAt, baseTime.Add(time.Hour))
	}
}

func TestTokenManagerRejectsTamperedToken(t *testing.T) {
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

	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		t.Fatalf("token has %d parts, want 3", len(parts))
	}

	signature, err := base64.RawURLEncoding.DecodeString(parts[2])
	if err != nil {
		t.Fatalf("decode signature: %v", err)
	}
	signature[0] ^= 0xff

	tampered := parts[0] + "." + parts[1] + "." + base64.RawURLEncoding.EncodeToString(signature)
	_, err = tokens.Verify(tampered)
	if !errors.Is(err, ErrInvalidToken) {
		t.Fatalf("Verify error = %v, want ErrInvalidToken", err)
	}
}

func TestTokenManagerRejectsExpiredToken(t *testing.T) {
	t.Parallel()

	baseTime := time.Date(2026, 4, 25, 12, 0, 0, 0, time.UTC)
	tokens := newTestTokenManager(t, baseTime, time.Minute)
	token, err := tokens.Issue(Identity{
		UserID:   42,
		TenantID: 7,
		Email:    "admin@example.com",
		Role:     models.UserRoleAdmin,
	})
	if err != nil {
		t.Fatalf("Issue returned error: %v", err)
	}

	tokens.now = func() time.Time { return baseTime.Add(2 * time.Minute) }

	_, err = tokens.Verify(token)
	if !errors.Is(err, ErrExpiredToken) {
		t.Fatalf("Verify error = %v, want ErrExpiredToken", err)
	}
}

func TestTokenManagerRejectsInvalidIdentity(t *testing.T) {
	t.Parallel()

	tokens := newTestTokenManager(t, time.Now().UTC(), time.Hour)
	_, err := tokens.Issue(Identity{
		UserID:   0,
		TenantID: 7,
		Email:    "admin@example.com",
		Role:     models.UserRoleAdmin,
	})
	if err == nil {
		t.Fatal("expected error for invalid identity")
	}
}

func newTestTokenManager(t *testing.T, now time.Time, ttl time.Duration) *TokenManager {
	t.Helper()

	tokens, err := NewTokenManager(testTokenSecret, "opslane-test", ttl)
	if err != nil {
		t.Fatalf("NewTokenManager returned error: %v", err)
	}

	tokens.now = func() time.Time { return now }
	return tokens
}
