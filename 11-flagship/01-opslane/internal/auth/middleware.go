// Copyright (c) 2026 Rasel Hossen
// See LICENSE for usage terms.

package auth

import (
	"net/http"
	"strings"
)

// RequireAuth rejects anonymous requests and attaches verified tenant identity to context.
func RequireAuth(tokens *TokenManager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			identity, err := IdentityFromRequest(tokens, r)
			if err != nil {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r.WithContext(WithIdentity(r.Context(), identity)))
		})
	}
}

// IdentityFromRequest verifies the Authorization header and returns trusted token claims.
func IdentityFromRequest(tokens *TokenManager, r *http.Request) (Identity, error) {
	if tokens == nil {
		return Identity{}, ErrInvalidToken
	}

	header := strings.TrimSpace(r.Header.Get("Authorization"))
	if header == "" {
		return Identity{}, ErrInvalidToken
	}

	parts := strings.Fields(header)
	if len(parts) != 2 {
		return Identity{}, ErrInvalidToken
	}

	if !strings.EqualFold(parts[0], "Bearer") {
		return Identity{}, ErrInvalidToken
	}

	token := strings.TrimSpace(parts[1])
	if token == "" {
		return Identity{}, ErrInvalidToken
	}

	return tokens.Verify(token)
}
