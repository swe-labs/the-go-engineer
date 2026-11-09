package main

import (
	"fmt"
)

// DiffLine represents a single line change in a pull request diff.
type DiffLine struct {
	LineNumber int
	Type       string // "+", "-", " "
	Content    string
}

// ReviewComment represents inline feedback left by a reviewer.
type ReviewComment struct {
	LineNumber int
	Reviewer   string
	Comment    string
	IsBlocking bool
}

func main() {
	fmt.Println("=== Code Review Workflow Simulation ===")

	diff := []DiffLine{
		{LineNumber: 10, Type: " ", Content: "func ConnectDB() (*sql.DB, error) {"},
		{LineNumber: 11, Type: "-", Content: "    db, _ := sql.Open(\"postgres\", connStr)"},
		{LineNumber: 11, Type: "+", Content: "    db, err := sql.Open(\"postgres\", connStr)"},
		{LineNumber: 12, Type: "+", Content: "    if err != nil { return nil, err }"},
		{LineNumber: 13, Type: " ", Content: "    return db, nil"},
	}

	fmt.Println("--- Diff Under Review ---")
	for _, l := range diff {
		fmt.Printf("%3d | %s %s\n", l.LineNumber, l.Type, l.Content)
	}

	comments := []ReviewComment{
		{
			LineNumber: 11,
			Reviewer:   "senior-gopher",
			Comment:    "Nice catch replacing the blank identifier with explicit error checking!",
			IsBlocking: false,
		},
		{
			LineNumber: 13,
			Reviewer:   "senior-gopher",
			Comment:    "Consider adding db.Ping() before returning to verify the connection is alive.",
			IsBlocking: false,
		},
	}

	fmt.Println("\n--- Review Comments ---")
	hasBlockingComments := false
	for _, c := range comments {
		tag := "SUGGESTION"
		if c.IsBlocking {
			tag = "BLOCKING"
			hasBlockingComments = true
		}
		fmt.Printf("[%s] Line %d by @%s: %s\n", tag, c.LineNumber, c.Reviewer, c.Comment)
	}

	fmt.Println("\n--- Merge Decision ---")
	if !hasBlockingComments {
		fmt.Println("Decision: APPROVED. No blocking review comments remaining; ready to merge.")
	} else {
		fmt.Println("Decision: CHANGES REQUESTED. Address blocking comments before merging.")
	}
}
