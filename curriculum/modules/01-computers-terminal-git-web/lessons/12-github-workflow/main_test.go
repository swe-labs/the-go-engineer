package main

import (
	"strings"
	"testing"
)

func TestPullRequest_CanMerge(t *testing.T) {
	pr := &PullRequest{
		ID:             42,
		Title:          "feat: add telemetry hooks",
		Author:         "dev-user",
		SourceBranch:   "feature/telemetry",
		TargetBranch:   "main",
		Status:         StatusPendingReview,
		ApprovalsCount: 0,
		CIPassing:      false,
	}

	// 1. CI failing -> cannot merge
	err := pr.CanMerge()
	if err == nil || !strings.Contains(err.Error(), "CI checks are failing") {
		t.Fatalf("expected CI failure error, got: %v", err)
	}

	// 2. CI passing but 0 approvals -> cannot merge
	pr.CIPassing = true
	err = pr.CanMerge()
	if err == nil || !strings.Contains(err.Error(), "approval is required") {
		t.Fatalf("expected approval required error, got: %v", err)
	}

	// 3. CI passing and 1 approval -> merge allowed
	pr.ApprovalsCount = 1
	pr.Status = StatusApproved
	if err := pr.CanMerge(); err != nil {
		t.Fatalf("expected CanMerge to pass, got: %v", err)
	}

	// 4. Execute Merge
	if err := pr.Merge(); err != nil {
		t.Fatalf("expected Merge to succeed, got: %v", err)
	}
	if pr.Status != StatusMerged {
		t.Errorf("expected PR status %s, got %s", StatusMerged, pr.Status)
	}
}

func TestMainExecution(t *testing.T) {
	main()
}

