package main

import (
	"fmt"
	"time"
)

type ADRStatus string

const (
	StatusProposed   ADRStatus = "proposed"
	StatusAccepted   ADRStatus = "accepted"
	StatusDeprecated ADRStatus = "deprecated"
	StatusSuperseded ADRStatus = "superseded"
)

type ArchitectureDecisionRecord struct {
	ID           string
	Title        string
	Status       ADRStatus
	Context      string
	Decision     string
	Consequences string
	Alternatives []string
	Pros         []string
	Cons         []string
	Date         time.Time
	Supersedes   string
}

func NewADR(id, title string) *ArchitectureDecisionRecord {
	return &ArchitectureDecisionRecord{
		ID:     id,
		Title:  title,
		Status: StatusProposed,
		Date:   time.Now(),
	}
}

func (a *ArchitectureDecisionRecord) Accept() {
	a.Status = StatusAccepted
}

func (a *ArchitectureDecisionRecord) Deprecate() {
	a.Status = StatusDeprecated
}

func (a *ArchitectureDecisionRecord) Supersede(newID string) {
	a.Status = StatusSuperseded
	a.Supersedes = newID
}

func (a *ArchitectureDecisionRecord) String() string {
	result := fmt.Sprintf("# %s - %s\n", a.ID, a.Title)
	result += fmt.Sprintf("Status: %s\n", a.Status)
	result += fmt.Sprintf("Date: %s\n\n", a.Date.Format("2006-01-02"))
	result += fmt.Sprintf("## Context\n%s\n\n", a.Context)
	result += fmt.Sprintf("## Decision\n%s\n\n", a.Decision)
	result += fmt.Sprintf("## Consequences\n%s\n\n", a.Consequences)
	if len(a.Alternatives) > 0 {
		result += "## Alternatives Considered\n"
		for _, alt := range a.Alternatives {
			result += fmt.Sprintf("- %s\n", alt)
		}
		result += "\n"
	}
	return result
}

type ADRManager struct {
	records map[string]*ArchitectureDecisionRecord
}

func NewADRManager() *ADRManager {
	return &ADRManager{records: make(map[string]*ArchitectureDecisionRecord)}
}

func (m *ADRManager) Add(adr *ArchitectureDecisionRecord) {
	m.records[adr.ID] = adr
}

func (m *ADRManager) Get(id string) *ArchitectureDecisionRecord {
	return m.records[id]
}

func (m *ADRManager) ByStatus(status ADRStatus) []*ArchitectureDecisionRecord {
	var result []*ArchitectureDecisionRecord
	for _, r := range m.records {
		if r.Status == status {
			result = append(result, r)
		}
	}
	return result
}

func (m *ADRManager) All() []*ArchitectureDecisionRecord {
	result := make([]*ArchitectureDecisionRecord, 0, len(m.records))
	for _, r := range m.records {
		result = append(result, r)
	}
	return result
}

func main() {
	manager := NewADRManager()

	adr1 := NewADR("ADR-001", "Use PostgreSQL for Primary Database")
	adr1.Context = "The team needs a primary database for the new order management system."
	adr1.Decision = "We will use PostgreSQL 16 as the primary database for all transactional workloads."
	adr1.Consequences = "PostgreSQL provides ACID compliance, full-text search, and JSONB support. The team has existing PostgreSQL expertise. We will use pgx as the Go driver."
	adr1.Alternatives = []string{"MySQL 8 - similar ACID compliance but weaker JSON support", "MongoDB - better scalability but eventual consistency is not suitable for payments", "SQLite - unsuitable for concurrent write workloads"}
	adr1.Accept()
	manager.Add(adr1)

	adr2 := NewADR("ADR-002", "Use Redis for Session Cache")
	adr2.Context = "The authentication service needs a fast, distributed cache for session tokens."
	adr2.Decision = "We will use Redis 7 with TTL-based key expiry. Session data is stored for 24 hours."
	adr2.Consequences = "Redis adds an operational dependency and requires a high-availability configuration with sentinel or cluster mode."
	adr2.Alternatives = []string{"In-memory Go map - not distributed, lost on restart", "Memcached - simpler but no data structures or TTL", "PostgreSQL - too slow for sub-millisecond session lookups"}
	adr2.Accept()
	manager.Add(adr2)

	adr3 := NewADR("ADR-003", "Use Kafka for Event Streaming")
	adr3.Context = "The order service needs to publish events for downstream consumers."
	adr3.Decision = "We will use Apache Kafka with Avro serialization and Schema Registry."
	adr3.Consequences = "Kafka provides at-least-once delivery guarantee. Consumers must be idempotent. Avro provides schema evolution."
	adr3.Alternatives = []string{"RabbitMQ - better for task queues, not event streaming", "AWS SQS - simpler but no ordered consumption", "gRPC streaming - no persistence or replay"}
	adr3.Accept()
	manager.Add(adr3)

	adr4 := NewADR("ADR-004", "Monolith First, Extract Services Later")
	adr4.Context = "The team is considering microservices architecture for the new platform."
	adr4.Decision = "We will start with a well-factored monolith. Services will only be extracted when the monolith causes measurable deployment pain."
	adr4.Consequences = "This avoids premature distributed complexity. We must enforce strict package boundaries within the monolith to make future extraction easier."
	adr4.Alternatives = []string{"Immediate microservices - adds overhead with no proven need", "Serverless - vendor lock-in and cold start latency"}
	adr4.Accept()
	manager.Add(adr4)

	for _, adr := range manager.All() {
		fmt.Println(adr.String())
	}

	fmt.Println("--- Accepted ADRs ---")
	for _, adr := range manager.ByStatus(StatusAccepted) {
		fmt.Printf("- %s: %s\n", adr.ID, adr.Title)
	}
}
