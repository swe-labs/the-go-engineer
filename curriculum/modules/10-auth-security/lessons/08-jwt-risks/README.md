# JWT risks

## Learning objective

Identify and mitigate common JWT vulnerabilities in Go services including the `alg=none` attack, key confusion, token theft, replay attacks, and the inherent difficulty of token revocation.

## Why this matters

JWTs are widely adopted but frequently misimplemented. The "none algorithm" vulnerability was present in every major JWT library in 2015, including `golang-jwt/jwt` (then `dgrijalva/jwt-go`). Key confusion attacks (forcing an RS256 public key to verify an HS256 token) have compromised real-world auth systems. Token theft via XSS or network interception is the most common attack on JWT-based systems because the stateless design makes revocation near-impossible without a denylist. Understanding these risks is essential before deploying JWTs in production — the convenience of stateless tokens comes with a distinct security tradeoff.

## Mental model

A JWT is a signed document. The signature proves the document has not been tampered with since signing. But the document itself is readable by anyone who possesses it. If an attacker steals the document, they can present it as their own until it expires. The issuer cannot "un-sign" a document — they can only publish a list of cancelled document IDs (denylist), which requires a database lookup, defeating the purpose of stateless tokens.

## Core idea

The major JWT risks fall into five categories:

| Risk | Description | Impact |
|---|---|---|
| Algorithm confusion | Attacker changes `alg` from `RS256` to `HS256` and signs with the public key | Token forgery |
| `alg: none` | Attacker sets `alg` to `none` and removes the signature | Complete authentication bypass |
| Token theft | Session hijacking via XSS, network interception, or device theft | Account takeover |
| Replay attacks | Same token used multiple times in different contexts | Unauthorized actions |
| Revocation difficulty | No built-in way to invalidate a token before expiry | Stolen tokens remain usable |
| Weak signing key | Poor entropy or leaked private key | All tokens can be forged |

## Under the hood

**Algorithm confusion attack**: The JWT header contains the algorithm. If the server's verification code accepts the algorithm from the header and selects the verification key based on it, an attacker can:

1. Obtain the server's public RSA key (often exposed via JWKS endpoint).
2. Create a JWT with `alg: HS256` (HMAC).
3. Sign the header+payload using the public key as the HMAC secret.
4. The server, seeing `HS256` in the header, uses the "secret" to verify — in many implementations, this is the same public key the server loaded for RSA verification.

Mitigation: always validate that the algorithm matches the expected one before verifying.

**`alg: none` attack**: The `none` algorithm means "no signature." Vulnerable libraries accept tokens without signatures when the algorithm is `none`. `golang-jwt/jwt` v5 explicitly rejects `none` by default, but older versions (v3 and earlier) were vulnerable.

**Token theft**: Since JWTs are sent with every request (typically as a Bearer header), any XSS vulnerability or insecure channel can leak the token. Unlike session cookies, JWTs cannot be marked `HttpOnly` because they are not cookies — they are in JavaScript-accessible storage or headers.

## How Go uses it

`golang-jwt/jwt` v5 includes built-in protections against many of these attacks:

- **`alg: none`**: the library rejects tokens with `alg: none` unless explicitly configured with `jwt.WithValidMethods`.
- **Algorithm validation**: the key function (`KeyFunc`) receives the parsed token and can inspect `token.Method` to verify the algorithm.
- **`jwt.ParseWithClaims`**: must be used instead of `jwt.Parse` to ensure claims are type-safe.

The Go security team recommends:

1. Always use `jwt.ParseWithClaims` with a concrete claims struct.
2. Validate `token.Method` against a whitelist in the key function.
3. Never use `jwt.Parse` without a key function (passing `nil` was the root cause of many `alg: none` vulnerabilities).
4. Use asymmetric algorithms (RS256/ES256) for multi-service architectures.

## Go example

```go
package main

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	rsaPrivateKey *rsa.PrivateKey
	rsaPublicKey  *rsa.PublicKey
	hmacSecret    = []byte("my-super-secret-key-change-in-production")
)

func init() {
	var err error
	rsaPrivateKey, err = rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Fatalf("key generation: %v", err)
	}
	rsaPublicKey = &rsaPrivateKey.PublicKey
}

// DenyList stores revoked token IDs.
type DenyList struct {
	mu sync.RWMutex
	m  map[string]time.Time
}

func NewDenyList() *DenyList {
	return &DenyList{m: make(map[string]time.Time)}
}

func (dl *DenyList) Revoke(jti string, expiry time.Time) {
	dl.mu.Lock()
	dl.m[jti] = expiry
	dl.mu.Unlock()
}

func (dl *DenyList) IsRevoked(jti string) bool {
	dl.mu.RLock()
	exp, ok := dl.m[jti]
	dl.mu.RUnlock()
	if !ok {
		return false
	}
	if time.Now().After(exp) {
		dl.mu.Lock()
		delete(dl.m, jti)
		dl.mu.Unlock()
		return false
	}
	return true
}

// SafeClaims includes a JTI for revocation support.
type SafeClaims struct {
	jwt.RegisteredClaims
	Role  string `json:"role"`
	Scope string `json:"scope"`
}

// IssueToken creates a JWT with RS256.
func IssueToken(userID, role string) (string, error) {
	now := time.Now()
	claims := SafeClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID,
			ExpiresAt: jwt.NewNumericDate(now.Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(now),
			ID:        fmt.Sprintf("%s-%d", userID, now.UnixNano()),
		},
		Role:  role,
		Scope: "api:read",
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(rsaPrivateKey)
}

// VerifyTokenSafe validates a JWT with algorithm whitelisting and denylist check.
func VerifyTokenSafe(tokenStr string, denyList *DenyList) (*SafeClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &SafeClaims{}, func(token *jwt.Token) (interface{}, error) {
		// SECURITY: algorithm validation prevents key confusion.
		if token.Method != jwt.SigningMethodRS256 {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return rsaPublicKey, nil
	})
	if err != nil {
		return nil, fmt.Errorf("parse failed: %w", err)
	}
	claims, ok := token.Claims.(*SafeClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}
	// Denylist check: has this token been revoked?
	if denyList.IsRevoked(claims.ID) {
		return nil, errors.New("token revoked")
	}
	return claims, nil
}

// VULNERABLE verification — for demonstration only.
func VerifyTokenVulnerable(tokenStr string) (string, error) {
	// VULNERABILITY: passes nil as key function — algorithm is not validated.
	token, err := jwt.Parse(tokenStr, nil)
	if err != nil {
		return "", err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return "", errors.New("invalid token")
	}
	return claims["sub"].(string), nil
}

// DemonstrateAlgorithmConfusion shows how key confusion works.
func DemonstrateAlgorithmConfusion() {
	// An attacker obtains the public key (usually via JWKS endpoint).
	attackerPublicKeyPEM := rsaPublicKey

	// Attacker creates a forged JWT with alg: HS256, signed with public key as secret.
	forged := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": "admin",
		"exp": time.Now().Add(1 * time.Hour).Unix(),
		"iat": time.Now().Unix(),
	})
	forgedStr, _ := forged.SignedString(attackerPublicKeyPEM)

	fmt.Println("Forged token (HS256 signed with public key):")
	fmt.Println(forgedStr[:80], "...")

	// Vulnerable server accepts it.
	sub, err := VerifyTokenVulnerable(forgedStr)
	if err != nil {
		fmt.Println("Vulnerable server rejected (good):", err)
	} else {
		fmt.Println("VULNERABLE: Server accepted forged token as user:", sub)
	}

	// Safe server rejects it.
	_, err = VerifyTokenSafe(forgedStr, NewDenyList())
	if err != nil {
		fmt.Println("Safe server rejected:", err)
	}
}

// DemonstrateNoneAlg shows the alg: none attack.
func DemonstrateNoneAlg() {
	// Attacker crafts a token with alg: none and no signature.
	header := base64.RawURLEncoding.EncodeToString([]byte(`{"alg":"none","typ":"JWT"}`))
	payload := base64.RawURLEncoding.EncodeToString([]byte(`{"sub":"admin","role":"admin","exp":9999999999}`))
	noneToken := fmt.Sprintf("%s.%s.", header, payload)

	fmt.Println("\nNone-alg token:")
	fmt.Println(noneToken[:80], "...")

	sub, err := VerifyTokenVulnerable(noneToken)
	if err != nil {
		fmt.Println("Safe: Server rejected none-alg token:", err)
	} else {
		fmt.Println("VULNERABLE: Server accepted none-alg token as user:", sub)
	}
}

func main() {
	fmt.Println("=== JWT vulnerability demonstrations ===\n")

	// 1. Issue a legitimate token.
	token, _ := IssueToken("alice", "user")
	fmt.Println("Legitimate token for alice:")
	fmt.Println(token[:80], "...\n")

	// 2. Normal verification.
	denyList := NewDenyList()
	claims, err := VerifyTokenSafe(token, denyList)
	if err != nil {
		log.Fatal("Unexpected:", err)
	}
	fmt.Printf("Verified: user=%s role=%s jti=%s\n\n", claims.Subject, claims.Role, claims.ID)

	// 3. Revoke and verify.
	denyList.Revoke(claims.ID, claims.ExpiresAt.Time)
	_, err = VerifyTokenSafe(token, denyList)
	if err != nil {
		fmt.Println("Revoked token correctly rejected:", err, "\n")
	}

	// 4. Algorithm confusion demo.
	DemonstrateAlgorithmConfusion()

	// 5. None alg demo.
	DemonstrateNoneAlg()
}
```

## Step-by-step execution

For the algorithm confusion demonstration:

1. The attacker obtains the server's RSA public key (often from a public JWKS endpoint).
2. The attacker creates a JWT with header `{"alg":"HS256"}` and payload `{"sub":"admin","exp":...}`.
3. The attacker signs the token using HMAC-SHA256 with the public key bytes as the secret.
4. The vulnerable server parses the token, sees `alg: HS256`, reads the "secret" from its key store — but the key store contains the public key, which is the same bytes the attacker used.
5. The server computes the HMAC signature using the public key as the secret and gets a match.
6. Server accepts the forged admin token.

The safe server prevents this by validating `token.Method` against `jwt.SigningMethodRS256` before selecting the key.

## Common mistakes

- Not validating the algorithm header — accepting `alg: none` or `alg: HS256` when expecting `RS256` bypasses signature verification.
- Using HS256 in multi-service architectures — every service must hold the shared secret, increasing the attack surface.
- Failing to validate the `exp` claim — tokens without expiration remain valid indefinitely if leaked.
- Storing sensitive data in JWT claims — the payload is base64-encoded, not encrypted. Anyone with the token can read claims.
- Not implementing a denylist for revocation — stolen tokens remain usable until expiry.
- Using short expiry times as the sole revocation strategy — force-logging-out all users every 5 minutes is poor UX.

## Debugging walkthrough

Consider this Go code that verifies JWTs:

```go
func authMiddleware(publicKey *rsa.PublicKey) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            tokenStr := extractBearer(r)
            token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
                return publicKey, nil
            })
            if err != nil || !token.Valid {
                http.Error(w, "unauthorized", http.StatusUnauthorized)
                return
            }
            next.ServeHTTP(w, r)
        })
    }
}
```

**Symptom**: An attacker authenticates as `admin` with a token that was never issued by the server.

**Investigation**: The `jwt.Parse` function is used instead of `jwt.ParseWithClaims`. The key function returns the public key without checking the algorithm. An attacker sends a token with `alg: none` — the library does not call the key function and accepts the token as valid.

**Root cause**: Using `jwt.Parse` (not `ParseWithClaims`) without algorithm validation. The key function is never called for `alg: none` tokens.

**Fix**:

```go
func authMiddleware(publicKey *rsa.PublicKey) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            tokenStr := extractBearer(r)
            token, err := jwt.ParseWithClaims(tokenStr, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
                if token.Method != jwt.SigningMethodRS256 {
                    return nil, fmt.Errorf("unexpected alg: %v", token.Header["alg"])
                }
                return publicKey, nil
            })
            if err != nil {
                http.Error(w, "unauthorized", http.StatusUnauthorized)
                return
            }
            claims, ok := token.Claims.(*CustomClaims)
            if !ok || !token.Valid {
                http.Error(w, "forbidden", http.StatusForbidden)
                return
            }
            // ... proceed with claims
        })
    }
}
```

## Production notes

- Always validate the algorithm in the key function. Whichever the expected algorithms, return an error for anything else.
- Use asymmetric signing (RS256/ES256) so the signing key is not distributed to every service.
- Implement a token denylist for high-value actions (password change, privilege escalation). Use Redis with the token `jti` as key and the expiry as TTL.
- Keep token lifetimes short (15 minutes for access tokens, 7 days for refresh tokens).
- Never expose the private key. Use a secrets manager and rotate keys quarterly.
- Monitor for tokens with unexpected algorithms in access logs — this signals an attack probe.

## Performance implications

- Denylist lookup adds 1-5ms per request (Redis). This partially defeats statelessness but is necessary for revocation.
- Algorithm validation is free — it is a struct comparison.
- Token parsing and signature verification is the same cost as any JWT implementation (~0.5ms for RS256).
- Storing `jti` in Redis with TTL matching the token expiry auto-cleans the denylist.

## Practice task

Write Go functions `VerifyTokenWithAlgCheck(tokenStr string, key interface{}) (jwt.MapClaims, error)` and `DemonstrateForgery(publicKey *rsa.PublicKey) string` where:
- `VerifyTokenWithAlgCheck` rejects any token whose algorithm is not `RS256`.
- `DemonstrateForgery` creates an HS256 token signed with the RSA public key (simulating attacker behavior).
- Then write a `main()` that shows the forged token is accepted by a vulnerable verifier but rejected by `VerifyTokenWithAlgCheck`.

## Tests / verification

```bash
go run ./curriculum/modules/10-auth-security/lessons/08-jwt-risks
go test ./curriculum/modules/10-auth-security/lessons/08-jwt-risks
```

The test file `main_test.go` contains table-driven tests that verify:
- `VerifyTokenWithAlgCheck` returns claims for a legitimate RS256 token.
- `VerifyTokenWithAlgCheck` rejects an HS256 token.
- `VerifyTokenWithAlgCheck` rejects a token with `alg: none`.
- `VerifyTokenWithAlgCheck` rejects a token signed with the wrong key.

## Review questions

1. Describe the JWT algorithm confusion attack. What makes a Go server vulnerable to it?
2. Why is the `alg: none` attack still a risk even though `golang-jwt/jwt` v5 rejects it by default?
3. What is a JTI claim and how does it enable token revocation?
4. If a JWT access token has a 15-minute expiry, what is the maximum window of vulnerability if the token is stolen?
5. Why can't JWTs have `HttpOnly` protection like session cookies, and what alternative mitigations exist?

## NEXT UP

OAuth2 and OIDC overview — delegated authorization and federated identity protocols that build on JWT concepts.
