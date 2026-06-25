# Microservices as a tradeoff, not a default

## Learning objective

Evaluate the concrete tradeoffs of microservices vs monoliths for a given system, calculate the latency cost of distributed calls, and apply a decision framework that recommends microservices only when the benefits outweigh the costs.

## Why this matters

Microservices are the default architecture in many organizations, not because they are the right choice, but because they are the fashionable choice. The industry has numerous cautionary tales: startups that spent 6 months splitting a monolith that worked perfectly well, enterprises with 200 microservices that nobody understands, and teams drowning in deployment complexity. Microservices are a tool, not a goal. Understanding the tradeoffs allows you to choose the right architecture for your context, defend the decision with data, and avoid paying the distributed complexity tax unnecessarily.

## Mental model

A monolith is a bicycle. Simple, cheap, one person can fix it in 5 minutes. A microservice architecture is a fleet of Formula 1 cars. Fast, specialized, requires a pit crew of engineers, mechanics, and strategists. If you are delivering pizza, use the bicycle. If you are racing at 300 km/h, use the F1 car. Using an F1 car to deliver pizza is absurdly expensive. Using a bicycle in an F1 race is suicide.

Microservices add overhead to every operation: every function call becomes a network call, every database query becomes an API call, every transaction becomes a saga. The overhead is only justified when the benefits -- independent deployability, team autonomy, independent scaling -- are needed and not achievable with a well-factored monolith.

## Core idea

The tradeoffs of microservices:

| Benefit | Cost |
|---|---|
| Independent deployability | Multiple CI/CD pipelines, release coordination |
| Team autonomy | Duplicated infrastructure and operational knowledge |
| Independent scaling | Network overhead per request (1-10ms added) |
| Technology flexibility | Polyglot complexity, integration overhead |
| Fault isolation | Distributed failure modes (partial failures, cascading) |
| Data ownership | Eventual consistency, saga complexity |

The decision framework: start with a monolith. Split when the monolith causes measurable pain:

- Deploy time exceeds 30 minutes.
- A single change requires coordination between 3+ teams.
- One part of the system needs dramatically different scaling.
- The team size exceeds 8-10 developers per codebase.

## Under the hood

Every cross-service call adds costs that are invisible in a monolith:

- **Network latency**: 1-10ms per call (vs 0.001ms for in-process call).
- **Serialization/deserialization**: JSON marshal/unmarshal adds 10-50us.
- **Connection management**: HTTP connection pools, keep-alive, retries.
- **Error handling**: network timeouts, 5xx responses, retry logic.
- **Observability**: distributed tracing, span propagation, log aggregation.

A monolith that makes 3 database calls in 50ms becomes a microservice that makes 3 HTTP calls in 65ms (50ms + 3 * 5ms network). If those 3 calls fan out to 10 more calls each, the latency multiplies. This is the n+1 problem of distributed systems.

## How Go uses it

- **Standard library is monolith-friendly**: `net/http`, `database/sql`, `html/template` all work out of the box in a single process. No special infrastructure needed.
- **`internal` packages**: Go's `internal` package visibility enforces module boundaries within a monorepo. This mirrors bounded context isolation without network calls.
- **`net/http/httptest`**: integration tests for a monolith are trivial. For microservices, you need Docker containers, service mocks, or test harnesses.
- **`protobuf` / `gRPC`**: if you do split, gRPC provides strongly typed contracts and efficient binary serialization. Protobuf is 3-5x faster than JSON.

## Go example

```go
package main

import (
	"fmt"
	"math"
)

type Tradeoff struct {
	Name      string
	Benefit   string
	Cost      string
	Severity  int
}

var MicroserviceTradeoffs = []Tradeoff{
	{"deployability", "independent deploys", "N CI/CD pipelines", 2},
	{"autonomy", "team ownership", "duplicated ops knowledge", 3},
	{"scalability", "scale per service", "network overhead", 4},
	{"flexibility", "polyglot freedom", "integration complexity", 4},
	{"consistency", "local consistency per service", "eventual consistency across services", 5},
}

type Decision struct {
	Name           string
	TeamSize       int
	Throughput     int
	Recommendation string
	Score          int
}

func EvaluateService(name string, teamSize, throughput int, latencySensitive bool) Decision {
	d := Decision{Name: name, TeamSize: teamSize, Throughput: throughput}
	score := 0
	if teamSize > 8 {
		d.Recommendation = "consider microservice"
		score += 10
	} else {
		score -= 5
	}
	if throughput > 10000 {
		score += 10
	} else {
		score -= 3
	}
	if !latencySensitive {
		score += 5
	} else {
		score -= 8
	}
	d.Score = score
	if score >= 10 {
		d.Recommendation = "consider microservice"
	} else if score >= 0 {
		d.Recommendation = "neutral -- evaluate further"
	} else {
		d.Recommendation = "start with monolith"
	}
	return d
}

type LatencyEstimate struct {
	MonolithMs float64
	SplitMs    float64
}

func EstimateLatency(calls int, svcMs, netMs float64) LatencyEstimate {
	return LatencyEstimate{
		MonolithMs: math.Round(float64(calls)*svcMs*100) / 100,
		SplitMs:    math.Round(float64(calls)*(svcMs+netMs)*100) / 100,
	}
}

func main() {
	dec := EvaluateService("Payment Service", 3, 500, true)
	fmt.Printf("Service: %s\n", dec.Name)
	fmt.Printf("Recommendation: %s (score=%d)\n", dec.Recommendation, dec.Score)

	est := EstimateLatency(3, 50, 5)
	fmt.Printf("\nMonolith: %.2fms | Split: %.2fms\n", est.MonolithMs, est.SplitMs)

	fmt.Println("\nKey tradeoffs:")
	for _, t := range MicroserviceTradeoffs {
		fmt.Printf("  - %s (severity %d): %s | %s\n", t.Name, t.Severity, t.Benefit, t.Cost)
	}
}
```

## Step-by-step execution

For `EvaluateService("Payment Service", 3, 500, true)`:

1. `teamSize=3` → not > 8 → score -= 5.
2. `throughput=500` → not > 10000 → score -= 3.
3. `latencySensitive=true` → score -= 8.
4. Total score: -16.
5. Score < 0 → `Recommendation = "start with monolith"`.

For `EvaluateService("Order Service", 12, 50000, false)`:

1. `teamSize=12` > 8 → score += 10.
2. `throughput=50000` > 10000 → score += 10.
3. `latencySensitive=false` → score += 5.
4. Total score: 25.
5. Score >= 10 → `Recommendation = "consider microservice"`.

## Common mistakes

- **Splitting for scalability too early**: 100 req/s does not need microservices. A monolith with a read replica handles that easily. Wait until you need to scale different parts independently.
- **Not counting operational overhead**: each microservice needs monitoring, logging, alerting, CI/CD, secrets management, and on-call rotation. A team of 4 cannot operate 10 services.
- **Forgetting about data consistency**: distributed transactions do not exist. Your atomic database transaction becomes a saga with compensating actions. This is a fundamental complexity increase.
- **Assuming network is reliable**: in a monolith, function calls never timeout. In microservices, every call can fail, timeout, or return garbage. Your code must handle all of these.
- **Splitting because everyone else does it**: the decision must be based on your context -- team size, throughput, latency requirements, data consistency needs. Not on industry trends.

## Debugging walkthrough

A team of 5 developers maintains 12 microservices. Deployments take 2 hours. A bug fix that would take 1 hour in a monolith takes 3 days across 4 services.

**Symptom**: The team is overwhelmed by operational overhead. They spend more time on CI/CD, monitoring, and deployment coordination than on feature development.

**Investigation**:
1. Count services vs team size: 12 services / 5 developers = 2.4 services per developer.
2. Measure deployment time: 2 hours per full release.
3. Measure cross-service changes: 40% of features touch 3+ services.

**Root cause**: The architecture was chosen before the team experienced monolith pain. The team does not have the capacity to operate 12 services.

**Fix**: Merge related services back into a monolith. Identify services that are always deployed together (e.g., user + auth, order + payment). Merge them into a single deployable unit. Reduce services to 4. Deployment time drops to 15 minutes. The team regains velocity.

## Production notes

- **Start with a modular monolith**: enforce package boundaries, separate concerns, use interfaces. If the monolith needs to be split later, the boundaries are already defined.
- **Extract only when painful**: use the "pain threshold" approach. Do not extract a service until the monolith's current architecture is causing measurable friction.
- **Shared platform team**: if the organization has 5+ microservices, invest in a platform team that provides common infrastructure (service mesh, observability, CI/CD).
- **Strangler fig for extraction**: when you do split, use the strangler fig pattern. Never do a big-bang rewrite.
- **Cost tracking**: track the cost per service (infrastructure, operations, developer time). If a service costs more to operate than the value it provides, consider merging it back.

## Performance implications

- **Network latency**: 1-10ms per cross-service call. A single user request that used to make 5 function calls (5us) now makes 5 API calls (25ms total added).
- **Serialization**: JSON adds 10-50us per call. Protobuf adds 1-5us. For 1000 req/s, this CPU adds up.
- **Connection overhead**: HTTP/2 multiplexing reduces connection overhead, but each connection uses memory (TCP buffer, goroutine). 10 services × 100 connections = overhead.
- **Memory**: each service has its own runtime, garbage collector, and connection pools. 10 Go services = 10 GCs running, each with its own heap.
- **Operational cost**: the real cost of microservices is not runtime but operations: monitoring, alerting, deploy pipelines, on-call rotation.

## Practice task

Write a function `RecommendArchitecture(teamSize, throughput, latencySensitive bool, deploysPerWeek int) string` that scores the tradeoffs and returns either "monolith" or "microservice". Use the scoring rubric above. Then call it with three different scenarios and print the recommendations.

## Tests / verification

```bash
go run ./curriculum/modules/14-architecture-distributed-systems/lessons/16-microservices-as-a-tradeoff-not-a-default
go test ./curriculum/modules/14-architecture-distributed-systems/lessons/16-microservices-as-a-tradeoff-not-a-default
```

The tests verify service evaluation for different team sizes and throughput levels, latency estimates, and the tradeoffs list.

## Review questions

1. What concrete problem does a microservice architecture solve that a well-factored monolith does not?
2. How much network latency does each cross-service HTTP call add to a request?
3. Why is data consistency harder in a microservice architecture than in a monolith?
4. What is the recommended team-to-service ratio for a microservice architecture?
5. When would you recommend a monolith over microservices for a new project?

## NEXT UP

Architecture decision records.
