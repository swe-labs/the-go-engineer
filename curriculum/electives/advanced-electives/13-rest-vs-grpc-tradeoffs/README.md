# REST vs gRPC tradeoffs

## When you need this

Go developers must choose between REST and gRPC when designing service APIs, but the decision is rarely clear-cut. REST is human-readable, browser-friendly, and has universal tooling. gRPC is faster, contract-enforced, and supports streaming. Choosing incorrectly leads to poor developer experience, performance bottlenecks, or maintenance overhead.

## When you do not need this

When simpler concrete types or standard core library primitives suffice without adding abstraction overhead.

## Prerequisites

- elective-12

## Mental model

REST and gRPC are different tools for the same job (service-to-service communication) with different tradeoffs. REST: resources as URLs, JSON payloads, HTTP methods, cacheable, browser-native. gRPC: services as .proto contracts, protobuf binary, HTTP/2, streaming, code-generated clients. REST prioritizes accessibility and simplicity; gRPC prioritizes performance and contract rigor.

## How it works

REST over HTTP/1.1 uses text-based JSON, one request per TCP connection (without keep-alive), and manual documentation of endpoints. gRPC over HTTP/2 uses binary protobuf, multiplexed streams over one connection, automatic code generation, and built-in service discovery. Performance difference: protobuf is 3-10x faster to serialize and 30-70% smaller payloads than JSON. HTTP/2 multiplexing eliminates connection overhead for multiple concurrent requests.

## Go example

```go
package main

import "fmt"

func main() {
    fmt.Println("Elective: REST vs gRPC tradeoffs")
}
```

## Tradeoffs

- gRPC is typically 5-10x faster than REST for equivalent payloads due to protobuf binary encoding and HTTP/2 multiplexing
- REST over HTTP/1.1 has head-of-line blocking — one slow response blocks subsequent requests on the same connection; HTTP/2 (gRPC) eliminates this
- gRPC's binary format is ~3-5x smaller than JSON for the same data, reducing bandwidth and deserialization time
- REST's text format (JSON) is debuggable with curl and browser dev tools — gRPC requires grpcurl or similar specialized tools
- For low-throughput APIs (<100 RPS), the performance difference is negligible — developer productivity and tooling matter more

## Common mistakes

- Using gRPC for browser-facing APIs without gRPC-Web — browsers do not natively support gRPC; use gRPC-Web proxy or REST for browser clients.
- Using REST for high-throughput internal microservice APIs — JSON serialization overhead adds latency and bandwidth costs at scale.
- Not leveraging HTTP caching for REST APIs — REST responses can be cached at the CDN or client level; gRPC responses are not cacheable without a proxy.
- Forgetting that gRPC requires protobuf schema management — .proto files must be shared between services; breaking changes require coordinated deployments.
- Choosing REST just because it is simpler — REST's simplicity comes at the cost of performance, contract enforcement, and streaming support.

## Practice task

The learner must design the API protocol for three scenarios: (1) a public API for a weather service (REST), (2) an internal service for processing real-time stock trades (gRPC streaming), (3) a CRUD API for a CMS admin panel (REST), and explain the tradeoff decision for each scenario.

## Production notes

Most Go microservices use both REST and gRPC: REST for external APIs (public-facing, browser clients, third-party integrations) and gRPC for internal APIs (service-to-service, high-throughput, streaming). Companies like Google, Uber, and Netflix use gRPC internally and REST externally. The grpc-gateway project generates REST endpoints from .proto files, allowing a single service contract for both protocols.

## Related core lessons

- Module 16.12 — gRPC interceptors
- Module 17.11 — gRPC streaming
