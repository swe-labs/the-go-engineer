# Authentication vs authorization

## Learning objective

Distinguish authentication (proving identity) from authorization (granting permission), implement them as separate middleware layers in Go, and use claims-based authorization to make permission decisions at the resource level.

## Why this matters

Confusing authentication with authorization is the root cause of countless security vulnerabilities. A system that authenticates a user but does not check authorization allows any authenticated user to perform any action. In 2021, a major cloud provider suffered a data breach because an internal API authenticated the caller but did not verify they were authorized to access the specific tenant's data. In Go services, separating authN and authZ into distinct middleware layers makes each independently testable, auditable, and replaceable.

## Mental model

Authentication is the bouncer checking your ID at the club door: "are you who you say you are?" Authorization is the host inside checking your reservation: "are you allowed at this table?" The bouncer does not decide which table you can sit at; the host does not verify your identity. Two separate roles, two separate checks. In HTTP terms: authN failure returns 401 (Unauthorized — the HTTP name is misleading). AuthZ failure returns 403 (Forbidden). A 401 means "who are you?" A 403 means "I know who you are, but you cannot do that."

## Core idea

- **Authentication (AuthN)**: verifies identity. Answers "who are you?" Methods: passwords, JWTs, session cookies, biometrics, OAuth2 tokens.
- **Authorization (AuthZ)**: verifies permission. Answers "what can you do?" Models: RBAC (roles), ABAC (attributes), ACLs (access control lists), claims-based.

The two must be separated because:

1. A user can be authenticated but not authorized (e.g., a free-tier user trying to access premium features).
2. AuthN mechanisms change (passwords -> SSO) independently of authZ policies.
3. AuthZ policies change (promotions, role changes) without affecting how users prove identity.

Claims-based authorization embeds identity attributes (roles, permissions, tenant ID) in the authentication token. After authN, the middleware extracts claims and passes them to the authZ layer.

## Under the hood

In Go, the separation is implemented as two middleware functions wrapping each other:

```
Request -> AuthN Middleware -> AuthZ Middleware -> Handler
   |                            |
   | - Validates token          | - Reads identity from context
   | - Extracts claims          | - Compares claims to required permissions
   | - Stores identity in ctx   | - Returns 403 if insufficient
   | - Returns 401 if invalid   |
```

AuthN middleware stores the verified identity in `context.Context`. AuthZ middleware reads it and makes a decision. This is the standard pattern in Go web frameworks:

```go
func authNMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        token := extractToken(r)
        claims, err := verifyToken(token)
        if err != nil {
            http.Error(w, "unauthorized", http.StatusUnauthorized)
            return
        }
        ctx := context.WithValue(r.Context(), "claims", claims)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}

func authZMiddleware(requiredPerm string) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            claims := r.Context().Value("claims").(Claims)
            if !hasPermission(claims, requiredPerm) {
                http.Error(w, "forbidden", http.StatusForbidden)
                return
            }
            next.ServeHTTP(w, r)
        })
    }
}
```

## How Go uses it

Go's `net/http` makes this separation natural. The standard library does not provide auth middleware, but the idiom is universal:

- **Echo**: `e.Use(middleware.JWT(...))` for authN, then `e.GET("/admin", adminHandler, roleMiddleware("admin"))` for authZ.
- **Gin**: `r.Use(authNMiddleware())`, `r.GET("/admin", authZMiddleware("admin"), adminHandler)`.
- **Chi**: `r.Group(func(r chi.Router) { r.Use(authNMiddleware); r.Use(adminOnlyMiddleware) })`.

The key insight: authN is applied globally (or per-route group), while authZ is applied per-route or per-resource.

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
)

type contextKey string

const claimsKey contextKey = "claims"

// Claims represent the authenticated user's identity and permissions.
type Claims struct {
	UserID string   `json:"user_id"`
	Role   string   `json:"role"`
	Perms  []string `json:"permissions"`
	Tenant string   `json:"tenant"`
}

// AuthN middleware: verifies identity.
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

// AuthZ middleware: verifies permission. Returns a 403 on failure.
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

// verifyToken simulates JWT verification. In production, use golang-jwt.
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

// hasPermission checks a specific permission in claims.
func hasPermission(c Claims, required string) bool {
	for _, p := range c.Perms {
		if p == required {
			return true
		}
		// Wildcard: doc:* matches doc:read, doc:write, etc.
		if strings.HasSuffix(p, ":*") {
			prefix := strings.TrimSuffix(p, ":*")
			if strings.HasPrefix(required, prefix+":") {
				return true
			}
		}
	}
	return false
}

// Handlers (no auth logic, just business logic).
func listDocsHandler(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value(claimsKey).(Claims)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"user":   claims.UserID,
		"tenant": claims.Tenant,
		"docs":   []string{"doc1.md", "doc2.md"},
	})
}

func createDocHandler(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value(claimsKey).(Claims)
	json.NewEncoder(w).Encode(map[string]string{
		"status":   "created",
		"by":       claims.UserID,
		"document": "new-doc.md",
	})
}

func deleteDocHandler(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value(claimsKey).(Claims)
	json.NewEncoder(w).Encode(map[string]string{
		"status":   "deleted",
		"by":       claims.UserID,
		"document": strings.TrimPrefix(r.URL.Path, "/api/docs/"),
	})
}

func main() {
	mux := http.NewServeMux()

	// Public endpoint (no authN or authZ).
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Protected API group: authN first, then per-route authZ.
	api := http.NewServeMux()
	api.Handle("/docs", authZMiddleware("doc:read")(http.HandlerFunc(listDocsHandler)))
	api.Handle("/docs/create", authZMiddleware("doc:write")(http.HandlerFunc(createDocHandler)))
	api.Handle("/docs/", authZMiddleware("doc:delete")(http.HandlerFunc(deleteDocHandler)))

	mux.Handle("/api/", authNMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = strings.TrimPrefix(r.URL.Path, "/api")
		api.ServeHTTP(w, r)
	})))

	log.Println("Server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
```

## Step-by-step execution

For a request `GET /api/docs` with `Authorization: Bearer admin-token`:

1. Request enters `mux`, routed to `/api/` handler.
2. `authNMiddleware` extracts the Bearer token, calls `verifyToken("admin-token")`.
3. `verifyToken` returns Claims `{UserID: "u1", Role: "admin", Perms: ["doc:read", "doc:write", "doc:delete", "user:manage"]}`.
4. Claims stored in context. AuthN passes.
5. Request forwarded to `api.ServeHTTP`, routes to `/docs`.
6. `authZMiddleware("doc:read")` reads claims from context.
7. `hasPermission` iterates permissions, finds `"doc:read"`. AuthZ passes.
8. `listDocsHandler` processes the request normally.

For `GET /api/docs` with `Bearer readonly-token`:
1-5. Same as above.
7. `hasPermission` finds `"doc:read"`. Authorized.

For `DELETE /api/docs/secret.doc` with `Bearer readonly-token`:
1-5. Same pattern.
6. `authZMiddleware("doc:delete")` reads claims: only `["doc:read"]`.
7. `hasPermission` returns false.
8. Response: 403 Forbidden.

## Common mistakes

- Confusing authentication (who you are) with authorization (what you can do) — returning 401 for permission failures instead of 403.
- Implementing authorization as an afterthought in handlers instead of a separate middleware layer — leads to copy-paste inconsistencies.
- Using the same middleware for both auth and authZ — mixing them makes either hard to test, audit, or replace independently.
- Making authorization decisions based on role alone without resource-level context — an admin should not automatically access all resources across all tenants.
- Storing authorization logic in the database layer — authorization is a middleware concern, not a query concern.

## Debugging walkthrough

Consider this code where authN and authZ are merged into a single function:

```go
func authMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        token := extractToken(r)
        user := findUserByToken(token)
        if user == nil {
            http.Error(w, "forbidden", http.StatusForbidden) // Wrong status!
            return
        }
        if user.Role != "admin" {
            http.Error(w, "forbidden", http.StatusForbidden)
            return
        }
        next.ServeHTTP(w, r)
    })
}
```

**Symptom**: All authenticated non-admin users receive 403 on every route, including public-facing ones. A valid user who has not authenticated (no token) also receives 403 instead of 401.

**Investigation**: The middleware mixes authN ("do you exist?") with authZ ("are you admin?"). Both failures return 403. A client library that sees 403 might retry with different credentials, but a 401 would trigger a fresh login. The merged middleware also prevents non-admin users from accessing user-level routes.

**Root cause**: Single middleware doing both jobs. No separation of concerns.

**Fix**: Split into two middlewares:

```go
func authNMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        token := extractToken(r)
        user := findUserByToken(token)
        if user == nil {
            http.Error(w, "unauthorized", http.StatusUnauthorized) // Correct status
            return
        }
        ctx := context.WithValue(r.Context(), "user", user)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}

func adminOnlyMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        user := r.Context().Value("user").(User)
        if user.Role != "admin" {
            http.Error(w, "forbidden", http.StatusForbidden)
            return
        }
        next.ServeHTTP(w, r)
    })
}
```

## Production notes

- AuthN middleware should run on every protected route. AuthZ middleware runs on specific routes with specific requirements.
- Use context keys with typed constants (e.g., `type contextKey string`) to avoid collisions between middlewares.
- Always return the correct HTTP status: 401 for missing/invalid credentials, 403 for valid credentials but insufficient permission.
- For claims-based authZ, embed all attributes needed for authorization (roles, permissions, tenant ID) in the JWT or session. Avoid database lookups in the authZ middleware.
- Audit all authZ failures at WARN level. A sudden increase in 403s may indicate a brute-force or privilege escalation attempt.

## Performance implications

- AuthN middleware: one token verification per request. JWT with RSA256 is ~0.5ms. Session token lookup in Redis is ~1-5ms.
- AuthZ middleware: in-memory permission check is O(n) where n is the number of permissions per user. With 50 permissions, this is sub-microsecond.
- Claims extraction from tokens is essentially free after the token is verified.
- If authZ requires a database lookup (e.g., fetch user's team memberships), add 5-50ms per request. Cache team membership data with a 5-minute TTL.

## Practice task

Write Go functions `AuthN(token string) (Claims, error)` and `AuthZ(claims Claims, requiredPerm string, resourceOwner string) bool` where:
- `AuthN` validates the token and returns claims containing `UserID`, `Role`, and `Perms`.
- `AuthZ` checks if the user has the required permission AND if they own the resource OR have an `admin` role.
- Then write a `main()` that simulates three scenarios: admin accessing any resource, owner accessing own resource, and non-owner accessing another's resource.

## Tests / verification

```bash
go run ./curriculum/modules/10-auth-security/lessons/04-authentication-vs-authorization
go test ./curriculum/modules/10-auth-security/lessons/04-authentication-vs-authorization
```

The test file `main_test.go` contains table-driven tests that verify:
- `AuthN` returns claims for a valid token and error for an invalid token.
- `AuthZ` returns true for admin on any resource.
- `AuthZ` returns true for resource owner with matching permission.
- `AuthZ` returns false for non-owner without admin role.

## Review questions

1. What HTTP status code should an authN failure return? What about an authZ failure?
2. Why is it dangerous to merge authN and authZ into a single middleware function?
3. In claims-based authorization, what information must the JWT contain to make authZ decisions without a database lookup?
4. A user authenticates successfully but receives a 403 on a resource they created yesterday. What could have changed?
5. How would you implement resource-level authorization (user A can edit document X but not document Y) in Go middleware?

## NEXT UP

Password hashing — why plaintext storage is catastrophic and how bcrypt, scrypt, and Argon2 provide cryptographic protection for user credentials.
