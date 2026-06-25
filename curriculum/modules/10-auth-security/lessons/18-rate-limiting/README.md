# Rate limiting

## Learning objective

Implement rate limiting in Go HTTP services using the token bucket algorithm, apply per-IP and per-user limits, propagate rate limit headers, and prevent resource exhaustion and abuse.

## Why this matters

Without rate limiting, a single client can monopolize server resources, exhaust database connections, run up cloud bills, or brute-force authentication endpoints. Rate limiting is a fundamental abuse-prevention mechanism used by every public API: GitHub (5000 req/hr), Stripe, Twitter, and all cloud providers. Go engineers must understand rate limiting to protect both the service and its users from denial-of-service and credential-stuffing attacks.

## Mental model

A token bucket is like a parking garage. The garage has a fixed number of spaces (tokens). Each arriving car (request) takes a space. Cars leave at a steady rate (refill), freeing spaces for new arrivals. If the garage is full, new cars are turned away.

Alternatively, a leaky bucket is like a funnel with a small drain. Water poured in at the top trickles out at a fixed rate. If you pour too fast, the funnel overflows.

The key insight: token buckets allow bursts (you can use all tokens at once) while capping the sustained rate. Leaky buckets enforce a strict processing rate but cannot handle bursts.

## Core idea

Rate limiting algorithms:

| Algorithm | Burst allowed | Behavior | Use case |
|---|---|---|---|
| Token bucket | Yes | Accumulates tokens up to a capacity, consumed per request, refilled at fixed rate | General-purpose, most common |
| Leaky bucket | No | Requests processed at fixed rate; excess queued or dropped | Strict traffic shaping |
| Fixed window | Yes (at boundary) | Counts requests per time window; resets at window boundary | Simple, but allows double-burst at boundaries |
| Sliding window | Yes | Counts requests in a rolling time window | More accurate than fixed window |
| Sliding window log | No | Tracks timestamp of each request in a log | Most accurate but memory-intensive |

Token bucket parameters:
- `capacity`: maximum number of tokens (burst size)
- `refillRate`: tokens added per second (sustained rate)

## Under the hood

The token bucket algorithm maintains a counter (`tokens`) and a timestamp of the last refill (`lastRefill`). On each request:

1. Calculate time elapsed since last refill: `elapsed = now - lastRefill`.
2. Add tokens: `tokens += elapsed * refillRate` (capped at `capacity`).
3. Update `lastRefill = now`.
4. If `tokens > 0`, decrement and allow. Otherwise, deny.

In production, rate limiting state is often stored in Redis using `INCR` + `EXPIRE` for the fixed window algorithm, or Redis Stack's `TDigest` and `TimeSeries` for sliding window. For distributed systems, the rate limiter must use a shared store (Redis, Memcached) rather than in-memory state.

The token bucket formula:

```text
tokens = min(capacity, tokens + (now - lastRefill) * refillRate)
```

## How Go uses it

Go's standard library does not include rate limiting middleware. The `golang.org/x/time/rate` package provides a token bucket implementation using the `rate.Limiter` type.

```go
import "golang.org/x/time/rate"

limiter := rate.NewLimiter(rate.Limit(10), 20) // 10 req/s, burst of 20
if limiter.Allow() {
    // handle request
} else {
    // rate limited
}
```

For HTTP middleware, rate limiters are typically keyed by IP address (for anonymous endpoints) or user ID (for authenticated endpoints). Headers should include `X-RateLimit-Limit` (capacity), `X-RateLimit-Remaining` (remaining tokens), and `Retry-After` (seconds until retry).

## Go example

```go
package main

import (
	"fmt"
	"sync"
	"time"
)

type TokenBucket struct {
	mu         sync.Mutex
	capacity   int
	tokens     int
	refillRate float64
	lastRefill time.Time
}

func NewTokenBucket(capacity int, refillRate float64) *TokenBucket {
	return &TokenBucket{
		capacity:   capacity,
		tokens:     capacity,
		refillRate: refillRate,
		lastRefill: time.Now(),
	}
}

func (tb *TokenBucket) Allow() bool {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(tb.lastRefill).Seconds()
	tb.lastRefill = now

	tb.tokens += int(elapsed * tb.refillRate)
	if tb.tokens > tb.capacity {
		tb.tokens = tb.capacity
	}

	if tb.tokens > 0 {
		tb.tokens--
		return true
	}
	return false
}

func main() {
	bucket := NewTokenBucket(3, 1.0)
	fmt.Printf("Capacity: %d, Refill: %.0f/s\n", bucket.capacity, bucket.refillRate)

	for i := 0; i < 10; i++ {
		allowed := bucket.Allow()
		fmt.Printf("Request %d: allowed=%v\n", i+1, allowed)
		if !allowed {
			break
		}
	}
}
```

## Step-by-step execution

For a token bucket with capacity=3 and refillRate=1/sec:

1. **t=0s**: Bucket starts with 3 tokens. Request 1 consumes 1 token (tokens=2). Allowed.
2. **t=0.1s**: Request 2 consumes 1 token (tokens=1). Allowed.
3. **t=0.2s**: Request 3 consumes 1 token (tokens=0). Allowed.
4. **t=0.3s**: Request 4: elapsed since last refill = 0.1s. Add 0.1 * 1 = 0 tokens (int truncation). Tokens=0. Denied.
5. **t=1.3s**: Request 5: elapsed since last refill = 1.0s. Add 1.0 * 1 = 1 token. Tokens=1. Consume it. Tokens=0. Allowed.
6. **t=2.3s**: Request 6: elapsed = 1.0s. Add 1 token. Tokens=1. Consume it. Allowed.

The pattern shows: first 3 requests are allowed immediately (burst), then requests are limited to approximately 1 per second (sustained rate).

## Common mistakes

- Mistake: Implementing per-endpoint rate limits without global limits.
  - Why it happens: Developers add rate limiters to each handler individually.
  - Fix: An attacker can rotate between endpoints to bypass per-endpoint limits. Always have a global limit per client.

- Mistake: Using IP-based rate limiting behind a reverse proxy.
  - Why it happens: The server reads `r.RemoteAddr` which is the proxy's IP, not the client's IP.
  - Fix: Use the `X-Forwarded-For` or `X-Real-IP` header. Validate and trust only proxy IPs from your infrastructure.

- Mistake: Rate limiting without a burst allowance.
  - Why it happens: Developers set capacity = rate, preventing any natural traffic spike.
  - Fix: Set capacity to 2-3x the sustained rate to allow short bursts. Legitimate clients have bursty traffic.

- Mistake: Applying rate limits after expensive operations.
  - Why it happens: The rate limiter check is placed deep in the handler, after database queries.
  - Fix: Apply rate limiting as early as possible in the middleware chain, before any I/O or computation.

- Mistake: Not propagating rate limit headers or `Retry-After`.
  - Why it happens: Developers assume clients will figure out when to retry.
  - Fix: Always set `X-RateLimit-Limit`, `X-RateLimit-Remaining`, and `Retry-After` headers. Clients need them for proper backoff.

## Debugging walkthrough

Consider this broken rate limiter:

```go
type FixedWindow struct {
	mu       sync.Mutex
	count    int
	limit    int
	window   time.Duration
	windowStart time.Time
}

func (fw *FixedWindow) Allow() bool {
	fw.mu.Lock()
	defer fw.mu.Unlock()

	if time.Since(fw.windowStart) > fw.window {
		fw.count = 0
		fw.windowStart = time.Now()
	}

	fw.count++
	return fw.count <= fw.limit
}
```

Symptom: A burst of 200 requests from the same IP at the exact second boundary all pass through, even though the limit is 100.

Investigation: Add logging for window boundaries:

```go
fmt.Printf("windowStart=%v, now=%v, elapsed=%v, count=%d\n",
    fw.windowStart, time.Now(), time.Since(fw.windowStart), fw.count)
```

Root cause: Fixed window resets at the start of each window. Two requests at the boundary (one at t=59.999, one at t=60.000) each see a fresh window with count=0. This is the "boundary burst" problem: the effective rate can be 2x the limit at window edges.

Fix: Use sliding window with sub-windows or a token bucket:

```go
// Use token bucket instead for smooth rate limiting
bucket := NewTokenBucket(100, 100.0/60.0) // 100 requests per 60 seconds
```

## Production notes

- Use Redis-based rate limiting for distributed services. In-memory rate limiting only works for single-instance deployments.
- For per-user rate limiting, extract user ID from the JWT or session, not from IP. Multiple users behind a NAT share an IP.
- Set different limits for different endpoints: login endpoints (lower limit), public read endpoints (higher limit).
- Always log rate limit violations. Monitor for patterns that indicate attacks.
- Use `golang.org/x/time/rate` for production token bucket implementations. It is maintained by the Go team.
- Expose rate limit metrics (Prometheus) for monitoring: `rate_limit_total`, `rate_limit_remaining`.

## Performance implications

- Token bucket operations are O(1) and complete in microseconds. The mutex contention is the main cost under high concurrency.
- For Redis-based rate limiting, each check adds 1-2ms of network latency. Batch or pipeline where possible.
- Fixed window is the least CPU-intensive (single counter + timestamp) but least accurate.
- Sliding window log is the most accurate but requires O(n) memory per client window.
- The 429 response should be lightweight: no database queries, no heavy serialization.

## Practice task

Implement a `SlidingWindow` rate limiter that counts requests in a rolling 60-second window. Use a slice of timestamps per key.

Type: `type SlidingWindow struct { mu sync.Mutex; windows map[string][]time.Time; limit int }`

Methods:
- `NewSlidingWindow(limit int) *SlidingWindow`
- `Allow(key string) bool` -- add current timestamp, remove timestamps older than 60s, check if count <= limit.

Then write a middleware that uses `SlidingWindow` and applies per-IP limits.

Write a `main()` that demonstrates: 5 rapid requests (first N allowed, rest denied), then wait and confirm allowance resumes.

## Tests / verification

```bash
go run ./curriculum/modules/10-auth-security/lessons/18-rate-limiting
go test ./curriculum/modules/10-auth-security/lessons/18-rate-limiting
```

The existing tests verify token bucket initial capacity, burst allowance up to capacity, refill behavior, middleware integration, and rate limit headers.

## Review questions

1. What is the difference between a token bucket and a leaky bucket algorithm?
2. Why does a fixed window rate limiter allow 2x traffic at window boundaries?
3. How would you implement per-user rate limiting in a JWT-authenticated Go API?
4. What headers should a rate limiter set on the response, and why is each important?
5. Why is IP-based rate limiting unreliable behind a reverse proxy or NAT?

## NEXT UP

TLS and HTTPS -- configuring TLS certificates, the TLS handshake process, and serving Go HTTP servers securely with HTTPS.
