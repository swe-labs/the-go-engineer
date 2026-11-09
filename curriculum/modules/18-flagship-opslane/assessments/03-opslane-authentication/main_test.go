package main

import (
	"errors"
	"testing"
	"time"

	"golang.org/x/crypto/bcrypt"
)

const testSecret = "test-secret-key-at-least-32-bytes-long!!"

func TestPasswordManager_HashAndVerify(t *testing.T) {
	pm := NewPasswordManager(bcrypt.MinCost)

	hash, err := pm.Hash("correct-horse-battery-staple")
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}

	if !pm.Verify("correct-horse-battery-staple", hash) {
		t.Errorf("expected correct password to verify successfully")
	}

	if pm.Verify("wrong-password", hash) {
		t.Errorf("expected incorrect password to fail verification")
	}
}

func TestTokenManager_GenerateAndValidate(t *testing.T) {
	tm := NewTokenManager(testSecret, 15*time.Minute)
	user := &UserAuth{
		ID:       "usr-42",
		TenantID: "tenant-globex",
		Email:    "engineer@globex.com",
		Role:     "admin",
	}

	tokenStr, expiresAt, err := tm.GenerateAccessToken(user)
	if err != nil {
		t.Fatalf("failed to generate access token: %v", err)
	}

	if time.Until(expiresAt) <= 0 {
		t.Errorf("expected expiry in future")
	}

	claims, err := tm.ValidateAccessToken(tokenStr)
	if err != nil {
		t.Fatalf("expected valid token, got: %v", err)
	}

	if claims.UserID != "usr-42" || claims.TenantID != "tenant-globex" || claims.Role != "admin" {
		t.Errorf("unexpected token claims: %+v", claims)
	}
}

func TestTokenManager_ExpiredToken(t *testing.T) {
	// Create token with negative TTL (already expired)
	tm := NewTokenManager(testSecret, -1*time.Minute)
	user := &UserAuth{
		ID:       "usr-exp",
		TenantID: "tenant-test",
		Role:     "member",
	}

	tokenStr, _, _ := tm.GenerateAccessToken(user)

	_, err := tm.ValidateAccessToken(tokenStr)
	if !errors.Is(err, ErrTokenExpired) {
		t.Errorf("expected ErrTokenExpired for expired token, got: %v", err)
	}
}

func TestTokenManager_TamperedSignature(t *testing.T) {
	tm := NewTokenManager(testSecret, 15*time.Minute)
	user := &UserAuth{ID: "usr-tamper", TenantID: "tenant-test", Role: "member"}
	tokenStr, _, _ := tm.GenerateAccessToken(user)

	// Validate with wrong secret key
	wrongTM := NewTokenManager("wrong-secret-key-different-bytes!!", 15*time.Minute)
	_, err := wrongTM.ValidateAccessToken(tokenStr)
	if !errors.Is(err, ErrInvalidToken) {
		t.Errorf("expected ErrInvalidToken when signature key does not match, got: %v", err)
	}
}

func TestAuthService_LoginSuccess(t *testing.T) {
	svc := NewAuthService(testSecret)

	res, err := svc.Login("admin@acme.corp", "opslane_secret_123")
	if err != nil {
		t.Fatalf("expected login to succeed, got: %v", err)
	}

	if res.AccessToken == "" || res.RefreshToken == "" {
		t.Errorf("expected non-empty access and refresh tokens")
	}
	if res.UserID != "usr-01" || res.TenantID != "tenant-acme" {
		t.Errorf("unexpected user or tenant in response: %+v", res)
	}
}

func TestAuthService_LoginInvalidPassword(t *testing.T) {
	// Core practice task requirement:
	// Verify that login with an incorrect password fails and confirms token is not returned.
	svc := NewAuthService(testSecret)

	res, err := svc.Login("admin@acme.corp", "completely_incorrect_password")
	if !errors.Is(err, ErrInvalidCredentials) {
		t.Fatalf("expected ErrInvalidCredentials for incorrect password, got: %v", err)
	}

	if res != nil {
		t.Fatalf("SECURITY VIOLATION: Token returned despite invalid password: %+v", res)
	}
}
