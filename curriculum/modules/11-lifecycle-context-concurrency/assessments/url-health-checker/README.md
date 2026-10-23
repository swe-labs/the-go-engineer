# Project Assessment — URL Health Checker

## Purpose

Evaluate your ability to build a concurrent HTTP health check probe engine in Go. This assessment verifies that you can implement bounded fan-out/fan-in worker patterns, enforce per-probe deadline contexts (`context.WithTimeout`), calculate HTTP latency, track state transitions (`UP` ↔ `DOWN`) thread-safely with `sync.RWMutex`, and eliminate data races under `go test -race`.

## Targets

- `project-url-health-checker`

## Scope

This assessment covers:

- Bounded fan-out/fan-in concurrency: distributing URL targets across a fixed worker pool (`min(workers, len(urls))`)
- Per-probe context timeouts: preventing hanging requests from blocking worker goroutines indefinitely
- HTTP request execution with `http.RequestWithContext` and latency measurement using `time.Since`
- State machine transition tracking: detecting status changes and logging events under `sync.RWMutex`
- Order-preserving result collation: indexing probe results to match original input URL order
- Error classification: distinguishing network failures, timeouts, and non-2xx status codes
- Data race elimination under concurrent probe execution verified by `go test -race`

## Evidence required

Submit:

1. Your completed `_starter/main.go` with single probe timeout execution, transition tracking, and bounded concurrent batch checks implemented
2. Output of `go test -v -race .` showing all unit and concurrent probe tests passing with zero data races
3. A written explanation answering: How did you size the bounded worker pool, and why is an unbounded goroutine-per-URL approach dangerous when scanning thousands of endpoints?

## Time expectation

Plan for 90 to 120 minutes.

## Rubric

| Criterion | Weight | Excellent (5) | Good (3) | Needs Work (0) |
|-----------|--------|--------------|----------|----------------|
| **Correctness** | 35% | Probes execute within timeouts; status codes and latencies recorded; results preserve input order; all automated tests pass | Most endpoints probed; minor latency calculation or ordering flaw | Non-functional probe engine or race detector failure |
| **Timeout & context control** | 25% | Every probe governed by timeout context; slow or hung endpoints aborted promptly without worker exhaustion | Timeout set globally on client without per-request context | Hangs indefinitely on slow endpoints |
| **State & transition tracking** | 25% | Thread-safe transition history recording under `sync.RWMutex`; detects flapping state changes (`UP` → `DOWN`) cleanly | State recorded but potential race conditions on concurrent map access | Misses status changes or panics on concurrent map access |
| **Concurrency & explanation** | 15% | Zero data races under `-race`; learner explains worker pool sizing, socket exhaustion risks, and synchronization mechanics | Race detector clean but explanation has minor conceptual gaps | Data races detected or unable to explain worker coordination |

**Passing score**: 80%

## Review questions

1. How did you size the bounded worker pool, and why is an unbounded goroutine-per-URL approach dangerous when scanning thousands of endpoints?
2. How do you differentiate between DNS resolution failure, connection timeout, TLS handshake failure, and HTTP status code errors in your health reports?
3. Walk me through your synchronization model: how do channels and `sync.WaitGroup` coordinate probe dispatches, result collection, and graceful shutdown?
4. Why must the HTTP response body be closed (`defer resp.Body.Close()`) even if you only read the status code and ignore the response body?
5. What is endpoint flapping, and how does your transition event tracker detect rapid status oscillations?

## Passing standard

You pass when all race detector tests pass, timeouts abort hanging probes promptly, transitions are tracked thread-safely, and you can explain worker sizing and synchronization models.

## Retake policy

Learner may revise code and resubmit after addressing all rubric gaps.

## NEXT UP

[Module 11 Checkpoint](../checkpoint/README.md)
