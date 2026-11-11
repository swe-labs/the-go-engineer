package main

import (
	"testing"
	"time"
)

func TestIssueTokens(t *testing.T) {
	store := NewRefreshTokenStore()
	pair := store.Issue("user-abc", 15*time.Minute, 30*24*time.Hour)

	if pair.AccessToken == "" {
		t.Error("expected non-empty access token")
	}
	if pair.RefreshToken == "" {
		t.Error("expected non-empty refresh token")
	}
	if pair.ExpiresIn != 900 {
		t.Errorf("ExpiresIn = %d, want 900", pair.ExpiresIn)
	}
}

func TestValidateToken(t *testing.T) {
	store := NewRefreshTokenStore()
	pair := store.Issue("user-abc", 15*time.Minute, 30*24*time.Hour)

	userID, err := store.Validate(pair.RefreshToken)
	if err != nil {
		t.Fatalf("Validate error: %v", err)
	}
	if userID != "user-abc" {
		t.Errorf("userID = %q, want %q", userID, "user-abc")
	}
}

func TestRotateToken(t *testing.T) {
	store := NewRefreshTokenStore()
	pair := store.Issue("user-abc", 15*time.Minute, 30*24*time.Hour)

	newPair, err := store.Rotate(pair.RefreshToken, 15*time.Minute, 30*24*time.Hour)
	if err != nil {
		t.Fatalf("Rotate error: %v", err)
	}

	if newPair.AccessToken == pair.AccessToken {
		t.Error("access token should change after rotation")
	}
	if newPair.RefreshToken == pair.RefreshToken {
		t.Error("refresh token should change after rotation")
	}

	_, err = store.Validate(pair.RefreshToken)
	if err == nil {
		t.Error("expected old refresh token to be invalid after rotation")
	}
}

func TestRotateTokenNotFound(t *testing.T) {
	store := NewRefreshTokenStore()
	_, err := store.Rotate("nonexistent-token", 15*time.Minute, 30*24*time.Hour)
	if err == nil {
		t.Error("expected error for nonexistent token")
	}
}

func TestRevokeToken(t *testing.T) {
	store := NewRefreshTokenStore()
	pair := store.Issue("user-abc", 15*time.Minute, 30*24*time.Hour)

	if !store.Revoke(pair.RefreshToken) {
		t.Error("expected Revoke to return true")
	}

	_, err := store.Validate(pair.RefreshToken)
	if err == nil {
		t.Error("expected revoked token to be invalid")
	}
}

func TestRevokeNonexistentToken(t *testing.T) {
	store := NewRefreshTokenStore()
	if store.Revoke("nonexistent") {
		t.Error("expected Revoke to return false for nonexistent token")
	}
}

func TestTokenExpiration(t *testing.T) {
	store := NewRefreshTokenStore()
	pair := store.Issue("user-abc", 15*time.Minute, -1*time.Second)

	time.Sleep(10 * time.Millisecond)

	_, err := store.Validate(pair.RefreshToken)
	if err == nil {
		t.Error("expected expired token to be invalid")
	}
}

func TestMultipleTokensPerUser(t *testing.T) {
	store := NewRefreshTokenStore()
	pair1 := store.Issue("user-abc", 15*time.Minute, 30*24*time.Hour)
	pair2 := store.Issue("user-abc", 15*time.Minute, 30*24*time.Hour)

	if _, err := store.Validate(pair1.RefreshToken); err != nil {
		t.Errorf("pair1 validation: %v", err)
	}
	if _, err := store.Validate(pair2.RefreshToken); err != nil {
		t.Errorf("pair2 validation: %v", err)
	}
}

func TestRotationReuseDetection(t *testing.T) {
	store := NewRefreshTokenStore()
	pair1 := store.Issue("user-abc", 15*time.Minute, 30*24*time.Hour)
	otherSession := store.Issue("user-abc", 15*time.Minute, 30*24*time.Hour)

	pair2, err := store.Rotate(pair1.RefreshToken, 15*time.Minute, 30*24*time.Hour)
	if err != nil {
		t.Fatalf("first Rotate failed: %v", err)
	}
	pair3, err := store.Rotate(pair2.RefreshToken, 15*time.Minute, 30*24*time.Hour)
	if err != nil {
		t.Fatalf("second Rotate failed: %v", err)
	}

	// Replaying an already-rotated token must fail and revoke its entire family.
	if _, err := store.Rotate(pair1.RefreshToken, 15*time.Minute, 30*24*time.Hour); err == nil {
		t.Fatal("expected error when reusing rotated refresh token")
	}

	if _, err := store.Validate(pair3.RefreshToken); err == nil {
		t.Error("expected replacement token to be revoked after family replay detection")
	}
	if _, err := store.Rotate(pair3.RefreshToken, 15*time.Minute, 30*24*time.Hour); err == nil {
		t.Error("expected replacement token rotation to fail after family revocation")
	}

	// Independent sessions for the same user belong to distinct token families.
	if _, err := store.Validate(otherSession.RefreshToken); err != nil {
		t.Errorf("expected independent session family to stay valid, got %v", err)
	}
}
