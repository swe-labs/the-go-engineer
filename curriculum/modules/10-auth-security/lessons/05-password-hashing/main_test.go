package main

import (
	"strings"
	"testing"
)

func TestHashPassword_Valid(t *testing.T) {
	hash, err := HashPassword("correct-horse-battery-staple", 4)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !strings.HasPrefix(string(hash), "$2a$") {
		t.Errorf("expected hash to start with $2a$, got %s", string(hash))
	}
}

func TestHashPassword_Empty(t *testing.T) {
	_, err := HashPassword("", 4)
	if err == nil {
		t.Fatal("expected error for empty password")
	}
	if !strings.Contains(err.Error(), "too short") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestHashPassword_CostTooLow(t *testing.T) {
	_, err := HashPassword("password123", 3)
	if err == nil {
		t.Fatal("expected error for cost < 4")
	}
	if !strings.Contains(err.Error(), "cost must be between") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestVerifyPassword_Match(t *testing.T) {
	hash, err := HashPassword("my-secret-password", 4)
	if err != nil {
		t.Fatal(err)
	}
	if !VerifyPassword(hash, "my-secret-password") {
		t.Error("expected true for matching password")
	}
}

func TestVerifyPassword_NoMatch(t *testing.T) {
	hash, err := HashPassword("my-secret-password", 4)
	if err != nil {
		t.Fatal(err)
	}
	if VerifyPassword(hash, "wrong-password") {
		t.Error("expected false for wrong password")
	}
}

func TestUserDBCreateAndAuth(t *testing.T) {
	db := newUserDB()
	if err := db.CreateUser("testuser", "test-password-123"); err != nil {
		t.Fatal(err)
	}
	if err := db.Authenticate("testuser", "test-password-123"); err != nil {
		t.Errorf("expected auth success, got %v", err)
	}
	if err := db.Authenticate("testuser", "wrong-password"); err == nil {
		t.Error("expected auth failure for wrong password")
	}
	if err := db.Authenticate("nonexistent", "any"); err == nil {
		t.Error("expected error for nonexistent user")
	}
}
