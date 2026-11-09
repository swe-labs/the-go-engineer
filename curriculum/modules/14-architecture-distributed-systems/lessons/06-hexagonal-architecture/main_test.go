package main

import (
	"errors"
	"testing"
)

type mockEmailAdapter struct {
	sentEmails []string
	failOnSend bool
}

func (m *mockEmailAdapter) SendWelcome(email string) error {
	if m.failOnSend {
		return errors.New("send failed")
	}
	m.sentEmails = append(m.sentEmails, email)
	return nil
}

func TestUserService_RegisterUser_Success(t *testing.T) {
	repo := NewInMemoryUserRepoAdapter()
	email := &mockEmailAdapter{}
	svc := NewUserService(repo, email)

	err := svc.RegisterUser("u1", "alice@example.com")
	if err != nil {
		t.Fatalf("RegisterUser failed: %v", err)
	}
	if len(email.sentEmails) != 1 || email.sentEmails[0] != "alice@example.com" {
		t.Errorf("welcome email not sent")
	}
	user, err := repo.FindByID("u1")
	if err != nil {
		t.Fatalf("FindByID failed: %v", err)
	}
	if user.Email != "alice@example.com" {
		t.Errorf("expected alice@example.com, got %s", user.Email)
	}
}

func TestUserService_RegisterUser_EmptyEmail(t *testing.T) {
	svc := NewUserService(NewInMemoryUserRepoAdapter(), &mockEmailAdapter{})
	err := svc.RegisterUser("u1", "")
	if err == nil {
		t.Fatal("expected error for empty email")
	}
}

func TestUserService_RegisterUser_EmailFails(t *testing.T) {
	repo := NewInMemoryUserRepoAdapter()
	email := &mockEmailAdapter{failOnSend: true}
	svc := NewUserService(repo, email)

	err := svc.RegisterUser("u1", "alice@example.com")
	if err == nil {
		t.Fatal("expected error when email fails")
	}
	_, err = repo.FindByID("u1")
	if err == nil {
		t.Fatal("expected user to not be saved when email fails")
	}
}

func TestLoggerMiddleware_LogsOperations(t *testing.T) {
	inner := NewInMemoryUserRepoAdapter()
	mw := &LoggerMiddleware{next: inner}
	mw.Save(User{ID: "u1", Email: "a@b.com"})
	mw.FindByID("u1")
	if len(mw.logs) != 2 {
		t.Errorf("expected 2 log entries, got %d", len(mw.logs))
	}
}

func TestInMemoryUserRepoAdapter_FindByID_NotFound(t *testing.T) {
	repo := NewInMemoryUserRepoAdapter()
	_, err := repo.FindByID("missing")
	if err == nil {
		t.Fatal("expected error for missing user")
	}
}
