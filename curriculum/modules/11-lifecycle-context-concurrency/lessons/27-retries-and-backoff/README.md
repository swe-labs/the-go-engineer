# Retries and backoff

## Learning objective

Implement exponential backoff with jitter, manage retry budgets, and make retries context-aware so they respect cancellation.

## Why this matters

Network calls, database queries, and RPCs fail transiently: a connection is reset, a server is momentarily overloaded, a DNS lookup times out. Retrying immediately after a failure is counterproductive — it piles more load onto an already strained system. Exponential backoff spreads retries over time, and jitter prevents all clients from retrying in lockstep (thundering herd). Every production Go service that talks to external systems needs a retry strategy.

## Mental model

Imagine you call a friend who is on another call. You wait a moment and call again. If they are still on the call, you wait longer. Each time you wait longer than the last time (exponential backoff). To avoid calling at the exact same time as other callers, you add a random delay (jitter). If you have a limited number of attempts (retry budget), you stop after your budget is exhausted. If your friend sends you a text saying "don't call back" (context cancellation), you stop.

## Core idea

Exponential backoff: after each failure, multiply the wait time by a constant factor (typically 2):

```
delay = baseDelay * 2^(attempt - 1)
```

Jitter: randomize the delay to desynchronize clients. Common strategies:

- **Full jitter**: `delay = random(0, delay)`
- **Equal jitter**: `delay = delay/2 + random(0, delay/2)`
- **Decorrelated jitter**: `delay = min(maxDelay, random(baseDelay, delay * 3))`

Retry budget: the maximum number of retry attempts or the maximum total time spent retrying.

Context awareness: every retry loop must check `ctx.Done()` between attempts so that cancellation stops retries immediately.

## Under the hood

A retry loop is a `for` loop with a `select` statement. After each failed attempt, the loop computes the sleep duration with exponential backoff and jitter, then does `select { case <-ctx.Done(): return ctx.Err(); case <-time.After(delay): }`. The `time.After` call creates a timer that is garbage collected after firing. In tight loops, cache the timer using `time.NewTimer` and call `timer.Stop()` on cancellation to avoid leaks.

## How Go uses it

The `net/http` transport has a default `Transport` with no retry. Higher-level packages like `aws-sdk-go`, `golang.org/x/oauth2`, and `github.com/hashicorp/go-retryablehttp` implement exponential backoff. The Google API Go client uses exponential backoff with jitter for all RPCs. The `grpc` library has built-in retry with exponential backoff configured via service config.

## Go example

```go
package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

type RetryConfig struct {
	MaxAttempts int
	BaseDelay   time.Duration
	MaxDelay    time.Duration
}

func retryWithBackoff(ctx context.Context, cfg RetryConfig, fn func(context.Context) error) error {
	var err error
	for attempt := 1; attempt <= cfg.MaxAttempts; attempt++ {
		err = fn(ctx)
		if err == nil {
			return nil
		}
		if attempt == cfg.MaxAttempts {
			return fmt.Errorf("all %d attempts failed: %w", cfg.MaxAttempts, err)
		}

		delay := cfg.BaseDelay * (1 << (attempt - 1))
		if delay > cfg.MaxDelay {
			delay = cfg.MaxDelay
		}
		jitter := time.Duration(rand.Int63n(int64(delay / 2)))
		delay = delay/2 + jitter

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(delay):
		}
	}
	return err
}

func unstableOp(ctx context.Context) error {
	if rand.Intn(100) < 70 {
		return fmt.Errorf("transient error")
	}
	return nil
}

func main() {
	ctx := context.Background()
	cfg := RetryConfig{
		MaxAttempts: 5,
		BaseDelay:   10 * time.Millisecond,
		MaxDelay:    1 * time.Second,
	}

	err := retryWithBackoff(ctx, cfg, unstableOp)
	if err != nil {
		fmt.Printf("Operation failed after retries: %v\n", err)
		return
	}
	fmt.Println("Operation succeeded")
}
```

## Step-by-step execution

1. First attempt: `unstableOp` is called. It fails (70% chance).
2. `attempt=1`, not last attempt. Compute delay: `10ms * 2^0 = 10ms`. Cap at 1s. Jitter: random up to 5ms. Total: `~5ms + jitter`.
3. `select` waits for the delay or context cancellation.
4. Second attempt: `unstableOp` called again. Fails.
5. `attempt=2`. Delay: `10ms * 2^1 = 20ms`. Jitter: up to 10ms. Total: `~10ms + jitter`.
6. Continues for 5 attempts. If all 5 fail, returns "all 5 attempts failed: transient error".

## Common mistakes

- Mistake: No jitter. Multiple clients retry at the same time.
  - Why it happens: Without jitter, every client with the same base delay retries at identical moments.
  - Fix: Always add jitter. Even a small random offset prevents thundering herds.

- Mistake: Retrying on errors that should not be retried (4xx, validation errors).
  - Why it happens: The retry loop treats all errors the same.
  - Fix: Classify errors into retryable (500, timeout, connection reset) and non-retryable (400, 404, 422).

- Mistake: Retrying without a context timeout.
  - Why it happens: The retry loop may take minutes to exhaust all attempts.
  - Fix: Pass a context with a timeout that covers the total retry budget.

- Mistake: Ignoring `ctx.Err()` after `fn` fails.
  - Why it happens: The function check `ctx.Done()` only between attempts, not after `fn` returns an error that may itself be a cancellation.
  - Fix: Check `ctx.Err()` after every failed attempt.

- Mistake: Using `time.Sleep` instead of `time.After` with `select`.
  - Why it happens: `time.Sleep` is not interruptible by context cancellation.
  - Fix: Always use `select` with `ctx.Done()` and `time.After`.

## Debugging walkthrough

Buggy program:

```go
func retryForever(fn func() error) {
	for {
		if err := fn(); err != nil {
			time.Sleep(1 * time.Second)
			continue
		}
		return
	}
}
```

Symptom: the function retries indefinitely if the operation keeps failing. No context cancellation, no backoff, no jitter.

Investigation: add a retry count and context parameter.

Fix: use `RetryConfig` with `MaxAttempts`, exponential backoff with jitter, and `select` on `ctx.Done()` as shown in the example.

## Production notes

Set retry budgets based on the service's SLO. A typical budget is 3 attempts for a 99.9% availability target. Use `context.WithTimeout` to set an overall deadline: `ctx, cancel := context.WithTimeout(parentCtx, 30*time.Second)`. For jitter, use `math/rand` seeded once at startup. The circuit breaker pattern complements retries: when failures exceed a threshold, the circuit breaker opens and rejects requests immediately without retrying. Log every retry attempt at `Debug` level with the attempt number, delay, and error message. Log the final failure at `Error` level with the total elapsed time. This makes it possible to tune retry parameters based on real traffic patterns without flooding the logs.

## Performance implications

Retries add latency proportional to the backoff delay: a 3-attempt retry with base 100ms adds ~700ms worst-case (100 + 200 + 400). Jitter increases variance but lowers peak load on the downstream system. The cost of a retry that succeeds is the sum of the failed attempts' latency plus the backoff delays. For high-throughput systems, avoid retrying on the hot path; use a background retry queue instead.

## Practice task

Implement a `DoWithRetry` function that accepts `ctx`, `maxAttempts`, `baseDelay`, and a function. The function should classify errors into retryable and non-retryable. Non-retryable errors should terminate immediately. Retryable errors should use full jitter (`delay = rand.Int63n(delay)`). Write table-driven tests that verify: success on first try, success after retry, exhaustion of retries, non-retryable error propagation, and context cancellation.

## Tests / verification

```bash
go run ./curriculum/modules/11-lifecycle-context-concurrency/lessons/27-retries-and-backoff
go test ./curriculum/modules/11-lifecycle-context-concurrency/lessons/27-retries-and-backoff
```

## Review questions

1. Why is exponential backoff better than fixed-interval retry?
2. What problem does jitter solve? Give a concrete scenario.
3. Why must a retry loop check `ctx.Done()` between attempts?
4. How do you decide which errors are retryable?
5. What is the relationship between retry budget, SLO, and deadline?

## NEXT UP

Idempotency keys
