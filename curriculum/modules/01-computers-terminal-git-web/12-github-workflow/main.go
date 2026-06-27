package main

import (
	"fmt"
)

type PRStatus int

const (
	PRDraft PRStatus = iota
	PROpen
	PRChangesRequested
	PRApproved
	PRMerged
	PRClosed
)

func (s PRStatus) String() string {
	switch s {
	case PRDraft:
		return "draft"
	case PROpen:
		return "open"
	case PRChangesRequested:
		return "changes requested"
	case PRApproved:
		return "approved"
	case PRMerged:
		return "merged"
	case PRClosed:
		return "closed"
	default:
		return "unknown"
	}
}

type PullRequest struct {
	ID          int
	Title       string
	Description string
	Source      string
	Target      string
	Status      PRStatus
	Comments    []string
	Approvals   int
}

type GitHubRepo struct {
	Name      string
	Branches  map[string]bool
	Protected bool
}

func StartPR(title, desc, source, target string) *PullRequest {
	return &PullRequest{
		ID:          1,
		Title:       title,
		Description: desc,
		Source:      source,
		Target:      target,
		Status:      PROpen,
		Comments:    []string{},
		Approvals:   0,
	}
}

func (pr *PullRequest) AddComment(author, body string) {
	pr.Comments = append(pr.Comments, fmt.Sprintf("%s: %s", author, body))
}

func (pr *PullRequest) RequestChanges() {
	pr.Status = PRChangesRequested
}

func (pr *PullRequest) Approve() {
	pr.Approvals++
	if pr.Approvals >= 1 {
		pr.Status = PRApproved
	}
}

func (pr *PullRequest) Merge(repo *GitHubRepo) (bool, string) {
	if pr.Status != PRApproved {
		return false, "PR must be approved before merging"
	}
	if repo.Protected && pr.Target == "main" {
		if len(pr.Comments) == 0 {
			return false, "protected branch requires at least one review comment"
		}
	}
	pr.Status = PRMerged
	return true, "Pull request merged successfully"
}

func (pr *PullRequest) Close() {
	if pr.Status != PRMerged && pr.Status != PRApproved {
		pr.Status = PRClosed
	}
}

func main() {
	fmt.Println("=== GitHub Workflow Simulation ===")
	fmt.Println()

	repo := &GitHubRepo{
		Name:      "my-project",
		Branches:  map[string]bool{"main": true, "feature-login": true},
		Protected: true,
	}

	pr := StartPR("Add login feature", "Implements user authentication with JWT", "feature-login", "main")
	fmt.Printf("PR #%d created: %s (%s -> %s)\n", pr.ID, pr.Title, pr.Source, pr.Target)
	fmt.Printf("Status: %s\n", pr.Status)
	fmt.Println()

	pr.AddComment("reviewer1", "Please add error handling for empty credentials")
	pr.RequestChanges()
	fmt.Printf("After review: %s\n", pr.Status)

	pr.AddComment("author", "Added error handling as requested")
	pr.Approve()
	fmt.Printf("After approval: %s\n", pr.Status)

	ok, msg := pr.Merge(repo)
	if ok {
		fmt.Printf("Merge result: %s\n", msg)
		fmt.Printf("Final status: %s\n", pr.Status)
	} else {
		fmt.Printf("Merge blocked: %s\n", msg)
	}
	fmt.Println()

	fmt.Println("Review comments:")
	for _, c := range pr.Comments {
		fmt.Println(" ", c)
	}
}
