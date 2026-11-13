package main

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"regexp"
	"strings"
	"testing"
)

func setupTestAuthServer(t *testing.T) *AuthServer {
	t.Helper()
	return NewAuthServer(defaultStaticFS)
}

// loginAs performs the full CSRF-protected login flow: it fetches /login to
// obtain a server-side single-use synchronizer token, then submits the
// credentials and returns the auth_session cookie.
func loginAs(t *testing.T, server *AuthServer, email, password string) *http.Cookie {
	t.Helper()

	reqView := httptest.NewRequest(http.MethodGet, "/login", nil)
	rrView := httptest.NewRecorder()
	server.Handler().ServeHTTP(rrView, reqView)

	m := regexp.MustCompile(`name="csrf_token" value="([^"]+)"`).FindStringSubmatch(rrView.Body.String())
	if m == nil {
		t.Fatalf("expected hidden csrf_token field in rendered login form")
	}

	form := url.Values{
		"email":      {email},
		"password":   {password},
		"csrf_token": {m[1]},
	}
	reqLogin := httptest.NewRequest(http.MethodPost, "/auth/login", strings.NewReader(form.Encode()))
	reqLogin.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rrLogin := httptest.NewRecorder()
	server.Handler().ServeHTTP(rrLogin, reqLogin)

	if rrLogin.Code != http.StatusSeeOther {
		t.Fatalf("expected 303 redirect on login for %s, got %d: %s", email, rrLogin.Code, rrLogin.Body.String())
	}

	for _, c := range rrLogin.Result().Cookies() {
		if c.Name == "auth_session" && c.Value != "" {
			return c
		}
	}
	t.Fatalf("expected non-empty auth_session cookie for %s", email)
	return nil
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

	// 1. POST /auth/login with an unknown/absent synchronizer token is
	// rejected before credentials are even evaluated.
	form := url.Values{
		"email":      {"alice@acme.corp"},
		"password":   {"opslane123"},
		"csrf_token": {"never-issued-token"},
	}
	reqNoCSRF := httptest.NewRequest(http.MethodPost, "/auth/login", strings.NewReader(form.Encode()))
	reqNoCSRF.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rrNoCSRF := httptest.NewRecorder()
	server.Handler().ServeHTTP(rrNoCSRF, reqNoCSRF)

	if rrNoCSRF.Code != http.StatusSeeOther {
		t.Fatalf("expected 303 redirect on missing login CSRF token, got %d", rrNoCSRF.Code)
	}
	if !strings.Contains(rrNoCSRF.Header().Get("Location"), "expired") {
		t.Errorf("expected 'session expired' redirect location, got %q", rrNoCSRF.Header().Get("Location"))
	}

	// 2. Fetch the login view to obtain a valid token, then submit wrong
	// credentials -> error redirect and no session cookie.
	reqView := httptest.NewRequest(http.MethodGet, "/login", nil)
	rrView := httptest.NewRecorder()
	server.Handler().ServeHTTP(rrView, reqView)

	m := regexp.MustCompile(`name="csrf_token" value="([^"]+)"`).FindStringSubmatch(rrView.Body.String())
	if m == nil {
		t.Fatalf("expected hidden csrf_token field in rendered login form")
	}

	badForm := url.Values{
		"email":      {"alice@acme.corp"},
		"password":   {"wrong-password"},
		"csrf_token": {m[1]},
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
	for _, c := range rrBad.Result().Cookies() {
		if c.Name == "auth_session" && c.Value != "" {
			t.Fatalf("SECURITY VIOLATION: session cookie issued despite wrong password")
		}
	}

	// 3. Valid credentials issue a session cookie and consume the CSRF token.
	authCookie := loginAs(t, server, "alice@acme.corp", "opslane123")
	if authCookie == nil || authCookie.Value == "" {
		t.Fatalf("expected non-empty auth_session cookie")
	}
}

func TestAuthServer_CSRFProtection(t *testing.T) {
	server := setupTestAuthServer(t)

	// Authenticate as Alice via the full CSRF-protected login flow
	authCookie := loginAs(t, server, "alice@acme.corp", "opslane123")

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

	// Authenticate as Bob (Viewer) via the full login flow
	bobCookie := loginAs(t, server, "bob@acme.corp", "opslane123")
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
	aliceCookie := loginAs(t, server, "alice@acme.corp", "opslane123")
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

	// Login via the full CSRF-protected flow
	authCookie := loginAs(t, server, "alice@acme.corp", "opslane123")

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
