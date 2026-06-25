# OAuth2 and OIDC overview

## Learning objective

Explain the four OAuth2 grant flows and understand how OIDC adds an identity layer on top. Implement an OAuth2 authorization code client in Go using `golang.org/x/oauth2` and validate an OIDC ID token.

## Why this matters

OAuth2 is the industry standard for delegated authorization — it is how Google, GitHub, Facebook, and every major platform let third-party applications access user data without sharing passwords. OIDC extends OAuth2 for authentication, enabling "Sign in with Google" and similar federated identity patterns. As a Go engineer, you will integrate with OAuth2/OIDC providers for authentication (login with GitHub), API access (calling Google APIs with user consent), and service-to-service authorization (machine-to-machine flows).

## Mental model

OAuth2 is a valet key system for APIs. The user (resource owner) gives a parking attendant (client application) a limited valet key (access token) instead of their car keys (password). The valet key only works for specific cars (scopes) and for a limited time (token expiry). The user gets their car keys back after the valet returns. OIDC adds an ID card (ID token) to the valet key: the parking garage can verify who the user is, not just that they have a key.

## Core idea

OAuth2 defines four roles:

| Role | Example | What they do |
|---|---|---|
| Resource owner | User | Owns the data |
| Client | Go web app | Wants to access the data |
| Authorization server | Google Auth | Issues tokens after user consent |
| Resource server | Google APIs | Accepts tokens to serve data |

The four grant flows (authorization grant types):

| Flow | Who uses it | Security |
|---|---|---|
| Authorization code | Web apps with backend | Most secure — token never reaches browser |
| Authorization code + PKCE | SPAs, mobile apps | Secure without client secret |
| Client credentials | Service-to-service | No user involved, machine-to-machine |
| Authorization code + PKCE is the modern replacement for the implicit flow (which is deprecated). |

OIDC extends OAuth2 with:

- The `openid` scope requests an ID token in addition to the access token.
- The ID token is a JWT containing the user's identity (sub, name, email, picture).
- The `/userinfo` endpoint returns additional claims about the user.
- OIDC discovery (`/.well-known/openid-configuration`) returns the provider's metadata.

## Under the hood

The authorization code flow with PKCE:

1. Client generates a `code_verifier` (random string) and `code_challenge` = SHA256(code_verifier).
2. Client redirects user to authorization server with `client_id`, `redirect_uri`, `scope`, `code_challenge`.
3. User authenticates and consents.
4. Authorization server redirects back to `redirect_uri` with a one-time `code`.
5. Client exchanges `code` + `code_verifier` for an access token and ID token.
6. Verification: server computes SHA256(code_verifier) and compares with `code_challenge`.

The `code` is valid for one use and short-lived (seconds). The `code_verifier` ensures that even if the `code` is intercepted, the attacker cannot exchange it without the `code_verifier`.

## How Go uses it

The `golang.org/x/oauth2` package provides the client side of OAuth2:

```go
import "golang.org/x/oauth2"

config := &oauth2.Config{
    ClientID:     "your-client-id",
    ClientSecret: "your-client-secret", // not needed with PKCE
    RedirectURL:  "http://localhost:8080/callback",
    Scopes:       []string{"openid", "profile", "email"},
    Endpoint:     google.Endpoint,
}
```

For OIDC, the `coreos/go-oidc/v3` package validates ID tokens and discovers provider configuration:

```go
provider, err := oidc.NewProvider(ctx, "https://accounts.google.com")
verifier := provider.Verifier(&oidc.Config{ClientID: "your-client-id"})
idToken, err := verifier.Verify(ctx, rawIDToken)
```

## Go example

```go
package main

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"golang.org/x/oauth2"
)

// PKCEParams holds the PKCE challenge and verifier values.
type PKCEParams struct {
	Verifier        string
	Challenge       string
	ChallengeMethod string
}

// GeneratePKCE creates a code verifier and challenge per RFC 7636.
func GeneratePKCE() (*PKCEParams, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return nil, fmt.Errorf("random: %w", err)
	}
	verifier := base64.RawURLEncoding.EncodeToString(b)
	hash := sha256.Sum256([]byte(verifier))
	challenge := base64.RawURLEncoding.EncodeToString(hash[:])
	return &PKCEParams{
		Verifier:        verifier,
		Challenge:       challenge,
		ChallengeMethod: "S256",
	}, nil
}

// OAuthClient simulates a minimal OAuth2 authorization server for demonstration.
// In production, you connect to Google, GitHub, Auth0, etc.
type OAuthClient struct {
	clientID     string
	redirectURI string
	authzServer string
	httpClient  *http.Client
	scopes      []string
}

// NewOAuthClient creates a client configured for a provider.
func NewOAuthClient(clientID, redirectURI, authzServer string, scopes []string) *OAuthClient {
	return &OAuthClient{
		clientID:     clientID,
		redirectURI: redirectURI,
		authzServer: authzServer,
		httpClient:  &http.Client{Timeout: 10 * time.Second},
		scopes:      scopes,
	}
}

// AuthURL returns the authorization URL with PKCE.
func (c *OAuthClient) AuthURL(state string, pkce *PKCEParams) string {
	params := fmt.Sprintf(
		"response_type=code&client_id=%s&redirect_uri=%s&scope=%s&state=%s&code_challenge=%s&code_challenge_method=%s",
		c.clientID, c.redirectURI, strings.Join(c.scopes, " "), state, pkce.Challenge, pkce.ChallengeMethod,
	)
	return c.authzServer + "/authorize?" + params
}

// ExchangeCode sends the authorization code and PKCE verifier to the token endpoint.
func (c *OAuthClient) ExchangeCode(ctx context.Context, code, verifier string) (*oauth2.Token, error) {
	body := fmt.Sprintf(
		"grant_type=authorization_code&code=%s&redirect_uri=%s&client_id=%s&code_verifier=%s",
		code, c.redirectURI, c.clientID, verifier,
	)
	req, err := http.NewRequestWithContext(ctx, "POST", c.authzServer+"/token", strings.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var token oauth2.Token
	if err := json.NewDecoder(resp.Body).Decode(&token); err != nil {
		return nil, err
	}
	return &token, nil
}

// IDToken represents a parsed OIDC ID token for demonstration.
type IDToken struct {
	Subject   string `json:"sub"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	Issuer    string `json:"iss"`
	ExpiresAt int64  `json:"exp"`
}

// ParseIDToken parses a raw ID token string (JWT) without verification.
// In production, use github.com/coreos/go-oidc/v3/oidc for verification.
func ParseIDToken(rawIDToken string) (*IDToken, error) {
	parts := strings.Split(rawIDToken, ".")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid JWT: expected 3 segments, got %d", len(parts))
	}
	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, fmt.Errorf("base64 decode: %w", err)
	}
	var token IDToken
	if err := json.Unmarshal(payload, &token); err != nil {
		return nil, fmt.Errorf("json decode: %w", err)
	}
	return &token, nil
}

// SimulatedTokenEndpoint handles token exchange in our mock server.
func SimulatedTokenEndpoint() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		code := r.FormValue("code")
		verifier := r.FormValue("code_verifier")

		// Verify PKCE: in production, server stores challenge from auth request.
		if code == "" || verifier == "" {
			http.Error(w, "missing code or verifier", http.StatusBadRequest)
			return
		}

		// Return a mock token response.
		tokenResp := struct {
			AccessToken  string `json:"access_token"`
			TokenType    string `json:"token_type"`
			ExpiresIn    int    `json:"expires_in"`
			IDToken      string `json:"id_token,omitempty"`
		}{
			AccessToken: "mock-access-token-" + code[:8],
			TokenType:   "Bearer",
			ExpiresIn:   3600,
			IDToken:     "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwiZW1haWwiOiJhbGljZUBleGFtcGxlLmNvbSIsIm5hbWUiOiJBbGljZSIsImlzcyI6Imh0dHBzOi8vYXV0aC5leGFtcGxlLmNvbSIsImV4cCI6OTk5OTk5OTk5OX0.signature",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(tokenResp)
	}
}

func main() {
	// Start a mock authorization and token server.
	http.HandleFunc("/authorize", func(w http.ResponseWriter, r *http.Request) {
		// Redirect back to the callback with a code.
		code := "auth-code-abc123"
		state := r.URL.Query().Get("state")
		redirectURI := r.URL.Query().Get("redirect_uri")
		http.Redirect(w, r, fmt.Sprintf("%s?code=%s&state=%s", redirectURI, code, state), http.StatusFound)
	})
	http.HandleFunc("/token", SimulatedTokenEndpoint())

	go func() {
		log.Println("OAuth2 mock server on :9090")
		log.Fatal(http.ListenAndServe(":9090", nil))
	}()
	time.Sleep(100 * time.Millisecond) // let server start

	// Client side: initiate OAuth2 flow.
	fmt.Println("=== OAuth2 Authorization Code Flow with PKCE ===\n")

	client := NewOAuthClient(
		"demo-client",
		"http://localhost:8080/callback",
		"http://localhost:9090",
		[]string{"openid", "profile", "email"},
	)

	// Step 1: Generate PKCE challenge.
	pkce, _ := GeneratePKCE()
	state := "random-state-123"

	// Step 2: Build authorization URL.
	authURL := client.AuthURL(state, pkce)
	fmt.Println("1. Redirect user to:", authURL[:100], "...\n")

	// Step 3: Simulate callback (user is redirected back with code).
	code := "auth-code-abc123"

	// Step 4: Exchange code for tokens.
	ctx := context.Background()
	token, err := client.ExchangeCode(ctx, code, pkce.Verifier)
	if err != nil {
		log.Fatal("Token exchange failed:", err)
	}
	fmt.Printf("2. Access token: %s...\n", token.AccessToken[:20])

	// Step 5: Use access token to call resource API.
	req, _ := http.NewRequestWithContext(ctx, "GET", "http://localhost:9090/userinfo", nil)
	req.Header.Set("Authorization", "Bearer "+token.AccessToken)
	// In production, this calls the provider's /userinfo endpoint.
	fmt.Println("3. Token type:", token.TokenType)
	fmt.Printf("4. Expires in: %.0f seconds\n\n", token.Expiry.Sub(time.Now()).Seconds())

	// Parse ID token (in production, verify with go-oidc).
	// For this demo we manually set an ID token.
	mockIDToken := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwiZW1haWwiOiJhbGljZUBleGFtcGxlLmNvbSIsIm5hbWUiOiJBbGljZSIsImlzcyI6Imh0dHBzOi8vYXV0aC5leGFtcGxlLmNvbSIsImV4cCI6OTk5OTk5OTk5OX0"
	idToken, _ := ParseIDToken(mockIDToken)
	fmt.Printf("5. ID Token claims:\n   Subject: %s\n   Email: %s\n   Name: %s\n   Issuer: %s\n",
		idToken.Subject, idToken.Email, idToken.Name, idToken.Issuer)
}

// Ensure io is used (for compilation).
var _ = io.Discard
```

## Step-by-step execution

For a user clicking "Sign in with Google" on a Go web app:

1. Server generates PKCE params: 32 random bytes -> base64url verifier, SHA256 -> challenge.
2. Server generates a random `state` value to prevent CSRF on the callback.
3. Server redirects user to `accounts.google.com/o/oauth2/auth?...` with `code_challenge`.
4. User logs in to Google, consents to scopes (email, profile).
5. Google redirects back to `https://yourapp.com/callback?code=xxx&state=yyy`.
6. Server validates `state` matches the original (prevents CSRF on auth callback).
7. Server POSTs to Google's token endpoint with `code`, `code_verifier`, `client_id`, `redirect_uri`.
8. Google verifies the `code_verifier` (SHA256 matches the stored challenge). Returns access token + ID token.
9. Server verifies the ID token signature (using Google's JWKS endpoint).
10. Server reads user identity from ID token claims (`sub`, `email`, `name`).
11. Server creates its own session for the user. Access token is stored for API calls.

## Common mistakes

- Implementing OAuth2 from scratch instead of using `golang.org/x/oauth2` — the specification has many edge cases.
- Confusing OAuth2 (authorization delegation) with OIDC (authentication) — OAuth2 alone does not verify user identity.
- Not validating the ID token signature and claims in OIDC — accepting unverified tokens allows impersonation.
- Storing client secrets in client-side code — public clients (SPAs, mobile apps) cannot keep secrets.
- Not using PKCE for public clients — without PKCE, the authorization code can be intercepted and exchanged by an attacker.
- Not validating the `state` parameter — without state validation, the callback is vulnerable to CSRF.

## Debugging walkthrough

Consider this OAuth2 callback handler:

```go
func callbackHandler(w http.ResponseWriter, r *http.Request) {
    code := r.URL.Query().Get("code")
    token, err := exchangeCode(code) // No PKCE, no state check
    if err != nil {
        http.Error(w, "auth failed", http.StatusInternalServerError)
        return
    }
    fmt.Fprintf(w, "Logged in with token: %s", token)
}
```

**Symptom**: Users occasionally get logged in as other users. The attacker intercepts the callback.

**Investigation**: There is no `state` parameter validation. An attacker can craft a URL that triggers a callback with a code pre-generated by the attacker. Without PKCE, the attacker's code can be exchanged for a token.

**Root cause**: Missing state validation and PKCE. The `state` parameter prevents CSRF on the callback. PKCE prevents code interception.

**Fix**:

```go
func callbackHandler(w http.ResponseWriter, r *http.Request) {
    state := r.URL.Query().Get("state")
    if state != sessionState { // Compare with stored state
        http.Error(w, "state mismatch", http.StatusForbidden)
        return
    }
    code := r.URL.Query().Get("code")
    verifier := sessionPKCEVerifier // Retrieve stored verifier
    token, err := exchangeCodeWithPKCE(code, verifier)
    // ...
}
```

## Production notes

- Use `golang.org/x/oauth2` for the client side. Never implement the OAuth2 flows manually.
- Use `github.com/coreos/go-oidc/v3/oidc` for OIDC. It handles JWKS fetching, key rotation, and ID token verification.
- Always use PKCE, even for confidential clients. It adds no complexity and prevents code interception.
- Store the `state` and `code_verifier` in the user's session before redirecting to the authorization server.
- Validate the `iss` (issuer) claim in ID tokens to prevent token reuse across different providers.
- For production, register your redirect URI with the provider exactly. Mismatches are a common source of integration failures.

## Performance implications

- OAuth2 flow only occurs at login (infrequent). Performance is not a primary concern.
- Token exchange: one HTTP POST to the provider (typical latency: 100-500ms).
- OIDC ID token verification: one JWKS fetch (cache for 24 hours), then JWT verification (~0.5ms).
- Access token validation on resource APIs: JWT verification per request (~0.5ms for RS256). If the provider uses opaque tokens, each request requires a token introspection call to the provider (10-50ms).

## Practice task

Write Go functions `BuildAuthURL(clientID, redirectURI, state, challenge string) string` and `ExchangeCodeForToken(ctx context.Context, code, verifier, tokenEndpoint string) (string, error)` where:
- `BuildAuthURL` constructs a valid OAuth2 authorization URL with PKCE.
- `ExchangeCodeForToken` POSTs to the token endpoint with the code and verifier.
- Then write a `main()` that simulates the full flow (generate PKCE, build URL, simulate callback, exchange code) and prints the access token.

## Tests / verification

```bash
go run ./curriculum/modules/10-auth-security/lessons/09-oauth2-and-oidc-overview
go test ./curriculum/modules/10-auth-security/lessons/09-oauth2-and-oidc-overview
```

The test file `main_test.go` contains table-driven tests that verify:
- `BuildAuthURL` includes `code_challenge` and `code_challenge_method=S256`.
- `BuildAuthURL` includes `state` and `redirect_uri`.
- `ExchangeCodeForToken` returns a token for valid inputs.
- `ParseIDToken` extracts subject, email, and name from a valid JWT payload.

## Review questions

1. What is the difference between OAuth2 and OIDC? Which one provides user identity information?
2. In the authorization code flow with PKCE, what is the purpose of the `code_verifier`?
3. What attack does the `state` parameter in OAuth2 prevent?
4. Why is the implicit flow (returning the access token directly in the URL fragment) deprecated in favor of PKCE?
5. When validating an OIDC ID token in Go, what claims must be verified beyond the signature?

## NEXT UP

API keys — generation, hashing, rotation, and scoping for service-to-service and developer API authentication.
