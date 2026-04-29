// Copyright (c) 2026 Rasel Hossen
// See LICENSE for usage terms.

package middleware

import (
	"fmt"
	"net/http"
	"time"
)

// CacheControl sets Cache-Control headers on responses. Use this for
// public endpoints like health or static assets. Protected API
// endpoints should use "no-store" to prevent proxies and browsers from
// caching authenticated data.
//
// maxAge controls how long clients and intermediary proxies are allowed
// to serve the cached response without revalidating with the origin.
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

// NoCache sets Cache-Control: no-store on responses. This is the right
// default for authenticated API endpoints where stale data could leak
// across tenants or users.
func NoCache(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store")
		next.ServeHTTP(w, r)
	})
}
