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
	Sentiment       string // "positive", "balanced", "negative"
	Actionable      int    // 1-10, how actionable the feedback is
}

func AnalyzeReview(r Review) ReviewAnalysis {
	a := ReviewAnalysis{Review: r}
	for _, c := range r.Comments {
		a.TotalComments++
		switch c.Category {
		case "positive":
			a.PositiveCount++
		case "critical":
			a.CriticalCount++
		case "suggestion":
			a.SuggestionCount++
		case "question":
			a.QuestionCount++
		}
	}
	a.Sentiment = computeSentiment(a)
	a.Actionable = computeActionability(r.Comments)
	return a
}

func computeSentiment(a ReviewAnalysis) string {
	if a.TotalComments == 0 {
		return "neutral"
	}
	pos := float64(a.PositiveCount) / float64(a.TotalComments)
	crit := float64(a.CriticalCount) / float64(a.TotalComments)
	switch {
	case pos > 0.5:
		return "positive"
	case crit > 0.4:
		return "negative"
	default:
		return "balanced"
	}
}

func computeActionability(comments []ReviewComment) int {
	if len(comments) == 0 {
		return 0
	}
	actionable := 0
	for _, c := range comments {
		words := strings.Fields(c.Text)
		for _, w := range words {
			lower := strings.ToLower(w)
			if lower == "consider" || lower == "suggest" || lower == "recommend" || lower == "try" || lower == "change" || lower == "rename" || lower == "move" || lower == "add" || lower == "remove" || lower == "extract" {
				actionable++
				break
			}
		}
	}
	score := (actionable * 10) / max(len(comments), 1)
	if score > 10 {
		score = 10
	}
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
		b.WriteString("  Tip: Add more specific suggestions. Instead of 'this is wrong', say 'consider using sync.Mutex here'.\n")
	}
	if a.CriticalCount > a.PositiveCount {
		b.WriteString("  Tip: Balance critical feedback with positive observations. Start with what works.\n")
	}
	if a.Sentiment == "negative" {
		b.WriteString("  Warning: This review may feel harsh. Consider softening the tone while keeping the substance.\n")
	}
	return b.String()
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

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
