package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

type contextKey string

const userKey contextKey = "user"

type User struct {
	ID    string
	Role  string
	Email string
}

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if auth == "" || !strings.HasPrefix(auth, "Bearer ") {
			http.Error(w, `{"error":"missing token"}`, http.StatusUnauthorized)
			return
		}
		token := strings.TrimPrefix(auth, "Bearer ")
		user, ok := validateToken(token)
		if !ok {
			http.Error(w, `{"error":"invalid token"}`, http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), userKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func requireRole(role string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, ok := r.Context().Value(userKey).(User)
			if !ok {
				http.Error(w, `{"error":"unauthenticated"}`, http.StatusUnauthorized)
				return
			}
			if user.Role != role && user.Role != "admin" {
				http.Error(w, `{"error":"insufficient permissions"}`, http.StatusForbidden)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func validateToken(token string) (User, bool) {
	tokens := map[string]User{
		"token-user":  {ID: "u1", Role: "user", Email: "alice@example.com"},
		"token-admin": {ID: "u2", Role: "admin", Email: "bob@example.com"},
	}
	u, ok := tokens[token]
	return u, ok
}

func BuildMiddlewareChain(boundaries []string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			for _, b := range boundaries {
				log.Printf("[trust-boundary] crossing: %s", b)
			}
			next.ServeHTTP(w, r)
		})
	}
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	protected := http.NewServeMux()
	protected.HandleFunc("/profile", func(w http.ResponseWriter, r *http.Request) {
		user, ok := r.Context().Value(userKey).(User)
		if !ok {
			http.Error(w, `{"error":"unauthenticated"}`, http.StatusUnauthorized)
			return
		}
		fmt.Fprintf(w, `{"user_id":"%s","role":"%s"}`, user.ID, user.Role)
	})

	mux.Handle("/api/", authMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = strings.TrimPrefix(r.URL.Path, "/api")
		protected.ServeHTTP(w, r)
	})))

	mux.Handle("/api/admin", authMiddleware(requireRole("admin")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"access":"admin_panel"}`)
	}))))

	server := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	fmt.Println("Trust boundaries server on :8080")
	server.ListenAndServe()
}
