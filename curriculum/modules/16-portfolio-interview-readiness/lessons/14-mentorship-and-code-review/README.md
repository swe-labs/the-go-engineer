# Mentorship and code review

## Learning objective

Provide constructive, actionable code review feedback that helps other engineers grow, and practice mentorship techniques that accelerate the development of junior and peer engineers in a professional Go environment.

## Why this matters

Code review is the highest-leverage activity in software engineering. A single thorough review can prevent a production outage, teach a junior engineer a pattern they will use for years, and establish coding standards across an entire team. At senior+ levels, your impact is measured not by the code you write but by the quality of code your team produces. Mentorship through code review is how staff and principal engineers scale their knowledge across organizations. Companies like Google, Stripe, and HashiCorp explicitly evaluate code review quality in promotion criteria.

## Mental model

Code review as mentorship is a teaching conversation, not an audit. The author has invested time and ego in their code. Your job is to improve the code while preserving the author's motivation and autonomy. Think of each comment as a teaching moment: explain not just what to change, but why the change matters and what principle it demonstrates. A review comment that teaches a concept ("this pattern separates concerns") is more valuable than one that just states a fix ("move this to a separate function"). The author should learn something from every review.

## Core idea

Effective code review feedback has four dimensions beyond correctness:

| Dimension | Description | Example |
|---|---|---|
| Specificity | Point to exact lines with concrete suggestions | Not "this is wrong" but "line 42: consider using `sync.Mutex` to protect `counter`" |
| Tone | Respectful and collaborative | "Have you considered..." vs "You should..." |
| Actionability | Author knows exactly what to do | "Suggest extracting lines 15-30 into a `parseConfig` function" |
| Teaching | Explains the underlying principle | "Go uses `error` as an interface so callers can use `errors.Is` to check specific error types" |

The ideal review ratio is 3:1:1 — three positive observations for every critical comment and every suggestion/question. This ratio maintains psychological safety while delivering substantive improvement.

## Under the hood

When you submit a code review, the author processes your feedback through a psychological filter:

1. **Safety check**: Is the reviewer attacking me or my code? Comments that start with "You" trigger defensiveness. Comments that start with "The code" or "This function" keep the focus on the artifact.
2. **Effort assessment**: Is this comment worth the fix? If the suggested change requires 20 minutes but provides marginal benefit, the author may ignore it. Be explicit about importance: "This is critical for correctness" vs "Consider for readability."
3. **Learning opportunity**: Does this comment teach something new? Comments that explain the rationale behind a suggestion turn a fix into a learning moment. "Using `io.ReadAll` with no size limit can OOM the server. Prefer `io.LimitReader` with a configurable max."

## How Go uses it

The Go project's code review culture is captured in the "Code Review Comments" guide (go.dev/wiki/CodeReviewComments). Key conventions:

- **Use Go's error handling idioms.** If a function returns an error, check it. If you see `if err != nil { return err }` without wrapping, suggest `fmt.Errorf("context: %w", err)`.
- **Prefer table-driven tests.** If you see repeated test functions with slightly different inputs, suggest converting to a table-driven test with subtests.
- **Avoid naked returns** in functions longer than 10 lines. They reduce readability.
- **Use `gofmt`.** If the code isn't formatted, the first comment should be "Run `gofmt -s` on this file."
- **Avoid package-level state.** If you see `var db *sql.DB` at the package level, suggest dependency injection.

Go's toolchain supports code review automation: `golangci-lint` catches 80% of common issues automatically, reducing the manual review burden so reviewers can focus on architecture and design.

## Go example

```go
package main

import (
	"fmt"
	"strings"
)

type ReviewComment struct {
	Text     string
	Category string // "positive", "critical", "suggestion", "question"
}

type Review struct {
	Reviewer string
	Comments []ReviewComment
}

type ReviewAnalysis struct {
	Review          Review
	TotalComments   int
	PositiveCount   int
	CriticalCount   int
	SuggestionCount int
	QuestionCount   int
	Sentiment       string
	Actionable      int
}

func AnalyzeReview(r Review) ReviewAnalysis {
	a := ReviewAnalysis{Review: r}
	for _, c := range r.Comments {
		a.TotalComments++
		switch c.Category {
		case "positive":   a.PositiveCount++
		case "critical":   a.CriticalCount++
		case "suggestion": a.SuggestionCount++
		case "question":   a.QuestionCount++
		}
	}
	a.Sentiment = computeSentiment(a)
	a.Actionable = computeActionability(r.Comments)
	return a
}

func computeSentiment(a ReviewAnalysis) string {
	if a.TotalComments == 0 { return "neutral" }
	pos := float64(a.PositiveCount) / float64(a.TotalComments)
	crit := float64(a.CriticalCount) / float64(a.TotalComments)
	switch {
	case pos > 0.5:  return "positive"
	case crit > 0.4: return "negative"
	default:         return "balanced"
	}
}

func computeActionability(comments []ReviewComment) int {
	if len(comments) == 0 { return 0 }
	actionable := 0
	for _, c := range comments {
		words := strings.Fields(c.Text)
		for _, w := range words {
			lower := strings.ToLower(w)
			if lower == "consider" || lower == "suggest" || lower == "recommend" ||
				lower == "try" || lower == "change" || lower == "rename" ||
				lower == "move" || lower == "add" || lower == "remove" || lower == "extract" {
				actionable++
				break
			}
		}
	}
	score := (actionable * 10) / max(len(comments), 1)
	if score > 10 { score = 10 }
	return score
}

func GenerateFeedback(a ReviewAnalysis) string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("Review by: %s\n", a.Review.Reviewer))
	b.WriteString(fmt.Sprintf("Sentiment: %s\n", a.Sentiment))
	b.WriteString(fmt.Sprintf("Total comments: %d\n", a.TotalComments))
	b.WriteString(fmt.Sprintf("  Positive:    %d\n", a.PositiveCount))
	b.WriteString(fmt.Sprintf("  Critical:    %d\n", a.CriticalCount))
	b.WriteString(fmt.Sprintf("  Suggestion:  %d\n", a.SuggestionCount))
	b.WriteString(fmt.Sprintf("  Question:    %d\n", a.QuestionCount))
	b.WriteString(fmt.Sprintf("Actionability: %d/10\n", a.Actionable))
	if a.Actionable < 4 {
		b.WriteString("  Tip: Add more specific suggestions.\n")
	}
	if a.CriticalCount > a.PositiveCount {
		b.WriteString("  Tip: Balance critical feedback with positive observations.\n")
	}
	if a.Sentiment == "negative" {
		b.WriteString("  Warning: This review may feel harsh.\n")
	}
	return b.String()
}

func max(a, b int) int { if a > b { return a }; return b }

func main() {
	review := Review{
		Reviewer: "Alice",
		Comments: []ReviewComment{
			{Text: "Great use of interfaces to decouple the storage layer.", Category: "positive"},
			{Text: "This function is too long. Consider extracting the validation logic into a separate method.", Category: "suggestion"},
			{Text: "The error handling here swallows the underlying error. Suggest using fmt.Errorf with %w to wrap it.", Category: "critical"},
			{Text: "What happens when the context is cancelled mid-request? Should we clean up resources?", Category: "question"},
			{Text: "The test coverage is excellent. The table-driven test pattern is exactly right for this.", Category: "positive"},
			{Text: "This variable name is unclear. Suggest renaming 'tmp' to 'decodedPayload'.", Category: "suggestion"},
		},
	}
	analysis := AnalyzeReview(review)
	fmt.Print(GenerateFeedback(analysis))
}
```

Run with `go run .` to analyze a code review and get sentiment, actionability score, and improvement tips.

## Step-by-step execution

When `AnalyzeReview` runs:

1. Each comment is categorized by its category field: positive, critical, suggestion, or question.
2. Total counts per category are accumulated as the comment list is iterated.
3. Sentiment is computed as the ratio of positive to total comments (>50% positive = "positive", >40% critical = "negative", otherwise "balanced").
4. Actionability scans each comment's text for action verbs (consider, suggest, rename, extract, etc.). The score is the percentage of comments containing at least one action verb, scaled to 1-10.
5. Feedback is generated with concrete tips: if actionability is low, suggest adding specific suggestions. If critical comments exceed positive ones, suggest balancing.

For the example review (6 comments, 2 positive, 1 critical, 2 suggestions, 1 question): sentiment = 33% positive, 16% critical → "balanced". Actionability = 67% (4 out of 6 have action words) → 7/10.

## Common mistakes

- **Only pointing out problems without solutions.** "This is wrong" without explanation teaches nothing. Always pair criticism with a suggestion or question that points toward the fix.
- **Reviewing style instead of substance.** Focus on correctness, test coverage, error handling, and architecture before formatting. Use linters for style and let `gofmt` handle formatting.
- **Overwhelming the author with 50 comments.** Limit reviews to 10-15 substantive comments. If there are more issues, schedule a live walkthrough instead of a written review.
- **Using absolute language.** "You must change this" triggers resistance. "Consider changing this for consistency with the rest of the codebase" invites collaboration.
- **Delaying reviews.** A review after three days is nearly useless — the author has context-switched and the code may already be merged. Respond within 24 hours for routine reviews, same-day for blocking reviews.
- **Reviewing code outside your expertise.** If you don't understand the domain, ask clarifying questions instead of making incorrect suggestions. "Can you explain why this approach was chosen?" is a valid review comment.

## Debugging walkthrough

A junior engineer receives a code review with 15 critical comments and zero positive ones. They feel discouraged and stop responding to the review thread.

**Symptom**: PR stalled for two weeks. Author disengaged.

**Investigation**: The review analysis shows 15 critical, 0 positive, 0 suggestions. Sentiment = "negative". Actionability = low (comments like "this is wrong" without alternatives). The review fails the 3:1:1 ratio entirely.

**Root cause**: The reviewer used their expertise to identify everything wrong but did not practice mentorship. The author perceived the review as a personal attack and disengaged from fear of additional criticism.

**Fix**: The reviewer should have:

1. Started with positive observations: "I like the overall structure. The test coverage is a good start."
2. Prioritized the top 3-5 critical issues and framed them as teaching moments: "This error handling pattern swallows the stack trace. In Go, we use `fmt.Errorf("context: %w", err)` to wrap errors so callers can use `errors.Is` to check for specific conditions."
3. Used suggestions and questions instead of commands: "Have you considered using `sync.Map` here for the concurrent access pattern?" instead of "You need to fix the race condition."

**Lesson learned**: Code review is a relationship, not a transaction. Every review is an opportunity to build trust and teach. If the author feels attacked, the review failed regardless of its technical accuracy.

## Production notes

- **Review async-first, sync-second.** Start with a written review. If there are architectural concerns, schedule a 30-minute video call to discuss. Written reviews give the author time to process; sync discussions allow for rapid clarification.
- **Use a review checklist.** Maintain a team checklist for common concerns: error wrapping, context propagation, goroutine lifecycle, test coverage, security (SQL injection, XSS), and dependency changes. Automate what you can with `golangci-lint`.
- **Mentor through review assignments.** Pair junior reviewers with senior reviewers. The junior writes the first pass, the senior reviews the review. Both learn: the junior learns what to look for, the senior learns to delegate.
- **Track review metrics.** Measure review turnaround time, comments per review, and comment resolution rate. Teams with fast, thorough reviews ship more reliably. Use GitHub's review analytics or a custom dashboard.
- **Document team conventions.** Create a `REVIEW-GUIDE.md` that documents your team's specific review priorities, common patterns, and style preferences. New team members reference this guide for both writing and reviewing code.

## Performance implications

- **Review latency** is the bottleneck for most teams. A PR waiting 3 days for review blocks the author's next task. Keep individual reviews under 2 hours of wall-clock time. Batch reviews at set times (morning and afternoon).
- **Comment volume** correlates with defect density but has diminishing returns. Research from Google shows that 10-15 comments per 200 lines is the sweet spot. More comments produce marginal defect reduction but significantly increase review time.
- **Async reviews scale better than sync.** A written review takes the reviewer 30 minutes and can be processed by the author at their convenience. A sync review requires scheduling, context switching, and often produces less thorough feedback.

## Practice task

Build a `ReviewCoach` tool that helps reviewers improve their feedback. Implement:

1. `ReviewImprovement` struct with `Original string`, `Improved string`, `Reason string`.
2. `ImproveComment(original ReviewComment) ReviewImprovement` that applies improvements:
   - "This is wrong" → "This error handling pattern could swallow the underlying error. Consider using `fmt.Errorf` with `%w` to preserve the error chain."
   - "You should..." → "Consider..."
3. `SuggestToneImprovements(comments []ReviewComment) []string` that detects and suggests fixes for:
   - Comments starting with "You" (rephrase to focus on code)
   - Comments without a question or suggestion
   - Comments with loaded words ("obviously", "clearly", "just")
4. `GenerateCoachReport(r Review) string` that runs the reviewer's review through the coach and prints before/after for the three worst comments.
5. `main()` that takes a sample review, runs it through the coach, and prints the improvement report.

## Tests / verification

```bash
go run ./curriculum/modules/16-portfolio-interview-readiness/lessons/14-mentorship-and-code-review
go test ./curriculum/modules/16-portfolio-interview-readiness/lessons/14-mentorship-and-code-review
```

The existing tests verify `AnalyzeReview` (basic counts, empty review, sentiment detection across positive/negative/balanced, actionability scoring). After completing the practice task, add tests for `ImproveComment`, `SuggestToneImprovements`, and `GenerateCoachReport`.

## Review questions

1. What is the recommended ratio of positive to critical to suggestion comments in a code review?
2. Why should code review comments focus on teaching principles rather than just stating what to change?
3. What is the maximum number of substantive comments you should leave in a single code review?
4. Why is "You must change this" less effective than "Consider changing this"?
5. Name three categories of review concerns that should be caught by automated linters rather than manual review.

## NEXT UP

Career roadmap — learn how to plan your career growth, choose between IC and management paths, set meaningful goals, and build a continuous learning practice.
