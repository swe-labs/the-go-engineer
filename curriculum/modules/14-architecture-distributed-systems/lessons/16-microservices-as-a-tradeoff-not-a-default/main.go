package main

import (
	"fmt"
	"math"
	"time"
)

type Tradeoff struct {
	Name       string
	Benefit    string
	Cost       string
	Severity   int
	Mitigation string
}

var MicroserviceTradeoffs = []Tradeoff{
	{
		Name:       "independent_deployability",
		Benefit:    "each service deploys independently",
		Cost:       "requires CI/CD pipelines per service, coordinated releases for cross-service changes",
		Severity:   2,
		Mitigation: "automated contract testing and backwards-compatible APIs",
	},
	{
		Name:       "team_autonomy",
		Benefit:    "teams own their services end-to-end",
		Cost:       "duplication of infrastructure, monitoring, and operational knowledge per team",
		Severity:   3,
		Mitigation: "shared platform team providing common infrastructure as a service",
	},
	{
		Name:       "scalability",
		Benefit:    "scale only the services that need it",
		Cost:       "network overhead per request, connection pool management",
		Severity:   2,
		Mitigation: "efficient serialization (protobuf), connection pooling, batching",
	},
	{
		Name:       "technology_flexibility",
		Benefit:    "each service can use the best language or database for its job",
		Cost:       "polyglot knowledge requirement, integration complexity, harder debugging",
		Severity:   4,
		Mitigation: "standardize on Go unless there is a compelling reason not to",
	},
	{
		Name:       "fault_isolation",
		Benefit:    "a crash in one service does not bring down others",
		Cost:       "distributed failure modes are harder to debug (partial failures, cascading failures)",
		Severity:   4,
		Mitigation: "circuit breakers, bulkheads, timeouts, and structured logging with tracing",
	},
	{
		Name:       "data_consistency",
		Benefit:    "each service owns its schema and can evolve independently",
		Cost:       "no distributed transactions across services; must use sagas or eventual consistency",
		Severity:   5,
		Mitigation: "idempotent operations, outbox pattern, saga orchestrator",
	},
}

type MicroserviceDecision struct {
	ServiceName        string
	TeamSize           int
	ExpectedThroughput int
	Pros               []string
	Cons               []string
	Recommendation     string
	Score              int
}

func EvaluateService(name string, teamSize int, throughput int, latencySensitive bool) MicroserviceDecision {
	dec := MicroserviceDecision{
		ServiceName:        name,
		TeamSize:           teamSize,
		ExpectedThroughput: throughput,
	}
	score := 0

	if teamSize > 8 {
		dec.Pros = append(dec.Pros, "team large enough to own a service")
		score += 10
	} else {
		dec.Cons = append(dec.Cons, "team too small to justify operational overhead")
		score -= 5
	}

	if throughput > 10000 {
		dec.Pros = append(dec.Pros, "high throughput justifies independent scaling")
		score += 10
	} else {
		dec.Cons = append(dec.Cons, "low throughput does not benefit from independent scaling")
		score -= 3
	}

	if !latencySensitive {
		dec.Pros = append(dec.Pros, "network latency is acceptable for this workload")
		score += 5
	} else {
		dec.Cons = append(dec.Cons, "network latency may impact latency-sensitive operations")
		score -= 8
	}

	dec.Score = score
	if score >= 10 {
		dec.Recommendation = "consider microservice"
	} else if score >= 0 {
		dec.Recommendation = "neutral -- could go either way"
	} else {
		dec.Recommendation = "start with monolith"
	}
	return dec
}

type ComplexityEstimate struct {
	MonolithMs     float64
	MicroserviceMs float64
	BreakEven      bool
}

func EstimateRequestLatency(calls int, avgServiceTimeMs float64, networkMs float64) ComplexityEstimate {
	mono := float64(calls) * avgServiceTimeMs
	micro := float64(calls) * (avgServiceTimeMs + networkMs)
	return ComplexityEstimate{
		MonolithMs:     math.Round(mono*100) / 100,
		MicroserviceMs: math.Round(micro*100) / 100,
		BreakEven:      micro < mono*1.5,
	}
}

func main() {
	dec := EvaluateService("Order Service", 3, 500, true)
	fmt.Printf("Service: %s\n", dec.ServiceName)
	fmt.Printf("Recommendation: %s (score=%d)\n", dec.Recommendation, dec.Score)
	for _, p := range dec.Pros {
		fmt.Printf("  + %s\n", p)
	}
	for _, c := range dec.Cons {
		fmt.Printf("  - %s\n", c)
	}

	est := EstimateRequestLatency(3, 50, time.Millisecond.Seconds()*5)
	fmt.Printf("\nMonolith latency: %.2fms\n", est.MonolithMs)
	fmt.Printf("Microservice latency: %.2fms\n", est.MicroserviceMs)
	fmt.Printf("Break-even acceptable: %v\n", est.BreakEven)

	fmt.Println("\nTop tradeoffs by severity:")
	for _, t := range MicroserviceTradeoffs {
		if t.Severity >= 4 {
			fmt.Printf("  %s (severity %d): %s\n", t.Name, t.Severity, t.Cost)
		}
	}
}
