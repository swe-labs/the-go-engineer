package main

import (
	"regexp"
	"testing"
)

func TestSecurityAuditAdd(t *testing.T) {
	audit := SecurityAudit{}
	audit.Add("Test", "check1", "description", true)

	if len(audit.Checks) != 1 {
		t.Errorf("expected 1 check, got %d", len(audit.Checks))
	}
	if audit.Checks[0].Category != "Test" {
		t.Errorf("category = %q, want Test", audit.Checks[0].Category)
	}
}

func TestSecurityAuditSummary(t *testing.T) {
	audit := SecurityAudit{}
	audit.Add("A", "pass", "passes", true)
	audit.Add("B", "fail", "fails", false)

	summary := audit.Summary()
	if !regexp.MustCompile(`1/2 checks passed, 1 failed`).MatchString(summary) {
		t.Errorf("unexpected summary: %q", summary)
	}
}

func TestCheckInputValidation(t *testing.T) {
	tests := []struct {
		name      string
		endpoints []string
		want      bool
	}{
		{
			name:      "all endpoints validated",
			endpoints: []string{"/api/user/{id}/validate", "/api/orders/validate"},
			want:      true,
		},
		{
			name:      "endpoint without validation",
			endpoints: []string{"/api/user/{id}"},
			want:      false,
		},
		{
			name:      "no path params",
			endpoints: []string{"/api/health"},
			want:      true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := checkInputValidation(tc.endpoints)
			if got != tc.want {
				t.Errorf("checkInputValidation(%v) = %v, want %v", tc.endpoints, got, tc.want)
			}
		})
	}
}

func TestCheckAuthRequired(t *testing.T) {
	tests := []struct {
		name   string
		routes map[string]bool
		want   bool
	}{
		{
			name:   "all routes require auth",
			routes: map[string]bool{"/api/user": true, "/api/admin": true},
			want:   true,
		},
		{
			name:   "some routes lack auth",
			routes: map[string]bool{"/api/user": true, "/api/health": false},
			want:   false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := checkAuthRequired(tc.routes)
			if got != tc.want {
				t.Errorf("checkAuthRequired(%v) = %v, want %v", tc.routes, got, tc.want)
			}
		})
	}
}

func TestCheckNoSecretsInCode(t *testing.T) {
	tests := []struct {
		name string
		code string
		want bool
	}{
		{
			name: "no secrets found",
			code: `const dbURL = os.Getenv("DB_URL")`,
			want: true,
		},
		{
			name: "hardcoded password",
			code: `password = "s3cret"`,
			want: false,
		},
		{
			name: "hardcoded api key",
			code: `api_key = "sk-live-abc123"`,
			want: false,
		},
		{
			name: "private key",
			code: "-----BEGIN RSA PRIVATE KEY-----\nMIIEpAIBAAKCAQEA...\n-----END RSA PRIVATE KEY-----",
			want: false,
		},
		{
			name: "safe usage of env var",
			code: `apiKey := os.Getenv("API_KEY")`,
			want: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := checkNoSecretsInCode(tc.code)
			if got != tc.want {
				t.Errorf("checkNoSecretsInCode(%q) = %v, want %v", tc.code, got, tc.want)
			}
		})
	}
}

func TestCheckRateLimited(t *testing.T) {
	tests := []struct {
		name      string
		endpoints []string
		want      bool
	}{
		{
			name:      "all endpoints rate limited",
			endpoints: []string{"/api/user rate-limited", "/api/orders rate-limited"},
			want:      true,
		},
		{
			name:      "some endpoints not rate limited",
			endpoints: []string{"/api/user", "/api/orders rate-limited"},
			want:      false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := checkRateLimited(tc.endpoints)
			if got != tc.want {
				t.Errorf("checkRateLimited(%v) = %v, want %v", tc.endpoints, got, tc.want)
			}
		})
	}
}
