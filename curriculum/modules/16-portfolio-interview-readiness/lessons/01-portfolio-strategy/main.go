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
