package auth

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/rasel9t6/the-go-engineer/11-flagship/01-opslane/internal/models"
)

func TestServiceLoginIssuesTokenForValidCredentials(t *testing.T) {
	t.Parallel()

	passwordHash, err := HashPassword("CorrectHorse7Battery")
	if err != nil {
		t.Fatalf("HashPassword returned error: %v", err)
	}

	users := stubUserLookup{
		user: models.User{
			ID:           42,
			TenantID:     7,
			Email:        "admin@example.com",
			PasswordHash: passwordHash,
			Role:         models.UserRoleAdmin,
		},
	}
	tokens := newTestTokenManager(t, time.Now().UTC(), time.Hour)
	service := NewService(users, tokens)

	result, err := service.Login(context.Background(), LoginRequest{
		TenantID: 7,
		Email:    "admin@example.com",
		Password: "CorrectHorse7Battery",
	})
	if err != nil {
		t.Fatalf("Login returned error: %v", err)
	}

	if result.Token == "" {
		t.Fatal("token should not be empty")
	}

	if result.Identity.UserID != 42 || result.Identity.TenantID != 7 {
		t.Fatalf("identity = %+v, want user 42 tenant 7", result.Identity)
	}
}

func TestServiceLoginRejectsWrongPassword(t *testing.T) {
	t.Parallel()

	passwordHash, err := HashPassword("CorrectHorse7Battery")
	if err != nil {
		t.Fatalf("HashPassword returned error: %v", err)
	}

	service := NewService(stubUserLookup{
		user: models.User{
			ID:           42,
			TenantID:     7,
			Email:        "admin@example.com",
			PasswordHash: passwordHash,
			Role:         models.UserRoleAdmin,
		},
	}, newTestTokenManager(t, time.Now().UTC(), time.Hour))

	_, err = service.Login(context.Background(), LoginRequest{
		TenantID: 7,
		Email:    "admin@example.com",
		Password: "WrongHorse7Battery",
	})
	if !errors.Is(err, ErrInvalidCredentials) {
		t.Fatalf("Login error = %v, want ErrInvalidCredentials", err)
	}
}

func TestServiceLoginHidesUserLookupFailures(t *testing.T) {
	t.Parallel()

	service := NewService(stubUserLookup{err: sql.ErrNoRows}, newTestTokenManager(t, time.Now().UTC(), time.Hour))
	_, err := service.Login(context.Background(), LoginRequest{
		TenantID: 7,
		Email:    "missing@example.com",
		Password: "CorrectHorse7Battery",
	})
	if !errors.Is(err, ErrInvalidCredentials) {
		t.Fatalf("Login error = %v, want ErrInvalidCredentials", err)
	}
}

func TestServiceLoginRejectsMismatchedTenantFromRepository(t *testing.T) {
	t.Parallel()

	passwordHash, err := HashPassword("CorrectHorse7Battery")
	if err != nil {
		t.Fatalf("HashPassword returned error: %v", err)
	}

	service := NewService(stubUserLookup{
		user: models.User{
			ID:           42,
			TenantID:     99,
			Email:        "admin@example.com",
			PasswordHash: passwordHash,
			Role:         models.UserRoleAdmin,
		},
	}, newTestTokenManager(t, time.Now().UTC(), time.Hour))

	_, err = service.Login(context.Background(), LoginRequest{
		TenantID: 7,
		Email:    "admin@example.com",
		Password: "CorrectHorse7Battery",
	})
	if !errors.Is(err, ErrInvalidCredentials) {
		t.Fatalf("Login error = %v, want ErrInvalidCredentials", err)
	}
}

type stubUserLookup struct {
	user models.User
	err  error
}

func (s stubUserLookup) GetUserByEmail(context.Context, int64, string) (models.User, error) {
	if s.err != nil {
		return models.User{}, s.err
	}

	return s.user, nil
}
