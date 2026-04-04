// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0
// Commercial use is prohibited without permission.

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"golang.org/x/crypto/bcrypt"
)

// ============================================================================
// Section 13: Web Masterclass — Authentication
// Level: Advanced
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Password hashing with bcrypt (NEVER store plaintext passwords!)
//   - Login and registration flows
//   - Authentication middleware to protect routes
//   - Storing user identity in request context
//
// ENGINEERING DEPTH:
//   A common misconception is that you must store a "Salt" in a separate
//   database column to secure passwords against Rainbow Tables. The `bcrypt`
//   algorithm natively generates a mathematically perfect 128-bit salt and
//   literally embeds it directly inside the final output string
//   (`$2y$10$SALT...HASH...`). When you call `bcrypt.CompareHashAndPassword`,
//   Go automatically extracts the salt from the string prefix, hashes the
//   incoming attempt with that *exact* salt, and performs a timing-attack resistant
//   byte-by-byte comparison. One string, absolute security.
//
// SECURITY RULES:
//   1. ALWAYS hash passwords with bcrypt (or argon2)
//   2. NEVER store plaintext passwords
//   3. Use HttpOnly + Secure cookies for sessions
//   4. Check authentication in middleware, not in each handler
//
// RUN: go run ./13-web-masterclass/6-auth
// ============================================================================

// contextKey is a custom type to avoid context key collisions.
// Using a plain string could collide with other packages.
type contextKey string

const contextKeyUserID = contextKey("userID")

// User represents a user in the system.
type User struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"-"` // json:"-" ensures password NEVER appears in JSON output
}

// In-memory user store (production: use a database)
var (
	users      = map[string]*User{}
	usersMutex = sync.RWMutex{}
)

func main() {
	mux := http.NewServeMux()

	// Public routes
	mux.HandleFunc("POST /register", handleRegister)
	mux.HandleFunc("POST /login", handleLogin)

	// Protected routes — wrapped with requireAuth middleware
	mux.Handle("GET /profile", requireAuth(http.HandlerFunc(handleProfile)))

	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}

func handleRegister(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")

	if email == "" || password == "" {
		http.Error(w, "Email and password required", http.StatusBadRequest)
		return
	}

	usersMutex.RLock()
	if _, exists := users[email]; exists {
		usersMutex.RUnlock()
		http.Error(w, "User already exists", http.StatusConflict)
		return
	}
	usersMutex.RUnlock()

	// bcrypt.GenerateFromPassword:
	//   - Adds a random salt automatically (no need to manage salts)
	//   - bcrypt.DefaultCost = 10 (2^10 = 1024 iterations)
	//   - Higher cost = more secure but slower
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	usersMutex.Lock()
	users[email] = &User{
		ID:       len(users) + 1,
		Email:    email,
		Password: string(hashedPassword),
	}
	usersMutex.Unlock()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "User registered successfully",
		"email":   email,
	})
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")

	usersMutex.RLock()
	user, exists := users[email]
	usersMutex.RUnlock()

	if !exists {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// bcrypt.CompareHashAndPassword:
	//   - Extracts the salt from the stored hash
	//   - Hashes the provided password with the same salt + cost
	//   - Compares the results in constant time (prevents timing attacks)
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// In production: create a session and set a session cookie.
	// For this example, we use a simple header-based approach.
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"message": "Login successful",
		"user_id": user.ID,
		"token":   fmt.Sprintf("user-%d", user.ID), // Simplified; use JWT in production
	})
}

// requireAuth is an authentication middleware.
// It checks for a valid auth token and adds the user ID to the request context.
func requireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// In production, validate a JWT or session cookie here.
		token := r.Header.Get("Authorization")
		if token == "" {
			http.Error(w, "Authorization required", http.StatusUnauthorized)
			return
		}

		// Extract user ID from token (simplified)
		var userID int
		fmt.Sscanf(token, "user-%d", &userID)
		if userID == 0 {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Add user ID to request context.
		// context.WithValue creates a NEW context with the added value.
		// This is how information flows through the middleware chain.
		ctx := context.WithValue(r.Context(), contextKeyUserID, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func handleProfile(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from context (set by requireAuth middleware)
	userID := r.Context().Value(contextKeyUserID).(int)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"message": "Welcome to your profile",
		"user_id": userID,
	})
}
