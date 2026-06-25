package main

import (
	"strings"
	"testing"
)

func TestBuildAuthURL(t *testing.T) {
	url := BuildAuthURL("client-1", "https://example.com/cb", "state-xyz", "challenge-abc")
	if !strings.Contains(url, "response_type=code") {
		t.Error("expected response_type=code")
	}
	if !strings.Contains(url, "client_id=client-1") {
		t.Error("expected client_id=client-1")
	}
	if !strings.Contains(url, "state=state-xyz") {
		t.Error("expected state=state-xyz")
	}
	if !strings.Contains(url, "code_challenge=challenge-abc") {
		t.Error("expected code_challenge=challenge-abc")
	}
	if !strings.Contains(url, "code_challenge_method=S256") {
		t.Error("expected code_challenge_method=S256")
	}
}

func TestBuildAuthURL_IncludesRedirectURI(t *testing.T) {
	url := BuildAuthURL("c1", "https://example.com/cb", "s1", "ch1")
	if !strings.Contains(url, "redirect_uri=") {
		t.Error("expected redirect_uri parameter")
	}
	if !strings.Contains(url, "example.com") {
		t.Error("expected example.com in redirect_uri")
	}
}

func TestGeneratePKCE(t *testing.T) {
	pkce, err := GeneratePKCE()
	if err != nil {
		t.Fatal(err)
	}
	if len(pkce.Verifier) == 0 {
		t.Error("expected non-empty verifier")
	}
	if len(pkce.Challenge) == 0 {
		t.Error("expected non-empty challenge")
	}
	if pkce.ChallengeMethod != "S256" {
		t.Errorf("expected S256, got %s", pkce.ChallengeMethod)
	}
}

func TestParseIDToken_Valid(t *testing.T) {
	mockToken := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwiZW1haWwiOiJhbGljZUBleGFtcGxlLmNvbSIsIm5hbWUiOiJBbGljZSIsImlzcyI6Imh0dHBzOi8vYXV0aC5leGFtcGxlLmNvbSIsImV4cCI6OTk5OTk5OTk5OX0.signature"
	idToken, err := ParseIDToken(mockToken)
	if err != nil {
		t.Fatal(err)
	}
	if idToken.Subject != "1234567890" {
		t.Errorf("expected 1234567890, got %s", idToken.Subject)
	}
	if idToken.Email != "alice@example.com" {
		t.Errorf("expected alice@example.com, got %s", idToken.Email)
	}
	if idToken.Name != "Alice" {
		t.Errorf("expected Alice, got %s", idToken.Name)
	}
}

func TestParseIDToken_InvalidFormat(t *testing.T) {
	_, err := ParseIDToken("invalid-token")
	if err == nil {
		t.Fatal("expected error for invalid JWT format")
	}
}

func TestGeneratePKCE_Unique(t *testing.T) {
	p1, _ := GeneratePKCE()
	p2, _ := GeneratePKCE()
	if p1.Verifier == p2.Verifier {
		t.Error("expected unique verifiers")
	}
}
