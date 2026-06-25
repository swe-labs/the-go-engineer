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
	b.WriteString(fmt.Sprintf("Resume Keyword Analysis\n"))
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
