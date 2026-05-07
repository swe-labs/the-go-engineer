# Known Limitations and Intentional Bounds

The Go Engineer is a **source-available educational curriculum**, not a drop-in Open Source template for arbitrary commercial systems. To maintain teaching clarity, several architectural and implementation trade-offs are made intentionally.

If you are using concepts from this curriculum in production, you should be aware of these bounds.

## Authentication and JWTs

The `internal/auth/token.go` implementation in the Opslane flagship provides a custom, HMAC-signed, JWT-compatible access token manager. 

**Why we built it this way:** To demystify cryptographic signing, base64url encoding, and identity extraction without hiding behind black-box frameworks or heavy third-party dependencies.

**Production Reality:** In most commercial production systems, you should not roll your own JWT infrastructure. Use mature, audited libraries (like `github.com/golang-jwt/jwt`) or delegate to managed identity providers (like Auth0, Clerk, or AWS Cognito) which handle key rotation, JWKS endpoints, and refresh tokens natively.

## In-Memory Observability Metrics

The `internal/metrics` package implements a custom, goroutine-safe registry with atomic counters and fixed-bucket latency histograms.

**Why we built it this way:** To teach synchronization primitives (`sync/atomic`, `sync.RWMutex`), bucket distributions, and the raw mechanics of how instrumentation impacts HTTP throughput.

**Production Reality:** You would typically use the official Prometheus Go client (`github.com/prometheus/client_golang`) or OpenTelemetry (`go.opentelemetry.io/otel`). They provide highly optimized metric exposition, standardized memory profiling, and dynamic label cardinality management.

## Worker Pools and Durability

The `internal/workers` pool processes jobs from an in-memory channel.

**Why we built it this way:** To teach bounded concurrency, worker goroutines, channel draining, and graceful shutdown sequencing.

**Production Reality:** If the process crashes abruptly, any jobs sitting in the Go channel are lost. In a real system, you would back this async work with a durable queue (like RabbitMQ, Amazon SQS, or Redis Streams), or implement the Outbox pattern in the PostgreSQL database to ensure at-least-once delivery semantics.

## Single-Node Event Bus

The `internal/events/bus.go` implementation is an in-memory channel router.

**Why we built it this way:** To decouple domain actions (like `CreateOrder`) from side effects (like `SendEmail` or `StartPayment`), demonstrating the Publisher/Subscriber pattern within a single binary.

**Production Reality:** This bus cannot route events between multiple microservices or scaled replicas. A real distributed system would use Kafka, AWS EventBridge, NATS, or similar distributed messaging backbones.

## Distributed Tracing

The `internal/otel` package now exports OTLP-compatible JSON payloads, supports trace context parsing (`traceparent`), random trace/span identifiers, and configurable sampling.

**Why we built it this way:** To keep tracing mechanics visible and teach propagation/export behavior without hiding all internals behind a full SDK wrapper.

**Production Reality:** Teams should still evaluate migration to the official OpenTelemetry SDK for richer semantics, ecosystem tooling, and standards coverage.

## Summary

The Opslane project teaches the *shape* of production code: clear boundaries, safe concurrency, and observable state. But it deliberately avoids importing the *weight* of production code when doing so would hide the Go mechanics you are here to learn.
