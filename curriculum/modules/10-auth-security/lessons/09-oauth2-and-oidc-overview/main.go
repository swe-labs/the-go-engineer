package main

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

type PKCEParams struct {
	Verifier        string
	Challenge       string
	ChallengeMethod string
}

func GeneratePKCE() (*PKCEParams, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return nil, fmt.Errorf("random: %w", err)
	}
	verifier := base64.RawURLEncoding.EncodeToString(b)
	hash := sha256.Sum256([]byte(verifier))
	challenge := base64.RawURLEncoding.EncodeToString(hash[:])
	return &PKCEParams{
		Verifier: verifier, Challenge: challenge, ChallengeMethod: "S256",
	}, nil
}

type OAuthClient struct {
	clientID    string
	redirectURI string
	authzServer string
	httpClient  *http.Client
	scopes      []string
}

func NewOAuthClient(clientID, redirectURI, authzServer string, scopes []string) *OAuthClient {
	return &OAuthClient{
		clientID: clientID, redirectURI: redirectURI,
		authzServer: authzServer, httpClient: &http.Client{Timeout: 10 * time.Second},
		scopes: scopes,
	}
}

func (c *OAuthClient) AuthURL(state string, pkce *PKCEParams) string {
	return fmt.Sprintf(
		"%s/authorize?response_type=code&client_id=%s&redirect_uri=%s&scope=%s&state=%s&code_challenge=%s&code_challenge_method=%s",
		c.authzServer, c.clientID, c.redirectURI, strings.Join(c.scopes, " "),
		state, pkce.Challenge, pkce.ChallengeMethod,
	)
}

func (c *OAuthClient) ExchangeCode(ctx context.Context, code, verifier string) (string, error) {
	body := fmt.Sprintf(
		"grant_type=authorization_code&code=%s&redirect_uri=%s&client_id=%s&code_verifier=%s",
		code, c.redirectURI, c.clientID, verifier,
	)
	req, err := http.NewRequestWithContext(ctx, "POST", c.authzServer+"/token", strings.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	var result struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}
	return result.AccessToken, nil
}

type IDToken struct {
	Subject   string `json:"sub"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	Issuer    string `json:"iss"`
	ExpiresAt int64  `json:"exp"`
}

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

func BuildAuthURL(clientID, redirectURI, state, challenge string) string {
	return fmt.Sprintf(
		"https://auth.example.com/authorize?response_type=code&client_id=%s&redirect_uri=%s&state=%s&code_challenge=%s&code_challenge_method=S256",
		clientID, redirectURI, state, challenge,
	)
}

func ExchangeCodeForToken(ctx context.Context, code, verifier, tokenEndpoint string) (string, error) {
	body := fmt.Sprintf(
		"grant_type=authorization_code&code=%s&code_verifier=%s", code, verifier,
	)
	req, err := http.NewRequestWithContext(ctx, "POST", tokenEndpoint, strings.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	var result struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}
	return result.AccessToken, nil
}

func main() {
	http.HandleFunc("/authorize", func(w http.ResponseWriter, r *http.Request) {
		code := "auth-code-abc123"
		state := r.URL.Query().Get("state")
		redirectURI := r.URL.Query().Get("redirect_uri")
		http.Redirect(w, r, fmt.Sprintf("%s?code=%s&state=%s", redirectURI, code, state), http.StatusFound)
	})
	http.HandleFunc("/token", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"access_token": "mock-access-token-abc123",
			"token_type":   "Bearer",
			"expires_in":   3600,
		})
	})

	go func() {
		fmt.Println("OAuth2 mock server on :9090")
		log.Fatal(http.ListenAndServe(":9090", nil))
	}()
	time.Sleep(100 * time.Millisecond)

	fmt.Println("=== OAuth2 Authorization Code Flow with PKCE ===")
	fmt.Println()

	client := NewOAuthClient("demo-client", "http://localhost:8080/callback", "http://localhost:9090", []string{"openid", "profile"})
	pkce, _ := GeneratePKCE()
	state := "random-state-123"

	authURL := client.AuthURL(state, pkce)
	fmt.Println("1. Authorization URL:", authURL[:80], "...")
	fmt.Println()

	code := "auth-code-abc123"
	ctx := context.Background()
	accessToken, err := client.ExchangeCode(ctx, code, pkce.Verifier)
	if err != nil {
		log.Fatal("Token exchange failed:", err)
	}
	fmt.Printf("2. Access token: %s\n", accessToken)

	builtURL := BuildAuthURL("demo-client", "http://localhost:8080/callback", "test-state", "test-challenge")
	fmt.Println("3. Built auth URL:", builtURL[:70], "...")

	mockIDToken := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwiZW1haWwiOiJhbGljZUBleGFtcGxlLmNvbSIsIm5hbWUiOiJBbGljZSIsImlzcyI6Imh0dHBzOi8vYXV0aC5leGFtcGxlLmNvbSIsImV4cCI6OTk5OTk5OTk5OX0"
	idToken, _ := ParseIDToken(mockIDToken)
	fmt.Printf("4. ID Token: sub=%s email=%s name=%s\n", idToken.Subject, idToken.Email, idToken.Name)
}
