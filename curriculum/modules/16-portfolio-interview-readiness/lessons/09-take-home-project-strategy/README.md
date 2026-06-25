# Take-home project strategy

## Learning objective

Develop a systematic approach to completing take-home projects that demonstrates production-quality software engineering, including requirements analysis, time management, testing strategy, documentation, and a structured submission checklist.

## Why this matters

Take-home projects are the single most common technical assessment in the Go ecosystem. Companies like Google, Uber, Shopify, and GitLab use them to evaluate how candidates structure real code. Unlike whiteboard or live-coding interviews, take-home projects test your ability to deliver complete, well-organized solutions. A strong submission can outweigh weaker performance in other interview stages. A weak submission — missing tests, no documentation, poor structure — can end your candidacy regardless of algorithmic skill.

## Mental model

Think of a take-home project as a miniature open-source contribution. The reviewer is a maintainer evaluating your PR. They check whether your code is correct, testable, documented, and maintainable. Your goal is to make their review as easy as possible. Every decision — naming, structure, comments, test coverage — should optimize for reviewer clarity. The clock is your constraint: allocate time deliberately across four phases: analysis (20%), implementation (50%), polish (20%), and review (10%).

## Core idea

A take-home project strategy is a repeatable framework that converts a prompt into a submission that maximizes signal per minute spent. The core elements are:

| Element | Purpose |
|---|---|
| Requirements breakdown | Decompose the prompt into explicit and implicit requirements |
| Architecture sketch | Plan packages, types, and data flow before writing code |
| Time budget | Allocate hours per phase with hard cutoffs |
| Testing strategy | Decide what to test at each level (unit, integration, end-to-end) |
| Documentation plan | Identify what needs README, API docs, or inline comments |
| Submission checklist | Verify completeness before sending |

## Under the hood

When you submit a take-home project, the reviewer evaluates along these axes:

1. **Correctness** — Does it solve the stated problem? Are edge cases handled?
2. **Code quality** — Is the code idiomatic Go? Are names clear? Are functions focused?
3. **Testing** — Is there meaningful test coverage? Do tests test behavior or implementation?
4. **Documentation** — Can the reviewer run it without asking questions?
5. **Project structure** — Does it follow Go conventions (cmd/, internal/, pkg/)?
6. **Error handling** — Are errors propagated and wrapped? Are panics avoided?
7. **Dependency management** — Are dependencies justified? Is `go.mod` clean?

Each axis is a score (1-10). The total weighted score determines pass/fail. Your strategy should optimize the weighted sum given your time budget.

## How Go uses it

Go's toolchain and conventions directly support take-home project excellence:

- `go mod init` / `go mod tidy` enforce clean dependency management.
- `go vet` and `go fmt` catch common mistakes and enforce style automatically.
- `go test -cover` measures test coverage as a quantitative target.
- The `testing` package provides table-driven tests, benchmarks, and fuzzing.
- `net/http/httptest` enables HTTP handler testing without external servers.
- Standard project layout (`cmd/`, `internal/`, `pkg/`) signals familiarity with Go conventions.

## Go example

```go
package main

import (
	"fmt"
	"strings"
)

type ScoreCategory struct {
	Name    string
	Weight  float64
	Score   int
	Comment string
}

type ProjectScore struct {
	ProjectName string
	Categories  []ScoreCategory
}

func (ps ProjectScore) TotalScore() float64 {
	var total, weightSum float64
	for _, c := range ps.Categories {
		total += float64(c.Score) * c.Weight
		weightSum += c.Weight
	}
	if weightSum == 0 {
		return 0
	}
	return total / weightSum * 10
}

func (ps ProjectScore) Summary() string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("Project: %s\n", ps.ProjectName))
	b.WriteString(fmt.Sprintf("Overall Score: %.1f/100\n", ps.TotalScore()))
	for _, c := range ps.Categories {
		bar := strings.Repeat("\u2588", c.Score) + strings.Repeat("\u2591", 10-c.Score)
		b.WriteString(fmt.Sprintf("  %s: %s %d/10", c.Name, bar, c.Score))
		if c.Comment != "" {
			b.WriteString(fmt.Sprintf(" (%s)", c.Comment))
		}
		b.WriteString("\n")
	}
	return b.String()
}

func main() {
	score := ProjectScore{
		ProjectName: "URL Shortener API",
		Categories: []ScoreCategory{
			{Name: "Testing", Weight: 1.0, Score: 8, Comment: "unit + integration tests present"},
			{Name: "Documentation", Weight: 0.8, Score: 7, Comment: "README clear, missing API docs"},
			{Name: "Code Clarity", Weight: 1.0, Score: 9, Comment: "clean structure, good naming"},
			{Name: "Error Handling", Weight: 0.9, Score: 7, Comment: "basic errors, no sentinel errors"},
			{Name: "Project Structure", Weight: 0.7, Score: 8, Comment: "standard Go layout"},
		},
	}
	fmt.Print(score.Summary())
}
```

Run with `go run .` to see the scored evaluation with visual bars.

## Step-by-step execution

When the `Summary()` method runs:

1. `TotalScore()` iterates over categories, multiplying each `Score` by `Weight` and summing into `total`.
2. `weightSum` accumulates the weights to normalize the result.
3. Dividing `total` by `weightSum` gives a weighted average, then multiplying by 10 scales it to 0-100.
4. `Summary()` builds a string: project name header, overall score line, then per-category lines.
5. Each category shows a Unicode bar (`█` repeated for score, `░` for remainder) for visual scanning.
6. Comments appear in parentheses after the bar.

For the example project, the calculation is:
- Total weighted: (8×1.0)+(7×0.8)+(9×1.0)+(7×0.9)+(8×0.7) = 8+5.6+9+6.3+5.6 = 34.5
- Weight sum: 1.0+0.8+1.0+0.9+0.7 = 4.4
- Final: 34.5 / 4.4 × 10 = 78.4

## Common mistakes

- **Starting implementation immediately.** Without requirements analysis, you build the wrong thing. Spend the first 20% of your budget reading and planning.
- **Over-investing in edge cases.** A perfect solution for an obscure edge case at the expense of core functionality signals poor prioritization. Cover the happy path first, then add edge cases if time permits.
- **Skipping tests.** Many reviewers weight testing at 30-40% of the score. A submission with no tests is almost an automatic rejection. Write at least one table-driven test per public function.
- **Neglecting the README.** If the reviewer cannot build and run your project in under two minutes, they will not evaluate your code. Include clear setup, run, and test instructions.
- **Over-engineering.** Using channels, goroutines, or dependency injection when a simple function suffices adds complexity without benefit. Write the simplest correct solution first.
- **Ignoring error messages.** Compiler warnings, `go vet` failures, or `go fmt` issues before submission suggest carelessness. Run all three before packaging.

## Debugging walkthrough

A candidate submits a take-home project for a URL shortener API. The reviewer runs it and finds:

```go
// BUG: handler always returns 500
func handleShorten(w http.ResponseWriter, r *http.Request) {
    url := r.URL.Query().Get("url")
    if url == "" {
        w.WriteHeader(400)
        return
    }
    id, err := store.Save(url)
    if err != nil {
        w.WriteHeader(500) // line 42
    }
    fmt.Fprint(w, id)
}
```

**Symptom**: Every valid request returns 500 with no body.

**Investigation**: Add a log line before line 42:
```go
log.Printf("save error: %v", err)
```

**Root cause**: `store.Save` checks if the URL already exists and returns `ErrDuplicate`, but the handler doesn't check the error type. The candidate assumed all errors from `Save` were fatal, but `ErrDuplicate` is a recoverable client error.

**Fix**:
```go
if errors.Is(err, store.ErrDuplicate) {
    w.WriteHeader(http.StatusConflict)
    fmt.Fprint(w, "URL already exists")
    return
}
if err != nil {
    w.WriteHeader(http.StatusInternalServerError)
    fmt.Fprint(w, "internal error")
    return
}
```

**Lesson learned**: Always inspect error types. Use `errors.Is` and `errors.As` to distinguish error classes. Return appropriate HTTP status codes for each class.

## Production notes

- **Time boxing is non-negotiable.** Use a timer. When the allocated block expires, move to the next phase. Perfecting one section while leaving others empty produces a worse overall impression.
- **Version control matters.** Commit early and often. A clean git history with descriptive messages signals professional habits. Push to a private GitHub repo and share the link.
- **Include a design decisions doc.** A brief `DECISIONS.md` explaining why you chose certain approaches (why `sqlx` over raw `database/sql`, why a `sync.Mutex` over a channel) demonstrates senior-level thinking.
- **Never use `context.Background()` in request-scoped work.** Always derive from the request context so cancellation propagates correctly. For long-running background tasks, use `context.WithTimeout` or `context.WithCancel`.

## Performance implications

- **Test execution time** matters. If your tests take 30 seconds to run, reviewers may skip them. Keep unit tests under 1ms per test. Use `-short` flag for quick validation and integration test build tags for slow tests.
- **Startup time** is a quality signal. A Go binary that starts in 2ms (no framework, no init complexity) is preferred over one that takes 500ms. Avoid unnecessary `init()` functions and heavy imports.
- **Binary size** reflects dependency discipline. Run `go build -o /dev/null && ls -lh` to check. A 5MB binary is fine; a 50MB binary suggests unnecessary dependencies or embedded assets.

## Practice task

Build a `ProjectChecklist` tool that validates a take-home project directory against a standard submission checklist. Implement:

1. `CheckItem` struct with fields: `Name string`, `Passed bool`, `Notes string`
2. `Checklist` struct with `Items []CheckItem` and methods:
   - `PassRate() float64` — returns percentage of passed items.
   - `Report() string` — returns formatted pass/fail per item.
3. A `ValidateProject(path string) Checklist` function that checks for:
   - `go.mod` exists
   - `main.go` or `cmd/` exists
   - `*_test.go` files exist
   - `README.md` exists
   - `go vet` passes (simulate with file existence check)
4. `main()` that validates a pretend project and prints the report.

## Tests / verification

```bash
go run ./curriculum/modules/16-portfolio-interview-readiness/lessons/09-take-home-project-strategy
go test ./curriculum/modules/16-portfolio-interview-readiness/lessons/09-take-home-project-strategy
```

The existing tests verify `TotalScore` (perfect, average, empty, single, uneven weights) and `Summary` output. After completing the practice task, add table-driven tests for `ValidateProject` covering existing project, missing README, missing tests, and empty directory.

## Review questions

1. What four phases should you allocate time across in a take-home project, and what percentage of total time should each receive?
2. Why is testing weighted so heavily by reviewers, and what is the minimum acceptable test coverage for a submission?
3. Name three Go toolchain commands you must run before submitting a take-home project. What does each check?
4. A handler returns 500 for every valid request. Describe your step-by-step debugging approach.
5. What is the purpose of a `DECISIONS.md` file in a take-home submission, and who benefits from it?

## NEXT UP

Negotiation fundamentals — learn how to evaluate offers, understand total compensation, and negotiate effectively when you receive the offer.
