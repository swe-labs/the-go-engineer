package main

import (
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
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

func TestVerifyJWT_RejectsInvalidClaims(t *testing.T) {
	secret := []byte("secret")
	now := time.Now()
	tests := []struct {
		name   string
		claims jwt.MapClaims
	}{
		{
			name: "numeric subject",
			claims: jwt.MapClaims{
				"sub":  12345,
				"role": "admin",
				"exp":  now.Add(1 * time.Hour).Unix(),
			},
		},
		{
			name: "missing subject",
			claims: jwt.MapClaims{
				"role": "admin",
				"exp":  now.Add(1 * time.Hour).Unix(),
			},
		},
		{
			name: "empty subject",
			claims: jwt.MapClaims{
				"sub":  "",
				"role": "admin",
				"exp":  now.Add(1 * time.Hour).Unix(),
			},
		},
		{
			name: "non-string role",
			claims: jwt.MapClaims{
				"sub":  "alice",
				"role": 42,
				"exp":  now.Add(1 * time.Hour).Unix(),
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tok, err := jwt.NewWithClaims(jwt.SigningMethodHS256, tc.claims).SignedString(secret)
			if err != nil {
				t.Fatalf("sign token: %v", err)
			}
			if _, _, err := VerifyJWT(tok, secret); err == nil {
				t.Errorf("VerifyJWT should fail for %s", tc.name)
			}
		})
	}
}
