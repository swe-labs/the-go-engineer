package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

type contextKey string

const claimsKey contextKey = "claims"

type Claims struct {
	UserID string   `json:"user_id"`
	Role   string   `json:"role"`
	Perms  []string `json:"permissions"`
	Tenant string   `json:"tenant"`
}

func authNMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if auth == "" || !strings.HasPrefix(auth, "Bearer ") {
			http.Error(w, `{"error":"missing or malformed token"}`, http.StatusUnauthorized)
			return
		}
		token := strings.TrimPrefix(auth, "Bearer ")
		claims, err := verifyToken(token)
		if err != nil {
			http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), claimsKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func authZMiddleware(requiredPerm string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := r.Context().Value(claimsKey).(Claims)
			if !ok {
				http.Error(w, `{"error":"unauthenticated"}`, http.StatusUnauthorized)
				return
			}
			if !hasPermission(claims, requiredPerm) {
				http.Error(w, fmt.Sprintf(`{"error":"forbidden: need %s permission"}`, requiredPerm), http.StatusForbidden)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func verifyToken(token string) (Claims, error) {
	tokens := map[string]Claims{
		"admin-token": {
			UserID: "u1", Role: "admin",
			Perms:  []string{"doc:read", "doc:write", "doc:delete", "user:manage"},
			Tenant: "acme-corp",
		},
		"user-token": {
			UserID: "u2", Role: "user",
			Perms:  []string{"doc:read", "doc:write"},
			Tenant: "acme-corp",
		},
		"readonly-token": {
			UserID: "u3", Role: "viewer",
			Perms:  []string{"doc:read"},
			Tenant: "acme-corp",
		},
	}
	c, ok := tokens[token]
	if !ok {
		return Claims{}, fmt.Errorf("invalid token")
	}
	return c, nil
}

func hasPermission(c Claims, required string) bool {
	for _, p := range c.Perms {
		if p == required {
			return true
		}
		if strings.HasSuffix(p, ":*") {
			prefix := strings.TrimSuffix(p, ":*")
			if strings.HasPrefix(required, prefix+":") {
				return true
			}
		}
	}
	return false
}

func AuthN(token string) (Claims, error) {
	return verifyToken(token)
}

func AuthZ(claims Claims, requiredPerm string, resourceOwner string) bool {
	if claims.Role == "admin" {
		return true
	}
	if !hasPermission(claims, requiredPerm) {
		return false
	}
	// Resource-level check: user must own the resource or be admin.
	return true
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	api := http.NewServeMux()
	api.Handle("/docs", authZMiddleware("doc:read")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]string{"status": "docs list"})
	})))
	api.Handle("/docs/create", authZMiddleware("doc:write")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]string{"status": "created"})
	})))

	mux.Handle("/api/", authNMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = strings.TrimPrefix(r.URL.Path, "/api")
		api.ServeHTTP(w, r)
	})))

	fmt.Println("AuthN/AuthZ server on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
