package auth

import (
	"context"
	"errors"
	"testing"

	"github.com/rasel9t6/the-go-engineer/11-flagship/01-opslane/internal/models"
)

func TestRequireIdentityReturnsContextIdentity(t *testing.T) {
	t.Parallel()

	want := Identity{
		UserID:   42,
		TenantID: 7,
		Email:    "admin@example.com",
		Role:     models.UserRoleAdmin,
	}

	got, err := RequireIdentity(WithIdentity(context.Background(), want))
	if err != nil {
		t.Fatalf("RequireIdentity returned error: %v", err)
	}

	if got != want {
		t.Fatalf("identity = %+v, want %+v", got, want)
	}
}

func TestRequireIdentityRejectsMissingIdentity(t *testing.T) {
	t.Parallel()

	_, err := RequireIdentity(context.Background())
	if !errors.Is(err, ErrMissingIdentity) {
		t.Fatalf("RequireIdentity error = %v, want ErrMissingIdentity", err)
	}
}
