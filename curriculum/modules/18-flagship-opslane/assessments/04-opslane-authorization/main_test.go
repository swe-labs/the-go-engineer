package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func setupTestRouter() http.Handler {
	am := NewAuthMiddleware()
	return BuildRouter(am)
}

func TestAuthorization_Unauthenticated(t *testing.T) {
	router := setupTestRouter()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/orders", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 Unauthorized for missing token, got %d", rr.Code)
	}
}

func TestAuthorization_ViewerPermissions(t *testing.T) {
	router := setupTestRouter()

	// 1. Viewer reading orders -> 200 OK
	reqRead := httptest.NewRequest(http.MethodGet, "/api/v1/orders", nil)
	reqRead.Header.Set("Authorization", "Bearer token-viewer")
	rrRead := httptest.NewRecorder()
	router.ServeHTTP(rrRead, reqRead)

	if rrRead.Code != http.StatusOK {
		t.Fatalf("expected Viewer to have read permissions (200 OK), got %d", rrRead.Code)
	}

	// 2. Viewer attempting to create an order -> 403 Forbidden
	reqCreate := httptest.NewRequest(http.MethodPost, "/api/v1/orders", strings.NewReader(`{}`))
	reqCreate.Header.Set("Authorization", "Bearer token-viewer")
	rrCreate := httptest.NewRecorder()
	router.ServeHTTP(rrCreate, reqCreate)

	if rrCreate.Code != http.StatusForbidden {
		t.Fatalf("expected 403 Forbidden when Viewer attempts order creation, got %d", rrCreate.Code)
	}
}

func TestAuthorization_MemberPermissions(t *testing.T) {
	router := setupTestRouter()

	// 1. Member creating order -> 201 Created
	reqCreate := httptest.NewRequest(http.MethodPost, "/api/v1/orders", strings.NewReader(`{}`))
	reqCreate.Header.Set("Authorization", "Bearer token-member")
	rrCreate := httptest.NewRecorder()
	router.ServeHTTP(rrCreate, reqCreate)

	if rrCreate.Code != http.StatusCreated {
		t.Fatalf("expected Member to create order (201 Created), got %d", rrCreate.Code)
	}

	// 2. Member modifying tenant settings -> 403 Forbidden
	reqSettings := httptest.NewRequest(http.MethodPost, "/api/v1/tenant/settings", strings.NewReader(`{}`))
	reqSettings.Header.Set("Authorization", "Bearer token-member")
	rrSettings := httptest.NewRecorder()
	router.ServeHTTP(rrSettings, reqSettings)

	if rrSettings.Code != http.StatusForbidden {
		t.Fatalf("expected 403 Forbidden when Member attempts tenant setting modification, got %d", rrSettings.Code)
	}
}

func TestAuthorization_OwnerPermissions(t *testing.T) {
	router := setupTestRouter()

	reqSettings := httptest.NewRequest(http.MethodPost, "/api/v1/tenant/settings", strings.NewReader(`{}`))
	reqSettings.Header.Set("Authorization", "Bearer token-owner")
	rrSettings := httptest.NewRecorder()
	router.ServeHTTP(rrSettings, reqSettings)

	if rrSettings.Code != http.StatusOK {
		t.Fatalf("expected Owner to have tenant management permissions (200 OK), got %d", rrSettings.Code)
	}
}
