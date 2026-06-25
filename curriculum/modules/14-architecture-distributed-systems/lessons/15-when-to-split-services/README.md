# When to split services

## Learning objective

Identify bounded contexts that are candidates for service splitting, evaluate team size and deployment coupling triggers, design splitting strategies (strangler fig, database separation, event-driven), and apply data ownership principles.

## Why this matters

Splitting a monolith into microservices is one of the highest-risk engineering decisions. Do it too early -- you pay the distributed complexity tax before reaping any benefit. Do it too late -- the monolith becomes undeployable, with code that no single developer understands. Knowing when to split is the difference between a well-factored system and a distributed big ball of mud. Go engineers leading architecture decisions must recognize the triggers and choose the right splitting strategy for each context.

## Mental model

A monolith is a single apartment building. As the building grows, more people share the same plumbing, elevator, and maintenance crew. When the building reaches 100 residents, the elevator is always crowded, the plumbing fails weekly, and the maintenance crew cannot keep up. The solution is to build separate buildings (services) for different teams. But building a new building costs time and money, and the old building still needs maintenance.

The trigger to split is not the size of the building itself but the friction of shared resources. If two teams can deploy their changes independently without stepping on each other, keep the monolith. If every deploy requires coordination between three teams, it is time to split.

## Core idea

Triggers for splitting:

| Trigger | Description |
|---|---|
| Team size > 8 (Dunbar number) | Beyond 8 developers on a single codebase, communication overhead dominates |
| Deployment coupling | Every deploy requires coordinated releases across teams |
| Data ownership conflict | Two teams frequently change the same database schema |
| Scaling requirements differ | One part of the system needs more resources than another |
| Technology divergence | A team wants to use a different database or language |

Splitting strategies:

| Strategy | Approach | Downtime | Complexity |
|---|---|---|---|
| Strangler Fig | Route traffic incrementally to new service | No | Low |
| Database Separation | Extract schema into own database | Yes | Medium |
| Event-driven Split | Decouple via events, build new consumer | No | High |

Data ownership: each service owns its data. No other service reads or writes it directly. Data sharing happens through APIs or events, not shared databases. This is the most important rule of service decomposition.

## Under the hood

The bounded context is the unit of splitting. In Domain-Driven Design, a bounded context is a boundary where a particular domain model applies. The Order Bounded Context owns orders, order items, and payments. The Inventory Bounded Context owns stock levels and warehouse locations. These contexts map to potential services.

To split, identify the bounded contexts in the codebase:
1. Look for packages or modules that change independently.
2. Look for database tables that are always modified together.
3. Look for teams that own specific features.

Once a bounded context is identified, apply the strangler fig pattern:
1. Add a routing layer (API gateway, load balancer) that can route requests to either the monolith or the new service.
2. Implement one endpoint in the new service.
3. Route traffic for that endpoint to the new service.
4. Verify correctness.
5. Repeat for each endpoint.
6. Remove the old endpoint code from the monolith.

## How Go uses it

- **`net/http` with `httputil.ReverseProxy`**: Go's standard library has everything needed to build the routing layer for a strangler fig. Route based on URL path, headers, or tenant ID.
- **`chi` router**: supports middleware-based routing decisions. A `ShouldRouteToNewService(tenantID)` middleware can redirect specific requests.
- **`wire` or `fx`**: dependency injection frameworks that make it easier to extract packages into separate modules before full service separation.
- **`protobuf` / `gRPC`**: service contracts between the monolith and the new service. Define the interface first, then implement on both sides.

## Go example

```go
package main

import (
	"fmt"
)

type BoundedContext struct {
	Name      string
	OwnsData  []string
	TeamSize  int
}

type SplittingAnalysis struct {
	Context           string
	Reasons           []string
	RecommendedAction string
	RiskLevel         string
}

func AnalyzeSplit(ctx BoundedContext) SplittingAnalysis {
	a := SplittingAnalysis{Context: ctx.Name}
	if ctx.TeamSize > 8 {
		a.Reasons = append(a.Reasons, "team exceeds 8 developers")
	}
	if len(ctx.OwnsData) > 5 {
		a.Reasons = append(a.Reasons, "owns more than 5 data entities")
	}

	if len(a.Reasons) >= 2 {
		a.RecommendedAction = fmt.Sprintf("split %s into separate service", ctx.Name)
		a.RiskLevel = "medium"
	} else if len(a.Reasons) >= 1 {
		a.RecommendedAction = fmt.Sprintf("monitor %s for further growth", ctx.Name)
		a.RiskLevel = "low"
	} else {
		a.RecommendedAction = "keep as part of monolith"
		a.RiskLevel = "none"
	}
	return a
}

type ServiceCoupling struct {
	ServiceA   string
	ServiceB   string
	SharedData []string
	CallPerSec int
}

func EvaluateCoupling(c ServiceCoupling) string {
	if len(c.SharedData) > 0 && c.CallPerSec > 100 {
		return fmt.Sprintf("high coupling: %s and %s should consider merging", c.ServiceA, c.ServiceB)
	}
	if len(c.SharedData) > 0 {
		return fmt.Sprintf("medium coupling: %s and %s share %v", c.ServiceA, c.ServiceB, c.SharedData)
	}
	return fmt.Sprintf("low coupling: %s and %s are independent", c.ServiceA, c.ServiceB)
}

var SplitStrategies = map[string][]string{
	"strangler-fig": {
		"add proxy routing layer",
		"migrate one endpoint at a time to new service",
		"route new service via proxy",
		"remove old endpoint code when traffic fully migrated",
	},
	"database-separation": {
		"extract schema for bounded context",
		"set up new database",
		"run dual-writes to both databases",
		"switch reads to new database",
		"remove old schema from shared database",
	},
}

func main() {
	orderCtx := BoundedContext{
		Name:     "Order Management",
		OwnsData: []string{"orders", "items", "payments", "refunds", "shipments", "history"},
		TeamSize: 9,
	}
	analysis := AnalyzeSplit(orderCtx)
	fmt.Printf("Context: %s\n", analysis.Context)
	for _, r := range analysis.Reasons {
		fmt.Printf("  - %s\n", r)
	}
	fmt.Printf("Action: %s (risk: %s)\n\n", analysis.RecommendedAction, analysis.RiskLevel)

	coupling := ServiceCoupling{
		ServiceA:   "Orders",
		ServiceB:   "Inventory",
		SharedData: []string{"product_sku"},
		CallPerSec: 300,
	}
	fmt.Println(EvaluateCoupling(coupling))

	fmt.Println("\nStrangler Fig strategy:")
	for _, s := range SplitStrategies["strangler-fig"] {
		fmt.Printf("  - %s\n", s)
	}
}
```

## Step-by-step execution

For `AnalyzeSplit(Order Management, 6 entities, team of 9)`:

1. Check `TeamSize`: 9 > 8 → append reason `"team exceeds 8 developers"`.
2. Check `OwnsData`: 6 > 5 → append reason `"owns more than 5 data entities"`.
3. `len(Reasons) == 2 >= 2` → recommended action: `"split Order Management into separate service"`.
4. Risk level: `"medium"`.

For `EvaluateCoupling(Orders, Inventory, shared=["product_sku"], 300 req/s)`:

1. `len(SharedData) > 0 AND CallPerSec > 100` → high coupling.
2. Return: `"high coupling: Orders and Inventory should consider merging"`.

## Common mistakes

- **Splitting because it is trendy**: microservices add latency, operational complexity, and debugging difficulty. Do not split without a concrete trigger like team growth or deployment coupling.
- **Splitting by technical layer instead of business domain**: `users-api`, `users-service`, `users-db` is three services for one feature. Split by bounded context: each service owns a complete business capability.
- **Shared database across services**: the most common anti-pattern in microservices. Two services reading the same table directly creates implicit coupling. Every schema change must be coordinated.
- **Forgetting data migration**: when splitting, existing data must be migrated to the new service's database. Plan the migration with dual-writes and backfill before cutting over.
- **Premature optimization**: a monolith serving 100 req/s does not need splitting. Wait until the monolith causes measurable pain -- slow deploys, team coordination overhead, or scaling bottlenecks.

## Debugging walkthrough

A team of 12 developers works on a single monolith. Every deploy takes 45 minutes and requires 3 teams to coordinate.

**Symptom**: The `auth` team changes the user table schema, breaking the `orders` team's queries. Deployments are scheduled weekly because any change risks breaking someone else's feature.

**Investigation**:
1. Identify bounded contexts: `auth`, `orders`, `payments`, `inventory`, `notifications`.
2. Measure deployment frequency per team: all teams deploy on the same schedule.
3. Count cross-team code changes: 60% of commits touch code owned by another team.

**Root cause**: No bounded context boundaries. Every team owns every table. A single deploy is a ball of mud.

**Fix**: Apply the strangler fig. Extract `payments` first because it has the clearest ownership boundary and the fewest shared tables. Route payment traffic to the new service. Deploy time drops from 45 minutes to 15 minutes for the monolith. The payments team deploys independently in 5 minutes.

## Production notes

- **Start with the most independent context**: identify the bounded context with the fewest dependencies on the rest of the monolith. Extract that one first. The risk is lowest and the team learns the splitting process.
- **Keep the monolith until the split is proven**: do not split more than one context at a time. Run both the monolith and the new service in production. Compare outputs for correctness.
- **Feature flags for routing**: use feature flags to control the percentage of traffic routed to the new service. Start with 1% of requests, verify correctness, ramp to 100%, then remove old code.
- **Observability across services**: distributed tracing is essential. Without tracing, debugging a request that spans 3 services is nearly impossible. Use OpenTelemetry.

## Performance implications

- **Splitting adds network latency**: every monolith-internal function call becomes an HTTP/gRPC call. Expect 1-10ms added per cross-service call.
- **Serialization overhead**: Go structs passed as function arguments become JSON or protobuf bytes. Protobuf serialization adds 1-5us; JSON adds 10-50us.
- **Connection overhead**: each service maintains connection pools to its database. N services means N connection pools, each consuming memory.
- **Transactional complexity**: a monolith can use a single database transaction for an order + payment + inventory update. Across services, you need a saga. Sagas are slower and more complex.

## Practice task

Given a monolith with four bounded contexts (Inventory, Orders, Shipping, Notifications) and 4 data entities per context, simulate the splitting decision. Write a function `RecommendSplit(contexts []BoundedContext) []SplittingAnalysis` and run it. Identify which context should be split first and why.

## Tests / verification

```bash
go run ./curriculum/modules/14-architecture-distributed-systems/lessons/15-when-to-split-services
go test ./curriculum/modules/14-architecture-distributed-systems/lessons/15-when-to-split-services
```

The tests verify split analysis for small/large teams, coupling evaluation, and strategy definitions.

## Review questions

1. What is the strongest predictor that a monolith should be split into services?
2. Why should a service own its data rather than sharing a database?
3. What is the strangler fig pattern and why is it preferred for gradual migration?
4. What is the risk of splitting a monolith before there is deployment coupling?
5. How would you determine which bounded context to extract first?

## NEXT UP

Microservices as a tradeoff, not a default.
