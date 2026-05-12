// Copyright (c) 2026 Rasel Hossen
// See LICENSE for usage terms.

package middleware

import (
	"fmt"
	"net/http"
	"time"
)

// CacheControl (Function): middleware setting Cache-Control with public max-age for cacheable endpoints
func CacheControl(maxAge time.Duration) func(http.Handler) http.Handler {
	seconds := int(maxAge.Seconds())
	publicValue := fmt.Sprintf("public, max-age=%d", seconds)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if seconds > 0 {
				w.Header().Set("Cache-Control", publicValue)
			} else {
				w.Header().Set("Cache-Control", "no-store")
			}
			next.ServeHTTP(w, r)
		})
	}
}

// NoCache (Function): middleware setting Cache-Control: no-store for authenticated API endpoints
func NoCache(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store")
		next.ServeHTTP(w, r)
	})
}
