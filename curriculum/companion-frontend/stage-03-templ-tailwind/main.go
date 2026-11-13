package main

import (
	"crypto/rand"
	"crypto/subtle"
	"embed"
	"encoding/hex"
	"io/fs"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/swe-labs/the-go-engineer/curriculum/companion-frontend/stage-03-templ-tailwind/components"
)

//go:embed static/*
var defaultStaticFS embed.FS

// User represents an authenticated identity. Password is stored as a bcrypt
// hash — the same credential-at-rest model Module 10 teaches.
type User struct {
	ID           string
	Email        string
	PasswordHash string
	Name         string
	Role         string // "admin" or "viewer"
	TenantName   string
}

// Session represents an active browser login session.
type Session struct {
	ID        string
	UserID    string
	CSRFToken string
	ExpiresAt time.Time
}

// TenantConfig stores workspace isolation parameters.
type TenantConfig struct {
	Name            string
	StrictIsolation bool
}

// AuthServer coordinates session authentication, CSRF, and templ UI rendering.
type AuthServer struct {
	mu           sync.RWMutex
	users        map[string]*User
	sessions     map[string]*Session
	loginTokens  map[string]issuedLoginToken
	tenantConfig TenantConfig
	staticFS     fs.FS
}

// NewAuthServer initializes the authentication and tenant management service.
func NewAuthServer(staticFS fs.FS) *AuthServer {
	s := &AuthServer{
		users:       make(map[string]*User),
		sessions:    make(map[string]*Session),
		loginTokens: make(map[string]issuedLoginToken),
		tenantConfig: TenantConfig{
			Name:            "Acme Cloud Infrastructure",
			StrictIsolation: true,
		},
		staticFS: staticFS,
	}

	// Seed accounts: 1 Admin and 1 Viewer. Demo passwords are hashed with a
	// deliberately low bcrypt cost so test suites stay fast; production code
	// should use bcrypt.DefaultCost or higher.
	for _, seed := range []struct {
		email  string
		id     string
		name   string
		role   string
		tenant string
	}{
		{"alice@acme.corp", "usr-admin-01", "Alice Admin", "admin", "Acme Corp"},
		{"bob@acme.corp", "usr-viewer-02", "Bob Viewer", "viewer", "Acme Corp"},
	} {
		hash, err := bcrypt.GenerateFromPassword([]byte("opslane123"), bcrypt.MinCost)
		if err != nil {
			log.Fatalf("failed to hash seed password: %v", err)
		}
		s.users[seed.email] = &User{
			ID:           seed.id,
			Email:        seed.email,
			PasswordHash: string(hash),
			Name:         seed.name,
			Role:         seed.role,
			TenantName:   seed.tenant,
		}
	}

	return s
}

func generateRandomToken(bytesLen int) string {
	b := make([]byte, bytesLen)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

// dummyHash is a fixed bcrypt hash used when the submitted email does not
// exist, so failed logins run the same expensive hash verification regardless
// of whether the account exists (user-enumeration timing defense).
var dummyHash = func() []byte {
	h, err := bcrypt.GenerateFromPassword([]byte("opslane-dummy-password"), bcrypt.MinCost)
	if err != nil {
		log.Fatalf("failed to generate dummy hash: %v", err)
	}
	return h
}()

// Handler registers all application routes and middleware.
func (s *AuthServer) Handler() http.Handler {
	mux := http.NewServeMux()

	// Static assets
	staticSub, err := fs.Sub(s.staticFS, "static")
	if err == nil {
		mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(staticSub))))
	}

	// Authentication Endpoints
	mux.HandleFunc("/login", s.handleLoginView)
	mux.HandleFunc("/auth/login", s.handleLoginSubmit)
	mux.HandleFunc("/auth/logout", s.handleLogout)

	// Protected Console Views
	mux.HandleFunc("/profile", s.requireAuth(s.handleProfile))
	mux.HandleFunc("/settings", s.requireAuth(s.handleSettings))

	// Root redirect
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		if s.getSession(r) != nil {
			http.Redirect(w, r, "/profile", http.StatusFound)
			return
		}
		http.Redirect(w, r, "/login", http.StatusFound)
	})

	return mux
}

// getSession inspects the incoming request for a valid auth_session cookie.
func (s *AuthServer) getSession(r *http.Request) *Session {
	cookie, err := r.Cookie("auth_session")
	if err != nil {
		return nil
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	session, exists := s.sessions[cookie.Value]
	if !exists || time.Now().After(session.ExpiresAt) {
		return nil
	}
	return session
}

// requireAuth ensures the request has an active session before invoking the handler.
func (s *AuthServer) requireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session := s.getSession(r)
		if session == nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		next(w, r)
	}
}

// validateCSRF ensures POST submissions include the session's active CSRF token.
// The comparison is constant-time so token bytes cannot be probed by timing.
func (s *AuthServer) validateCSRF(r *http.Request, session *Session) bool {
	token := r.FormValue("csrf_token")
	if token == "" {
		token = r.Header.Get("X-CSRF-Token")
	}
	if token == "" {
		return false
	}
	return subtle.ConstantTimeCompare([]byte(token), []byte(session.CSRFToken)) == 1
}

// loginTokenTTL bounds how long an issued login-form token stays valid.
const loginTokenTTL = 5 * time.Minute

// anonymousLoginCookie identifies the browser across the unauthenticated
// login flow so issued synchronizer tokens are bound to the browser that
// requested them — the classic login-CSRF defense. Two tabs of the same
// browser share the cookie, so both tabs' tokens stay valid.
const anonymousLoginCookie = "login_anon"

// issuedLoginToken couples a form token to the anonymous browser identity
// it was issued for.
type issuedLoginToken struct {
	anonID    string
	expiresAt time.Time
}

// issueLoginToken registers a fresh single-use synchronizer token bound to
// the anonymous browser identity. Tokens live server-side (never embedded in
// the cookie itself) so a token captured by an attacker is useless without
// the victim's matching anonymous cookie.
func (s *AuthServer) issueLoginToken(w http.ResponseWriter, r *http.Request) string {
	anonID := ""
	if c, err := r.Cookie(anonymousLoginCookie); err == nil && c.Value != "" {
		anonID = c.Value
	} else {
		anonID = generateRandomToken(16)
	}
	// Re-set on every view so the cookie outlives long-lived tabs.
	http.SetCookie(w, &http.Cookie{
		Name:     anonymousLoginCookie,
		Value:    anonID,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   86400,
	})

	token := generateRandomToken(16)
	now := time.Now()

	s.mu.Lock()
	for t, issued := range s.loginTokens {
		if now.After(issued.expiresAt) {
			delete(s.loginTokens, t)
		}
	}
	s.loginTokens[token] = issuedLoginToken{anonID: anonID, expiresAt: now.Add(loginTokenTTL)}
	s.mu.Unlock()

	return token
}

// consumeLoginToken validates and burns a login-form token (single-use).
// The token must have been issued for the anonymous identity presented by
// the calling browser; a token captured from a different browser is refused.
func (s *AuthServer) consumeLoginToken(token, anonID string) bool {
	if token == "" {
		return false
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	issued, ok := s.loginTokens[token]
	if !ok || time.Now().After(issued.expiresAt) {
		return false
	}
	if subtle.ConstantTimeCompare([]byte(issued.anonID), []byte(anonID)) != 1 {
		return false
	}
	delete(s.loginTokens, token)
	return true
}

// handleLoginView displays the sign-in form.
func (s *AuthServer) handleLoginView(w http.ResponseWriter, r *http.Request) {
	if s.getSession(r) != nil {
		http.Redirect(w, r, "/profile", http.StatusFound)
		return
	}

	flash := r.URL.Query().Get("error")
	csrf := s.issueLoginToken(w, r)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_ = components.Login(flash, csrf).Render(r.Context(), w)
}

// handleLoginSubmit processes credential verification. The login form is
// CSRF-protected with a single-use synchronizer token bound to the calling
// browser's anonymous identity (issued by handleLoginView); the submitted
// token must exist, be unexpired, and match that identity.
func (s *AuthServer) handleLoginSubmit(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 1. CSRF verification for the unauthenticated login form
	anonID := ""
	if c, err := r.Cookie(anonymousLoginCookie); err == nil {
		anonID = c.Value
	}
	if !s.consumeLoginToken(r.FormValue("csrf_token"), anonID) {
		http.Redirect(w, r, "/login?error=Login+session+expired%2C+please+retry", http.StatusSeeOther)
		return
	}

	email := r.FormValue("email")
	password := r.FormValue("password")

	s.mu.RLock()
	user, exists := s.users[email]
	s.mu.RUnlock()

	// bcrypt.CompareHashAndPassword runs the full hash verification even for a
	// missing user (via a fixed dummy hash), so response timing does not reveal
	// whether the email exists.
	var hash []byte
	if exists {
		hash = []byte(user.PasswordHash)
	} else {
		hash = dummyHash
	}
	if bcrypt.CompareHashAndPassword(hash, []byte(password)) != nil || !exists {
		http.Redirect(w, r, "/login?error=Invalid+credentials+provided", http.StatusSeeOther)
		return
	}

	// Create new authenticated session
	sessionID := generateRandomToken(32)
	csrfToken := generateRandomToken(32)

	s.mu.Lock()
	s.sessions[sessionID] = &Session{
		ID:        sessionID,
		UserID:    user.ID,
		CSRFToken: csrfToken,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}
	s.mu.Unlock()

	http.SetCookie(w, &http.Cookie{
		Name:     "auth_session",
		Value:    sessionID,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().Add(24 * time.Hour),
	})

	http.Redirect(w, r, "/profile", http.StatusSeeOther)
}

// handleLogout clears the session.
func (s *AuthServer) handleLogout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	cookie, err := r.Cookie("auth_session")
	if err == nil {
		s.mu.Lock()
		delete(s.sessions, cookie.Value)
		s.mu.Unlock()
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "auth_session",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1,
	})

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

// handleProfile handles viewing and updating personal profile info.
func (s *AuthServer) handleProfile(w http.ResponseWriter, r *http.Request) {
	session := s.getSession(r)
	s.mu.RLock()
	var user *User
	for _, u := range s.users {
		if u.ID == session.UserID {
			user = u
			break
		}
	}
	s.mu.RUnlock()

	if user == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	flashSuccess := ""
	flashError := ""

	if r.Method == http.MethodPost {
		if !s.validateCSRF(r, session) {
			http.Error(w, "403 Forbidden: Invalid or missing CSRF token", http.StatusForbidden)
			return
		}

		newName := r.FormValue("name")
		if len(newName) < 2 {
			flashError = "Name must be at least 2 characters long."
		} else {
			s.mu.Lock()
			user.Name = newName
			s.mu.Unlock()
			flashSuccess = "Profile successfully updated!"
		}
	}

	vm := components.ProfileViewModel{
		User: components.UserContext{
			ID:         user.ID,
			Email:      user.Email,
			Name:       user.Name,
			Role:       user.Role,
			TenantName: user.TenantName,
		},
		CSRFToken:    session.CSRFToken,
		FlashSuccess: flashSuccess,
		FlashError:   flashError,
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_ = components.Profile(vm).Render(r.Context(), w)
}

// handleSettings manages workspace and tenant parameters with role guards.
func (s *AuthServer) handleSettings(w http.ResponseWriter, r *http.Request) {
	session := s.getSession(r)
	s.mu.RLock()
	var user *User
	for _, u := range s.users {
		if u.ID == session.UserID {
			user = u
			break
		}
	}
	currentTenant := s.tenantConfig
	s.mu.RUnlock()

	if user == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	flashSuccess := ""
	flashError := ""

	if r.Method == http.MethodPost {
		// 1. CSRF Verification
		if !s.validateCSRF(r, session) {
			http.Error(w, "403 Forbidden: Invalid or missing CSRF token", http.StatusForbidden)
			return
		}

		// 2. Role-Based Access Control Guard: Only Admin may modify tenant settings
		if user.Role != "admin" {
			http.Error(w, "403 Forbidden: Insufficient privileges. Admin role required.", http.StatusForbidden)
			return
		}

		newName := r.FormValue("tenant_name")
		strictIsolation := r.FormValue("strict_isolation") == "true"

		if len(newName) < 3 {
			flashError = "Organization name must be at least 3 characters."
		} else {
			s.mu.Lock()
			s.tenantConfig.Name = newName
			s.tenantConfig.StrictIsolation = strictIsolation
			currentTenant = s.tenantConfig
			s.mu.Unlock()
			flashSuccess = "Tenant configuration successfully updated!"
		}
	}

	vm := components.SettingsViewModel{
		User: components.UserContext{
			ID:         user.ID,
			Email:      user.Email,
			Name:       user.Name,
			Role:       user.Role,
			TenantName: user.TenantName,
		},
		TenantName:   currentTenant.Name,
		StrictMulti:  currentTenant.StrictIsolation,
		CSRFToken:    session.CSRFToken,
		FlashSuccess: flashSuccess,
		FlashError:   flashError,
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_ = components.Settings(vm).Render(r.Context(), w)
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	server := NewAuthServer(defaultStaticFS)

	addr := ":" + port
	log.Printf("[Stage 3 Starter] Opslane Console listening at http://localhost%s", addr)
	srv := &http.Server{
		Addr:         addr,
		Handler:      server.Handler(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server stopped: %v", err)
	}
}
