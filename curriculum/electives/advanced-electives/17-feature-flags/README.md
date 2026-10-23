# Feature flags

## When you need this

Deploying a new feature requires merging code and redeploying, making it impossible to test in production or roll back without another deployment. Feature flags solve this by allowing features to be toggled at runtime without redeployment — enabling gradual rollouts, A/B testing, and instant kill switches.

## When you do not need this

When simpler concrete types or standard core library primitives suffice without adding abstraction overhead.

## Prerequisites

- elective-16

## Mental model

A feature flag is a conditional branch controlled by an external configuration system. The flag evaluates to true or false based on context (user ID, request attributes, environment, random percentage). The code path for the new feature is wrapped in an `if flag.Enabled(ctx, "feature-name")` check. The flag can be changed at runtime: enabling the feature for 1% of users, disabling it if errors spike, and removing the flag entirely once the feature is stable.

## How it works

A feature flag system evaluates flag values based on rules: percentage rollout (enabled for X% of users based on hash(user_id) % 100), user targeting (enabled for specific user IDs or email domains), environment (enabled in staging, disabled in production), and gradual rollout (enabled for increasing percentages over time). Flags are stored in a configuration service (etcd, Consul, LaunchDarkly) and cached locally with a TTL. The local cache avoids per-request network calls. Flag removal is a code cleanup step: after the feature is validated at 100% for a stable period, the flag check and old code path are removed.

## Go example

```go
package main

import "fmt"

func main() {
    fmt.Println("Elective: Feature flags")
}
```

## Tradeoffs

- Flag evaluation adds latency: each check is typically a config lookup (map access or RPC) — 100ns-1ms per flag
- Caching flag values: TTL-based cache reduces evaluation to an atomic load — ~2ns vs uncached ~500ns
- Dead flag branches compile into the binary: even disabled flags have code compiled in — use build tags for truly experimental features
- Removing a flag removes dead code paths and simplifies the call graph — measurable improvement in binary size and compile time

## Common mistakes

- Not cleaning up flags after rollout — dead flag branches accumulate, making code harder to read and maintain.
- Using flags for configuration — flags are for features, not for configuration values (database DSN, port number).
- Not monitoring flag usage — without tracking which flags are checked and their values, operator confusion about which features are active.
- Using flags that are checked in too many places — a flag should control a single feature in a single layer, not scattered across handlers, services, and frontend.
- Not rolling back flags — if the new feature causes data corruption, disabling the flag does not undo the corruption (compensating actions needed).

## Practice task

The learner must implement a feature flag system supporting percentage-based rollouts and user targeting, wrap a new code path in a flag check, demonstrate that the flag can be toggled at runtime without redeployment, and document the flag removal process.

## Production notes

Feature flags are used in every continuous delivery pipeline. Companies like Netflix, Facebook, and Google use flags for every feature launch. LaunchDarkly is the leading SaaS feature flag platform. Open-source alternatives: Unleash, Flipt. Go services use flags for gradual rollouts, A/B testing, operational kill switches, and environment-specific behavior.

## Related core lessons

- Module 16.16 — CQRS
- Module 17.15 — sync.Pool
