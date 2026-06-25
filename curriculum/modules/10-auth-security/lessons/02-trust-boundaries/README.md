# Trust boundaries

## Learning objective

Identify trust boundaries in Go web service architectures, implement middleware-based gate patterns at each boundary, and apply a zero-trust mindset where every request is authenticated and authorized regardless of origin.

## Why this matters

Every data flow in a distributed system crosses trust boundaries: from the public internet to your load balancer, from your load balancer to your application, from your application to your database, and from your application to third-party APIs. Each boundary is an opportunity for an attacker to inject, spoof, or escalate. Engineers who understand trust boundaries build defenses at the right layer instead of relying on a single perimeter. The Capital One breach (2019) exploited a missing trust boundary between a WAF and a metadata service — a single trust boundary mistake leaked 100 million records.

## Mental model

Picture a medieval castle with concentric walls. The outermost wall faces the wilderness (public internet). Inside that wall is the bailey (DMZ), where merchants and travellers are screened. The inner wall surrounds the keep (your application logic), and the innermost chamber (the vault) holds the treasure (user data). Each gate between these zones is a trust boundary: anyone entering must present credentials, and the gatekeeper checks them. In modern architecture, every service-to-service call, every database query, and every goroutine that handles user data crosses a trust boundary.

## Core idea

A trust boundary is any point where data passes from a lower trust level to a higher trust level. In a typical Go web service, trust boundaries exist at:

- **HTTP listener**: data enters from the public network (untrusted) into your process (semi-trusted).
- **Auth middleware**: unauthenticated requests (untrusted) become authenticated requests (trusted but not authorized).
- **Database driver**: application queries (trusted logic) become SQL statements (executed with database permissions).
- **RPC client**: your service (trusted) calls another service (separate trust domain).
- **File system**: uploaded files (untrusted) are read by your application (trusted).

The zero-trust principle says: do not trust any boundary implicitly. Verify every request, even from inside the network.

## Under the hood

Trust boundaries in Go are implemented using the middleware pattern. A middleware is a function that wraps an `http.Handler` and returns a new `http.Handler`. Each middleware enforces a trust check:

```
Request -> Middleware 1 (TLS) -> Middleware 2 (AuthN) -> Middleware 3 (AuthZ) -> Handler
```

The `context.Context` carries identity across boundaries. Once a middleware authenticates a request, it stores the user identity in the context. Downstream handlers and middlewares read the identity without re-authenticating.

Go's `net/http` server creates a new goroutine per connection. Each goroutine has its own stack and context — this is a lightweight trust boundary between concurrent requests.

## How Go uses it

Go's standard library provides the building blocks for trust boundary enforcement:

- **`http.Handler` and `http.HandlerFunc`**: the composable unit for middleware chains.
- **`context.Context`**: carries authenticated identity across boundaries.
- **`crypto/tls`**: enforces a trust boundary at the transport layer.
- **`net/http/httputil.ReverseProxy`**: a trust boundary between your service and upstream services.
- **`database/sql`**: parameterized queries enforce a trust boundary between application and database — the query text is never concatenated with user input.

Production Go frameworks like Chi, Echo, and Gin all use middleware chains that map directly to trust boundary layers.

## Go example

```go
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

type contextKey string

const userKey contextKey = "user"

// User represents an authenticated identity.
type User struct {
	ID    string
	Role  string
	Email string
}

// Trust boundary 1: TLS termination (assumed at load balancer).
// Trust boundary 2: Authentication middleware.
func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract token from Authorization header.
		auth := r.Header.Get("Authorization")
		if auth == "" || !strings.HasPrefix(auth, "Bearer ") {
			http.Error(w, `{"error":"missing token"}`, http.StatusUnauthorized)
			return
		}
		token := strings.TrimPrefix(auth, "Bearer ")
		// In production: validate JWT or session token here.
		// For this example, a simple token-to-user mapping.
		user, ok := validateToken(token)
		if !ok {
			http.Error(w, `{"error":"invalid token"}`, http.StatusUnauthorized)
			return
		}
		// Store user in context (crosses trust boundary into handler).
		ctx := context.WithValue(r.Context(), userKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Trust boundary 3: Authorization middleware (separate from auth).
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

// Trust boundary 4: handler processes authenticated, authorized request.
func profileHandler(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(userKey).(User)
	if !ok {
		http.Error(w, `{"error":"unauthenticated"}`, http.StatusUnauthorized)
		return
	}
	// Database trust boundary: parameterized query prevents injection.
	json.NewEncoder(w).Encode(map[string]string{
		"user_id": user.ID,
		"email":   user.Email,
		"role":    user.Role,
	})
}

func adminHandler(w http.ResponseWriter, r *http.Request) {
	user, _ := r.Context().Value(userKey).(User)
	json.NewEncoder(w).Encode(map[string]string{
		"message": fmt.Sprintf("Welcome admin %s", user.Email),
		"access":  "admin_panel",
	})
}

func validateToken(token string) (User, bool) {
	// Simulated token validation.
	tokens := map[string]User{
		"token-user":  {ID: "u1", Role: "user", Email: "alice@example.com"},
		"token-admin": {ID: "u2", Role: "admin", Email: "bob@example.com"},
	}
	u, ok := tokens[token]
	return u, ok
}

func main() {
	mux := http.NewServeMux()
	// Public endpoint (no trust boundary check).
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Protected endpoints with auth boundary.
	protected := http.NewServeMux()
	protected.HandleFunc("/profile", profileHandler)
	protected.HandleFunc("/admin", adminHandler)

	// Middleware chain: auth -> optional role check.
	mux.Handle("/api/", authMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Strip /api prefix and route to protected mux.
		r.URL.Path = strings.TrimPrefix(r.URL.Path, "/api")
		protected.ServeHTTP(w, r)
	})))

	// Admin-only sub-boundary.
	mux.Handle("/api/admin", authMiddleware(requireRole("admin")(http.HandlerFunc(adminHandler))))

	server := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	log.Printf("Server starting on %s", server.Addr)
	log.Fatal(server.ListenAndServe())
}
```

## Step-by-step execution

For a request `GET /api/profile` with `Authorization: Bearer token-user`:

1. Request arrives at `:8080`. TLS termination trust boundary (assumed at load balancer).
2. `mux` routes to `/api/` handler, which delegates to `authMiddleware`.
3. `authMiddleware` extracts the token, calls `validateToken`. Trust boundary crossed: raw HTTP -> authenticated identity.
4. Valid user `{ID: "u1", Role: "user"}` is stored in `context.WithValue`.
5. Request passes to `protected.ServeHTTP` with authenticated context. Another trust boundary: middleware context -> handler.
6. `profileHandler` reads user from context. Trust boundary between handler and data access: no raw SQL, uses parameterized queries.
7. Response flows back through the same boundaries in reverse.

For a request to `GET /api/admin` with `token-user`:

1-3. Same as above.
4. `authMiddleware` authenticates user (Role: "user").
5. `requireRole("admin")` middleware checks role: user has `"user"`, not `"admin"`.
6. Returns 403 Forbidden. Trust boundary enforced: authorization denied.

## Common mistakes

- Drawing trust boundaries around deployment boundaries instead of data flow boundaries — a Kubernetes namespace boundary does not automatically make data trustable.
- Assuming the network perimeter is the only trust boundary — in zero-trust architectures, boundaries exist between services, between pods, and even between goroutines.
- Placing the trust boundary inside the request handler instead of at the middleware layer — every handler reimplements auth checks, leading to inconsistent enforcement.
- Mixing trust levels in the same handler — same function handling public data and admin-only mutations.
- Not propagating context across goroutine boundaries — a `go func()` that does not receive the parent context loses the authenticated identity.

## Debugging walkthrough

Consider this code where the trust boundary is misplaced:

```go
func handler(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	// Trust boundary expected: user_id should come from auth, not URL!
	rows, _ := db.Query("SELECT * FROM documents WHERE owner_id = $1", userID)
}
```

**Symptom**: Any authenticated user can see any other user's documents by changing the `user_id` query parameter.

**Investigation**: The trust boundary between authentication and data access is broken. The handler reads `user_id` from the URL query (untrusted input) instead of from the authenticated session context.

**Root cause**: The developer placed the trust boundary at the database query (parameterized query prevents injection) but forgot the authorization boundary: the `owner_id` value must come from the authenticated session, not from user input.

**Fix**: Remove `user_id` from the URL and read it from the authenticated context:

```go
func handler(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(userKey).(User)
	rows, _ := db.Query("SELECT * FROM documents WHERE owner_id = $1", user.ID)
}
```

## Production notes

- Every handler should be behind at least one trust boundary middleware. The only exception is the health check endpoint.
- Use a middleware chain library like `chi` or `alice` to compose boundaries declaratively.
- In gRPC services, use unary and stream interceptors as trust boundaries. Each RPC crosses a trust boundary.
- For database trust boundaries, always use parameterized queries. Raw string concatenation across a trust boundary is a vulnerability.
- Log trust boundary crossings for audit: "user X crossed boundary Y at time Z". Do not log the data itself.

## Performance implications

Middleware-based trust boundaries add overhead proportional to the number of middlewares:

- Auth middleware: one token decode/verify per request. JWT verification with RSA is ~0.5ms. bcrypt verification is ~250ms.
- TLS handshake: ~1-5ms for full handshake, negligible for session resumption.
- Database trust boundary: parameterized queries have no overhead compared to string concatenation.

For high-throughput services, pre-validate tokens at the gateway/ingress level so internal services trust the already-validated identity header.

## Practice task

Write a Go function `BuildMiddlewareChain(trustBoundaries []string) func(http.Handler) http.Handler` that:
- Accepts a list of trust boundary names (e.g., `["tls", "auth", "rate_limit"]`).
- Returns a middleware that chains functions for each boundary.
- Each boundary logs a message when crossed (use `log.Printf`).
- Then write a `main()` that creates a server with at least four trust boundaries and serves a protected `/data` endpoint.

## Tests / verification

```bash
go run ./curriculum/modules/10-auth-security/lessons/02-trust-boundaries
go test ./curriculum/modules/10-auth-security/lessons/02-trust-boundaries
```

The test file `main_test.go` contains table-driven tests that verify:
- `BuildMiddlewareChain` returns a non-nil middleware for any non-empty boundary list.
- A request without a token receives a 401.
- A request with an invalid token receives a 401.
- A request with a valid token and insufficient role receives a 403.
- A request with a valid token and sufficient role receives a 200.

## Review questions

1. What is a trust boundary? Give three examples in a typical Go web service.
2. Why is `context.Context` essential for trust boundary enforcement across goroutines?
3. In the debugging walkthrough, why does reading `user_id` from the URL query parameter violate the trust boundary?
4. How does the zero-trust principle change where you place trust boundaries compared to a perimeter-based model?
5. A middleware authenticates a request and stores the user in context. A second middleware reads the user from context. What trust boundary exists between them, and how is it enforced?

## NEXT UP

Input validation for security — how allowlists, denylists, and sanitization gates prevent malicious data from crossing trust boundaries.
