# Architecture decision records

## Learning objective

Write architecture decision records (ADRs) using the standard format, manage the ADR lifecycle from proposal to supersession, document tradeoffs with pros and cons, and implement a Go-based ADR management tool.

## Why this matters

Architecture decisions are the most expensive decisions a software team makes. A wrong database choice or service boundary error costs months of rework. Yet most teams make these decisions in Slack threads, hallway conversations, or meeting notes that are lost within weeks. ADRs capture the context, decision, and consequences in a version-controlled document that lives alongside the code. New team members read ADRs to understand why things are the way they are. Future engineers see that the decision was deliberate, not accidental. Every professional Go engineer should write ADRs for non-trivial architecture decisions.

## Mental model

An ADR is a commit message for an architecture decision. It records what was decided, why, what alternatives were considered, and what the consequences are. Like a commit, an ADR has a lifecycle: it is proposed, accepted, and possibly deprecated or superseded by a later ADR.

Think of ADRs as a decision journal. When you face a fork in the road (monolith vs microservices, PostgreSQL vs MySQL, REST vs gRPC), you write an ADR documenting which path you took and why. Six months later, when someone asks "why did we choose Kafka over RabbitMQ?", the ADR has the answer.

## Core idea

The standard ADR format (by Michael Nygard):

```
# ADR-NNN: Title

## Status
proposed | accepted | deprecated | superseded

## Context
The situation that motivated the decision.
What is the problem being solved?
What constraints exist?

## Decision
The chosen approach. Full sentences, active voice.
"This ADR proposes that we use X because Y."

## Consequences
What becomes easier or harder after this decision?
What tradeoffs were accepted?
```

Optional sections: `Alternatives Considered`, `Pros and Cons`, `Consequences`.

ADR lifecycle:
1. **Proposed**: written by any team member, not yet agreed.
2. **Accepted**: agreed by the team, the decision is in effect.
3. **Deprecated**: no longer recommended but still in use.
4. **Superseded**: replaced by a newer ADR (references the new one).

## Under the hood

ADRs should live in the repository at `docs/adr/` with sequential numbering. Each ADR is a markdown file: `docs/adr/ADR-001-use-postgresql.md`. The sequence number must never be reused, even if an ADR is deprecated.

The ADR number is permanent. It becomes a stable reference in code comments, pull request descriptions, and team discussions. When an ADR is superseded, the old ADR includes a `Superseded by ADR-NNN` line in its status, and the new ADR includes a `Supersedes ADR-MMM` line.

ADR review process:
1. Author writes the ADR under `docs/adr/` in a branch.
2. Team reviews in a pull request.
3. Discussion happens in PR comments (which are also documented).
4. When consensus is reached, the PR merges and the ADR becomes accepted.
5. If no consensus, the ADR is not merged or is marked as rejected (a valid state).

## How Go uses it

- **Kubernetes**: the Kubernetes project uses ADRs (called "KEPs" -- Kubernetes Enhancement Proposals) for every significant architecture change.
- **Go standard library**: the Go proposal process (`golang.org/issue`) serves a similar purpose. Proposals include motivation, design, and alternatives.
- **`docs/adr/` directory**: GitHub and GitLab render markdown files automatically, making ADRs discoverable in the repository browser.
- **`adr-tools`**: a command-line tool (written in Bash) for creating and managing ADRs. Several Go reimplementations exist.

## Go example

```go
package main

import (
	"fmt"
	"time"
)

type ADRStatus string

const (
	Proposed   ADRStatus = "proposed"
	Accepted   ADRStatus = "accepted"
	Deprecated ADRStatus = "deprecated"
	Superseded ADRStatus = "superseded"
)

type ADR struct {
	ID           string
	Title        string
	Status       ADRStatus
	Context      string
	Decision     string
	Consequences string
	Alternatives []string
	Date         time.Time
	Supersedes   string
}

func NewADR(id, title string) *ADR {
	return &ADR{ID: id, Title: title, Status: Proposed, Date: time.Now()}
}

func (a *ADR) Accept()     { a.Status = Accepted }
func (a *ADR) Deprecate()  { a.Status = Deprecated }
func (a *ADR) SupersedeBy(newID string) {
	a.Status = Superseded
	a.Supersedes = newID
}

type Manager struct {
	records map[string]*ADR
}

func NewManager() *Manager {
	return &Manager{records: make(map[string]*ADR)}
}

func (m *Manager) Add(a *ADR)            { m.records[a.ID] = a }
func (m *Manager) Get(id string) *ADR    { return m.records[id] }
func (m *Manager) All() []*ADR {
	out := make([]*ADR, 0, len(m.records))
	for _, a := range m.records {
		out = append(out, a)
	}
	return out
}
func (m *Manager) ByStatus(s ADRStatus) []*ADR {
	var out []*ADR
	for _, a := range m.records {
		if a.Status == s {
			out = append(out, a)
		}
	}
	return out
}

func main() {
	m := NewManager()

	adr1 := NewADR("ADR-001", "Use PostgreSQL")
	adr1.Context = "The team needs an ACID-compliant database for transactional workloads."
	adr1.Decision = "We will use PostgreSQL 16 as the primary database."
	adr1.Consequences = "PostgreSQL provides ACID, JSONB, and full-text search."
	adr1.Alternatives = []string{"MySQL 8 - similar ACID but weaker JSON", "MongoDB - eventual consistency unsuitable for payments"}
	adr1.Accept()
	m.Add(adr1)

	adr2 := NewADR("ADR-002", "Use Redis for Caching")
	adr2.Context = "The auth service needs fast session storage."
	adr2.Decision = "We will use Redis 7 with TTL-based expiry."
	adr2.Consequences = "Redis adds operational complexity but provides sub-ms session lookups."
	adr2.Accept()
	m.Add(adr2)

	adr3 := NewADR("ADR-003", "Monolith First")
	adr3.Context = "Team considering microservices vs monolith."
	adr3.Decision = "Start with modular monolith; extract services when the monolith causes measurable pain."
	adr3.Consequences = "Avoids premature distributed complexity. Strict package boundaries required."
	adr3.Accept()
	m.Add(adr3)

	for _, a := range m.All() {
		fmt.Printf("%s: %s [%s]\n", a.ID, a.Title, a.Status)
	}
	fmt.Println("\nAccepted decisions:")
	for _, a := range m.ByStatus(Accepted) {
		fmt.Printf("  %s: %s\n", a.ID, a.Title)
	}
}
```

## Step-by-step execution

For creating ADR-001:

1. `NewADR("ADR-001", "Use PostgreSQL")` creates an ADR with Proposed status and current date.
2. Set `Context`, `Decision`, `Consequences`, `Alternatives`.
3. `adr1.Accept()` sets Status to Accepted.
4. `m.Add(adr1)` stores it in the manager's map.
5. `m.Get("ADR-001")` retrieves the ADR.
6. `m.ByStatus(Accepted)` returns all accepted ADRs including ADR-001.

For ADR lifecycle:
1. ADR-001 is proposed → accepted (after team review).
2. ADR-001 is later deprecated when the team migrates to a new database.
3. ADR-004 (Use CockroachDB) supersedes ADR-001: `adr1.SupersedeBy("ADR-004")`.

## Common mistakes

- **Writing ADRs after the decision is made**: an ADR should be written as part of the decision process, not as post-hoc documentation. The alternatives and tradeoffs are fresh only during the decision.
- **Too much detail**: ADRs are not design documents. Keep them to 1-2 pages. The goal is to capture the rationale, not the full implementation spec.
- **No alternatives section**: listing "alternatives considered" shows the team explored options. A reader in the future who disagrees with the decision can understand what was ruled out and why.
- **Skipping consequences**: every decision has downsides. Documenting them helps future engineers decide whether to revisit the decision.
- **Not updating ADRs**: when a decision is reversed or superseded, update the ADR status. A stale ADR with "accepted" status that no longer applies is misleading.

## Debugging walkthrough

A new engineer asks: "Why are we using RabbitMQ instead of Kafka?"

**Symptom**: No one on the team remembers the rationale. The decision was made 2 years ago.

**Investigation**: Check `docs/adr/`. Find `ADR-007-use-rabbitmq.md`.

The ADR says:
```
Status: accepted

Context: The order service needs to publish events to 2 downstream consumers.
Throughput is 100 msg/s. We do not need event replay or long-term retention.

Decision: Use RabbitMQ for event publishing.

Alternatives:
- Kafka: overkill for 100 msg/s, operational complexity is too high for a 3-person team.
- Redis Pub/Sub: no persistence, messages lost on crash.
```

**Root cause**: There was no ADR. The decision was made in a Slack message. The team does not practice ADR documentation.

**Fix**: Start writing ADRs. Write `ADR-001-use-adrs.md` as the first ADR to formalize the practice. Write ADRs for the last 3 significant architecture decisions retroactively (with context from the Slack history).

## Production notes

- **ADR template**: create a `docs/adr/template.md` with the standard sections. Everyone starts from the same template.
- **ADR review in PRs**: require an ADR (or reference to an existing ADR) in every pull request that makes an architecture- or infrastructure-level change.
- **Numbering**: use sequential numbering (ADR-001, ADR-002). Do not use dates or feature names -- the number is a stable identifier.
- **Supersession chain**: ADR-001 can be superseded by ADR-005. Keep the chain intact: ADR-001 links to ADR-005, and ADR-005 links back to ADR-001.
- **ADR expiry**: periodically review all accepted ADRs. Are they still valid? If not, mark them as deprecated and write a new ADR.

## Performance implications

- **Storage**: ADRs are markdown text files. Storage cost is negligible.
- **Review time**: an ADR review takes 15-30 minutes of team time. This is an investment that pays for itself when it prevents a bad architecture decision.
- **Discovery cost**: finding ADRs in `docs/adr/` is instant. Finding the same information in a wiki, a Notion page, or a Slack message takes 10-60 minutes.
- **Onboarding time**: a new engineer reading 10 ADRs (30 minutes of reading) understands the key architecture decisions of a 5-year-old codebase. Without ADRs, this understanding takes weeks of mentorship and code archaeology.

## Practice task

Write an ADR for a decision you have faced recently (real or hypothetical). Use the standard format with Status, Context, Decision, Consequences, and Alternatives Considered. Then implement a Go function `RenderADR(adr *ADR) string` that returns the ADR as a formatted markdown string.

## Tests / verification

```bash
go run ./curriculum/modules/14-architecture-distributed-systems/lessons/17-architecture-decision-records
go test ./curriculum/modules/14-architecture-distributed-systems/lessons/17-architecture-decision-records
```

The tests verify ADR creation, status transitions (proposed, accepted, deprecated, superseded), manager CRUD operations, and filtering by status.

## Review questions

1. What are the four ADR lifecycle states, and when is each used?
2. Why should ADRs be stored in the repository rather than in a wiki or Notion?
3. What is the purpose of the `Alternatives Considered` section in an ADR?
4. How do you handle a situation where an accepted ADR is no longer the right decision?
5. Why is the ADR number permanent even after the ADR is deprecated?

## NEXT UP

Congratulations on completing Module 14! You now understand architecture patterns, distributed systems, event-driven design, multi-tenancy, and architecture decision records. Next up: Module 15 -- Deployment, CI/CD, DevOps.
