package main

import (
	"strings"
	"testing"
	"time"
)

func TestCreateJWT_Valid(t *testing.T) {
	token, err := CreateJWT("alice", "admin", []byte("secret"), 1*time.Hour)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		t.Fatalf("expected 3 dot-separated segments, got %d", len(parts))
	}
}

func TestVerifyJWT_CorrectSecret(t *testing.T) {
	token, _ := CreateJWT("alice", "admin", []byte("secret"), 1*time.Hour)
	sub, role, err := VerifyJWT(token, []byte("secret"))
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if sub != "alice" {
		t.Errorf("expected alice, got %s", sub)
	}
	if role != "admin" {
		t.Errorf("expected admin, got %s", role)
	}
}

func TestVerifyJWT_WrongSecret(t *testing.T) {
	token, _ := CreateJWT("alice", "admin", []byte("secret"), 1*time.Hour)
	_, _, err := VerifyJWT(token, []byte("wrong-secret"))
	if err == nil {
		t.Fatal("expected error for wrong secret")
	}
}

func TestVerifyJWT_Expired(t *testing.T) {
	token, _ := CreateJWT("alice", "admin", []byte("secret"), -1*time.Hour) // expired
	_, _, err := VerifyJWT(token, []byte("secret"))
	if err == nil {
		t.Fatal("expected error for expired token")
	}
}

func TestIssueAndValidateRSA(t *testing.T) {
	token, err := IssueToken("bob", "editor", "tenant-1", 1*time.Hour)
	if err != nil {
		t.Fatal(err)
	}
	claims, err := ValidateToken(token)
	if err != nil {
		t.Fatal(err)
	}
	if claims.Subject != "bob" {
		t.Errorf("expected bob, got %s", claims.Subject)
	}
	if claims.Role != "editor" {
		t.Errorf("expected editor, got %s", claims.Role)
	}
}
