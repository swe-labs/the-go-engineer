# One deployment target

## Learning objective

Model deployment strategies (blue-green, canary, rolling, recreate) in Go, implement deploy and rollback logic for each strategy, and determine which strategies provide zero-downtime deployments.

## Why this matters

Deploying a new version of a Go service without downtime is a core operational requirement. Users expect 24/7 availability; a deployment that takes the site down for seconds costs revenue and trust. Choosing the right deployment strategy depends on the service's architecture, traffic patterns, and risk tolerance. A single-binary Go service is particularly well-suited to immutable deployment strategies because the binary is self-contained and statically linked. Understanding blue-green, canary, and rolling deployments lets you ship updates continuously without disruption.

## Mental model

A deployment target is the environment where your Go binary runs. In immutable infrastructure, you never modify a running server. Instead, you create a new server with the new binary and redirect traffic to it. The old server remains as a rollback target.

Think of it like replacing an aircraft engine in flight: you do not shut down the engine and replace it while flying. You bring a second engine online, verify it works, then switch over. If the new engine fails, you switch back instantly.

- **Blue-green**: Two identical environments (blue = old, green = new). Switch traffic instantly.
- **Canary**: Route a small percentage of traffic to the new version. Monitor, then gradually increase.
- **Rolling**: Replace instances one at a time. Each instance is taken out of service, updated, and returned.
- **Recreate**: Stop all instances, start new ones. Requires downtime.

## Core idea

| Strategy | Traffic switch | Rollback speed | Zero downtime | Complexity |
|---|---|---|---|---|
| Recreate | Immediate (all at once) | Slow (redeploy old) | No | Low |
| Rolling | Gradual (per instance) | Medium (per instance) | Yes | Medium |
| Blue-green | Instant (load balancer) | Instant (flip back) | Yes | High (double resources) |
| Canary | Gradual (percentage) | Instant (reset traffic) | Yes | High (monitoring needed) |

For Go services, blue-green is the most common zero-downtime strategy. Since the Go binary is statically compiled, the green environment can be fully provisioned and tested before switching traffic. The blue environment is kept warm for instant rollback.

## Under the hood

Blue-green deployment at the infrastructure level:

1. The load balancer (nginx, ELB, HAProxy) has two target groups: blue and green.
2. Initially, all traffic goes to blue (`TargetGroupA`).
3. A new green instance is provisioned with the new Go binary. Health checks pass.
4. The load balancer switches the listener rule to point to `TargetGroupB` (green).
5. Connection draining on the blue target group completes outstanding requests.
6. Blue instances are terminated or kept as rollback.

Canary deployment adds a traffic splitting layer:

1. The load balancer routes 10% of requests to the new version, 90% to the old version.
2. Error rate, latency, and throughput are monitored for a stabilization period.
3. If metrics are healthy, traffic is increased to 25%, then 50%, then 100%.
4. If metrics degrade, traffic is immediately reverted to 0% on canary.

## How Go uses it

Go's single-binary deployment model makes blue-green and canary strategies particularly effective:

- The binary is self-contained with no runtime dependencies. No need to install Python, Ruby, or a JVM.
- A simple health endpoint (`GET /healthz`) returns 200 when the server is ready.
- Graceful shutdown (`signal.Notify` for SIGTERM) lets the load balancer drain connections before terminating.

Typical Go deployment flow:

```go
func main() {
    go func() {
        sigCh := make(chan os.Signal, 1)
        signal.Notify(sigCh, syscall.SIGTERM)
        <-sigCh
        // Stop accepting new requests, drain existing ones
        server.Shutdown(context.Background())
    }()
    log.Fatal(server.ListenAndServe())
}
```

## Go example

```go
package main

import (
	"fmt"
)

type DeploymentStrategy string

const (
	StrategyBlueGreen DeploymentStrategy = "blue-green"
	StrategyCanary    DeploymentStrategy = "canary"
	StrategyRolling   DeploymentStrategy = "rolling"
	StrategyRecreate  DeploymentStrategy = "recreate"
)

type Target struct {
	Name    string
	Version string
	Healthy bool
}

type Deployment struct {
	Strategy      DeploymentStrategy
	Current       *Target
	New           *Target
	CanaryPercent int
	Done          bool
}

func NewBlueGreenDeployment(currentVersion, newVersion string) *Deployment {
	return &Deployment{
		Strategy: StrategyBlueGreen,
		Current:  &Target{Name: "blue", Version: currentVersion, Healthy: true},
		New:      &Target{Name: "green", Version: newVersion, Healthy: false},
	}
}

func NewCanaryDeployment(currentVersion, newVersion string, pct int) *Deployment {
	return &Deployment{
		Strategy:      StrategyCanary,
		Current:       &Target{Name: "production", Version: currentVersion, Healthy: true},
		New:           &Target{Name: "canary", Version: newVersion, Healthy: false},
		CanaryPercent: pct,
	}
}

func (d *Deployment) Deploy() ([]string, error) {
	var events []string
	switch d.Strategy {
	case StrategyBlueGreen:
		events = append(events,
			fmt.Sprintf("Spinning up %s (%s)", d.New.Name, d.New.Version),
			"Running health checks")
		d.New.Healthy = true
		events = append(events,
			"Health checks passed, switching traffic",
			fmt.Sprintf("Traffic switched to %s", d.New.Name),
			fmt.Sprintf("Keeping %s as rollback", d.Current.Name))
	case StrategyCanary:
		events = append(events,
			fmt.Sprintf("Deploying canary (%s) with %d%% traffic", d.New.Version, d.CanaryPercent),
			"Monitoring error rates...",
			"Canary healthy, rolling out to 100%")
		d.New.Healthy = true
	default:
		return nil, fmt.Errorf("unknown strategy: %s", d.Strategy)
	}
	d.Done = true
	return events, nil
}

func (d *Deployment) Rollback() []string {
	if !d.Done {
		return []string{"Deployment not complete, nothing to rollback"}
	}
	switch d.Strategy {
	case StrategyBlueGreen:
		return []string{fmt.Sprintf(
			"Switching traffic back to %s (%s)", d.Current.Name, d.Current.Version)}
	default:
		return []string{fmt.Sprintf("Redeploying version %s", d.Current.Version)}
	}
}

func ZeroDowntimeRequired(s DeploymentStrategy) bool {
	return s == StrategyBlueGreen || s == StrategyCanary
}

func main() {
	blueGreen := NewBlueGreenDeployment("v1.0.0", "v2.0.0")
	events, _ := blueGreen.Deploy()
	for _, e := range events {
		fmt.Println("  ", e)
	}
	fmt.Println("Zero downtime:", ZeroDowntimeRequired(blueGreen.Strategy))
}
```

## Step-by-step execution

For a blue-green deployment from v1.0.0 to v2.0.0:

1. `NewBlueGreenDeployment("v1.0.0", "v2.0.0")` creates a deployment with blue=v1 (healthy) and green=v2 (not healthy).
2. `Deploy()` is called with strategy `blue-green`.
3. Step: "Spinning up green environment (v2.0.0)" — simulates provisioning a new server or container with v2.
4. Step: "Running health checks" — the new environment passes its `/healthz` endpoint check.
5. `d.New.Healthy = true` — the new target is marked healthy.
6. Step: "Health checks passed, switching traffic" — the load balancer is updated.
7. Step: "Traffic switched to green" — all requests now go to v2.0.0.
8. Step: "Keeping blue as rollback" — the blue environment remains running.

For rollback:
1. `Rollback()` is called.
2. Step: "Switching traffic back to blue (v1.0.0)" — load balancer switches back.
3. The current and new targets are swapped.

## Common mistakes

- **Not health-checking before switching traffic**: If the new version fails health checks and traffic is switched anyway, all requests fail. Always wait for a passing health check.
- **Terminating the old environment too quickly**: In blue-green, keep the old environment running for at least 10-15 minutes after switching to allow for quick rollback if issues are discovered.
- **Not draining connections on shutdown**: When a server is terminated, in-flight requests fail. Use graceful shutdown with a drain timeout (typically 30 seconds).
- **Canary without monitoring**: Deploying a canary without monitoring error rates and latency defeats the purpose. The canary must be observable to make a go/no-go decision.
- **Forgetting database migrations**: Zero-downtime code deployment does not help if the database migration breaks the old version. Migrations must be backward-compatible.

## Debugging walkthrough

Consider a blue-green deployment where traffic is switched but users report 502 errors.

**Symptom**: After switching traffic to green, all requests return HTTP 502.

**Investigation**:
1. Check the green environment's health endpoint: `curl http://green:8080/healthz`.
2. If it returns non-200, the health check step in the deployment logic did not actually verify before switching.
3. Check the Go server logs on the green instance for startup errors (port already in use, database connection failure).

**Root cause**: The health check was a shell command that returned 0 (success) even though the server was not listening. For example, `curl --fail http://localhost:8080/healthz` might succeed if curl defaults to connecting to a different port.

**Fix**: Update the health check to verify the correct port and expect a 200 response:

```go
func healthCheck(url string) bool {
    resp, err := http.Get(url)
    return err == nil && resp.StatusCode == http.StatusOK
}
```

Another scenario: The canary deployment reports elevated error rates, but the deployment script continues to increase traffic.

**Root cause**: The deployment script does not have automated rollback on metrics degradation. The canary decision was manual.

**Fix**: Implement automated canary analysis: if error rate increases by more than 1% or p99 latency increases by more than 100ms, automatically rollback.

## Production notes

In production deployment pipelines:

- **Immutable artifacts**: The same binary that passed CI is deployed to every environment. Never rebuild for a specific environment.
- **Blue-green requires double capacity**: You need enough resources to run two full environments simultaneously. For cost-sensitive services, rolling updates are more practical.
- **Canary analysis automation**: Use tools like Flagger, Argo Rollouts, or Google Cloud Deploy to automate canary analysis and promotion.
- **Feature flags alongside deployment**: Separate deployment (new binary) from release (feature enabled). Use feature flags to gradually expose new functionality without redeploying.
- **Database schema compatibility**: The old version must work with the new schema during a rolling or blue-green deployment. This means additive-only schema changes (add columns, never remove or rename).

## Performance implications

- **Blue-green doubles infrastructure cost** during the transition period (both environments running). After verification, the old environment is terminated.
- **Canary deployments add latency to user requests** if the traffic split happens at the application layer. Load balancer-level splitting has negligible overhead.
- **Rolling updates minimize cost** but take longer (N instances × startup time). For 10 instances with 30-second startup, a rolling update takes 5 minutes.
- **Health check frequency**: Every health check adds load. Checking every 5 seconds per instance is standard. Rapid health checks (every 1 second) during deployment are acceptable.

## Practice task

Write a function `SelectStrategy(needsZeroDowntime bool, hasExtraCapacity bool, hasMonitoring bool) DeploymentStrategy` that recommends a deployment strategy based on requirements. Then write a `func (d *Deployment) Describe() string` that produces a human-readable summary of the deployment plan. Write a `main()` that creates deployments with each strategy and prints their descriptions.

## Tests / verification

```bash
go run ./curriculum/modules/15-docker-cicd-deployment/lessons/17-one-deployment-target
go test ./curriculum/modules/15-docker-cicd-deployment/lessons/17-one-deployment-target
```

The existing tests verify blue-green and canary deployment creation, deploy execution, health state changes, rollback behavior (including before deploy), unknown strategy errors, and zero-downtime strategy detection. After completing the practice task, add tests for `SelectStrategy` covering all strategy combinations.

## Review questions

1. What is immutable infrastructure, and why is a single Go binary well-suited for it?
2. How does blue-green deployment achieve zero downtime? What is the cost trade-off?
3. In a canary deployment, what metrics should be monitored to decide whether to proceed or rollback?
4. Why must database schema changes be backward-compatible during a zero-downtime deployment?
5. A rolling update replaces instances one at a time. How does this differ from blue-green in terms of resource usage and risk?

## NEXT UP

Rollback basics: strategies for reverting to a previous version when a deployment goes wrong.
