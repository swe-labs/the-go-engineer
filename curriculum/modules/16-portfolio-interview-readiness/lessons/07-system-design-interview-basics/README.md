# System design interview basics

## Learning objective

Approach system design interviews with a structured framework: parse requirements, estimate capacity, identify tradeoffs, and communicate a coherent design under time pressure.

## Why this matters

System design interviews are the highest-weighted round for senior engineering roles (L5+ at top tech companies). Unlike coding interviews, there is no single correct answer — the interviewer evaluates your process: how you gather requirements, make tradeoffs, and communicate your reasoning. A structured approach to system design separates candidates who ramble from candidates who lead a productive design discussion.

## Mental model

Think of a system design interview as a collaborative whiteboarding session with a senior engineer who is your future teammate. Your goal is not to design the perfect system in 45 minutes — it is to demonstrate how you think about distributed systems problems. The interviewer wants to see that you can scope the problem, make reasonable assumptions, identify bottlenecks, and iterate on feedback. It is a conversation, not a presentation.

## Core idea

The system design interview follows a predictable four-phase structure:

**Phase 1: Requirements gathering (5 minutes)**
- Functional requirements: What should the system do? List 3-5 core features.
- Non-functional requirements: Latency, durability, availability, consistency, scalability.
- Traffic estimates: DAU, read/write ratio, peak vs. average load.
- Storage estimates: Data per entity, retention period, growth rate.

**Phase 2: High-level design (10 minutes)**
- Draw the main components: client, API gateway, application servers, databases, caches, queues.
- Define the data model: entities, relationships, primary keys.
- Choose communication patterns: synchronous (REST/gRPC) vs. asynchronous (queues/events).

**Phase 3: Deep dive (20 minutes)**
- Pick the most interesting component and go deep.
- Discuss tradeoffs: SQL vs. NoSQL, consistency model, partitioning strategy, caching strategy.
- Address the non-functional requirements: how do you meet the latency and availability targets?

**Phase 4: Wrap-up and bottlenecks (10 minutes)**
- Identify the top 3 bottlenecks in the current design.
- Propose specific improvements.
- Summarize the key decisions and tradeoffs.

**Common design questions**:
- Design a URL shortener (bit.ly)
- Design a chat system (WhatsApp)
- Design a news feed (Facebook)
- Design a rate limiter
- Design a distributed key-value store
- Design a web crawler

## Under the hood

The capacity calculator in this lesson translates vague requirements ("1 million daily active users") into concrete infrastructure numbers ("11,500 reads/second, 2.5 GB storage/day"). These calculations are based on a few key formulas:

**Reads per second**: `DAU * requests_per_user_per_day / 86400`
**Writes per second**: `Reads / read_write_ratio`
**Storage per year**: `records_per_day * record_size * 365`
**Peak traffic**: `average_traffic * 2` (general rule of thumb)

These estimates are intentionally rough. Interviewers accept order-of-magnitude accuracy (±2x). The goal is to demonstrate that you can reason about scale, not to compute exact numbers.

**Storage units reference**:

| Unit | Bytes | Approximate |
|---|---|---|
| 1 KB | 2^10 | One paragraph of text |
| 1 MB | 2^20 | A high-res photo |
| 1 GB | 2^30 | A full-length movie |
| 1 TB | 2^40 | 250,000 songs |

## How Go uses it

Go excels at building exactly the kind of systems that appear in design interviews:

- `net/http` for REST APIs and reverse proxies.
- `database/sql` for relational database access.
- `sync` package for concurrency control.
- Goroutines and channels for async processing.
- `context` for request-scoped cancellation and deadlines.

When discussing design tradeoffs, mention Go-specific considerations:

- Goroutines are lightweight (2KB stack) vs. OS threads (1MB+), enabling high concurrency with low overhead.
- Go channels provide built-in synchronization, reducing the need for external message queues in some designs.
- The standard library HTTP server handles connection pooling and keep-alive, simplifying proxy and gateway design.

## Go example

The requirement parser and capacity calculator translates system design requirements into storage and throughput estimates with assumptions and recommendations.

```go
package main

import (
	"fmt"
	"math"
	"strings"
)

type DataUnit int

const (
	Bytes DataUnit = iota
	Kilobytes
	Megabytes
	Gigabytes
	Terabytes
	RequestsPerSecond
	Records
)

func (d DataUnit) String() string {
	switch d {
	case Bytes:
		return "bytes"
	case Kilobytes:
		return "KB"
	case Megabytes:
		return "MB"
	case Gigabytes:
		return "GB"
	case Terabytes:
		return "TB"
	case RequestsPerSecond:
		return "req/s"
	case Records:
		return "records"
	default:
		return "units"
	}
}

type Requirement struct {
	Name        string
	Value       float64
	Unit        DataUnit
	Description string
}

type CapacityEstimate struct {
	Requirement     Requirement
	EstimatedValue  float64
	EstimatedUnit   DataUnit
	Assumptions     []string
	Recommendations []string
}

type SystemDesign struct {
	Name         string
	Requirements []Requirement
	Estimates    []CapacityEstimate
}

func (sd *SystemDesign) AddRequirement(name string, value float64, unit DataUnit, desc string) {
	sd.Requirements = append(sd.Requirements, Requirement{name, value, unit, desc})
}

func bytesToUnit(bytes float64) (float64, DataUnit) {
	switch {
	case bytes >= 1<<40:
		return bytes / (1 << 40), Terabytes
	case bytes >= 1<<30:
		return bytes / (1 << 30), Gigabytes
	case bytes >= 1<<20:
		return bytes / (1 << 20), Megabytes
	case bytes >= 1<<10:
		return bytes / (1 << 10), Kilobytes
	default:
		return bytes, Bytes
	}
}

func (sd *SystemDesign) EstimateStorage() {
	var totalBytes float64
	for _, req := range sd.Requirements {
		if req.Unit == Records || req.Unit == RequestsPerSecond {
			continue
		}
		switch req.Unit {
		case Bytes:
			totalBytes += req.Value
		case Kilobytes:
			totalBytes += req.Value * (1 << 10)
		case Megabytes:
			totalBytes += req.Value * (1 << 20)
		case Gigabytes:
			totalBytes += req.Value * (1 << 30)
		case Terabytes:
			totalBytes += req.Value * (1 << 40)
		}
	}
	estimated, unit := bytesToUnit(totalBytes)
	sd.Estimates = append(sd.Estimates, CapacityEstimate{
		Requirement:    Requirement{Name: "Total Storage", Unit: Bytes},
		EstimatedValue: estimated,
		EstimatedUnit:  unit,
		Assumptions:    []string{"assumes single-copy storage", "no compression applied"},
		Recommendations: []string{"add replication factor (3x for durability)",
			"consider compression (2-5x reduction)", "plan for 20% annual growth"},
	})
}

func (sd *SystemDesign) EstimateThroughput() {
	var totalRPS float64
	growth := 1.20
	for _, req := range sd.Requirements {
		if req.Unit == RequestsPerSecond {
			totalRPS += req.Value
		}
	}
	peakRPS := totalRPS * 2.0

	sd.Estimates = append(sd.Estimates, CapacityEstimate{
		Requirement:    Requirement{Name: "Throughput", Value: totalRPS, Unit: RequestsPerSecond},
		EstimatedValue: peakRPS,
		EstimatedUnit:  RequestsPerSecond,
		Assumptions:    []string{"peak traffic is 2x average", fmt.Sprintf("annual growth rate: %.0f%%", (growth-1)*100)},
		Recommendations: []string{"design for 2x peak headroom", "use connection pooling",
			"consider CDN for static content"},
	})
}

func (sd *SystemDesign) PrintReport() string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("System Design: %s\n", sd.Name))
	b.WriteString(strings.Repeat("-", 60) + "\n")

	b.WriteString("Requirements:\n")
	for _, req := range sd.Requirements {
		val := req.Value
		if val == math.Trunc(val) {
			b.WriteString(fmt.Sprintf("  - %s: %.0f %s\n", req.Name, val, req.Unit))
		} else {
			b.WriteString(fmt.Sprintf("  - %s: %.2f %s\n", req.Name, val, req.Unit))
		}
		if req.Description != "" {
			b.WriteString(fmt.Sprintf("    %s\n", req.Description))
		}
	}

	b.WriteString("\nCapacity Estimates:\n")
	for _, est := range sd.Estimates {
		b.WriteString(fmt.Sprintf("  - %s: %.0f %s\n", est.Requirement.Name, est.EstimatedValue, est.EstimatedUnit))
		b.WriteString("    Assumptions:\n")
		for _, a := range est.Assumptions {
			b.WriteString(fmt.Sprintf("      * %s\n", a))
		}
		b.WriteString("    Recommendations:\n")
		for _, r := range est.Recommendations {
			b.WriteString(fmt.Sprintf("      * %s\n", r))
		}
	}

	return b.String()
}

func main() {
	design := SystemDesign{Name: "URL Shortener"}

	design.AddRequirement("Daily Active Users", 1000000, RequestsPerSecond, "expected MAU after first year")
	design.AddRequirement("Write RPS", 115, RequestsPerSecond, "short URL creations")
	design.AddRequirement("Read RPS", 11500, RequestsPerSecond, "redirect requests (100:1 read/write ratio)")
	design.AddRequirement("URL Record Size", 512, Bytes, "average entry including metadata")
	design.AddRequirement("Total URL Records", 36500000*5, Records, "36.5M URLs/year x 5 years")

	design.EstimateStorage()
	design.EstimateThroughput()

	fmt.Print(design.PrintReport())
}
```

## Step-by-step execution

1. Create a `SystemDesign` with a name and add requirements using `AddRequirement`. Each requirement has a name, numeric value, unit, and description.
2. `EstimateStorage()` aggregates all byte-based requirements using the appropriate multiplier for each unit (KB, MB, GB, TB).
3. The total bytes are converted to the most readable unit using `bytesToUnit` (e.g., 5 TB instead of 5,497,558,138,880 bytes).
4. Storage assumptions and recommendations are added to help the engineer plan for real-world factors like replication and growth.
5. `EstimateThroughput()` calculates peak RPS at 2x average and notes the annual growth rate.
6. `PrintReport()` formats the complete analysis.

## Common mistakes

- Mistake: Jumping into the design without clarifying requirements. "Let's design a URL shortener" and immediately drawing boxes.
  - Why it happens: Engineers want to show technical depth early.
  - Fix: Spend the first 5 minutes clarifying: "What are the core features? How many users? What is the read/write ratio? What latency do we need?"

- Mistake: Using unrealistic numbers. Saying "10 billion requests per second" when the requirement is 1 million DAU.
  - Why it happens: Engineers underestimate the magnitude of large numbers.
  - Fix: Practice back-of-envelope calculations. 1 million DAU doing 10 requests/day = 10 million requests/day = ~115 requests/second. Always sanity-check.

- Mistake: Ignoring tradeoffs. Proposing a design without mentioning alternatives.
  - Why it happens: The design sounds correct to the author.
  - Fix: For every decision, mention one alternative and why you chose the current approach. "I chose PostgreSQL over MongoDB because our data is highly relational with joins."

- Mistake: Designing in a silo. Not engaging the interviewer or incorporating their feedback.
  - Why it happens: Interview pressure makes candidates go into delivery mode.
  - Fix: Pause every few minutes. Ask "Does that align with what you were thinking? Should I go deeper on any part?"

## Debugging walkthrough

An engineer designs a URL shortener and calculates:

```
- Storage: 92 GB (5 years of URLs at 512 bytes each)
- Throughput: 23,000 peak req/s
```

The interviewer asks: "How does this change if we add user accounts with analytics?" The engineer adds requirements for user profiles (2KB each), analytics events (256 bytes per click), and click tracking (10 clicks per URL on average). Re-running the calculator:

```
- Storage: 92 GB + 20 GB (users) + 470 GB (analytics) = ~582 GB
- Throughput: 23,000 + 1,150 (analytics writes) = 24,150 peak req/s
```

The storage jumps significantly due to analytics. The engineer realizes analytics storage is the dominant cost and proposes a separate time-series database for analytics events, reducing pressure on the primary key-value store.

## Production notes

In a real system design interview, the capacity calculator is a mental model, not a tool you run. Practice the calculations until they are automatic. Key benchmarks to memorize:

- 1 million requests/day ≈ 12 req/s
- 1 billion requests/day ≈ 11,500 req/s
- 1 TB can store ~2 billion 512-byte records
- A single PostgreSQL instance can handle ~10,000 read req/s or ~5,000 write req/s
- A single Redis instance can handle ~100,000 operations/s

These numbers help you make quick, credible estimates during the interview without a calculator.

## Performance implications

The capacity estimates inform performance-critical design decisions:

- Throughput estimates determine whether a single database instance suffices or sharding is needed.
- Storage estimates determine whether SSD, HDD, or cloud object storage is cost-effective.
- Peak traffic estimates determine autoscaling thresholds and buffer capacity.
- Read/write ratios determine caching strategy: read-heavy workloads benefit from aggressive caching; write-heavy workloads need write-behind queues.

Understanding these implications is the difference between a design that looks good on paper and one that works in production.

## Practice task

Add a `CacheSize` estimation method to `SystemDesign`. It should calculate the cache size needed to serve 80% of read requests from memory. Assume each cached record is the same size as a URL record (512 bytes), and the cache should hold the top 20% of URLs by access frequency. For 365M URLs, the working set is ~73M URLs. Calculate cache memory in GB and add it to the report with recommendations (e.g., Redis cluster size, eviction policy).

## Tests / verification

```bash
go test ./curriculum/modules/16-portfolio-interview-readiness/lessons/07-system-design-interview-basics
go run ./curriculum/modules/16-portfolio-interview-readiness/lessons/07-system-design-interview-basics
```

The tests verify basic storage estimation, multi-unit aggregation, throughput estimate assumptions, byte-to-unit conversion, and requirement storage. After the practice task, tests verify cache size calculation.

## Review questions

1. What are the four phases of a system design interview, and how long should each take?
2. Why is it important to clarify non-functional requirements before drawing the architecture?
3. Given 5 million daily active users making 20 requests per day, what is the approximate requests per second?
4. What tradeoff does SQL vs. NoSQL represent in a system design context?
5. Why should you state assumptions explicitly during a design interview?

## NEXT UP

Live coding interview practice — learn strategies for coding under pressure, verbalizing your thought process, and solving problems methodically in an interview setting.
