# Portfolio strategy

## Learning objective

Build a standout Go engineering portfolio by selecting, analyzing, and presenting projects that demonstrate technical depth, clean code practices, and real-world impact.

## Why this matters

A GitHub profile is the new resume. Every hiring manager and technical screener will look at your portfolio before or during the interview process. A strong portfolio proves you can ship production-quality Go code, collaborate on open source, and solve meaningful problems. Without a deliberate strategy, your projects can appear scattered, unfinished, or shallow. A portfolio strategy turns your work into a coherent narrative that lands you the interview.

## Mental model

Think of your portfolio as a product you are building for a specific audience: hiring managers, technical screeners, and your future teammates. Each project is a feature that showcases one or more dimensions of your engineering ability. Just as you would prioritize features by impact, you prioritize projects by how strongly they signal the skills your target roles demand. The portfolio strategy is your roadmap: it tells you which projects to build, which to retire, and how to present each one for maximum signal.

## Core idea

There are three portfolio types that engineers mix and match:

**Project portfolio**: 3-5 polished side projects that solve real problems. Each project has a clear README, tests, CI, and documentation. These demonstrate your ability to ship end-to-end.

**Problem-solving portfolio**: Solutions to algorithmic challenges, competitive programming, or coding challenge platforms. These show raw problem-solving ability but do not demonstrate production engineering skills on their own.

**Open source portfolio**: Contributions to established Go projects. These prove you can work within existing codebases, follow contributor guidelines, and collaborate with other engineers.

A balanced portfolio includes at least one project from each category but emphasizes the project portfolio because it signals the most relevant skills for a production engineering role.

**Project selection criteria**:

| Criterion | Why it matters | Signal strength |
|---|---|---|
| Complexity | Shows you can handle non-trivial architecture | Medium |
| Technology relevance | Aligns with industry demand (gRPC, Kubernetes, distributed systems) | High |
| Testing | Proves you value correctness and maintainability | High |
| Documentation | Shows communication skills and developer empathy | High |
| Real-world use | Adoption by others validates the idea | Very high |
| Maintainability | CI/CD, clean code, modular design | Medium |

## Under the hood

Hiring managers spend an average of 30 seconds scanning a GitHub profile before deciding whether to dive deeper. They look for:

1. A pinned repository with a clear README that explains the problem, architecture, and how to run it.
2. Evidence of testing (a `_test.go` file visible in the file list).
3. A green CI badge or workflow file.
4. Recent activity — stale repos suggest the engineer stopped learning.
5. Contributions to other repos — signals collaboration.

The Go community specifically values projects that use idiomatic Go patterns: `context` for cancellation, `error` handling, `sync` primitives, and the standard library. A project that showcases these patterns signals that you understand Go conventions.

## How Go uses it

Go's tooling and ecosystem make it particularly well-suited for portfolio projects:

- `go test` built-in — no excuses for untested projects.
- `gofmt` enforces consistent style automatically.
- `go mod` provides clear dependency management.
- `godoc` renders documentation from comments.
- `go vet` and `staticcheck` catch common mistakes.

These tools mean a Go portfolio project can demonstrate professional engineering practices with minimal setup. The barrier to a production-quality project is lower in Go than in most languages.

The standard library also enables building meaningful projects without heavy frameworks: `net/http` for APIs, `database/sql` for persistence, `encoding/json` for data exchange, and `testing` for verification.

## Go example

The following tool analyzes a project across six dimensions and generates a score. It demonstrates the kind of tool you might build as part of a portfolio to self-evaluate project quality.

```go
package main

import (
	"fmt"
	"strings"
)

type Dimension int

const (
	Complexity Dimension = iota
	TechRelevance
	Testing
	Documentation
	RealWorldUse
	Maintainability
)

func (d Dimension) String() string {
	switch d {
	case Complexity:
		return "Complexity"
	case TechRelevance:
		return "Tech Relevance"
	case Testing:
		return "Testing"
	case Documentation:
		return "Documentation"
	case RealWorldUse:
		return "Real-World Use"
	case Maintainability:
		return "Maintainability"
	default:
		return "Unknown"
	}
}

type Project struct {
	Name        string
	Description string
	Language    string
	Stars       int
	HasTests    bool
	HasDocs     bool
	HasCI       bool
	LinesOfCode int
}

type Score struct {
	Dimension Dimension
	Score     int
	Reason    string
}

type ProjectAnalysis struct {
	Project Project
	Scores  []Score
	Total   int
}

func AnalyzeProject(p Project) ProjectAnalysis {
	analysis := ProjectAnalysis{Project: p}

	compScore := 0
	reason := "basic"
	switch {
	case p.LinesOfCode > 50000:
		compScore = 3
		reason = "large codebase shows ability to navigate complexity"
	case p.LinesOfCode > 10000:
		compScore = 2
		reason = "moderate codebase with meaningful structure"
	default:
		compScore = 1
		reason = "small codebase, manageable scope"
	}
	analysis.Scores = append(analysis.Scores, Score{Complexity, compScore, reason})

	techScore := 4
	techReason := "uses Go, a sought-after systems language"
	if strings.Contains(strings.ToLower(p.Description), "grpc") ||
		strings.Contains(strings.ToLower(p.Description), "microservice") ||
		strings.Contains(strings.ToLower(p.Description), "distributed") ||
		strings.Contains(strings.ToLower(p.Description), "cloud") {
		techScore = 5
		techReason = "uses Go with modern architecture patterns"
	}
	analysis.Scores = append(analysis.Scores, Score{TechRelevance, techScore, techReason})

	testScore := 1
	testReason := "no tests detected"
	if p.HasTests {
		testScore = 4
		testReason = "tests present, demonstrates quality mindset"
	}
	analysis.Scores = append(analysis.Scores, Score{Testing, testScore, testReason})

	docScore := 1
	docReason := "no documentation detected"
	if p.HasDocs {
		docScore = 4
		docReason = "documented, shows communication skills"
	}
	analysis.Scores = append(analysis.Scores, Score{Documentation, docScore, docReason})

	realScore := 2
	realReason := "personal or experimental project"
	if p.Stars > 100 {
		realScore = 5
		realReason = "significant community adoption, real-world impact"
	} else if p.Stars > 10 {
		realScore = 4
		realReason = "some adoption, demonstrates value"
	} else if p.Stars > 0 {
		realScore = 3
		realReason = "initial traction, potential for growth"
	}
	analysis.Scores = append(analysis.Scores, Score{RealWorldUse, realScore, realReason})

	maintainScore := 1
	maintainReason := "no CI detected"
	if p.HasCI {
		maintainScore = 4
		maintainReason = "CI configured, shows engineering discipline"
	}
	analysis.Scores = append(analysis.Scores, Score{Maintainability, maintainScore, maintainReason})

	for _, s := range analysis.Scores {
		analysis.Total += s.Score
	}

	return analysis
}

func (a ProjectAnalysis) Summary() string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("Analysis for %s:\n", a.Project.Name))
	for _, s := range a.Scores {
		b.WriteString(fmt.Sprintf("  %-18s %d/5 — %s\n", s.Dimension.String()+":", s.Score, s.Reason))
	}
	b.WriteString(fmt.Sprintf("  Total score: %d/%d\n", a.Total, len(a.Scores)*5))
	return b.String()
}

func main() {
	projects := []Project{
		{
			Name:        "go-url-shortener",
			Description: "A URL shortening service built with Go, Redis, and gRPC",
			Language:    "Go",
			Stars:       245,
			HasTests:    true,
			HasDocs:     true,
			HasCI:       true,
			LinesOfCode: 8500,
		},
		{
			Name:        "hello-world-cli",
			Description: "Simple CLI tool in Go",
			Language:    "Go",
			Stars:       3,
			HasTests:    false,
			HasDocs:     false,
			HasCI:       false,
			LinesOfCode: 120,
		},
	}

	for _, p := range projects {
		analysis := AnalyzeProject(p)
		fmt.Println(analysis.Summary())
	}
}
```

## Step-by-step execution

1. Define a `Project` struct with metadata including name, description, stars, and boolean flags for tests, docs, and CI.
2. `AnalyzeProject` receives a `Project` and builds a `ProjectAnalysis` containing six dimension scores.
3. Each dimension has its own scoring logic: complexity uses lines of code, tech relevance scans the description for keywords, and real-world use uses GitHub stars.
4. The function appends each `Score` with a human-readable reason to `analysis.Scores`.
5. The total is the sum of all six dimension scores, each out of 5, for a maximum of 30.
6. `Summary()` formats the analysis into a readable string.
7. `main()` creates two sample projects and prints their analyses for comparison.

## Common mistakes

- Mistake: Building projects that are too similar (three CLI tools).
  - Why it happens: It is easier to build the same type of project repeatedly.
  - Fix: Diversify across API services, CLI tools, libraries, and open source contributions.

- Mistake: Skipping tests to ship faster.
  - Why it happens: Testing feels like overhead in a side project.
  - Fix: Use `go test` from day one. Tests are the strongest signal of engineering quality.

- Mistake: Neglecting the README.
  - Why it happens: The code feels self-explanatory to the author.
  - Fix: Write the README as if for a new team member joining your project cold.

- Mistake: Leaving projects unfinished with outstanding work-in-progress comments visible.
  - Why it happens: Side projects naturally compete for limited time.
  - Fix: Either finish a minimal viable version or archive the repository with a note.

## Debugging walkthrough

A developer runs the portfolio analyzer on a project called "go-redis-cache" with 50 stars, tests, documentation, CI, and 3000 lines of code:

```
Analysis for go-redis-cache:
  Complexity:        1/5 — small codebase, manageable scope
  Tech Relevance:    5/5 — uses Go with modern architecture patterns
  Testing:           4/5 — tests present, demonstrates quality mindset
  Documentation:     4/5 — documented, shows communication skills
  Real-World Use:    4/5 — some adoption, demonstrates value
  Maintainability:   4/5 — CI configured, shows engineering discipline
  Total score: 22/30
```

The developer expected a higher complexity score. The description does not mention gRPC, microservices, or distributed systems, so the tech relevance scored 4 instead of 5. After updating the project description to mention "distributed caching layer for microservices," the tech relevance score improves to 5, and the total becomes 23/30. This shows how the analyzer can guide both project selection and presentation.

## Production notes

In a real portfolio review tool, you would supplement automated analysis with manual review. Automated scoring is useful for quick triage but cannot evaluate code quality, architectural decisions, or problem originality. A production tool might also:

- Integrate with the GitHub API to fetch real star counts, last commit dates, and language breakdowns.
- Parse `go.mod` to check dependency hygiene.
- Run `go vet` and `staticcheck` against the project source.
- Track projects over time to visualize improvement.

## Performance implications

The portfolio analyzer runs locally and handles small datasets, so performance is not a concern. If extended to analyze hundreds of GitHub repositories via the API, consider:

- Batching API calls to avoid rate limits.
- Caching analysis results per project.
- Running the code review checks concurrently using goroutines.
- Using context timeouts for external API calls.

## Practice task

Add a new dimension `Innovation` to the portfolio analyzer. It should score projects based on whether the description contains keywords like "novel", "new approach", "alternative", or "unique". Add the dimension constant, the scoring logic, and include it in the total. Run the tests to confirm the result.

## Tests / verification

```bash
go test ./curriculum/modules/16-portfolio-interview-readiness/lessons/01-portfolio-strategy
go run ./curriculum/modules/16-portfolio-interview-readiness/lessons/01-portfolio-strategy
```

The tests verify that high-quality projects score at least 20, low-quality projects score at most 15, and all six dimensions appear in every analysis. After completing the practice task, the tests should still pass with the new dimension included.

## Review questions

1. What three portfolio types exist, and which carries the most signal for a production engineering role?
2. Why do hiring managers care about testing in a portfolio project?
3. What is the minimum viable README structure for a portfolio project?
4. How does Go's tooling lower the barrier to a professional-quality portfolio compared to other languages?
5. A project has 500 GitHub stars but no tests. How does the analyzer score it, and what does that tell you about its portfolio readiness?

## NEXT UP

Project selection and scoping — learn how to define the right project, avoid over-engineering, and estimate a realistic timeline before writing a single line of code.
