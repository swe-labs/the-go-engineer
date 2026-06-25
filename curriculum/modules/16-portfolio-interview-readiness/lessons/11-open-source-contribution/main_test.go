package main

import (
	"testing"
	"time"
)

func TestContributor_MergeRate(t *testing.T) {
	tests := []struct {
		name        string
		contributor Contributor
		expected    float64
	}{
		{
			name: "all merged",
			contributor: Contributor{
				PullRequests: []PullRequest{
					{ID: 1, MergedAt: timePtr(time.Now())},
					{ID: 2, MergedAt: timePtr(time.Now())},
				},
			},
			expected: 100.0,
		},
		{
			name: "half merged",
			contributor: Contributor{
				PullRequests: []PullRequest{
					{ID: 1, MergedAt: timePtr(time.Now())},
					{ID: 2, MergedAt: nil},
				},
			},
			expected: 50.0,
		},
		{
			name: "none merged",
			contributor: Contributor{
				PullRequests: []PullRequest{
					{ID: 1, MergedAt: nil},
				},
			},
			expected: 0.0,
		},
		{
			name:        "no PRs",
			contributor: Contributor{},
			expected:    0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.contributor.MergeRate()
			if got != tt.expected {
				t.Errorf("MergeRate() = %.1f, want %.1f", got, tt.expected)
			}
		})
	}
}

func TestContributor_TotalChanges(t *testing.T) {
	c := Contributor{
		PullRequests: []PullRequest{
			{Additions: 100, Deletions: 50},
			{Additions: 30, Deletions: 10},
		},
	}
	add, del := c.TotalChanges()
	if add != 130 {
		t.Errorf("additions = %d, want 130", add)
	}
	if del != 60 {
		t.Errorf("deletions = %d, want 60", del)
	}
}

func TestGenerateReport(t *testing.T) {
	now := time.Now()
	c := Contributor{
		Username: "testuser",
		PullRequests: []PullRequest{
			{ID: 1, MergedAt: timePtr(now), Reviewed: true},
			{ID: 2, MergedAt: timePtr(now), Reviewed: true},
			{ID: 3, MergedAt: nil, Reviewed: false},
		},
	}
	r := GenerateReport(c)
	if r.TotalPRs != 3 {
		t.Errorf("TotalPRs = %d, want 3", r.TotalPRs)
	}
	if r.MergedPRs != 2 {
		t.Errorf("MergedPRs = %d, want 2", r.MergedPRs)
	}
	if r.MergeRate != 66.66666666666666 {
		t.Errorf("MergeRate = %f, want 66.67", r.MergeRate)
	}
	if r.ReviewedCount != 2 {
		t.Errorf("ReviewedCount = %d, want 2", r.ReviewedCount)
	}
}
