package main

import (
	"testing"
)

type mockUserRepo struct {
	saveFunc     func(user User) error
	findByIDFunc func(id string) (User, error)
}

func (m *mockUserRepo) Save(user User) error {
	return m.saveFunc(user)
}

func (m *mockUserRepo) FindByID(id string) (User, error) {
	return m.findByIDFunc(id)
}

func TestUserService_Register_Success(t *testing.T) {
	repo := &mockUserRepo{
		saveFunc: func(user User) error { return nil },
	}
	svc := NewUserService(repo)
	user, err := svc.Register("Alice", "alice@example.com")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if user.Name != "Alice" {
		t.Errorf("expected Alice, got %s", user.Name)
	}
	if user.Email != "alice@example.com" {
		t.Errorf("expected alice@example.com, got %s", user.Email)
	}
}

func TestUserService_Register_EmptyName(t *testing.T) {
	svc := NewUserService(&mockUserRepo{})
	_, err := svc.Register("", "alice@example.com")
	if err == nil {
		t.Fatal("expected error for empty name")
	}
}

func TestUserService_Register_EmptyEmail(t *testing.T) {
	svc := NewUserService(&mockUserRepo{})
	_, err := svc.Register("Alice", "")
	if err == nil {
		t.Fatal("expected error for empty email")
	}
}

func TestUserService_GetUser_Success(t *testing.T) {
	repo := &mockUserRepo{
		findByIDFunc: func(id string) (User, error) {
			return User{ID: id, Name: "Alice"}, nil
		},
	}
	svc := NewUserService(repo)
	user, err := svc.GetUser("usr_1")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if user.Name != "Alice" {
		t.Errorf("expected Alice, got %s", user.Name)
	}
}

func TestUserService_GetUser_NotFound(t *testing.T) {
	repo := &mockUserRepo{
		findByIDFunc: func(id string) (User, error) {
			return User{}, errUserNotFound
		},
	}
	svc := NewUserService(repo)
	_, err := svc.GetUser("usr_missing")
	if err == nil {
		t.Fatal("expected error for missing user")
	}
}
