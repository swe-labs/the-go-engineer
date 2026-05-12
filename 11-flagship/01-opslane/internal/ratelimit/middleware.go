// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

package ratelimit

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/netip"
	"sync"
	"time"

	"log/slog"

	"github.com/swe-labs/the-go-engineer/11-flagship/01-opslane/internal/middleware"
)

// LimiterKey (Function): extracts a rate limit key from an HTTP request
type LimiterKey func(r *http.Request) string

// ByIP (Function): creates a LimiterKey that extracts the client IP address
func ByIP(trustedProxyCIDRs []netip.Prefix) LimiterKey {
	return func(r *http.Request) string {
		return middleware.ClientAddress(r, trustedProxyCIDRs)
	}
}

// ByHeader (Function): creates a LimiterKey that extracts a specific HTTP header value
func ByHeader(header string) func(r *http.Request) string {
	return func(r *http.Request) string {
		return r.Header.Get(header)
	}
}

// ByTenantAndUser (Function): creates a LimiterKey combining tenant and user identifiers
func ByTenantAndUser(tenantID, userID func(*http.Request) string) func(r *http.Request) string {
	return func(r *http.Request) string {
		return tenantID(r) + ":" + userID(r)
	}
}

// Middleware (Function): HTTP rate limiting middleware with fail-open and rate limit headers
func Middleware(limiter *Limiter, keyFunc LimiterKey, logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			key := keyFunc(r)
			decision, err := limiter.AllowWithDecision(r.Context(), key)
			if err != nil {
				if logger != nil {
					logger.Warn("rate limit check failed, allowing request",
						slog.String("key", key),
						slog.Any("error", err))
				}
				decision = Decision{Allowed: true, Limit: limiter.cfg.BurstSize, Remaining: limiter.cfg.BurstSize}
			}

			setRateLimitHeaders(w, decision)
			if !decision.Allowed {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusTooManyRequests)
				_ = json.NewEncoder(w).Encode(map[string]any{
					"error": map[string]string{
						"code":    "rate_limited",
						"message": "rate limit exceeded",
					},
				})
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// setRateLimitHeaders (Function): sets X-RateLimit-* and Retry-After response headers
func setRateLimitHeaders(w http.ResponseWriter, d Decision) {
	limit := d.Limit
	if limit <= 0 {
		limit = 1
	}
	remaining := d.Remaining
	if remaining < 0 {
		remaining = 0
	}
	w.Header().Set("X-RateLimit-Limit", fmt.Sprintf("%d", limit))
	w.Header().Set("X-RateLimit-Remaining", fmt.Sprintf("%d", remaining))
	if !d.ResetAt.IsZero() {
		reset := int(time.Until(d.ResetAt).Seconds())
		if reset < 0 {
			reset = 0
		}
		w.Header().Set("X-RateLimit-Reset", fmt.Sprintf("%d", reset))
		if !d.Allowed {
			w.Header().Set("Retry-After", fmt.Sprintf("%d", reset))
		}
	}
}

// PerIPMiddleware (Function): convenience middleware for per-IP rate limiting
func PerIPMiddleware(limiter *Limiter, trustedProxyCIDRs []netip.Prefix, logger *slog.Logger) func(http.Handler) http.Handler {
	return Middleware(limiter, ByIP(trustedProxyCIDRs), logger)
}

// AuthenticatedMiddleware (Function): convenience middleware for per-tenant-user rate limiting
func AuthenticatedMiddleware(limiter *Limiter, tenantIDFunc, userIDFunc func(*http.Request) string, logger *slog.Logger) func(http.Handler) http.Handler {
	return Middleware(limiter, ByTenantAndUser(tenantIDFunc, userIDFunc), logger)
}

// LimiterPool (Struct): thread-safe registry of named limiters
type LimiterPool struct {
	mu       sync.RWMutex
	limiters map[string]*Limiter
}

// pool (Struct): package-level default limiter pool
var pool = &LimiterPool{
	limiters: make(map[string]*Limiter),
}

// RegisterLimiter (Function): registers a named limiter in the default pool
func RegisterLimiter(name string, limiter *Limiter) {
	pool.mu.Lock()
	defer pool.mu.Unlock()
	pool.limiters[name] = limiter
}

// GetLimiter (Function): retrieves a named limiter from the default pool
func GetLimiter(name string) *Limiter {
	pool.mu.RLock()
	defer pool.mu.RUnlock()
	return pool.limiters[name]
}
