package main

import (
	"fmt"
	"strings"
)

type Asset struct {
	ID      string
	OwnerID string
	Payload string
}

type Threat struct {
	ID          string
	Category    string
	Asset       string
	Description string
	Severity    int
	Mitigation  string
}

func analyzeThreats(assets []Asset) []Threat {
	var threats []Threat
	for _, a := range assets {
		threats = append(threats, Threat{
			ID:          "S-" + a.ID,
			Category:    "Spoofing",
			Asset:       a.ID,
			Description: fmt.Sprintf("Attacker modifies OwnerID of asset %s", a.ID),
			Severity:    8,
			Mitigation:  "Validate JWT claims before any mutation; store OwnerID server-side only",
		})
		threats = append(threats, Threat{
			ID:          "T-" + a.ID,
			Category:    "Tampering",
			Asset:       a.ID,
			Description: fmt.Sprintf("Attacker modifies payload of asset %s via man-in-the-middle", a.ID),
			Severity:    9,
			Mitigation:  "Enforce TLS; sign payloads with HMAC for idempotency",
		})
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

func prioritizeBySeverity(threats []Threat) []Threat {
	out := make([]Threat, len(threats))
	copy(out, threats)
	for i := 1; i < len(out); i++ {
		for j := i; j > 0 && out[j-1].Severity < out[j].Severity; j-- {
			out[j-1], out[j] = out[j], out[j-1]
		}
	}
	return out
}

func BuildThreatModel(endpoints []string, assets []string) map[string][]Threat {
	result := make(map[string][]Threat)
	for _, ep := range endpoints {
		var endpointThreats []Threat
		for _, asset := range assets {
			epLower := strings.ToLower(ep)
			if strings.Contains(epLower, "login") {
				endpointThreats = append(endpointThreats, Threat{
					ID:          "S-" + asset,
					Category:    "Spoofing",
					Asset:       asset,
					Description: fmt.Sprintf("Attacker brute-forces credentials at %s to access %s", ep, asset),
					Severity:    9,
					Mitigation:  "Rate-limit login; enforce account lockout after N failures",
				})
			}
			if strings.Contains(epLower, "documents") || strings.Contains(epLower, "admin") {
				endpointThreats = append(endpointThreats, Threat{
					ID:          "I-" + asset,
					Category:    "Information Disclosure",
					Asset:       asset,
					Description: fmt.Sprintf("Unauthenticated user reads %s via %s", asset, ep),
					Severity:    8,
					Mitigation:  "Require auth middleware on " + ep,
				})
				endpointThreats = append(endpointThreats, Threat{
					ID:          "E-" + asset,
					Category:    "Elevation of Privilege",
					Asset:       asset,
					Description: fmt.Sprintf("Low-privilege user accesses %s through %s", asset, ep),
					Severity:    10,
					Mitigation:  "Enforce RBAC middleware on " + ep,
				})
			}
		}
		if len(endpointThreats) > 0 {
			result[ep] = endpointThreats
		}
	}
	return result
}

func main() {
	assets := []Asset{
		{ID: "doc-1", OwnerID: "user-42", Payload: "confidential report"},
		{ID: "doc-2", OwnerID: "user-99", Payload: "financial data"},
	}
	threats := analyzeThreats(assets)
	for _, t := range threats {
		fmt.Printf("[%s] %s (severity: %d)\n  Mitigation: %s\n\n", t.Category, t.Description, t.Severity, t.Mitigation)
	}

	fmt.Println("=== Prioritized threats ===")
	ordered := prioritizeBySeverity(threats)
	for _, t := range ordered {
		fmt.Printf("[%d] %s: %s\n", t.Severity, t.ID, t.Description)
	}

	fmt.Println("\n=== BuildThreatModel ===")
	model := BuildThreatModel(
		[]string{"POST /login", "GET /documents", "DELETE /admin/users"},
		[]string{"password_hash", "document_content", "admin_session"},
	)
	for ep, thrs := range model {
		fmt.Printf("Endpoint: %s\n", ep)
		for _, t := range thrs {
			fmt.Printf("  [%s] %s\n", t.Category, t.Description)
		}
	}
}
