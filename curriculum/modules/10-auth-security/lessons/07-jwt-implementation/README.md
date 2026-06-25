# JWT implementation

## Learning objective

Issue and verify JSON Web Tokens in Go using the `golang-jwt/jwt` library, configure RSA or HMAC signing, embed custom claims, and compose JWT verification middleware.

## Why this matters

JWTs enable stateless authentication: the token itself contains the user's identity and permissions, verified by a digital signature. No server-side session store is needed — every service that has the public key can verify a token independently. This makes JWTs the standard for microservice authentication, OAuth2 access tokens, OIDC ID tokens, and API authentication. Stripe, GitHub, Google, and Auth0 all use JWTs. Go's performance and concurrency model make it an ideal language for JWT issuing and verification at scale.

## Mental model

A JWT is a digitally signed passport. The passport contains your photo and personal details (claims), is stamped with an official seal (signature), and has an expiry date. Any official (service) in any country (microservice) can verify the seal using the public key book (JWKS endpoint) without calling the issuing country's embassy (central auth server). The passport is self-sufficient: its validity is proven by the signature, not by a database lookup. But if the passport is stolen, the thief can use it until it expires — revocation requires a denylist.

## Core idea

A JWT is three base64url-encoded segments separated by dots:

```
header.payload.signature
```

- **Header**: algorithm and token type `{"alg":"RS256","typ":"JWT"}`.
- **Payload**: claims — registered (`iss`, `sub`, `exp`, `iat`), public (`name`, `role`), private (`tenant_id`).
- **Signature**: cryptographic output of `HMACSHA256(base64(header) + "." + base64(payload), secret)` for HMAC, or `RSASHA256(...)` for RSA.

Two signing approaches:

| Method | Key type | Verification | Use case |
|---|---|---|---|
| HMAC (HS256) | Single shared secret | Same secret signs and verifies | Single service, trust boundary |
| RSA (RS256) | Private/public key pair | Public key verifies, private key signs | Multi-service, third-party verification |
| ECDSA (ES256) | Private/public key pair | Smaller signatures than RSA | Mobile, constrained environments |

`golang-jwt/jwt` (formerly `dgrijalva/jwt-go`) is the de facto Go library. The v5 API uses functional options for token parsing:

```go
token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
    return publicKey, nil
})
```

## Under the hood

When a JWT is verified, the library decodes the base64url header and payload without validating them, extracts the `alg` from the header, looks up the key via the key function, and recomputes the signature. If the computed signature matches the provided signature and the token is not expired, the token is valid. For RS256, signature verification uses RSA with SHA-256 under the `crypto/rsa` and `crypto/sha256` packages. The library does not decrypt anything — JWTs are signed, not encrypted. To protect sensitive claims, use JWE (JSON Web Encryption) or transport encryption (TLS).

## How Go uses it

The `golang-jwt/jwt/v5` package provides:

- `jwt.NewWithClaims(method, claims)`: create a new token.
- `token.SignedString(key)`: produce the signed token string.
- `jwt.Parse(tokenString, keyFunc)`: parse and verify a token.
- Registered claims via `jwt.RegisteredClaims`: `Issuer`, `Subject`, `Audience`, `ExpiresAt`, `IssuedAt`, `NotBefore`, `ID`.
- Custom claims via embedding: `type CustomClaims struct { jwt.RegisteredClaims; Role string; TenantID string }`.

## Go example

```go
package main

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// CustomClaims extends the standard registered claims.
type CustomClaims struct {
	jwt.RegisteredClaims
	Role     string `json:"role"`
	TenantID string `json:"tenant_id"`
}

var (
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
)

func init() {
	// Generate an RSA key pair for signing.
	// In production, load from a secure file or env variable.
	var err error
	privateKey, err = rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Fatalf("Failed to generate key: %v", err)
	}
	publicKey = &privateKey.PublicKey
}

// IssueToken creates a signed JWT for a user.
func IssueToken(userID, role, tenantID string, ttl time.Duration) (string, error) {
	now := time.Now()
	claims := CustomClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "the-go-engineer",
			Subject:   userID,
			Audience:  jwt.ClaimStrings{"api.thegoengineer.com"},
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
			IssuedAt:  jwt.NewNumericDate(now),
			ID:        fmt.Sprintf("%s-%d", userID, now.UnixNano()),
		},
		Role:     role,
		TenantID: tenantID,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(privateKey)
}

// ValidateToken parses and verifies a JWT, returning the claims.
func ValidateToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method.
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return publicKey, nil
	})
	if err != nil {
		return nil, fmt.Errorf("token validation failed: %w", err)
	}
	claims, ok := token.Claims.(*CustomClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token claims")
	}
	return claims, nil
}

// authMiddleware verifies JWTs on protected routes.
func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if auth == "" || !strings.HasPrefix(auth, "Bearer ") {
			http.Error(w, `{"error":"missing token"}`, http.StatusUnauthorized)
			return
		}
		tokenStr := strings.TrimPrefix(auth, "Bearer ")
		claims, err := ValidateToken(tokenStr)
		if err != nil {
			http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusUnauthorized)
			return
		}
		// In production, store claims in context.
		_ = claims
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message":   "authenticated",
			"user_id":   claims.Subject,
			"role":      claims.Role,
			"tenant_id": claims.TenantID,
			"expires":   claims.ExpiresAt,
		})
	})
}

// tokenHandler issues a token for demo purposes.
func tokenHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		UserID   string `json:"user_id"`
		Role     string `json:"role"`
		TenantID string `json:"tenant_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	token, err := IssueToken(req.UserID, req.Role, req.TenantID, 1*time.Hour)
	if err != nil {
		http.Error(w, "token generation failed", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/token", tokenHandler)
	mux.Handle("/protected", authMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"data": "sensitive information"})
	})))

	log.Println("Server on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}

// Standalone functions for testing.

// CreateToken is a convenience wrapper for testing: returns a signed JWT string.
func CreateToken(userID, role string, ttl time.Duration) (string, error) {
	return IssueToken(userID, role, "default-tenant", ttl)
}

// VerifyToken returns the subject (user ID) from a valid token, or an error.
func VerifyToken(tokenStr string) (string, error) {
	claims, err := ValidateToken(tokenStr)
	if err != nil {
		return "", err
	}
	return claims.Subject, nil
}
```

## Step-by-step execution

For `POST /token` with `{"user_id":"alice","role":"admin","tenant_id":"acme-corp"}`:

1. `tokenHandler` decodes the JSON request.
2. `IssueToken("alice", "admin", "acme-corp", 1*time.Hour)` creates a `CustomClaims` struct:
   - `Subject`: "alice", `ExpiresAt`: now + 1 hour, `IssuedAt`: now, `ID`: unique.
   - `Role`: "admin", `TenantID`: "acme-corp".
3. `jwt.NewWithClaims(jwt.SigningMethodRS256, claims)` creates a token object.
4. `token.SignedString(privateKey)` signs the header+payload with RSA private key, producing the three-part token.
5. Token returned to client.

For `GET /protected` with `Authorization: Bearer <token>`:

1. `authMiddleware` extracts the token from the Authorization header.
2. `ValidateToken(tokenStr)` splits the token, base64-decodes header and payload, checks expiry.
3. Verifies the RSA signature against `publicKey`. If the token was tampered with, signature verification fails.
4. If valid, claims are extracted and the handler processes the request.
5. If invalid, returns 401.

## Common mistakes

- Not validating the algorithm header — accepting `alg: none` tokens bypasses verification entirely. `golang-jwt/jwt` v5 explicitly requires a signing method.
- Using symmetric signing (HS256) for multi-service architectures — the shared secret must be distributed to every service, increasing exposure.
- Not setting expiration (`exp` claim) — tokens that never expire remain valid indefinitely if leaked.
- Storing sensitive data in JWT claims — the payload is base64-encoded, not encrypted. Anyone with the token can read claims.
- Not using `jwt.ParseWithClaims` — using `jwt.Parse` and casting claims manually misses claim validation.

## Debugging walkthrough

Consider this token verification code:

```go
func ValidateToken(tokenString string) (jwt.MapClaims, error) {
    token, err := jwt.Parse(tokenString, nil)
    if err != nil {
        return nil, err
    }
    return token.Claims.(jwt.MapClaims), nil
}
```

**Symptom**: Tokens with `alg: none` are accepted. The server is vulnerable to the "none algorithm" attack.

**Investigation**: Passing `nil` as the key function means the library does not enforce any signing method. An attacker can set `alg` to `none`, omit the signature, and present any claims.

**Root cause**: The key function must validate the signing algorithm. `nil` disables verification.

**Fix**:

```go
func ValidateToken(tokenString string) (*CustomClaims, error) {
    token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }
        return publicKey, nil
    })
    // ...
}
```

## Production notes

- Store private keys in a secrets manager (HashiCorp Vault, AWS Secrets Manager, GCP Secret Manager). Never commit them to a repository.
- Use RS256 (RSA) for multi-service architectures. Use ES256 (ECDSA) for smaller tokens in mobile environments. Use HS256 only for single-service applications.
- Set `exp`, `iat`, `nbf`, and `jti` claims. The `jti` (unique identifier) enables token denylisting.
- Rotate signing keys regularly. Maintain a grace period where old public keys are still accepted for verification.
- Load public keys from a JWKS endpoint for distributed verification. The `golang-jwt/jwt` package supports JWKS via `github.com/MicahParks/keyfunc/v3`.

## Performance implications

- RSA 2048-bit signing: ~1-5ms per token (private key operation).
- RSA 2048-bit verification: ~0.1-0.5ms per token (public key operation).
- ECDSA P-256 signing: ~0.3ms, verification: ~0.1ms.
- HMAC-SHA256 signing/verification: ~0.001ms (very fast, symmetric).
- Token parsing (base64 decode, JSON unmarshal): ~0.01ms.

For high-throughput services, cache the parsed token for the duration of the request. Token parsing is cheap, but the signature verification is the bottleneck.

## Practice task

Write Go functions `CreateJWT(userID, role string, secret []byte, ttl time.Duration) (string, error)` and `VerifyJWT(tokenStr string, secret []byte) (userID string, role string, err error)` that:
- Use HMAC-SHA256 (HS256) for signing.
- Include `sub`, `role`, `exp`, `iat`, and `jti` claims.
- `VerifyJWT` validates the signature, expiry, and returns the subject and role.
- Then write a `main()` that creates a token, verifies it, tries to verify with a wrong secret (must fail), and tries to verify an expired token (must fail).

## Tests / verification

```bash
go run ./curriculum/modules/10-auth-security/lessons/07-jwt-implementation
go test ./curriculum/modules/10-auth-security/lessons/07-jwt-implementation
```

The test file `main_test.go` contains table-driven tests that verify:
- `CreateJWT` returns a token with three dot-separated segments.
- `VerifyJWT` returns the correct user ID and role for a valid token.
- `VerifyJWT` returns an error for a token signed with a different secret.
- `VerifyJWT` returns an error for an expired token.

## Review questions

1. What are the three parts of a JWT and what does each contain?
2. Why is the `alg: none` vulnerability dangerous, and how does `golang-jwt/jwt` v5 prevent it?
3. What is the difference between HS256 and RS256? When would you choose one over the other?
4. Why should sensitive data (like passwords) never be stored in JWT claims?
5. How does JWT revocation work without a server-side session store?

## NEXT UP

JWT risks — common vulnerabilities including algorithm confusion, key theft, and the difficulty of token revocation.
