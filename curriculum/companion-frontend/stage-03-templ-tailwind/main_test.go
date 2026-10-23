package main

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func setupTestAuthServer(t *testing.T) *AuthServer {
	t.Helper()
	return NewAuthServer(defaultStaticFS)
}

func TestAuthServer_UnauthenticatedRedirect(t *testing.T) {
	server := setupTestAuthServer(t)

	req := httptest.NewRequest(http.MethodGet, "/profile", nil)
	rr := httptest.NewRecorder()
	server.Handler().ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Fatalf("expected redirect 303 See Other for unauthenticated /profile, got %d", rr.Code)
	}

	loc := rr.Header().Get("Location")
	if loc != "/login" {
		t.Errorf("expected redirect to /login, got %q", loc)
	}
}

func TestAuthServer_LoginFlow(t *testing.T) {
	server := setupTestAuthServer(t)

	// 1. Invalid credentials
	badForm := url.Values{
		"email":    {"alice@acme.corp"},
		"password": {"wrong-password"},
	}
	reqBad := httptest.NewRequest(http.MethodPost, "/auth/login", strings.NewReader(badForm.Encode()))
	reqBad.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rrBad := httptest.NewRecorder()
	server.Handler().ServeHTTP(rrBad, reqBad)

	if rrBad.Code != http.StatusSeeOther {
		t.Fatalf("expected 303 redirect on bad login, got %d", rrBad.Code)
	}
	if !strings.Contains(rrBad.Header().Get("Location"), "error=") {
		t.Errorf("expected error parameter in redirect location")
	}

	// 2. Valid credentials
	goodForm := url.Values{
		"email":    {"alice@acme.corp"},
		"password": {"opslane123"},
	}
	reqGood := httptest.NewRequest(http.MethodPost, "/auth/login", strings.NewReader(goodForm.Encode()))
	reqGood.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rrGood := httptest.NewRecorder()
	server.Handler().ServeHTTP(rrGood, reqGood)

	if rrGood.Code != http.StatusSeeOther {
		t.Fatalf("expected 303 redirect on successful login, got %d", rrGood.Code)
	}

	cookies := rrGood.Result().Cookies()
	var authCookie *http.Cookie
	for _, c := range cookies {
		if c.Name == "auth_session" {
			authCookie = c
			break
		}
	}

	if authCookie == nil || authCookie.Value == "" {
		t.Fatalf("expected non-empty auth_session cookie")
	}
}

func TestAuthServer_CSRFProtection(t *testing.T) {
	server := setupTestAuthServer(t)

	// Authenticate as Alice
	loginForm := url.Values{
		"email":    {"alice@acme.corp"},
		"password": {"opslane123"},
	}
	reqLogin := httptest.NewRequest(http.MethodPost, "/auth/login", strings.NewReader(loginForm.Encode()))
	reqLogin.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rrLogin := httptest.NewRecorder()
	server.Handler().ServeHTTP(rrLogin, reqLogin)

	authCookie := rrLogin.Result().Cookies()[0]

	// 1. POST /profile without CSRF token -> 403 Forbidden
	noCSRFForm := url.Values{
		"name": {"Alice Updated"},
	}
	reqNoCSRF := httptest.NewRequest(http.MethodPost, "/profile", strings.NewReader(noCSRFForm.Encode()))
	reqNoCSRF.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	reqNoCSRF.AddCookie(authCookie)
	rrNoCSRF := httptest.NewRecorder()
	server.Handler().ServeHTTP(rrNoCSRF, reqNoCSRF)

	if rrNoCSRF.Code != http.StatusForbidden {
		t.Fatalf("expected 403 Forbidden when CSRF token is omitted, got %d", rrNoCSRF.Code)
	}

	// 2. Retrieve valid CSRF token from active session
	session := server.sessions[authCookie.Value]
	if session == nil {
		t.Fatalf("session not found in server state")
	}

	// 3. POST /profile with valid CSRF token -> 200 OK & update succeeds
	validForm := url.Values{
		"name":       {"Alice Wonderland"},
		"csrf_token": {session.CSRFToken},
	}
	reqValid := httptest.NewRequest(http.MethodPost, "/profile", strings.NewReader(validForm.Encode()))
	reqValid.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	reqValid.AddCookie(authCookie)
	rrValid := httptest.NewRecorder()
	server.Handler().ServeHTTP(rrValid, reqValid)

	if rrValid.Code != http.StatusOK {
		t.Fatalf("expected 200 OK on profile update, got %d: %s", rrValid.Code, rrValid.Body.String())
	}

	body := rrValid.Body.String()
	if !strings.Contains(body, "Alice Wonderland") {
		t.Errorf("expected rendered page to reflect updated name 'Alice Wonderland'")
	}
	if !strings.Contains(body, "Profile successfully updated!") {
		t.Errorf("expected success notification banner")
	}
}

func TestAuthServer_RoleBasedAccessControl(t *testing.T) {
	server := setupTestAuthServer(t)

	// Authenticate as Bob (Viewer)
	loginForm := url.Values{
		"email":    {"bob@acme.corp"},
		"password": {"opslane123"},
	}
	reqLogin := httptest.NewRequest(http.MethodPost, "/auth/login", strings.NewReader(loginForm.Encode()))
	reqLogin.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rrLogin := httptest.NewRecorder()
	server.Handler().ServeHTTP(rrLogin, reqLogin)

	bobCookie := rrLogin.Result().Cookies()[0]
	bobSession := server.sessions[bobCookie.Value]

	// 1. GET /settings should display Read-Only Viewer banner
	reqGetSettings := httptest.NewRequest(http.MethodGet, "/settings", nil)
	reqGetSettings.AddCookie(bobCookie)
	rrGetSettings := httptest.NewRecorder()
	server.Handler().ServeHTTP(rrGetSettings, reqGetSettings)

	if rrGetSettings.Code != http.StatusOK {
		t.Fatalf("expected 200 OK for GET /settings, got %d", rrGetSettings.Code)
	}
	if !strings.Contains(rrGetSettings.Body.String(), "Read-Only Mode") {
		t.Errorf("expected Viewer to see 'Read-Only Mode' notice")
	}

	// 2. Viewer attempts POST /settings -> 403 Forbidden
	settingsForm := url.Values{
		"tenant_name": {"Hacked Tenant"},
		"csrf_token":  {bobSession.CSRFToken},
	}
	reqPostSettings := httptest.NewRequest(http.MethodPost, "/settings", strings.NewReader(settingsForm.Encode()))
	reqPostSettings.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	reqPostSettings.AddCookie(bobCookie)
	rrPostSettings := httptest.NewRecorder()
	server.Handler().ServeHTTP(rrPostSettings, reqPostSettings)

	if rrPostSettings.Code != http.StatusForbidden {
		t.Fatalf("expected 403 Forbidden when Viewer attempts to modify tenant settings, got %d", rrPostSettings.Code)
	}

	// 3. Admin (Alice) attempts POST /settings -> 200 OK
	adminLoginForm := url.Values{
		"email":    {"alice@acme.corp"},
		"password": {"opslane123"},
	}
	reqAdminLogin := httptest.NewRequest(http.MethodPost, "/auth/login", strings.NewReader(adminLoginForm.Encode()))
	reqAdminLogin.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rrAdminLogin := httptest.NewRecorder()
	server.Handler().ServeHTTP(rrAdminLogin, reqAdminLogin)

	aliceCookie := rrAdminLogin.Result().Cookies()[0]
	aliceSession := server.sessions[aliceCookie.Value]

	adminSettingsForm := url.Values{
		"tenant_name":      {"Acme Global Tech"},
		"strict_isolation": {"true"},
		"csrf_token":       {aliceSession.CSRFToken},
	}
	reqAdminSettings := httptest.NewRequest(http.MethodPost, "/settings", strings.NewReader(adminSettingsForm.Encode()))
	reqAdminSettings.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	reqAdminSettings.AddCookie(aliceCookie)
	rrAdminSettings := httptest.NewRecorder()
	server.Handler().ServeHTTP(rrAdminSettings, reqAdminSettings)

	if rrAdminSettings.Code != http.StatusOK {
		t.Fatalf("expected 200 OK for admin settings update, got %d", rrAdminSettings.Code)
	}
	if server.tenantConfig.Name != "Acme Global Tech" {
		t.Errorf("expected tenant name to be updated in server state")
	}
}

func TestAuthServer_Logout(t *testing.T) {
	server := setupTestAuthServer(t)

	// Login
	loginForm := url.Values{
		"email":    {"alice@acme.corp"},
		"password": {"opslane123"},
	}
	reqLogin := httptest.NewRequest(http.MethodPost, "/auth/login", strings.NewReader(loginForm.Encode()))
	reqLogin.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rrLogin := httptest.NewRecorder()
	server.Handler().ServeHTTP(rrLogin, reqLogin)

	authCookie := rrLogin.Result().Cookies()[0]

	// Logout
	reqLogout := httptest.NewRequest(http.MethodPost, "/auth/logout", nil)
	reqLogout.AddCookie(authCookie)
	rrLogout := httptest.NewRecorder()
	server.Handler().ServeHTTP(rrLogout, reqLogout)

	if rrLogout.Code != http.StatusSeeOther {
		t.Fatalf("expected 303 redirect on logout, got %d", rrLogout.Code)
	}

	// Verify session destroyed in server
	if _, exists := server.sessions[authCookie.Value]; exists {
		t.Errorf("expected session to be deleted from server state")
	}
}
