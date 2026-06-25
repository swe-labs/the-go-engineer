package main

import (
	"testing"
)

func TestAuthN_ValidToken(t *testing.T) {
	claims, err := AuthN("admin-token")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if claims.UserID != "u1" {
		t.Errorf("expected u1, got %s", claims.UserID)
	}
	if claims.Role != "admin" {
		t.Errorf("expected admin, got %s", claims.Role)
	}
}

func TestAuthN_InvalidToken(t *testing.T) {
	_, err := AuthN("invalid-token")
	if err == nil {
		t.Fatal("expected error for invalid token")
	}
}

func TestAuthZ_AdminAnyResource(t *testing.T) {
	claims := Claims{UserID: "u1", Role: "admin", Perms: []string{"doc:*"}}
	if !AuthZ(claims, "doc:delete", "any-resource") {
		t.Error("expected admin to have access to any resource")
	}
}

func TestAuthZ_OwnerWithPerm(t *testing.T) {
	claims := Claims{UserID: "u2", Role: "user", Perms: []string{"doc:read"}}
	if !AuthZ(claims, "doc:read", "u2") {
		t.Error("expected owner with read perm to have access")
	}
}

func TestAuthZ_NonOwnerNoAdmin(t *testing.T) {
	claims := Claims{UserID: "u2", Role: "user", Perms: []string{"doc:read"}}
	if AuthZ(claims, "doc:delete", "u1") {
		t.Error("expected non-owner without delete perm to be denied")
	}
}

func TestHasPermission_ExactMatch(t *testing.T) {
	c := Claims{Perms: []string{"doc:read"}}
	if !hasPermission(c, "doc:read") {
		t.Error("expected exact match")
	}
}

func TestHasPermission_Wildcard(t *testing.T) {
	c := Claims{Perms: []string{"doc:*"}}
	if !hasPermission(c, "doc:delete") {
		t.Error("expected wildcard match")
	}
}

func TestHasPermission_NoMatch(t *testing.T) {
	c := Claims{Perms: []string{"doc:read"}}
	if hasPermission(c, "user:manage") {
		t.Error("expected no match")
	}
}
