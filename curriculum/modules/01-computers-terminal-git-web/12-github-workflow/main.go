package main

import "fmt"

const (
	statusOpen             = 0
	statusChangesRequested = 1
	statusApproved         = 2
	statusMerged           = 3
	statusClosed           = 4
)

var prID int
var prTitle string
var prSource string
var prTarget string
var prStatus int
var prComments [30]string
var prCommentCount int

func prStatusString(s int) string {
	if s == statusOpen {
		return "open"
	}
	if s == statusChangesRequested {
		return "changes requested"
	}
	if s == statusApproved {
		return "approved"
	}
	if s == statusMerged {
		return "merged"
	}
	if s == statusClosed {
		return "closed"
	}
	return "unknown"
}

func createPR(title, desc, source, target string) int {
	prID = 1
	prTitle = title
	prSource = source
	prTarget = target
	prStatus = statusOpen
	prCommentCount = 0
	return prID
}

func addComment(author, body string) {
	prComments[prCommentCount] = author + ": " + body
	prCommentCount++
}

func requestChanges() {
	prStatus = statusChangesRequested
}

func approvePR() {
	prStatus = statusApproved
}

func mergePR() (bool, string) {
	if prStatus != statusApproved {
		return false, "PR must be approved before merging"
	}
	prStatus = statusMerged
	return true, "Pull request merged successfully"
}

func closePR() {
	if prStatus != statusMerged && prStatus != statusApproved {
		prStatus = statusClosed
	}
}

func main() {
	createPR("Add login feature", "Implements user authentication with JWT", "feature-login", "main")

	fmt.Println("=== GitHub Workflow Simulation ===")
	fmt.Println()
	fmt.Printf("PR #%d created: %s (%s -> %s)\n", prID, prTitle, prSource, prTarget)
	fmt.Printf("Status: %s\n", prStatusString(prStatus))
	fmt.Println()

	addComment("reviewer1", "Please add error handling for empty credentials")
	requestChanges()
	fmt.Printf("After review: %s\n", prStatusString(prStatus))

	addComment("author", "Added error handling as requested")
	approvePR()
	fmt.Printf("After approval: %s\n", prStatusString(prStatus))

	ok, msg := mergePR()
	if ok {
		fmt.Printf("Merge result: %s\n", msg)
		fmt.Printf("Final status: %s\n", prStatusString(prStatus))
	} else {
		fmt.Printf("Merge blocked: %s\n", msg)
	}
	fmt.Println()

	fmt.Println("Review comments:")
	for i := 0; i < prCommentCount; i++ {
		fmt.Println(" ", prComments[i])
	}
}
