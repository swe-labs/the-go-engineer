# Resume and LinkedIn optimization

## Learning objective

Craft a resume and LinkedIn profile optimized for Go engineering roles by analyzing keyword coverage, structuring achievement-based bullets, and understanding how applicant tracking systems score your profile.

## Why this matters

Your resume and LinkedIn profile are the first filter in every job application. Before a human reads your resume, an applicant tracking system (ATS) likely scans it for keywords. If your profile scores below the threshold, no human ever sees it. A keyword-optimized resume that properly reflects your Go skills can triple your callback rate. This lesson gives you a tool to audit your resume and fix gaps before applying.

## Mental model

Think of your resume as a search engine optimization (SEO) problem. Recruiters and ATS systems search for specific terms: "Go", "Kubernetes", "gRPC", "distributed systems". Your resume is a web page that needs to rank for these keywords. But unlike SEO, keyword stuffing is penalized — every keyword must appear in the context of a real achievement. The keyword analyzer measures your coverage while the bullet-point structure ensures the keywords are backed by evidence.

## Core idea

A resume for a Go engineering role breaks down into four keyword categories:

**Languages**: Go, Golang, Rust, Python, TypeScript. These are hard filters for most roles.

**Technologies**: Docker, Kubernetes, AWS, GCP, PostgreSQL, Redis, gRPC, REST, microservices, CI/CD. These demonstrate your tooling experience.

**Soft skills**: Leadership, mentoring, collaboration, communication. These differentiate senior candidates.

**Methodologies**: Agile, Scrum, TDD. These show you can work in a team environment.

The resume analyzer scores your resume against these keywords and breaks down the coverage by category. A balanced resume has keywords from all four categories. A resume with only language keywords will not pass the screen for senior roles that require system design and leadership.

**Achievement-based bullets** follow the format: `Accomplished [X] by [Y] resulting in [Z]`. Example: `Reduced API latency by 40% by implementing an in-memory cache layer, improving p99 response times from 200ms to 120ms.` This format packs a keyword (cache, latency, API), an action (implemented), and a measurable result (40% reduction).

## Under the hood

The keyword analyzer uses Go's `regexp` package with word boundary matching (`\b`) to find exact keyword matches. Case-insensitive matching is used so "go" in "golang" does not produce false positives. The analyzer counts occurrences (frequency matters lightly) and categorizes matches.

In production, ATS systems like Greenhouse, Lever, and Workday use proprietary algorithms. Some use TF-IDF (term frequency-inverse document frequency) to weigh rare keywords more heavily — "gRPC" scores higher than "JavaScript" because fewer candidates list it. Others use embeddings to match semantic meaning. The simple regex-based analyzer approximates the keyword matching step but cannot capture semantic relevance.

## How Go uses it

Go's regex package (`regexp`) is exceptionally fast — it uses a RE2-based engine that guarantees linear time performance. This makes it ideal for scanning resume-length text (500-2000 words) in microseconds. The `strings` package provides additional utilities for case folding and substring checks.

The analyzer could be extended to use Go's NLP libraries:
- `github.com/jdkato/prose` for part-of-speech tagging to extract action verbs.
- `github.com/bbalet/stopwords` for filtering common words.
- Custom word frequency analysis with `map[string]int`.

## Go example

The resume keyword analyzer scans text for predefined Go ecosystem keywords and produces a categorized breakdown with match counts.

```go
package main

import (
	"fmt"
	"regexp"
	"strings"
)

type KeywordCategory int

const (
	LanguageKeyword KeywordCategory = iota
	TechnologyKeyword
	SoftSkillKeyword
	MethodologyKeyword
)

func (k KeywordCategory) String() string {
	switch k {
	case LanguageKeyword:
		return "Language"
	case TechnologyKeyword:
		return "Technology"
	case SoftSkillKeyword:
		return "Soft Skill"
	case MethodologyKeyword:
		return "Methodology"
	default:
		return "Unknown"
	}
}

type Keyword struct {
	Word     string
	Category KeywordCategory
}

type MatchResult struct {
	Keyword  string
	Category KeywordCategory
	Count    int
}

type ResumeScore struct {
	Text       string
	Results    []MatchResult
	TotalScore int
	MaxScore   int
}

var Keywords = []Keyword{
	{"Go", LanguageKeyword},
	{"Golang", LanguageKeyword},
	{"Rust", LanguageKeyword},
	{"Python", LanguageKeyword},
	{"Java", LanguageKeyword},
	{"JavaScript", LanguageKeyword},
	{"TypeScript", LanguageKeyword},
	{"Docker", TechnologyKeyword},
	{"Kubernetes", TechnologyKeyword},
	{"AWS", TechnologyKeyword},
	{"GCP", TechnologyKeyword},
	{"Azure", TechnologyKeyword},
	{"PostgreSQL", TechnologyKeyword},
	{"Redis", TechnologyKeyword},
	{"gRPC", TechnologyKeyword},
	{"REST", TechnologyKeyword},
	{"microservice", TechnologyKeyword},
	{"distributed", TechnologyKeyword},
	{"CI/CD", TechnologyKeyword},
	{"leadership", SoftSkillKeyword},
	{"mentor", SoftSkillKeyword},
	{"collaboration", SoftSkillKeyword},
	{"communication", SoftSkillKeyword},
	{"agile", MethodologyKeyword},
	{"scrum", MethodologyKeyword},
	{"TDD", MethodologyKeyword},
}

type ResumeAnalyzer struct {
	Keywords []Keyword
}

func NewResumeAnalyzer() *ResumeAnalyzer {
	return &ResumeAnalyzer{Keywords: Keywords}
}

func (ra *ResumeAnalyzer) Analyze(text string) ResumeScore {
	score := ResumeScore{Text: text}
	lower := strings.ToLower(text)

	for _, kw := range ra.Keywords {
		re := regexp.MustCompile(fmt.Sprintf(`(?i)\b%s\b`, regexp.QuoteMeta(kw.Word)))
		matches := re.FindAllString(lower, -1)
		count := len(matches)
		if count > 0 {
			score.Results = append(score.Results, MatchResult{
				Keyword:  kw.Word,
				Category: kw.Category,
				Count:    count,
			})
			score.TotalScore += count
		}
		score.MaxScore += 1
	}

	return score
}

func (rs ResumeScore) CategoryBreakdown() map[string]int {
	breakdown := make(map[string]int)
	for _, r := range rs.Results {
		cat := r.Category.String()
		breakdown[cat] += r.Count
	}
	return breakdown
}

func (rs ResumeScore) Summary() string {
	var b strings.Builder
	b.WriteString("Resume Keyword Analysis\n")
	b.WriteString(fmt.Sprintf("Matched: %d/%d keywords\n", len(rs.Results), rs.MaxScore))
	b.WriteString(fmt.Sprintf("Total matches: %d\n\n", rs.TotalScore))

	b.WriteString("Category Breakdown:\n")
	for cat, count := range rs.CategoryBreakdown() {
		b.WriteString(fmt.Sprintf("  %s: %d matches\n", cat, count))
	}

	if len(rs.Results) > 0 {
		b.WriteString("\nMatched Keywords:\n")
		for _, r := range rs.Results {
			b.WriteString(fmt.Sprintf("  %s (%s) x%d\n", r.Keyword, r.Category, r.Count))
		}
	}

	return b.String()
}

func main() {
	resume := `Experienced Go backend engineer with strong knowledge of Kubernetes and Docker.
Worked on distributed microservices using gRPC and REST APIs.
Used PostgreSQL and Redis for data storage. Deployed on AWS with CI/CD pipelines.
Led a team of 5 engineers, mentoring junior developers and fostering collaboration.`

	analyzer := NewResumeAnalyzer()
	score := analyzer.Analyze(resume)
	fmt.Print(score.Summary())
}
```

## Step-by-step execution

1. Define a set of keywords organized into four categories: Language, Technology, Soft Skill, Methodology.
2. `NewResumeAnalyzer()` initializes the analyzer with the default keyword set.
3. `Analyze()` receives resume text, converts to lowercase, and scans each keyword using a case-insensitive word-boundary regex.
4. For each match, it records the keyword, category, and occurrence count.
5. `TotalScore` accumulates the total match count. `MaxScore` tracks how many keywords exist in the dictionary.
6. `CategoryBreakdown()` groups matches by category for a high-level view of resume balance.
7. `Summary()` formats both the detailed match list and the category breakdown.

## Common mistakes

- Mistake: Listing "Go" in the skills section but never using it in a bullet point.
  - Why it happens: Skills section and experience section are treated separately.
  - Fix: Every Language keyword in the skills section should appear in at least one bullet point describing what you built with it.

- Mistake: Writing "Responsible for the backend" instead of "Designed and built the Go backend API serving 10K requests per second."
  - Why it happens: Passive voice is common in older resume formats.
  - Fix: Use active verbs: designed, built, implemented, optimized, led, architected.

- Mistake: Using acronyms without spelling them out (gRPC but never "gRPC (remote procedure call)").
  - Why it happens: Acronyms are second nature to engineers.
  - Fix: ATS systems match literal strings. Spell out acronyms at least once in the resume.

- Mistake: Including every technology you have ever touched.
  - Why it happens: Listing everything seems safer.
  - Fix: Only list technologies you can defend in a technical interview. Every keyword is fair game for follow-up questions.

## Debugging walkthrough

A developer runs the analyzer on their resume and gets:

```
Matched: 8/26 keywords
Category Breakdown:
  Technology: 6 matches
  Language: 3 matches
  Soft Skill: 0 matches
  Methodology: 0 matches
```

The missing soft skills and methodology categories suggest the resume reads like a technical spec, not a career narrative. The developer adds two bullet points: one about mentoring a junior engineer (soft skill) and one about adopting TDD on the team (methodology). Re-running the analyzer shows 11/26 matches with coverage in all four categories.

## Production notes

In a real job search, the keyword analyzer is one of several tools:

- ATS simulation: paste your resume into job descriptions to see keyword overlap.
- LinkedIn profile audit: check that your headline, summary, and experience sections are consistent.
- Bullet point power: every bullet should pass the "so what?" test with a measurable result.
- Customization: adjust the keyword weights per job description. A role at a cloud infrastructure company weights container and orchestration keywords higher.

Use the analyzer before submitting each application. Different roles require different keyword emphasis. A backend API role needs more gRPC and REST. An SRE role needs more Kubernetes and observability.

## Performance implications

The regex-based scanner runs in O(n*m) where n is the number of keywords and m is the resume length. With our 26 keywords and typical resumes of 1000 words, execution time is under 1 millisecond. For optimization, compile all regexes once in the constructor and reuse them. The `Regexp` type is safe for concurrent use once compiled.

## Practice task

Extend the analyzer to support custom keyword lists per job description. Add a `JobKeywords` method that accepts a job description string, extracts key terms (words appearing 3+ times), and generates a custom keyword list. Then run the analyzer against both the resume and the custom keywords to produce a "match score" for that specific job.

## Tests / verification

```bash
go test ./curriculum/modules/16-portfolio-interview-readiness/lessons/05-resume-and-linkedin-optimization
go run ./curriculum/modules/16-portfolio-interview-readiness/lessons/05-resume-and-linkedin-optimization
```

The tests verify keyword detection, correct category breakdown, accurate match counting, word boundary respect, and summary output format. After the practice task, test the custom keyword extraction against a sample job description.

## Review questions

1. What are the four keyword categories for a Go engineering resume, and why must each be represented?
2. What is the recommended bullet point format for achievement-based statements?
3. Why does "gRPC" score higher than "JavaScript" in an ATS for a Go role?
4. What is the risk of listing a technology you have only used once?
5. How would you optimize a resume that scores well on Technology keywords but misses Soft Skill keywords entirely?

## NEXT UP

Behavioral interview preparation — master the STAR method to tell compelling stories about your engineering experience.
