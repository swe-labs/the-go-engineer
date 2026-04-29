// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// ============================================================================
// Section 06: Backend, APIs & Databases
// Title: Web Masterclass - Sessions
// Level: Advanced
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - How to manage state across requests using Cookies and Sessions.
//   - How to set secure, HttpOnly cookies in Go.
//   - How to implement a simple in-memory session store.
//   - The difference between "Stateful Sessions" and "Stateless Tokens".
//
// WHY THIS MATTERS:
//   - HTTP is a stateless protocol. To build features like login, shopping
//     carts, or user preferences, you need a way to "remember" who a user
//     is from one request to the next. Sessions are the standard way to
//     bridge this gap safely.
//
// RUN:
//   go run ./06-backend-db/01-web-and-database/web-masterclass/5-sessions
//
// KEY TAKEAWAY:
//   - Never store sensitive data directly in a cookie. Store a random ID
//     and keep the real data on the server.
// ============================================================================

package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"sync"
)

// Stage 06: Web Masterclass - Sessions
//
//   - Cookie Security: HttpOnly, Secure, SameSite
//   - Session Storage: In-memory map (Thread-safe)
//   - Authentication Flow: Login -> Session -> Access
//
// ENGINEERING DEPTH:
//   A session cookie should never contain actual user data (like
//   their email or balance). Instead, it should contain a long,
//   randomly generated string (a "Session ID"). The server uses
//   this ID as a key to look up the user's data in its own memory
//   or a database like Redis. This prevents users from "tampering"
//   with their session data in the browser.

// application holds our in-memory session store.
type application struct {
	sessionStore map[string]string
	mu           sync.RWMutex
}

func main() {
	app := &application{
		sessionStore: make(map[string]string),
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /", app.handleHome)
	mux.HandleFunc("GET /login", app.handleLogin)
	mux.HandleFunc("GET /secret", app.handleSecret)

	fmt.Println("=== Web Masterclass: Sessions ===")
	fmt.Println("  🚀 Server starting on http://localhost:8084")
	fmt.Println()
	fmt.Println("  1. Visit http://localhost:8084/secret (Access Denied)")
	fmt.Println("  2. Visit http://localhost:8084/login (Login)")
	fmt.Println("  3. Visit http://localhost:8084/secret (Access Granted)")

	log.Fatal(http.ListenAndServe(":8084", mux))

	fmt.Println("\n---------------------------------------------------")
	fmt.Println("NEXT UP: MC.6 authentication")
	fmt.Println("Current: MC.5 (sessions)")
	fmt.Println("Previous: MC.4 (middleware)")
	fmt.Println("---------------------------------------------------")
}

func (app *application) handleHome(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome to the Session Masterclass.")
}

func (app *application) handleLogin(w http.ResponseWriter, r *http.Request) {
	// 1. Generate a secure, random Session ID
	b := make([]byte, 16)
	rand.Read(b)
	sessionID := hex.EncodeToString(b)

	// 2. Store it in our server-side map
	app.mu.Lock()
	app.sessionStore[sessionID] = "Alice" // We "remember" this ID belongs to Alice
	app.mu.Unlock()

	// 3. Set a secure cookie in the response
	cookie := &http.Cookie{
		Name:     "session_token",
		Value:    sessionID,
		Path:     "/",
		HttpOnly: true,                 // Important: Prevent JavaScript access (XSS)
		Secure:   false,                // Set to true in production (HTTPS only)
		SameSite: http.SameSiteLaxMode, // CSRF protection
	}
	http.SetCookie(w, cookie)

	fmt.Fprintln(w, "Successfully logged in! A session cookie has been set.")
}

func (app *application) handleSecret(w http.ResponseWriter, r *http.Request) {
	// 1. Read the cookie from the request
	cookie, err := r.Cookie("session_token")
	if err != nil {
		http.Error(w, "Unauthorized: No session cookie found.", http.StatusUnauthorized)
		return
	}

	// 2. Look up the Session ID in our store
	app.mu.RLock()
	username, ok := app.sessionStore[cookie.Value]
	app.mu.RUnlock()

	if !ok {
		http.Error(w, "Unauthorized: Invalid session.", http.StatusUnauthorized)
		return
	}

	// 3. Access granted!
	fmt.Fprintf(w, "Welcome to the Secret Dashboard, %s!", username)
}
