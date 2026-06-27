package main

import (
	"fmt"
	"strings"
)

type DiffLine struct {
	Number    int
	Content   string
	Type      string
}

type ReviewComment struct {
	Line      int
	Author    string
	Body      string
	Blocking  bool
}

type CodeReview struct {
	Diff         []DiffLine
	Comments     []ReviewComment
	Approved     bool
	ChangesReq   bool
}

func NewCodeReview(diffContent string) *CodeReview {
	lines := strings.Split(diffContent, "\n")
	cr := &CodeReview{Diff: make([]DiffLine, 0, len(lines))}
	for i, line := range lines {
		dt := "context"
		if strings.HasPrefix(line, "+") {
			dt = "addition"
		} else if strings.HasPrefix(line, "-") {
			dt = "deletion"
		}
		cr.Diff = append(cr.Diff, DiffLine{Number: i + 1, Content: line, Type: dt})
	}
	return cr
}

func (cr *CodeReview) LeaveComment(line int, author, body string, blocking bool) {
	cr.Comments = append(cr.Comments, ReviewComment{Line: line, Author: author, Body: body, Blocking: blocking})
	if blocking {
		cr.ChangesReq = true
		cr.Approved = false
	}
}

func (cr *CodeReview) Approve() {
	if cr.ChangesReq {
		return
	}
	cr.Approved = true
}

func (cr *CodeReview) ResolveBlockingComments() {
	cr.ChangesReq = false
	for i := range cr.Comments {
		cr.Comments[i].Blocking = false
	}
}

func (cr *CodeReview) Summary() string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("Diff: %d lines (%d additions, %d deletions)\n", len(cr.Diff), cr.countType("addition"), cr.countType("deletion")))
	b.WriteString(fmt.Sprintf("Comments: %d (%d blocking)\n", len(cr.Comments), cr.countBlocking()))
	if cr.Approved {
		b.WriteString("Status: APPROVED\n")
	} else if cr.ChangesReq {
		b.WriteString("Status: CHANGES REQUESTED\n")
	} else {
		b.WriteString("Status: PENDING REVIEW\n")
	}
	return b.String()
}

func (cr *CodeReview) countType(t string) int {
	count := 0
	for _, d := range cr.Diff {
		if d.Type == t {
			count++
		}
	}
	return count
}

func (cr *CodeReview) countBlocking() int {
	count := 0
	for _, c := range cr.Comments {
		if c.Blocking {
			count++
		}
	}
	return count
}

func main() {
	diff := ` package main

 import "fmt"

-func old() {
-    fmt.Println("old")
+func new() {
+    fmt.Println("new")
+    // TODO: add error handling
 }

 func main() {
-    old()
+    new()
 }`

	review := NewCodeReview(diff)

	fmt.Println("=== Code Review Simulation ===")
	fmt.Println()
	fmt.Println("Submitted Diff:")
	for _, d := range review.Diff {
		mark := " "
		if d.Type == "addition" {
			mark = "+"
		} else if d.Type == "deletion" {
			mark = "-"
		}
		fmt.Printf("  %s %s\n", mark, d.Content)
	}
	fmt.Println()

	review.LeaveComment(5, "alice", "Missing error handling on new() call", true)
	review.LeaveComment(9, "bob", "Consider adding a unit test for this path", false)

	fmt.Println("Review Comments:")
	for _, c := range review.Comments {
		blocking := ""
		if c.Blocking {
			blocking = " [BLOCKING]"
		}
		fmt.Printf("  Line %d - %s%s: %s\n", c.Line, c.Author, blocking, c.Body)
	}
	fmt.Println()

	fmt.Println(review.Summary())

	review.ResolveBlockingComments()
	review.Approve()
	fmt.Println("After resolving comments and approving:")
	fmt.Println(review.Summary())
}
