# Threat modeling

## Learning objective

Build and analyze threat models using STRIDE, DREAD, and attack trees to identify security risks in Go web services, and translate each identified threat into a concrete security requirement.

## Why this matters

Security vulnerabilities are cheapest to fix during design and most expensive after deployment. Threat modeling shifts security left: instead of finding bugs in production, you find design flaws before writing code. Every production breach at major companies can be traced to a threat that was identified but not modeled, or modeled but not mitigated. Professional Go engineers who threat-model produce services that survive penetration testing and pass security reviews on the first pass.

## Mental model

Imagine your Go service as a medieval castle. STRIDE is a checklist of ways an attacker can breach it: Spoofing (fake identity), Tampering (alter data), Repudiation (deny acting), Information disclosure (leak secrets), Denial of service (crash the castle), Elevation of privilege (guard becomes king). DREAD helps you rank which threats to fix first. Attack trees are the paths an attacker could walk to reach the treasure room. Every wall, gate, and guard post is a trust boundary where you must ask: what could go wrong here?

## Core idea

A threat model has four components:

- **Assets**: what you protect (user credentials, payment data, session tokens)
- **Trust boundaries**: where data crosses between trust levels (network edge, database client, internal RPC)
- **Attack surface**: all entry points an attacker can touch (HTTP endpoints, WebSocket connections, file uploads)
- **Threats**: specific things that can go wrong (SQL injection at login, JWT theft via XSS, privilege escalation at `/admin`)

The output of threat modeling is a prioritized list of mitigations, not a diagram that collects dust.

## Under the hood

STRIDE was created by Microsoft in the 1990s and remains the industry standard for classifying threats. Each letter maps to a security property:

| Letter | Threat | Violates |
|---|---|---|
| S | Spoofing | Authentication |
| T | Tampering | Integrity |
| R | Repudiation | Non-repudiation |
| I | Information disclosure | Confidentiality |
| D | Denial of service | Availability |
| E | Elevation of privilege | Authorization |

DREAD scores each threat on five dimensions (Damage, Reproducibility, Exploitability, Affected users, Discoverability), each rated 1-10. The sum (or average) prioritizes fixes.

Attack trees start with a goal ("read another user's private messages") and branch into sub-goals ("bypass auth", "guess session token", "SQL injection on messages endpoint"). Each leaf is a potential exploit path that maps to one or more STRIDE categories.

## How Go uses it

Go's standard library and ecosystem support threat modeling at the code level:

- **`net/http` middleware chains** implement trust boundaries: auth middleware before handler logic
- **`crypto/subtle`** prevents timing attacks on comparison operations
- **`database/sql` parameterized queries** prevent SQL injection (a tampering threat)
- **`context.Context`** propagates identity and deadline information across goroutines, making data flow visible
- **`io.Reader`/`io.Writer` interfaces** make data flow tracing straightforward for DFDs

In a typical Go service, you model threats at each handler, each database call, and each external API call.

## Go example

```go
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

// Asset represents a protected resource in the system.
type Asset struct {
	ID       string `json:"id"`
	OwnerID  string `json:"owner_id"`
	Payload  string `json:"payload,omitempty"`
}

// Threat models a single identified threat.
type Threat struct {
	ID          string
	Category    string // STRIDE category
	Asset       string
	Description string
	Severity    int // DREAD score 1-10
	Mitigation  string
}

// analyzeThreats applies STRIDE categories to a set of assets.
func analyzeThreats(assets []Asset) []Threat {
	var threats []Threat
	for _, a := range assets {
		// Spoofing: can an attacker claim ownership?
		threats = append(threats, Threat{
			ID:          "S-" + a.ID,
			Category:    "Spoofing",
			Asset:       a.ID,
			Description: fmt.Sprintf("Attacker modifies OwnerID of asset %s", a.ID),
			Severity:    8,
			Mitigation:  "Validate JWT claims before any mutation; store OwnerID server-side only",
		})
		// Tampering: can an attacker modify the payload?
		threats = append(threats, Threat{
			ID:          "T-" + a.ID,
			Category:    "Tampering",
			Asset:       a.ID,
			Description: fmt.Sprintf("Attacker modifies payload of asset %s via man-in-the-middle", a.ID),
			Severity:    9,
			Mitigation:  "Enforce TLS; sign payloads with HMAC for idempotency",
		})
		// Information disclosure: can an attacker read another user's asset?
		threats = append(threats, Threat{
			ID:          "I-" + a.ID,
			Category:    "Information Disclosure",
			Asset:       a.ID,
			Description: fmt.Sprintf("Attacker reads asset %s belonging to another user", a.ID),
			Severity:    7,
			Mitigation:  "Row-level authorization: SELECT queries must include owner_id filter",
		})
	}
	return threats
}

func main() {
	assets := []Asset{
		{ID: "doc-1", OwnerID: "user-42", Payload: "confidential report"},
		{ID: "doc-2", OwnerID: "user-99", Payload: "financial data"},
	}
	threats := analyzeThreats(assets)
	for _, t := range threats {
		fmt.Printf("[%s] %s (severity: %d): %s\n  Mitigation: %s\n\n", t.Category, t.Description, t.Severity, t.Mitigation)
	}

	// Simulate DREAD prioritization: sort by severity descending
	fmt.Println("=== Prioritized threats (highest severity first) ===")
	ordered := prioritizeBySeverity(threats)
	for _, t := range ordered {
		fmt.Printf("[%d] %s: %s\n", t.Severity, t.ID, t.Description)
	}

	// Start a simple server (trust boundary demonstration)
	http.HandleFunc("/api/assets/", assetHandler(assets))
	log.Println("Server listening on :8080")
	// In production, use log.Fatal(http.ListenAndServeTLS(...))
}

func prioritizeBySeverity(threats []Threat) []Threat {
	// Simple insertion sort by severity descending
	out := make([]Threat, len(threats))
	copy(out, threats)
	for i := 1; i < len(out); i++ {
		for j := i; j > 0 && out[j-1].Severity < out[j].Severity; j-- {
			out[j-1], out[j] = out[j], out[j-1]
		}
	}
	return out
}

func assetHandler(assets []Asset) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Trust boundary: every handler must authenticate the caller
		userID := r.Header.Get("X-User-ID")
		if userID == "" {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		// Authorization check: user can only access their own assets
		assetID := strings.TrimPrefix(r.URL.Path, "/api/assets/")
		for _, a := range assets {
			if a.ID == assetID {
				if a.OwnerID != userID {
					http.Error(w, "forbidden", http.StatusForbidden)
					return
				}
				json.NewEncoder(w).Encode(a)
				return
			}
		}
		http.Error(w, "not found", http.StatusNotFound)
	}
}
```

## Step-by-step execution

For the threat analysis above with assets `doc-1` (owner user-42) and `doc-2` (owner user-99):

1. `analyzeThreats` iterates over each asset.
2. For `doc-1`, three threats are created: spoofing (modify owner), tampering (modify payload), information disclosure (read by another user).
3. Each threat receives a DREAD-like severity score.
4. `prioritizeBySeverity` sorts threats: tampering (9), spoofing (8), information disclosure (7), repeated for `doc-2`.
5. The output prints each threat with its mitigation strategy.
6. The HTTP handler demonstrates the trust boundary: `X-User-ID` header is the gate. If missing `->` 401. If present but not the owner `->` 403.

## Common mistakes

- Drawing threat models after an incident instead of during design — post-hoc models miss the assumptions that led to the breach.
- Modeling only external threats and ignoring insider risk — most breaches involve compromised credentials.
- Stopping at diagram creation without tracing trust boundaries, data flows, and attack surfaces through the codebase.
- Rating all threats as high severity — without prioritization, everything is urgent and nothing gets fixed.
- Forgetting the human element: social engineering, phishing, and credential stuffing are threats too.

## Debugging walkthrough

Consider this Go handler that leaks tenant data:

```go
func searchHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query().Get("q")
		rows, _ := db.Query("SELECT data FROM documents WHERE content LIKE '%" + query + "%'")
		// ...
	}
}
```

**Threat model analysis**: This handler has an SQL injection threat (Tampering) and an information disclosure threat (no tenant filter).

**Investigation**: Apply STRIDE. The `q` parameter crosses a trust boundary (HTTP -> SQL). There is no input validation (allowlist). The query string is concatenated directly, enabling SQL injection. There is no `WHERE tenant_id = ?` clause, so any user can search all documents.

**Root cause**: No threat model existed for this handler. The developer did not ask "what could go wrong at this trust boundary?"

**Fix**: Parameterized queries and tenant-scoped filtering:

```go
func searchHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tenantID := r.Context().Value("tenant_id")
		query := r.URL.Query().Get("q")
		rows, _ := db.Query("SELECT data FROM documents WHERE tenant_id = $1 AND content LIKE '%' || $2 || '%'", tenantID, query)
		// ...
	}
}
```

## Production notes

- Every endpoint should have an associated threat model entry. Link threat IDs in code comments for audit trails.
- Automate threat model reviews in CI: require a THREAT_MODEL.md update for any new endpoint.
- Use tools like Threat Dragon, OWASP Threat Dragon, or Microsoft TMT for diagram-based modeling.
- For Go gRPC services, model threats at each RPC call boundary. gRPC interceptor chains are trust boundaries.
- Integrate threat modeling with your incident response runbooks: each threat should map to a detection rule.

## Performance implications

Threat modeling itself has no runtime cost — it is a design-time activity. However, mitigations derived from threat modeling do:

- Input validation adds CPU overhead per request but prevents injection attacks.
- Authorization checks add database queries or cache lookups per request.
- Rate limiting uses memory for token buckets and CPU for time comparisons.
- Encryption (TLS, AEAD) adds 5-15% CPU overhead but is mandatory for data in transit.

Profile your mitigations: a bcrypt hash (cost 12) takes ~250ms on typical hardware — placing it on every request is a DoS vector itself.

## Practice task

Write a Go function `BuildThreatModel(endpoints []string, assets []string) map[string][]Threat` that:
- Accepts a list of API endpoints and a list of asset names.
- For each endpoint, identifies which assets are reachable through it.
- Applies STRIDE to each endpoint-asset pair and returns at least one threat per category (S, T, I, D, E).
- Skips repudiation (R) for simplicity.
- Returns a map keyed by endpoint with a slice of threats.

Then write a `main()` that defines three endpoints (`POST /login`, `GET /documents`, `DELETE /admin/users`) and three assets (`password_hash`, `document_content`, `admin_session`) and prints the threat model as formatted text.

## Tests / verification

```bash
go run ./curriculum/modules/10-auth-security/lessons/01-threat-modeling
go test ./curriculum/modules/10-auth-security/lessons/01-threat-modeling
```

The test file `main_test.go` contains table-driven tests that verify:
- `BuildThreatModel` returns threats for every endpoint-asset pair.
- Each threat has a non-empty `Category` that is one of the STRIDE letters.
- Each threat has a `Severity` between 1 and 10.

## Review questions

1. What is the difference between a trust boundary and an attack surface in a threat model?
2. Given a `POST /transfer` endpoint that moves money between accounts, walk through each STRIDE category and name one threat.
3. Why is repudiation (R) particularly important for financial applications? How would you mitigate it in Go?
4. If your DREAD analysis gives all threats a severity of 9-10, what does that indicate about your threat model?
5. An API returns `UserID` in the URL path (`GET /api/users/{id}`) and uses it directly in SQL queries. What STRIDE threats apply?

## NEXT UP

Trust boundaries and how they define the lines between different trust levels in your Go architecture.
