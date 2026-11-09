package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type contextKey string

const userContextKey = contextKey("user_context")

// Role definitions
type Role string

const (
	RoleOwner  Role = "owner"
	RoleAdmin  Role = "admin"
	RoleMember Role = "member"
	RoleViewer Role = "viewer"
)

// Permission definitions
type Permission string

const (
	PermOrderRead    Permission = "order:read"
	PermOrderCreate  Permission = "order:create"
	PermOrderCancel  Permission = "order:cancel"
	PermTenantManage Permission = "tenant:manage"
	PermAuditView    Permission = "audit:view"
)

// RolePermissions defines the authoritative RBAC access control matrix.
var RolePermissions = map[Role][]Permission{
	RoleOwner: {
		PermOrderRead,
		PermOrderCreate,
		PermOrderCancel,
		PermTenantManage,
		PermAuditView,
	},
	RoleAdmin: {
		PermOrderRead,
		PermOrderCreate,
		PermOrderCancel,
		PermAuditView,
	},
	RoleMember: {
		PermOrderRead,
		PermOrderCreate,
	},
	RoleViewer: {
		PermOrderRead,
	},
}

// UserContext carries the authenticated caller's identity and assigned role.
type UserContext struct {
	UserID   string `json:"user_id"`
	TenantID string `json:"tenant_id"`
	Role     Role   `json:"role"`
}

// HasPermission checks if the user's role grants the requested permission.
func (u *UserContext) HasPermission(perm Permission) bool {
	perms, exists := RolePermissions[u.Role]
	if !exists {
		return false
	}
	for _, p := range perms {
		if p == perm {
			return true
		}
	}
	return false
}

// UserFromContext retrieves the UserContext from the request context.
func UserFromContext(ctx context.Context) (*UserContext, bool) {
	u, ok := ctx.Value(userContextKey).(*UserContext)
	return u, ok && u != nil
}

// ContextWithUser injects a UserContext into the context.
func ContextWithUser(ctx context.Context, u *UserContext) context.Context {
	return context.WithValue(ctx, userContextKey, u)
}

// AuthMiddleware simulates JWT authentication (or inspects an in-memory/mock validator).
type AuthMiddleware struct {
	tokenStore map[string]*UserContext
}

func NewAuthMiddleware() *AuthMiddleware {
	// Seed demo tokens for RBAC testing
	return &AuthMiddleware{
		tokenStore: map[string]*UserContext{
			"token-owner":  {UserID: "usr-owner-01", TenantID: "tenant-acme", Role: RoleOwner},
			"token-admin":  {UserID: "usr-admin-02", TenantID: "tenant-acme", Role: RoleAdmin},
			"token-member": {UserID: "usr-member-03", TenantID: "tenant-acme", Role: RoleMember},
			"token-viewer": {UserID: "usr-viewer-04", TenantID: "tenant-acme", Role: RoleViewer},
		},
	}
}

// Authenticate extracts the Bearer token and places the verified user into context.
func (am *AuthMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, `{"error":"unauthorized","message":"missing Authorization header"}`, http.StatusUnauthorized)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, `{"error":"unauthorized","message":"invalid Authorization format, expected 'Bearer <token>'"}`, http.StatusUnauthorized)
			return
		}

		user, exists := am.tokenStore[parts[1]]
		if !exists {
			http.Error(w, `{"error":"unauthorized","message":"invalid or expired token"}`, http.StatusUnauthorized)
			return
		}

		ctx := ContextWithUser(r.Context(), user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RequirePermission enforces that the authenticated user possesses the specified RBAC permission.
func RequirePermission(perm Permission, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, ok := UserFromContext(r.Context())
		if !ok {
			http.Error(w, `{"error":"unauthorized","message":"authentication required"}`, http.StatusUnauthorized)
			return
		}

		if !user.HasPermission(perm) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusForbidden)
			_ = json.NewEncoder(w).Encode(map[string]any{
				"error":     "forbidden",
				"message":   fmt.Sprintf("insufficient permissions: required permission '%s'", perm),
				"user_role": user.Role,
			})
			return
		}

		next(w, r)
	}
}

// BuildRouter registers protected routes with permission guards.
func BuildRouter(am *AuthMiddleware) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/v1/orders", RequirePermission(PermOrderRead, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"status":"success","orders":[]}`))
	}))

	mux.HandleFunc("POST /api/v1/orders", RequirePermission(PermOrderCreate, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{"status":"created","order_id":"ord-new-1"}`))
	}))

	mux.HandleFunc("POST /api/v1/orders/cancel", RequirePermission(PermOrderCancel, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"status":"cancelled"}`))
	}))

	mux.HandleFunc("POST /api/v1/tenant/settings", RequirePermission(PermTenantManage, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"status":"updated","tenant":"settings_saved"}`))
	}))

	return am.Authenticate(mux)
}

func main() {
	am := NewAuthMiddleware()
	router := BuildRouter(am)

	fmt.Println("Opslane RBAC authorization router listening on :8080...")
	_ = http.ListenAndServe(":8080", router)
}
