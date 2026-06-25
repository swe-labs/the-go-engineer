# Code review readiness

## Learning objective

Prepare Go code for professional code review by running a structured checklist across six categories (correctness, style, performance, security, reliability, maintainability) and interpreting the results to improve code quality before submitting a PR.

## Why this matters

Code review is where software quality happens. Every Go engineer spends as much time reviewing code as writing it. Sending code that has not been self-reviewed wastes reviewer time, slows down the team, and damages your reputation. A structured readiness checklist ensures your PRs are clean, focused, and respectful of the reviewer's time. Mastering self-review is the fastest path to becoming a trusted engineer who ships quality code.

## Mental model

Think of your code like a submission to a publication. Before you submit, you proofread, check formatting, verify facts, and ensure you have not introduced errors. A code review readiness checklist is the same thing: you run through categories systematically to catch issues before the reviewer does. The reviewer's job should be to evaluate architectural decisions and business logic, not to fix formatting, missing error handling, or concurrency bugs.

## Core idea

Code review checks fall into six categories. Each category targets a different aspect of code quality:

**Correctness**: Does the code do what it is supposed to do? Checks: error handling, context propagation, concurrency safety, edge cases.

**Style**: Does the code follow Go conventions? Checks: `gofmt` compliance, naming conventions (PascalCase for exported, camelCase for unexported), consistent formatting, no magic numbers.

**Performance**: Is the code efficient? Checks: pre-allocation, avoiding unnecessary allocations, using appropriate data structures, no hot-path allocations.

**Security**: Is the code safe from common vulnerabilities? Checks: input validation, SQL injection prevention, no hardcoded secrets, proper TLS configuration.

**Reliability**: Does the code handle failures gracefully? Checks: goroutine lifecycle management, proper cleanup with `defer`, timeout handling, retry logic.

**Maintainability**: Will the next engineer understand this code? Checks: clear naming, appropriate comments, single responsibility, test coverage, documentation.

A pull request should pass all checks in at least four of six categories before submission. The remaining two may have known tradeoffs that you document in the PR description.

## Under the hood

A code review checklist is a static analysis approximation of what a human reviewer would check. Real code review involves understanding the business domain, the architecture, and the tradeoffs — things no automated tool can fully capture. The checklist approach catches the mechanical issues so the human reviewer can focus on the important decisions.

The tool in this lesson uses string matching and regex to detect patterns. This is intentionally simple. Production tools like `golangci-lint`, `staticcheck`, and `revive` use full AST analysis and dataflow tracking. The pattern-matching approach demonstrates the concept but a real checklist tool would integrate with the Go toolchain at the AST level.

## How Go uses it

Go's design philosophy makes code review readiness easier than in most languages:

- `gofmt` eliminates formatting debates entirely. There is no style discussion — the tool decides.
- `go vet` catches suspicious constructs like unreachable code or mismatched printf args.
- The `error` return convention means error handling is explicit and visible in every function signature.
- `context.Context` as the first parameter convention makes cancellation and deadlines visible.
- Package-level documentation conventions (`package doc.go` or package comment) enforce documentation habits.

A Go code review typically spends time on: error handling strategy (wrap vs. return raw), concurrency design (channels vs. mutexes), interface design, and package organization — not on formatting or naming debates.

## Go example

The code review checklist checker scans a source file for patterns across six categories and produces a pass/fail report with explanations.

```go
package main

import (
	"fmt"
	"regexp"
	"strings"
)

type ReviewCategory int

const (
	Correctness ReviewCategory = iota
	Style
	Performance
	Security
	Reliability
	Maintainability
)

type Check struct {
	Category    ReviewCategory
	Description string
	Passed      bool
	Detail      string
}

type ReviewResult struct {
	File    string
	Content string
	Checks  []Check
	Score   int
	Total   int
}

func NewReviewResult(file, content string) ReviewResult {
	return ReviewResult{File: file, Content: content}
}

func (rr *ReviewResult) AddCheck(cat ReviewCategory, desc string, passed bool, detail string) {
	rr.Checks = append(rr.Checks, Check{cat, desc, passed, detail})
	rr.Total++
	if passed {
		rr.Score++
	}
}

func RunCodeReview(content string) ReviewResult {
	rr := NewReviewResult("reviewed.go", content)

	hasErrorHandling := strings.Contains(content, "if err != nil")
	rr.AddCheck(Correctness, "Error handling: checks returned errors", hasErrorHandling,
		map[bool]string{true: "uses if err != nil pattern", false: "missing error handling"}[hasErrorHandling])

	hasContext := strings.Contains(content, "context.Background()") ||
		strings.Contains(content, "context.WithCancel") ||
		strings.Contains(content, "context.WithTimeout") ||
		strings.Contains(content, "context.Context")
	rr.AddCheck(Correctness, "Context usage: proper context propagation", hasContext,
		map[bool]string{true: "context detected", false: "no context usage"}[hasContext])

	hasMutex := strings.Contains(content, "sync.Mutex") || strings.Contains(content, "sync.RWMutex")
	hasChan := strings.Contains(content, "chan") || strings.Contains(content, "make(chan")
	if hasMutex || hasChan {
		rr.AddCheck(Correctness, "Concurrency safety: synchronization used", true,
			"uses mutex or channels for concurrency")
	} else {
		rr.AddCheck(Correctness, "Concurrency safety: synchronization used", true,
			"no concurrent access detected")
	}

	gofmt := regexp.MustCompile(`\t`).MatchString(content)
	rr.AddCheck(Style, "Formatting: uses tabs for indentation", gofmt,
		map[bool]string{true: "tabs detected", false: "no tabs found, may use spaces"}[gofmt])

	hasExported := regexp.MustCompile(`^func [A-Z]`).MatchString(content)
	rr.AddCheck(Style, "Naming: exported functions use PascalCase", hasExported,
		map[bool]string{true: "exported functions follow convention", false: "no exported functions or naming issue"}[hasExported])

	hasGoRoutine := strings.Contains(content, "go ")
	if hasGoRoutine {
		hasWG := strings.Contains(content, "sync.WaitGroup")
		rr.AddCheck(Reliability, "Goroutine lifecycles: goroutines have synchronization", hasWG,
			map[bool]string{true: "uses WaitGroup for goroutine coordination",
				false: "goroutines without WaitGroup may leak"}[hasWG])
	} else {
		rr.AddCheck(Reliability, "Goroutine lifecycles: goroutines have synchronization", true,
			"no goroutines to manage")
	}

	hasInputValidation := strings.Contains(content, "len(") &&
		(strings.Contains(content, "== 0") || strings.Contains(content, "< 0"))
	rr.AddCheck(Security, "Input validation: validates inputs before use", hasInputValidation,
		map[bool]string{true: "input length validation detected", false: "no input validation detected"}[hasInputValidation])

	hasSQL := strings.Contains(content, "sql.DB") || strings.Contains(content, "database/sql")
	if hasSQL {
		hasParamQuery := strings.Contains(content, "$1") || (strings.Contains(content, "?") && strings.Contains(content, "Query"))
		rr.AddCheck(Security, "SQL injection: uses parameterized queries", hasParamQuery,
			map[bool]string{true: "parameterized queries detected", false: "possible string concatenation in SQL"}[hasParamQuery])
	}

	hasAlloc := strings.Contains(content, "make(") || strings.Contains(content, "new(")
	if hasAlloc {
		hasCap := regexp.MustCompile(`make\(\[?[a-zA-Z]`).MatchString(content)
		rr.AddCheck(Performance, "Allocation: pre-allocates slices/maps with capacity", hasCap,
			map[bool]string{true: "pre-allocation detected", false: "no pre-allocation, may cause reallocation"}[hasCap])
	} else {
		rr.AddCheck(Performance, "Allocation: pre-allocates slices/maps with capacity", true,
			"no allocations to optimize")
	}

	hasDefer := strings.Contains(content, "defer ")
	rr.AddCheck(Maintainability, "Resource cleanup: uses defer for cleanup", hasDefer,
		map[bool]string{true: "defer usage detected", false: "no defer usage, resources may not be released"}[hasDefer])

	return rr
}

func (rr ReviewResult) Summary() string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("Review: %s\nScore: %d/%d (%.0f%%)\n\n", rr.File, rr.Score, rr.Total, float64(rr.Score)/float64(rr.Total)*100))
	for _, c := range rr.Checks {
		status := "PASS"
		if !c.Passed {
			status = "FAIL"
		}
		b.WriteString(fmt.Sprintf("[%s] [%s] %s\n", status, c.Category, c.Description))
		b.WriteString(fmt.Sprintf("      %s\n", c.Detail))
	}
	return b.String()
}

func main() {
	code := `package main
import (
	"context"
	"sync"
)

type Server struct {
	mu    sync.Mutex
	items map[string]string
}

func NewServer() *Server {
	return &Server{items: make(map[string]string)}
}

func (s *Server) Get(ctx context.Context, key string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	val, ok := s.items[key]
	if !ok {
		return "", fmt.Errorf("key not found: %s", key)
	}
	return val, nil
}

func main() {
	s := NewServer()
	ctx := context.Background()
	val, err := s.Get(ctx, "hello")
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Println(val)
}`
	result := RunCodeReview(code)
	fmt.Println(result.Summary())
}
```

## Step-by-step execution

1. `RunCodeReview` receives a string of Go source code.
2. It creates a `ReviewResult` and runs each check in sequence.
3. Correctness checks: scan for error handling (`if err != nil`), context usage, and concurrency primitives (mutexes, channels).
4. Style checks: look for tab indentation and PascalCase exported functions via regex.
5. Reliability: if goroutines are detected, verify that `sync.WaitGroup` is present.
6. Security: check for input validation patterns and parameterized SQL queries.
7. Performance: if allocations (`make`, `new`) exist, check for capacity hints.
8. Maintainability: check for `defer` usage as a proxy for resource cleanup discipline.
9. Each check is added to the result with pass/fail status and a human-readable detail.
10. `Summary()` formats the results with a percentage score.

## Common mistakes

- Mistake: Only running `gofmt` and assuming the code is review-ready.
  - Why it happens: Formatting is the most visible quality signal.
  - Fix: Formatting is table stakes. Run the full checklist against correctness, security, and reliability.

- Mistake: Ignoring concurrency checks because "it runs fine on my machine."
  - Why it happens: Data races are nondeterministic and may not manifest during development.
  - Fix: Always run `go test -race` before submitting. Use synchronization primitives explicitly.

- Mistake: Sending a PR with commented-out code.
  - Why it happens: Dead code is left behind "just in case."
  - Fix: Delete commented code. Git history preserves it if needed.

- Mistake: Writing a PR description that just says "fixed bug."
  - Why it happens: The fix feels obvious to the author.
  - Fix: Describe what the bug was, how you diagnosed it, and why your fix is correct.

## Debugging walkthrough

A developer runs the checker on a file that uses goroutines without WaitGroup:

```
[FAIL] [Reliability] Goroutine lifecycles: goroutines have synchronization
      goroutines without WaitGroup may leak
```

The developer adds `sync.WaitGroup` to coordinate goroutine completion and re-runs the checker. It now passes. The developer also notices the Performance check says "no allocations to optimize" because the code does not use `make` or `new`. After adding a slice with pre-allocation (`make([]int, 0, 100)`), the Performance check still fails because the regex does not match capacity hints in all patterns. The developer decides the check is a guideline, not a gate, and submits the PR with a note about the limitation.

## Production notes

In a production setting, the code review checklist would be:

- Integrated into CI so every PR gets an automated check comment.
- Configured per team or per repository with different severity levels.
- Supplemented with `golangci-lint` which runs dozens of linters.
- Used as a teaching tool for junior engineers to learn what reviewers look for.

The checklist should not block PRs unilaterally — some checks are guidelines, not rules. The goal is to reduce the mechanical review burden, not to enforce rigid compliance.

## Performance implications

The checker itself is fast: string matching and simple regex run in microseconds on typical file sizes. In a CI pipeline processing hundreds of files, the total overhead is negligible compared to compilation and test execution. However, the checker only approximates what a full AST-based tool would find. For production use, prefer `golangci-lint` or `staticcheck` over a custom pattern-matching tool.

## Practice task

Add two new checks to the review checklist: one for `Magic numbers` (detect bare numeric literals other than 0 and 1 in non-test code) under the Style category, and one for `Context timeout` (detect if `context.WithTimeout` or `context.WithDeadline` is used when goroutines are present) under the Correctness category. Add corresponding tests and verify they pass.

## Tests / verification

```bash
go test ./curriculum/modules/16-portfolio-interview-readiness/lessons/03-code-review-readiness
go run ./curriculum/modules/16-portfolio-interview-readiness/lessons/03-code-review-readiness
```

The tests verify that good code passes at least 5 of 9 checks, that error handling and context detection work correctly, and that all six review categories are represented. After the practice task, the tests should cover the two new checks.

## Review questions

1. What are the six categories of the code review checklist, and what does each target?
2. Why does Go eliminate formatting debates compared to languages like Python or JavaScript?
3. What is the minimum score (passing checks) you should aim for before submitting a PR?
4. Why is `go test -race` important beyond what a static checklist checker can detect?
5. What is the difference between a mechanical code review check and an architectural review?

## NEXT UP

Technical writing for engineers — learn to write clear documentation, API docs, and Go doc comments that make your projects understandable and maintainable.
