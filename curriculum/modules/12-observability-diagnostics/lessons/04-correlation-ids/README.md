# Correlation IDs

## Learning objective

Propagate a correlation ID across service boundaries using context.Context for in-process calls and HTTP headers for out-of-process calls, enabling end-to-end request tracing in a distributed system.

## Why this matters

A single user request in a microservice architecture touches 5–20 services: API gateway, auth, product catalog, inventory, payments, notifications. When one of these services fails, you need to see the entire chain — not just the failing service's logs. A correlation ID is a unique identifier generated at the entry point (API gateway or client middleware) and forwarded to every downstream service. Every service logs the correlation ID alongside its own events. When debugging, you search the log aggregator for the correlation ID and see every log line from every service in chronological order. Without correlation IDs, each service's logs are an island and you must manually stitch timestamps across 5 dashboards.

## Mental model

A correlation ID is a single UUID that travels with a request across every service boundary — like a unique ticket number attached to a customer's journey through a multi-department process. Every service that touches the request logs the correlation ID alongside its own events. When debugging a production issue, the operator searches the log aggregator for the correlation ID and sees every log line from every service in chronological order — a complete timeline of the request.

The flow: client (or API gateway) generates a correlation ID → passes it as an HTTP header (X-Correlation-Id) → service A extracts it, stores in context, logs it, forwards it to service B → service B extracts it, stores in context, logs it, forwards to service C → ... → when any service fails, the entire chain is visible by searching the correlation ID.

## Core idea

A correlation ID is a globally unique identifier that is generated once at the entry point of a request and propagated through every service, every goroutine, and every async operation that the request touches. It is never regenerated — the original ID travels unchanged through the entire system.

The three rules of correlation IDs:

1. **Generate once**: Create the ID at the first entry point (API gateway, load balancer, or client SDK).
2. **Propagate everywhere**: Forward the ID in every downstream call (HTTP headers, gRPC metadata, message queue headers).
3. **Log always**: Include the correlation ID in every log line, every metric tag, and every trace span.

## Under the hood

Correlation ID propagation in Go depends on `context.Context`. The context carries the correlation ID through the call chain without explicit function parameters. The key is an unexported type to prevent key collisions between packages:

```go
type cidKey struct{}

func WithCorrelationID(ctx context.Context, cid string) context.Context {
    return context.WithValue(ctx, cidKey{}, cid)
}

func GetCorrelationID(ctx context.Context) string {
    cid, _ := ctx.Value(cidKey{}).(string)
    return cid
}
```

For HTTP propagation, the standard header is `X-Correlation-Id` (also `X-Request-Id`). The middleware extracts it from the incoming request, and the HTTP client sets it on outgoing requests:

- **Incoming**: `r.Header.Get("X-Correlation-Id")` → store in context.
- **Outgoing**: `req.Header.Set("X-Correlation-Id", GetCorrelationID(ctx))` on every outgoing HTTP request.

For gRPC, it is carried in metadata (key-value pairs in the gRPC context). For message queues, it is carried in message headers (Kafka headers, RabbitMQ message properties). Go's `net/http` Transport does not automatically propagate context — the caller must explicitly copy the correlation ID from context to the outgoing request header.

## How Go uses it

In Go microservices, correlation ID propagation follows a consistent pattern:

```go
// Middleware to extract incoming correlation ID
func CorrelationMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        cid := r.Header.Get("X-Correlation-Id")
        if cid == "" {
            cid = uuid.New().String()
        }
        ctx := WithCorrelationID(r.Context(), cid)
        w.Header().Set("X-Correlation-Id", cid)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}

// HTTP client that propagates correlation ID
type Client struct {
    base *http.Client
}

func (c *Client) Do(ctx context.Context, req *http.Request) (*http.Response, error) {
    if cid := GetCorrelationID(ctx); cid != "" {
        req.Header.Set("X-Correlation-Id", cid)
    }
    return c.base.Do(req)
}
```

Every service in the chain uses the same middleware and client, so the correlation ID flows end-to-end.

## Go example

```go
package main

import (
	"context"
	"log/slog"
	"os"
)

type cidKey struct{}

func WithCorrelationID(ctx context.Context, cid string) context.Context {
	return context.WithValue(ctx, cidKey{}, cid)
}

func GetCorrelationID(ctx context.Context) string {
	cid, _ := ctx.Value(cidKey{}).(string)
	return cid
}

func SimulateServiceCall(ctx context.Context, service, action string) {
	slog.Info("service call",
		"service", service,
		"action", action,
		"correlation_id", GetCorrelationID(ctx),
	)
}

func CallServices(ctx context.Context, services []string) []string {
	var chain []string
	cid := GetCorrelationID(ctx)
	for _, svc := range services {
		chain = append(chain, cid)
		slog.Info("entering service", "service", svc, "correlation_id", cid)
	}
	return chain
}

func main() {
	handler := slog.NewTextHandler(os.Stdout, nil)
	slog.SetDefault(slog.New(handler))

	ctx := WithCorrelationID(context.Background(), "corr-abc-123")
	SimulateServiceCall(ctx, "api-gateway", "authenticate")
	SimulateServiceCall(ctx, "user-service", "get_profile")
	SimulateServiceCall(ctx, "payment-service", "charge")
	CallServices(ctx, []string{"auth", "orders", "inventory", "notifications"})
}
```

Run with `go run .` to see every log line carries the same `correlation_id=corr-abc-123` across all simulated service calls. In a log aggregator, searching for `corr-abc-123` returns all these lines in chronological order.

## Step-by-step execution

For a request that hits three services (gateway → auth → payments):

1. Client sends request with header `X-Correlation-Id: corr-xyz`.
2. **Gateway service**: middleware extracts `corr-xyz` from header, stores in context via `WithCorrelationID`. Calls `auth` service with header `X-Correlation-Id: corr-xyz`. Logs: `"incoming request" correlation_id=corr-xyz path=/api/orders`.
3. **Auth service**: middleware extracts `corr-xyz`, stores in context. Verifies token. Calls `payments` service with header. Logs: `"token verified" correlation_id=corr-xyz user_id=u-42`.
4. **Payments service**: middleware extracts `corr-xyz`. Processes charge. Logs: `"charge processed" correlation_id=corr-xyz amount=59.99`.
5. Response propagates back through the chain. All three services logged with the same `correlation_id=corr-xyz`.
6. In the log aggregator, searching `corr-xyz` returns three log lines in order — a complete timeline of the request.

## Common mistakes

- **Generating a new correlation ID at each service boundary**: The trace is broken into disconnected segments that cannot be correlated during debugging. Fix: always propagate the incoming ID; generate a new ID only if one does not exist.

- **Storing the correlation ID in a global variable**: Concurrent requests overwrite each other's correlation IDs, mixing traces. Every request appears to have the same ID. Fix: always store in context, never in a global.

- **Logging the correlation ID in only some log lines**: Searching for a trace produces partial results. Fix: use request-scoped logging (lesson 03) to stamp every log line automatically with the correlation ID.

- **Not forwarding to downstream services**: The trace stops at the first service boundary. Downstream operations appear as orphan events with no correlation ID. Fix: always copy the correlation ID from context to outgoing request headers/metadata.

- **Using a mutable field that can be overwritten**: Some developers use a `context.WithValue` with a string key or a mutable pointer, which can be overwritten by intermediate code. Fix: use an unexported struct key type, which prevents any code outside the package from setting or reading the value.

## Debugging walkthrough

Consider a distributed system with three services (api-gateway, orders, payments). Users report that some orders fail silently — the order is created but never charged.

**Symptom**: Order service shows `"order created"` log lines, but payment service shows no corresponding `"charge processed"` lines.

**Investigation**: Without correlation IDs, you see:
- Order logs: `"order created" order_id=1234 user_id=42`
- Payment logs: nothing matching

You cannot tell which order log line corresponds to which payment call. Was the payment never made, or was it made but with a different correlation?

**With correlation IDs**, the order logs show:
```
level=INFO msg="order created" order_id=1234 correlation_id=corr-xyz
level=INFO msg="calling payment service" correlation_id=corr-xyz
```

The payment logs show nothing for `corr-xyz`. Searching further, you find a log line 100ms later:
```
level=ERROR msg="payment service timeout" correlation_id=corr-xyz timeout_ms=5000
```

**Root cause**: The payment service timed out, the order service did not handle the timeout gracefully (no retry, no error propagation to the user), and the order was left in a created-but-unpaid state.

**Fix**: Add correlation ID propagation, then implement a retry-with-backoff in the order service's payment client. With the correlation ID, the bug was identified in 2 minutes instead of 30.

## Production notes

Correlation IDs are the foundation of observability in distributed systems. Every major production Go service uses them:

- **Kubernetes audit logs**: Every API request carries a `requestID` that correlates the API server logs with the etcd logs.
- **Docker Hub API**: Every HTTP request carries a correlation ID across all microservices.
- **Uber's observability pipeline**: Correlation IDs are part of the tracing infrastructure, processing billions of spans per day.

In incident response, the first question is always: "what is the correlation ID?" Without it, finding the root cause of a multi-service failure requires manually stitching together timestamps from 5 different log dashboards — a 10-minute investigation becomes a 2-hour slog.

Correlation IDs also enable automated analysis: an incident response tool can fetch all logs for a given correlation ID, build a timeline, and suggest the root cause based on which service logged an error first.

## Performance implications

Correlation ID propagation adds negligible overhead:

- **Context storage**: One `context.WithValue` call per request boundary, costing ~1 allocation of ~16 bytes.
- **Header extraction**: One `Header.Get` call in middleware, costing a map lookup in the HTTP header map.
- **Header setting**: One `Header.Set` call per downstream HTTP call, costing a map insert.
- **Log inclusion**: One `slog.String` attribute per log line, costing ~50ns for serialization.

The total cost per request is well under 1μs — free for all practical purposes.

The far greater cost is the _absence_ of correlation IDs, which causes hours of wasted engineering time per incident. The implementation cost of correlation IDs is ~50 lines of middleware code.

## Practice task

Complete the correlation ID propagation functions:

1. `WithCorrelationID(ctx, cid)` stores the correlation ID in context.
2. `GetCorrelationID(ctx)` retrieves it from context.
3. `CallServices(ctx, services)` logs entry into each service with the correlation ID and returns the chain.

The tests verify that all services in the chain receive the same correlation ID, and that `GetCorrelationID` returns `""` when no ID is set.

## Tests / verification

```bash
go run ./curriculum/modules/12-observability-diagnostics/lessons/04-correlation-ids
go test ./curriculum/modules/12-observability-diagnostics/lessons/04-correlation-ids
```

## Review questions

1. Why must a correlation ID be generated once and never regenerated at service boundaries?
2. What could go wrong if you store the correlation ID in a global variable instead of context.Context?
3. What HTTP header is conventionally used to propagate correlation IDs between services?
4. Describe the debugging workflow for a multi-service failure with and without correlation IDs.
5. How would you propagate a correlation ID through an asynchronous message queue (e.g., Kafka, RabbitMQ)?

## NEXT UP

PII redaction — ensuring sensitive user data (passwords, emails, SSNs) is never written to logs.
