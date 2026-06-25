package main

import (
	"fmt"
	"strings"
)

type BlogSection struct {
	Heading    string
	WordCount  int
	CodeBlocks int
}

type BlogPost struct {
	Title    string
	Author   string
	Sections []BlogSection
}

type BlogQuality struct {
	Post            BlogPost
	TotalWords      int
	AvgSectionWords float64
	CodeRatio       float64
	HasIntro        bool
	HasConclusion   bool
	Readability     string
	Issues          []string
}

func AnalyzePost(post BlogPost) BlogQuality {
	q := BlogQuality{Post: post}
	totalCodeWords := 0
	for _, s := range post.Sections {
		q.TotalWords += s.WordCount
		totalCodeWords += s.CodeBlocks * 30
	}
	if len(post.Sections) > 0 {
		q.AvgSectionWords = float64(q.TotalWords) / float64(len(post.Sections))
	}
	if q.TotalWords > 0 {
		q.CodeRatio = float64(totalCodeWords) / float64(q.TotalWords) * 100
	}
	if len(post.Sections) > 0 && strings.Contains(strings.ToLower(post.Sections[0].Heading), "intro") {
		q.HasIntro = true
	}
	if len(post.Sections) > 0 && strings.Contains(strings.ToLower(post.Sections[len(post.Sections)-1].Heading), "conclusion") {
		q.HasConclusion = true
	}
	q.Readability = scoreReadability(q.TotalWords, len(post.Sections), totalCodeWords)
	q.Issues = checkIssues(post, q)
	return q
}

func scoreReadability(totalWords, sectionCount, codeWords int) string {
	avg := float64(totalWords) / float64(max(sectionCount, 1))
	switch {
	case totalWords < 300:
		return "too short"
	case totalWords > 5000:
		return "too long"
	case avg < 100:
		return "sections too short"
	case avg > 1500:
		return "sections too long"
	case float64(codeWords)/float64(max(totalWords, 1)) > 0.5:
		return "too much code"
	default:
		return "good"
	}
}

func checkIssues(post BlogPost, q BlogQuality) []string {
	var issues []string
	if !q.HasIntro {
		issues = append(issues, "Missing introduction section")
	}
	if !q.HasConclusion {
		issues = append(issues, "Missing conclusion section")
	}
	if q.CodeRatio > 50 {
		issues = append(issues, "Code blocks exceed 50% of content")
	}
	if q.TotalWords < 500 {
		issues = append(issues, "Post is under 500 words")
	}
	return issues
}

func GenerateOutline(post BlogPost) string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("# %s\n", post.Title))
	b.WriteString(fmt.Sprintf("By %s\n\n", post.Author))
	for i, s := range post.Sections {
		b.WriteString(fmt.Sprintf("## %s\n", s.Heading))
		b.WriteString(fmt.Sprintf("~%d words", s.WordCount))
		if s.CodeBlocks > 0 {
			b.WriteString(fmt.Sprintf(", %d code example(s)", s.CodeBlocks))
		}
		b.WriteString("\n\n")
		if i == 0 {
			b.WriteString("> Hook: Open with a problem statement or surprising fact\n\n")
		} else if i == len(post.Sections)-1 {
			b.WriteString("> Key takeaway: Summarize main points\n\n")
		} else {
			b.WriteString("> Explain concept with code\n\n")
		}
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
	post := BlogPost{
		Title:  "Building Production-Ready HTTP Services in Go",
		Author: "Gopher Engineer",
		Sections: []BlogSection{
			{Heading: "Introduction", WordCount: 150, CodeBlocks: 0},
			{Heading: "Project Structure", WordCount: 400, CodeBlocks: 3},
			{Heading: "Handler Patterns", WordCount: 350, CodeBlocks: 4},
			{Heading: "Middleware Chains", WordCount: 300, CodeBlocks: 3},
			{Heading: "Testing Strategies", WordCount: 450, CodeBlocks: 5},
			{Heading: "Conclusion", WordCount: 100, CodeBlocks: 0},
		},
	}

	q := AnalyzePost(post)
	fmt.Printf("Analysis for: %s\n", q.Post.Title)
	fmt.Printf("  Total words:     %d\n", q.TotalWords)
	fmt.Printf("  Avg per section: %.0f\n", q.AvgSectionWords)
	fmt.Printf("  Code ratio:      %.1f%%\n", q.CodeRatio)
	fmt.Printf("  Readability:     %s\n", q.Readability)
	fmt.Printf("  Has intro:       %v\n", q.HasIntro)
	fmt.Printf("  Has conclusion:  %v\n", q.HasConclusion)
	if len(q.Issues) > 0 {
		fmt.Println("  Issues:")
		for _, iss := range q.Issues {
			fmt.Printf("    - %s\n", iss)
		}
	} else {
		fmt.Println("  No issues found")
	}

	fmt.Println("\nGenerated Outline:")
	fmt.Print(GenerateOutline(post))
}
