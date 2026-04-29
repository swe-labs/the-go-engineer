// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// ============================================================================
// Section 06: Backend, APIs & Databases
// Title: Web Masterclass - Authentication
// Level: Advanced
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - How to securely handle user passwords using bcrypt.
//   - The difference between hashing and encryption.
//   - How to use 'context.Context' to store user identity across middleware.
//   - How to protect routes using Authentication Middleware.
//
// WHY THIS MATTERS:
//   - Authentication is the gatekeeper of your application. Storing
//     passwords in plain text is a critical security failure. Using
//     industry-standard algorithms like 'bcrypt' ensures that even if
//     your database is compromised, your users' passwords remain safe.
//
// RUN:
//   go run ./06-backend-db/01-web-and-database/web-masterclass/6-auth
//
// KEY TAKEAWAY:
//   - Never roll your own crypto. Use trusted, battle-tested libraries.
// ============================================================================

package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"

	"golang.org/x/crypto/bcrypt"
)

// Stage 06: Web Masterclass - Authentication
//
//   - Password Hashing: Using bcrypt (Adaptive Hashing)
//   - Context Values: Passing identity through the chain
//   - Auth Middleware: Guarding the gates
//
// ENGINEERING DEPTH:
//   Bcrypt is the gold standard for password hashing because it
//   is "Slow by Design." It includes a "Cost Factor" that you can
//   increase as computers get faster, ensuring that "Brute Force"
//   attacks remain computationally impossible. It also handles
//   salting automatically, so you don't have to manage salts in
//   separate database columns.

type contextKey string

const userCtxKey contextKey = "user"

type User struct {
	ID           int
	Username     string
	PasswordHash string
}

type application struct {
	users map[string]*User
	mu    sync.RWMutex
}

func main() {
	app := &application{
		users: make(map[string]*User),
	}

	mux := http.NewServeMux()

	// Public routes
	mux.HandleFunc("POST /register", app.handleRegister)
	mux.HandleFunc("POST /login", app.handleLogin)

	// Protected route (wrapped in middleware)
	protected := app.authMiddleware(http.HandlerFunc(app.handleProfile))
	mux.Handle("GET /profile", protected)

	fmt.Println("=== Web Masterclass: Authentication ===")
	fmt.Println("  🚀 Server starting on http://localhost:8085")
	fmt.Println()

	log.Fatal(http.ListenAndServe(":8085", mux))

	fmt.Println("\n---------------------------------------------------")
	fmt.Println("NEXT UP: MC.7 forms")
	fmt.Println("Current: MC.6 (authentication)")
	fmt.Println("Previous: MC.5 (sessions)")
	fmt.Println("---------------------------------------------------")
}

func (app *application) handleRegister(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	// 1. Hash the password using bcrypt
	// DefaultCost is currently 10 (2^10 iterations)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	// 2. Store the user
	app.mu.Lock()
	app.users[username] = &User{
		ID:           len(app.users) + 1,
		Username:     username,
		PasswordHash: string(hashedPassword),
	}
	app.mu.Unlock()

	fmt.Fprintf(w, "User %s registered successfully using bcrypt!", username)
}

func (app *application) handleLogin(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	app.mu.RLock()
	user, ok := app.users[username]
	app.mu.RUnlock()

	if !ok {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// 3. Compare the provided password with the stored hash
	err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	fmt.Fprintln(w, "Login successful!")
}

func (app *application) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// SIMULATION: Assume any request with ?user=... is "authenticated" for this demo.
		username := r.URL.Query().Get("user")
		if username == "" {
			http.Error(w, "Forbidden: Missing user query param.", http.StatusForbidden)
			return
		}

		app.mu.RLock()
		user, ok := app.users[username]
		app.mu.RUnlock()

		if !ok {
			http.Error(w, "Forbidden: Unknown user.", http.StatusForbidden)
			return
		}

		// 4. Attach the user identity to the request context
		ctx := context.WithValue(r.Context(), userCtxKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (app *application) handleProfile(w http.ResponseWriter, r *http.Request) {
	// 5. Retrieve the user identity from the context
	user := r.Context().Value(userCtxKey).(*User)
	fmt.Fprintf(w, "Welcome to your profile, %s! (Authenticated via bcrypt)", user.Username)
}
