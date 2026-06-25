# Open source contribution

## Learning objective

Contribute effectively to open source Go projects by finding suitable projects, making well-scoped first contributions, following PR etiquette, participating in issue triage, and building a professional reputation through sustained quality contributions.

## Why this matters

Open source contributions are the most authentic signal of Go engineering ability. Hiring managers at top Go shops (Google, Tailscale, HashiCorp, Docker, Grafana) actively review candidates' GitHub profiles. A history of well-structured PRs, thoughtful code reviews, and issue triage demonstrates exactly the skills these companies need. Open source work also expands your network, exposes you to production-scale codebases, and builds a public portfolio that never needs a resume bullet point to explain.

## Mental model

Think of open source contribution as a tiered career ladder: Consumer (download and use packages) → Reporter (file issues with reproduction steps) → Contributor (submit PRs with tests) → Reviewer (provide feedback on others' PRs) → Maintainer (merge PRs and set project direction). Each tier requires more context and responsibility. Start at the tier above your current comfort level and progress upward. The code is the project's asset, but community and communication are its operating system.

## Core idea

Open source contribution is a repeatable process of finding, understanding, fixing, and submitting improvements to public codebases. The key areas:

| Area | Description | Entry point |
|---|---|---|
| Finding projects | Locate active Go projects with good first issues | `good first issue` label, `help wanted` |
| First contribution | Make a small, low-risk change with high value | Documentation, test coverage, bug fixes |
| PR etiquette | Communicate clearly and respect maintainer time | Descriptive title, small diff, tests included |
| Issue triage | Reproduce bugs, add labels, suggest priorities | Read CONTRIBUTING.md, comment on issues |
| Building reputation | Earn trust through consistent quality over time | Review others' code, write thorough PR descriptions |

## Under the hood

When you submit a PR to an active Go project, the maintainer evaluates:

1. **Does the PR address a real need?** — Is there an issue linking to the change? Does the PR description explain the motivation?
2. **Is the code idiomatic Go?** — `gofmt` compliance, proper error handling, no unused variables, correct use of interfaces.
3. **Are tests included?** — Go projects often require 80%+ test coverage for new code. Table-driven tests are the standard.
4. **Does the PR compile?** — CI runs `go build ./...`, `go vet ./...`, `go test ./...` automatically.
5. **Is the commit history clean?** — Squash commits, write descriptive messages, sign commits if required.

Projects use CI/CD pipelines (GitHub Actions, CircleCI, Drone) that run the full Go toolchain on every PR. A red CI build is the fastest way to get a PR closed without review.

## How Go uses it

Go itself is an open source project that embodies ideal contribution practices. The standard library's contribution guidelines (`https://go.dev/doc/contribute`) require:

- Signing the Contributor License Agreement (CLA).
- Running `git codereview` for Gerrit-based reviews.
- Writing tests in `*_test.go` files using table-driven patterns.
- Following the Go Code Review Comments guide.
- Participating in the Gerrit review process with iterative feedback.

The Go community has a strong culture of respectful, thorough code review. This culture extends to the broader Go ecosystem: `github.com/golang/go` sets the tone for how Go projects handle contributions.

## Go example

```go
package main

import (
	"fmt"
	"time"
)

type PullRequest struct {
	ID        int
	Title     string
	Repository string
	Author    string
	CreatedAt time.Time
	MergedAt  *time.Time
	Additions int
	Deletions int
	Comments  int
	Reviewed  bool
}

type Contributor struct {
	Username     string
	PullRequests []PullRequest
}

func (c Contributor) TotalPRs() int           { return len(c.PullRequests) }
func (c Contributor) MergedPRs() int           { count := 0; for _, pr := range c.PullRequests { if pr.MergedAt != nil { count++ } }; return count }
func (c Contributor) MergeRate() float64       { if len(c.PullRequests) == 0 { return 0 }; return float64(c.MergedPRs()) / float64(len(c.PullRequests)) * 100 }
func (c Contributor) TotalChanges() (additions, deletions int) { for _, pr := range c.PullRequests { additions += pr.Additions; deletions += pr.Deletions }; return }
func (c Contributor) AverageTimeToMerge() time.Duration { var total time.Duration; var count int; for _, pr := range c.PullRequests { if pr.MergedAt != nil { total += pr.MergedAt.Sub(pr.CreatedAt); count++ } }; if count == 0 { return 0 }; return total / time.Duration(count) }

type ContributionReport struct {
	Contributor    Contributor
	TotalPRs       int
	MergedPRs      int
	MergeRate      float64
	Additions      int
	Deletions      int
	AvgMergeTime   time.Duration
	ReviewedCount  int
}

func GenerateReport(c Contributor) ContributionReport {
	additions, deletions := c.TotalChanges()
	reviewed := 0
	for _, pr := range c.PullRequests {
		if pr.Reviewed { reviewed++ }
	}
	return ContributionReport{
		Contributor:   c,
		TotalPRs:      c.TotalPRs(),
		MergedPRs:     c.MergedPRs(),
		MergeRate:     c.MergeRate(),
		Additions:     additions,
		Deletions:     deletions,
		AvgMergeTime:  c.AverageTimeToMerge(),
		ReviewedCount: reviewed,
	}
}

func main() {
	now := time.Now()
	contributor := Contributor{
		Username: "gopher-engineer",
		PullRequests: []PullRequest{
			{ID: 101, Title: "Add request validation middleware", Repository: "go-api-starter", Author: "gopher-engineer",
				CreatedAt: now.Add(-72 * time.Hour), MergedAt: timePtr(now.Add(-48 * time.Hour)), Additions: 120, Deletions: 30, Comments: 4, Reviewed: true},
			{ID: 102, Title: "Fix race condition in worker pool", Repository: "go-api-starter", Author: "gopher-engineer",
				CreatedAt: now.Add(-120 * time.Hour), MergedAt: timePtr(now.Add(-96 * time.Hour)), Additions: 45, Deletions: 15, Comments: 6, Reviewed: true},
			{ID: 103, Title: "Update README with deployment guide", Repository: "docs", Author: "gopher-engineer",
				CreatedAt: now.Add(-24 * time.Hour), MergedAt: timePtr(now.Add(-12 * time.Hour)), Additions: 200, Deletions: 10, Comments: 2, Reviewed: true},
			{ID: 104, Title: "Implement rate limiting middleware", Repository: "go-api-starter", Author: "gopher-engineer",
				CreatedAt: now.Add(-12 * time.Hour), MergedAt: nil, Additions: 80, Deletions: 5, Comments: 3, Reviewed: false},
		},
	}
	report := GenerateReport(contributor)
	fmt.Printf("Contribution Report for %s\n", report.Contributor.Username)
	fmt.Printf("  Total PRs:        %d\n", report.TotalPRs)
	fmt.Printf("  Merged PRs:       %d\n", report.MergedPRs)
	fmt.Printf("  Merge Rate:       %.1f%%\n", report.MergeRate)
	fmt.Printf("  Additions:        %d\n", report.Additions)
	fmt.Printf("  Deletions:        %d\n", report.Deletions)
	fmt.Printf("  Avg Time to Merge: %s\n", report.AvgMergeTime.Round(time.Hour))
	fmt.Printf("  Reviewed PRs:     %d\n", report.ReviewedCount)
}

func timePtr(t time.Time) *time.Time { return &t }
```

Run with `go run .` to see a contribution report with PR statistics, merge rate, and timeline.

## Step-by-step execution

When `GenerateReport` runs:

1. The contributor's PR list is scanned to count total versus merged PRs. Merged PRs are identified by a non-nil `MergedAt` pointer.
2. `MergeRate` divides merged by total and converts to percentage. Four PRs with three merged yields 75%.
3. `TotalChanges` sums additions and deletions across all PRs, revealing the contributor's net impact on the codebase.
4. `AverageTimeToMerge` computes mean time from creation to merge for all merged PRs, filtering out open PRs.
5. The report also counts how many PRs received a reviewed flag (indicating the contributor also reviewed other code).
6. Results print as a formatted table showing key contribution metrics.

For the example: 4 total PRs, 3 merged (75%), 445 additions + 60 deletions, average ~2 days to merge, 3 PRs reviewed.

## Common mistakes

- **Submitting a large PR as a first contribution.** A 2000-line change is hard to review and likely to be rejected. Start with documentation, test coverage, or a small bug fix (under 100 lines).
- **Not searching for existing issues.** Duplicate PRs waste everyone's time. Before starting, search for existing issues and PRs. Comment on the issue to signal your intent.
- **Skipping the CONTRIBUTING.md.** Every well-run project has contribution guidelines. Read them. They describe commit message format, coding style, testing requirements, and CLA signing.
- **Ignoring CI failures.** A red build is the maintainer's first signal. Fix CI before requesting review. Run `go build ./...`, `go vet ./...`, and `go test ./...` locally.
- **Getting defensive about feedback.** Code review feedback is about the code, not the person. Respond graciously, make the requested changes, and explain your reasoning if you disagree.
- **Abandoning PRs after feedback.** A PR that gets review comments and then sits for two weeks may be closed. Respond within 48 hours or communicate your timeline.

## Debugging walkthrough

A contributor opens a PR with a failing test on CI. The local tests pass but the CI build fails:

```
--- FAIL: TestParseConfig (0.01s)
    config_test.go:25: ParseConfig("testdata/config.yaml") = nil, want non-nil
```

**Symptom**: Local machine passes, CI fails. The config file exists locally but not in CI.

**Investigation**: The contributor checks the CI workflow file `.github/workflows/ci.yml`:

```yaml
- name: Run tests
  run: go test ./...
```

The test reads `testdata/config.yaml`, but the CI checkout step doesn't include `testdata/` because the `paths` filter was too restrictive.

**Root cause**: The contributor added the `testdata/` directory but the CI pipeline only triggers on changes to `*.go` files. The testdata directory was never checked out.

**Fix**: Update the CI trigger or add a checkout step that includes all files. Also fix the test to use `os.Open` with proper error checking:

```go
func TestParseConfig(t *testing.T) {
    f, err := os.Open("testdata/config.yaml")
    if err != nil {
        t.Skip("testdata not available:", err)
    }
    defer f.Close()
    // ...
}
```

**Lesson learned**: Always check CI configuration before submitting. Add `t.Skip` for environment-dependent tests so they fail gracefully rather than hard-fail.

## Production notes

- **Start with documentation.** Fixing a typo, improving a README, or adding an example teaches you the contribution workflow with minimal risk. It also earns goodwill from maintainers.
- **Look for projects with active maintainers.** A project with no commits in 6 months will likely ignore your PR. Check the merge latency before investing effort.
- **Use the project's communication channels.** Join the Slack, Discord, or mailing list. Introduce yourself before submitting a PR. Context and relationships accelerate review.
- **Keep a contribution log.** Track your PRs, their status, and what you learned. This becomes material for performance reviews, resume bullets, and interview stories.
- **Review others' PRs first.** Code review is a skill that transfers directly to your day job. Leave constructive comments on 3-5 PRs before submitting your first code change.

## Performance implications

- **CI pipeline time** affects maintainer willingness to merge. Each CI run costs compute and takes 2-10 minutes. Minimize CI iterations by testing thoroughly before pushing.
- **Test suite speed** matters in large projects. A package adding 5 seconds to the test suite may be asked to optimize. Keep unit tests focused and use build tags for integration tests.
- **Binary size impact** from dependencies is scrutinized. Adding a new dependency to a popular Go project requires justification. Prefer standard library solutions when possible.

## Practice task

Build a `ProjectFinder` tool that helps identify good open source projects for contribution. Implement:

1. `Project` struct with fields: `Name string`, `Stars int`, `Issues []Issue`, `LastCommit time.Time`
2. `Issue` struct with fields: `Title string`, `Labels []string`, `URL string`
3. `FilterByLabel(issues []Issue, label string) []Issue` — returns issues matching a label.
4. `FindGoodFirstIssues(project Project) []Issue` — returns issues tagged "good first issue" or "help wanted".
5. `ScoreProject(p Project) float64` — scores projects by activity: stars (0-10), recent commit (0-10), open issue count (0-10), weighted equally.
6. `main()` that creates sample projects with issues, finds good first issues, and prints ranked projects.

## Tests / verification

```bash
go run ./curriculum/modules/16-portfolio-interview-readiness/lessons/11-open-source-contribution
go test ./curriculum/modules/16-portfolio-interview-readiness/lessons/11-open-source-contribution
```

The existing tests verify `MergeRate` (all, half, none, empty), `TotalChanges`, and `GenerateReport` (counts and rates). After completing the practice task, add tests for `FilterByLabel`, `FindGoodFirstIssues`, and `ScoreProject`.

## Review questions

1. What are the five tiers of the open source contribution ladder, from lowest to highest?
2. Name three things a maintainer evaluates when reviewing your first PR to their project.
3. Why should you search for existing issues before starting work on a contribution?
4. What is the recommended size and scope for a first contribution to an unfamiliar project?
5. What should you do if your local tests pass but CI fails on your PR?

## NEXT UP

Technical blog writing — learn how to write clear, engaging technical blog posts about Go that build your professional brand and share knowledge effectively.
