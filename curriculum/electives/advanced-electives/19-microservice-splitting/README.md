# Microservice splitting

## When you need this

Monolith applications become difficult to maintain as they grow — deployment takes hours, a single bug affects all features, and teams step on each other's code. Microservice splitting solves this by dividing the monolith into independently deployable services owned by separate teams, each with its own data store and deployment pipeline.

## When you do not need this

When simpler concrete types or standard core library primitives suffice without adding abstraction overhead.

## Prerequisites

- elective-18

## Mental model

A monolith is a single deployable unit with all features and data in one codebase and database. Microservice splitting extracts bounded contexts (groups of related features) into separate services. Each service has its own codebase, its own database, its own CI/CD pipeline, and its own on-call team. Services communicate via APIs (HTTP, gRPC, message queues). The monolith shrinks as services are extracted.

## How it works

Bounded contexts are identified by analyzing the monolith's domain: which features change together? Which features share data? Which features have different scaling or latency requirements? Each bounded context becomes a candidate for extraction. The extraction process: (1) identify the bounded context, (2) define the API contract (proto file or OpenAPI spec), (3) extract the data (new database with migrated schema), (4) extract the code (new service with its own repo), (5) route traffic from monolith to new service.

## Go example

```go
package main

import "fmt"

func main() {
    fmt.Println("Elective: Microservice splitting")
}
```

## Tradeoffs

- In-process function call: ~1-50ns (direct call or interface dispatch)
- Same-datacenter HTTP call: ~1-10ms (serialization + network + deserialization) — 1000-10000x slower than in-process
- Distributed transactions are 10-100x slower than local transactions — prefer sagas or eventual consistency across service boundaries
- Network calls fail: retries, timeouts, circuit breakers add latency — budget 100-500ms for any cross-service call to allow for retries

## Common mistakes

- Splitting too early — the bounded context is not well understood; frequent API changes couple the services more tightly than the monolith did.
- Sharing a database between services — services share data via API calls, not database queries; a shared database creates coupling that prevents independent deployment.
- Creating a distributed monolith — services that call each other synchronously in a chain; a failure in one service cascades through all services.
- Not handling network failures — a service call can fail, timeout, or return garbage; every cross-service call must have retries, timeouts, and circuit breakers.
- Splitting by technical layer instead of business capability — having a 'data service', 'logic service', and 'presentation service' instead of 'order service', 'payment service', 'inventory service'.

## Practice task

The learner must identify three bounded contexts in a sample e-commerce monolith, define the API contract for one extracted service, and explain the migration strategy from monolith to services: data migration, traffic routing, and rollback plan.

## Production notes

Microservice splitting is used by every large-scale organization: Netflix (500+ services), Amazon (1000+ services), Uber (2200+ services). Go is a popular choice for microservices due to its small binary size, fast startup, and efficient resource usage. Opslane uses Go microservices for orders, payments, inventory, and notifications — each with its own PostgreSQL database and independent deployment pipeline.

## Related core lessons

- Module 16.18 — Blue/green deployments
- Module 17.17 — Feature flags
