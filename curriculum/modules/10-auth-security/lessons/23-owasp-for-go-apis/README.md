# OWASP for Go APIs

## Learning objective

Apply the OWASP API Security Top 10 to Go API development, identify common Go-specific vulnerabilities, and maintain a security checklist for API development and code review.

## Why this matters

The OWASP API Security Top 10 is the definitive guide to API-specific security risks. While general OWASP Top 10 covers web applications, API security has unique concerns: excessive data exposure, broken object-level authorization, and mass assignment. Go APIs are not immune -- in fact, Go's simplicity can lead developers to skip security controls they would add in more verbose frameworks. Understanding the OWASP API Security Top 10 helps Go engineers build APIs that are secure by default and pass security audits.

## Mental model

OWASP is the periodic table of web security risks. It categorizes the elements of insecurity: injection, broken auth, data exposure, and so on. Each element has known properties, detection methods, and remediation. The OWASP API Security Top 10 is the subset most relevant to API development.

Think of it as a pre-flight checklist for pilots. You do not skip items because "nothing bad happened last time." Each item addresses a real risk that has caused real breaches. Following the checklist makes security a repeatable process rather than a reactive scramble.

## Core idea

OWASP API Security Top 10 (2023):

| Rank | Category | Go-specific risk |
|---|---|---|
| API1 | Broken Object Level Authorization | Missing user ID validation in handlers |
| API2 | Broken Authentication | JWT without constant-time comparison |
| API3 | Broken Object Property Level Authorization | Mass assignment via JSON unmarshal |
| API4 | Unrestricted Resource Consumption | No rate limiting on expensive endpoints |
| API5 | Broken Function Level Authorization | Missing admin endpoint checks |
| API6 | Unrestricted Access to Sensitive Business Flows | No CAPTCHA on critical flows |
| API7 | Server Side Request Forgery | Unsanitized URL in http.Get |
| API8 | Security Misconfiguration | CORS wildcard, debug enabled |
| API9 | Improper Inventory Management | Shadow APIs, outdated endpoints |
| API10 | Unsafe Consumption of APIs | No TLS verification on downstream calls |

## Under the hood

OWASP API Security Top 10 is developed by the OWASP community based on real-world incident data. Each category includes:

- **Description**: What the risk is and how it manifests.
- **Example attack scenario**: Realistic walkthrough of an exploitation.
- **Prevention measures**: Concrete actions to prevent the vulnerability.
- **References**: Links to standards, tools, and further reading.

The list is updated every 3-4 years based on breach data, community surveys, and expert analysis. The API Security Top 10 specifically focuses on RESTful APIs, GraphQL, gRPC, and other API protocols, not traditional web applications.

## How Go uses it

Go's standard library and ecosystem provide tools that address many OWASP categories:

| OWASP category | Go defense |
|---|---|
| BOLA | Custom middleware checking user ID against session |
| Broken Auth | `golang.org/x/crypto` for bcrypt, constant-time compare |
| Data Exposure | Struct tags with `json:"-"` to hide fields |
| Rate Limiting | `golang.org/x/time/rate` or custom middleware |
| Security Misconfig | `crypto/tls` with `MinVersion`, CORS allowlist |
| SSRF | URL validation with allowlist of hostnames |
| Injection | Parameterized queries with `database/sql` |
| Mass Assignment | Explicit binding (only accept allowed fields) |

## Go example

```go
package main

import (
	"fmt"
	"regexp"
	"strings"
)

type SecurityCheck struct {
	Category    string
	Check       string
	Passed      bool
	Description string
}

type SecurityAudit struct {
	Checks []SecurityCheck
}

func (sa *SecurityAudit) Add(category, check, description string, passed bool) {
	sa.Checks = append(sa.Checks, SecurityCheck{category, check, passed, description})
}

func (sa *SecurityAudit) Summary() string {
	passed, failed := 0, 0
	for _, c := range sa.Checks {
		if c.Passed {
			passed++
		} else {
			failed++
		}
	}
	return fmt.Sprintf("%d/%d checks passed, %d failed", passed, len(sa.Checks), failed)
}

func checkNoSecretsInCode(code string) bool {
	patterns := []string{
		`password\s*=\s*["'][^"']+["']`,
		`api[_-]?key\s*=\s*["'][^"']+["']`,
		`sk-(live|test)-[a-zA-Z0-9]+`,
	}
	for _, p := range patterns {
		if regexp.MustCompile(`(?i)` + p).MatchString(code) {
			return false
		}
	}
	return true
}

func main() {
	audit := SecurityAudit{}
	audit.Add("BOLA", "user ID validation", "Validate user ID against session", false)
	audit.Add("Secrets", "no hardcoded secrets", checkNoSecretsInCode(`password = "s3cret"`), false)
	fmt.Println(audit.Summary())
}
```

## Step-by-step execution

For an OWASP API Security audit of a Go API:

1. **Inventory**: List all API endpoints, methods, authentication requirements, and rate limits.
2. **BOLA check**: For each endpoint that takes a user ID, verify that the authenticated user can only access their own data. Example: `GET /api/user/{id}` must check `id == session.userID`.
3. **Authentication check**: Verify that password hashing uses bcrypt, JWTs are validated with constant-time comparison (`crypto/subtle`), and sessions use secure cookies.
4. **Data exposure check**: Review JSON response structs for fields tagged with `json:"-"` or use view models that exclude sensitive fields.
5. **Rate limiting check**: Confirm that all production endpoints have rate limiting middleware, especially login, registration, and password reset.
6. **Security config check**: Verify CORS allowlist, TLS minimum version, CSP headers, and that debug mode is disabled.
7. **Injection check**: Confirm all database queries use parameterized statements.
8. **Secrets check**: Search codebase for hardcoded credentials using `grep` or `gitleaks`.
9. **SSRF check**: Review all `http.Get`, `http.Post`, and similar calls that use user-supplied URLs.
10. **Inventory check**: Verify there are no undocumented endpoints (shadow APIs), outdated versions, or deprecated endpoints that are still accessible.

## Common mistakes

- Mistake: Assuming OWASP is only relevant for traditional web applications, not JSON APIs.
  - Why it happens: Developers think "we only serve JSON, not HTML, so XSS and injection do not apply."
  - Fix: Injection attacks work on any input channel (SQL injection, NoSQL injection, LDAP injection). BOLA and data exposure are API-specific.

- Mistake: Implementing security controls in isolation without referencing a threat model.
  - Why it happens: Developers add a security feature (JWT, CORS, rate limiting) without understanding what threat it addresses.
  - Fix: Create a threat model for your API. For each attack vector, identify which controls address it.

- Mistake: Focusing only on authentication while ignoring authorization, rate limiting, and input validation.
  - Why it happens: Authentication is the most visible security control (login pages, JWT, OAuth).
  - Fix: Use the OWASP API Security Top 10 as a checklist to ensure all categories are addressed.

- Mistake: Treating OWASP recommendations as a compliance checkbox instead of a risk-based framework.
  - Why it happens: Auditors require OWASP compliance, so teams implement controls without understanding them.
  - Fix: Each control has a purpose. Understand the attack it prevents and test that it actually works.

- Mistake: Not securing the API gateway while having per-service security.
  - Why it happens: Teams secure individual microservices but leave the gateway (which routes to all services) unprotected.
  - Fix: The API gateway should enforce authentication, rate limiting, and input validation before requests reach backend services.

## Debugging walkthrough

Consider this Go handler that is vulnerable to Broken Object Level Authorization:

```go
func getUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["id"]
	user, err := db.GetUserByID(userID)
	if err != nil {
		http.Error(w, "not found", 404)
		return
	}
	json.NewEncoder(w).Encode(user)
}
```

Symptom: Any authenticated user can access any other user's data by changing the ID in the URL.

Investigation: Check if there is any authorization check between getting the user ID from the URL and querying the database.

Root cause: The handler uses `mux.Vars(r)["id"]` directly without verifying that the authenticated user owns that ID.

Fix: Add an authorization check:

```go
func getUserHandler(w http.ResponseWriter, r *http.Request) {
	authUserID := r.Context().Value("userID").(string)
	vars := mux.Vars(r)
	requestedUserID := vars["id"]

	if authUserID != requestedUserID {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	user, err := db.GetUserByID(requestedUserID)
	if err != nil {
		http.Error(w, "not found", 404)
		return
	}
	json.NewEncoder(w).Encode(user)
}
```

## Production notes

- Create an OWASP API Security checklist specific to your organization's tech stack. Include Go-specific items: parameterized queries, `html/template` usage, crypto/rand for tokens.
- Automate security checks: use `gosec` for static analysis, `govulncheck` for dependency vulnerabilities, and custom linters for company-specific rules.
- Conduct regular security reviews using the OWASP checklist. Document findings and track remediation.
- Monitor the OWASP API Security Top 10 updates (typically every 3-4 years) and update your checklist accordingly.
- Train developers on OWASP Top 10 with Go-specific examples during onboarding.
- Use security headers middleware (CSP, X-Frame-Options, X-Content-Type-Options, HSTS) as a standard part of your HTTP server configuration.

## Performance implications

- Most OWASP controls have negligible performance impact: header validation, authorization checks, and input validation are O(1) operations.
- Rate limiting middleware adds microseconds per request.
- The main cost is development and review time, not runtime performance.
- Some controls (e.g., input sanitization, payload size limits) can add latency if scanning large payloads. Apply them early in the middleware chain.

## Practice task

Write a function `owaspAudit(routes []Route) *SecurityAudit` that performs an automated OWASP API Security audit on a list of API routes.

Define a `Route` type: `type Route struct { Method, Path string; AuthRequired bool; RateLimited bool; InputValidated bool }`.

The audit function should check:
1. All routes with path parameters have input validation.
2. All non-GET routes require authentication.
3. All routes accept rate limiting headers.
4. No route exposes a "secret" or "password" field in its response (simulate by checking if the path contains "secret").

Return a `SecurityAudit` with individual check results.

Write a `main()` that creates several routes (some with issues), runs the audit, and prints detailed results.

## Tests / verification

```bash
go run ./curriculum/modules/10-auth-security/lessons/23-owasp-for-go-apis
go test ./curriculum/modules/10-auth-security/lessons/23-owasp-for-go-apis
```

The existing tests verify the security audit mechanism, input validation detection, authentication requirement checking, rate limiting detection, and hardcoded secret detection.

## Review questions

1. What is the difference between Broken Object Level Authorization (BOLA) and Broken Function Level Authorization?
2. How would you prevent mass assignment in a Go API that accepts JSON request bodies?
3. Why is excessive data exposure particularly dangerous in APIs compared to traditional web applications?
4. What is a shadow API and how does it relate to the Improper Inventory Management category?
5. How does the OWASP API Security Top 10 differ from the general OWASP Top 10?

## NEXT UP

Secure logging -- logging security events without exposing sensitive data, structured logging with redaction, audit trails, and log levels in Go.
