package main

import "fmt"

const maxDiffLines = 50
const maxComments = 30

var diffLines [maxDiffLines]string
var diffTypes [maxDiffLines]string
var diffCount int

var commentLines [maxComments]int
var commentAuthors [maxComments]string
var commentBodies [maxComments]string
var commentBlocking [maxComments]bool
var commentCount int

var changesRequested bool
var approved bool

func reviewDiff(diffContent string) int {
	diffCount = 0
	changesRequested = false
	approved = false
	commentCount = 0
	start := 0
	for i := 0; i < len(diffContent); i++ {
		if diffContent[i] == '\n' {
			line := diffContent[start:i]
			addDiffLine(line)
			start = i + 1
		}
	}
	if start < len(diffContent) {
		line := diffContent[start:]
		addDiffLine(line)
	}
	return diffCount
}

func addDiffLine(content string) {
	diffLines[diffCount] = content
	dt := "context"
	if len(content) > 0 && content[0] == '+' {
		dt = "addition"
	} else if len(content) > 0 && content[0] == '-' {
		dt = "deletion"
	}
	diffTypes[diffCount] = dt
	diffCount++
}

func addComment(line int, author, body string, blocking bool) {
	commentLines[commentCount] = line
	commentAuthors[commentCount] = author
	commentBodies[commentCount] = body
	commentBlocking[commentCount] = blocking
	if blocking {
		changesRequested = true
		approved = false
	}
	commentCount++
}

func resolveBlocking() {
	changesRequested = false
	for i := 0; i < commentCount; i++ {
		commentBlocking[i] = false
	}
}

func submitReview() bool {
	if changesRequested {
		return false
	}
	approved = true
	return true
}

func summary() string {
	additions := 0
	deletions := 0
	for i := 0; i < diffCount; i++ {
		if diffTypes[i] == "addition" {
			additions++
		} else if diffTypes[i] == "deletion" {
			deletions++
		}
	}
	blocking := 0
	for i := 0; i < commentCount; i++ {
		if commentBlocking[i] {
			blocking++
		}
	}
	status := "PENDING REVIEW"
	if approved {
		status = "APPROVED"
	} else if changesRequested {
		status = "CHANGES REQUESTED"
	}
	return fmt.Sprintf("Diff: %d lines (%d additions, %d deletions)\nComments: %d (%d blocking)\nStatus: %s\n",
		diffCount, additions, deletions, commentCount, blocking, status)
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

	reviewDiff(diff)

	fmt.Println("=== Code Review Simulation ===")
	fmt.Println()
	fmt.Println("Submitted Diff:")
	for i := 0; i < diffCount; i++ {
		mark := " "
		if diffTypes[i] == "addition" {
			mark = "+"
		} else if diffTypes[i] == "deletion" {
			mark = "-"
		}
		fmt.Printf("  %s %s\n", mark, diffLines[i])
	}
	fmt.Println()

	addComment(5, "alice", "Missing error handling on new() call", true)
	addComment(9, "bob", "Consider adding a unit test for this path", false)

	fmt.Println("Review Comments:")
	for i := 0; i < commentCount; i++ {
		blocking := ""
		if commentBlocking[i] {
			blocking = " [BLOCKING]"
		}
		fmt.Printf("  Line %d - %s%s: %s\n", commentLines[i], commentAuthors[i], blocking, commentBodies[i])
	}
	fmt.Println()

	fmt.Print(summary())

	resolveBlocking()
	submitReview()
	fmt.Println("After resolving comments and approving:")
	fmt.Print(summary())
}
