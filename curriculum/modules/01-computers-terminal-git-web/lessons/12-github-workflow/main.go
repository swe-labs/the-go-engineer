package main

import (
	"errors"
	"fmt"
)

type PRStatus string

const (
	StatusDraft            PRStatus = "Draft"
	StatusPendingReview    PRStatus = "Pending Review"
	StatusChangesRequested PRStatus = "Changes Requested"
	StatusApproved         PRStatus = "Approved"
	StatusMerged           PRStatus = "Merged"
)

// PullRequest simulates a GitHub pull request with branch protection rules.
type PullRequest struct {
	ID             int
	Title          string
	Author         string
	SourceBranch   string
	TargetBranch   string
	Status         PRStatus
	ApprovalsCount int
	CIPassing      bool
}

// CanMerge validates GitHub branch protection rules: CI must pass and at least 1 approval.
func (pr *PullRequest) CanMerge() error {
	if !pr.CIPassing {
		return errors.New("cannot merge: CI checks are failing")
	}
	if pr.ApprovalsCount < 1 {
		return errors.New("cannot merge: at least 1 approval is required by branch protection")
	}
	return nil
}

func (pr *PullRequest) Merge() error {
	if err := pr.CanMerge(); err != nil {
		return err
	}
	pr.Status = StatusMerged
	return nil
}

func main() {
	fmt.Println("=== GitHub Workflow State Machine Simulation ===")

	pr := &PullRequest{
		ID:           101,
		Title:        "feat: add rate limiting middleware",
		Author:       "alice",
		SourceBranch: "feature/rate-limiter",
		TargetBranch: "main",
		Status:       StatusPendingReview,
		CIPassing:    true,
	}

	fmt.Printf("PR #%d: %s (Status: %s)\n", pr.ID, pr.Title, pr.Status)
	fmt.Printf("Initial merge attempt: %v\n", pr.Merge())

	// Reviewer adds approval
	fmt.Println("\nReviewer 'bob' reviews and submits Approval...")
	pr.ApprovalsCount++
	pr.Status = StatusApproved
	fmt.Printf("PR Status: %s (Approvals: %d, CI: %t)\n", pr.Status, pr.ApprovalsCount, pr.CIPassing)

	// Second merge attempt
	if err := pr.Merge(); err != nil {
		fmt.Printf("Merge failed: %v\n", err)
	} else {
		fmt.Printf("✓ Merge successful! PR #%d is now merged into %s.\n", pr.ID, pr.TargetBranch)
	}
}
