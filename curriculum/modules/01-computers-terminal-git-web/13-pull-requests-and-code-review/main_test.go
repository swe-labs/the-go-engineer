package main

import (
	"strings"
	"testing"
)

func TestNewCodeReview(t *testing.T) {
	diff := "+func new() {\n-func old() {\n func main() {"
	n := reviewDiff(diff)
	if n != 3 {
		t.Errorf("expected 3 diff lines, got %d", n)
	}
}

func TestLeaveComment(t *testing.T) {
	reviewDiff("+line1\n-line2\n line3")
	addComment(1, "reviewer", "Good change", false)
	if commentCount != 1 {
		t.Errorf("expected 1 comment, got %d", commentCount)
	}
	if commentAuthors[0] != "reviewer" {
		t.Errorf("expected author 'reviewer', got %q", commentAuthors[0])
	}
}

func TestBlockingCommentsPreventApproval(t *testing.T) {
	reviewDiff("+line1")
	addComment(1, "reviewer", "Must fix", true)
	ok := submitReview()
	if ok {
		t.Error("expected approval to be blocked by unresolved changes")
	}
}

func TestResolveAndApprove(t *testing.T) {
	reviewDiff("+line1")
	addComment(1, "reviewer", "Must fix", true)
	resolveBlocking()
	ok := submitReview()
	if !ok {
		t.Error("expected approval after resolving blocking comments")
	}
}

func TestDiffLineTypes(t *testing.T) {
	diff := "+add\n-remove\n context"
	n := reviewDiff(diff)
	if n != 3 {
		t.Errorf("expected 3 lines, got %d", n)
	}
	if diffTypes[0] != "addition" {
		t.Errorf("expected 'addition', got %q", diffTypes[0])
	}
	if diffTypes[1] != "deletion" {
		t.Errorf("expected 'deletion', got %q", diffTypes[1])
	}
	if diffTypes[2] != "context" {
		t.Errorf("expected 'context', got %q", diffTypes[2])
	}
}

func TestSummary(t *testing.T) {
	reviewDiff("+a\n-b\n c")
	s := summary()
	if !strings.Contains(s, "APPROVED") && !strings.Contains(s, "PENDING") {
		t.Errorf("unexpected summary: %q", s)
	}
}
