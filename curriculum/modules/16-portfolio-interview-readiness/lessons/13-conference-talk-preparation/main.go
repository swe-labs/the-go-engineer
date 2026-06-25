package main

import (
	"fmt"
	"strings"
)

type TalkProposal struct {
	Title          string
	Abstract       string
	Category       string // "beginner", "intermediate", "advanced"
	DurationMin    int    // 20, 30, 45, or 60
	HasDemo        bool
	HasSlides      bool
	TargetAudience string
}

type ProposalScore struct {
	Proposal      TalkProposal
	Abstract      int // 1-10
	Relevance     int // 1-10
	Novelty       int // 1-10
	Structure     int // 1-10
	Actionability int // 1-10
	Total         int // weighted sum
	Feedback      []string
}

var scoreNames = []string{"Abstract", "Relevance", "Novelty", "Structure", "Actionability"}

func EvaluateProposal(p TalkProposal) ProposalScore {
	s := ProposalScore{Proposal: p}
	s.Abstract = scoreAbstract(p.Abstract)
	s.Relevance = scoreRelevance(p)
	s.Novelty = scoreNovelty(p)
	s.Structure = scoreStructure(p)
	s.Actionability = scoreActionability(p)
	s.Total = s.Abstract*3 + s.Relevance*2 + s.Novelty*2 + s.Structure*2 + s.Actionability*1
	s.Feedback = generateFeedback(s)
	return s
}

func scoreAbstract(abstract string) int {
	words := strings.Fields(abstract)
	wordCount := len(words)
	switch {
	case wordCount < 50:
		return 3
	case wordCount > 300:
		return 5
	default:
		return 8
	}
}

func scoreRelevance(p TalkProposal) int {
	keywords := []string{"go", "golang", "concurrency", "performance", "testing", "api", "microservice", "cloud", "kubernetes", "observability"}
	count := 0
	lower := strings.ToLower(p.Abstract + " " + p.Title)
	for _, kw := range keywords {
		if strings.Contains(lower, kw) {
			count++
		}
	}
	switch {
	case count >= 4:
		return 9
	case count >= 2:
		return 6
	default:
		return 3
	}
}

func scoreNovelty(p TalkProposal) int {
	noveltyWords := []string{"new", "novel", "real-world", "lessons", "case study", "deep dive", "under the hood", "internals", "unexpected", "surprising"}
	count := 0
	lower := strings.ToLower(p.Abstract + " " + p.Title)
	for _, w := range noveltyWords {
		if strings.Contains(lower, w) {
			count++
		}
	}
	switch {
	case count >= 3:
		return 9
	case count >= 1:
		return 6
	default:
		return 3
	}
}

func scoreStructure(p TalkProposal) int {
	structureWords := []string{"first", "second", "then", "finally", "part", "section", "step", "walk through", "overview", "agenda", "learn"}
	count := 0
	lower := strings.ToLower(p.Abstract)
	for _, w := range structureWords {
		if strings.Contains(lower, w) {
			count++
		}
	}
	switch {
	case count >= 4:
		return 9
	case count >= 2:
		return 6
	default:
		return 4
	}
}

func scoreActionability(p TalkProposal) int {
	actionWords := []string{"you will learn", "takeaway", "practical", "implement", "code", "example", "template", "pattern", "strategy", "guide"}
	count := 0
	lower := strings.ToLower(p.Abstract)
	for _, w := range actionWords {
		if strings.Contains(lower, w) {
			count++
		}
	}
	switch {
	case count >= 3:
		return 9
	case count >= 1:
		return 6
	default:
		return 3
	}
}

func generateFeedback(s ProposalScore) []string {
	var fb []string
	if s.Abstract < 6 {
		fb = append(fb, "Abstract is too short or too long. Aim for 100-250 words.")
	}
	if s.Novelty < 6 {
		fb = append(fb, "Emphasize what makes this talk unique. Use words like 'real-world', 'lessons learned', or 'case study'.")
	}
	if s.Structure < 6 {
		fb = append(fb, "Add structural signposts: 'First we'll cover X, then Y, finally Z.'")
	}
	if s.Actionability < 6 {
		fb = append(fb, "Include concrete takeaways. Tell the audience what they will learn or build.")
	}
	if s.Total < 50 {
		fb = append(fb, "Overall score is low. Consider revising the abstract to be more specific and actionable.")
	}
	if s.Proposal.DurationMin < 30 && s.Proposal.HasDemo {
		fb = append(fb, "Demo may be too long for a short talk slot. Consider a pre-recorded demo or trimmed walkthrough.")
	}
	return fb
}

func (s ProposalScore) Summary() string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("Proposal: %s\n", s.Proposal.Title))
	b.WriteString(fmt.Sprintf("Total Score: %d/100\n", s.Total))
	b.WriteString(fmt.Sprintf("  Abstract:       %d/10\n", s.Abstract))
	b.WriteString(fmt.Sprintf("  Relevance:      %d/10\n", s.Relevance))
	b.WriteString(fmt.Sprintf("  Novelty:        %d/10\n", s.Novelty))
	b.WriteString(fmt.Sprintf("  Structure:      %d/10\n", s.Structure))
	b.WriteString(fmt.Sprintf("  Actionability:  %d/10\n", s.Actionability))
	if len(s.Feedback) > 0 {
		b.WriteString("  Feedback:\n")
		for _, f := range s.Feedback {
			b.WriteString(fmt.Sprintf("    - %s\n", f))
		}
	}
	return b.String()
}

func main() {
	proposal := TalkProposal{
		Title:          "Building Observability into Go Microservices",
		Abstract:       "This talk covers real-world lessons from implementing distributed tracing, structured logging, and metrics in Go microservices at scale. First, we'll examine the three pillars of observability and how Go's standard library supports each. Then, we'll walk through instrumenting an HTTP service with OpenTelemetry, including context propagation and span attributes. Finally, we'll discuss common pitfalls and patterns for production-grade observability. You will learn concrete strategies for adding observability to existing Go services without major rewrites.",
		Category:       "intermediate",
		DurationMin:    45,
		HasDemo:        true,
		HasSlides:      true,
		TargetAudience: "Go developers building or maintaining microservices",
	}
	score := EvaluateProposal(proposal)
	fmt.Print(score.Summary())
}
