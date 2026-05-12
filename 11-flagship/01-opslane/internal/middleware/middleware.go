// Copyright (c) 2026 Rasel Hossen
// See LICENSE for usage terms.

// Package middleware provides HTTP middleware components for the Opslane server.
// It includes panic recovery, request logging, CORS, rate limiting, and security headers.
package middleware

import (
	"log/slog"
	"net"
	"net/http"
	"net/netip"
	"runtime/debug"
	"strings"
	"sync"
	"time"
)

// ============================================================================
// Internal Package: Middleware
// Level: Advanced
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Wrapping Handlers in structured logging (`slog`)
//   - Defending the server from Panics
//   - Adding critical Content-Type security headers to prevent Cross-site scripting (XSS)
//
// ENGINEERING DEPTH:
//   "Layering" middleware wraps your core application logic in an onion-like shell.
//   Every request *must* pass through these layers first before hitting `ServeMux`.
// ============================================================================

// RecoverPanic (Function): outermost middleware that catches panics, logs stack trace, and returns 500
func RecoverPanic(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					// Hard-close the TCP connection since the state is unstable.
					w.Header().Set("Connection", "close")
					logger.Error("panic caught by middleware",
						slog.Any("error", err),
						slog.String("stack", string(debug.Stack())),
					)
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				}
			}()
			next.ServeHTTP(w, r) // Execute downstream chain
		})
	}
}

// LogRequest (Function): lightweight access log middleware logging method, path, and latency
func LogRequest(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			next.ServeHTTP(w, r)
			logger.Info("request",
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.Duration("latency", time.Since(start)),
			)
		})
	}
}

// CORS (Function): middleware allowing browser-based API calls during local development
func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// clientWindow (Struct): tracks request count and reset time for a single client's rate limit window
type clientWindow struct {
	count   int
	resetAt time.Time
}

// RateLimit (Function): per-IP rate limiting middleware with fixed window and proxy-aware client address
func RateLimit(maxRequests int, window time.Duration, trustedProxyCIDRs []netip.Prefix) func(http.Handler) http.Handler {
	var mu sync.Mutex
	clients := make(map[string]clientWindow)
	nextCleanup := time.Now().Add(window)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if maxRequests <= 0 || window <= 0 {
				next.ServeHTTP(w, r)
				return
			}

			clientIP := ClientAddress(r, trustedProxyCIDRs)
			now := time.Now()

			mu.Lock()
			if !now.Before(nextCleanup) {
				pruneExpiredClientWindows(clients, now)
				nextCleanup = now.Add(window)
			}

			state := clients[clientIP]
			if state.resetAt.IsZero() || now.After(state.resetAt) {
				state = clientWindow{resetAt: now.Add(window)}
			}
			state.count++
			clients[clientIP] = state
			allowed := state.count <= maxRequests
			mu.Unlock()

			if !allowed {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusTooManyRequests)
				_, _ = w.Write([]byte(`{"error":{"code":"rate_limited","message":"rate limit exceeded"}}`))
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// pruneExpiredClientWindows (Function): removes rate limit entries whose window has expired
func pruneExpiredClientWindows(clients map[string]clientWindow, now time.Time) {
	for clientIP, state := range clients {
		if !now.Before(state.resetAt) {
			delete(clients, clientIP)
		}
	}
}

// SecureHeaders (Function): middleware setting mandatory browser security headers (X-Frame-Options, X-Content-Type-Options)
func SecureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Prevent Clickjacking (disallow being loaded in an iframe)
		w.Header().Set("X-Frame-Options", "deny")
		// Prevent MIME-sniffing bypasses
		w.Header().Set("X-Content-Type-Options", "nosniff")

		next.ServeHTTP(w, r)
	})
}

// ClientAddress (Function): extracts the client IP address from the request.
// If trustedProxyCIDRs is configured and the request comes from a trusted proxy,
// it extracts the client IP from X-Forwarded-For or X-Real-IP headers.
// Otherwise, it uses the direct remote address.
func ClientAddress(r *http.Request, trustedProxyCIDRs []netip.Prefix) string {
	peerIP := remotePeerIP(r.RemoteAddr)
	if peerIP.IsValid() && isTrustedProxy(peerIP, trustedProxyCIDRs) {
		if forwardedIP, ok := forwardedClientIP(r); ok {
			return forwardedIP.String()
		}
	}

	if peerIP.IsValid() {
		return peerIP.String()
	}

	return r.RemoteAddr
}

// remotePeerIP (Function): extracts the IP address from a remote address string
func remotePeerIP(remoteAddr string) netip.Addr {
	host, _, err := net.SplitHostPort(remoteAddr)
	if err == nil {
		if ip, parseErr := netip.ParseAddr(host); parseErr == nil {
			return ip.Unmap()
		}
	}

	ip, err := netip.ParseAddr(remoteAddr)
	if err != nil {
		return netip.Addr{}
	}

	return ip.Unmap()
}

// forwardedClientIP (Function): extracts the client IP from X-Forwarded-For or X-Real-IP headers
func forwardedClientIP(r *http.Request) (netip.Addr, bool) {
	if forwardedFor := strings.TrimSpace(r.Header.Get("X-Forwarded-For")); forwardedFor != "" {
		firstHop := strings.TrimSpace(strings.Split(forwardedFor, ",")[0])
		if ip, err := netip.ParseAddr(firstHop); err == nil {
			return ip.Unmap(), true
		}
	}

	if realIP := strings.TrimSpace(r.Header.Get("X-Real-Ip")); realIP != "" {
		if ip, err := netip.ParseAddr(realIP); err == nil {
			return ip.Unmap(), true
		}
	}

	return netip.Addr{}, false
}

// isTrustedProxy (Function): checks if an IP address is within any of the trusted proxy CIDR ranges
func isTrustedProxy(ip netip.Addr, trustedProxyCIDRs []netip.Prefix) bool {
	for _, prefix := range trustedProxyCIDRs {
		if prefix.Contains(ip) {
			return true
		}
	}

	return false
}
