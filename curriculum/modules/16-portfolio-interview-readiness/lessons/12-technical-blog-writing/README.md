# Technical blog writing

## Learning objective

Write clear, engaging technical blog posts about Go that effectively explain concepts, demonstrate code, and connect with your target audience — using a repeatable structure that maximizes readability and impact.

## Why this matters

Technical blogging is one of the highest-leverage activities for a Go engineer. A well-written post can reach thousands of engineers, generate job offers, establish speaking invitations, and build your professional brand. Blog posts serve as extended interview answers: when an interviewer asks about your experience with goroutines or HTTP middleware, you can point to a post that demonstrates both depth and communication skill. Companies like Honeycomb, Grafana, and Tailscale actively recruit engineers who blog about their technical work.

## Mental model

A technical blog post is a teaching conversation with a motivated but uninformed reader. Your job is to take them from confusion to clarity in 5-10 minutes. The post has three phases: hook (why should I care?), teach (here is how it works), and land (here is what you can do now). Each section serves one of these phases. Code is not decoration — it is the primary teaching mechanism. Every code block should compile, demonstrate one idea, and be referenced by the surrounding text.

## Core idea

A technical blog post follows a standard structure optimized for developer reading patterns:

| Element | Purpose | Length |
|---|---|---|
| Title | Attract the right reader, signal the topic | 6-12 words |
| Introduction | Hook with a problem or surprising fact | 100-200 words |
| Background | Establish prerequisite knowledge | 200-400 words |
| Core explanation | Teach the concept with code | 500-2000 words |
| Code walkthrough | Step through the example | 300-800 words |
| Best practices / pitfalls | Practical guidance | 200-400 words |
| Conclusion | Key takeaways and next steps | 100-200 words |

The most important rule: one post, one concept. If you try to cover goroutines, channels, select statements, and the worker pool pattern in a single post, the reader learns nothing. Scope down to the smallest complete teaching unit.

## Under the hood

Readers scan before they read. Eye-tracking studies show developers jump to code blocks first, then headings, then bullet lists. Design your post for scanning:

1. **Code blocks are anchors.** Every code block should be independently understandable. Add comments inside code blocks for readers who skip the prose.
2. **Headings tell the story.** A reader should understand the post's argument by reading only the headings. Make them descriptive: "Why buffered channels can hide bugs" not "Channels."
3. **Paragraphs are one idea.** Three sentences max per paragraph. Developers read on small screens with limited attention. Short paragraphs keep them moving.
4. **Images explain relationships.** Architecture diagrams, data flow charts, and benchmark comparisons communicate in seconds what text takes minutes. Use Go's `image` package or a diagram-as-code tool like Mermaid.

## How Go uses it

The Go blog (`go.dev/blog`) is the gold standard for technical writing. Every post follows the same pattern: a real problem, minimal code, clear explanation, and a conclusion that points to further reading. Notable examples:

- **"Go Concurrency Patterns: Pipelines and cancellation"** — Teaches complex channel patterns through a single running example (MD5 hashing).
- **"Using Go Modules"** — A tutorial that walks through the exact commands and expected output, building from simple to complex.
- **"The Go Memory Model"** — A reference-grade post that uses precise language and minimal code to describe the happens-before relationship.

These posts succeed because they respect the reader's time, use code as evidence, and never assume the reader has context they haven't provided.

## Go example

```go
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
	switch {
	case totalWords < 300:
		return "too short"
	case totalWords > 5000:
		return "too long"
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
			b.WriteString("> Hook: Open with a problem statement\n\n")
		} else if i == len(post.Sections)-1 {
			b.WriteString("> Key takeaway: Summarize main points\n\n")
		} else {
			b.WriteString("> Explain concept with code\n\n")
		}
	}
	return b.String()
}

func max(a, b int) int { if a > b { return a }; return b }

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
	fmt.Printf("  Code ratio:      %.1f%%\n", q.CodeRatio)
	fmt.Printf("  Readability:     %s\n", q.Readability)
	fmt.Printf("  Has intro:       %v\n", q.HasIntro)
	fmt.Printf("  Has conclusion:  %v\n", q.HasConclusion)
	if len(q.Issues) > 0 {
		fmt.Println("  Issues:")
		for _, iss := range q.Issues {
			fmt.Printf("    - %s\n", iss)
		}
	}
	fmt.Println("\nGenerated Outline:")
	fmt.Print(GenerateOutline(post))
}
```

Run with `go run .` to analyze a sample blog post and generate a structured outline.

## Step-by-step execution

When `AnalyzePost` runs:

1. The post's sections are iterated. Each section's word count is summed to compute `TotalWords`. Code blocks are estimated at 30 words each for ratio calculation.
2. `CodeRatio` is computed as code-words divided by total-words, multiplied by 100. A ratio above 50% triggers a code-overload issue.
3. The first section heading is checked for "intro" keywords (case-insensitive) to set `HasIntro`. The last section heading is checked for "conclusion" keywords.
4. `scoreReadability` applies quality thresholds: posts under 300 words are too short, over 5000 are too long, code-dominant posts are flagged.
5. Issues are collected as a list of human-readable strings for the author to address.
6. `GenerateOutline` produces a markdown outline with estimated word counts, code block counts, and structural guidance for each section.

For the example post (1750 words, 15 code blocks): total words = 1750, code ratio = (450/1750)*100 = 25.7%, readability = "good", no issues.

## Common mistakes

- **Writing without an outline first.** Without structure, posts meander. Always start with headings and bullet points for each section. The outline becomes the table of contents.
- **Including non-compiling code.** A single typo in a code block erodes trust. Copy-paste your code from a running test. Better: write the code as a test first, then include it in the post.
- **Explaining the obvious.** Do not explain `fmt.Println` in a post about gRPC. Know your audience and skip the basics they already understand. Assume they know Go syntax but not your specific topic.
- **No real examples.** Abstract explanations without concrete code leave the reader unable to apply the concept. Every paragraph should either teach a concept or show it in code.
- **Skipping the edit pass.** First drafts are for discovery. Second drafts are for clarity. Third drafts are for cutting. Cut every word that does not serve the teaching goal. Aim to remove 20% of the words in each edit pass.

## Debugging walkthrough

An engineer posts a blog titled "Understanding Go Interfaces" and gets feedback that readers find it confusing. They run the quality checker:

```go
post := BlogPost{
    Title:  "Understanding Go Interfaces",
    Author: "Dev",
    Sections: []BlogSection{
        {Heading: "What is an Interface", WordCount: 500, CodeBlocks: 1},
        {Heading: "Interface Values", WordCount: 600, CodeBlocks: 2},
        {Heading: "Empty Interface", WordCount: 400, CodeBlocks: 1},
        {Heading: "Type Assertions", WordCount: 700, CodeBlocks: 3},
    },
}
```

**Analysis output**: Total words = 2200, Code ratio = (7*30)/2200 = 9.5%, Readability = "good", Has intro = false, Has conclusion = false.

**Issues detected**: Missing introduction section, Missing conclusion section.

**Root cause**: The post jumps directly into technical content without hooking the reader or explaining why interfaces matter. The reader needs motivation before mechanics. Without a conclusion, the reader has no summary of key takeaways.

**Fix**: Add an introduction section titled "Introduction" with 150-200 words on why interfaces solve real problems (testability, decoupling, polymorphism). Add a conclusion titled "Conclusion" with the top three lessons.

## Production notes

- **Publish on multiple platforms.** Write in markdown, publish on your own site (using Hugo, a Go-based SSG), then cross-post to Dev.to, Medium, or Go Time. Each platform reaches a different audience.
- **Include a "Run it yourself" section.** Link to a Go Playground or GitHub repo so readers can execute the code without copying. This dramatically increases engagement.
- **SEO matters for discoverability.** Use the primary keyword in the title, first paragraph, and one H2 heading. Include `alt` text on images. Write a meta-description under 160 characters.
- **Respond to comments.** Comments are free editing. A clarifying question reveals a gap in your explanation. Update the post to address it. This builds both the post and your reputation.
- **Repurpose content.** A blog post can become a conference talk, a workshop exercise, a newsletter issue, or a chapter in an ebook. Plan the post so it can be extended in multiple directions.

## Performance implications

- **Page load time** affects bounce rate. A blog with heavy JavaScript or large images loses mobile readers. Use static site generation (Hugo) and optimized images.
- **Code syntax highlighting** should be server-side rendered, not JavaScript. Highlight.js with server-side rendering avoids layout shift and works without JavaScript.
- **Go Playground snippets** load asynchronously and may fail if the playground is down. Embed the code directly as a fallback so the post is self-contained.

## Practice task

Build a `BlogScorer` that provides detailed quality metrics. Implement:

1. `ReadTime(totalWords int) string` — returns estimated read time (avg 200 words/min).
2. `HeadingHierarchy(sections []BlogSection) ([]string, error)` — checks that headings follow a logical H1→H2→H3 hierarchy with no jumps.
3. `CodeCompiles(section BlogSection, code string) bool` — attempts to parse the code string with `go/parser`, returns whether it's syntactically valid Go.
4. `GeneratePostSummary(post BlogPost) string` — produces a 2-3 sentence summary by concatenating first and last section content and key code concepts.
5. `main()` that scores a sample post, generates a summary, and prints improvement suggestions.

## Tests / verification

```bash
go run ./curriculum/modules/16-portfolio-interview-readiness/lessons/12-technical-blog-writing
go test ./curriculum/modules/16-portfolio-interview-readiness/lessons/12-technical-blog-writing
```

The existing tests verify `AnalyzePost` (basic, missing sections, issues detection, readability scoring) and `GenerateOutline` (title, headings, word counts). After completing the practice task, add tests for `ReadTime`, `HeadingHierarchy`, `CodeCompiles`, and `GeneratePostSummary`.

## Review questions

1. What are the three phases of a technical blog post, and what does each accomplish?
2. Why is it important to limit a post to one concept rather than covering everything you know?
3. What is the most important rule for code blocks in a technical blog post?
4. Name three issues detected by the blog quality checker in the Go example.
5. Why should you write the outline before writing the full post?

## NEXT UP

Conference talk preparation — learn how to craft compelling talk proposals, design effective slides, and deliver presentations that showcase your expertise.
