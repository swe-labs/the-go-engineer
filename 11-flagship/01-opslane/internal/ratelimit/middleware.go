// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

package ratelimit

import (
	"net/http"
	"strings"
	"sync"
	"time"

	"log/slog"
)

type LimiterKey func(r *http.Request) string

func ByIP(r *http.Request) string {
	ip := r.Header.Get("X-Forwarded-For")
	if ip == "" {
		ip = r.RemoteAddr
		if idx := strings.LastIndex(ip, ":"); idx != -1 {
			ip = ip[:idx]
		}
	}
	return ip
}

func ByHeader(header string) func(r *http.Request) string {
	return func(r *http.Request) string {
		return r.Header.Get(header)
	}
}

func ByTenantAndUser(tenantID, userID func(*http.Request) string) func(r *http.Request) string {
	return func(r *http.Request) string {
		return tenantID(r) + ":" + userID(r)
	}
}

func Middleware(limiter *Limiter, keyFunc LimiterKey, logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			key := keyFunc(r)

			allowed, err := limiter.Allow(key)
			if err != nil {
				if logger != nil {
					logger.Warn("rate limit check failed, allowing request",
						slog.String("key", key),
						slog.Any("error", err))
				}
			}

			if !allowed {
				w.Header().Set("Retry-After", "1")
				w.Header().Set("X-RateLimit-Limit", "10")
				w.Header().Set("X-RateLimit-Remaining", "0")
				http.Error(w, "rate limit exceeded", http.StatusTooManyRequests)
				return
			}

			w.Header().Set("X-RateLimit-Limit", "10")
			w.Header().Set("X-RateLimit-Remaining", "0")

			next.ServeHTTP(w, r)
		})
	}
}

func PerIPMiddleware(limiter *Limiter, logger *slog.Logger) func(http.Handler) http.Handler {
	return Middleware(limiter, ByIP, logger)
}

func AuthenticatedMiddleware(limiter *Limiter, tenantIDFunc, userIDFunc func(*http.Request) string, logger *slog.Logger) func(http.Handler) http.Handler {
	return Middleware(limiter, ByTenantAndUser(tenantIDFunc, userIDFunc), logger)
}

type LimiterPool struct {
	mu       sync.RWMutex
	limiters map[string]*Limiter
}

var pool = &LimiterPool{
	limiters: make(map[string]*Limiter),
}

func RegisterLimiter(name string, limiter *Limiter) {
	pool.mu.Lock()
	defer pool.mu.Unlock()
	pool.limiters[name] = limiter
}

func GetLimiter(name string) *Limiter {
	pool.mu.RLock()
	defer pool.mu.RUnlock()
	return pool.limiters[name]
}

func Middleger(limiter *Limiter, keyFunc LimiterKey, logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			key := keyFunc(r)

			allowed, err := limiter.Allow(key)
			if err != nil {
				if logger != nil {
					logger.Warn("rate limit check failed, allowing request",
						slog.String("key", key),
						slog.Any("error", err))
				}
			}

			if !allowed {
				retryAfter := "1"
				w.Header().Set("Retry-After", retryAfter)
				w.Header().Set("X-RateLimit-Limit", "10")
				w.Header().Set("X-RateLimit-Remaining", "0")
				w.Header().Set("X-RateLimit-Reset", time.Now().Add(time.Second).Format(time.RFC1123))
				http.Error(w, "rate limit exceeded", http.StatusTooManyRequests)
				return
			}

			w.Header().Set("X-RateLimit-Limit", "10")
			w.Header().Set("X-RateLimit-Remaining", "0")
			w.Header().Set("X-RateLimit-Reset", time.Now().Add(time.Second).Format(time.RFC1123))

			next.ServeHTTP(w, r)
		})
	}
}