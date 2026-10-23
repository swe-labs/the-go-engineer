# gRPC interceptors

## When you need this

gRPC services need cross-cutting concerns — logging, authentication, rate limiting, error handling — for every RPC. Adding these to every handler duplicates code and is easy to miss on new endpoints. gRPC interceptors solve this by providing middleware hooks that run before (or after) every RPC, keeping handlers focused on business logic.

## When you do not need this

When simpler concrete types or standard core library primitives suffice without adding abstraction overhead.

## Prerequisites

- elective-11

## Mental model

An interceptor is a function that wraps an RPC handler. The server creates an interceptor chain at startup. For each incoming RPC, gRPC calls the interceptor chain instead of the handler directly. Each interceptor can inspect the request, modify the context, call the next interceptor in the chain, inspect the response, and modify or replace the response. Client interceptors wrap outgoing RPC calls with the same pattern.

## How it works

A unary server interceptor has the signature: func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error). The interceptor calls handler(ctx, req) to invoke the next interceptor or the final handler. Stream interceptors wrap streaming RPCs with a similar pattern using grpc.StreamServerInfo and grpc.StreamHandler. Client interceptors have analogous signatures for outgoing calls.

## Go example

```go
package main

import "fmt"

func main() {
    fmt.Println("Elective: gRPC interceptors")
}
```

## Tradeoffs

- Each interceptor adds function call overhead — a chain of 5 interceptors adds ~50ns per call, negligible for typical RPCs
- Interceptors that allocate (log entries, metric tags) increase GC pressure per RPC — batch or sample in high-throughput services
- Stream interceptors wrap the entire stream lifecycle — stream creation, each message send/recv, and stream close — cost scales with message count
- Unary interceptor allocation: each interceptor call creates a closure (the handler wrapper) — pool interceptor instances to reduce allocation

## Common mistakes

- Forgetting to call handler(ctx, req) in the interceptor — the handler never executes, and the client gets no response.
- Modifying the request after calling handler — the request is already processed; the interceptor should modify the request before calling handler.
- Using a unary interceptor for streaming RPCs — stream RPCs must use stream interceptors, not unary interceptors.
- Not propagating context values set by interceptors — context values are passed via context.WithValue and must be extracted by downstream interceptors and handlers.
- Blocking in an interceptor — an interceptor that calls time.Sleep or blocks on I/O blocks the entire gRPC server's handler goroutine.

## Practice task

The learner must implement a unary server interceptor that logs each RPC with method name, request size, duration, and response status, register it with the gRPC server, and verify that every RPC produces a log line with the expected fields.

## Production notes

gRPC interceptors are used in every production gRPC service: auth interceptors validate JWT tokens from gRPC metadata, logging interceptors log every RPC with duration and status, recovery interceptors catch panics and return gRPC Internal errors, rate limit interceptors reject requests when the rate limit is exceeded, and metrics interceptors record request counts and latency histograms.

## Related core lessons

- Module 16.11 — gRPC streaming
- Module 17.10 — gRPC fundamentals
