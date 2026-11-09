package main

import (
	"fmt"
)

type Skill struct {
	Name     string
	Category string // "language", "tooling", "testing", "architecture", "soft", "devops"
	Current  int    // 1-10
	Desired  int    // 1-10
	Priority string // "high", "medium", "low"
}

type CareerPlan struct {
	Name   string
	Role   string // "junior", "mid", "senior", "staff", "principal"
	Track  string // "ic" or "management"
	Skills []Skill
}

type GapAnalysis struct {
	Plan          CareerPlan
	Gaps          []SkillGap
	TotalGap      int
	ReadinessPct  float64
	PriorityAreas []string
}

type SkillGap struct {
	Skill     Skill
	Gap       int
	Urgency   string
	Resources []string
}

func AnalyzeGaps(plan CareerPlan) GapAnalysis {
	ga := GapAnalysis{Plan: plan}
	totalCurrent := 0
	totalDesired := 0
	for _, s := range plan.Skills {
		gap := s.Desired - s.Current
		if gap < 0 {
			gap = 0
		}
		sg := SkillGap{
			Skill:   s,
			Gap:     gap,
			Urgency: computeUrgency(gap, s.Priority),
		}
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
	case gap >= 5 && priority == "high":
		return "critical"
	case gap >= 3:
		return "high"
	case gap >= 1:
		return "medium"
	default:
		return "none"
	}
}

func suggestResources(gap SkillGap) []string {
	resourceMap := map[string][]string{
		"language":     {"Effective Go", "Go by Example", "Go Standard Library docs"},
		"tooling":      {"Go Tooling in Action (blog)", "golangci-lint docs", "Go Modules wiki"},
		"testing":      {"Go Test: Testing Techniques (talk)", "table-driven tests guide", "testing/slogans"},
		"architecture": {"Designing Data-Intensive Applications", "System Design Interview (book)", "Go with Domain-Driven Design"},
		"soft":         {"Crucial Conversations", "The Manager's Path", "Radical Candor"},
		"devops":       {"Docker Deep Dive", "Kubernetes in Action", "Terraform: Up and Running"},
	}
	base := resourceMap[gap.Skill.Category]
	if base == nil {
		base = []string{"General online courses", "Community mentorship", "Conference talks"}
	}
	return base
}

func GenerateRoadmap(ga GapAnalysis) string {
	plan := ga.Plan
	roadmap := fmt.Sprintf("Career Roadmap for %s\n", plan.Name)
	roadmap += fmt.Sprintf("Target: %s %s\n\n", plan.Track, plan.Role)
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
		roadmap += fmt.Sprintf("%-30s %-10d %-10d %-10d %-10s\n",
			g.Skill.Name, g.Skill.Current, g.Skill.Desired, g.Gap, g.Urgency)
	}

	roadmap += "\nRecommended Resources:\n"
	for _, g := range ga.Gaps {
		if g.Gap > 0 {
			roadmap += fmt.Sprintf("  %s:\n", g.Skill.Name)
			for _, r := range g.Resources {
				roadmap += fmt.Sprintf("    - %s\n", r)
			}
		}
	}

	roadmap += "\nSuggested Timeline:\n"
	roadmap += "  0-3 months: Address critical gaps\n"
	roadmap += "  3-6 months: Close high-urgency gaps\n"
	roadmap += "  6-12 months: Improve medium-urgency gaps\n"

	return roadmap
}

func repeat(s string, n int) string {
	result := ""
	for i := 0; i < n; i++ {
		result += s
	}
	return result
}

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
