package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

var testSecret = []byte("test-signing-secret-key-32-chars-long")

func setupTestApp() (http.Handler, *MemoryStore) {
	store := NewMemoryStore()
	handler := NewAPIHandler(store, testSecret)
	return handler, store
}

func TestPasswordHashingAndJsonExclusion(t *testing.T) {
	_, store := setupTestApp()

	user, err := store.Register(RegisterRequest{
		TenantID: "tenant-1",
		Email:    "alice@example.com",
		Password: "supersecretpassword",
		Role:     RoleAdmin,
	})
	if err != nil {
		t.Fatalf("registration failed: %v", err)
	}

	if user.PasswordHash == "supersecretpassword" {
		t.Fatal("password was stored in plaintext")
	}
	if !CheckPassword(user.PasswordHash, "supersecretpassword") {
		t.Fatal("bcrypt verification failed for valid password")
	}

	// Verify JSON serialization omits password hash
	jsonBytes, err := json.Marshal(user)
	if err != nil {
		t.Fatalf("json marshal failed: %v", err)
	}
	if bytes.Contains(jsonBytes, []byte("password")) || bytes.Contains(jsonBytes, []byte(user.PasswordHash)) {
		t.Fatalf("password hash leaked into json serialization: %s", string(jsonBytes))
	}
}

func TestLoginAuthentication(t *testing.T) {
	handler, store := setupTestApp()

	_, _ = store.Register(RegisterRequest{
		TenantID: "tenant-1",
		Email:    "bob@example.com",
		Password: "password123",
		Role:     RoleMember,
	})

	// Valid login
	validBody, _ := json.Marshal(LoginRequest{Email: "bob@example.com", Password: "password123"})
	req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewReader(validBody))
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	var resp JSONResponse
	_ = json.Unmarshal(rec.Body.Bytes(), &resp)
	data := resp.Data.(map[string]any)
	tokenStr, ok := data["token"].(string)
	if !ok || tokenStr == "" {
		t.Fatal("expected JWT token string in login response")
	}

	// Invalid password
	invalidBody, _ := json.Marshal(LoginRequest{Email: "bob@example.com", Password: "wrongpassword"})
	reqBad := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewReader(invalidBody))
	recBad := httptest.NewRecorder()
	handler.ServeHTTP(recBad, reqBad)

	if recBad.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 for wrong password, got %d", recBad.Code)
	}
}

func TestNegativeAuthSuite(t *testing.T) {
	handler, store := setupTestApp()

	user, _ := store.Register(RegisterRequest{
		TenantID: "tenant-1",
		Email:    "charlie@example.com",
		Password: "password123",
		Role:     RoleMember,
	})

	// 1. Missing Authorization header
	reqNoAuth := httptest.NewRequest(http.MethodGet, "/auth/me", nil)
	recNoAuth := httptest.NewRecorder()
	handler.ServeHTTP(recNoAuth, reqNoAuth)
	if recNoAuth.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 for missing auth header, got %d", recNoAuth.Code)
	}

	// 2. Expired JWT
	expiredToken, _ := GenerateJWT(user, testSecret, -1*time.Hour)
	reqExpired := httptest.NewRequest(http.MethodGet, "/auth/me", nil)
	reqExpired.Header.Set("Authorization", "Bearer "+expiredToken)
	recExpired := httptest.NewRecorder()
	handler.ServeHTTP(recExpired, reqExpired)
	if recExpired.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 for expired token, got %d", recExpired.Code)
	}

	// 3. Tampered JWT signature
	validToken, _ := GenerateJWT(user, testSecret, 1*time.Hour)
	tamperedToken := validToken[:len(validToken)-4] + "fake"
	reqTampered := httptest.NewRequest(http.MethodGet, "/auth/me", nil)
	reqTampered.Header.Set("Authorization", "Bearer "+tamperedToken)
	recTampered := httptest.NewRecorder()
	handler.ServeHTTP(recTampered, reqTampered)
	if recTampered.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 for tampered token, got %d", recTampered.Code)
	}
}

func TestRBACHierarchy(t *testing.T) {
	handler, store := setupTestApp()

	viewer, _ := store.Register(RegisterRequest{TenantID: "tenant-1", Email: "viewer@example.com", Password: "password123", Role: RoleViewer})
	member, _ := store.Register(RegisterRequest{TenantID: "tenant-1", Email: "member@example.com", Password: "password123", Role: RoleMember})
	admin, _ := store.Register(RegisterRequest{TenantID: "tenant-1", Email: "admin@example.com", Password: "password123", Role: RoleAdmin})

	viewerToken, _ := GenerateJWT(viewer, testSecret, 1*time.Hour)
	memberToken, _ := GenerateJWT(member, testSecret, 1*time.Hour)
	adminToken, _ := GenerateJWT(admin, testSecret, 1*time.Hour)

	// Viewer attempts to POST /documents (requires MEMBER) -> 403
	docBody, _ := json.Marshal(CreateDocumentRequest{Title: "Doc 1", Content: "Secret Content"})
	reqViewerPost := httptest.NewRequest(http.MethodPost, "/documents", bytes.NewReader(docBody))
	reqViewerPost.Header.Set("Authorization", "Bearer "+viewerToken)
	recViewerPost := httptest.NewRecorder()
	handler.ServeHTTP(recViewerPost, reqViewerPost)
	if recViewerPost.Code != http.StatusForbidden {
		t.Fatalf("expected 403 Forbidden for Viewer creating document, got %d", recViewerPost.Code)
	}

	// Member creates document -> 201 Created
	reqMemberPost := httptest.NewRequest(http.MethodPost, "/documents", bytes.NewReader(docBody))
	reqMemberPost.Header.Set("Authorization", "Bearer "+memberToken)
	recMemberPost := httptest.NewRecorder()
	handler.ServeHTTP(recMemberPost, reqMemberPost)
	if recMemberPost.Code != http.StatusCreated {
		t.Fatalf("expected 201 Created for Member creating document, got %d", recMemberPost.Code)
	}

	var resp JSONResponse
	_ = json.Unmarshal(recMemberPost.Body.Bytes(), &resp)
	raw, _ := json.Marshal(resp.Data)
	var doc Document
	_ = json.Unmarshal(raw, &doc)

	// Member attempts to DELETE document (requires ADMIN) -> 403
	reqMemberDel := httptest.NewRequest(http.MethodDelete, "/documents/"+doc.ID, nil)
	reqMemberDel.Header.Set("Authorization", "Bearer "+memberToken)
	recMemberDel := httptest.NewRecorder()
	handler.ServeHTTP(recMemberDel, reqMemberDel)
	if recMemberDel.Code != http.StatusForbidden {
		t.Fatalf("expected 403 Forbidden for Member deleting document, got %d", recMemberDel.Code)
	}

	// Admin deletes document -> 204 No Content
	reqAdminDel := httptest.NewRequest(http.MethodDelete, "/documents/"+doc.ID, nil)
	reqAdminDel.Header.Set("Authorization", "Bearer "+adminToken)
	recAdminDel := httptest.NewRecorder()
	handler.ServeHTTP(recAdminDel, reqAdminDel)
	if recAdminDel.Code != http.StatusNoContent {
		t.Fatalf("expected 204 No Content for Admin deleting document, got %d", recAdminDel.Code)
	}
}

func TestCrossTenantIsolation(t *testing.T) {
	handler, store := setupTestApp()

	// Tenant A user
	userA, _ := store.Register(RegisterRequest{TenantID: "tenant-A", Email: "userA@a.com", Password: "password123", Role: RoleMember})
	tokenA, _ := GenerateJWT(userA, testSecret, 1*time.Hour)

	// Tenant B user (even as Admin of Tenant B)
	adminB, _ := store.Register(RegisterRequest{TenantID: "tenant-B", Email: "adminB@b.com", Password: "password123", Role: RoleAdmin})
	tokenB, _ := GenerateJWT(adminB, testSecret, 1*time.Hour)

	// User A creates doc in Tenant A
	docA, _ := store.CreateDocument("tenant-A", userA.ID, "Alpha Specs", "Alpha Confidential")

	// User A successfully accesses docA (same tenant)
	reqOwnerGet := httptest.NewRequest(http.MethodGet, "/documents/"+docA.ID, nil)
	reqOwnerGet.Header.Set("Authorization", "Bearer "+tokenA)
	recOwnerGet := httptest.NewRecorder()
	handler.ServeHTTP(recOwnerGet, reqOwnerGet)
	if recOwnerGet.Code != http.StatusOK {
		t.Fatalf("expected 200 OK for tenant-A owner accessing docA, got %d", recOwnerGet.Code)
	}

	// User B attempts to access docA -> Must receive 403 Forbidden
	reqCrossGet := httptest.NewRequest(http.MethodGet, "/documents/"+docA.ID, nil)
	reqCrossGet.Header.Set("Authorization", "Bearer "+tokenB)
	recCrossGet := httptest.NewRecorder()
	handler.ServeHTTP(recCrossGet, reqCrossGet)

	if recCrossGet.Code != http.StatusForbidden {
		t.Fatalf("expected 403 Forbidden for cross-tenant access attempt, got %d", recCrossGet.Code)
	}

	// User B attempts to DELETE docA -> Must receive 403 Forbidden
	reqCrossDel := httptest.NewRequest(http.MethodDelete, "/documents/"+docA.ID, nil)
	reqCrossDel.Header.Set("Authorization", "Bearer "+tokenB)
	recCrossDel := httptest.NewRecorder()
	handler.ServeHTTP(recCrossDel, reqCrossDel)

	if recCrossDel.Code != http.StatusForbidden {
		t.Fatalf("expected 403 Forbidden for cross-tenant delete attempt, got %d", recCrossDel.Code)
	}

	// Verify docA still exists
	stillThere, err := store.GetDocument(docA.ID)
	if err != nil || stillThere.Title != "Alpha Specs" {
		t.Fatal("cross-tenant delete improperly affected tenant-A document")
	}
}

func TestSecurityHeaders(t *testing.T) {
	handler, _ := setupTestApp()

	req := httptest.NewRequest(http.MethodGet, "/auth/me", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Header().Get("X-Content-Type-Options") != "nosniff" {
		t.Fatalf("expected X-Content-Type-Options: nosniff, got %q", rec.Header().Get("X-Content-Type-Options"))
	}
	if rec.Header().Get("X-Frame-Options") != "DENY" {
		t.Fatalf("expected X-Frame-Options: DENY, got %q", rec.Header().Get("X-Frame-Options"))
	}
}
