# Project selection and scoping

## Learning objective

Select a portfolio project that matches your skill level and goals, define its MVP scope, estimate timeline and complexity, and avoid common scope creep traps.

## Why this matters

The most common reason portfolio projects fail is not lack of skill — it is poor scoping. Engineers start too big, burn out, and abandon the project three weeks in. Others start too small and produce something that does not impress. Learning to scope a project correctly means you ship a complete, polished project that demonstrates your best work within a realistic timeframe. This skill transfers directly to your day job, where scoping is a daily activity.

## Mental model

Think of a project scope like a camera lens. A wide-angle lens (broad scope) captures everything but with less detail. A zoom lens (narrow scope) captures one thing sharply. The best portfolio projects use a zoom lens: they pick one valuable feature set and execute it exceptionally well rather than building ten features poorly. Scope is the art of choosing what not to build.

## Core idea

Project scoping breaks down into four dimensions:

**Feature scope**: How many features the project includes. More features mean more code, more testing, more documentation, and more time.

**Technical complexity**: The difficulty of the implementation. A REST API is simpler than a distributed consensus system. Adding authentication, databases, external APIs, and UIs multiplies complexity.

**Team size**: A solo project has zero coordination overhead but limited throughput. A team project introduces communication overhead that slows velocity.

**Time horizon**: How many weeks or months you have. A two-week project must be scoped differently from a three-month project.

**MVP definition**: The Minimum Viable Product is the smallest set of features that still demonstrates the core value. For a portfolio project, the MVP must demonstrate enough engineering depth to impress. An MVP for a URL shortener might be: create short links, redirect to original URLs, track click counts, with tests and a README. Nice-to-haves like custom aliases, analytics dashboards, and team management come after.

**Complexity estimation framework**:

| Factor | Weight |
|---|---|
| Base project | 2 weeks |
| Per feature | +1.5 weeks |
| Per integration | +2 weeks |
| Authentication | +2 weeks |
| Database | +1.5 weeks |
| External API | +1.5 weeks |
| UI/frontend | +3 weeks |
| Per additional team member | +15% overhead |

## Under the hood

The estimator uses a linear additive model with a team overhead multiplier. This is a simplification — real project estimation uses techniques like:

- **Three-point estimation**: Best case, worst case, most likely case. The weighted average gives a more realistic timeline.
- **PERT**: Program Evaluation and Review Technique uses `(Optimistic + 4*MostLikely + Pessimistic) / 6`.
- **Historical data**: Past project velocity is the best predictor of future velocity.

The simple model here works well for small, solo projects. For larger efforts, use a spreadsheet or dedicated estimation tool.

## How Go uses it

Go's standard library is extensive enough to build many portfolio projects without external dependencies. This simplifies scoping because you spend zero time on framework selection. You can build:

- An HTTP API with `net/http` instead of a framework.
- A CLI with `flag` instead of cobra (or use cobra for more complex needs).
- A database-backed service with `database/sql`.
- A concurrent worker with goroutines and channels.

When scoping a Go project, consider which standard library packages you can use and where you genuinely need an external dependency. Each dependency adds integration complexity and maintenance cost.

## Go example

The project scoping tool estimates complexity and timeline based on project characteristics. It helps you make data-driven scoping decisions before committing to a project.

```go
package main

import (
	"fmt"
	"math"
)

type ComplexityLevel int

const (
	Trivial ComplexityLevel = iota + 1
	Easy
	Moderate
	Complex
	VeryComplex
)

func (c ComplexityLevel) String() string {
	switch c {
	case Trivial:
		return "Trivial"
	case Easy:
		return "Easy"
	case Moderate:
		return "Moderate"
	case Complex:
		return "Complex"
	case VeryComplex:
		return "Very Complex"
	default:
		return "Unknown"
	}
}

type ProjectScope struct {
	Name             string
	Description      string
	NumFeatures      int
	NumIntegrations  int
	HasAuth          bool
	HasDatabase      bool
	HasExternalAPI   bool
	NeedsUI          bool
	TeamSize         int
	EstimatedWeeks   float64
	Complexity       ComplexityLevel
	RiskFactors      []string
}

type ComplexityEstimator struct {
	BaseWeeks     float64
	FeatureWeight float64
	IntegrationW  float64
	AuthWeight    float64
	DatabaseW     float64
	APIWeight     float64
	UIWeight      float64
	TeamOverhead  float64
}

func DefaultEstimator() ComplexityEstimator {
	return ComplexityEstimator{
		BaseWeeks:     2.0,
		FeatureWeight: 1.5,
		IntegrationW:  2.0,
		AuthWeight:    2.0,
		DatabaseW:     1.5,
		APIWeight:     1.5,
		UIWeight:      3.0,
		TeamOverhead:  0.15,
	}
}

func EstimateProject(scope ProjectScope, est ComplexityEstimator) ProjectScope {
	weeks := est.BaseWeeks
	weeks += float64(scope.NumFeatures) * est.FeatureWeight
	weeks += float64(scope.NumIntegrations) * est.IntegrationW

	if scope.HasAuth {
		weeks += est.AuthWeight
		scope.RiskFactors = append(scope.RiskFactors, "authentication adds security complexity")
	}
	if scope.HasDatabase {
		weeks += est.DatabaseW
		scope.RiskFactors = append(scope.RiskFactors, "database schema design and migration overhead")
	}
	if scope.HasExternalAPI {
		weeks += est.APIWeight
		scope.RiskFactors = append(scope.RiskFactors, "external API integration has reliability dependencies")
	}
	if scope.NeedsUI {
		weeks += est.UIWeight
		scope.RiskFactors = append(scope.RiskFactors, "UI development adds frontend complexity")
	}

	if scope.TeamSize > 0 {
		weeks *= 1.0 + est.TeamOverhead*float64(scope.TeamSize-1)
	}

	estimatedWeeks := math.Ceil(weeks)
	scope.EstimatedWeeks = estimatedWeeks

	switch {
	case estimatedWeeks <= 2:
		scope.Complexity = Trivial
	case estimatedWeeks <= 4:
		scope.Complexity = Easy
	case estimatedWeeks <= 8:
		scope.Complexity = Moderate
	case estimatedWeeks <= 16:
		scope.Complexity = Complex
	default:
		scope.Complexity = VeryComplex
	}

	return scope
}

func (s ProjectScope) Summary() string {
	return fmt.Sprintf(
		"Project: %s\n  Description: %s\n  Features: %d, Integrations: %d\n  Auth: %v, DB: %v, API: %v, UI: %v\n  Team Size: %d\n  Estimated Complexity: %s\n  Estimated Timeline: %.0f weeks\n  Risk Factors: %v\n",
		s.Name, s.Description, s.NumFeatures, s.NumIntegrations,
		s.HasAuth, s.HasDatabase, s.HasExternalAPI, s.NeedsUI,
		s.TeamSize, s.Complexity, s.EstimatedWeeks, s.RiskFactors,
	)
}

func main() {
	est := DefaultEstimator()
	projects := []ProjectScope{
		{
			Name: "URL Shortener API",
			Description: "A REST API that shortens URLs with redirect tracking and analytics",
			NumFeatures: 4, NumIntegrations: 1,
			HasAuth: true, HasDatabase: true, HasExternalAPI: false,
			NeedsUI: false, TeamSize: 1,
		},
		{
			Name: "CLI Task Manager",
			Description: "A command-line task management tool with local storage",
			NumFeatures: 3, NumIntegrations: 0,
			HasAuth: false, HasDatabase: true, HasExternalAPI: false,
			NeedsUI: false, TeamSize: 1,
		},
		{
			Name: "Full SaaS Dashboard",
			Description: "A full-featured SaaS dashboard with multi-tenant auth and analytics",
			NumFeatures: 8, NumIntegrations: 3,
			HasAuth: true, HasDatabase: true, HasExternalAPI: true,
			NeedsUI: true, TeamSize: 3,
		},
	}

	for _, p := range projects {
		result := EstimateProject(p, est)
		fmt.Println(result.Summary())
	}
}
```

## Step-by-step execution

1. Define a `ProjectScope` struct with all relevant project characteristics: features, integrations, team size, and boolean flags for auth, database, API, and UI.
2. `DefaultEstimator()` returns a `ComplexityEstimator` with sensible default weights based on industry experience building Go projects.
3. `EstimateProject` takes a `ProjectScope` and an `Estimator`. It starts with the base weeks and adds weighted contributions for each factor.
4. Risk factors accumulate as each technical decision is added, providing a natural language explanation of the timeline drivers.
5. Team overhead is applied as a multiplier: a team of three has 30% overhead (2 members x 15%).
6. The final weeks are ceiling-rounded and mapped to a `ComplexityLevel` from Trivial to Very Complex.
7. `Summary()` formats the complete estimate for review.

## Common mistakes

- Mistake: Building a project that requires learning Go, a database, Docker, Kubernetes, gRPC, and React simultaneously.
  - Why it happens: Engineers want to learn everything at once.
  - Fix: Scope the project to use at most one or two new technologies. Mastery requires focus.

- Mistake: Defining MVP as "all features but buggy" instead of "core features finished well."
  - Why it happens: It is tempting to keep adding features because each one seems small.
  - Fix: Write the MVP feature list on a sticky note. Do not add anything until the list is complete and polished.

- Mistake: Ignoring the cost of integrations.
  - Why it happens: An API call looks like one line of code.
  - Fix: Every integration requires error handling, retries, testing, and documentation. Count that cost.

- Mistake: Estimating optimistically because you will "work harder later."
  - Why it happens: Optimism bias is human nature.
  - Fix: Estimate based on your past velocity, not your ideal velocity.

## Debugging walkthrough

An engineer scopes a "Simple Blog Engine" with 6 features, a database, authentication, and a UI, with 1 team member. The estimator returns:

```
Estimated Complexity: Very Complex
Estimated Timeline: 19 weeks
```

The engineer expected 6-8 weeks. Checking the breakdown: 3 features (4.5 weeks) + database (1.5 weeks) + auth (2 weeks) + UI (3 weeks) + integrations (0) + base (2 weeks) = 13 weeks. The "+1 integration" added 2 more, and the team size of 1 meant no overhead. The total is 15 weeks, rounded to 19. The engineer realizes the UI weight is the biggest driver and decides to build a CLI-only version first, reducing the estimate to 11 weeks.

## Production notes

In a real project, scoping is iterative. You estimate, start building, learn, and re-estimate. The estimator here is a starting point. Adjust the weights based on your own experience:

- If you are new to Go, multiply all estimates by 1.5.
- If you are experienced, the default weights are reasonable.
- If the project includes a technology you have never used, add 2-4 weeks learning time.

The risk factors section is as important as the timeline. Each risk factor should trigger a mitigation strategy: "database schema design" means you should spend a day on schema modeling before writing code.

## Performance implications

The scoping tool itself has negligible performance cost — it runs in microseconds. However, the concepts it models have real performance implications:

- Adding a database introduces query latency you must design for.
- Adding external APIs introduces network latency and failure modes.
- Adding a UI means your backend must serve a frontend, potentially affecting API response times.
- Each team member adds communication overhead but can increase throughput for independent work.

The scoping tool helps you anticipate these costs before they affect your project.

## Practice task

Add a `LearningCurve` field to `ProjectScope` that models how many new technologies the project requires. Each new technology adds 2 weeks to the estimate. Update `EstimateProject` to include this factor. Add a corresponding test that verifies a project with 3 new technologies gets 6 extra weeks.

## Tests / verification

```bash
go test ./curriculum/modules/16-portfolio-interview-readiness/lessons/02-project-selection-and-scoping
go run ./curriculum/modules/16-portfolio-interview-readiness/lessons/02-project-selection-and-scoping
```

The tests verify that trivial projects estimate 1-4 weeks, complex projects estimate longer than simple ones, complexity levels are correctly assigned, and risk factors accumulate. After the practice task, verify that the learning curve factor is correctly applied.

## Review questions

1. What four dimensions define a project scope, and how do they interact?
2. Why does adding a UI add more time than adding a database to a Go project?
3. How does team size affect project timeline, and why is the relationship not linear?
4. What is the difference between MVP and "all features but buggy?"
5. If you have 8 weeks to build a portfolio project, which factors should you prioritize and which should you eliminate?

## NEXT UP

Code review readiness — prepare your code for professional review by understanding what reviewers look for and how to self-review before submitting.
