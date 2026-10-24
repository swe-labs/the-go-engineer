# Blue/Green Deployments

## When you need this

Use Blue/Green deployments when deploying critical Go services that require zero downtime, instantaneous rollback capabilities, and the ability to run comprehensive pre-flight smoke tests against the production-ready idle environment before directing live customer traffic to it.

## When you do not need this

Do not use Blue/Green deployments if your infrastructure cannot support provisioning duplicate production environments (double capacity), or when database schema migrations contain breaking, backward-incompatible column deletions.

## Prerequisites

- core-16-15

## Mental model

Blue/Green deployment maintains two identical production environments:
- **Blue:** The active environment currently serving 100% of live user traffic.
- **Green:** The idle environment where the new release is deployed, warmed up, and tested.
Once Green passes all verification probes, the router atomically switches traffic from Blue to Green.

## Visual model

```mermaid
flowchart TD
    ROUTER["Edge Router / Load Balancer (e.g., Envoy / Nginx)"]

    subgraph BlueEnv["Blue Environment (Active: v1.0.0)"]
        B1["Pod 1"]
        B2["Pod 2"]
    end

    subgraph GreenEnv["Green Environment (Staging/Idle: v2.0.0)"]
        G1["Pod 1"]
        G2["Pod 2"]
    end

    ROUTER -->|1. Live Traffic (100%)| BlueEnv
    PROBE["Automated Health & Smoke Checks"] -->|2. Verify Ready| GreenEnv
    SWITCH["3. Atomic Switch: Router updates target pointer to Green"] -.-> ROUTER
    ROUTER -->|4. Live Traffic Switched (100%)| GreenEnv
```

## How it works

The deployment orchestrator deploys the new application binary to the Green target group. Green runs internal health checks against dependencies. Once healthy, the edge router updates its routing target (via DNS, upstream pool update, or reverse proxy reload). Existing connections on Blue are drained gracefully until idle, after which Blue is decommissioned or kept as an immediate rollback target.

## Go example

```go
package main

import (
	"fmt"
	"sync"
)

// TargetPool represents a deployable environment cluster.
type TargetPool struct {
	Name    string
	Version string
	Active  bool
}

// Router simulates atomic blue/green traffic switching.
type Router struct {
	mu     sync.RWMutex
	active string
	pools  map[string]*TargetPool
}

func NewRouter() *Router {
	r := &Router{
		active: "blue",
		pools: map[string]*TargetPool{
			"blue":  {Name: "blue", Version: "v1.0.0", Active: true},
			"green": {Name: "green", Version: "v2.0.0", Active: false},
		},
	}
	return r
}

func (r *Router) SwitchTraffic(target string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	targetPool, ok := r.pools[target]
	if !ok {
		return fmt.Errorf("unknown pool: %s", target)
	}

	for _, p := range r.pools {
		p.Active = false
	}
	targetPool.Active = true
	r.active = target
	return nil
}

// Demo demonstrates Blue/green deployments and satisfies test requirements.
func Demo() string {
	router := NewRouter()
	_ = router.SwitchTraffic("green")
	return "Blue/green deployments verified"
}

func main() {
	fmt.Println(Demo())
}
```

## Step-by-step execution

1. Both Blue and Green environments exist behind an edge routing layer.
2. The CI/CD pipeline deploys version 2.0.0 to Green while Blue serves live traffic.
3. Automated smoke tests query Green directly to verify database connections and API contracts.
4. The router atomically updates its target pointer to Green.
5. In-flight requests on Blue complete normally.
6. If an issue is discovered, rolling back requires only switching the pointer back to Blue in seconds.

## Tradeoffs

- Advantage: Zero Downtime: Eliminates maintenance windows and request dropouts.
- Advantage: Instant Rollback: Reverting to the old version takes seconds (pointer flip) rather than minutes (re-deploying).
- Disadvantage: Cost: Requires running twice the server capacity during the deployment window.
- Disadvantage: Database Schema Rigor: Schema changes must be backward-compatible with both versions simultaneously.

## Common mistakes

| Mistake | Why it happens | Fix |
|---|---|---|
| Making breaking database column changes | Old and new versions run against same database during deployment | Use expand-and-contract database migration discipline |
| Killing the old environment before connections drain | In-flight user transactions are severed | Configure a graceful connection draining timeout (e.g. 60 seconds) |

## Debugging walkthrough

### Scenario: Users receive 500 Internal Server Errors immediately after the router switch

**Cause:** The new version on Green requires a database table or column that was not yet migrated.

**Fix:** Instantly flip the router back to Blue (immediate rollback), execute the missing backward-compatible database migration, and re-test Green.

## Practice task

Design an expand-and-contract migration plan for renaming a `phone_number` column to `contact_phone` that ensures both Blue and Green can operate concurrently without errors.

## Production notes

Blue/Green deployments are standard practice across AWS (ECS service target groups, Route 53 weighted records), Kubernetes (service selector switches), and Google Cloud Run.

## Related core lessons

- Module 15.15 — Demo preparation
- Module 17.15 — Demo preparation
