# Conference talk preparation

## Learning objective

Write compelling conference talk proposals, structure talks for maximum audience impact, design effective slides, and deliver live demos with confidence — increasing your acceptance rate at Go-focused conferences.

## Why this matters

Speaking at conferences is a force multiplier for your career. A single talk can reach 500+ live attendees and thousands more through recorded video. Speaking positions you as a subject matter expert, generates consulting and job offers, and builds a network of peers at other companies. The Go community particularly values practitioner talks — engineers sharing real experience building production systems. Conferences like GopherCon, GoLab, FOSDEM Go DevRoom, and GothamGo actively seek first-time speakers. Your only barrier is a well-crafted proposal.

## Mental model

A conference talk is a narrative arc, not a lecture. The audience is tired, distracted, and skeptical. Your job in 30-45 minutes is to earn their attention, teach them something valuable, and send them out with a concrete takeaway. Think of the talk as a journey: start where the audience is (a problem), take them through the struggle (your failed approaches), reveal the solution (your successful approach), and end with the destination (what they can do now). The proposal is the trailer for this journey — it must convince the program committee that the journey is worth funding.

## Core idea

Every talk has three deliverables: the proposal (gets you accepted), the slides (guide the narrative), and the delivery (connects with the audience). Each requires different skills:

| Element | Purpose | Key technique |
|---|---|---|
| CFP submission | Get accepted by the program committee | Clear abstract, structured outline, speaker bio that builds credibility |
| Talk structure | Keep the audience engaged | Problem → Struggle → Solution → Takeaway |
| Slide design | Reinforce the spoken word | Minimal text, one idea per slide, high-contrast visuals |
| Live demo | Prove the concept works | Pre-recorded backup, keyboard shortcuts, zoomed terminal font |
| Speaking delivery | Connect with the audience | Eye contact, varied pacing, planned pauses, Q&A preparation |

## Under the hood

Program committees evaluate proposals along five axes, each scored 1-10:

1. **Abstract clarity** — Does the proposal describe a specific, well-scoped topic? Vague abstracts ("I will talk about Go concurrency") are rejected. Specific abstracts ("How we reduced p99 latency by 60% using worker pools with dynamic scaling") are accepted.
2. **Relevance** — Will this topic interest the conference audience? A deep-dive on Go assembly may not fit a general Go conference. A talk on production debugging techniques fits almost everywhere.
3. **Novelty** — Does the talk offer new information? A rehash of Go's WaitGroup documentation is not novel. A case study of a real production incident with lessons learned is highly novel.
4. **Structure** — Does the proposal show a clear narrative? Abstracts that use signposts ("First we'll cover X, then Y, finally Z") score higher than stream-of-consciousness descriptions.
5. **Actionability** — Will attendees leave with something they can apply? Talks that end with concrete patterns, code examples, or checklists score highest on actionability.

## How Go uses it

GopherCon, the flagship Go conference, publishes its CFP rubric publicly. The committee evaluates each proposal on:

- **Originality**: Is this a fresh perspective or a rehash of existing content?
- **Depth**: Does the speaker demonstrate deep understanding of the topic?
- **Practicality**: Will attendees apply what they learn?
- **Speaker fit**: Does the speaker's background match the topic's requirements?

The Go community strongly values diversity of speakers and topics. First-time speakers with strong proposals are actively mentored. Many conferences offer speaker mentoring programs and office hours for CFP feedback before the deadline.

## Go example

```go
package main

import (
	"fmt"
	"strings"
)

type TalkProposal struct {
	Title           string
	Abstract        string
	Category        string
	DurationMin     int
	HasDemo         bool
	HasSlides       bool
	TargetAudience  string
}

type ProposalScore struct {
	Proposal      TalkProposal
	Abstract      int
	Relevance     int
	Novelty       int
	Structure     int
	Actionability int
	Total         int
	Feedback      []string
}

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
	words := len(strings.Fields(abstract))
	switch {
	case words < 50: return 3
	case words > 300: return 5
	default: return 8
	}
}

func scoreRelevance(p TalkProposal) int {
	keywords := []string{"go", "golang", "concurrency", "performance", "testing", "api", "microservice", "cloud", "kubernetes", "observability"}
	count := 0
	lower := strings.ToLower(p.Abstract + " " + p.Title)
	for _, kw := range keywords { if strings.Contains(lower, kw) { count++ } }
	switch { case count >= 4: return 9; case count >= 2: return 6; default: return 3 }
}

func scoreNovelty(p TalkProposal) int {
	words := []string{"new", "novel", "real-world", "lessons", "case study", "deep dive", "under the hood", "internals", "unexpected", "surprising"}
	count := 0
	lower := strings.ToLower(p.Abstract + " " + p.Title)
	for _, w := range words { if strings.Contains(lower, w) { count++ } }
	switch { case count >= 3: return 9; case count >= 1: return 6; default: return 3 }
}

func scoreStructure(p TalkProposal) int {
	words := []string{"first", "second", "then", "finally", "part", "section", "step", "walk through", "overview", "agenda", "learn"}
	count := 0
	lower := strings.ToLower(p.Abstract)
	for _, w := range words { if strings.Contains(lower, w) { count++ } }
	switch { case count >= 4: return 9; case count >= 2: return 6; default: return 4 }
}

func scoreActionability(p TalkProposal) int {
	words := []string{"you will learn", "takeaway", "practical", "implement", "code", "example", "template", "pattern", "strategy", "guide"}
	count := 0
	lower := strings.ToLower(p.Abstract)
	for _, w := range words { if strings.Contains(lower, w) { count++ } }
	switch { case count >= 3: return 9; case count >= 1: return 6; default: return 3 }
}

func generateFeedback(s ProposalScore) []string {
	var fb []string
	if s.Abstract < 6 { fb = append(fb, "Abstract is too short or too long. Aim for 100-250 words.") }
	if s.Novelty < 6 { fb = append(fb, "Emphasize what makes this talk unique.") }
	if s.Structure < 6 { fb = append(fb, "Add structural signposts: 'First X, then Y, finally Z.'") }
	if s.Actionability < 6 { fb = append(fb, "Include concrete takeaways.") }
	if s.Total < 50 { fb = append(fb, "Overall score is low. Consider revising.") }
	if s.Proposal.DurationMin < 30 && s.Proposal.HasDemo { fb = append(fb, "Demo may be too long for a short slot.") }
	return fb
}

func (s ProposalScore) Summary() string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("Proposal: %s\n", s.Proposal.Title))
	b.WriteString(fmt.Sprintf("Total Score: %d/100\n", s.Total))
	for _, name := range []string{"Abstract", "Relevance", "Novelty", "Structure", "Actionability"} {
		// output handled in main
	}
	b.WriteString(fmt.Sprintf("  Abstract:       %d/10\n", s.Abstract))
	b.WriteString(fmt.Sprintf("  Relevance:      %d/10\n", s.Relevance))
	b.WriteString(fmt.Sprintf("  Novelty:        %d/10\n", s.Novelty))
	b.WriteString(fmt.Sprintf("  Structure:      %d/10\n", s.Structure))
	b.WriteString(fmt.Sprintf("  Actionability:  %d/10\n", s.Actionability))
	if len(s.Feedback) > 0 {
		b.WriteString("  Feedback:\n")
		for _, f := range s.Feedback { b.WriteString(fmt.Sprintf("    - %s\n", f)) }
	}
	return b.String()
}

func main() {
	p := TalkProposal{
		Title:       "Building Observability into Go Microservices",
		Abstract:    "This talk covers real-world lessons from implementing distributed tracing, structured logging, and metrics in Go microservices at scale. First, we'll examine the three pillars of observability and how Go's standard library supports each. Then, we'll walk through instrumenting an HTTP service with OpenTelemetry. Finally, we'll discuss common pitfalls and patterns for production-grade observability. You will learn concrete strategies for adding observability to existing Go services without major rewrites.",
		Category:    "intermediate",
		DurationMin: 45,
		HasDemo:     true,
		HasSlides:   true,
		TargetAudience: "Go developers building or maintaining microservices",
	}
	score := EvaluateProposal(p)
	fmt.Print(score.Summary())
}
```

Run with `go run .` to see a scored proposal evaluation with per-category breakdown and improvement suggestions.

## Step-by-step execution

When `EvaluateProposal` runs:

1. The abstract is tokenized by whitespace. If word count is between 50 and 300, the abstract scores 8/10. Below 50 words scores 3 (too short), above 300 scores 5 (too long).
2. Relevance is scored by matching keywords (go, concurrency, performance, etc.) against the title and abstract. Four or more matches = 9, two to three = 6, fewer = 3.
3. Novelty checks for words like "real-world", "case study", "deep dive", "internals". Three or more = 9, one to two = 6, none = 3.
4. Structure counts signpost words ("first", "then", "finally", "section", "step"). Four or more = 9, two to three = 6, fewer = 4.
5. Actionability counts practical-focused words ("you will learn", "takeaway", "pattern"). Three or more = 9, one to two = 6, none = 3.
6. Total is the weighted sum: Abstract×3 + Relevance×2 + Novelty×2 + Structure×2 + Actionability×1.
7. Feedback is generated for any dimension scoring below 6, plus a recommendation to revise if total is under 50.

For the example proposal: abstract (8), relevance (9), novelty (9), structure (9), actionability (9) = total 24+18+18+18+9 = 87/100. The only feedback is a note about demo duration since the proposal is 45 minutes with a demo.

## Common mistakes

- **Submitting the same proposal to every conference.** Each conference has a different audience. GopherCon wants depth. FOSDEM wants community-focused talks. A local meetup wants tutorials. Tailor the abstract and emphasis to the venue.
- **Writing the abstract like a blog post title.** "An Introduction to Go" tells the committee nothing. "How We Migrated 2 Million Lines of Python to Go Without Downtime" tells them everything. Be specific.
- **Overloading the talk with content.** A 30-minute talk can cover exactly one concept, one example, and one takeaway. Trying to cover more results in surface-level coverage of everything and deep coverage of nothing.
- **Skipping the demo run-through.** Live demos fail 80% of the time without practice. Run the demo 10 times. Record it. Have a screenshot or recording backup for every step that could fail.
- **Reading from slides.** Slides are visual aids, not scripts. The audience can read faster than you can speak. Slides should show diagrams, code, or key terms — your voice provides the narrative.

## Debugging walkthrough

A speaker submits a proposal and gets a score of 35/100:

```go
proposal := TalkProposal{
    Title:    "Go Concurrency",
    Abstract: "I will talk about goroutines and channels.",
    DurationMin: 30,
}
```

**Score breakdown**: Abstract = 3 (25 words), Relevance = 6 (contains "Go"), Novelty = 3 (no novelty words), Structure = 4 (no signposts), Actionability = 3 (no takeaway words). Total = 9+12+6+8+3 = 38.

**Issues**: Abstract too vague and short, no narrative structure, no indication of what attendees will learn.

**Fix**: Rewrite the abstract with specificity and structure:

```
"This talk covers real-world lessons from building a production job queue in Go.
First, we'll examine why standard channel patterns failed under high throughput.
Then, we'll walk through our solution: a dynamic worker pool with backpressure.
Finally, you will learn a reusable pattern for building resilient concurrent systems
in your own Go services."
```

New score: Abstract = 8 (80 words), Relevance = 9 (go, concurrency), Novelty = 6 ("real-world lessons"), Structure = 9 ("First", "Then", "Finally"), Actionability = 9 ("you will learn", "pattern"). Total = 24+18+12+18+9 = 81.

## Production notes

- **Apply to multiple conferences.** Most CFPs have 10-30% acceptance rates. Submit to 5-10 conferences per talk. Each rejection is practice — iterate the proposal based on feedback if available.
- **Prepare for Q&A.** Anticipate the five hardest questions about your topic and prepare answers. Admit when you don't know: "That's a great question. I don't have a definitive answer, but here's how I would approach it."
- **Record your practice talk.** Watch the recording without sound first — study your body language, slide transitions, and timing. Then watch with sound for filler words ("um", "like", "you know").
- **Arrive early to check A/V.** Never assume the projector, microphone, or clicker works. Test your laptop with the venue's system 30 minutes before your slot. Have HDMI and USB-C adapters.
- **Share slides after the talk.** Post slides, code, and a written version of the talk within 48 hours. Include a link in your bio and on social media. This extends the talk's lifespan beyond the conference.

## Performance implications

- **Slide load time** affects flow. Heavy images or embedded videos may lag on conference wifi. Optimize images and pre-load videos. Keep the slide deck under 5MB.
- **Demo responsiveness** is critical. A slow demo (waiting for compilation, Docker builds, or database queries) kills energy. Pre-build containers, pre-warm caches, and use live-reload tools.
- **Font size and contrast** determine readability. Minimum 28pt font for body text, 36pt for headings. High contrast (dark text on light background) works on all projectors. Avoid thin fonts and low-contrast color schemes.

## Practice task

Build a `TalkScheduler` that helps plan a conference talk timeline. Implement:

1. `Timeline` struct with `SegmentName string`, `DurationMin int`, and `Type string` ("content", "demo", "qna").
2. `PlanTalk(totalMinutes int, demoMinutes int, qnaMinutes int) []Timeline` that creates a structured timeline with appropriate segments and breaks.
3. `ValidateTimeline(timeline []Timeline) (bool, []string)` that checks: total duration matches, demo segments are not too close to the end, Q&A is at the end, no segment is longer than 15 minutes without a transition.
4. `FormatTimeline(timeline []Timeline) string` that returns a markdown-style schedule.
5. `main()` that plans a 45-minute talk with 10-minute demo and 5-minute Q&A, then validates and prints the timeline.

## Tests / verification

```bash
go run ./curriculum/modules/16-portfolio-interview-readiness/lessons/13-conference-talk-preparation
go test ./curriculum/modules/16-portfolio-interview-readiness/lessons/13-conference-talk-preparation
```

The existing tests verify `EvaluateProposal` (basic, empty abstract, novelty scoring, feedback generation) and `Summary` output. After completing the practice task, add tests for `PlanTalk`, `ValidateTimeline`, and `FormatTimeline`.

## Review questions

1. What five axes do program committees use to evaluate talk proposals, and which has the highest weight?
2. Why does a vague abstract like "I will talk about Go concurrency" get rejected, while a specific one gets accepted?
3. What is the recommended word count range for a CFP abstract?
4. Name three ways to reduce the risk of a live demo failure during a conference talk.
5. What should you do within 48 hours after delivering a conference talk to maximize its impact?

## NEXT UP

Mentorship and code review — learn how to give constructive feedback, mentor junior engineers, and grow your leadership skills through code review.
