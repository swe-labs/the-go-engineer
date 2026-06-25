# Career roadmap

## Learning objective

Create a structured career development plan by analyzing your current skills against target role requirements, identifying gaps, prioritizing learning areas, and building a timeline for professional growth as a Go engineer.

## Why this matters

Most engineers drift through their careers, reacting to opportunities rather than creating them. A deliberate career roadmap transforms luck into strategy. Engineers with written career plans grow faster, earn more, and report higher satisfaction than those without. In the Go ecosystem specifically, the landscape changes rapidly: new frameworks, deployment patterns, and architectural paradigms emerge yearly. A roadmap ensures you learn the right things at the right time rather than chasing every trend. Whether you aspire to Staff Engineer at a Go-heavy company or CTO of a startup, a written plan increases your odds.

## Mental model

Your career is a product. You are the engineer building it. Skills are features. Experience is the changelog. Your resume is the README. A career roadmap is the product roadmap — it defines what features to build, in what order, and by when. The market (employers) provides feedback through interviews, promotions, and salary offers. Your job is to iterate: measure where you are, identify the gap to where you want to be, and execute a learning plan to close that gap.

## Core idea

Career progression has two tracks: Individual Contributor (IC) and Engineering Management. Each has distinct skill requirements and promotion criteria:

| Level | IC focus | Management focus | Key skills |
|---|---|---|---|
| Junior (0-2 yrs) | Write correct code under guidance | — | Language fundamentals, testing, debugging |
| Mid (2-5 yrs) | Own features independently | Lead small projects | System design, code review, project estimation |
| Senior (5-8 yrs) | Design and deliver complex systems | Mentor 3-5 engineers | Architecture, cross-team collaboration, technical strategy |
| Staff (8-12 yrs) | Set technical direction for org | Manage 1-2 teams | Organizational influence, risk management, long-term planning |
| Principal (12+ yrs) | Industry-wide technical leadership | Director of engineering | Industry standards, open source leadership, executive communication |

The skill gap analysis evaluates where you are now vs where you need to be for your target level. Each skill is scored 1-10 for current and desired proficiency. The gap is the delta.

## Under the hood

Engineering promotion committees evaluate candidates along four axes:

1. **Impact** — What did you deliver? Quantify with metrics (reduced latency by 40%, shipped feature used by 1M users).
2. **Scope** — How complex was the work? Did it span teams, systems, or quarters?
3. **Skills** — What technical and leadership capabilities did you demonstrate? This is the skill gap analysis.
4. **Behaviors** — How did you work? Collaboration, mentorship, ownership, communication.

The promotion threshold is typically 80% readiness across all axes. Below 80%, the committee sees risk. Above 80%, they see readiness. Your career roadmap should target 80%+ on every dimension for your desired level.

## How Go uses it

The Go community has well-defined career levels at major employers:

- **Google** (Go's home): T3 (Junior SWE) → T4 (Mid) → T5 (Senior) → T6 (Staff) → T7 (Principal) — each level requires demonstrating the next level's skills for 6+ months before promotion.
- **Uber**: L3 (Junior) → L4 (Mid) → L5 (Senior) → L6 (Staff) → L7 (Principal) — Go engineers at Uber focus on microservices, distributed systems, and high-throughput APIs.
- **HashiCorp**: IC1-IC5 with clear rubrics for each level. Go is the primary language for all product engineering.

Each company publishes leveling guides internally. The common thread: senior+ requires system design, code review impact, and cross-team collaboration beyond individual coding.

## Go example

```go
package main

import (
	"fmt"
)

type Skill struct {
	Name     string
	Category string
	Current  int
	Desired  int
	Priority string
}

type CareerPlan struct {
	Name   string
	Role   string
	Track  string
	Skills []Skill
}

type SkillGap struct {
	Skill     Skill
	Gap       int
	Urgency   string
	Resources []string
}

type GapAnalysis struct {
	Plan          CareerPlan
	Gaps          []SkillGap
	TotalGap      int
	ReadinessPct  float64
	PriorityAreas []string
}

func AnalyzeGaps(plan CareerPlan) GapAnalysis {
	ga := GapAnalysis{Plan: plan}
	totalCurrent := 0
	totalDesired := 0
	for _, s := range plan.Skills {
		gap := s.Desired - s.Current
		if gap < 0 { gap = 0 }
		sg := SkillGap{Skill: s, Gap: gap, Urgency: computeUrgency(gap, s.Priority)}
		sg.Resources = suggestResources(sg)
		ga.Gaps = append(ga.Gaps, sg)
		totalCurrent += s.Current
		totalDesired += s.Desired
	}
	if len(plan.Skills) > 0 {
		ga.ReadinessPct = float64(totalCurrent) / float64(totalDesired) * 100
	}
	for _, g := range ga.Gaps {
		ga.TotalGap += g.Gap
		if g.Urgency == "critical" || g.Urgency == "high" {
			ga.PriorityAreas = append(ga.PriorityAreas, g.Skill.Name)
		}
	}
	return ga
}

func computeUrgency(gap int, priority string) string {
	switch {
	case gap >= 5 && priority == "high": return "critical"
	case gap >= 3: return "high"
	case gap >= 1: return "medium"
	default: return "none"
	}
}

func suggestResources(gap SkillGap) []string {
	resources := map[string][]string{
		"language":     {"Effective Go", "Go by Example", "Go Standard Library docs"},
		"tooling":      {"Go Tooling in Action (blog)", "golangci-lint docs", "Go Modules wiki"},
		"testing":      {"Go Test: Testing Techniques (talk)", "table-driven tests guide", "testing/slogans"},
		"architecture": {"Designing Data-Intensive Applications", "System Design Interview", "Go with Domain-Driven Design"},
		"soft":         {"Crucial Conversations", "The Manager's Path", "Radical Candor"},
		"devops":       {"Docker Deep Dive", "Kubernetes in Action", "Terraform: Up and Running"},
	}
	base := resources[gap.Skill.Category]
	if base == nil { base = []string{"Online courses", "Community mentorship", "Conference talks"} }
	return base
}

func GenerateRoadmap(ga GapAnalysis) string {
	roadmap := fmt.Sprintf("Career Roadmap for %s\n", ga.Plan.Name)
	roadmap += fmt.Sprintf("Target: %s %s\n\n", ga.Plan.Track, ga.Plan.Role)
	roadmap += fmt.Sprintf("Overall Readiness: %.0f%%\n", ga.ReadinessPct)
	roadmap += fmt.Sprintf("Total Skill Gap: %d points\n\n", ga.TotalGap)
	roadmap += "Priority Improvement Areas:\n"
	for _, area := range ga.PriorityAreas {
		roadmap += fmt.Sprintf("  - %s\n", area)
	}
	roadmap += "\nDetailed Gap Analysis:\n"
	roadmap += fmt.Sprintf("%-30s %-10s %-10s %-10s %-10s\n", "Skill", "Current", "Desired", "Gap", "Urgency")
	roadmap += fmt.Sprintf("%s\n", "------"+repeat("-", 70))
	for _, g := range ga.Gaps {
		roadmap += fmt.Sprintf("%-30s %-10d %-10d %-10d %-10s\n", g.Skill.Name, g.Skill.Current, g.Skill.Desired, g.Gap, g.Urgency)
	}
	roadmap += "\nRecommended Resources:\n"
	for _, g := range ga.Gaps {
		if g.Gap > 0 {
			roadmap += fmt.Sprintf("  %s:\n", g.Skill.Name)
			for _, r := range g.Resources { roadmap += fmt.Sprintf("    - %s\n", r) }
		}
	}
	roadmap += "\nSuggested Timeline:\n  0-3 months: Address critical gaps\n  3-6 months: Close high-urgency gaps\n  6-12 months: Improve medium-urgency gaps\n"
	return roadmap
}

func repeat(s string, n int) string { r := ""; for i := 0; i < n; i++ { r += s }; return r }

func main() {
	plan := CareerPlan{
		Name:  "Alex Gopher",
		Role:  "senior",
		Track: "ic",
		Skills: []Skill{
			{Name: "Go concurrency", Category: "language", Current: 7, Desired: 9, Priority: "high"},
			{Name: "Testing & fuzzing", Category: "testing", Current: 6, Desired: 8, Priority: "high"},
			{Name: "System design", Category: "architecture", Current: 5, Desired: 8, Priority: "high"},
			{Name: "Docker & Kubernetes", Category: "devops", Current: 4, Desired: 7, Priority: "medium"},
			{Name: "Mentoring & code review", Category: "soft", Current: 6, Desired: 8, Priority: "medium"},
			{Name: "Go tooling & modules", Category: "tooling", Current: 8, Desired: 9, Priority: "low"},
			{Name: "gRPC & protobuf", Category: "architecture", Current: 3, Desired: 6, Priority: "medium"},
		},
	}
	ga := AnalyzeGaps(plan)
	fmt.Print(GenerateRoadmap(ga))
}
```

Run with `go run .` to generate a personalized career roadmap with gap analysis, priority areas, and recommended learning resources.

## Step-by-step execution

When `AnalyzeGaps` runs:

1. The skills list is iterated. For each skill, a `SkillGap` is computed as `Desired - Current` (capped at 0 to avoid negative gaps for overqualified skills).
2. Urgency is determined by gap magnitude combined with priority: gap ≥ 5 with high priority = critical; gap ≥ 3 = high; gap ≥ 1 = medium; gap = 0 = none.
3. Resources are mapped from the skill's category (language, tooling, testing, architecture, soft, devops) to a curated list of books, talks, and guides.
4. Readiness percentage is total current divided by total desired across all skills. For Alex (Current = 7+6+5+4+6+8+3 = 39, Desired = 9+8+8+7+8+9+6 = 55): readiness = 39/55 = 70.9%.
5. Priority areas are collected from gaps with urgency "critical" or "high" — these become the focus of the 0-3 month timeline.
6. `GenerateRoadmap` formats the analysis as a structured roadmap with timeline phases.

For the example: total gap = (2+2+3+3+2+1+3) = 16 points. Critical: system design (gap 3, high priority). High: testing, Docker, mentoring, gRPC. The roadmap suggests focusing on system design first (0-3 months) followed by testing and Docker (3-6 months).

## Common mistakes

- **Setting vague goals.** "Become a Senior Go Engineer" is not a goal. It is a wish. A goal is: "Design and ship a distributed system feature end-to-end, lead its code review, and document the architecture decisions by Q3 2026."
- **Ignoring soft skills.** Technical skills get you to senior. Communication, mentorship, and leadership get you to staff and beyond. Half of your skill development should target soft skills after mid-level.
- **Learning without applying.** Watching conference talks and reading books without building projects produces knowledge without ability. For every learning resource, have a project that practices that skill.
- **Comparing to others.** Your career is not a race. Comparison creates anxiety without insight. Measure against your own past self and your target role, not against a peer's promotion timeline.
- **Staying too long in one role.** A role that no longer challenges you is costing you growth. If you have not learned something new in 6 months, it is time to move — internally or externally.
- **Neglecting the network.** Promotions and opportunities flow through relationships. Invest in your network before you need it. Attend meetups, comment on PRs, write blog posts, and stay in touch with former colleagues.

## Debugging walkthrough

An engineer runs the career roadmap tool and sees readiness at 45% with 10 priority areas.

**Symptom**: Too many priorities. The roadmap is overwhelming and suggests learning everything at once.

**Investigation**: The skill list contains 15 skills, all with high priority and desired level 10. The gap analysis shows every skill needs improvement. The urgency distribution: 6 critical, 4 high, 5 medium.

```go
plan := CareerPlan{
    Name: "Overwhelmed",
    Role: "senior",
    Skills: []Skill{
        {Name: "Go", Current: 5, Desired: 10, Priority: "high"},
        {Name: "Kubernetes", Current: 2, Desired: 10, Priority: "high"},
        {Name: "System Design", Current: 3, Desired: 10, Priority: "high"},
        // ... 12 more skills
    },
}
```

**Root cause**: The engineer set every desired level to 10 (expert) and every priority to high. This is unrealistic. Even Staff Engineers rarely score 10 on more than 3-4 skills. Setting 15 skills to desired 10 creates analysis paralysis.

**Fix**: Prioritize ruthlessly. For a senior role, identify the 5-7 most impactful skills and set realistic desired levels:

- Go concurrency: 8 (not 10) — strong proficiency, expert is unnecessary
- System design: 7 — solid understanding of distributed systems patterns
- Testing: 7 — comprehensive test strategies
- Kubernetes: 5 — enough to deploy and debug, no need to administer clusters
- Mentoring: 6 — start practicing with junior colleagues

Revised readiness: total current / total desired = 23/33 = 70% with only 2 priority areas. The roadmap becomes actionable.

**Lesson learned**: A roadmap with 10 priorities is a roadmap to burnout. Focus on 3-5 skills per quarter. Reach 80% readiness on the most important skills before expanding scope.

## Production notes

- **Write your roadmap down.** A goal written down is 42% more likely to be achieved according to a Dominican University study. Use a document, Notion, or a private GitHub repo. Review it quarterly.
- **Create a learning budget.** Dedicate 5 hours per week to deliberate skill development. 2 hours for theory (books, courses, talks), 2 hours for practice (side projects, open source), 1 hour for reflection (writing, mentoring).
- **Find a sponsor, not just a mentor.** A mentor gives advice. A sponsor gives opportunities — they recommend you for stretch assignments, promotions, and leadership roles. Cultivate sponsors at work.
- **Build public proof.** Your career roadmap means nothing without evidence. Write blog posts, speak at meetups, contribute to open source, and maintain a portfolio. The best interview is a URL to your work.
- **Re-evaluate annually.** Career goals change. What you wanted at 25 is different at 35. Schedule a yearly career strategy session. Archive old roadmaps and create new ones based on your current values and market conditions.

## Performance implications

- **Skill development has compounding returns.** Early investment in fundamentals (Go concurrency, system design, testing) pays exponential dividends. A deep understanding of channels and goroutines affects every concurrent system you build for the rest of your career.
- **Generalization vs specialization.** Early career (0-5 years): generalize across Go, databases, networking, and DevOps. Mid career (5-10 years): specialize in 2-3 areas (e.g., distributed systems + performance engineering). Late career (10+): specialize further or generalize into management.
- **The 10,000-hour rule applies.** True expertise in an area requires deliberate practice over years. Do not expect to become a Go performance expert in 3 months. Plan for multi-year skill development.
- **Market timing matters.** Learning a hot skill (Kubernetes in 2018, LLMs in 2024, WebAssembly in 2025) during its growth phase creates disproportionate career value. Watch the Go community trends (go.dev/blog, Go Time podcast, GopherCon talks) to identify emerging areas.

## Practice task

Build a `GoalTracker` that helps set and track quarterly career goals. Implement:

1. `Goal` struct with fields: `Description string`, `Category string`, `TargetDate string`, `Status string` ("not_started", "in_progress", "completed", "abandoned").
2. `QuarterlyPlan` struct with `Quarter string` (e.g., "2026-Q2"), `Goals []Goal`, and methods:
   - `CompletionRate() float64` — percentage of goals completed.
   - `CategoryBreakdown() map[string]int` — count of goals per category.
3. `SuggestNextQuarter(previous QuarterPlan) QuarterPlan` that:
   - Carries over incomplete goals from the previous quarter.
   - Adds 2-3 new goals based on common career gaps.
   - Adjusts deadlines for carried-over goals.
4. `FormatPlan(plan QuarterPlan) string` — returns a formatted markdown report.
5. `main()` that creates a Q1 plan, marks some goals complete, generates Q2 plan, and prints both.

## Tests / verification

```bash
go run ./curriculum/modules/16-portfolio-interview-readiness/lessons/15-career-roadmap
go test ./curriculum/modules/16-portfolio-interview-readiness/lessons/15-career-roadmap
```

The existing tests verify `AnalyzeGaps` (basic, no gaps, empty skills), `computeUrgency` (critical/high/medium/none), and `GenerateRoadmap` (formatting). After completing the practice task, add tests for `CompletionRate`, `CategoryBreakdown`, `SuggestNextQuarter`, and `FormatPlan`.

## Review questions

1. What is the difference between an IC track and a management track, and at what level do the tracks typically diverge?
2. What readiness percentage is typically required for promotion consideration, and across how many axes?
3. Why is it a mistake to set every skill's desired level to 10 when creating a career roadmap?
4. Name three ways to build public proof of your career development beyond your day job.
5. What three time horizons should a career roadmap span, and what should each focus on?

## NEXT UP

Congratulations on completing Module 16 and the entire core curriculum! You have built a comprehensive foundation in Go engineering — from fundamentals through production deployment and career readiness. Next up: Elective modules, starting with Module 18 — Flagship Opslane.
