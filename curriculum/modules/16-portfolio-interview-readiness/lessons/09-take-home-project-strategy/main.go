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
