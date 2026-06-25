# Idempotency keys

## Learning objective

Design an idempotency-key-based deduplication system for HTTP APIs, implement concurrent-safe dedup storage with a mutex, and ensure retry safety.

## Why this matters

When a client retries a payment or order creation, the server must process the request exactly once. Without idempotency, a retry might charge the customer twice, create duplicate orders, or corrupt state. Idempotency keys are the standard solution: the client sends a unique key with each request, and the server uses the key to detect and reject duplicates. Every payment API (Stripe, Square, PayPal) and most modern REST APIs enforce idempotency on mutation endpoints.

## Mental model

An idempotency key is like a locker at a train station. The client puts their luggage in a locker (first request), and the station gives them a key. If the client returns with the same key (retry request), the station opens the same locker and shows them the same luggage — no duplicate luggage is created. The locker is the dedup store, the key is the idempotency key, and the luggage is the result.

## Core idea

Idempotency means applying an operation multiple times produces the same result as applying it once. For HTTP APIs:

1. The client generates a unique key (typically a UUID v4) and sends it in the `Idempotency-Key` header.
2. The server checks if the key exists in the dedup store.
3. If it does not exist: the server processes the request, stores the key and result atomically.
4. If it exists: the server returns the stored result without processing again.

The dedup store must be thread-safe (use a mutex for in-memory storage, or atomic upsert for Redis/PostgreSQL). The key observation period must exceed the maximum retry interval.

## Under the hood

A dedup store is logically a `map[string]StoredResponse` protected by a `sync.Mutex` or `sync.RWMutex`. On concurrent requests with the same key, exactly one request wins the race to insert, and all others see the stored result. This is a compare-and-swap pattern: "insert if not exists". In Redis this is `SETNX`. In PostgreSQL it is `INSERT ... ON CONFLICT DO NOTHING` or `INSERT ... ON CONFLICT DO UPDATE`. The mutex ensures that two concurrent requests with the same key do not both execute the operation.

## How Go uses it

Stripe's API client libraries generate idempotency keys automatically. The `net/http` package does not implement idempotency at the transport layer; it delegates to the application. Many Go API frameworks (chi, gin) provide middleware for automatic idempotency key handling. The Go `google.golang.org/api` client uses idempotency tokens for mutating Google APIs.

## Go example

```go
package main

import (
	"fmt"
	"sync"
)

type IdempotencyStore struct {
	mu    sync.Mutex
	store map[string]string
}

func NewIdempotencyStore() *IdempotencyStore {
	return &IdempotencyStore{store: make(map[string]string)}
}

func (s *IdempotencyStore) Get(key string) (string, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	val, ok := s.store[key]
	return val, ok
}

func (s *IdempotencyStore) Set(key, value string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.store[key]; exists {
		return false
	}
	s.store[key] = value
	return true
}

type PaymentRequest struct {
	ID       string
	Amount   int
	Currency string
}

func processPayment(store *IdempotencyStore, req PaymentRequest) (string, bool) {
	if result, ok := store.Get(req.ID); ok {
		return result, true
	}
	result := fmt.Sprintf("charged %d %s", req.Amount, req.Currency)
	if store.Set(req.ID, result) {
		return result, false
	}
	if result, ok := store.Get(req.ID); ok {
		return result, true
	}
	return "", false
}

func main() {
	store := NewIdempotencyStore()
	req := PaymentRequest{ID: "pay_123", Amount: 5000, Currency: "USD"}

	result1, dedup1 := processPayment(store, req)
	fmt.Printf("First call: %s (dedup=%v)\n", result1, dedup1)

	result2, dedup2 := processPayment(store, req)
	fmt.Printf("Second call: %s (dedup=%v)\n", result2, dedup2)
}
```

## Step-by-step execution

1. Client calls `processPayment` with `ID: "pay_123"`. Store does not have "pay_123".
2. `store.Get("pay_123")` returns `("", false)`.
3. The payment is processed, producing `"charged 5000 USD"`.
4. `store.Set("pay_123", "charged 5000 USD")` acquires the mutex, checks the map (still empty), inserts the key/value, returns `true`.
5. Client calls `processPayment` again with `ID: "pay_123"` (retry).
6. `store.Get("pay_123")` returns `("charged 5000 USD", true)`.
7. The function returns `("charged 5000 USD", true)` without processing a second payment.
8. The second `dedup` flag tells the server this was a duplicate.

## Common mistakes

- Mistake: Storing only the fact that a key was seen, not the result.
  - Why it happens: The store records "key present" but deletes or overwrites the result on partial failure.
  - Fix: Store the complete response. When a retry arrives, return the stored response.

- Mistake: Using `sync.RWMutex` and checking + inserting in two separate lock acquisitions.
  - Why it happens: Between the read lock and write lock, another goroutine may insert the key.
  - Fix: Use a single `sync.Mutex` that covers both the check and insert, or use `Set` that returns whether it was the first insert.

- Mistake: Not expiring old keys.
  - Why it happens: The in-memory map grows without bound as clients send new keys.
  - Fix: Implement a TTL-based eviction (background goroutine that deletes keys older than the retention window, typically 24 hours).

- Mistake: Using the request body hash as the idempotency key.
  - Why it happens: Retries may resend the same body, but the server cannot distinguish a retry from a distinct request with the same body.
  - Fix: The client must generate a unique key per request attempt, regardless of the body content.

## Debugging walkthrough

Buggy program:

```go
type BadStore struct {
	store map[string]string
}

func (s *BadStore) Set(key, val string) bool {
	if _, ok := s.store[key]; ok {
		return false
	}
	s.store[key] = val
	return true
}
```

Symptom: under concurrent access, two goroutines both see that the key does not exist (both get `ok == false`) and both insert. The function processes the payment twice.

Investigation: add a mutex. Without synchronization, the `if _, ok := s.store[key]; ok` check in one goroutine is not atomic with the `s.store[key] = val` in the same goroutine.

Fix: wrap the check-and-insert in a mutex lock as shown in the example.

## Production notes

Use a persistent dedup store (Redis or PostgreSQL) for production. In-memory stores are lost on restart, and restarts after a crash will accept duplicate requests. Set the idempotency key retention to at least the maximum retry interval plus a safety margin (24 hours is standard). For Redis, use `SET key value NX EX ttl` for atomic insert-if-not-exists with TTL. For PostgreSQL, use `INSERT INTO idempotency (key, response, created_at) VALUES ($1, $2, NOW()) ON CONFLICT (key) DO NOTHING`. The dedup store is a single point of consistency; make sure it is replicated and durable.

## Performance implications

An in-memory idempotency check with a mutex costs ~50 ns for the read (uncontended) and ~100 ns for the insert. A Redis-based check costs ~1 ms of network round-trip. The dedup store is on the critical path for every mutation request, so latency matters. Use a local cache (e.g., `sync.Map` with TTL) in front of Redis for hot keys. The store size must be bounded: a 24-hour retention at 1000 keys/second requires ~86 million entries. Use Redis with expiry or a PostgreSQL table with a TTL-based partition scheme.

## Practice task

Implement a `DedupStore` that supports `Set(key, value string, ttl time.Duration) bool` and `Get(key string) (string, bool)`. Use a `map[string]storedEntry` where each entry has a value and an expiration time. Launch a background goroutine that evicts expired entries every minute. Write table-driven tests for: set and get, duplicate rejection, concurrent set, TTL expiry, and the eviction goroutine.

## Tests / verification

```bash
go run ./curriculum/modules/11-lifecycle-context-concurrency/lessons/28-idempotency-keys
go test ./curriculum/modules/11-lifecycle-context-concurrency/lessons/28-idempotency-keys
```

## Review questions

1. Why must the idempotency key be generated by the client, not the server?
2. What happens if two concurrent requests arrive with the same idempotency key?
3. Why should the dedup store store the response, not just a boolean indicating the key was seen?
4. What is a reasonable TTL for an idempotency key, and why?
5. How do you ensure thread safety of a dedup store in Go?

## NEXT UP

Sync.Once
