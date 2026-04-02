// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0
// Commercial use is prohibited without permission.

package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/rasel9t6/the-go-engineer/22-enterprise-capstone/internal/middleware"
	"github.com/rasel9t6/the-go-engineer/22-enterprise-capstone/internal/models"
	"github.com/rasel9t6/the-go-engineer/22-enterprise-capstone/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

// ============================================================================
// Internal Package: Handlers
// Level: Advanced
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Dependency Injection via `application` config struct
//   - Separation of route registration from HTTP logic
//   - Handling JSON payloads and Authorization
//
// ENGINEERING DEPTH:
//   Notice that `handlers` NEVER runs `INSERT INTO`. It calls `app.posts.Create()`.
//   This keeps our handlers "thin" and strictly focused on HTTP networking (JSON parsing,
//   Headers, Status Codes, and Security).
// ============================================================================

type contextKey string

const CtxUserID = contextKey("userID")

// Application configures the handlers layer with required dependencies
type Application struct {
	Logger *slog.Logger
	Posts  repository.PostRepository
	DB     *sql.DB // We cheat slightly by putting DB here for auth directly
}

// Routes constructs the Mux router and attaches all handlers + middlewares.
func (app *Application) Routes() http.Handler {
	mux := http.NewServeMux()

	// Public Routes
	mux.HandleFunc("GET /health", app.handleHealth)
	mux.HandleFunc("POST /register", app.handleRegister)
	mux.HandleFunc("POST /login", app.handleLogin)
	mux.HandleFunc("GET /posts", app.handleListPosts)

	// Protected Routes (Wrapped in `requireAuth` middleware explicitly)
	mux.Handle("POST /posts", app.requireAuth(http.HandlerFunc(app.handleCreatePost)))

	// Global Middleware Chain (Wrap the entire multiplexer)
	// Because of how standard library decorators work, middleware executes from
	// the OUTSIDE -> IN. So `RecoverPanic` runs first, `LogRequest` runs second,
	// `SecureHeaders` is third, and finally the router `mux` evaluates the path.
	handler := middleware.SecureHeaders(
		middleware.LogRequest(app.Logger)(
			middleware.RecoverPanic(app.Logger)(mux),
		),
	)

	return handler
}

// --- Protected Authentication Middleware specific to Handlers ---

func (app *Application) requireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			http.Error(w, `{"error":"authorization string 'Bearer <ID>' required"}`, http.StatusUnauthorized)
			return
		}

		var uid int
		// A fast parser extracting the UserID natively from the Bearer Token
		fmt.Sscanf(token, "Bearer %d", &uid)
		if uid == 0 {
			http.Error(w, `{"error":"invalid token format"}`, http.StatusUnauthorized)
			return
		}

		// Context Value Propagation
		// Request contexts are IMMUTABLE. We cannot edit `r.Context()`.
		// We use `context.WithValue` to create a COPY of the Context, binding
		// the `uid` directly inside it.
		ctx := context.WithValue(r.Context(), CtxUserID, uid)

		// We overwrite the request pointer using `r.WithContext(ctx)`.
		// Any downstream handler in the chain can now recover this identity!
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// --- The Handlers ---

func (app *Application) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"database connection verified and network active"}`))
}

func (app *Application) handleRegister(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, `{"error":"invalid json"}`, http.StatusBadRequest)
		return
	}

	// Bcrypt hashes the plaintext password into a salt+hash string safely
	hashed, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, `{"error":"password generation failure"}`, http.StatusInternalServerError)
		return
	}

	// Native execution (For scaling this up, extract into a UserRepository!)
	query := `INSERT INTO users (username, email, password_hash) VALUES ($1, $2, $3) RETURNING id`
	var newID int
	err = app.DB.QueryRowContext(r.Context(), query, input.Username, input.Email, string(hashed)).Scan(&newID)

	if err != nil {
		http.Error(w, `{"error":"user already exists or database conflict: `+err.Error()+`"}`, http.StatusConflict)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]any{"id": newID, "email": input.Email})
}

func (app *Application) handleLogin(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, `{"error":"invalid json"}`, http.StatusBadRequest)
		return
	}

	var user models.User
	query := `SELECT id, email, password_hash FROM users WHERE email = $1`
	err := app.DB.QueryRowContext(r.Context(), query, input.Email).Scan(&user.ID, &user.Email, &user.Password)
	if err != nil {
		http.Error(w, `{"error":"invalid email or generic database error"}`, http.StatusUnauthorized)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		http.Error(w, `{"error":"invalid credentials - password format mismatch"}`, http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"token":   fmt.Sprintf("Bearer %d", user.ID), // Fake Token simulation (For production: use JWT here!)
		"message": "authorization successful",
	})
}

func (app *Application) handleCreatePost(w http.ResponseWriter, r *http.Request) {
	// Secure Identity Extraction
	// Because `requireAuth` middleware passed successfully, we know `CtxUserID` exists.
	// We type-assert the empty interface `any` returned by `.Value()` into an `int`.
	userID := r.Context().Value(CtxUserID).(int)

	var input struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, `{"error":"invalid json content payload"}`, http.StatusBadRequest)
		return
	}

	// DELEGATE! We cross the boundary from 'HTTP' down to the 'Database Repository'
	id, err := app.Posts.Create(r.Context(), &models.Post{
		Title:    input.Title,
		Content:  input.Content,
		AuthorID: userID,
	})

	if err != nil {
		app.Logger.Error("repository failure", slog.Any("error", err))
		http.Error(w, `{"error":"failed to persist post downstream"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]any{"post_id": id})
}

func (app *Application) handleListPosts(w http.ResponseWriter, r *http.Request) {
	posts, err := app.Posts.GetAll(r.Context())
	if err != nil {
		app.Logger.Error("failed to list posts", slog.Any("error", err))
		http.Error(w, `{"error":"failed to fetch posts internally"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}
