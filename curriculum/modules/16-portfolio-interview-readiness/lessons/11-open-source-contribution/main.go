package main

import (
	"fmt"
	"time"
)

type PullRequest struct {
	ID         int
	Title      string
	Repository string
	Author     string
	CreatedAt  time.Time
	MergedAt   *time.Time
	Additions  int
	Deletions  int
	Comments   int
	Reviewed   bool
}

type Contributor struct {
	Username     string
	PullRequests []PullRequest
}

func (c Contributor) TotalPRs() int {
	return len(c.PullRequests)
}

func (c Contributor) MergedPRs() int {
	count := 0
	for _, pr := range c.PullRequests {
		if pr.MergedAt != nil {
			count++
		}
	}
	return count
}

func (c Contributor) MergeRate() float64 {
	if len(c.PullRequests) == 0 {
		return 0
	}
	return float64(c.MergedPRs()) / float64(len(c.PullRequests)) * 100
}

func (c Contributor) TotalChanges() (additions, deletions int) {
	for _, pr := range c.PullRequests {
		additions += pr.Additions
		deletions += pr.Deletions
	}
	return
}

func (c Contributor) AverageTimeToMerge() time.Duration {
	var total time.Duration
	var count int
	for _, pr := range c.PullRequests {
		if pr.MergedAt != nil {
			total += pr.MergedAt.Sub(pr.CreatedAt)
			count++
		}
	}
	if count == 0 {
		return 0
	}
	return total / time.Duration(count)
}

type ContributionReport struct {
	Contributor   Contributor
	TotalPRs      int
	MergedPRs     int
	MergeRate     float64
	Additions     int
	Deletions     int
	AvgMergeTime  time.Duration
	ReviewedCount int
}

func GenerateReport(c Contributor) ContributionReport {
	additions, deletions := c.TotalChanges()
	reviewed := 0
	for _, pr := range c.PullRequests {
		if pr.Reviewed {
			reviewed++
		}
	}
	return ContributionReport{
		Contributor:   c,
		TotalPRs:      c.TotalPRs(),
		MergedPRs:     c.MergedPRs(),
		MergeRate:     c.MergeRate(),
		Additions:     additions,
		Deletions:     deletions,
		AvgMergeTime:  c.AverageTimeToMerge(),
		ReviewedCount: reviewed,
	}
}

func main() {
	now := time.Now()
	contributor := Contributor{
		Username: "gopher-engineer",
		PullRequests: []PullRequest{
			{ID: 101, Title: "Add request validation middleware", Repository: "go-api-starter", Author: "gopher-engineer", CreatedAt: now.Add(-72 * time.Hour), MergedAt: timePtr(now.Add(-48 * time.Hour)), Additions: 120, Deletions: 30, Comments: 4, Reviewed: true},
			{ID: 102, Title: "Fix race condition in worker pool", Repository: "go-api-starter", Author: "gopher-engineer", CreatedAt: now.Add(-120 * time.Hour), MergedAt: timePtr(now.Add(-96 * time.Hour)), Additions: 45, Deletions: 15, Comments: 6, Reviewed: true},
			{ID: 103, Title: "Update README with deployment guide", Repository: "docs", Author: "gopher-engineer", CreatedAt: now.Add(-24 * time.Hour), MergedAt: timePtr(now.Add(-12 * time.Hour)), Additions: 200, Deletions: 10, Comments: 2, Reviewed: true},
			{ID: 104, Title: "Implement rate limiting middleware", Repository: "go-api-starter", Author: "gopher-engineer", CreatedAt: now.Add(-12 * time.Hour), MergedAt: nil, Additions: 80, Deletions: 5, Comments: 3, Reviewed: false},
		},
	}

	report := GenerateReport(contributor)
	fmt.Printf("Contribution Report for %s\n", report.Contributor.Username)
	fmt.Printf("  Total PRs:        %d\n", report.TotalPRs)
	fmt.Printf("  Merged PRs:       %d\n", report.MergedPRs)
	fmt.Printf("  Merge Rate:       %.1f%%\n", report.MergeRate)
	fmt.Printf("  Additions:        %d\n", report.Additions)
	fmt.Printf("  Deletions:        %d\n", report.Deletions)
	fmt.Printf("  Avg Time to Merge: %s\n", report.AvgMergeTime.Round(time.Hour))
	fmt.Printf("  Reviewed PRs:     %d\n", report.ReviewedCount)
}

func timePtr(t time.Time) *time.Time {
	return &t
}
