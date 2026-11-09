package main

import (
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func TestVerifyTokenWithAlgCheck_ValidRS256(t *testing.T) {
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"sub": "alice",
		"exp": time.Now().Add(1 * time.Hour).Unix(),
	})
	tokenStr, _ := token.SignedString(rsaPrivateKey)
	_, err := VerifyTokenWithAlgCheck(tokenStr, rsaPublicKey)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestVerifyTokenWithAlgCheck_RejectsHS256(t *testing.T) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": "admin",
		"exp": time.Now().Add(1 * time.Hour).Unix(),
	})
	tokenStr, _ := token.SignedString([]byte("any-secret"))
	_, err := VerifyTokenWithAlgCheck(tokenStr, rsaPublicKey)
	if err == nil {
		t.Fatal("expected error for HS256 when expecting RS256")
	}
	if !strings.Contains(err.Error(), "unexpected signing method") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestVerifyTokenWithAlgCheck_RejectsNoneAlg(t *testing.T) {
	// Craft a "none" algorithm token.
	noneToken := "eyJhbGciOiJub25lIiwidHlwIjoiSldUIn0.eyJzdWIiOiJhZG1pbiJ9."
	_, err := VerifyTokenWithAlgCheck(noneToken, rsaPublicKey)
	if err == nil {
		t.Fatal("expected error for none alg")
	}
}

func TestVerifyTokenSafe_Valid(t *testing.T) {
	tokenStr, _ := IssueToken("bob", "editor")
	dl := NewDenyList()
	claims, err := VerifyTokenSafe(tokenStr, dl)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if claims.Subject != "bob" {
		t.Errorf("expected bob, got %s", claims.Subject)
	}
	if claims.Role != "editor" {
		t.Errorf("expected editor, got %s", claims.Role)
	}
}

func TestVerifyTokenSafe_Revoked(t *testing.T) {
	tokenStr, _ := IssueToken("bob", "editor")
	dl := NewDenyList()
	claims, _ := VerifyTokenSafe(tokenStr, dl)
	dl.Revoke(claims.ID, claims.ExpiresAt.Time)
	_, err := VerifyTokenSafe(tokenStr, dl)
	if err == nil {
		t.Fatal("expected error for revoked token")
	}
	if !strings.Contains(err.Error(), "revoked") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestDemonstrateForgery(t *testing.T) {
	forged := DemonstrateForgery(rsaPublicKey)
	parts := strings.Split(forged, ".")
	if len(parts) != 3 {
		t.Fatalf("expected 3 segments, got %d", len(parts))
	}
}
