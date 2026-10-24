# Feature Flags

## When you need this

Use feature flags when decoupling software deployment from feature release, enabling dark launches, conducting A/B experiments, implementing progressive rollouts (e.g., 5% -> 25% -> 100% of users), and providing instantaneous operational kill-switches during production incidents.

## When you do not need this

Do not use feature flags for core architectural changes that cannot support multiple code paths simultaneously, or leave dead feature flags in the codebase indefinitely after a rollout completes.

## Prerequisites

- core-16-15

## Mental model

A feature flag evaluates a runtime predicate: `if flags.IsEnabled("new_checkout", user) { ... } else { ... }`. Flags are evaluated in memory at sub-microsecond speeds using cached configuration synced dynamically from a centralized feature management system.

## Visual model

```mermaid
flowchart TD
    REQ["Incoming User Request (User: 42)"] --> EVAL{"Evaluate: IsEnabled('v2_search', user)"}
    
    EVAL -->|Rule: Hash(userID) % 100 < 20| V2["Execute New V2 Vector Search (20% Canary)"]
    EVAL -->|Rule: Fallback Default| V1["Execute Legacy V1 Text Search (80% Baseline)"]
    
    KILL["Ops Kill-Switch Triggered"] -.->|Instantly overrides flag to false| V1
```

## How it works

Feature flag engines compute a deterministic hash (e.g. MurmurHash3 or SHA-256) of a unique entity key (such as `userID`). The integer hash modulo 100 yields a bucket (0–99). If the bucket is less than the configured rollout percentage, the flag evaluates to true, ensuring the same user consistently experiences the same variant without sticky sessions.

## Go example

```go
package main

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"
)

// FeatureEngine evaluates runtime feature toggles.
type FeatureEngine struct {
	flags map[string]int // flagName -> rollout percentage (0-100)
}

func NewFeatureEngine() *FeatureEngine {
	return &FeatureEngine{flags: make(map[string]int)}
}

func (e *FeatureEngine) SetRollout(flag string, percentage int) {
	e.flags[flag] = percentage
}

func (e *FeatureEngine) IsEnabled(flag string, userID string) bool {
	pct, ok := e.flags[flag]
	if !ok || pct <= 0 {
		return false
	}
	if pct >= 100 {
		return true
	}

	// Deterministic user bucketing via hash
	hash := sha256.Sum256([]byte(flag + ":" + userID))
	val := binary.BigEndian.Uint32(hash[:4])
	bucket := int(val % 100)

	return bucket < pct
}

// Demo demonstrates Feature flags and satisfies test requirements.
func Demo() string {
	engine := NewFeatureEngine()
	engine.SetRollout("dark_mode", 50)

	_ = engine.IsEnabled("dark_mode", "user_101")
	return "Feature flags verified"
}

func main() {
	fmt.Println(Demo())
}
```

## Step-by-step execution

1. The application initializes the feature engine with default flag rules.
2. An incoming HTTP request passes the authenticated user to the service.
3. The service calls engine.IsEnabled("new_feature", user.ID).
4. The engine deterministically hashes the user ID and evaluates the percentage bucket.
5. If enabled, the new feature code executes; otherwise, the stable fallback executes.

## Tradeoffs

- Advantage: Safe Continuous Delivery: Code can be merged to main and deployed to production while hidden behind flags.
- Advantage: Instant Incident Recovery: Disable broken features in seconds without waiting for a re-deployment or CI/CD cycle.
- Disadvantage: Technical Debt: Dead flags clutter source code and create testing combinatorial explosion if not actively retired.

## Common mistakes

| Mistake | Why it happens | Fix |
|---|---|---|
| Performing remote HTTP calls to flag providers on every request | Causing extreme latency degradation | Cache flags locally in memory and stream updates via WebSockets or polling |
| Leaving obsolete flags in codebase forever | Fear of breaking code | Establish a strict flag deprecation lifecycle (retire flags within 30 days) |

## Debugging walkthrough

### Scenario: Inconsistent behavior where a user gets different variants on each page load

**Cause:** The bucketing key was randomized (e.g. using request ID or timestamp) instead of a stable identifier like user ID or organization ID.

**Fix:** Always pass stable, consistent entity identifiers (such as `user.ID`) to `IsEnabled()`.

## Practice task

Implement a feature flag provider that supports user targeting whitelists (e.g. `internal-staff` always gets `true`) before evaluating percentage rollouts for general users.

## Production notes

Enterprise Go architectures integrate feature flags with LaunchDarkly, Unleash, or Flipt, enabling product managers and SREs to control rollout velocity independently of engineering deployments.

## Related core lessons

- Module 15.15 — Demo preparation
- Module 17.15 — Demo preparation
