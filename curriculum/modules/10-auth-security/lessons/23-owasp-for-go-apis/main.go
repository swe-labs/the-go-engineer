package main

import (
	"fmt"
	"regexp"
	"strings"
)

type SecurityCheck struct {
	Category    string `json:"category"`
	Check       string `json:"check"`
	Passed      bool   `json:"passed"`
	Description string `json:"description"`
}

type SecurityAudit struct {
	Checks []SecurityCheck `json:"checks"`
}

func (sa *SecurityAudit) Add(category, check, description string, passed bool) {
	sa.Checks = append(sa.Checks, SecurityCheck{
		Category:    category,
		Check:       check,
		Passed:      passed,
		Description: description,
	})
}

func (sa *SecurityAudit) Summary() string {
	passed := 0
	failed := 0
	for _, c := range sa.Checks {
		if c.Passed {
			passed++
		} else {
			failed++
		}
	}
	return fmt.Sprintf("%d/%d checks passed, %d failed", passed, len(sa.Checks), failed)
}

func checkInputValidation(endpoints []string) bool {
	for _, e := range endpoints {
		if strings.Contains(e, "{") && !strings.Contains(e, "validate") {
			return false
		}
	}
	return true
}

func checkAuthRequired(routes map[string]bool) bool {
	for _, requiresAuth := range routes {
		if !requiresAuth {
			return false
		}
	}
	return true
}

func checkRateLimited(endpoints []string) bool {
	for _, e := range endpoints {
		if !strings.Contains(e, "rate-limited") {
			return false
		}
	}
	return true
}

func checkNoSecretsInCode(code string) bool {
	patterns := []string{
		`password\s*=\s*["'][^"']+["']`,
		`api[_-]?key\s*=\s*["'][^"']+["']`,
		`secret\s*=\s*["'][^"']+["']`,
		`sk-(live|test)-[a-zA-Z0-9]+`,
		`-----BEGIN (RSA|EC) PRIVATE KEY-----`,
	}
	for _, p := range patterns {
		re := regexp.MustCompile(`(?i)` + p)
		if re.MatchString(code) {
			return false
		}
	}
	return true
}

func main() {
	fmt.Println("=== OWASP API Security Checklist for Go ===")

	audit := SecurityAudit{}

	fmt.Println("\n1. Broken Object Level Authorization (BOLA)")
	fmt.Println("   Check: Are user IDs validated against the authenticated user?")
	audit.Add("BOLA", "user ID validation", "Ensure users can only access their own resources", false)

	fmt.Println("\n2. Broken Authentication")
	fmt.Println("   Check: Is JWT validation using constant-time comparison?")
	audit.Add("Authentication", "constant-time comparison", "Use crypto/subtle.ConstantTimeCompare", true)

	fmt.Println("\n3. Excessive Data Exposure")
	fmt.Println("   Check: Are API responses limited to necessary fields?")
	audit.Add("Data Exposure", "response field filtering", "Use struct tags or view models", false)

	fmt.Println("\n4. Rate Limiting")
	routes := []string{"/api/user rate-limited", "/api/orders rate-limited"}
	rateLimited := checkRateLimited(routes)
	audit.Add("Rate Limiting", "rate limit on all endpoints", "Apply rate limiting middleware globally", rateLimited)

	fmt.Println("\n5. Input Validation")
	endpoints := []string{"/api/user/{id}/validate", "/api/orders"}
	inputValid := checkInputValidation(endpoints)
	audit.Add("Input Validation", "validate all inputs", "Validate types, lengths, and ranges", inputValid)

	fmt.Println("\n6. Security Misconfiguration")
	audit.Add("Security Config", "CORS allowlist", "Do not use Access-Control-Allow-Origin: *", true)
	audit.Add("Security Config", "TLS min version", "Set MinVersion: tls.VersionTLS12", true)

	fmt.Println("\n7. Secrets in Code")
	code := `clientSecret := "sk-live-abc123def456"`
	noSecrets := checkNoSecretsInCode(code)
	audit.Add("Secrets Management", "no hardcoded secrets", "Use environment variables or secret store", noSecrets)

	fmt.Println("\n=== Audit Summary ===")
	for _, c := range audit.Checks {
		status := "PASS"
		if !c.Passed {
			status = "FAIL"
		}
		fmt.Printf("  [%s] %s: %s\n", status, c.Category, c.Description)
	}
	fmt.Println("\n", audit.Summary())

	fmt.Println("\nNote: OWASP API Security Top 10 is available at owasp.org.")
}
