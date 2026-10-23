package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrNotFound     = errors.New("not found")
	ErrUnauthorized = errors.New("unauthorized")
	ErrForbidden    = errors.New("forbidden")
	ErrInvalidInput = errors.New("invalid input")
	ErrUserExists   = errors.New("user already exists")
	ErrCrossTenant  = errors.New("cross-tenant access prohibited")
)

type Role string

const (
	RoleAdmin  Role = "ADMIN"
	RoleMember Role = "MEMBER"
	RoleViewer Role = "VIEWER"
)

func RoleLevel(r Role) int {
	switch r {
	case RoleAdmin:
		return 3
	case RoleMember:
		return 2
	case RoleViewer:
		return 1
	default:
		return 0
	}
}

type User struct {
	ID           string    `json:"id"`
	TenantID     string    `json:"tenant_id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"` // Never serialize to JSON
	Role         Role      `json:"role"`
	CreatedAt    time.Time `json:"created_at"`
}

type Document struct {
	ID        string    `json:"id"`
	TenantID  string    `json:"tenant_id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	OwnerID   string    `json:"owner_id"`
	CreatedAt time.Time `json:"created_at"`
}

type Claims struct {
	UserID   string `json:"uid"`
	TenantID string `json:"tid"`
	Role     Role   `json:"role"`
	jwt.RegisteredClaims
}

type RegisterRequest struct {
	TenantID string `json:"tenant_id"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     Role   `json:"role"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type CreateDocumentRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type JSONResponse struct {
	Data  any       `json:"data,omitempty"`
	Error *APIError `json:"error,omitempty"`
}

type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func GenerateID() string {
	b := make([]byte, 8)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

func HashPassword(password string) (string, error) {
	// TODO: Hash password using bcrypt.GenerateFromPassword with DefaultCost
	return "", errors.New("not implemented")
}

func CheckPassword(hash, password string) bool {
	// TODO: Verify password using bcrypt.CompareHashAndPassword
	return false
}

func GenerateJWT(user User, secret []byte, ttl time.Duration) (string, error) {
	// TODO: Create Claims with UserID, TenantID, Role, and RegisteredClaims (ExpiresAt, IssuedAt)
	// TODO: Sign using HMAC-SHA256 (jwt.SigningMethodHS256) and secret key
	return "", errors.New("not implemented")
}

func ParseJWT(tokenStr string, secret []byte) (*Claims, error) {
	// TODO: Parse token with claims, verify HMAC algorithm, return parsed claims
	return nil, errors.New("not implemented")
}

type MemoryStore struct {
	mu        sync.RWMutex
	users     map[string]User
	usersByID map[string]User
	docs      map[string]Document
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		users:     make(map[string]User),
		usersByID: make(map[string]User),
		docs:      make(map[string]Document),
	}
}

func (s *MemoryStore) Register(req RegisterRequest) (User, error) {
	if strings.TrimSpace(req.TenantID) == "" {
		return User{}, fmt.Errorf("%w: tenant_id is required", ErrInvalidInput)
	}
	if strings.TrimSpace(req.Email) == "" || !strings.Contains(req.Email, "@") {
		return User{}, fmt.Errorf("%w: valid email is required", ErrInvalidInput)
	}
	if len(req.Password) < 8 {
		return User{}, fmt.Errorf("%w: password must be at least 8 characters", ErrInvalidInput)
	}
	if RoleLevel(req.Role) == 0 {
		return User{}, fmt.Errorf("%w: invalid role (must be ADMIN, MEMBER, or VIEWER)", ErrInvalidInput)
	}

	hash, err := HashPassword(req.Password)
	if err != nil {
		return User{}, err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	emailKey := strings.ToLower(strings.TrimSpace(req.Email))
	if _, exists := s.users[emailKey]; exists {
		return User{}, ErrUserExists
	}

	user := User{
		ID:           GenerateID(),
		TenantID:     strings.TrimSpace(req.TenantID),
		Email:        emailKey,
		PasswordHash: hash,
		Role:         req.Role,
		CreatedAt:    time.Now().UTC(),
	}

	s.users[emailKey] = user
	s.usersByID[user.ID] = user
	return user, nil
}

func (s *MemoryStore) Authenticate(email, password string) (User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	user, exists := s.users[strings.ToLower(strings.TrimSpace(email))]
	if !exists {
		return User{}, ErrUnauthorized
	}

	if !CheckPassword(user.PasswordHash, password) {
		return User{}, ErrUnauthorized
	}

	return user, nil
}

func (s *MemoryStore) CreateDocument(tenantID, ownerID, title, content string) (Document, error) {
	if strings.TrimSpace(title) == "" {
		return Document{}, fmt.Errorf("%w: title is required", ErrInvalidInput)
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	doc := Document{
		ID:        GenerateID(),
		TenantID:  tenantID,
		Title:     strings.TrimSpace(title),
		Content:   content,
		OwnerID:   ownerID,
		CreatedAt: time.Now().UTC(),
	}
	s.docs[doc.ID] = doc
	return doc, nil
}

func (s *MemoryStore) GetDocument(docID string) (Document, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	doc, exists := s.docs[docID]
	if !exists {
		return Document{}, ErrNotFound
	}
	return doc, nil
}

func (s *MemoryStore) DeleteDocument(docID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.docs[docID]; !exists {
		return ErrNotFound
	}
	delete(s.docs, docID)
	return nil
}

type authContextKey string

const userClaimsKey authContextKey = "user_claims"

func ClaimsFromContext(ctx context.Context) (*Claims, bool) {
	claims, ok := ctx.Value(userClaimsKey).(*Claims)
	return claims, ok
}

// Security Headers Middleware
func SecurityHeadersMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("Content-Security-Policy", "default-src 'self'")
		next.ServeHTTP(w, r)
	})
}

// Auth Middleware: validates JWT and injects claims into context
func AuthMiddleware(secret []byte) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// TODO: Extract Bearer token from Authorization header
			// TODO: Parse and validate JWT, return 401 on error
			// TODO: Inject claims into request context with userClaimsKey
			next.ServeHTTP(w, r)
		})
	}
}

// RBAC Middleware: ensures minimum role requirement
func RequireRole(minRole Role) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// TODO: Retrieve claims from context, return 401 if missing
			// TODO: Compare RoleLevel(claims.Role) < RoleLevel(minRole), return 403 Forbidden if insufficient
			next.ServeHTTP(w, r)
		})
	}
}

type APIHandler struct {
	store  *MemoryStore
	secret []byte
	router *http.ServeMux
}

func NewAPIHandler(store *MemoryStore, secret []byte) http.Handler {
	api := &APIHandler{
		store:  store,
		secret: secret,
		router: http.NewServeMux(),
	}

	authMW := AuthMiddleware(secret)
	requireMember := RequireRole(RoleMember)
	requireAdmin := RequireRole(RoleAdmin)

	api.router.HandleFunc("POST /auth/register", api.handleRegister)
	api.router.HandleFunc("POST /auth/login", api.handleLogin)

	api.router.Handle("GET /auth/me", authMW(http.HandlerFunc(api.handleMe)))
	api.router.Handle("POST /documents", authMW(requireMember(http.HandlerFunc(api.handleCreateDoc))))
	api.router.Handle("GET /documents/{id}", authMW(http.HandlerFunc(api.handleGetDoc)))
	api.router.Handle("DELETE /documents/{id}", authMW(requireAdmin(http.HandlerFunc(api.handleDeleteDoc))))

	return SecurityHeadersMiddleware(api.router)
}

func (a *APIHandler) handleRegister(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, JSONResponse{
			Error: &APIError{Code: "BAD_REQUEST", Message: "Malformed JSON payload"},
		})
		return
	}

	user, err := a.store.Register(req)
	if err != nil {
		if errors.Is(err, ErrUserExists) {
			writeJSON(w, http.StatusConflict, JSONResponse{
				Error: &APIError{Code: "CONFLICT", Message: "User email already registered"},
			})
			return
		}
		if errors.Is(err, ErrInvalidInput) {
			writeJSON(w, http.StatusUnprocessableEntity, JSONResponse{
				Error: &APIError{Code: "VALIDATION_FAILED", Message: err.Error()},
			})
			return
		}
		writeJSON(w, http.StatusInternalServerError, JSONResponse{
			Error: &APIError{Code: "INTERNAL_ERROR", Message: "Failed to register user"},
		})
		return
	}

	writeJSON(w, http.StatusCreated, JSONResponse{Data: user})
}

func (a *APIHandler) handleLogin(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, JSONResponse{
			Error: &APIError{Code: "BAD_REQUEST", Message: "Malformed JSON payload"},
		})
		return
	}

	user, err := a.store.Authenticate(req.Email, req.Password)
	if err != nil {
		writeJSON(w, http.StatusUnauthorized, JSONResponse{
			Error: &APIError{Code: "UNAUTHORIZED", Message: "Invalid credentials"},
		})
		return
	}

	token, err := GenerateJWT(user, a.secret, 2*time.Hour)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, JSONResponse{
			Error: &APIError{Code: "INTERNAL_ERROR", Message: "Failed to issue token"},
		})
		return
	}

	writeJSON(w, http.StatusOK, JSONResponse{
		Data: map[string]any{
			"token":      token,
			"token_type": "Bearer",
			"expires_in": 7200,
			"user":       user,
		},
	})
}

func (a *APIHandler) handleMe(w http.ResponseWriter, r *http.Request) {
	claims, _ := ClaimsFromContext(r.Context())
	writeJSON(w, http.StatusOK, JSONResponse{Data: claims})
}

func (a *APIHandler) handleCreateDoc(w http.ResponseWriter, r *http.Request) {
	claims, _ := ClaimsFromContext(r.Context())
	var req CreateDocumentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, JSONResponse{
			Error: &APIError{Code: "BAD_REQUEST", Message: "Malformed JSON payload"},
		})
		return
	}

	doc, err := a.store.CreateDocument(claims.TenantID, claims.UserID, req.Title, req.Content)
	if err != nil {
		writeJSON(w, http.StatusUnprocessableEntity, JSONResponse{
			Error: &APIError{Code: "VALIDATION_FAILED", Message: err.Error()},
		})
		return
	}

	writeJSON(w, http.StatusCreated, JSONResponse{Data: doc})
}

func (a *APIHandler) handleGetDoc(w http.ResponseWriter, r *http.Request) {
	claims, _ := ClaimsFromContext(r.Context())
	id := r.PathValue("id")

	doc, err := a.store.GetDocument(id)
	if err != nil {
		writeJSON(w, http.StatusNotFound, JSONResponse{
			Error: &APIError{Code: "NOT_FOUND", Message: fmt.Sprintf("Document %s not found", id)},
		})
		return
	}

	// Tenant Isolation: strictly enforce tenant boundary
	if doc.TenantID != claims.TenantID {
		writeJSON(w, http.StatusForbidden, JSONResponse{
			Error: &APIError{Code: "FORBIDDEN", Message: "Cross-tenant access prohibited"},
		})
		return
	}

	writeJSON(w, http.StatusOK, JSONResponse{Data: doc})
}

func (a *APIHandler) handleDeleteDoc(w http.ResponseWriter, r *http.Request) {
	claims, _ := ClaimsFromContext(r.Context())
	id := r.PathValue("id")

	doc, err := a.store.GetDocument(id)
	if err != nil {
		writeJSON(w, http.StatusNotFound, JSONResponse{
			Error: &APIError{Code: "NOT_FOUND", Message: fmt.Sprintf("Document %s not found", id)},
		})
		return
	}

	// Tenant Isolation: cannot delete other tenant's document
	if doc.TenantID != claims.TenantID {
		writeJSON(w, http.StatusForbidden, JSONResponse{
			Error: &APIError{Code: "FORBIDDEN", Message: "Cross-tenant access prohibited"},
		})
		return
	}

	_ = a.store.DeleteDocument(id)
	w.WriteHeader(http.StatusNoContent)
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func main() {
	secret := []byte(os.Getenv("JWT_SECRET"))
	if len(secret) == 0 {
		secret = []byte("super-secure-production-jwt-signing-secret-key-256")
	}

	store := NewMemoryStore()
	handler := NewAPIHandler(store, secret)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8082"
	}

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      handler,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	go func() {
		slog.Info("secure auth api running", "port", port)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("server error", "error", err)
			os.Exit(1)
		}
	}()

	<-shutdown
	slog.Info("shutting down secure api gracefully")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		slog.Error("shutdown failed", "error", err)
	}
}
