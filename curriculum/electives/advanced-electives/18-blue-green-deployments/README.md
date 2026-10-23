# Blue/green deployments

## When you need this

Rolling updates minimize downtime but can cause issues when a bad deployment is gradually rolled out — some users see the bad version while others see the good version, and rollback takes multiple cycles. Blue/green deployments solve this by maintaining two identical environments and switching traffic instantly — deployment is either fully on the new version or fully on the old version.

## When you do not need this

When simpler concrete types or standard core library primitives suffice without adding abstraction overhead.

## Prerequisites

- elective-17

## Mental model

Two identical environments (blue and green) run the application. At any time, one environment serves production traffic (active) and the other is idle or running the previous version (standby). To deploy: deploy the new version to the standby environment, run smoke tests, switch the load balancer to the standby environment (making it active), and keep the old environment as a rollback target.

## How it works

Blue/green requires: (1) infrastructure that supports two identical environments (separate compute, database, cache), (2) a load balancer or proxy that can switch traffic instantly (HAProxy, Nginx, AWS ALB, Kubernetes Service), (3) database migrations that are backward-compatible (old version can read/write the schema that the new version creates), and (4) a shared session store so users do not lose sessions on switch.

## Go example

```go
package main

import "fmt"

func main() {
    fmt.Println("Elective: Blue/green deployments")
}
```

## Tradeoffs

- Blue/green requires 2x infrastructure capacity during deployment — both environments run simultaneously, doubling compute and database costs
- Traffic switch is instantaneous with a load balancer (no cold start) — no performance degradation during switch
- Green's caches are cold after switch — first few seconds of requests may be 2-5x slower until caches warm
- Database connection pool on green starts empty — connection establishment adds latency to initial requests after switch

## Common mistakes

- Not making database migrations backward-compatible — the old environment cannot read the new schema after migration, making rollback impossible without a reverse migration.
- Skipping smoke tests on the standby environment — the new version is deployed but never tested; the switch puts a broken version into production.
- Not pre-warming caches — the standby environment starts with cold caches; the first requests after switch are slow (cache stampede).
- Using DNS-based traffic switching — DNS changes take minutes to propagate; use a load balancer for instant switching.
- Keeping the old environment for too long — the old environment costs 2x infrastructure; decommission after the new version is verified stable (typically 1-7 days).

## Practice task

The learner must design a blue/green deployment strategy for a Go service with a PostgreSQL database, including: (1) backward-compatible migration, (2) traffic switch mechanism, (3) health check on the standby environment before the switch, (4) cache pre-warming, and (5) rollback procedure.

## Production notes

Blue/green deployments are used by organizations requiring instant rollback and zero-downtime deployments: financial services, healthcare, e-commerce. Cloud providers support blue/green natively: AWS CodeDeploy, Azure Deployment Slots, Google Cloud App Engine. Kubernetes can implement blue/green with two Deployments and label-based Service switching.

## Related core lessons

- Module 16.17 — Feature flags
- Module 17.16 — CQRS
