# CSRF

## Learning objective

Identify Cross-Site Request Forgery (CSRF) attacks and implement server-side protections using anti-CSRF tokens, SameSite cookies, origin/referer header validation, and middleware patterns in Go.

## Why this matters

CSRF attacks exploit the browser's automatic cookie inclusion to forge authenticated requests. A logged-in user visits an attacker's site, which submits a form to your bank's transfer endpoint. The browser automatically includes the session cookie, and the server processes the request as if the user intended it. CSRF vulnerabilities have been used to change email addresses, transfer funds, and modify account settings on major platforms. Go applications using cookie-based authentication must implement CSRF protection unless all state-changing endpoints use token-based auth (e.g., Bearer tokens in Authorization headers).

## Mental model

CSRF is a "confused deputy" attack. The browser is the deputy: it has authority (the user's session cookie) and is tricked by the attacker into using that authority against the user's interest.

Anti-CSRF tokens are a secret handshake. The server embeds a unique, unpredictable token in every form or response. When the client submits a state-changing request, it must include this token. The attacker's site cannot read the token (same-origin policy blocks cross-origin reads), so it cannot include it in a forged request.

SameSite cookies are a complementary defense. A cookie with `SameSite=Strict` is not sent with cross-origin requests, preventing CSRF entirely. `SameSite=Lax` allows top-level navigations (like clicking a link) but blocks POST forms from other sites.

## Core idea

A CSRF attack has three requirements:
1. The victim is authenticated on the target site (has a valid session cookie).
2. The target site has a state-changing endpoint (POST /transfer, POST /change-email).
3. The attacker can craft a valid request to that endpoint.

The defense is to ensure that each state-changing request contains a value that the attacker cannot forge or obtain.

| Defense | How it works | Strength |
|---|---|---|
| Synchronizer token pattern | Server generates token, embedded in form, validated on submission | Strong |
| Double-submit cookie | Token set as cookie + custom header, validated server-side | Strong |
| SameSite cookie attribute | Browser prevents cookie from being sent cross-origin | Strong (modern browsers) |
| Origin/Referer header validation | Server checks that the request originated from its own origin | Moderate |
| Custom request header | API requires a non-standard header (e.g., X-Requested-By) | Weak alone |

## Under the hood

The synchronizer token pattern works as follows:

1. Server generates a cryptographically random token (at least 128 bits) and stores it associated with the user's session.
2. Server renders the token into forms as a hidden field, or returns it in a response header.
3. Client submits the token with each state-changing request (usually in a header or form field).
4. Server looks up the expected token for this session and compares. If they match, the request is accepted. The token is single-use (consumed after validation).

For the double-submit cookie pattern:
1. Server sets a CSRF token as a cookie (not HttpOnly, since JavaScript needs to read it).
2. JavaScript reads the cookie and includes the same value in a custom header (e.g., `X-CSRF-Token`).
3. Server checks that the header value matches the cookie value.

SameSite cookies work at the browser level:
- `SameSite=Strict`: Cookie is never sent for cross-origin requests.
- `SameSite=Lax`: Cookie is sent for top-level navigation GET requests but not for POST.
- `SameSite=None; Secure`: Cookie is sent for all cross-origin requests (requires Secure flag).

## How Go uses it

Go standard library does not include built-in CSRF middleware. Popular frameworks provide it:

- **Chi**: `chi/middleware` provides `RealIP`, `Logger`, etc. CSRF via `github.com/gorilla/csrf`.
- **Gin**: Built-in CSRF middleware or third-party `github.com/utrack/gin-csrf`.
- **Echo**: `github.com/labstack/echo/v4/middleware` includes `CSRFWithConfig`.
- **Gorilla CSRF**: `github.com/gorilla/csrf` is a standalone package compatible with `net/http`.

For custom implementations, a CSRF protector can be a simple struct with a sync-protected map of tokens.

## Go example

```go
package main

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sync"
)

type CSRFProtector struct {
	mu     sync.RWMutex
	tokens map[string]string
	key    []byte
}

func NewCSRFProtector() *CSRFProtector {
	key := make([]byte, 32)
	rand.Read(key)
	return &CSRFProtector{
		tokens: make(map[string]string),
		key:    key,
	}
}

func (c *CSRFProtector) GenerateToken(sessionID string) string {
	c.mu.Lock()
	defer c.mu.Unlock()
	token := make([]byte, 32)
	rand.Read(token)
	encoded := hex.EncodeToString(token)
	hash := sha256.Sum256([]byte(encoded + c.key))
	c.tokens[string(hash[:])] = sessionID
	return encoded
}

func (c *CSRFProtector) ValidateToken(token, sessionID string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	hash := sha256.Sum256([]byte(token + c.key))
	key := string(hash[:])
	storedSession, ok := c.tokens[key]
	if !ok {
		return false
	}
	delete(c.tokens, key)
	return storedSession == sessionID
}

func main() {
	p := NewCSRFProtector()
	session := "session-abc-123"
	token := p.GenerateToken(session)
	fmt.Println("Generated token:", token)

	valid := p.ValidateToken(token, session)
	fmt.Println("First use valid:", valid)

	valid = p.ValidateToken(token, session)
	fmt.Println("Second use valid (replay):", valid)
}
```

## Step-by-step execution

For a CSRF attack on a bank transfer endpoint:

1. User logs in to `bank.com`. Server sets a session cookie.
2. User visits `evil.com` in another tab.
3. `evil.com` has an auto-submitting form: `<form action="https://bank.com/transfer" method="POST"><input type="hidden" name="amount" value="10000"><input type="hidden" name="to" value="attacker"></form><script>document.forms[0].submit()</script>`.
4. Browser submits the form to `bank.com`. The session cookie is automatically included.
5. `bank.com` processes the transfer. The request is authenticated (cookie is valid) but unauthorized (user did not intend to transfer money).

With CSRF protection:

1. `bank.com` generates a CSRF token on login and stores it server-side.
2. The transfer form includes the token: `<input type="hidden" name="csrf_token" value="random-hex">`.
3. `evil.com` cannot read the token because it is on a different origin.
4. The forged form submission lacks the CSRF token.
5. Server rejects the request with HTTP 403.

## Common mistakes

- Mistake: Building APIs without CSRF protection, assuming JSON requests are not vulnerable.
  - Why it happens: Developers think JSON content type prevents form submission. But `XMLHttpRequest` and `fetch` with credentials: include can send JSON cross-origin.
  - Fix: CORS must be configured to reject cross-origin requests with credentials. Always protect state-changing endpoints.

- Mistake: Using the wrong token pattern: storing tokens in cookies without SameSite attribute.
  - Why it happens: Developers implement double-submit cookie but the attacker's site can also set cookies. The attacker could set a known token as a cookie and include it in a header.
  - Fix: Use SameSite=Strict or SameSite=Lax for session cookies. Do not rely solely on double-submit cookies.

- Mistake: Implementing CSRF tokens manually without one-time use (token replay).
  - Why it happens: Developers validate the token but do not consume it. An attacker who intercepts a token can reuse it.
  - Fix: Consume (delete) the token after successful validation. Generate a new token for the next request.

- Mistake: Disabling CSRF protection for authenticated endpoints by mistake.
  - Why it happens: Developers skip CSRF on GET endpoints (correct) but also on some POST endpoints that "don't need it."
  - Fix: CSRF protection applies to all state-changing methods: POST, PUT, PATCH, DELETE.

- Mistake: Exposing CSRF tokens in URLs via query parameters.
  - Why it happens: Developers include the token in a GET request for convenience. The token is leaked in referer headers and server logs.
  - Fix: CSRF tokens should never appear in URLs. Use headers or form bodies.

## Debugging walkthrough

Consider this broken Go middleware:

```go
func csrfMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			next.ServeHTTP(w, r)
			return
		}
		token := r.Header.Get("X-CSRF-Token")
		cookie, err := r.Cookie("csrf_token")
		if err != nil || cookie.Value != token {
			http.Error(w, "CSRF validation failed", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}
```

Symptom: CSRF validation passes for legitimate requests but also passes for forged requests in some browsers.

Investigation: The code compares a custom header value with a cookie value. But the cookie `csrf_token` is not set with the `HttpOnly` flag (good, JS needs to read it) and most importantly, the code does not check that the request origin matches the site's origin.

Root cause: The double-submit cookie pattern without origin validation. While `evil.com` cannot read `bank.com`'s cookies, it can set its own cookies for `bank.com` (if the cookie domain is misconfigured). More importantly, this pattern does not validate the origin of the request.

Fix: Add Origin/Referer header validation, use SameSite cookies, and consider using the synchronizer token pattern instead:

```go
func csrfMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" || r.Method == "HEAD" || r.Method == "OPTIONS" {
			next.ServeHTTP(w, r)
			return
		}
		origin := r.Header.Get("Origin")
		if origin != "" && !strings.HasPrefix(origin, "https://bank.com") {
			http.Error(w, "CSRF: invalid origin", http.StatusForbidden)
			return
		}
		// ... token validation
		next.ServeHTTP(w, r)
	})
}
```

## Production notes

- Use a proven library (gorilla/csrf, chi's middleware, gin's CSRF) rather than writing your own. Cryptography and token storage are easy to get wrong.
- Set `SameSite=Strict` on session cookies by default. This alone prevents CSRF on modern browsers.
- Use origin/referer header validation as a defense-in-depth measure. Never trust the referer alone (it can be missing or stripped), but use it when present.
- CSRF tokens should be cryptographically random (at least 128 bits) generated by `crypto/rand`. Never use `math/rand` for security-sensitive values.
- For APIs that use Bearer token authentication (Authorization header), CSRF protection is not needed because the browser does not automatically include Authorization headers.
- Log CSRF validation failures at WARN level. Consistent failures from a single IP may indicate a scanning attack.

## Performance implications

- CSRF token generation uses `crypto/rand`, which is slightly slower than `math/rand` but still negligible (microseconds per token).
- Token storage (in-memory map or session store) is fast O(1) lookup.
- The HTTP overhead of an extra header or cookie is negligible.
- The main performance cost is the complexity of cookie parsing and the extra HTTP round trip if you fail early without a token.

## Practice task

Write a `CSRFMiddleware` function that implements the synchronizer token pattern. It should:

1. Generate a CSRF token on GET requests (returned in `X-CSRF-Token` response header).
2. Store the token in an in-memory store, keyed by a session ID from the `X-Session-ID` header.
3. Validate the token on POST/PUT/DELETE requests from the `X-CSRF-Token` request header.
4. Consume the token after successful validation (one-time use).
5. Return HTTP 403 on validation failure.

Write a `main()` that creates a test server with this middleware, makes a GET request to get a token, then a POST request with the valid token, and a POST request with an invalid token. Print the responses.

## Tests / verification

```bash
go run ./curriculum/modules/10-auth-security/lessons/16-csrf
go test ./curriculum/modules/10-auth-security/lessons/16-csrf
```

The existing tests verify token generation produces unique values, validation succeeds for correct session/token pairs, tokens are single-use (replay prevention), and tampered tokens are rejected.

## Review questions

1. How does the synchronizer token pattern prevent CSRF when SameSite cookies would also work?
2. Why can't the attacker's website read the CSRF token from the victim's session to include it in a forged request?
3. What is the difference between `SameSite=Strict` and `SameSite=Lax`? When would you use each?
4. Why is CSRF protection unnecessary for APIs that use Bearer token authentication in the Authorization header?
5. What is token replay and how does consuming tokens after validation prevent it?

## NEXT UP

CORS -- Cross-Origin Resource Sharing, how browser-enforced access control works, and how to configure CORS headers in Go APIs.
