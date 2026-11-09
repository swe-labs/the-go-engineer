# Answer Key — Module 16 Checkpoint: Portfolio & Interview Readiness

## Multiple Choice Answers

### Question 1: C
**Explanation:** Senior engineering evaluators and hiring managers look for verifiable engineering rigor. Anyone can clone a tutorial project, but an engineer who provides thorough architectural contracts, automated unit/integration test suites passing under `go test -race`, reproducible throughput and memory benchmarks, and production containerization artifacts demonstrates professional production readiness.

### Question 2: B
**Explanation:** In senior live coding interviews, jumping straight into code without clarifying ambiguous requirements is an immediate red flag indicating a lack of systems thinking. Clarifying read/write ratios, concurrency constraints, memory bounds, and failure modes shows architectural maturity and prevents costly refactoring mid-interview.

### Question 3: B
**Explanation:** For senior roles, correct functional output is only the baseline. Submissions fail primarily due to lack of production engineering hygiene: absence of tests, missing context cancellation, unhandled panics or ignored errors, data races detected under `go test -race`, and failure to document operational trade-offs and edge case handling.

### Question 4: B
**Explanation:** The "Result" in a STAR-formatted behavioral response must quantify engineering impact with concrete business and operational metrics (e.g., P99 latency reduction, error rate drop, cost reduction, MTTR decrease) and highlight lessons learned and organizational guardrails instituted.

### Question 5: B
**Explanation:** Production systems always involve trade-offs and operational challenges. Discussing architectural trade-offs, unexpected edge cases, and failure modes demonstrates intellectual humility, engineering depth, and authentic operational experience, setting the engineer apart from candidates who only understand happy-path tutorials.

### Question 6: B
**Explanation:** Initial verbal offers are typically anchored at the lower or middle end of a hiring band. Expressing excitement and appreciation while asking for the written details and a brief review window allows the candidate to benchmark against market compensation data and formulate a professional, value-driven counter-offer.

---

## Technical & Architectural Answers

### Question 7: System Design Walkthrough (Global Rate Limiter at 500k RPS)
1. **Requirements & Scale**:
   - 500,000 RPS globally across 4 regions (US-East, US-West, EU-West, AP-East).
   - Low latency: P99 < 5ms.
   - High availability: Rate limiting failures must fail open (do not block user traffic if rate limiter degrades).
2. **Architecture**:
   - Edge layer (Cloudflare / Envoy) with local token buckets for coarse-grained DDoS protection.
   - Regional Go Rate Limiting Microservice with local in-memory lock-striped ring buffers (Module 13 pattern) handling 99% of requests without network roundtrips.
   - Asynchronous batching: Local instances batch consumption counts every 100ms and sync with a regional Redis Cluster via pipelined Lua scripts or atomic increments.
   - Inter-region synchronization: Asynchronous gossip or CRDT-based counter synchronization for global quota enforcement.
3. **Failure Containment (Redis Unreachable)**:
   - Circuit breaker trips to "fail-open" mode: allow traffic through while logging warning metrics.
   - Fall back to purely local memory limits per instance, preventing cascading outages of upstream services.

### Question 8: Live Coding Interview Solution (Generic Concurrent Cache with TTL)

```go
package cache

import (
	"context"
	"sync"
	"time"
)

type cacheItem[V any] struct {
	value     V
	expiresAt time.Time
}

func (i cacheItem[V]) isExpired(now time.Time) bool {
	if i.expiresAt.IsZero() {
		return false
	}
	return now.After(i.expiresAt)
}

type Cache[K comparable, V any] struct {
	mu    sync.RWMutex
	items map[K]cacheItem[V]
}

func NewCache[K comparable, V any](ctx context.Context, cleanupInterval time.Duration) *Cache[K, V] {
	c := &Cache[K, V]{
		items: make(map[K]cacheItem[V]),
	}
	go c.cleanupLoop(ctx, cleanupInterval)
	return c
}

func (c *Cache[K, V]) Set(key K, val V, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	var exp time.Time
	if ttl > 0 {
		exp = time.Now().Add(ttl)
	}
	c.items[key] = cacheItem[V]{
		value:     val,
		expiresAt: exp,
	}
}

func (c *Cache[K, V]) Get(key K) (V, bool) {
	c.mu.RLock()
	item, ok := c.items[key]
	c.mu.RUnlock()

	if !ok {
		var zero V
		return zero, false
	}

	if item.isExpired(time.Now()) {
		// Clean up lazily on read
		c.mu.Lock()
		delete(c.items, key)
		c.mu.Unlock()

		var zero V
		return zero, false
	}

	return item.value, true
}

func (c *Cache[K, V]) cleanupLoop(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			now := time.Now()
			c.mu.Lock()
			for k, item := range c.items {
				if item.isExpired(now) {
					delete(c.items, k)
				}
			}
			c.mu.Unlock()
		}
	}
}
```

### Question 9: Incident Response & Postmortem Defense (STAR Method)
1. **Situation**: A production Go notification service handling 25,000 push events/sec started experiencing Kubernetes OOMKilled crashes during traffic spikes, causing a 12% drop in delivery success.
2. **Task**: As lead backend engineer, I was paged to halt the crash loop, isolate the memory leak, and restore SLA guarantees.
3. **Action**:
   - Retrieved a heap dump via `curl -o heap.pprof http://<pod>:6060/debug/pprof/heap` and analyzed with `go tool pprof -inuse_space`.
   - Identified that an unbuffered task channel was abandoned when client connections timed out; producer goroutines blocked waiting to send on the channel and could not be garbage collected, accumulating 80,000 leaked goroutines holding request buffers.
   - Deployed an emergency patch introducing a bounded worker pool with `select { case ch <- task: default: // drop or return backpressure }` and passed `ctx.Done()` to all producers.
4. **Result**:
   - Memory usage dropped from 4GB to 180MB; MTTR was 28 minutes with zero further OOM crashes.
   - Added automated goroutine leak detection tests (`uber-go/goleak`) to the CI pipeline to prevent future regressions.

### Question 10: Senior Offer Negotiation Strategy
1. **RSUs vs ISOs**:
   - RSUs (public companies): Actual shares granted with liquid monetary value upon vesting (subject to ordinary income tax upon vesting).
   - ISOs (startups): Options to purchase shares at a fixed strike price. Value is theoretical until a liquidity event (IPO or acquisition), and exercising can trigger Alternative Minimum Tax (AMT).
2. **Essential Startup Equity Questions**:
   - What is the total fully diluted share count, and what percentage do my options represent?
   - What was the most recent 409A valuation and preferred share price from the last funding round?
   - Are there liquidation preferences ahead of common stock (e.g., 1x non-participating)?
   - What is the post-termination exercise window (90 days vs 7–10 years)?
3. **Professional Counter-Offer Script**:
   > *"Thank you so much for extending this offer. I am genuinely excited about the team's technical mission and confident that my background in distributed Go architectures will make an immediate impact on platform reliability. Based on my current interview pipeline and market compensation data for senior backend roles requiring high-scale distributed systems expertise, I was targeting a total package with a base salary of $185,000. If we can adjust the base to $185,000 and increase the equity grant by 15%, I am ready to sign immediately and decline other ongoing interview processes."*
