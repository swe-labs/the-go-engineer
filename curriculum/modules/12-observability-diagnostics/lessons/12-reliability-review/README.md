# Reliability review

## Learning objective

Build resilient Go services by applying reliability patterns: circuit breakers, bulkheads, graceful degradation, redundancy, and a comprehensive reliability checklist.

## Why this matters

Software fails. Networks partition. Databases time out. Dependencies crash. A service that assumes everything works will fail catastrophically when something breaks. The difference between a service that survives a dependency outage and one that cascades into a multi-service meltdown is deliberate reliability engineering. Netflix's Simian Army, Amazon's cell-based architecture, and Google's redundancy patterns all follow the same principles: assume failure, isolate faults, degrade gracefully, and never lose data.

## Mental model

Reliability patterns are like building fire safety into a skyscraper:

- **Circuit breaker**: a fire door that automatically closes when a fire is detected. Once closed, it stays closed until someone confirms the fire is out. This prevents the fire (errors) from spreading to other parts of the building (services).
- **Bulkhead**: a watertight compartment in a ship. If one compartment floods, the others stay dry. In software, this means isolating resources (connection pools, goroutines) per dependency so one failing dependency cannot exhaust resources for others.
- **Graceful degradation**: when the elevator is broken, people take the stairs. The building still functions, just with reduced capability. A service that returns stale cached data when the database is down is better than a service that returns HTTP 500.
- **Redundancy**: having two fire escapes instead of one. Multiple replicas, multiple availability zones, multiple data centers.

The analogy breaks because software is not physical — bulkheads can be resized at runtime, circuit breakers can be half-open (testing if the fire is out), and redundancy is limited by cost and consistency constraints.

## Core idea

Reliability patterns work together to create systems that survive failures:

**Circuit Breaker** (states: Closed, Open, Half-Open):
- **Closed**: normal operation. Requests pass through. Failures are counted.
- **Open**: failure threshold exceeded. Requests fail fast (without calling the dependency).
- **Half-Open**: after a timeout, a probe request is allowed to test if the dependency has recovered.
- Success moves to Closed; failure returns to Open.

**Bulkhead**:
- Allocate separate resource pools (goroutines, connections, memory) for each dependency.
- If one pool is exhausted, other dependencies continue unaffected.

**Graceful Degradation**:
- When a dependency fails, return a degraded response instead of an error.
- Examples: stale cache data, simplified UI, read-only mode.

**Redundancy**:
- Run multiple replicas behind a load balancer.
- Deploy across availability zones.
- Use leader election for stateful services.

## Under the hood

A circuit breaker wraps each external call. The implementation tracks failures as a sliding window counter (e.g., last 60 seconds). When the failure rate exceeds a threshold (e.g., 50% of requests), the breaker opens. All subsequent calls return an error immediately without making the network call — typically within microseconds instead of seconds (the timeout duration). After a configurable cooldown (e.g., 30 seconds), the breaker transitions to half-open and allows one probe request. If the probe succeeds, the breaker closes. If it fails, the cooldown resets.

The `golang.org/x/sync/semaphore` package provides weighted semaphores useful for bulkhead implementation. A bulkhead acquires a semaphore before making a call; if the semaphore is full, the call fails fast instead of queuing indefinitely.

## How Go uses it

Go's standard library and ecosystem provide building blocks for reliability:

- `net/http.Transport` has `MaxIdleConns` and `MaxIdleConnsPerHost` for connection pooling.
- `database/sql` has `SetMaxOpenConns`, `SetMaxIdleConns`, and `SetConnMaxLifetime` for DB connection bulkheads.
- `context.Context` with `context.WithTimeout` and `context.WithDeadline` for request-level timeouts.
- `golang.org/x/sync/errgroup` for goroutine lifecycle management with cancellation.
- `github.com/sony/gobreaker` is the most popular circuit breaker library for Go.

For production, many teams build lightweight circuit breakers inline rather than adding a library dependency.

## Go example

```go
package main

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

type State int

const (
	StateClosed   State = 0
	StateOpen     State = 1
	StateHalfOpen State = 2
)

type CircuitBreaker struct {
	mu               sync.Mutex
	state            State
	failureCount     int
	successCount     int
	threshold        int
	halfOpenMax      int
	cooldown         time.Duration
	lastFailureTime  time.Time
}

func NewCircuitBreaker(threshold int, cooldown time.Duration) *CircuitBreaker {
	return &CircuitBreaker{
		state:       StateClosed,
		threshold:   threshold,
		halfOpenMax: 3,
		cooldown:    cooldown,
	}
}

func (cb *CircuitBreaker) Call(ctx context.Context, fn func(context.Context) error) error {
	cb.mu.Lock()

	switch cb.state {
	case StateOpen:
		if time.Since(cb.lastFailureTime) < cb.cooldown {
			cb.mu.Unlock()
			return errors.New("circuit breaker: open")
		}
		cb.state = StateHalfOpen
		cb.successCount = 0

	case StateHalfOpen:
		if cb.successCount >= cb.halfOpenMax {
			cb.state = StateClosed
			cb.failureCount = 0
			cb.successCount = 0
		}
	}

	cb.mu.Unlock()

	err := fn(ctx)

	cb.mu.Lock()
	defer cb.mu.Unlock()

	if err != nil {
		cb.failureCount++
		cb.lastFailureTime = time.Now()
		if cb.failureCount >= cb.threshold {
			cb.state = StateOpen
		}
		return err
	}

	cb.successCount++
	if cb.state == StateHalfOpen {
		cb.successCount++
	}
	return nil
}

type Bulkhead struct {
	sem chan struct{}
}

func NewBulkhead(maxConcurrent int) *Bulkhead {
	return &Bulkhead{sem: make(chan struct{}, maxConcurrent)}
}

func (b *Bulkhead) Execute(ctx context.Context, fn func(context.Context) error) error {
	select {
	case b.sem <- struct{}{}:
		defer func() { <-b.sem }()
		return fn(ctx)
	case <-ctx.Done():
		return ctx.Err()
	}
}

type DegradationFunc func(ctx context.Context) (string, error)

type ResilientClient struct {
	primary   DegradationFunc
	fallback  DegradationFunc
	breaker   *CircuitBreaker
	bulkhead  *Bulkhead
}

func NewResilientClient(primary, fallback DegradationFunc) *ResilientClient {
	return &ResilientClient{
		primary:  primary,
		fallback: fallback,
		breaker:  NewCircuitBreaker(5, 10*time.Second),
		bulkhead: NewBulkhead(10),
	}
}

func (c *ResilientClient) Fetch(ctx context.Context) (string, error) {
	result, err := c.breaker.Call(ctx, func(ctx context.Context) error {
		res, err := c.bulkhead.Execute(ctx, func(ctx context.Context) error {
			r, e := c.primary(ctx)
			result = r
			return e
		})
		return err
	})

	if err != nil {
		fallbackResult, fallbackErr := c.fallback(ctx)
		if fallbackErr != nil {
			return "", fmt.Errorf("primary and fallback both failed: %w", fallbackErr)
		}
		return fallbackResult + " (degraded)", nil
	}

	return result, nil
}

func main() {
	primaryCallCount := 0
	client := NewResilientClient(
		func(ctx context.Context) (string, error) {
			primaryCallCount++
			if primaryCallCount <= 3 {
				return "", errors.New("primary failure")
			}
			return "primary data", nil
		},
		func(ctx context.Context) (string, error) {
			return "cached data", nil
		},
	)

	for i := 0; i < 8; i++ {
		result, err := client.Fetch(context.Background())
		if err != nil {
			fmt.Printf("Request %d: error — %v\n", i+1, err)
		} else {
			fmt.Printf("Request %d: %s\n", i+1, result)
		}
	}

	fmt.Println("\nReliability patterns applied:")
	fmt.Println("  - Circuit breaker: stops cascading failures")
	fmt.Println("  - Bulkhead: isolates dependency resource usage")
	fmt.Println("  - Degradation: falls back to cached data")
}
```

## Step-by-step execution

1. `CircuitBreaker` starts in `Closed` state. Every failure increments the counter.
2. After `threshold` failures (5), breaker transitions to `Open` — all calls fail fast with "circuit breaker: open".
3. After `cooldown` (10s), breaker transitions to `HalfOpen` — allows probe requests.
4. If `halfOpenMax` (3) consecutive probes succeed, breaker transitions back to `Closed`.
5. `Bulkhead` limits concurrent calls to the primary dependency to `maxConcurrent` (10).
6. `ResilientClient.Fetch` calls through breaker → bulkhead → primary.
7. If primary fails, the fallback function returns cached data with "(degraded)" suffix.
8. If both primary and fallback fail, the error propagates to the caller.

## Common mistakes

- Mistake: Setting the circuit breaker threshold too low, causing it to open during normal traffic spikes.
  - Why it happens: A 3-second timeout on a 100 req/s service means even 1 slow request creates 100 failures in the window.
  - Fix: Set threshold based on the service's normal error rate plus a safety margin. Use a sliding window rather than absolute count.

- Mistake: Using circuit breakers without bulkheads.
  - Why it happens: The circuit breaker prevents calls to a failing dependency, but the connection pool from previous calls may be exhausted, and the breaker's probe requests also use connections.
  - Fix: Always pair circuit breakers with bulkheads. The bulkhead limits resource usage per dependency.

- Mistake: No degradation path — when the circuit breaker opens, the service returns HTTP 503 to all callers.
  - Why it happens: The team only implemented the circuit breaker without planning what to serve when it opens.
  - Fix: Define a degradation strategy for every dependency: stale cache, default values, or simplified responses.

## Debugging walkthrough

A service suddenly starts returning HTTP 503 errors to all clients. The circuit breaker was implemented last month.

**Symptom**: All requests to the service return "circuit breaker: open" after the first 5 failures.

**Investigation**:
1. Check the circuit breaker's failure count and threshold.
2. Check the dependency (database, external API) — is it healthy? Yes, the database is responding normally.
3. Check the circuit breaker threshold: 5 failures. Check the deployment log: a new version was deployed 10 minutes ago with a bug that causes all database queries to time out.
4. The bug was in a different code path, but every request hits the database at least once, so every request fails and opens the breaker.

**Root cause**: A single bug in the database access layer caused all queries to time out. The circuit breaker opened after 5 failures, preventing all subsequent requests (even those that did not need the database).

**Fix**:
1. Roll back the buggy deployment immediately.
2. Split the circuit breaker into per-endpoint breakers so a bug in one code path does not block unrelated endpoints.
3. Add a health check endpoint that bypasses the circuit breaker for monitoring.

## Production notes

- Circuit breakers should be per-dependency, not global. Each external service, database, and queue should have its own breaker.
- Set circuit breaker timeouts: open timeout (how long to wait before half-open) should be longer than the dependency's typical recovery time. For a database failover, use 30-60 seconds. For a flaky API, use 5-10 seconds.
- Bulkhead sizes: set `maxConcurrent` to the dependency's connection pool size. For a database with `MaxOpenConns=25`, the bulkhead should also be 25.
- Monitor circuit breaker state as a Prometheus gauge (`circuit_breaker_state{service="orders"}`). Alert on state transitions from closed to open.
- Test reliability patterns with chaos engineering: use tools like `toxiproxy` to inject latency and failures in staging, then verify the circuit breaker opens and degradation works correctly.

## Performance implications

- Circuit breaker in closed state: adds ~100ns per call (mutex lock + counter increment). In open state: returns immediately (~100ns) instead of waiting for a timeout (seconds).
- Bulkhead: acquiring/releasing a semaphore adds ~200ns. If the semaphore is full, the caller fails fast (~200ns) instead of queuing indefinitely.
- Graceful degradation: returning cached data is typically faster than the primary call (microseconds vs milliseconds).
- The cost of not having these patterns: a cascading failure across 10 services can cause hours of downtime, costing orders of magnitude more than the pattern overhead.

## Practice task

Implement a `MultiBreaker` that manages per-dependency circuit breakers:
1. `Register(name string, threshold int, cooldown time.Duration)` — registers a named breaker.
2. `Call(name string, ctx context.Context, fn func(context.Context) error)` — calls through the named breaker.
3. `State(name string) State` — returns the current state of a breaker.
4. `Metrics() map[string]BreakerMetrics` — returns failure/success counts for all breakers.

Write tests that verify:
- A breaker opens after threshold failures.
- A breaker closes after successful probes in half-open state.
- Different breakers operate independently (one failing does not affect another).

## Tests / verification

```bash
go test ./curriculum/modules/12-observability-diagnostics/lessons/12-reliability-review -v
go run ./curriculum/modules/12-observability-diagnostics/lessons/12-reliability-review
```

## Review questions

1. What are the three states of a circuit breaker and what triggers each transition?
2. Why does a bulkhead pattern need a circuit breaker, and vice versa?
3. What is the difference between graceful degradation and failing fast? When would you choose each?
4. How does redundancy improve reliability but also increase operational cost?
5. A service calls three external APIs. One API starts returning 5-second timeouts. Without any reliability patterns, what happens to the service's own clients?

## NEXT UP

Congratulations on completing Module 12! You now understand observability, metrics, logging, tracing, and reliability patterns for production Go services. Next up: Module 13 — Deployment, CI/CD, DevOps.
