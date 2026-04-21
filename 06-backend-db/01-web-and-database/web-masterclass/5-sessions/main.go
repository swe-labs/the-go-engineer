// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0
// Commercial use is prohibited without permission.

package main

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

// ============================================================================
// Stage 06: Web Masterclass — Sessions
// Level: Advanced
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Cookie-based session management
//   - Secure session tokens
//   - Session data storage (in-memory for learning)
//   - Flash messages pattern
//
// ENGINEERING DEPTH:
//   A massive architectural debate in backend engineering is "Stateful Sessions"
//   vs "Stateless JWTs". Because HTTP is inherently stateless, the server must
//   remember who you are. JWTs (Stateless) offload the memory to the client's
//   browser, requiring zero database lookups, but they CANNOT be instantly
//   revoked if a hacker steals the token. Stateful Sessions (like this file)
//   store a random string in the cookie and keep the actual data in Redis/Memory
//   on the server. This requires a fast database lookup on every single request,
//   but allows the server to instantly kill a session (logout) with zero latency.
//
// IMPORTANT: This uses a simple in-memory store for learning.
// In production, use a library like gorilla/sessions or scs.
//
// RUN: go run ./06-backend-db/01-web-and-database/web-masterclass/5-sessions
// ============================================================================

// Session stores data for a single user session.
type Session struct {
	Data      map[string]any
	CreatedAt time.Time
	ExpiresAt time.Time
}

// SessionStore manages all active sessions in memory.
// In production, use Redis, PostgreSQL, or another persistent store.
type SessionStore struct {
	mu       sync.RWMutex
	sessions map[string]*Session
}

func NewSessionStore() *SessionStore {
	return &SessionStore{
		sessions: make(map[string]*Session),
	}
}

// Create generates a new session and returns its ID.
func (s *SessionStore) Create() (string, *Session) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// In production, use crypto/rand for secure token generation.
	b := make([]byte, 16)
	rand.Read(b)
	id := hex.EncodeToString(b)
	sess := &Session{
		Data:      make(map[string]any),
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}
	s.sessions[id] = sess
	return id, sess
}

// Get retrieves a session by ID.
func (s *SessionStore) Get(id string) (*Session, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	sess, ok := s.sessions[id]
	if !ok || time.Now().After(sess.ExpiresAt) {
		return nil, false
	}
	return sess, true
}

type application struct {
	sessions *SessionStore
}

func main() {
	app := &application{
		sessions: NewSessionStore(),
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /", app.handleHome)
	mux.HandleFunc("POST /login", app.handleLogin)
	mux.HandleFunc("GET /dashboard", app.handleDashboard)
	mux.HandleFunc("GET /flash", app.handleFlash)

	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
	fmt.Println("\n---------------------------------------------------")
	fmt.Println("🚀 NEXT UP: WM.6 authentication")
	fmt.Println("   Current: WM.5 (sessions)")
	fmt.Println("---------------------------------------------------")
}

func (app *application) handleHome(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Visit POST /login to create a session")
}

func (app *application) handleLogin(w http.ResponseWriter, r *http.Request) {
	// Create a new session
	id, sess := app.sessions.Create()
	sess.Data["username"] = "gopher"
	sess.Data["flash"] = "Welcome back, gopher!"

	// Set the session cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    id,
		Path:     "/",
		HttpOnly: true,                 // JavaScript can't access this cookie
		Secure:   false,                // Set true in production (HTTPS only)
		SameSite: http.SameSiteLaxMode, // CSRF protection
		MaxAge:   86400,                // 24 hours
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message":    "Login successful",
		"session_id": id,
	})
}

func (app *application) handleDashboard(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		http.Error(w, "Not logged in", http.StatusUnauthorized)
		return
	}

	sess, ok := app.sessions.Get(cookie.Value)
	if !ok {
		http.Error(w, "Session expired", http.StatusUnauthorized)
		return
	}

	username := sess.Data["username"]
	fmt.Fprintf(w, "Welcome to your dashboard, %s!\n", username)
}

// handleFlash demonstrates the flash message pattern.
// Flash messages are shown once and then deleted from the session.
func (app *application) handleFlash(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		http.Error(w, "No session", http.StatusUnauthorized)
		return
	}

	sess, ok := app.sessions.Get(cookie.Value)
	if !ok {
		http.Error(w, "Session expired", http.StatusUnauthorized)
		return
	}

	// Read and delete flash message (one-time display)
	if flash, exists := sess.Data["flash"]; exists {
		fmt.Fprintf(w, "Flash message: %s\n", flash)
		delete(sess.Data, "flash") // Remove after reading
	} else {
		fmt.Fprintln(w, "No flash messages")
	}
}
