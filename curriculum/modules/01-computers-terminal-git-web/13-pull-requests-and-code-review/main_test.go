package main

import (
	"strings"
	"testing"
)

func TestNewCodeReview(t *testing.T) {
	diff := "+func new() {\n-func old() {\n func main() {"
	cr := NewCodeReview(diff)
	if cr == nil {
		t.Fatal("NewCodeReview returned nil")
	}
	if len(cr.Diff) != 3 {
		t.Errorf("expected 3 diff lines, got %d", len(cr.Diff))
	}
}

func TestLeaveComment(t *testing.T) {
	cr := NewCodeReview("+line1\n-line2\n line3")
	cr.LeaveComment(1, "reviewer", "Good change", false)
	if len(cr.Comments) != 1 {
		t.Errorf("expected 1 comment, got %d", len(cr.Comments))
	}
	if cr.Comments[0].Author != "reviewer" {
		t.Errorf("expected author 'reviewer', got %q", cr.Comments[0].Author)
	}
}

func TestBlockingCommentsPreventApproval(t *testing.T) {
	cr := NewCodeReview("+line1")
	cr.LeaveComment(1, "reviewer", "Must fix", true)
	cr.Approve()
	if cr.Approved {
		t.Error("expected approval to be blocked by unresolved changes")
	}
}

func TestResolveAndApprove(t *testing.T) {
	cr := NewCodeReview("+line1")
	cr.LeaveComment(1, "reviewer", "Must fix", true)
	cr.ResolveBlockingComments()
	cr.Approve()
	if !cr.Approved {
		t.Error("expected approval after resolving blocking comments")
	}
}

func TestDiffLineTypes(t *testing.T) {
	diff := "+add\n-remove\n context"
	cr := NewCodeReview(diff)
	if cr.Diff[0].Type != "addition" {
		t.Errorf("expected 'addition', got %q", cr.Diff[0].Type)
	}
	if cr.Diff[1].Type != "deletion" {
		t.Errorf("expected 'deletion', got %q", cr.Diff[1].Type)
	}
	if cr.Diff[2].Type != "context" {
		t.Errorf("expected 'context', got %q", cr.Diff[2].Type)
	}
}

func TestSummary(t *testing.T) {
	cr := NewCodeReview("+a\n-b\n c")
	s := cr.Summary()
	if !strings.Contains(s, "APPROVED") && !strings.Contains(s, "PENDING") {
		t.Errorf("unexpected summary: %q", s)
	}
}
